# AGENTS.md

## Project Overview
TreeForge is a modern directory tree visualization CLI tool in Go. It replaces the classic `tree` command with additional features like file sizes, permissions, git status integration, glob filtering, and multiple output formats.

## Architecture

### Core Modules

1. **pkg/scanner/** - Directory walking and tree building
   - `models.go` - Data structures (FileNode, TreeStats, ScannerConfig)
   - `scanner.go` - Recursive directory walker with glob filtering

2. **pkg/output/** - Output formatters
   - `writer.go` - Tree, JSON, compact, and list formatters

3. **pkg/git/** - Git integration
   - `git.go` - Git status detection via `git status --porcelain`

4. **cmd/treeforge/** - CLI entry point
   - `main.go` - Cobra-based CLI with subcommands

### CLI Commands
- `treeforge [path]` - Display directory tree
- `treeforge stats [path]` - Show directory statistics
- `treeforge version` - Print version

## Testing
- Run all tests: `go test ./tests/... -v`
- Run specific test: `go test ./tests/... -run TestScanDirectory -v`
- Test coverage: `go test ./tests/... -cover`

## Key Design Decisions
- Tree output uses Unicode box-drawing characters for clean visualization
- Git status is detected by running `git status --porcelain` and building a status map
- Glob patterns use Go's `filepath.Match` for inclusion/exclusion
- Statistics are computed during the scan pass for efficiency
- Color output uses ANSI escape codes (no external dependencies)

## Adding New Features
- New output formats: add to `pkg/output/writer.go`
- New scanner options: add to `pkg/scanner/models.go`
- New CLI commands: add to `cmd/treeforge/main.go`
