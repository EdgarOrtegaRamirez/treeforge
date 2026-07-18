package output

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/EdgarOrtegaRamirez/treeforge/pkg/scanner"
)

// Format constants
const (
	FormatTree    = "tree"
	FormatJSON    = "json"
	FormatCompact = "compact"
	FormatList    = "list"
)

// Writer writes the tree output
type Writer struct {
	writer        io.Writer
	showSizes     bool
	showPerms     bool
	showModTime   bool
	showGitStatus bool
	totalSize     int64
	totalFiles    int
	totalDirs     int
}

// NewWriter creates a new output writer
func NewWriter(w io.Writer, opts ...Option) *Writer {
	ww := &Writer{writer: w}
	for _, opt := range opts {
		opt(ww)
	}
	return ww
}

// Option configures the writer
type Option func(*Writer)

func WithSizes(show bool) Option   { return func(w *Writer) { w.showSizes = show } }
func WithPerms(show bool) Option   { return func(w *Writer) { w.showPerms = show } }
func WithModTime(show bool) Option { return func(w *Writer) { w.showModTime = show } }
func WithGitStatus(show bool) Option {
	return func(w *Writer) { w.showGitStatus = show }
}

// WriteTree outputs the tree in tree format
func (w *Writer) WriteTree(node *scanner.FileNode) error {
	w.totalFiles = 0
	w.totalDirs = 0
	w.totalSize = 0

	fmt.Fprintln(w.writer, node.Name)
	w.writeNode(node, "", true)

	return nil
}

