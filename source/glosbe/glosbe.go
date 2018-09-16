// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package glosbe provides a dictionary source via the Glosbe API
package glosbe

import (
	"encoding/json"
	"io/ioutil"
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

// apiResult is a struct that defines the data structure for Glosbe API results
type apiResult struct {
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

// glosbeEntry is a struct that contains the entry types for this API
type glosbeEntry struct {
	source.DictionaryEntryValue
	source.ThesaurusEntryValue
}

// Initialize the package
func init() {
	var err error

	apiURL, err = url.Parse(baseURLString)

	if nil != err {
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

// Name returns the name of the source
func (g *api) Name() string {
	return Name
}

// Define takes a word string and returns a dictionary source.Result
func (g *api) Define(word string) (source.Result, error) {
	// Prepare our URL
	queryParams := apiURL.Query()
	queryParams.Set(wordParameter, word)
	apiURL.RawQuery = queryParams.Encode()

	httpRequest, err := http.NewRequest(http.MethodGet, apiURL.String(), nil)

	if nil != err {
		return nil, err
	}

	httpRequest.Header.Set(httpRequestAcceptHeaderName, jsonMIMEType)

	httpResponse, err := g.httpClient.Do(httpRequest)

	if nil != err {
		return nil, err
	}

	if err = source.ValidateHTTPResponse(httpResponse, validMIMETypes, nil); nil != err {
		return nil, err
	}

	body, err := ioutil.ReadAll(httpResponse.Body)

	if nil != err {
		return nil, err
	}

	var result apiResult

	if err = json.Unmarshal(body, &result); nil != err {
		return nil, err
	}

	if len(result.TUC) < 1 {
		return nil, &source.EmptyResultError{Word: word}
	}

	return source.ValidateAndReturnResult(result.toResult())
}

// toResult converts the proprietary API result to a generic source.Result
func (r apiResult) toResult() source.Result {
	entry := glosbeEntry{
		source.DictionaryEntryValue{},
		source.ThesaurusEntryValue{},
	}

	senses := make([]source.SenseValue, 0)

	for _, item := range r.TUC {
		// Entries are only valid definitions if they don't have a separate
		// phrase, or their phrase matches the looked-up phrase
		if nil == item.Phrase || strings.EqualFold(item.Phrase.Text, r.Phrase) {
			for _, meaning := range item.Meanings {
				definition := sanitize(meaning.Text)

				sense := source.SenseValue{DefinitionVals: []string{definition}}

				senses = append(senses, sense)
			}
		} else if nil != item.Phrase && "" != item.Phrase.Text {
			entry.SynonymVals = append(entry.SynonymVals, item.Phrase.Text)
		}
	}

	entry.SenseVals = senses

	return source.ResultValue{
		Head:      r.Phrase,
		Lang:      r.Dest,
		EntryVals: []interface{}{entry},
	}
}

// sanitize cleans a string of any formatting identifiers or markup
func sanitize(str string) string {
	str = htmlCleaner.Sanitize(str)
	str = html.UnescapeString(str)
	str = stringCleaner.Replace(str)

	return str
}
