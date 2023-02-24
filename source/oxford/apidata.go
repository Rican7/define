package oxford

import (
	"sort"
	"strings"

	"github.com/Rican7/define/source"
)

const (
	apiSearchResultMatchTypeInflection = "inflection"
)

// apiDefinitionResponse defines the structure of an Oxford API define response
type apiDefinitionResponse struct {
	Metadata struct {
		Operation string `json:"operation"`
		Provider  string `json:"provider"`
		Schema    string `json:"schema"`
	} `json:"metadata"`
	Results []apiDefinitionResult `json:"results"`
}

// apiSearchResponse defines the structure of an Oxford API search response
type apiSearchResponse struct {
	Metadata struct {
		Limit          string `json:"limit"`
		Offset         string `json:"offset"`
		Operation      string `json:"operation"`
		Provider       string `json:"provider"`
		Schema         string `json:"schema"`
		SourceLanguage string `json:"sourceLanguage"`
		Total          string `json:"total"`
	} `json:"metadata"`
	Results []apiSearchResult `json:"results"`
}

// apiDefinitionResult defines the structure of an Oxford API definition result
type apiDefinitionResult struct {
	ID             string             `json:"id"`
	Language       string             `json:"language"`
	LexicalEntries []apiLexicalEntry  `json:"lexicalEntries"`
	Pronunciations []apiPronunciation `json:"pronunciations"`
	Type           string             `json:"type"`
	Word           string             `json:"word"`
}

// apiSearchResult defines the structure of an Oxford API search result
type apiSearchResult struct {
	ID          string  `json:"id"`
	Label       string  `json:"label"`
	MatchString string  `json:"matchString"`
	MatchType   string  `json:"matchType"`
	Region      string  `json:"region"`
	Score       float64 `json:"score"`
	Word        string  `json:"word"`
}

// apiLexicalEntry defines the structure of an Oxford API lexical entry
type apiLexicalEntry struct {
	Compounds    []apiWordReference `json:"compounds"`
	DerivativeOf []apiWordReference `json:"derivativeOf"`
	Derivatives  []apiWordReference `json:"derivatives"`
	Entries      []struct {
		CrossReferenceMarkers []string           `json:"crossReferenceMarkers"`
		CrossReferences       []apiTypedIDText   `json:"crossReferences"`
		Etymologies           []string           `json:"etymologies"`
		GrammaticalFeatures   []apiTypedIDText   `json:"grammaticalFeatures"`
		HomographNumber       string             `json:"homographNumber"`
		Inflections           []apiInflection    `json:"inflections"`
		Notes                 []apiTypedIDText   `json:"notes"`
		Pronunciations        []apiPronunciation `json:"pronunciations"`
		Senses                []apiSense         `json:"senses"`
		VariantForms          []apiVariantForm   `json:"variantForms"`
	} `json:"entries"`
	GrammaticalFeatures []apiTypedIDText   `json:"grammaticalFeatures"`
	Language            string             `json:"language"`
	LexicalCategory     apiIDText          `json:"lexicalCategory"`
	Notes               []apiTypedIDText   `json:"notes"`
	PhrasalVerbs        []apiWordReference `json:"phrasalVerbs"`
	Phrases             []apiWordReference `json:"phrases"`
	Pronunciations      []apiPronunciation `json:"pronunciations"`
	Root                string             `json:"root"`
	Text                string             `json:"text"`
	VariantForms        []apiVariantForm   `json:"variantForms"`
}

// apiWordReference defines the structure of an Oxford API word reference
type apiWordReference struct {
	Domains   []apiIDText `json:"domains"`
	ID        string      `json:"id"`
	Language  string      `json:"language"`
	Regions   []apiIDText `json:"regions"`
	Registers []apiIDText `json:"registers"`
	Text      string      `json:"text"`
}

// apiIDText defines the structure of an Oxford API text with ID
type apiIDText struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// apiTypedIDText defines the structure of an Oxford API typed, ID'd text
type apiTypedIDText struct {
	apiIDText

	Type string `json:"type"`
}

// apiInflection defines the structure of an Oxford API inflection
type apiInflection struct {
	Domains             []apiIDText        `json:"domains"`
	GrammaticalFeatures []apiTypedIDText   `json:"grammaticalFeatures"`
	InflectedForm       string             `json:"inflectedForm"`
	LexicalCategory     apiIDText          `json:"lexicalCategory"`
	Pronunciations      []apiPronunciation `json:"pronunciations"`
	Regions             []apiIDText        `json:"regions"`
	Registers           []apiIDText        `json:"registers"`
}

// apiVariantForm defines the structure of an Oxford API variant form
type apiVariantForm struct {
	Domains        []apiIDText        `json:"domains"`
	Notes          []apiTypedIDText   `json:"notes"`
	Pronunciations []apiPronunciation `json:"pronunciations"`
	Regions        []apiIDText        `json:"regions"`
	Registers      []apiIDText        `json:"registers"`
	Text           string             `json:"text"`
}

