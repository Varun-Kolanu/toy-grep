package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2)
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin)
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

}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	var ok bool = false
	var toMatch strings.Builder
	if pattern == "\\d" {
		for c := '0'; c <= '9'; c++ {
			toMatch.WriteString(string(c))
		}
	} else if pattern == "\\w" {
		for c := '0'; c <= '9'; c++ {
			toMatch.WriteString(string(c))
		}
		for c := 'a'; c <= 'z'; c++ {
			toMatch.WriteString(string(c))
		}
		for c := 'A'; c <= 'Z'; c++ {
			toMatch.WriteString(string(c))
		}
	} else if pattern[0] == '[' && pattern[len(pattern)-1] == ']' {
		toMatch.WriteString(pattern[1 : len(pattern)-1])
	} else {
		toMatch.WriteString(pattern)
	}

	matchString := toMatch.String()
	ok = bytes.ContainsAny(line, matchString)

	return ok, nil
}
