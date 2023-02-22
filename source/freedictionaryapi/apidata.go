// Copyright © 2023 Trevor N. Suarez (Rican7)

package freedictionaryapi

import (
	"strings"

	"github.com/Rican7/define/source"
)

const (
	// apiPhoneticsWrapper defines the character used to wrap phonetics strings
	// in the Free Dictionary API
	apiPhoneticsWrapper = '/'
)

// apiResponse defines the data structure for a Free Dictionary API response
type apiResponse []apiResult

// apiResult defines the data structure for a Free Dictionary API result
type apiResult struct {
	Word       string         `json:"word"`
	Phonetic   string         `json:"phonetic"`
	Phonetics  []apiPhonetics `json:"phonetics"`
	Meanings   []apiMeaning   `json:"meanings"`
	License    apiLicense     `json:"license"`
	SourceUrls []string       `json:"sourceUrls"`
}

// apiPhonetics defines the data structure for Free Dictionary API phonetics
type apiPhonetics struct {
	Text      string     `json:"text"`
	Audio     string     `json:"audio"`
	SourceURL string     `json:"sourceUrl"`
	License   apiLicense `json:"license"`
}

// apiPhonetics defines the data structure for Free Dictionary API phonetics
type apiMeaning struct {
	PartOfSpeech string          `json:"partOfSpeech"`
	Definitions  []apiDefinition `json:"definitions"`

	apiThesaurusValues
}

// apiDefinition defines the data structure for a Free Dictionary API definition
type apiDefinition struct {
	Definition string `json:"definition"`
	Example    string `json:"example"`

	apiThesaurusValues
}

// apiThesaurusValues defines the data structure for Free Dictionary API
// thesaurus values
type apiThesaurusValues struct {
	Synonyms []string `json:"synonyms"`
	Antonyms []string `json:"antonyms"`
}

// apiLicense defines the data structure for a Free Dictionary API license
type apiLicense struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// toResult converts the API response to the results that a source expects to
// return.
func (r apiResponse) toResults() []source.DictionaryResult {
	sourceResults := make([]source.DictionaryResult, 0, len(r))

	for _, result := range r {
		sourceEntries := make([]source.DictionaryEntry, 0, len(result.Meanings))

		var pronunciations []string
		if result.Phonetic != "" {
			pronunciation := cleanPhoneticText(result.Phonetic)

			pronunciations = append(pronunciations, pronunciation)
		}

		for _, phonetic := range result.Phonetics {
			if phonetic.Text == "" {
				continue
			}

			pronunciation := cleanPhoneticText(phonetic.Text)

			if len(pronunciations) < 1 || pronunciations[0] != pronunciation {
				pronunciations = append(pronunciations, pronunciation)
			}
		}

		for _, apiMeaning := range result.Meanings {
			sourceEntry := apiMeaning.toEntry()

			sourceEntry.Word = result.Word
			sourceEntry.Pronunciations = pronunciations

			sourceEntries = append(sourceEntries, sourceEntry)
		}

		sourceResults = append(
			sourceResults,
			source.DictionaryResult{
				Language: "en", // TODO
				Entries:  sourceEntries,
			},
		)
	}

	return sourceResults
}

// toEntry converts the API meaning to a source.DictionaryEntry
func (m *apiMeaning) toEntry() source.DictionaryEntry {
	sourceEntry := source.DictionaryEntry{}

	sourceEntry.LexicalCategory = m.PartOfSpeech

	for _, apiDefinition := range m.Definitions {
		sourceSense := apiDefinition.toSense()

		sourceEntry.Senses = append(sourceEntry.Senses, sourceSense)
	}

	sourceEntry.ThesaurusValues = m.apiThesaurusValues.toThesaurusValues()

	return sourceEntry
}

// toThesaurusValues converts API thesaurus values to a source.ThesaurusValues
func (v *apiThesaurusValues) toThesaurusValues() source.ThesaurusValues {
	return source.ThesaurusValues{
		Synonyms: v.Synonyms,
		Antonyms: v.Antonyms,
	}
}

// toSense converts the API definition to a source.Sense
func (d *apiDefinition) toSense() source.Sense {
	var definitions []string
	var examples []string

	if d.Definition != "" {
		definitions = []string{d.Definition}
	}

	if d.Example != "" {
		examples = []string{d.Example}
	}

	return source.Sense{
		Definitions: definitions,
		Examples:    examples,

		ThesaurusValues: d.apiThesaurusValues.toThesaurusValues(),
	}
}

func cleanPhoneticText(text string) string {
	return strings.Trim(text, string(apiPhoneticsWrapper))
}
