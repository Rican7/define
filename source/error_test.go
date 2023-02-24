// Copyright © 2018 Trevor N. Suarez (Rican7)

package source

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// Enforce interface contracts
var (
	_ error = (*EmptyResultError)(nil)
	_ error = (*InvalidResponseError)(nil)
)

func TestValidateDictionaryResults(t *testing.T) {
	testData := []struct {
		word   string
		result DictionaryResults
		want   error
	}{
		{word: "", result: nil, want: &EmptyResultError{}},
		{word: "", result: DictionaryResults{}, want: &EmptyResultError{}},
		{word: "test", result: DictionaryResults{}, want: &EmptyResultError{Word: "test"}},
		{word: "test", result: DictionaryResults{{Language: "test"}}, want: nil},
	}

	for _, tt := range testData {
		if got := ValidateDictionaryResults(tt.word, tt.result); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ValidateDictionaryResults returned wrong value. Got %#v. Want %#v.", got, tt.want)
		}
	}
}

func TestValidateAndReturnDictionaryResults(t *testing.T) {
	testData := []struct {
		word    string
		result  DictionaryResults
		wantErr bool
	}{
		{word: "", result: nil, wantErr: true},
		{word: "", result: DictionaryResults{}, wantErr: true},
		{word: "test", result: DictionaryResults{}, wantErr: true},
		{word: "test", result: DictionaryResults{{Language: "test"}}, wantErr: false},
	}

	for _, tt := range testData {
		got, err := ValidateAndReturnDictionaryResults(tt.word, tt.result)

		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateAndReturnDictionaryResults returned an error when not expected. Got %#v.", err)
		}

		if !tt.wantErr && !reflect.DeepEqual(got, tt.result) {
			t.Errorf("ValidateAndReturnDictionaryResults returned wrong value. Got %#v. Want %#v.", got, tt.result)
		}
	}
}

func TestValidateSearchResults(t *testing.T) {
	testData := []struct {
		word   string
		result SearchResults
		want   error
	}{
		{word: "", result: nil, want: &EmptyResultError{}},
		{word: "", result: SearchResults{}, want: &EmptyResultError{}},
		{word: "test", result: SearchResults{}, want: &EmptyResultError{Word: "test"}},
		{word: "test", result: SearchResults{"test"}, want: nil},
	}

	for _, tt := range testData {
		if got := ValidateSearchResults(tt.word, tt.result); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ValidateSearchResults returned wrong value. Got %#v. Want %#v.", got, tt.want)
		}
	}
}

func TestValidateAndReturnSearchResults(t *testing.T) {
	testData := []struct {
		word    string
		result  SearchResults
		wantErr bool
	}{
		{word: "", result: nil, wantErr: true},
		{word: "", result: SearchResults{}, wantErr: true},
		{word: "test", result: SearchResults{}, wantErr: true},
		{word: "test", result: SearchResults{"test"}, wantErr: false},
	}

	for _, tt := range testData {
		got, err := ValidateAndReturnSearchResults(tt.word, tt.result)

		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateAndReturnSearchResults returned an error when not expected. Got %#v.", err)
		}

		if !tt.wantErr && !reflect.DeepEqual(got, tt.result) {
			t.Errorf("ValidateAndReturnSearchResults returned wrong value. Got %#v. Want %#v.", got, tt.result)
		}
	}
}

func TestValidateHTTPResponse(t *testing.T) {
	testData := []struct {
		httpResponse     *http.Response
		validMIMETypes   []string
		validStatusCodes []int
		wantErr          bool
	}{
		{httpResponse: nil, wantErr: true},
		{httpResponse: &http.Response{StatusCode: 400}, wantErr: true},
		{httpResponse: &http.Response{StatusCode: 500}, wantErr: true},
		{httpResponse: &http.Response{StatusCode: 500}, validStatusCodes: []int{500}, wantErr: false},
		{
			httpResponse: &http.Response{
				StatusCode: 200,
				Header:     http.Header{contentTypeHeaderName: []string{"application/test;charset=UTF-8"}},
			},
			validMIMETypes: []string{"application/test", "foo"},
			wantErr:        false,
		},
		{
			httpResponse: &http.Response{
				StatusCode: 500,
				Header:     http.Header{contentTypeHeaderName: []string{"application/test;charset=UTF-8"}},
			},
			validMIMETypes:   []string{"application/test", "foo"},
			validStatusCodes: []int{500},
			wantErr:          false,
		},
		{httpResponse: &http.Response{StatusCode: 200}, wantErr: false},
	}

	expectedErrType := reflect.TypeOf(&InvalidResponseError{}).Elem().Name()

	for _, tt := range testData {
		err := ValidateHTTPResponse(tt.httpResponse, tt.validMIMETypes, tt.validStatusCodes)
		invalidRespErr, ok := err.(*InvalidResponseError)

		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateHTTPResponse returned an error when not expected. Got %#v.", err)
		}

		if err != nil && !ok {
			errType := reflect.TypeOf(err).Elem().Name()

			t.Errorf("ValidateHTTPResponse returned an unexpected error type. Got %q. Want %q.", errType, expectedErrType)
		}

		if tt.wantErr && invalidRespErr.httpResponse != tt.httpResponse {
			t.Errorf("ValidateHTTPResponse returned wrong value. Got %#v. Want %#v.", err, tt.httpResponse)
		}
	}
}

func TestEmptyResultError_Error(t *testing.T) {
	word := "test"
	msg := (&EmptyResultError{Word: word}).Error()

	if msg == "" {
		t.Errorf("Error returned an empty message")
	}

	if !strings.Contains(msg, word) {
		t.Errorf("Error message %q didn't contain word %q", msg, word)
	}
}

func TestAuthenticationError_Error(t *testing.T) {
	msg := (&AuthenticationError{}).Error()

	if msg == "" {
		t.Errorf("Error returned an empty message")
	}
}

func TestInvalidResponseError_Error(t *testing.T) {
	msg := (&InvalidResponseError{}).Error()

	if msg == "" {
		t.Errorf("Error returned an empty message")
	}
}
