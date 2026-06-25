package main

import (
	"strings"
)

func cleanInput(text string) []string {
	var result []string
	sep := " "
	i := strings.Index(text, sep)

	for i > -1 {
		word := strings.ToLower(text[:i])
		if word != "" {
			result = append(result, strings.ToLower(text[:i]))
		}
		text = text[i+len(sep):]
		i = strings.Index(text, sep)
	}

	if text != "" {
		result = append(result, strings.ToLower(text))
	}

	return result
}
