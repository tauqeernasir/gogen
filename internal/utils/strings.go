package utils

import (
	"regexp"
	"strings"
)

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	words := regexp.MustCompile(`[^a-zA-Z0-9]+`).Split(s, -1)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if len(pascal) > 0 {
		return strings.ToLower(pascal[:1]) + pascal[1:]
	}
	return pascal
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	words := regexp.MustCompile(`[^a-zA-Z0-9]+`).Split(s, -1)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "_")
}

// Contains checks if a slice contains a specific string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
