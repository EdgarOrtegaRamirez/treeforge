package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/EdgarOrtegaRamirez/treeforge/pkg/scanner"
)

func TestNewScanner(t *testing.T) {
	config := scanner.ScannerConfig{
		RootPath: ".",
		MaxDepth: 3,
	}
	s := scanner.NewScanner(config)
	if s == nil {
		t.Fatal("NewScanner returned nil")
	}
}

func TestScanDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.go"), []byte("package main"), 0644)

	subDir := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "file3.txt"), []byte("world"), 0644)

	os.WriteFile(filepath.Join(tmpDir, ".hidden"), []byte("secret"), 0644)

	config := scanner.ScannerConfig{
		RootPath:   tmpDir,
		MaxDepth:   -1,
		ShowHidden: false,
	}

	s := scanner.NewScanner(config)
	root, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if root == nil {
		t.Fatal("root is nil")
	}

	if !root.IsDir {
		t.Fatal("root should be a directory")
	}

	if len(root.Children) != 3 {
		t.Errorf("expected 3 children, got %d", len(root.Children))
	}

	stats := s.GetStats()
	if stats.TotalFiles != 3 {
		t.Errorf("expected 3 files, got %d", stats.TotalFiles)
	}
	if stats.TotalDirs != 2 {
		t.Errorf("expected 2 dirs, got %d", stats.TotalDirs)
	}
}

func TestScanWithHidden(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".hidden"), []byte("secret"), 0644)

	config := scanner.ScannerConfig{
		RootPath:   tmpDir,
		MaxDepth:   -1,
		ShowHidden: true,
	}

	s := scanner.NewScanner(config)
	root, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(root.Children) != 2 {
		t.Errorf("expected 2 children with hidden, got %d", len(root.Children))
	}
}

func TestScanMaxDepth(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello"), 0644)

	subDir := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("world"), 0644)

	config := scanner.ScannerConfig{
		RootPath: tmpDir,
		MaxDepth: 1,
	}

	s := scanner.NewScanner(config)
	root, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(root.Children) != 2 {
		t.Errorf("expected 2 children at depth 1, got %d", len(root.Children))
	}
}

func TestScanExcludePattern(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file.py"), []byte("print('hi')"), 0644)

	config := scanner.ScannerConfig{
		RootPath:       tmpDir,
		MaxDepth:       -1,
		ExcludePattern: "*.go",
	}

	s := scanner.NewScanner(config)
	root, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	for _, child := range root.Children {
		if filepath.Ext(child.Name) == ".go" {
			t.Errorf("should have excluded .go files, but found %s", child.Name)
		}
	}
}

func TestScanIncludePattern(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file.py"), []byte("print('hi')"), 0644)

	config := scanner.ScannerConfig{
		RootPath:       tmpDir,
		MaxDepth:       -1,
		IncludePattern: "*.go",
	}

	s := scanner.NewScanner(config)
	root, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(root.Children) != 1 {
		t.Errorf("expected 1 child with include pattern, got %d", len(root.Children))
	}
	if root.Children[0].Name != "file.go" {
		t.Errorf("expected file.go, got %s", root.Children[0].Name)
	}
}

func TestScanEmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	config := scanner.ScannerConfig{
		RootPath: tmpDir,
		MaxDepth: -1,
	}

	s := scanner.NewScanner(config)
	root, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(root.Children) != 0 {
		t.Errorf("expected 0 children for empty dir, got %d", len(root.Children))
	}
}

func TestScanNonExistentPath(t *testing.T) {
	config := scanner.ScannerConfig{
		RootPath: "/nonexistent/path/that/does/not/exist",
		MaxDepth: -1,
	}

	s := scanner.NewScanner(config)
	_, err := s.Scan()
	if err == nil {
		t.Fatal("expected error for non-existent path")
	}
}

