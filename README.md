# plexize
[![Go Report Card](https://goreportcard.com/badge/github.com/m4ns0ur/plexize)](https://goreportcard.com/report/github.com/m4ns0ur/plexize) [![Travis CI Build](https://travis-ci.org/m4ns0ur/plexize.svg?branch=master)](https://travis-ci.com/github/m4ns0ur/plexize)

Movie and TV show files, Plex friendly maker.

[Plex](https://www.plex.tv/) is a media server where you can keep all your media (movies, TV show and more) in one centralized place, and access them from different devices. Also Plex is really good in fetching media information and metada, and it has a nice dashboard to show all these information.

Plex is doing its best effort to find media information based on media file name. If you follow the file name conventions and file structures ([movie](https://support.plex.tv/articles/naming-and-organizing-your-movie-media-files/)/[TV show](https://support.plex.tv/articles/naming-and-organizing-your-tv-show-files/)) it helps a lot to get best results. Plexize will help to convert downloaded movie/TV show file name in the Plex way.

Plex also has guidlines regarding to [Linux permissions for media files](https://support.plex.tv/articles/200288596-linux-permissions-guide/). Plexize will help you to set this up too.

## Install
`$ GO111MODULE=on go get github.com/m4ns0ur/plexize`

## Run
`$ plexize -h`

Note that `$GOPATH/bin` should be in the path.

## Usage
```
$ plexize -h
Movie and TV show files, Plex friendly maker.

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
  $ plexize -m -o The.Flash.2014.S01E01.HDTV.mkv   # change mode/owner a TV show file (would be separated in its own folder)
```

## License
MIT - see [LICENSE][license]

[license]: https://github.com/m4ns0ur/covid/blob/master/LICENSE
