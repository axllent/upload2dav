# Upload2dav

[![Go Report Card](https://goreportcard.com/badge/github.com/axllent/upload2dav)](https://goreportcard.com/report/github.com/axllent/upload2dav)

A simple utility written in Go to upload files to a webdav server such as Nextcloud.


## Usage options

```shell
Usage: upload2dav [options] <file(s)>

Options:
  -d, --dir string     Alternative upload directory
  -c, --conf string    Specify config file (default "~/.config/upload2dav.json")
  -w, --write-config   Write config
  -q, --quiet          Quiet (do not show upload progress)
  -v, --version        Show version
  -u, --update         Update to latest version
  -h, --help           Show help
```


## Installation

You can use of the pre-built binaries (see releases).

If you prefer to build it from source `go get github.com/axllent/upload2dav`
