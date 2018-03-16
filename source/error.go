package source

import (
	"fmt"
	"net/http"
)

const (
	emptyResultErrorMessage         = "The source returned an empty result"
	invalidResponseErrorMessage     = "The source returned an invalid response"
	errorMessageForWordSuffixFormat = " for word: %q"
)

var acceptableStatusCodes = []int{http.StatusOK}

// EmptyResultError represents an error caused by an empty result
type EmptyResultError struct {
	Word string
}

// InvalidResponseError represents an error caused by an invalid response
type InvalidResponseError struct {
	httpResponse http.Response
}

// ValidateHTTPResponse valides an HTTP response and returns an error if the
// response is invalid
func ValidateHTTPResponse(httpResponse *http.Response, validStatusCodes ...int) error {
	validStatusCodes = append(acceptableStatusCodes, validStatusCodes...)

	isValidStatusCode := false

	// Check if the HTTP response code is valid
	for _, statusCode := range validStatusCodes {
		if statusCode == httpResponse.StatusCode {
			isValidStatusCode = true
			break
		}
	}

	if !isValidStatusCode {
		return &InvalidResponseError{*httpResponse}
	}

	return nil
}

func (e *EmptyResultError) Error() string {
	msg := emptyResultErrorMessage

	if "" != e.Word {
		msg = msg + fmt.Sprintf(errorMessageForWordSuffixFormat, e.Word)
	}

	return msg
}

func (e *InvalidResponseError) Error() string {
	return invalidResponseErrorMessage
}
