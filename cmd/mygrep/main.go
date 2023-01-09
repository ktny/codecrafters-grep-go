package main

import (
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
	for i := range string(line) {
		ok, err := matchHere(line[i:], pattern)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func matchHere(line []byte, pattern string) (bool, error) {
	if pattern == "" {
		return true, nil
	}

	if len(line) == 0 {
		return false, nil
	}

	char, size := utf8.DecodeRune(line)

	switch {
	// digits (\d)
	case strings.HasPrefix(pattern, `\d`):
		if unicode.IsDigit(char) {
			return matchHere(line[size:], pattern[2:])
		}
		return false, nil

	// alphanumerice characters (\w)
	case strings.HasPrefix(pattern, `\w`):
		if unicode.IsDigit(char) || unicode.IsLetter(char) {
			// fmt.Printf("char: %v", char)
			return matchHere(line[size:], pattern[2:])
		}
		return false, nil

	// negative charcter groups (e.g. [^abc])
	case strings.HasPrefix(pattern, "[^"):
		end := strings.IndexByte(pattern, ']')
		negative_chars := pattern[2:end]
		if !strings.ContainsRune(negative_chars, char) {
			return matchHere(line[size:], pattern[end+1:])
		}
		return false, nil

	// positive charcter groups (e.g. [abc])
	case strings.HasPrefix(pattern, "["):
		end := strings.IndexByte(pattern, ']')
		positive_chars := pattern[1:end]
		if strings.ContainsRune(positive_chars, char) {
			return matchHere(line[size:], pattern[end+1:])
		}
		return false, nil

	// non regexp chars
	default:
		patternChar, patternCharSize := utf8.DecodeRuneInString(pattern)
		if char == patternChar {
			return matchHere(line[size:], pattern[patternCharSize:])
		}
	}

	return false, nil
}
