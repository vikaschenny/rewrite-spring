package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/openrewrite/rewrite-spring-go/pkg/core"
)

// MatchGlob checks if a file path matches a glob pattern
func MatchGlob(path, pattern string) (bool, error) {
	// Convert glob pattern to regex
	regexPattern := globToRegex(pattern)
	matched, err := regexp.MatchString(regexPattern, path)
	return matched, err
}

// globToRegex converts a glob pattern to a regular expression
func globToRegex(glob string) string {
	// Escape special regex characters except * and ?
	escaped := regexp.QuoteMeta(glob)

	// Replace escaped glob characters with regex equivalents
	escaped = strings.ReplaceAll(escaped, "\\*", ".*")
	escaped = strings.ReplaceAll(escaped, "\\?", ".")

	// Handle ** for recursive directory matching
	escaped = strings.ReplaceAll(escaped, "\\.\\*\\.\\*", ".*")

	return "^" + escaped + "$"
}

// LoadSourceFile loads a source file from disk
func LoadSourceFile(filePath string) (core.SourceFile, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	fileType := DetermineFileType(filePath)

	return &core.SpringConfigFile{
		Path:    filePath,
		Content: string(content),
		Type:    fileType,
	}, nil
}

// SaveSourceFile saves a source file to disk
func SaveSourceFile(sourceFile core.SourceFile, outputPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file content
	err := ioutil.WriteFile(outputPath, []byte(sourceFile.GetContent()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", outputPath, err)
	}

	return nil
}

// DetermineFileType determines the file type based on extension
func DetermineFileType(filePath string) core.FileType {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".properties":
		return core.Properties
	case ".yml", ".yaml":
		return core.YAML
	case ".java":
		return core.Java
	default:
		// Check if it's a properties file without extension
		if strings.Contains(filepath.Base(filePath), "application") {
			return core.Properties
		}
		return core.Properties // Default fallback
	}
}

// FindSpringConfigFiles finds Spring configuration files in a directory
func FindSpringConfigFiles(rootDir string, patterns []string) ([]string, error) {
	var configFiles []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file matches any of the patterns
		for _, pattern := range patterns {
			if matched, _ := MatchGlob(path, pattern); matched {
				configFiles = append(configFiles, path)
				break
			}
		}

		return nil
	})

	return configFiles, err
}

// ReadFileLines reads a file and returns its lines
func ReadFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// WriteFileLines writes lines to a file
func WriteFileLines(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// IsSpringConfigFile checks if a file is likely a Spring configuration file
func IsSpringConfigFile(filePath string) bool {
	fileName := strings.ToLower(filepath.Base(filePath))

	// Check for common Spring configuration file patterns
	patterns := []string{
		"application.properties",
		"application.yml",
		"application.yaml",
		"application-*.properties",
		"application-*.yml",
		"application-*.yaml",
		"bootstrap.properties",
		"bootstrap.yml",
		"bootstrap.yaml",
	}

	for _, pattern := range patterns {
		if matched, _ := MatchGlob(fileName, pattern); matched {
			return true
		}
	}

	return false
}

// BackupFile creates a backup of a file
func BackupFile(filePath string) error {
	backupPath := filePath + ".backup"

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read original file: %w", err)
	}

	err = ioutil.WriteFile(backupPath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}

	return nil
}
