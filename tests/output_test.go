package tests

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/treeforge/pkg/output"
	"github.com/EdgarOrtegaRamirez/treeforge/pkg/scanner"
)

func TestWriteTree(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello"), 0644)

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

	var buf bytes.Buffer
	w := output.NewWriter(&buf)
	err = w.WriteTree(root)
	if err != nil {
		t.Fatalf("WriteTree failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "file.txt") {
		t.Errorf("expected output to contain file.txt, got:\n%s", output)
	}
}

func TestWriteTreeWithSizes(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello world"), 0644)

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

	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.WithSizes(true))
	err = w.WriteTree(root)
	if err != nil {
		t.Fatalf("WriteTree failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[11 B]") {
		t.Errorf("expected output to contain [11 B], got:\n%s", output)
	}
}

func TestWriteCompact(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello"), 0644)
	os.MkdirAll(filepath.Join(tmpDir, "sub"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "sub", "nested.txt"), []byte("world"), 0644)

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

	var buf bytes.Buffer
	w := output.NewWriter(&buf)
	err = w.WriteCompact(root)
	if err != nil {
		t.Fatalf("WriteCompact failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "files") {
		t.Errorf("expected output to contain 'files', got:\n%s", output)
	}
	if !strings.Contains(output, "dirs") {
		t.Errorf("expected output to contain 'dirs', got:\n%s", output)
	}
}

func TestWriteJSON(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello"), 0644)

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

	var buf bytes.Buffer
	w := output.NewWriter(&buf)
	err = w.WriteJSON(root)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"name"`) {
		t.Errorf("expected JSON output to contain 'name', got:\n%s", output)
	}
	if !strings.Contains(output, `"type"`) {
		t.Errorf("expected JSON output to contain 'type', got:\n%s", output)
	}
}

func TestWriteList(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello"), 0644)
	os.MkdirAll(filepath.Join(tmpDir, "sub"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "sub", "nested.txt"), []byte("world"), 0644)

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

	var buf bytes.Buffer
	w := output.NewWriter(&buf)
	err = w.WriteList(root)
	if err != nil {
		t.Fatalf("WriteList failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "file.txt") {
		t.Errorf("expected list output to contain file.txt, got:\n%s", output)
	}
	if !strings.Contains(output, "sub/") {
		t.Errorf("expected list output to contain sub/, got:\n%s", output)
	}
}

func TestWriteListWithSizes(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello"), 0644)

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

	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.WithSizes(true))
	err = w.WriteList(root)
	if err != nil {
		t.Fatalf("WriteList failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "5 B") {
		t.Errorf("expected list output to contain size, got:\n%s", output)
	}
}

func TestFormatSize(t *testing.T) {
	// Test through compact output
	tmpDir := t.TempDir()

	// Create a 1KB file
	kbContent := make([]byte, 1024)
	os.WriteFile(filepath.Join(tmpDir, "1kb.txt"), kbContent, 0644)

	// Create a 1MB file
	mbContent := make([]byte, 1024*1024)
	os.WriteFile(filepath.Join(tmpDir, "1mb.txt"), mbContent, 0644)

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

	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.WithSizes(true))
	err = w.WriteCompact(root)
	if err != nil {
		t.Fatalf("WriteCompact failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "MB") {
		t.Errorf("expected output to contain MB, got:\n%s", output)
	}
}

func TestWriteEmptyTree(t *testing.T) {
	tmpDir := t.TempDir()

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

	var buf bytes.Buffer
	w := output.NewWriter(&buf)
	err = w.WriteTree(root)
	if err != nil {
		t.Fatalf("WriteTree failed: %v", err)
	}

	// Should only contain the root directory name
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line for empty tree, got %d: %s", len(lines), output)
	}
}

func TestWriteJSONEmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

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

	var buf bytes.Buffer
	w := output.NewWriter(&buf)
	err = w.WriteJSON(root)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"items":`) || !strings.Contains(output, `]`) {
		t.Errorf("expected items array in output, got:\n%s", output)
	}
}

func TestWriteTreeMultipleFormats(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.py"), []byte("print('hi')"), 0644)

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

	formats := []string{"tree", "json", "compact", "list"}
	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			var buf bytes.Buffer
			w := output.NewWriter(&buf)

			switch format {
			case "tree":
				err = w.WriteTree(root)
			case "json":
				err = w.WriteJSON(root)
			case "compact":
				err = w.WriteCompact(root)
			case "list":
				err = w.WriteList(root)
			}

			if err != nil {
				t.Fatalf("Write%s failed: %v", format, err)
			}

			if buf.Len() == 0 {
				t.Errorf("Write%s produced empty output", format)
			}
		})
	}
}

func TestWriterWithGitStatus(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello"), 0644)

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

	// Simulate git status
	for _, child := range root.Children {
		child.GitStatus = scanner.GitModified
	}

	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.WithGitStatus(true))
	err = w.WriteTree(root)
	if err != nil {
		t.Fatalf("WriteTree failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "<modified>") {
		t.Errorf("expected output to contain <modified>, got:\n%s", output)
	}
}
