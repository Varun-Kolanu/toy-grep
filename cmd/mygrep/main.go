package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
	if utf8.RuneCountInString(pattern) > 2 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	var ok bool = false
	if pattern == "\\d" {
		ok = bytes.ContainsAny(line, "0123456789")
	} else if pattern == "\\w" {
		for _, c := range line {
			if (c > '0' && c < '9') || (c > 'A' && c < 'Z') || (c > 'a' && c < 'z') || c == '_' {
				fmt.Println(string(c))
				ok = true
				break
			}
		}
	} else {
		ok = bytes.ContainsAny(line, pattern)
	}

	return ok, nil
}
