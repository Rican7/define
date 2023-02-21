package source

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// RemoveDiacritics takes a text and returns the same text with any diacritics
// removed. If there's an issue with cleaning the string, the original text is
// returned unchanged.
func RemoveDiacritics(text string) string {
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalized, _, err := transform.String(transformer, text)

	if err != nil {
		return text
	}

	return normalized
}

// EqualFoldPlain reports whether s and t, interpreted as UTF-8 strings, are
// equal under Unicode case-folding AFTER first having each string normalized.
func EqualFoldPlain(s, t string) bool {
	return strings.EqualFold(
		RemoveDiacritics(s),
		RemoveDiacritics(t),
	)
}
