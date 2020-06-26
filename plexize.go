package main

import (
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
		log.Panicf("Error: cannot parse the number: %s.\n", u.Uid)
	}
	gid, err = strconv.Atoi(u.Gid)
	if err != nil {
		log.Panicf("Error: cannot parse the number: %s.\n", u.Gid)
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
	regexp.MustCompile(`(?:PPV\.)?[HP]DTV|(?:HD)?CAM|hd-?ts|(?:PPV )?WEB-?DL(?: DVDRip)?|(?:DVD|H[dD]|Cam|W[EB]B|B[DR])(?:(?i)rip)|(?:(?i)blu[-]?ray)|(?:(?i)telesync)|DvDScr|hdtv|PPV`),
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
	sep    string
	name   string
	year   string
	season string
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
		seasonRe      = regexp.MustCompile(`[sS]?(\d{1,2})[eExX]\d{1,2}`)
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
	for _, s := range [4]string{" ", ".", "-", "_"} {
		c := strings.Count(n, s)
		if c > max {
			max = c
			p.mov.sep = s
			continue
		}
	}

	if p.mov.sep == "" {
		p.mov.name = n
		return
	}

	ts := strings.Split(n, p.mov.sep)
	done := false
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
				p.mov.name += s[0] + " "
				p.mov.season = fmt.Sprintf("%02v", s[1])
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
			p.mov.name += s[0] + " "
			p.mov.season = fmt.Sprintf("%02v", s[1])
			continue
		}

		if commonPatterns.match(t) {
			break
		}

		p.mov.name += t + " "
	}

	p.mov.name = strings.TrimSpace(p.mov.name)
	p.mov.name = strings.Title(p.mov.name)
}

func (p *plexFile) plexName() string {
	if p.mov.name == "" {
		return ""
	}
	if p.mov.year == "" {
		return p.mov.name
	}
	return fmt.Sprintf("%s (%s)", p.mov.name, p.mov.year)
}

func usage() {
	fmt.Fprintln(flag.CommandLine.Output(), `Movies file, Plex friendly maker.

Usage:
  plexize [OPTION]... [FILE]...

Options:
  -d, --dry-run             Show result without running
  -m, --change-mode         Change file mode to 660
  -o, --change-owner        Change file owner to plex:plex (sudo might be needed)`)
}

func main() {
	log.SetFlags(0)

	var (
		dryRun bool
		chmod  bool
		chown  bool
	)

	flag.Usage = usage
	flag.BoolVar(&dryRun, "d", false, "Show result without running")
	flag.BoolVar(&dryRun, "dry-run", false, "Show result without running")
	flag.BoolVar(&chmod, "m", false, "Change file mode to 660")
	flag.BoolVar(&chmod, "change-mode", false, "Change file mode to 660")
	flag.BoolVar(&chown, "o", false, "Change file owner (default is plex:plex)")
	flag.BoolVar(&chown, "change-owner", false, "Change file owner (default is plex:plex)")
	flag.Parse()

	if flag.NArg() < 1 {
		// TODO: support stdin.
		log.Fatalln("Error: file (path) is required")
	}

	if dryRun {
		log.Println("Dry run...")
	}

	canChown := true
	if chown {
		if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
			log.Printf("Error: OS does not support changing the file owner.\n")
		} else if uid == -1 {
			log.Printf("Error: user plex does not exist. Cannot change the file owner.\n")
		}
		canChown = false
	}

	// TODO: support recursive path walkthrough.
	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		dir, file := filepath.Split(path)
		ext := filepath.Ext(file)
		name := strings.TrimSuffix(file, ext)

		pf := &plexFile{mov: movie{}}
		pf.dir = dir
		pf.name = name
		pf.ext = strings.ToLower(ext)
		pf.parse()

		np := fmt.Sprintf("%s%s", filepath.Join(dir, pf.plexName()), pf.ext)
		log.Printf("%s -> %s\n", path, np)

		if !dryRun {
			err := os.Rename(path, np)
			if err != nil {
				log.Printf("Error: cannot rename the file.\n")
			}

			if chmod {
				err := os.Chmod(np, 0660)
				if err != nil {
					log.Printf("Error: cannot change the file mode.\n")
				}
			}

			if chown && canChown {
				err = os.Chown(np, uid, gid)
				if os.IsPermission(err) {
					log.Println("Error: you don't have permission to change owner of the file (you can retry with sudo).")
				}
			}

			// TODO: support copy to server (delete local).
		}
	}
}
