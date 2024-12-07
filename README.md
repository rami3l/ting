# ting

`ting` is a simple `tcping` implementation in Zig.

<!--
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
-->

## Usage

```console
> ting --help
ting, yet another TCPing
    -h, --help
            Display this help and exit.
    -v, --version
            Output version information and exit.
    -i, --interval <f32>
            Interval between pings, in seconds (default: 1.0)
    -c, --count <u16>
            Number of tries (default: null)
    -p, --port <u16>
            Numeric TCP port (default: 80)
    -w, --timeout <f32>
            Maximum time to wait for a response, in seconds (default: 5.0)
    <str>
            Host to reach
```
