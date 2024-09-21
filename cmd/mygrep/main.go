package main

import (
	"fmt"
	"io"
	"os"
	"strings"
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
	var ok bool = false

	lineLength := len(line)

	for ind := 0; ind < lineLength; ind++ {
		if matchChar(line, pattern, ind) {
			ok = true
			break
		}
		if pattern[0] == '^' {
			break
		}
	}
	return ok, nil
}

func matchDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func matchAlphaNumeric(char byte) bool {
	return matchDigit(char) || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_'
}

func matchPositiveCharacterGroup(char byte, str string) bool {
	for _, c := range str {
		if c == rune(char) {
			return true
		}
	}
	return false
}

func matchNegativeCharacterGroup(char byte, str string) bool {
	for _, c := range str {
		if c == rune(char) {
			return false
		}
	}
	return true
}

func matchChar(line []byte, pattern string, ind int) bool {
	if ind >= len(line) {
		return false
	}
	char := line[ind]
	patternIndex := 0
	if pattern[0] == '^' {
		patternIndex = 1
	} else if len(pattern) >= 2 && pattern[:2] == "\\d" {
		if !matchDigit(char) {
			return false
		}
		patternIndex = 2
		ind++
	} else if len(pattern) >= 2 && pattern[:2] == "\\w" {
		if !matchAlphaNumeric(char) {
			return false
		}
		patternIndex = 2
		ind++
	} else if pattern[0] == '[' {
		rightBracket := strings.Index(pattern, "]")
		if rightBracket == -1 {
			fmt.Fprintf(os.Stderr, "unsupported pattern\n")
			os.Exit(2)
		}
		charGroup := pattern[1:rightBracket]
		if charGroup[0] == '^' {
			charGroup = charGroup[1:]
			if !matchNegativeCharacterGroup(char, charGroup) {
				return false
			}
		} else {
			if !matchPositiveCharacterGroup(char, charGroup) {
				return false
			}
		}
		patternIndex = rightBracket + 1
		ind++
	} else {
		if char != byte(pattern[0]) {
			return false
		}
		patternIndex = 1
		ind++
	}
	if patternIndex >= len(pattern) {
		return true
	}
	return matchChar(line, pattern[patternIndex:], ind)
}
