package main

import (
	"bufio"
	"os"
	"strings"
)

// parseIgnoreFile reads the ignore file and returns a slice of patterns to ignore
func parseIgnoreFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return patterns, nil
}

// main function for demonstration
func main() {
	patterns, err := parseIgnoreFile(".projecttextignore")
	if err != nil {
		panic(err)
	}

	for _, pattern := range patterns {
		println(pattern)
	}
}
