# D2M

[![Go](https://github.com/vuon9/d2m/actions/workflows/go.yml/badge.svg)](https://github.com/vuon9/d2m/actions/workflows/go.yml)

Dota2 scheduled matches tracker

## Installation

Install from source with go (requires go 1.18+)

```bash
# install
go install github.com/vuon9/d2m@latest
```

Or download the binary from [release page](https://github.com/vuon9/d2m/releases)

## To use it

```bash
❯ d2m
```

Get help
```
❯ d2m help
NAME:
   d2m - Dota2 matches schedule on terminal

USAGE:
   d2m [global options] command [command options] [arguments...]

COMMANDS:
   live, l       Live matches
   coming, u     Upcoming matches
   finished, f   Finished matches
   today, t      Matches today
   tomorrow, m   Matches tomorrow
   yesterday, y  Matches yesterday
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```


## Screenshots

![Main](./screenshots/main.png)