# D2M

[![GitHub release](https://img.shields.io/github/release/vuon9/d2m.svg)](https://GitHub.com/vuon9/d2m/releases/)
[![GitHub license](https://badgen.net/github/license/vuon9/d2m)](https://github.com/vuon9/d2m/blob/master/LICENSE)
[![Go](https://github.com/vuon9/d2m/actions/workflows/go.yml/badge.svg)](https://github.com/vuon9/d2m/actions/workflows/go.yml)
[![GitHub commits](https://badgen.net/github/commits/vuon9/d2m)](https://github.com/vuon9/d2m/commit/)

Dota2 scheduled matches tracker

![Main](./screenshots/main-with-details.gif)

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

## Features
- View details of teams
- Type `o` to open Twitch streaming link in web browser
- Type `?` to see all available filter commands
- Type `/` to quickly filter with [Regular expression](https://en.wikipedia.org/wiki/Regular_expression) then your regex, e.g. `team1|team2`
- Display icons to help you quickly identify the status:
    - Team has info: ◆ (but not 100% sure if it has roster, some teams have no roster or no info at all)
    - Team has no info (e.g TBD): ◇
    - Live match has streaming page: Twitch icon (has to install Nerdfont to display correctly)
