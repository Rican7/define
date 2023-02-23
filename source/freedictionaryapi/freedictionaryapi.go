// Copyright © 2023 Trevor N. Suarez (Rican7)

// Package freedictionaryapi provides a dictionary source via the "Free
// Dictionary API"
package freedictionaryapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/Rican7/define/source"
)

// Name defines the name of the source
const Name = "Free Dictionary API"

const (
	// baseURLString is the base URL for all Free Dictionary API interactions
	baseURLString = "https://api.dictionaryapi.dev/api/v2/"

	entriesURLString = baseURLString + "entries/"

	httpRequestAcceptHeaderName = "Accept"

	jsonMIMEType = "application/json"
)

// apiURL is the URL instance used for Free Dictionary API calls
var apiURL *url.URL

// validMIMETypes is the list of valid response MIME types
var validMIMETypes = []string{jsonMIMEType}

// api is a struct containing a configured HTTP client for Free Dictionary API operations
type api struct {
	httpClient *http.Client
}

// Initialize the package
func init() {
	var err error

	apiURL, err = url.Parse(baseURLString)

	if err != nil {
		panic(err)
	}
}

// New returns a new Free Dictionary API dictionary source
func New(httpClient http.Client) source.Source {
	return &api{&httpClient}
}

// Name returns the printable, human-readable name of the source.
func (a *api) Name() string {
	return Name
}

// Define takes a word string and returns a list of dictionary results, and
// an error if any occurred.
func (a *api) Define(word string) (source.DictionaryResults, error) {
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

	httpResponse, err := a.httpClient.Do(httpRequest)

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

	var response apiResponse

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response) < 1 {
		return nil, &source.EmptyResultError{Word: word}
	}

	return source.ValidateAndReturnDictionaryResults(word, response.toResults())
}
