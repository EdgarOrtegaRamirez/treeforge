package git

import (
	"os/exec"
	"strings"
)

// StatusMap maps file paths to their git status
type StatusMap map[string]string

// GetStatus runs git status --porcelain and returns a map of file -> status
func GetStatus(repoPath string) StatusMap {
	result := make(StatusMap)

	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return result
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		status := line[:2]
		filePath := strings.TrimSpace(line[3:])

		var gitStatus string
		switch {
		case status[0] == 'M' || status[1] == 'M':
			gitStatus = "modified"
		case status[0] == 'A':
			gitStatus = "added"
		case status[0] == 'D' || status[1] == 'D':
			gitStatus = "deleted"
		case status[0] == 'R':
			gitStatus = "renamed"
		case status[0] == 'C':
			gitStatus = "copied"
		case status == "??":
			gitStatus = "untracked"
		case status[0] == '!' || status[1] == '!':
			gitStatus = "ignored"
		}

		if gitStatus != "" {
			result[filePath] = gitStatus
		}
	}

	return result
}

// IsRepo checks if the given path is inside a git repository
func IsRepo(path string) bool {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}
