package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var uid int = -1
var gid int

func init() {
	u, _ := user.Lookup("plex")
	if u == nil {
		return
	}
	var err error
	uid, err = strconv.Atoi(u.Uid)
	if err != nil {
		log.Panicf("cannot parse the number: %v\n", err)
	}
	gid, err = strconv.Atoi(u.Gid)
	if err != nil {
		log.Panicf("cannot parse the number: %v\n", err)
	}
}

type patterns [10]*regexp.Regexp

func (p *patterns) match(s string) bool {
	for _, r := range p {
		if r.MatchString(s) {
			return true
		}
	}
	return false
}

// Inspired by PTN project (https://github.com/divijbindlish/parse-torrent-name/blob/master/PTN/patterns.py).
var commonPatterns = patterns{
	regexp.MustCompile(`(?:PPV\.)?[HP]DTV|(?:HD)?CAM|hd-?ts|(?:PPV )?WEB-?DL(?: DVDRip)?|(?:D[vV][dD]|H[dD]|Cam|W[EB]B|B[DR])(?:(?i)rip)|(?:(?i)blu[-]?ray)|(?:(?i)telesync)|DvDScr|hdtv|PPV`),
	regexp.MustCompile(`MP3|DD5\.?1|Dual[\- ]Audio|LiNE|DTS|AAC[.-]LC|AAC(?:\.?2\.0|2)?|AC3(?:\.5\.1)?|Dual|Audio`),
	regexp.MustCompile(`xvid|[hx]\.?26[45](?:(?i)-fov|-w4f)?`),
	regexp.MustCompile(`(?i)hindi|(?:rus|ita)(?:\.eng|$)|eng$`),
	regexp.MustCompile(`[1-9]\d+(?:\.\d+)?(?:(?i)gb|mb)`),
	regexp.MustCompile(`(?i)EXTENDED(:?.CUT)?`),
	regexp.MustCompile(`[1-9]\d{2,3}p`),
	regexp.MustCompile(`(?:Half-)?SBS`),
	regexp.MustCompile(`MKV|AVI|MP4`),
	regexp.MustCompile(`unknown_release_type|UpScaled|iNTERNAL|CONVERT|[hH]ard[sS]ub|READNFO|PROPER|REPACK|UNRATED|(?:(?i)rarbg)|(?:(?i)hevc)|AMZN|PDTV|1CD|WEB|NBY|R[0-9]|TS|HC|WS|3D`),
}

type movie struct {
	sep     string
	name    string
	year    string
	season  string
	episode string
	epiName string
}

type plexFile struct {
	dir  string
	name string
	ext  string
	mov  movie
}

func (p *plexFile) parse() {
	var (
		// First movie ever was in 1888, so lets check movie years since 1800.
		yearRe        = regexp.MustCompile(`((?:1[8-9]|[2-9]\d)\d{2})`)
		seasonRe      = regexp.MustCompile(`[sS]?(\d{1,2})[eExX](\d{1,2})`)
		domainRe      = regexp.MustCompile(`^[wW]{2,3}\.[^.]*\.[^.]{3,4}(.*)$`)
		bracePrefixRe = regexp.MustCompile(`^[\[\(🃏].*[\]\)🃏](.*)$`)
		prefixRe      = regexp.MustCompile(`^[^0-9a-zA-Z]*(.*)$`)
	)

	n := p.name

	n = bracePrefixRe.ReplaceAllString(n, "$1")
	n = domainRe.ReplaceAllString(n, "$1")
	n = prefixRe.ReplaceAllString(n, "$1")
	n = strings.ReplaceAll(n, " - ", " ")

	max := 0
	for _, s := range [...]string{" ", ".", "-", "_"} {
		c := strings.Count(n, s)
		if c > max {
			max = c
			p.mov.sep = s
			continue
		}
	}

	if p.mov.sep == "" {
		p.mov.name = strings.Title(n)
		return
	}

	ts := strings.Split(n, p.mov.sep)
	done := false
	seasoned := false
	for _, t := range ts {
		t = strings.Trim(t, " -[]()")
		if !done {
			if y := yearRe.FindString(t); y != "" && p.mov.name != "" {
				done = true
				p.mov.year = y
				continue
			}

			if s := seasonRe.FindStringSubmatch(t); len(s) != 0 {
				done = true
				seasoned = true
				p.mov.season = fmt.Sprintf("%02v", s[1])
				p.mov.episode = fmt.Sprintf("%02v", s[2])
				continue
			}

			if commonPatterns.match(t) {
				done = true
				continue
			}

			p.mov.name += t + " "
			continue
		}

		if y := yearRe.FindString(t); y != "" {
			if p.mov.year == "" {
				p.mov.year = y
			} else {
				p.mov.name += p.mov.year + " "
				p.mov.year = y
			}
			continue
		}

		if s := seasonRe.FindStringSubmatch(t); len(s) != 0 {
			seasoned = true
			p.mov.season = fmt.Sprintf("%02v", s[1])
			p.mov.episode = fmt.Sprintf("%02v", s[2])
			continue
		}

		if commonPatterns.match(t) {
			break
		}

		if seasoned {
			p.mov.epiName += t + " "
		} else {
			p.mov.name += t + " "
		}
	}

	p.mov.name = strings.TrimSpace(p.mov.name)
	p.mov.name = strings.Title(p.mov.name)

	if seasoned {
		p.mov.epiName = strings.TrimSpace(p.mov.epiName)
		p.mov.epiName = strings.Title(p.mov.epiName)
	}
}