func TestScanFileSizes(t *testing.T) {
	tmpDir := t.TempDir()
	content := []byte("hello world")
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), content, 0644)

	config := scanner.ScannerConfig{
		RootPath: tmpDir,
		MaxDepth: -1,
	}

	s := scanner.NewScanner(config)
	root, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(root.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(root.Children))
	}

	if root.Children[0].Size != 11 {
		t.Errorf("expected size 11, got %d", root.Children[0].Size)
	}
}

func TestScanStats(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("aaa"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.go"), []byte("bbbbb"), 0644)

	subDir := filepath.Join(tmpDir, "sub")
	os.MkdirAll(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "c.txt"), []byte("c"), 0644)

	config := scanner.ScannerConfig{
		RootPath: tmpDir,
		MaxDepth: -1,
	}

	s := scanner.NewScanner(config)
	_, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	stats := s.GetStats()

	if stats.TotalFiles != 3 {
		t.Errorf("expected 3 files, got %d", stats.TotalFiles)
	}
	if stats.TotalDirs != 2 {
		t.Errorf("expected 2 dirs, got %d", stats.TotalDirs)
	}
	if stats.TotalSize != 9 {
		t.Errorf("expected total size 9, got %d", stats.TotalSize)
	}
	if stats.LargestFile != "b.go" {
		t.Errorf("expected largest file b.go, got %s", stats.LargestFile)
	}
	if stats.LargestSize != 5 {
		t.Errorf("expected largest size 5, got %d", stats.LargestSize)
	}

	if stats.ByExtension[".txt"] != 2 {
		t.Errorf("expected 2 .txt files, got %d", stats.ByExtension[".txt"])
	}
	if stats.ByExtension[".go"] != 1 {
		t.Errorf("expected 1 .go file, got %d", stats.ByExtension[".go"])
	}
}

func TestScanMaxFiles(t *testing.T) {
	tmpDir := t.TempDir()
	for i := 0; i < 10; i++ {
		os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("x"), 0644)
	}

	config := scanner.ScannerConfig{
		RootPath: tmpDir,
		MaxDepth: -1,
		MaxFiles: 3,
	}

	s := scanner.NewScanner(config)
	root, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(root.Children) > 3 {
		t.Errorf("expected at most 3 children with MaxFiles=3, got %d", len(root.Children))
	}
}

func TestScanDirectoryTotalSize(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("aaa"), 0644) // 3 bytes
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("bb"), 0644)  // 2 bytes

	config := scanner.ScannerConfig{
		RootPath: tmpDir,
		MaxDepth: -1,
	}

	s := scanner.NewScanner(config)
	root, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if root.TotalSize != 5 {
		t.Errorf("expected root total size 5, got %d", root.TotalSize)
	}
}

func TestScanGitStatusConstants(t *testing.T) {
	if scanner.GitModified != "modified" {
		t.Errorf("expected GitModified to be 'modified', got %q", scanner.GitModified)
	}
	if scanner.GitAdded != "added" {
		t.Errorf("expected GitAdded to be 'added', got %q", scanner.GitAdded)
	}
	if scanner.GitUntracked != "untracked" {
		t.Errorf("expected GitUntracked to be 'untracked', got %q", scanner.GitUntracked)
	}
}

func TestScanDefaultMaxDepth(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "a", "b", "c"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "a", "b", "c", "deep.txt"), []byte("deep"), 0644)

	config := scanner.ScannerConfig{
		RootPath: tmpDir,
		MaxDepth: 0, // Should default to unlimited
	}

	s := scanner.NewScanner(config)
	root, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Should find the deep file since 0 defaults to unlimited
	found := false
	var walk func(node *scanner.FileNode)
	walk = func(node *scanner.FileNode) {
		if node.Name == "deep.txt" {
			found = true
		}
		for _, child := range node.Children {
			walk(child)
		}
	}
	walk(root)

	if !found {
		t.Error("expected to find deep.txt with unlimited depth")
	}
}
