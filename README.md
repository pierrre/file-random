# File Random

[![Go Reference](https://pkg.go.dev/badge/github.com/pierrre/file-random.svg)](https://pkg.go.dev/github.com/pierrre/file-random)

## Features

- Get random files (and open them)
- Command line ready to use
- Library that can be integrated in a project

## Usage

```bash
# Local build
make build
./build/file-random -h

# Remote install
go install github.com/pierrre/file-random/cmd/file-random@latest

# Module install
go get github.com/pierrre/file-random@latest
```

## Flags

| Flag | Default | Description |
| --- | --- | --- |
| `-min-size` | `1` | Minimum file size in bytes. Files smaller than this are excluded. |
| `-open` | `false` | Open the selected file with the default system application. |
| `-loop` | `false` | Continuously select a new random file. Waits for Enter between each selection. |
| `-continue-on-error` | `false` | Continue on errors instead of exiting. |
| `-v` | `false` | Verbose. Logs errors during file collection. Only effective with `-continue-on-error`. |

Root directories are passed as positional arguments.
