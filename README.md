# EVE CLI

A simple CLI for useful [eve-os](https://github.com/lf-edge/eve) actions.
Likely eventually to be merged into eve-os itself or [eden](https://github.com/lf-edge/eden).

It does several useful things. For now, all it does is generate an iPXE script either to
stdout or to a file. Defaults to pulling assets from github, amd64 and the latest available
version, but you can override all of the above.

## Usage

```
eve-cli help
```

That will tell you everything you need to know.

## Installation

Download the command for your operating system and architecture on the
[releases page](https://github.com/deitch/eve-cli/releases).

## Building

1. Clone this repository
1. Build it with [go](https://golang.org)


