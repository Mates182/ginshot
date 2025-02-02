package formatter

import (
	"strings"
	"unicode"
)

// ToPascalCase converts a lowercase string with hyphens to PascalCase
func ToPascalCase(input string) string {
	// Split the string by hyphens
	words := strings.Split(input, "_")

	// Capitalize each word and concatenate them
	for i, word := range words {
		words[i] = strings.Title(word)
	}

	// Join the words together to form the PascalCase string
	return strings.Join(words, "")
}

// ToLowerCase converts a PascalCase string to lowercase with hyphens
func ToLowerCase(input string) string {
	var result []rune

	for i, r := range input {
		// If it's an uppercase letter and not the first character, add a hyphen before it
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		// Append the lowercase version of the character
		result = append(result, unicode.ToLower(r))
	}

	return string(result)
}
