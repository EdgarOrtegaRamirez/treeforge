package scanner

import (
	"os"
	"time"
)

// GitStatus represents the git status of a file
type GitStatus string

const (
	GitModified   GitStatus = "modified"
	GitAdded      GitStatus = "added"
	GitDeleted    GitStatus = "deleted"
	GitRenamed    GitStatus = "renamed"
	GitCopied     GitStatus = "copied"
	GitUntracked  GitStatus = "untracked"
	GitIgnored    GitStatus = "ignored"
	GitUnmodified GitStatus = ""
)

// FileNode represents a file or directory in the tree
type FileNode struct {
	Name         string
	Path         string
	RelativePath string
	IsDir        bool
	Size         int64
	Mode         os.FileMode
	ModTime      time.Time
	Children     []*FileNode
	GitStatus    GitStatus
	Depth        int
	NumChildren  int   // for directories: total file count
	TotalSize    int64 // for directories: total size
}

// TreeStats holds summary statistics
type TreeStats struct {
	TotalFiles  int
	TotalDirs   int
	TotalSize   int64
	LargestFile string
	LargestSize int64
	OldestFile  string
	OldestTime  time.Time
	NewestFile  string
	NewestTime  time.Time
	ByExtension map[string]int
	ByGitStatus map[GitStatus]int
}

// ScannerConfig holds configuration for the scanner
type ScannerConfig struct {
	RootPath        string
	MaxDepth        int
	ShowHidden      bool
	ShowSizes       bool
	ShowPermissions bool
	ShowModTime     bool
	ShowGitStatus   bool
	IncludePattern  string // glob pattern for inclusion
	ExcludePattern  string // glob pattern for exclusion
	SortBy          string // "name", "size", "time"
	SortReverse     bool
	FollowSymlinks  bool
	MaxFiles        int // 0 = unlimited
}
