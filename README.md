# D2M

[![GitHub release](https://img.shields.io/github/release/vuon9/d2m.svg)](https://GitHub.com/vuon9/d2m/releases/)
[![GitHub license](https://badgen.net/github/license/vuon9/d2m)](https://github.com/vuon9/d2m/blob/master/LICENSE)
[![Go](https://github.com/vuon9/d2m/actions/workflows/go.yml/badge.svg)](https://github.com/vuon9/d2m/actions/workflows/go.yml)
[![GitHub commits](https://badgen.net/github/commits/vuon9/d2m)](https://github.com/vuon9/d2m/commit/)

Dota2 scheduled matches tracker

![Main](./screenshots/main.gif)

## Installation

Install from source with go (requires go 1.18+)

```bash
go install github.com/vuon9/d2m@latest
```

Or download the binary from [release page](https://github.com/vuon9/d2m/releases)

## Usage

```bash
❯ d2m
```

## Key features
- Filter matches by time, status:
    - Upcoming
    - Live
    - Finished
    - Tomorrow
    - Today
    - Yesterday
    - From today
- Open Twitch streaming link in browser

## Tasks

- Async handle for the details, to ensure can handle async call well
    - make a call to a function which has sleep 3, return text when end
    - handle loading during the sleeping time
    - show the returned text
