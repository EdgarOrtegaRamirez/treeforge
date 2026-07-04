package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/EdgarOrtegaRamirez/treeforge/pkg/git"
	"github.com/EdgarOrtegaRamirez/treeforge/pkg/output"
	"github.com/EdgarOrtegaRamirez/treeforge/pkg/scanner"
	"github.com/spf13/cobra"
)

var (
	showSizes    bool
	showPerms    bool
	showModTime  bool
	showGit      bool
	maxDepth     int
	showHidden   bool
	excludePat   string
	includePat   string
	sortBy       string
	sortReverse  bool
	outputFormat string
	version      = "1.0.0"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "treeforge [path]",
		Short: "A modern, feature-rich directory tree visualization tool",
		Long: `TreeForge is a modern directory tree visualization tool with features like
file sizes, permissions, git status integration, glob filtering, and multiple
output formats. It's a modern replacement for the classic 'tree' command.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runTree,
	}

	// Global flags
	rootCmd.Flags().BoolVarP(&showSizes, "sizes", "s", false, "Show file sizes")
	rootCmd.Flags().BoolVarP(&showPerms, "perms", "p", false, "Show file permissions")
	rootCmd.Flags().BoolVarP(&showModTime, "time", "t", false, "Show modification times")
	rootCmd.Flags().BoolVarP(&showGit, "git", "g", false, "Show git status colors")
	rootCmd.Flags().IntVarP(&maxDepth, "depth", "d", -1, "Max depth (0 = root only, -1 = unlimited)")
	rootCmd.Flags().BoolVarP(&showHidden, "hidden", "a", false, "Show hidden files (starting with .)")
	rootCmd.Flags().StringVarP(&excludePat, "exclude", "e", "", "Exclude pattern (glob)")
	rootCmd.Flags().StringVarP(&includePat, "include", "i", "", "Include pattern (glob)")
	rootCmd.Flags().StringVar(&sortBy, "sort", "name", "Sort by: name, size, time")
	rootCmd.Flags().BoolVar(&sortReverse, "reverse", false, "Reverse sort order")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "tree", "Output format: tree, json, compact, list")

	// Version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("treeforge v%s\n", version)
		},
	}
	rootCmd.AddCommand(versionCmd)

	// Stats command
	statsCmd := &cobra.Command{
		Use:   "stats [path]",
		Short: "Show directory statistics",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runStats,
	}
	rootCmd.AddCommand(statsCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runTree(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	config := scanner.ScannerConfig{
		RootPath:        absPath,
		MaxDepth:        maxDepth,
		ShowHidden:      showHidden,
		ShowSizes:       showSizes,
		ShowPermissions: showPerms,
		ShowModTime:     showModTime,
		IncludePattern:  includePat,
		ExcludePattern:  excludePat,
		SortBy:          sortBy,
		SortReverse:     sortReverse,
	}

	// Check for git status
	if showGit && git.IsRepo(absPath) {
		config.ShowGitStatus = true
	}

	s := scanner.NewScanner(config)
	root, err := s.Scan()
	if err != nil {
		return fmt.Errorf("scan error: %w", err)
	}

	// Apply git status if requested
	if config.ShowGitStatus {
		statusMap := git.GetStatus(absPath)
		applyGitStatus(root, statusMap)
	}

	// Choose output format
	w := output.NewWriter(os.Stdout,
		output.WithSizes(showSizes),
		output.WithPerms(showPerms),
		output.WithModTime(showModTime),
		output.WithGitStatus(config.ShowGitStatus),
	)

	switch outputFormat {
	case "json":
		err = w.WriteJSON(root)
	case "compact":
		err = w.WriteCompact(root)
	case "list":
		err = w.WriteList(root)
	default:
		err = w.WriteTree(root)
	}

	if err != nil {
		return fmt.Errorf("output error: %w", err)
	}

	return nil
}

func runStats(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	config := scanner.ScannerConfig{
		RootPath:   absPath,
		MaxDepth:   -1,
		ShowHidden: showHidden,
	}

	s := scanner.NewScanner(config)
	_, err = s.Scan()
	if err != nil {
		return fmt.Errorf("scan error: %w", err)
	}

	stats := s.GetStats()

	fmt.Printf("Statistics for: %s\n", absPath)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("  Files:      %d\n", stats.TotalFiles)
	fmt.Printf("  Directories: %d\n", stats.TotalDirs)
	fmt.Printf("  Total Size: %s\n", formatSize(stats.TotalSize))

	if stats.LargestFile != "" {
		fmt.Printf("  Largest:    %s (%s)\n", stats.LargestFile, formatSize(stats.LargestSize))
	}

	if !stats.OldestTime.IsZero() {
		fmt.Printf("  Oldest:     %s (%s)\n", stats.OldestFile, stats.OldestTime.Format("2006-01-02"))
	}
	if !stats.NewestTime.IsZero() {
		fmt.Printf("  Newest:     %s (%s)\n", stats.NewestFile, stats.NewestTime.Format("2006-01-02"))
	}

	if len(stats.ByExtension) > 0 {
		fmt.Println("\nBy Extension:")
		for ext, count := range stats.ByExtension {
			fmt.Printf("  %-10s %d\n", ext, count)
		}
	}

	if len(stats.ByGitStatus) > 0 {
		fmt.Println("\nBy Git Status:")
		for status, count := range stats.ByGitStatus {
			if status != "" {
				fmt.Printf("  %-12s %d\n", status, count)
			}
		}
	}

	return nil
}

func applyGitStatus(node *scanner.FileNode, statusMap git.StatusMap) {
	relPath := node.RelativePath
	if status, ok := statusMap[relPath]; ok {
		node.GitStatus = scanner.GitStatus(status)
	}
	for _, child := range node.Children {
		applyGitStatus(child, statusMap)
	}
}

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
