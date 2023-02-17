// Copyright © 2018 Trevor N. Suarez (Rican7)

package source

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	emptyResultErrorMessage         = "the source returned an empty result"
	authenticationErrorMessage      = "the source returned an authentication error"
	invalidResponseErrorMessage     = "the source returned an invalid response"
	errorMessageForWordSuffixFormat = " for word: %q"

	contentTypeHeaderName = "Content-Type"
)

var acceptableStatusCodes = []int{http.StatusOK}

// EmptyResultError represents an error caused by an empty result
type EmptyResultError struct {
	Word string
}

// AuthenticationError represents an error caused by an authentication problem
type AuthenticationError struct {
}

// InvalidResponseError represents an error caused by an invalid response
type InvalidResponseError struct {
	httpResponse *http.Response
}

// ValidateDictionaryResults validates the results of a define operation and
// returns an error if they're invalid
func ValidateDictionaryResults(word string, results []DictionaryResult) error {
	if len(results) < 1 {
		return &EmptyResultError{word}
	}

	return nil
}

// ValidateAndReturnDictionaryResults validates the results of a define
// operation and returns the results and a nil error if valid. If invalid, it'll
// return nil results and an error.
func ValidateAndReturnDictionaryResults(word string, results []DictionaryResult) ([]DictionaryResult, error) {
	if err := ValidateDictionaryResults(word, results); err != nil {
		return nil, err
	}

	return results, nil
}

// ValidateHTTPResponse validates an HTTP response and returns an error if the
// response is invalid
func ValidateHTTPResponse(httpResponse *http.Response, validContentTypes []string, validStatusCodes []int) error {
	if httpResponse == nil {
		return &InvalidResponseError{}
	}

	validStatusCodes = append(acceptableStatusCodes, validStatusCodes...)

	isValidStatusCode := false

	// Check if the HTTP response code is valid
	for _, validStatusCode := range validStatusCodes {
		if validStatusCode == httpResponse.StatusCode {
			isValidStatusCode = true
			break
		}
	}

	contentType := strings.ToLower(httpResponse.Header.Get(contentTypeHeaderName))
	isValidContentType := false

	if len(validContentTypes) < 1 {
		isValidContentType = true
	}

	// Check if the HTTP content-type is valid
	for _, validContentType := range validContentTypes {
		if strings.Contains(contentType, strings.ToLower(validContentType)) {
			isValidContentType = true
			break
		}
	}

	if !isValidStatusCode || !isValidContentType {
		return &InvalidResponseError{httpResponse}
	}

	return nil
}

func (e *EmptyResultError) Error() string {
	msg := emptyResultErrorMessage

	if e.Word != "" {
		msg = msg + fmt.Sprintf(errorMessageForWordSuffixFormat, e.Word)
	}

	return msg
}

func (e *AuthenticationError) Error() string {
	return authenticationErrorMessage
}

func (e *InvalidResponseError) Error() string {
	return invalidResponseErrorMessage
}
