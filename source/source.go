// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package source defines interfaces to be implemented by dictionary sources and
// common structures and operations for those implementations to use.
package source

// Source defines an interface for interacting with different dictionaries
type Source interface {
	Name() string

	Define(word string) (Result, error)
}

// Result defines an interface for the results of a dictionary lookup
type Result interface {
	Headword() string
	Language() string
	Entries() []DictionaryEntry
}

// Entry defines a composite interface for the complete account of a word
type Entry interface {
	WordEntry
	DictionaryEntry
	EtymologyEntry
	ThesaurusEntry
}

// ComprehensiveDictionaryEntry defines a composite interface for a
// comprehensive dictionary entry of a word
type ComprehensiveDictionaryEntry interface {
	DictionaryEntry
	EtymologyEntry
}

// VersatileDictionaryEntry defines a composite interface for a versatile
// dictionary entry of a word
type VersatileDictionaryEntry interface {
	DictionaryEntry
	ThesaurusEntry
}

// WordEntry defines an interface for an entry of a specific word
type WordEntry interface {
	Word() string
	Category() string
}

// DictionaryEntry defines an interface for a dictionary entry of a word
type DictionaryEntry interface {
	Pronunciation() string
	Senses() []Sense
}

// EtymologyEntry defines an interface for an etymological entry of a word
type EtymologyEntry interface {
	Etymologies() []string
}

// ThesaurusEntry defines an interface for a thesaurus entry of a word
type ThesaurusEntry interface {
	Synonyms() []string
	Antonyms() []string
}

// Sense defines an interface for the different meanings of a word
type Sense interface {
	Definitions() []string
	Examples() []string
	Notes() []string

	Subsenses() []Sense
}