func (w *Writer) writeNode(node *scanner.FileNode, prefix string, isLast bool) {
	if !node.IsDir {
		return
	}

	children := node.Children
	for i, child := range children {
		isLastChild := i == len(children)-1

		// Build prefix
		connector := "├── "
		if isLastChild {
			connector = "└── "
		}

		// Build name
		name := child.Name
		if child.IsDir {
			name += "/"
		}

		// Build metadata
		var meta []string
		if w.showSizes {
			if child.IsDir {
				meta = append(meta, fmt.Sprintf("[%s]", formatSize(child.TotalSize)))
			} else {
				meta = append(meta, fmt.Sprintf("[%s]", formatSize(child.Size)))
			}
		}
		if w.showPerms {
			meta = append(meta, fmt.Sprintf("(%s)", child.Mode.Perm().String()))
		}
		if w.showModTime {
			meta = append(meta, fmt.Sprintf("{%s}", child.ModTime.Format("2006-01-02 15:04")))
		}
		if w.showGitStatus && child.GitStatus != "" {
			meta = append(meta, fmt.Sprintf("<%s>", child.GitStatus))
		}

		// Color the name
		coloredName := colorize(name, child)

		// Print line
		metaStr := ""
		if len(meta) > 0 {
			metaStr = " " + strings.Join(meta, " ")
		}
		fmt.Fprintf(w.writer, "%s%s%s%s\n", prefix, connector, coloredName, metaStr)

		// Update stats
		if child.IsDir {
			w.totalDirs++
		} else {
			w.totalFiles++
			w.totalSize += child.Size
		}

		// Recurse for directories
		if child.IsDir {
			newPrefix := prefix
			if isLastChild {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			w.writeNode(child, newPrefix, isLastChild)
		}
	}
}

// WriteCompact outputs a compact single-line summary
func (w *Writer) WriteCompact(node *scanner.FileNode) error {
	w.walkCount(node)
	fmt.Fprintf(w.writer, "%d files, %d dirs, %s total\n",
		w.totalFiles, w.totalDirs, formatSize(w.totalSize))
	return nil
}

func (w *Writer) walkCount(node *scanner.FileNode) {
	if node.IsDir {
		w.totalDirs++
		for _, child := range node.Children {
			w.walkCount(child)
		}
	} else {
		w.totalFiles++
		w.totalSize += node.Size
	}
}

// WriteJSON outputs the tree as JSON
func (w *Writer) WriteJSON(node *scanner.FileNode) error {
	w.walkCount(node)
	w.writeJSONNode(node, 0)
	return nil
}

func (w *Writer) writeJSONNode(node *scanner.FileNode, indent int) {
	prefix := strings.Repeat("  ", indent)
	if node.IsDir {
		fmt.Fprintf(w.writer, "%s{\n", prefix)
		fmt.Fprintf(w.writer, "%s  \"name\": %q,\n", prefix, node.Name)
		fmt.Fprintf(w.writer, "%s  \"type\": \"directory\",\n", prefix)
		fmt.Fprintf(w.writer, "%s  \"children\": %d,\n", prefix, node.NumChildren)
		fmt.Fprintf(w.writer, "%s  \"totalSize\": %d,\n", prefix, node.TotalSize)
		if w.showGitStatus && node.GitStatus != "" {
			fmt.Fprintf(w.writer, "%s  \"gitStatus\": %q,\n", prefix, node.GitStatus)
		}
		fmt.Fprintf(w.writer, "%s  \"items\": [\n", prefix)
		for i, child := range node.Children {
			w.writeJSONNode(child, indent+2)
			if i < len(node.Children)-1 {
				fmt.Fprint(w.writer, ",")
			}
			fmt.Fprintln(w.writer)
		}
		fmt.Fprintf(w.writer, "%s  ]\n", prefix)
		fmt.Fprintf(w.writer, "%s}", prefix)
	} else {
		fmt.Fprintf(w.writer, "%s{\n", prefix)
		fmt.Fprintf(w.writer, "%s  \"name\": %q,\n", prefix, node.Name)
		fmt.Fprintf(w.writer, "%s  \"type\": \"file\",\n", prefix)
		fmt.Fprintf(w.writer, "%s  \"size\": %d,\n", prefix, node.Size)
		fmt.Fprintf(w.writer, "%s  \"mode\": %q,\n", prefix, node.Mode.Perm().String())
		fmt.Fprintf(w.writer, "%s  \"modTime\": %q,\n", prefix, node.ModTime.Format("2006-01-02T15:04:05Z"))
		if w.showGitStatus && node.GitStatus != "" {
			fmt.Fprintf(w.writer, "%s  \"gitStatus\": %q,\n", prefix, node.GitStatus)
		}
		fmt.Fprintf(w.writer, "%s}", prefix)
	}
}

// WriteList outputs a flat list of files
func (w *Writer) WriteList(node *scanner.FileNode) error {
	w.writeListNode(node, "")
	return nil
}

func (w *Writer) writeListNode(node *scanner.FileNode, prefix string) {
	if !node.IsDir {
		return
	}

	for _, child := range node.Children {
		path := filepath.Join(prefix, child.Name)
		if child.IsDir {
			fmt.Fprintf(w.writer, "%s/\n", path)
			w.writeListNode(child, path)
		} else {
			var meta []string
			if w.showSizes {
				meta = append(meta, formatSize(child.Size))
			}
			if w.showGitStatus && child.GitStatus != "" {
				meta = append(meta, string(child.GitStatus))
			}
			metaStr := ""
			if len(meta) > 0 {
				metaStr = " (" + strings.Join(meta, ", ") + ")"
			}
			fmt.Fprintf(w.writer, "%s%s\n", path, metaStr)
		}
	}
}

// GetStats returns the stats
func (w *Writer) GetStats() (files, dirs int, size int64) {
	return w.totalFiles, w.totalDirs, w.totalSize
}

// colorize adds color to the name based on type and git status
func colorize(name string, node *scanner.FileNode) string {
	if node.IsDir {
		return "\033[1;34m" + name + "\033[0m" // Bold blue for dirs
	}

	// Color by git status
	if node.GitStatus != "" {
		switch node.GitStatus {
		case "modified":
			return "\033[33m" + name + "\033[0m" // Yellow
		case "added":
			return "\033[32m" + name + "\033[0m" // Green
		case "deleted":
			return "\033[31m" + name + "\033[0m" // Red
		case "untracked":
			return "\033[36m" + name + "\033[0m" // Cyan
		}
	}

	// Color by extension
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".go", ".rs", ".py", ".js", ".ts", ".java", ".c", ".cpp", ".h":
		return "\033[32m" + name + "\033[0m" // Green for code
	case ".md", ".txt", ".rst":
		return "\033[36m" + name + "\033[0m" // Cyan for docs
	case ".json", ".yaml", ".yml", ".toml", ".xml":
		return "\033[35m" + name + "\033[0m" // Magenta for config
	case ".png", ".jpg", ".jpeg", ".gif", ".svg":
		return "\033[33m" + name + "\033[0m" // Yellow for images
	}

	return name
}

// formatSize formats bytes into human-readable size
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
