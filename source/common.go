package source

// A ResultValue contains the common attributes of a dictionary lookup result
type ResultValue struct {
	Head      string
	Lang      string
	EntryVals []interface{}
}

// An EntryValue contains the common attributes of a complete account of a
// particular word
type EntryValue struct {
	DictionaryEntryValue
	EtymologyEntryValue
	ThesaurusEntryValue
}

// An EtymologyEntryValue contains the common attributes of an etymological
// entry of a word
type EtymologyEntryValue struct {
	EtymologyVals []string
}

// A DictionaryEntryValue contains the common attributes of a dictionary entry
// of a word
type DictionaryEntryValue struct {
	PronunciationVal string
	SenseVals        []SenseValue
}

// A ThesaurusEntryValue contains the common attributes of a thesaurus entry
// of a word
type ThesaurusEntryValue struct {
	SynonymVals []string
	AntonymVals []string
}

// A SenseValue contains the common attributes of a word's meanings
type SenseValue struct {
	DefinitionVals []string
	ExampleVals    []string
	NoteVals       []string
}

// Headword returns the result's headword
func (r ResultValue) Headword() string {
	return r.Head
}

// Language returns the result's language
func (r ResultValue) Language() string {
	return r.Lang
}

// Entries returns the result's entries
func (r ResultValue) Entries() []DictionaryEntry {
	entries := make([]DictionaryEntry, len(r.EntryVals))

	for i, entry := range r.EntryVals {
		if dictionaryEntry, ok := entry.(DictionaryEntry); ok {
			entries[i] = dictionaryEntry
		} else {
			panic("Invalid type in set")
		}
	}

	return entries
}

// Etymologies returns the entry's etymology strings
func (e EtymologyEntryValue) Etymologies() []string {
	return e.EtymologyVals
}

// Pronunciation returns the entry's pronunciation representation
func (e DictionaryEntryValue) Pronunciation() string {
	return e.PronunciationVal
}

// Senses returns the entry's senses
func (e DictionaryEntryValue) Senses() []Sense {
	meanings := make([]Sense, len(e.SenseVals))

	for i, meaning := range e.SenseVals {
		meanings[i] = meaning
	}

	return meanings
}

// Synonyms returns the entry's synonyms
func (e ThesaurusEntryValue) Synonyms() []string {
	return e.SynonymVals
}

// Antonyms returns the entry's antonyms
func (e ThesaurusEntryValue) Antonyms() []string {
	return e.AntonymVals
}

// Definitions returns the sense's definitions
func (s SenseValue) Definitions() []string {
	return s.DefinitionVals
}

// Examples returns the sense's examples
func (s SenseValue) Examples() []string {
	return s.ExampleVals
}

// Notes returns the sense's notes
func (s SenseValue) Notes() []string {
	return s.NoteVals
}
