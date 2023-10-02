package util

import (
	"strings"
)

func EscapeSpecialChars(input string) string {
	input = strings.ReplaceAll(input, `\`, `\\`)
	replacements := map[string]string{
		`:`: `\:`,
		`*`: `\*`,
		`'`: `''`,
	}
	for char, replacement := range replacements {
		input = strings.ReplaceAll(input, char, replacement)
	}
	return input
}
