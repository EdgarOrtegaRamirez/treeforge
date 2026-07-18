package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Scanner walks a directory tree and builds FileNode structures
type Scanner struct {
	config    ScannerConfig
	stats     TreeStats
	fileCount int
}

// NewScanner creates a new Scanner with the given config
func NewScanner(config ScannerConfig) *Scanner {
	if config.MaxDepth == 0 {
		config.MaxDepth = -1 // unlimited
	}
	return &Scanner{
		config: config,
		stats: TreeStats{
			ByExtension: make(map[string]int),
			ByGitStatus: make(map[GitStatus]int),
		},
	}
}

// Scan walks the directory tree and returns the root FileNode
func (s *Scanner) Scan() (*FileNode, error) {
	info, err := os.Stat(s.config.RootPath)
	if err != nil {
		return nil, err
	}

	root := &FileNode{
		Name:         filepath.Base(s.config.RootPath),
		Path:         s.config.RootPath,
		RelativePath: ".",
		IsDir:        info.IsDir(),
		Size:         info.Size(),
		Mode:         info.Mode(),
		ModTime:      info.ModTime(),
		Depth:        0,
	}

	if info.IsDir() {
		s.scanDir(root, 0)
	}

	s.computeStats(root)
	return root, nil
}

// GetStats returns the statistics from the last scan
func (s *Scanner) GetStats() TreeStats {
	return s.stats
}

func (s *Scanner) scanDir(node *FileNode, depth int) {
	if s.config.MaxDepth >= 0 && depth >= s.config.MaxDepth {
		return
	}
	if s.config.MaxFiles > 0 && s.fileCount >= s.config.MaxFiles {
		return
	}

	entries, err := os.ReadDir(node.Path)
	if err != nil {
		return
	}

	// Sort entries
	s.sortEntries(entries)

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files if configured
		if !s.config.ShowHidden && strings.HasPrefix(name, ".") {
			continue
		}

		// Check exclusion pattern
		if s.config.ExcludePattern != "" {
			matched, _ := filepath.Match(s.config.ExcludePattern, name)
			if matched {
				continue
			}
		}

		// Check inclusion pattern
		if s.config.IncludePattern != "" {
			matched, _ := filepath.Match(s.config.IncludePattern, name)
			if !matched {
				continue
			}
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		childPath := filepath.Join(node.Path, name)
		relPath := filepath.Join(node.RelativePath, name)

		child := &FileNode{
			Name:         name,
			Path:         childPath,
			RelativePath: relPath,
			IsDir:        entry.IsDir(),
			Size:         info.Size(),
			Mode:         info.Mode(),
			ModTime:      info.ModTime(),
			Depth:        depth + 1,
		}

		s.fileCount++

		if entry.IsDir() {
			s.scanDir(child, depth+1)
			node.NumChildren += child.NumChildren + 1
			node.TotalSize += child.TotalSize + child.Size
		} else {
			node.NumChildren++
			node.TotalSize += child.Size
		}

		node.Children = append(node.Children, child)
	}
}

func (s *Scanner) sortEntries(entries []os.DirEntry) {
	sort.Slice(entries, func(i, j int) bool {
		// Directories first
		iDir := entries[i].IsDir()
		jDir := entries[j].IsDir()
		if iDir != jDir {
			return iDir
		}

		nameI := strings.ToLower(entries[i].Name())
		nameJ := strings.ToLower(entries[j].Name())

		if s.config.SortReverse {
			return nameI > nameJ
		}
		return nameI < nameJ
	})
}

func (s *Scanner) computeStats(root *FileNode) {
	s.walkStats(root, ".")
}

func (s *Scanner) walkStats(node *FileNode, basePath string) {
	if node.IsDir {
		s.stats.TotalDirs++
		for _, child := range node.Children {
			s.walkStats(child, basePath+"/"+child.Name)
		}
	} else {
		s.stats.TotalFiles++
		s.stats.TotalSize += node.Size

		// Track by extension
		ext := strings.ToLower(filepath.Ext(node.Name))
		if ext == "" {
			ext = "(no ext)"
		}
		s.stats.ByExtension[ext]++

		// Track by git status
		s.stats.ByGitStatus[node.GitStatus]++

		// Track largest
		if node.Size > s.stats.LargestSize {
			s.stats.LargestSize = node.Size
			s.stats.LargestFile = node.RelativePath
		}

		// Track oldest/newest
		if s.stats.OldestTime.IsZero() || node.ModTime.Before(s.stats.OldestTime) {
			s.stats.OldestTime = node.ModTime
			s.stats.OldestFile = node.RelativePath
		}
		if s.stats.NewestTime.IsZero() || node.ModTime.After(s.stats.NewestTime) {
			s.stats.NewestTime = node.ModTime
			s.stats.NewestFile = node.RelativePath
		}
	}
}
