// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package oxford provides a dictionary source via the Oxford Dictionaries API
package oxford

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Rican7/define/source"
)

// Name defines the name of the source
const Name = "Oxford Dictionaries API"

const (
	// baseURLString is the base URL for all Oxford API interactions
	baseURLString = "https://od-api.oxforddictionaries.com/api/v1/"

	entriesURLString = baseURLString + "entries/"

	httpRequestAcceptHeaderName = "Accept"
	httpRequestAppIDHeaderName  = "app_id"
	httpRequestAppKeyHeaderName = "app_key"

	jsonMIMEType = "application/json"

	phoneticNotationIPAIdentifier = "IPA"
)

// apiURL is the URL instance used for Oxford API calls
var apiURL *url.URL

// validMIMETypes is the list of valid response MIME types
var validMIMETypes = []string{jsonMIMEType}

// api is a struct containing a configured HTTP client for Oxford API operations
type api struct {
	httpClient *http.Client
	appID      string
	appKey     string
}

// apiResult is a struct that defines the data structure for Oxford API results
type apiResult struct {
	Metadata struct {
		Provider string
	}
	Results []struct {
		ID             string
		Language       string
		LexicalEntries []struct {
			DerivativeOf []struct {
				Domains   []string
				ID        string
				Language  string
				Regions   []string
				Registers []string
				Text      string
			}
			Derivatives []struct {
				Domains   []string
				ID        string
				Language  string
				Regions   []string
				Registers []string
				Text      string
			}
			Entries []struct {
				Etymologies         []string
				GrammaticalFeatures []struct {
					Text string
					Type string
				}
				HomographNumber string
				Notes           []struct {
					ID   string
					Text string
					Type string
				}
				Pronunciations []struct {
					AudioFile        string
					Dialects         []string
					PhoneticNotation string
					PhoneticSpelling string
					Regions          []string
				}
				Senses       []apiSense
				VariantForms []struct {
					Regions []string
					Text    string
				}
			}
			GrammaticalFeatures []struct {
				Text string
				Type string
			}
			Language        string
			LexicalCategory string
			Notes           []struct {
				ID   string
				Text string
				Type string
			}
			Pronunciations []struct {
				AudioFile        string
				Dialects         []string
				PhoneticNotation string
				PhoneticSpelling string
				Regions          []string
			}
			Text         string
			VariantForms []struct {
				Regions []string
				Text    string
			}
		}
		Pronunciations []struct {
			AudioFile        string
			Dialects         []string
			PhoneticNotation string
			PhoneticSpelling string
			Regions          []string
		}
		Type string
		Word string
	}
}

// apiSense is a struct that defines the data structure for Oxford API senses
type apiSense struct {
	CrossReferenceMarkers []string
	CrossReferences       []struct {
		ID   string
		Text string
		Type string
	}
	Definitions []string
	Domains     []string
	Examples    []struct {
		Definitions []string
		Domains     []string
		Notes       []struct {
			ID   string
			Text string
			Type string
		}
		Regions      []string
		Registers    []string
		SenseIds     []string
		Text         string
		Translations []struct {
			Domains             []string
			GrammaticalFeatures []struct {
				Text string
				Type string
			}
			Language string
			Notes    []struct {
				ID   string
				Text string
				Type string
			}
			Regions   []string
			Registers []string
			Text      string
		}
	}
	ID    string
	Notes []struct {
		ID   string
		Text string
		Type string
	}
	Pronunciations []struct {
		AudioFile        string
		Dialects         []string
		PhoneticNotation string
		PhoneticSpelling string
		Regions          []string
	}
	Regions      []string
	Registers    []string
	Subsenses    []apiSense
	Translations []struct {
		Domains             []string
		GrammaticalFeatures []struct {
			Text string
			Type string
		}
		Language string
		Notes    []struct {
			ID   string
			Text string
			Type string
		}
		Regions   []string
		Registers []string
		Text      string
	}
	VariantForms []struct {
		Regions []string
		Text    string
	}
}