// apiSense defines the structure of an Oxford API "sense"
type apiSense struct {
	Antonyms      []apiWordReference `json:"antonyms"`
	Constructions []struct {
		Domains   []apiIDText      `json:"domains"`
		Examples  [][]string       `json:"examples"`
		Notes     []apiTypedIDText `json:"notes"`
		Regions   []apiIDText      `json:"regions"`
		Registers []apiIDText      `json:"registers"`
		Text      string           `json:"text"`
	} `json:"constructions"`
	CrossReferenceMarkers []string            `json:"crossReferenceMarkers"`
	CrossReferences       []apiTypedIDText    `json:"crossReferences"`
	Definitions           []string            `json:"definitions"`
	DomainClasses         []apiIDText         `json:"domainClasses"`
	Domains               []apiIDText         `json:"domains"`
	Etymologies           []string            `json:"etymologies"`
	Examples              []apiComplexExample `json:"examples"`
	ID                    string              `json:"id"`
	Inflections           []apiInflection     `json:"inflections"`
	Notes                 []apiTypedIDText    `json:"notes"`
	Pronunciations        []apiPronunciation  `json:"pronunciations"`
	Regions               []apiIDText         `json:"regions"`
	Registers             []apiIDText         `json:"registers"`
	SemanticClasses       []apiIDText         `json:"semanticClasses"`
	ShortDefinitions      []string            `json:"shortDefinitions"`
	Subsenses             []apiSense          `json:"subsenses"`
	Synonyms              []apiWordReference  `json:"synonyms"`
	ThesaurusLinks        []struct {
		EntryID string `json:"entry_id"`
		SenseID string `json:"sense_id"`
	} `json:"thesaurusLinks"`
	VariantForms []apiVariantForm `json:"variantForms"`
}

// apiComplexExample defines the structure of an Oxford API "example"
type apiComplexExample struct {
	Definitions []string         `json:"definitions"`
	Domains     []apiIDText      `json:"domains"`
	Notes       []apiTypedIDText `json:"notes"`
	Regions     []apiIDText      `json:"regions"`
	Registers   []apiIDText      `json:"registers"`
	SenseIds    []string         `json:"senseIds"`
	Text        string           `json:"text"`
}

// apiPronunciation defines the structure of an Oxford API "pronunciation"
type apiPronunciation struct {
	AudioFile        string      `json:"audioFile"`
	Dialects         []string    `json:"dialects"`
	PhoneticNotation string      `json:"phoneticNotation"`
	PhoneticSpelling string      `json:"phoneticSpelling"`
	Regions          []apiIDText `json:"regions"`
	Registers        []apiIDText `json:"registers"`
}

// toResults converts the API response to the results that a source expects to
// return.
func (r *apiDefinitionResponse) toResults() source.DictionaryResults {
	sourceResults := make(source.DictionaryResults, 0, len(r.Results))

	for _, result := range r.Results {
		word := result.Word
		sourceEntries := make([]source.DictionaryEntry, 0, len(result.LexicalEntries))

		for _, lexicalEntry := range result.LexicalEntries {
			sourceEntry := lexicalEntry.toEntry()

			if word == "" {
				word = sourceEntry.Word
			}

			sourceEntries = append(sourceEntries, sourceEntry)
		}

		sourceResults = append(
			sourceResults,
			source.DictionaryResult{
				Language: result.Language,
				Word:     word,
				Entries:  sourceEntries,
			},
		)
	}

	return sourceResults
}

// toResults converts the API response to the results that a source expects to
// return.
func (r *apiSearchResponse) toResults() source.SearchResults {
	apiResults := r.Results
	sourceResults := make(source.SearchResults, 0, len(r.Results))

	// Sort the results by score
	sort.Slice(
		apiResults,
		func(i, j int) bool {
			return apiResults[i].Score < apiResults[i].Score
		},
	)

	for _, apiResult := range apiResults {
		sourceResult := source.SearchResult(apiResult.Label)

		sourceResults = append(sourceResults, sourceResult)
	}

	return sourceResults
}

// toEntry converts the API lexical entry to a source.DictionaryEntry
func (e *apiLexicalEntry) toEntry() source.DictionaryEntry {
	sourceEntry := source.DictionaryEntry{}

	for _, pronunciation := range e.Pronunciations {
		if strings.EqualFold(phoneticNotationIPAIdentifier, pronunciation.PhoneticNotation) {
			sourceEntry.Pronunciations = append(sourceEntry.Pronunciations, source.Pronunciation(pronunciation.PhoneticSpelling))
		}
	}

	sourceEntry.Word = e.Text
	sourceEntry.LexicalCategory = e.LexicalCategory.Text

	for _, subEntry := range e.Entries {
		sourceEntry.Etymologies = append(sourceEntry.Etymologies, subEntry.Etymologies...)

		for _, pronunciation := range subEntry.Pronunciations {
			if strings.EqualFold(phoneticNotationIPAIdentifier, pronunciation.PhoneticNotation) {
				sourceEntry.Pronunciations = append(sourceEntry.Pronunciations, source.Pronunciation(pronunciation.PhoneticSpelling))
			}
		}

		for _, sense := range subEntry.Senses {
			sourceSense := sense.toSense()

			// Only go one level deep of sub-senses
			for _, subSense := range sense.Subsenses {
				sourceSense.SubSenses = append(sourceSense.SubSenses, subSense.toSense())
			}

			sourceEntry.Senses = append(sourceEntry.Senses, sourceSense)
		}
	}

	return sourceEntry
}

// toSense converts the API sense to a source.Sense
func (s *apiSense) toSense() source.Sense {
	examples := make([]source.AttributedText, 0, len(s.Examples))
	notes := make([]string, 0, len(s.Notes))

	for _, example := range s.Examples {
		examples = append(examples, example.toAttributedText())
	}

	for _, note := range s.Notes {
		notes = append(notes, note.Text)
	}

	return source.Sense{
		Definitions: s.Definitions,
		Examples:    examples,
		Notes:       notes,
	}
}

// toAttributedText converts the API example to a source.AttributedText
func (e *apiComplexExample) toAttributedText() source.AttributedText {
	return source.AttributedText{
		Text: e.Text,
	}
}
