// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package oxford provides a dictionary source via the Oxford Dictionaries API
package oxford

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Rican7/define/source"
)

// Name defines the name of the source
const Name = "Oxford Dictionaries API"

const (
	// baseURLString is the base URL for all Oxford API interactions
	baseURLString = "https://od-api.oxforddictionaries.com/api/v2/"

	entriesURLString = baseURLString + "entries/"
	searchURLString  = baseURLString + "search/"

	httpRequestAcceptHeaderName           = "Accept"
	httpRequestAppIDHeaderName            = "app_id"
	httpRequestAppKeyHeaderName           = "app_key"
	httpRequestSearchStringQueryParamName = "q"
	httpRequestLimitQueryParamName        = "limit"

	jsonMIMEType = "application/json"

	fallbackSearchResultLimit = 10

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

// Name returns the printable, human-readable name of the source.
func (a *api) Name() string {
	return Name
}

// Define takes a word string and returns a list of dictionary results, and
// an error if any occurred.
func (a *api) Define(word string) (source.DictionaryResults, error) {
	// Prepare our URL
	requestURL, err := url.Parse(entriesURLString + "en-us/" + word)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodGet, apiURL.ResolveReference(requestURL).String(), nil)
	if err != nil {
		return nil, err
	}

	a.signRequest(httpRequest)

	httpResponse, err := a.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	defer httpResponse.Body.Close()

	if err = validateResponse(word, httpResponse); err != nil {
		if _, isEmptyResult := err.(*source.EmptyResultError); isEmptyResult {
			return a.apiSearchFallback(word)
		}

		return nil, err
	}

	var response apiDefinitionResponse

	if err = decodeResponseData(httpResponse.Body, &response); err != nil {
		return nil, err
	}

	if len(response.Results) < 1 {
		return a.apiSearchFallback(word)
	}

	return source.ValidateAndReturnDictionaryResults(word, response.toResults())
}

// Search takes a word string and returns a list of found words, and an
// error if any occurred.
func (a *api) Search(word string, limit uint) (source.SearchResults, error) {
	response, err := a.apiSearch(word, limit)
	if err != nil {
		return nil, err
	}

	results := response.toResults()

	if limit > 1 && limit < uint(len(results)) {
		results = results[:limit]
	}

	return source.ValidateAndReturnSearchResults(word, results)
}

func (a *api) apiSearch(word string, limit uint) (*apiSearchResponse, error) {
	// Prepare our URL
	requestURL, err := url.Parse(searchURLString + "en-us")

	queryParams := apiURL.Query()
	queryParams.Set(httpRequestSearchStringQueryParamName, word)

	if limit > 0 {
		queryParams.Set(httpRequestLimitQueryParamName, strconv.FormatUint(uint64(limit), 10))
	}

	requestURL.RawQuery = queryParams.Encode()

	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodGet, apiURL.ResolveReference(requestURL).String(), nil)
	if err != nil {
		return nil, err
	}

	a.signRequest(httpRequest)

	httpResponse, err := a.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	defer httpResponse.Body.Close()

	if err = validateResponse(word, httpResponse); err != nil {
		return nil, err
	}

	var response apiSearchResponse

	if err = decodeResponseData(httpResponse.Body, &response); err != nil {
		return nil, err
	}

	if len(response.Results) < 1 {
		return nil, &source.EmptyResultError{Word: word}
	}

	return &response, nil
}

func (a *api) apiSearchFallback(word string) (source.DictionaryResults, error) {
	response, err := a.apiSearch(word, fallbackSearchResultLimit)
	if err != nil {
		return nil, err
	}

	if len(response.Results) < 1 {
		return nil, &source.EmptyResultError{Word: word}
	}

	var fallbackWord string

	for _, apiSearchResult := range response.Results {
		// Look for our first inflection result
		if apiSearchResult.MatchType == apiSearchResultMatchTypeInflection {
			fallbackWord = apiSearchResult.Label
			break
		}
	}

	if strings.EqualFold(word, fallbackWord) {
		// Prevent matching against the same word
		return nil, &source.EmptyResultError{Word: word}
	}

	return a.Define(fallbackWord)
}

func (a *api) signRequest(request *http.Request) {
	request.Header.Set(httpRequestAcceptHeaderName, jsonMIMEType)
	request.Header.Set(httpRequestAppIDHeaderName, a.appID)
	request.Header.Set(httpRequestAppKeyHeaderName, a.appKey)
}

func validateResponse(word string, response *http.Response) error {
	switch response.StatusCode {
	case http.StatusNotFound:
		return &source.EmptyResultError{Word: word}
	case http.StatusForbidden:
		return &source.AuthenticationError{}
	}

	if err := source.ValidateHTTPResponse(response, validMIMETypes, nil); err != nil {
		return err
	}

	return nil
}

func decodeResponseData(data io.Reader, into any) error {
	body, err := io.ReadAll(data)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, into); err != nil {
		return err
	}

	return nil
}
