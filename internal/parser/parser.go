package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ReadFile reads the code from a file
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DetectLanguage detects the programming language from file extension
func DetectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	
	// Map common extensions to chroma language names
	langMap := map[string]string{
		".go":     "go",
		".py":     "python",
		".js":     "javascript",
		".jsx":    "jsx",
		".ts":     "typescript",
		".tsx":    "tsx",
		".rs":     "rust",
		".rb":     "ruby",
		".java":   "java",
		".c":      "c",
		".cpp":    "cpp",
		".cc":     "cpp",
		".h":      "c",
		".hpp":    "cpp",
		".cs":     "csharp",
		".php":    "php",
		".swift":  "swift",
		".kt":     "kotlin",
		".scala":  "scala",
		".sh":     "bash",
		".bash":   "bash",
		".zsh":    "bash",
		".fish":   "fish",
		".ps1":    "powershell",
		".r":      "r",
		".sql":    "sql",
		".html":   "html",
		".css":    "css",
		".scss":   "scss",
		".sass":   "sass",
		".json":   "json",
		".yaml":   "yaml",
		".yml":    "yaml",
		".toml":   "toml",
		".xml":    "xml",
		".md":     "markdown",
		".vim":    "vim",
		".lua":    "lua",
		".pl":     "perl",
		".ex":     "elixir",
		".exs":    "elixir",
		".erl":    "erlang",
		".hs":     "haskell",
		".clj":    "clojure",
		".lisp":   "commonlisp",
		".dart":   "dart",
		".vue":    "vue",
		".svelte": "svelte",
	}

	if lang, ok := langMap[ext]; ok {
		return lang
	}

	// Default to text if unknown
	return "text"
}

// ParseHighlightLines parses highlight line specification (e.g., "5,7-9" -> [5,7,8,9])
func ParseHighlightLines(spec string) ([]int, error) {
	if spec == "" {
		return nil, nil
	}

	var lines []int
	parts := strings.Split(spec, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		
		// Check for range (e.g., "7-9")
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start line: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end line: %s", rangeParts[1])
			}

			if start > end {
				return nil, fmt.Errorf("start line (%d) must be <= end line (%d)", start, end)
			}

			// Add all lines in range
			for i := start; i <= end; i++ {
				lines = append(lines, i)
			}
		} else {
			// Single line number
			line, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid line number: %s", part)
			}
			lines = append(lines, line)
		}
	}

	return lines, nil
}
