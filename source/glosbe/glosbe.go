// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package glosbe provides a dictionary source via the Glosbe API
package glosbe

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"html"

	"github.com/Rican7/define/source"
	"github.com/microcosm-cc/bluemonday"
)

// Name defines the name of the source
const Name = "Glosbe API"

const (
	// baseURLString is the base URL for all Glosbe API interactions
	baseURLString = "https://glosbe.com/gapi/translate?format=json&from=en&dest=en"

	// wordParameter defines the HTTP parameter for the word to define
	wordParameter = "phrase"

	httpRequestAcceptHeaderName = "Accept"
	jsonMIMEType                = "application/json"
)

// apiURL is the URL instance used for Glosbe API calls
var apiURL *url.URL

// validMIMETypes is the list of valid response MIME types
var validMIMETypes = []string{jsonMIMEType}

// htmlCleaner is used to clean the strings returned from the API
var htmlCleaner = bluemonday.StrictPolicy()

// stringCleaner is used to clean the strings returned from the API
var stringCleaner *strings.Replacer

// itemsToClean is a list of itemsto clean/remove from the strings returned
// from the API
var itemsToClean = []string{
	"[i]",
	"[/i]",
}

// api is a struct containing a configured HTTP client for Glosbe API operations
type api struct {
	httpClient *http.Client
}

// apiResponse defines the data structure for a Glosbe API response
type apiResponse struct {
	Result string
	TUC    []*struct {
		Meanings []*struct {
			Language string
			Text     string
		}
		Phrase *struct {
			Language string
			Text     string
		}
		Authors []int
	}
	Phrase string
	Dest   string
}

// Initialize the package
func init() {
	var err error

	apiURL, err = url.Parse(baseURLString)

	if err != nil {
		panic(err)
	}

	stringCleanerPairs := make([]string, len(itemsToClean)*2)

	for _, items := range itemsToClean {
		stringCleanerPairs = append(stringCleanerPairs, []string{items, ""}...)
	}

	stringCleaner = strings.NewReplacer(stringCleanerPairs...)
}

// New returns a new Glosbe API dictionary source
func New(httpClient http.Client) source.Source {
	return &api{&httpClient}
}

// Name returns the printable, human-readable name of the source.
func (g *api) Name() string {
	return Name
}

// Define takes a word string and returns a list of dictionary results, and
// an error if any occurred.
func (g *api) Define(word string) ([]source.DictionaryResult, error) {
	// Prepare our URL
	queryParams := apiURL.Query()
	queryParams.Set(wordParameter, word)
	apiURL.RawQuery = queryParams.Encode()

	httpRequest, err := http.NewRequest(http.MethodGet, apiURL.String(), nil)

	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(httpRequestAcceptHeaderName, jsonMIMEType)

	httpResponse, err := g.httpClient.Do(httpRequest)

	if err != nil {
		return nil, err
	}

	if err = source.ValidateHTTPResponse(httpResponse, validMIMETypes, nil); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(httpResponse.Body)

	if err != nil {
		return nil, err
	}

	var response apiResponse

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response.TUC) < 1 {
		return nil, &source.EmptyResultError{Word: word}
	}

	return source.ValidateAndReturnDictionaryResults(word, response.toResults())
}

// toResult converts the API response to the results that a source expects to
// return.
func (r *apiResponse) toResults() []source.DictionaryResult {
	entry := source.DictionaryEntry{}

	senses := make([]source.Sense, 0)

	for _, item := range r.TUC {
		// Entries are only valid definitions if they don't have a separate
		// phrase, or their phrase matches the looked-up phrase
		if item.Phrase == nil || strings.EqualFold(item.Phrase.Text, r.Phrase) {
			for _, meaning := range item.Meanings {
				definition := sanitize(meaning.Text)

				sense := source.Sense{Definitions: []string{definition}}

				senses = append(senses, sense)
			}
		} else if item.Phrase != nil && item.Phrase.Text != "" {
			entry.Synonyms = append(entry.Synonyms, item.Phrase.Text)
		}
	}

	entry.Senses = senses

	return []source.DictionaryResult{
		{
			Language: r.Dest,
			Entries:  []source.DictionaryEntry{entry},
		},
	}
}

// sanitize cleans a string of any formatting identifiers or markup
func sanitize(str string) string {
	str = htmlCleaner.Sanitize(str)
	str = html.UnescapeString(str)
	str = stringCleaner.Replace(str)

	return str
}
