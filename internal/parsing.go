package internal

import (
	"strings"
)

func ParseInput(input string) ([]string, error) {
	words := strings.Fields(input)

	for i, w := range words {
		words[i] = strings.ToLower(w)
	}

	return words, nil
}