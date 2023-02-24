// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package source defines interfaces to be implemented by dictionary sources and
// common structures and operations for those implementations to use.
package source

import (
	"fmt"
	"strings"
)

// Source defines an interface for interacting with different dictionaries
type Source interface {
	// Name returns the printable, human-readable name of the source.
	Name() string

	// Define takes a word string and returns a list of dictionary results, and
	// an error if any occurred.
	Define(word string) (DictionaryResults, error)
}

// Searcher defines an interface for a source that supports search capabilities
type Searcher interface {
	// Search takes a word string and returns a list of found words, and an
	// error if any occurred.
	Search(word string, limit uint) (SearchResults, error)
}

// DictionaryResults defines the structure of a list of dictionary word results
type DictionaryResults []DictionaryResult

// SearchResults defines the structure of a list of word search results
type SearchResults []SearchResult

// DictionaryResult defines the structure of a dictionary word result in a
// specific language
type DictionaryResult struct {
	Language string
	Word     string
	Entries  []DictionaryEntry
}

// SearchResult defines the structure of a word search result
type SearchResult string

// Entry defines the structure of an entry of a specific word
type Entry struct {
	Word            string
	LexicalCategory string
}

// DictionaryEntry defines the structure of a dictionary entry of a word
type DictionaryEntry struct {
	Entry

	Senses      []Sense
	Etymologies []string // Origins of the word

	Pronunciations
	ThesaurusValues
}

// Pronunciations defines the structure of a collection of pronunciations
type Pronunciations []Pronunciation

// Pronunciation defines the structure of a pronunciation of a word
type Pronunciation string

// Sense defines the structure of a particular meaning of a word
type Sense struct {
	Definitions []string
	Categories  []string
	Examples    []AttributedText
	Notes       []string

	ThesaurusValues

	SubSenses []Sense
}

// AttributedText defines the structure of a general text with attribution
type AttributedText struct {
	Text string

	Attribution
}

// Attribution defines the structure of a general attribution of a data piece
type Attribution struct {
	Author string
	Source string
}

// ThesaurusValues defines the structure of the thesaurus values of a word
type ThesaurusValues struct {
	Synonyms []string // Words with similar meaning
	Antonyms []string // Words with the opposite meaning
}

// IsSortedForPrimaryResult takes a word and returns true if the first result
// is a direct match.
func (r DictionaryResults) IsSortedForPrimaryResult(word string) bool {
	return len(r) < 1 || r[0].Word == word
}

// SortForPrimaryResult takes a word and attempts to find any direct match
// within the results. If any match is found, it moves that result to the
// first position of the results (while retaining the other results).
func (r *DictionaryResults) SortForPrimaryResult(word string) {
	if r.IsSortedForPrimaryResult(word) {
		// Short-circuit if we're already sorted
		return
	}

	var matchIndex *int

	for i, result := range *r {
		if result.Word == word {
			matchIndex = &i
			break
		}
	}

	// Make sure the matchIndex isn't 0, because then the results are already
	// sorted (the match is already in the first position)
	if matchIndex != nil && *matchIndex != 0 {
		i := *matchIndex
		match := (*r)[i]

		// Sort the slice by placing the matched result at the first position
		*r = append(DictionaryResults{match}, append((*r)[:i], (*r)[i+1:]...)...)
	}
}

// String satisfies fmt.Stringer and dictates the string format of the value
func (p Pronunciations) String() string {
	var pronunciationText string

	if len(p) > 0 {
		pronunciationText = p[0].String()
	}

	if len(p) > 1 {
		var pronunciationStrings []string
		for _, pronunciation := range p {
			pronunciationStrings = append(pronunciationStrings, pronunciation.String())
		}

		pronunciationText = fmt.Sprintf("%s (%s)", pronunciationText, strings.Join(pronunciationStrings[1:], " "))
	}

	return pronunciationText
}

// String satisfies fmt.Stringer and dictates the string format of the value
func (p Pronunciation) String() string {
	return fmt.Sprintf("/%s/", string(p))
}

// String satisfies fmt.Stringer and dictates the string format of the value
func (t AttributedText) String() string {
	text := fmt.Sprintf("%q", t.Text)

	if t.Author != "" {
		text = fmt.Sprintf("%s - %s", text, t.Author)
	}

	if t.Source != "" {
		text = fmt.Sprintf("%s (%s)", text, t.Source)
	}

	return text
}
