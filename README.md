# ting

`ting` is yet another `tcping` implementation.

## Contents

- [ting](#ting)
  - [Contents](#contents)
  - [Motivation & Current Status](#motivation--current-status)
  - [Build & Installation](#build--installation)
  - [Usage & Options](#usage--options)

## Motivation & Current Status

This project is heavily inspired by [zhengxiaowai/tcping], which is working pretty fine most of the times,
but gets potentially broken when the `Python` interpreter gets updated.  

Thus, using `Golang` enables me to solve the problem almost as efficiently, with the additional benefit of
being able to easily distribute the binaries.

## Build & Installation

- `homebrew` install:

  ```bash
  brew install rami3l/tap/ting
  ```

- Build from source:

  ```bash
  # To install:
  go install github.com/rami3l/ting

  # To uninstall:
  go clean -i "github.com/rami3l/ting"
  ```

## Usage & Options

Usage: `ting [hosts...] [flags]`

- `-i, --interval float32`: Interval between pings, in seconds (default `1`)
- `-n, --count int`: Number of tries (default `5`)
- `-p, --port int`: Numeric TCP port (default `80`)
- `-w, --timeout float32`: Maximum time to wait for a response, in seconds (default `5`)

[zhengxiaowai/tcping]: https://github.com/zhengxiaowai/tcping
