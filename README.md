# plexize
[![Go Report Card](https://goreportcard.com/badge/github.com/m4ns0ur/plexize)](https://goreportcard.com/report/github.com/m4ns0ur/plexize)

Movie file, [Plex](https://www.plex.tv/) friendly maker.

## Install
`$ GO111MODULE=on go get github.com/m4ns0ur/plexize`

## Run
`$ plexize -h`

Note that `$GOPATH/bin` should be in the path.

## Usage
```
$ plexize -h
Movie file, Plex friendly maker.

Usage:
  plexize [OPTION]... FILE...

Options:
  -d, --dry-run             Show result without running
  -m, --change-mode         Change file mode to 660
  -o, --change-owner        Change file owner to plex:plex (sudo might be needed)
```

## License
MIT - see [LICENSE][license]

[license]: https://github.com/m4ns0ur/covid/blob/master/LICENSE