func (p *plexFile) plexName() string {
	if p.mov.name == "" {
		return ""
	}

	if p.mov.season == "" {
		if p.mov.year == "" {
			return p.mov.name
		}
		return fmt.Sprintf("%s (%s)", p.mov.name, p.mov.year)
	}

	if p.mov.epiName == "" {
		if p.mov.year == "" {
			return fmt.Sprintf("%s - s%se%s", p.mov.name, p.mov.season, p.mov.episode)
		}
		return fmt.Sprintf("%s (%s) - s%se%s", p.mov.name, p.mov.year, p.mov.season, p.mov.episode)
	}

	if p.mov.year == "" {
		return fmt.Sprintf("%s - s%se%s - %s", p.mov.name, p.mov.season, p.mov.episode, p.mov.epiName)
	}

	return fmt.Sprintf("%s (%s) - s%se%s - %s", p.mov.name, p.mov.year, p.mov.season, p.mov.episode, p.mov.epiName)
}

func (p *plexFile) plexDir() string {
	if p.mov.name == "" {
		return ""
	}

	if p.mov.year == "" {
		return p.mov.name
	}

	return fmt.Sprintf("%s (%s)", p.mov.name, p.mov.year)
}

func (p *plexFile) seasonDir() string {
	if p.mov.season == "" {
		return ""
	}

	return fmt.Sprintf("Season %s", p.mov.season)
}

func usage() {
	fmt.Fprintln(flag.CommandLine.Output(), `Movie and TV show files, Plex friendly maker.

Usage:
  plexize [-]
  plexize [OPTION]... FILE...

Options:
  -d, --dry-run             Show result without running
  -m, --change-mode         Change file mode to 660
  -o, --change-owner        Change file owner to plex:plex (sudo might be needed)
  -p, --path PATH           Output path (move file to the path and then refactor)
  -s, --separate            Separate movie files in their own folders (not required for TV series)

Example:
  $ plexize                                        # start in interactive mode to convert file(s) name
  $ cat movie_list.txt | plexize                   # convert file(s) name with piping
  $ plexize trainwreck.mkv war.dogs.2016.mkv       # convert multiple files
  $ plexize the*.mkv                               # convert multiple files with wildcard
  $ plexize -d The.Platform.2019.720p.mkv          # dry run
  $ plexize -p ~/plex The.Platform.2019.720p.mkv   # move the file to ~/plex and convert
  $ plexize -m -o -s The.Platform.2019.720p.mkv    # change mode/owner and move the movie file to its own folder
  $ plexize -m -o The.Flash.2014.S01E01.HDTV.mkv   # change mode/owner a TV show file (would be separated in its own folder)`)
}

