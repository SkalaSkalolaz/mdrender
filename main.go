package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/yuin/goldmark"
)

const version = "1.0.0"

// ReadMarkdownFile reads the content of a markdown file and returns it as a string.
func ReadMarkdownFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	return string(content), nil
}

// ConvertMarkdownToHTML takes a markdown string and returns the converted HTML.
func ConvertMarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown: %w", err)
	}
	return buf.String(), nil
}

// RenderMarkdownToTerminal takes a markdown string and renders it with terminal formatting.
func RenderMarkdownToTerminal(markdown string) (string, error) {
	out, err := glamour.Render(markdown, "dark")
	if err != nil {
		return "", fmt.Errorf("failed to render markdown for terminal: %w", err)
	}
	return out, nil
}

func main() {
	var filePath string

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--help":
			fmt.Println("Usage: go run main.go [options] [file.md]\n\nOptions:\n  --help     Show this help message\n  --version Show the version")
			os.Exit(0)
		case "--version":
			fmt.Printf("Version: %s\n", version)
			os.Exit(0)
		default:
			filePath = arg
		}
	}

	if filePath == "" {
		fmt.Println("Error: No markdown file provided.")
		fmt.Println("Usage: go run main.go [options] [file.md]")
		os.Exit(1)
	}

	markdownContent, err := ReadMarkdownFile(filePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	terminalOutput, err := RenderMarkdownToTerminal(markdownContent)
	if err != nil {
		fmt.Printf("Error rendering for terminal: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(terminalOutput)
}
