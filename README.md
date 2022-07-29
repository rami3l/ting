# ting

`ting` is a simple `tcping` implementation, heavily inspired by [zhengxiaowai/tcping].

## Contents

- [ting](#ting)
  - [Contents](#contents)
  - [Installation](#installation)
  - [Usage](#usage)

## Installation

- `homebrew`/`linuxbrew` install:

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

## Usage

`ting [hosts...] [flags]`

- `-i, --interval float32`: Interval between pings, in seconds (default `1`)
- `-n, --count int`: Number of tries (default `5`)
- `-p, --port int`: Numeric TCP port (default `80`)
- `-w, --timeout float32`: Maximum time to wait for a response, in seconds (default `5`)

[zhengxiaowai/tcping]: https://github.com/zhengxiaowai/tcping
