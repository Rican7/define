// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package webster provides a dictionary source via the Webster Dictionaries API
package webster

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/Rican7/define/source"
)

// Name defines the name of the source
const Name = "Merriam-Webster's Dictionary API"

const (
	// baseURLString is the base URL for all Webster API interactions
	baseURLString = "https://www.dictionaryapi.com/api/v3/"

	entriesURLString = baseURLString + "references/collegiate/json/"

	httpRequestAcceptHeaderName  = "Accept"
	httpRequestKeyQueryParamName = "key"

	jsonMIMEType = "application/json"
)

// apiURL is the URL instance used for Webster API calls
var apiURL *url.URL

// validMIMETypes is the list of valid response MIME types
var validMIMETypes = []string{jsonMIMEType}

// api contains a configured HTTP client for Webster API operations
type api struct {
	httpClient *http.Client
	appKey     string
}

// Initialize the package
func init() {
	var err error

	apiURL, err = url.Parse(baseURLString)

	if err != nil {
		panic(err)
	}
}

// New returns a new Webster API dictionary source
func New(httpClient http.Client, appKey string) source.Source {
	return &api{&httpClient, appKey}
}

// Name returns the printable, human-readable name of the source.
func (g *api) Name() string {
	return Name
}

// Define takes a word string and returns a list of dictionary results, and
// an error if any occurred.
func (g *api) Define(word string) ([]source.DictionaryResult, error) {
	// Prepare our URL
	requestURL, err := url.Parse(entriesURLString + word)
	queryParams := apiURL.Query()
	queryParams.Set(httpRequestKeyQueryParamName, g.appKey)
	requestURL.RawQuery = queryParams.Encode()

	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodGet, apiURL.ResolveReference(requestURL).String(), nil)

	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(httpRequestAcceptHeaderName, jsonMIMEType)

	httpResponse, err := g.httpClient.Do(httpRequest)

	if err != nil {
		return nil, err
	}

	defer httpResponse.Body.Close()

	if err = source.ValidateHTTPResponse(httpResponse, validMIMETypes, nil); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(httpResponse.Body)

	if err != nil {
		return nil, err
	}

	var rawResponse apiRawResponse

	if err = json.Unmarshal(body, &rawResponse); err != nil {
		return nil, err
	}

	if len(rawResponse) < 1 {
		return nil, &source.EmptyResultError{Word: word}
	}

	switch rawResponse[0].(type) {
	case apiDefinitionResult:
		response := apiResponseFromRaw[apiDefinitionResult](rawResponse)
		results := apiDefinitionResults(response)

		return source.ValidateAndReturnDictionaryResults(word, results.toResults())
	case apiSearchResult:
		response := apiResponseFromRaw[apiSearchResult](rawResponse)
		_ = apiSearchResults(response)

		// TODO: Handle a fallback search with a "did you mean?" error of sorts
		return nil, &source.EmptyResultError{Word: word}
	}

	return nil, &source.EmptyResultError{Word: word}
}