func main() {
	log.SetFlags(0)

	var (
		dryRun, chmod, chown, separate bool
		outDir                         string
	)

	flag.Usage = usage
	flag.BoolVar(&dryRun, "d", false, "Show result without running")
	flag.BoolVar(&dryRun, "dry-run", false, "Show result without running")
	flag.BoolVar(&chmod, "m", false, "Change file mode to 660")
	flag.BoolVar(&chmod, "change-mode", false, "Change file mode to 660")
	flag.BoolVar(&chown, "o", false, "Change file owner (default is plex:plex)")
	flag.BoolVar(&chown, "change-owner", false, "Change file owner (default is plex:plex)")
	flag.StringVar(&outDir, "p", "", "Output path (move file to the path and then refactor)")
	flag.StringVar(&outDir, "path", "", "Output path (move file to the path and then refactor)")
	flag.BoolVar(&separate, "s", false, "Separate movie files in their own folders (not required for TV series)")
	flag.BoolVar(&separate, "separate", false, "Separate movie files in their own folders (not required for TV series)")
	flag.Parse()

	if flag.Arg(0) == "" || flag.Arg(0) == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			l := scanner.Text()
			log.Printf("%s\n", convert(l, true, false, ""))
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("cannot read from stdin: %v\n", err)
		}
		os.Exit(0)
	}

	if dryRun {
		log.Println("Dry run...")
	}

	canChown := true
	if chown {
		if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
			canChown = false
			log.Println("the OS does not support changing the file owner")
		} else if uid == -1 {
			canChown = false
			log.Println("user plex does not exist, cannot change the file owner")
		}
	}

	for i := 0; i < flag.NArg(); i++ {
		var paths []string
		var err error
		if strings.Contains(flag.Arg(i), "*") {
			paths, err = filepath.Glob(flag.Arg(i))
			if err != nil {
				log.Fatalf("invalid path format: %v\n", err)
			}
		} else {
			paths = append(paths, flag.Arg(i))
		}
		for _, path := range paths {
			np := convert(path, dryRun, separate, outDir)
			log.Printf("%s -> %s\n", path, np)

			if !dryRun {
				err := os.Rename(path, np)
				if err != nil {
					if os.IsPermission(err) {
						log.Printf("you don't have permission to move/rename the file (you can retry with sudo): %v\n", err)
					} else {
						log.Printf("cannot move/rename the file: %v\n", err)
					}
				}

				if chmod {
					err := os.Chmod(np, 0660)
					if err != nil {
						log.Printf("cannot change the file mode: %v\n", err)
					}
				}

				if chown && canChown {
					err = os.Chown(np, uid, gid)
					if os.IsPermission(err) {
						log.Printf("you don't have permission to change owner of the file (you can retry with sudo): %v\n", err)
					}
				}

				// TODO: support copy to server (delete local).
				// TODO: support fixing title in metadata.
			}
		}
	}
}

func convert(path string, dryRun, separate bool, outDir string) (newPath string) {
	dir, file := filepath.Split(path)
	ext := filepath.Ext(file)

	pf := &plexFile{
		dir:  dir,
		name: strings.TrimSuffix(file, ext),
		ext:  strings.ToLower(ext),
		mov:  movie{},
	}
	pf.parse()

	ps := make([]string, 0, 4)
	ps = append(ps, dir)
	if outDir != "" {
		ps[0] = outDir
	}
	if separate || pf.mov.season != "" {
		ps = append(ps, pf.plexDir())
		if !dryRun {
			err := os.MkdirAll(filepath.Join(ps...), os.ModePerm)
			if err != nil && !os.IsExist(err) {
				log.Printf("cannot make separate movie or TV serie folder: %v\n", err)
			}
		}
	}
	if pf.mov.season != "" {
		ps = append(ps, pf.seasonDir())
		if !dryRun {
			err := os.MkdirAll(filepath.Join(ps...), os.ModePerm)
			if err != nil && !os.IsExist(err) {
				log.Printf("cannot make TV serie season folder: %v\n", err)
			}
		}
	}
	ps = append(ps, pf.plexName())
	return fmt.Sprintf("%s%s", filepath.Join(ps...), pf.ext)
}