// oxfordEntry is a struct that contains the entry types for this API
type oxfordEntry struct {
	source.WordEntryValue
	source.DictionaryEntryValue
	source.EtymologyEntryValue
}

// Initialize the package
func init() {
	var err error

	apiURL, err = url.Parse(baseURLString)

	if err != nil {
		panic(err)
	}
}

// New returns a new Oxford API dictionary source
func New(httpClient http.Client, appID, appKey string) source.Source {
	return &api{&httpClient, appID, appKey}
}

// Name returns the name of the source
func (g *api) Name() string {
	return Name
}

// Define takes a word string and returns a dictionary source.Result
func (g *api) Define(word string) (source.Result, error) {
	// Prepare our URL
	requestURL, err := url.Parse(entriesURLString + "en/" + word)

	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodGet, apiURL.ResolveReference(requestURL).String(), nil)

	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(httpRequestAcceptHeaderName, jsonMIMEType)
	httpRequest.Header.Set(httpRequestAppIDHeaderName, g.appID)
	httpRequest.Header.Set(httpRequestAppKeyHeaderName, g.appKey)

	httpResponse, err := g.httpClient.Do(httpRequest)

	if err != nil {
		return nil, err
	}

	defer httpResponse.Body.Close()

	if http.StatusNotFound == httpResponse.StatusCode {
		return nil, &source.EmptyResultError{Word: word}
	}

	if http.StatusForbidden == httpResponse.StatusCode {
		return nil, &source.AuthenticationError{}
	}

	if err = source.ValidateHTTPResponse(httpResponse, validMIMETypes, nil); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(httpResponse.Body)

	if err != nil {
		return nil, err
	}

	var result apiResult

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if len(result.Results) < 1 {
		return nil, &source.EmptyResultError{Word: word}
	}

	return source.ValidateAndReturnResult(result.toResult())
}

// toResult converts the proprietary API result to a generic source.Result
func (r apiResult) toResult() source.Result {
	mainResult := r.Results[0]

	entries := make([]interface{}, len(mainResult.LexicalEntries))

	for i, lexicalEntry := range mainResult.LexicalEntries {
		entry := oxfordEntry{}

		for _, pronunciation := range lexicalEntry.Pronunciations {
			if strings.EqualFold(phoneticNotationIPAIdentifier, pronunciation.PhoneticNotation) {
				entry.PronunciationVal = pronunciation.PhoneticSpelling
			}
		}

		entry.WordVal = lexicalEntry.Text
		entry.CategoryVal = lexicalEntry.LexicalCategory

		for _, subEntry := range lexicalEntry.Entries {
			entry.EtymologyVals = append(entry.EtymologyVals, subEntry.Etymologies...)

			for _, sense := range subEntry.Senses {
				senseVal := sense.toSenseValue()

				// Only go one level deep of sub-senses
				for _, subSense := range sense.Subsenses {
					senseVal.SubsenseVals = append(senseVal.SubsenseVals, subSense.toSenseValue())
				}

				entry.SenseVals = append(entry.SenseVals, senseVal)
			}
		}

		entries[i] = entry
	}

	return source.ResultValue{
		Head:      mainResult.Word,
		Lang:      mainResult.Language,
		EntryVals: entries,
	}
}

// toSenseValue converts the proprietary API sense to a source.SenseValue
func (s apiSense) toSenseValue() source.SenseValue {
	examples := make([]string, len(s.Examples))
	notes := make([]string, len(s.Notes))

	for i, example := range s.Examples {
		examples[i] = example.Text
	}

	for i, note := range s.Notes {
		notes[i] = note.Text
	}

	return source.SenseValue{
		DefinitionVals: s.Definitions,
		ExampleVals:    examples,
		NoteVals:       notes,
	}
}
