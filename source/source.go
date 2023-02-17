// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package source defines interfaces to be implemented by dictionary sources and
// common structures and operations for those implementations to use.
package source

// Source defines an interface for interacting with different dictionaries
type Source interface {
	// Name returns the printable, human-readable name of the source.
	Name() string

	// Define takes a word string and returns a list of dictionary results, and
	// an error if any occurred.
	Define(word string) ([]DictionaryResult, error)
}

// DictionaryResult defines the structure of a dictionary word result in a
// specific language
type DictionaryResult struct {
	Language string
	Entries  []DictionaryEntry
}

// Entry defines the structure for an entry of a specific word
type Entry struct {
	Word            string
	LexicalCategory string
}

// DictionaryEntry defines the structure for a dictionary entry of a word
type DictionaryEntry struct {
	Entry

	Pronunciations []string
	Senses         []Sense
	Etymologies    []string // Origins of the word

	ThesaurusValues
}

// Sense defines the structure of a particular meaning of a word
type Sense struct {
	Definitions []string
	Examples    []string
	Notes       []string

	ThesaurusValues

	SubSenses []Sense
}

// ThesaurusValues defines the structure for the thesaurus values of a word
type ThesaurusValues struct {
	Synonyms []string // Words with similar meaning
	Antonyms []string // Words with the opposite meaning
}
