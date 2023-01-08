package main

import (
	// Uncomment this to pass the first stage
	// "bytes"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Usage: echo <input_text> | your_grep.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	switch {
	// digits (\d)
	case pattern == `\d`:
		for _, char := range string(line) {
			if unicode.IsDigit(char) {
				return true, nil
			}
		}
		return false, nil

	// alphanumerice characters (\w)
	case pattern == `\w`:
		for _, char := range string(line) {
			if unicode.IsDigit(char) || unicode.IsLetter(char) {
				return true, nil
			}
		}
		return false, nil

	// negative charcter groups (e.g. [^abc])
	case strings.HasPrefix(pattern, "[^") && strings.HasSuffix(pattern, "]"):
		negative_chars := pattern[2 : len(pattern)-1]
		for _, char := range negative_chars {
			if bytes.ContainsAny(line, string(char)) {
				return false, nil
			}
		}
		return true, nil

	// positive charcter groups (e.g. [abc])
	case strings.HasPrefix(pattern, "[") && strings.HasSuffix(pattern, "]"):
		positive_chars := pattern[1 : len(pattern)-1]
		for _, char := range positive_chars {
			if bytes.ContainsAny(line, string(char)) {
				return true, nil
			}
		}
		return false, nil

	// single character (e.g. a)
	case utf8.RuneCountInString(pattern) == 1:
		return bytes.ContainsAny(line, pattern), nil
	}

	return false, fmt.Errorf("unsupported pattern: %q", pattern)
}
