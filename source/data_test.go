// Copyright © 2018 Trevor N. Suarez (Rican7)

package source

import (
	"reflect"
	"testing"
)

// Enforce interface contracts
var (
	_ Result          = (*ResultValue)(nil)
	_ Entry           = (*EntryValue)(nil)
	_ WordEntry       = (*WordEntryValue)(nil)
	_ DictionaryEntry = (*DictionaryEntryValue)(nil)
	_ EtymologyEntry  = (*EtymologyEntryValue)(nil)
	_ ThesaurusEntry  = (*ThesaurusEntryValue)(nil)
	_ Sense           = (*SenseValue)(nil)
)

func TestHeadword(t *testing.T) {
	r := ResultValue{Head: "test"}

	got := r.Headword()
	want := r.Head

	if got != want {
		t.Errorf("Headword returned wrong value. Got %v. Want %v.", got, want)
	}
}

func TestLanguage(t *testing.T) {
	r := ResultValue{Lang: "test"}

	got := r.Language()
	want := r.Lang

	if got != want {
		t.Errorf("Language returned wrong value. Got %v. Want %v.", got, want)
	}
}

func TestEntries(t *testing.T) {
	entries := []interface{}{
		EntryValue{},
	}
	r := ResultValue{EntryVals: entries}

	for i, entry := range r.Entries() {
		got := entry
		want := entries[i].(EntryValue)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Entries returned wrong value. Got %v. Want %v.", got, want)
		}
	}
}

func TestEntriesPanicsOnInvalidType(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Errorf("Entries with an invalid type did not panic.")
		}
	}()

	ResultValue{
		EntryVals: []interface{}{
			1234,
		},
	}.Entries()
}

func TestWord(t *testing.T) {
	e := WordEntryValue{WordVal: "test"}

	got := e.Word()
	want := e.WordVal

	if got != want {
		t.Errorf("Word returned wrong value. Got %v. Want %v.", got, want)
	}
}

func TestCategory(t *testing.T) {
	e := WordEntryValue{CategoryVal: "test"}

	got := e.Category()
	want := e.CategoryVal

	if got != want {
		t.Errorf("Category returned wrong value. Got %v. Want %v.", got, want)
	}
}

func TestPronunciation(t *testing.T) {
	e := DictionaryEntryValue{PronunciationVal: "test"}

	got := e.Pronunciation()
	want := e.PronunciationVal

	if got != want {
		t.Errorf("Pronunciation returned wrong value. Got %v. Want %v.", got, want)
	}
}

func TestSenses(t *testing.T) {
	senses := []SenseValue{
		{},
	}
	e := DictionaryEntryValue{SenseVals: senses}

	for i, sense := range e.Senses() {
		got := sense
		want := senses[i]

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Senses returned wrong value. Got %v. Want %v.", got, want)
		}
	}
}

func TestEtymologies(t *testing.T) {
	etymologies := []string{
		"test",
	}
	e := EtymologyEntryValue{EtymologyVals: etymologies}

	for i, etymology := range e.Etymologies() {
		got := etymology
		want := etymologies[i]

		if got != want {
			t.Errorf("Etymologies returned wrong value. Got %v. Want %v.", got, want)
		}
	}
}

func TestSynonyms(t *testing.T) {
	synonyms := []string{
		"test",
	}
	e := ThesaurusEntryValue{SynonymVals: synonyms}

	for i, synonym := range e.Synonyms() {
		got := synonym
		want := synonyms[i]

		if got != want {
			t.Errorf("Synonyms returned wrong value. Got %v. Want %v.", got, want)
		}
	}
}

func TestAntonyms(t *testing.T) {
	antonyms := []string{
		"test",
	}
	e := ThesaurusEntryValue{AntonymVals: antonyms}

	for i, antonym := range e.Antonyms() {
		got := antonym
		want := antonyms[i]

		if got != want {
			t.Errorf("Antonyms returned wrong value. Got %v. Want %v.", got, want)
		}
	}
}

func TestDefinitions(t *testing.T) {
	definitions := []string{
		"test",
	}
	s := SenseValue{DefinitionVals: definitions}

	for i, definition := range s.Definitions() {
		got := definition
		want := definitions[i]

		if got != want {
			t.Errorf("Definitions returned wrong value. Got %v. Want %v.", got, want)
		}
	}
}

func TestExamples(t *testing.T) {
	examples := []string{
		"test",
	}
	s := SenseValue{ExampleVals: examples}

	for i, example := range s.Examples() {
		got := example
		want := examples[i]

		if got != want {
			t.Errorf("Examples returned wrong value. Got %v. Want %v.", got, want)
		}
	}
}

func TestNotes(t *testing.T) {
	notes := []string{
		"test",
	}
	s := SenseValue{NoteVals: notes}

	for i, note := range s.Notes() {
		got := note
		want := notes[i]

		if got != want {
			t.Errorf("Notes returned wrong value. Got %v. Want %v.", got, want)
		}
	}
}

func TestSubsenses(t *testing.T) {
	senses := []SenseValue{
		{},
	}
	s := SenseValue{SubsenseVals: senses}

	for i, sense := range s.Subsenses() {
		got := sense
		want := senses[i]

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Subsenses returned wrong value. Got %v. Want %v.", got, want)
		}
	}
}
