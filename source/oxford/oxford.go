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
	baseURLString = "https://od-api.oxforddictionaries.com/api/v2/"

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
	requestURL, err := url.Parse(entriesURLString + "en-us/" + word)

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

	var result apiResponse

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if len(result.Results) < 1 {
		return nil, &source.EmptyResultError{Word: word}
	}

	return source.ValidateAndReturnResult(result.toResult())
}

// toResult converts the proprietary API response to a generic source.Result
func (r apiResponse) toResult() source.Result {
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
		entry.CategoryVal = lexicalEntry.LexicalCategory.Text

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
