# TreeForge 🔨

A modern, feature-rich directory tree visualization CLI tool in Go. A replacement for the classic `tree` command with additional features like file sizes, permissions, git status integration, glob filtering, and multiple output formats.

[![CI](https://github.com/EdgarOrtegaRamirez/treeforge/actions/workflows/ci.yml/badge.svg)](https://github.com/EdgarOrtegaRamirez/treeforge/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/EdgarOrtegaRamirez/treeforge)](https://goreportcard.com/report/github.com/EdgarOrtegaRamirez/treeforge)

## Features

- 📁 **Directory tree visualization** with Unicode box-drawing characters
- 📊 **File sizes** with human-readable formatting (B, KB, MB, GB)
- 🔐 **File permissions** display
- ⏰ **Modification times** 
- 🎨 **Color-coded output** by file type and git status
- 🌿 **Git status integration** — see modified, added, untracked files at a glance
- 🔍 **Glob pattern filtering** — include/exclude files by pattern
- 📏 **Depth limiting** — control how deep the tree goes
- 📄 **Multiple output formats** — tree, JSON, compact, list
- 📈 **Statistics** — file counts, sizes, extensions breakdown
- 🚫 **Hidden file control** — show or hide dotfiles

## Installation

```bash
go install github.com/EdgarOrtegaRamirez/treeforge/cmd/treeforge@latest
```

Or build from source:

```bash
git clone https://github.com/EdgarOrtegaRamirez/treeforge.git
cd treeforge
go build -o treeforge ./cmd/treeforge
```

## Quick Start

```bash
# Basic tree view
treeforge

# With file sizes
treeforge -s

# With sizes and permissions
treeforge -s -p

# Show git status
treeforge -g

# Limit depth
treeforge -d 2

# Show hidden files
treeforge -a

# Exclude patterns
treeforge -e "*.log"

# JSON output
treeforge -f json

# Directory statistics
treeforge stats
```

## Usage

```
Usage:
  treeforge [path] [flags]

Flags:
  -d, --depth int        Max depth (0 = root only, -1 = unlimited) (default -1)
  -e, --exclude string   Exclude pattern (glob)
  -f, --format string    Output format: tree, json, compact, list (default "tree")
  -g, --git              Show git status colors
  -h, --help             help for treeforge
  -a, --hidden           Show hidden files (starting with .)
  -i, --include string   Include pattern (glob)
  -p, --perms            Show file permissions
      --reverse          Reverse sort order
  -s, --sizes            Show file sizes
      --sort string      Sort by: name, size, time (default "name")
  -t, --time             Show modification times
```

## Output Formats

### Tree (default)
```
myproject
├── cmd/
│   └── treeforge/
├── pkg/
│   ├── git/
│   ├── output/
│   └── scanner/
├── go.mod [211 B]
├── go.sum [900 B]
└── treeforge [3.6 MB]
```

### JSON
```json
{
  "name": "myproject",
  "type": "directory",
  "children": 3,
  "totalSize": 3813592,
  "items": [...]
}
```

### Compact
```
4 files, 8 dirs, 3.6 MB total
```

### List
```
cmd/
pkg/
go.mod (211 B)
go.sum (900 B)
treeforge (3.6 MB)
```

## Git Status Colors

When using `-g` flag, files are color-coded by git status:
- 🟢 **Green** — Added files
- 🟡 **Yellow** — Modified files
- 🔴 **Red** — Deleted files
- 🔵 **Blue** — Directories
- 🩵 **Cyan** — Untracked files

## File Type Colors

- **Green** — Source code (`.go`, `.py`, `.js`, `.ts`, `.rs`, `.java`, `.c`, `.cpp`, `.h`)
- **Cyan** — Documentation (`.md`, `.txt`, `.rst`)
- **Magenta** — Config files (`.json`, `.yaml`, `.yml`, `.toml`, `.xml`)
- **Yellow** — Images (`.png`, `.jpg`, `.jpeg`, `.gif`, `.svg`)

## Examples

### Show a project with sizes and git status
```bash
treeforge -s -g
```

### Filter for Go files only
```bash
treeforge -i "*.go" -s
```

### Get project statistics
```bash
treeforge stats .
```

### Export as JSON for processing
```bash
treeforge -f json | jq '.items[] | select(.type == "file") | .name'
```

### Limit depth and exclude vendor
```bash
treeforge -d 3 -e "vendor"
```

## Architecture

```
treeforge/
├── cmd/treeforge/     # CLI entry point with cobra commands
├── pkg/
│   ├── scanner/       # Directory walking and tree building
│   ├── output/        # Output formatters (tree, JSON, compact, list)
│   └── git/           # Git status integration
└── tests/             # Test suite
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.
