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

func TestValidateResult(t *testing.T) {
	testData := []struct {
		result Result
		want   error
	}{
		{result: nil, want: &EmptyResultError{}},
		{result: ResultValue{EntryVals: []interface{}{DictionaryEntryValue{}}}, want: &EmptyResultError{}},
		{result: ResultValue{Head: "test"}, want: &EmptyResultError{Word: "test"}},
		{result: ResultValue{EntryVals: []interface{}{DictionaryEntryValue{}}, Head: "test"}, want: nil},
	}

	for _, tt := range testData {
		if got := ValidateResult(tt.result); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ValidateResult returned wrong value. Got %#v. Want %#v.", got, tt.want)
		}
	}
}

func TestValidateAndReturnResult(t *testing.T) {
	testData := []struct {
		result  Result
		wantErr bool
	}{
		{result: nil, wantErr: true},
		{result: &ResultValue{EntryVals: []interface{}{DictionaryEntryValue{}}}, wantErr: true},
		{result: &ResultValue{Head: "test"}, wantErr: true},
		{result: &ResultValue{EntryVals: []interface{}{DictionaryEntryValue{}}, Head: "test"}, wantErr: false},
	}

	for _, tt := range testData {
		got, err := ValidateAndReturnResult(tt.result)

		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateAndReturnResult returned an error when not expected. Got %#v.", err)
		}

		if !tt.wantErr && got != tt.result {
			t.Errorf("ValidateAndReturnResult returned wrong value. Got %#v. Want %#v.", got, tt.result)
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

	if "" == msg {
		t.Errorf("Error returned an empty message")
	}

	if !strings.Contains(msg, word) {
		t.Errorf("Error message %q didn't contain word %q", msg, word)
	}
}

func TestAuthenticationError_Error(t *testing.T) {
	msg := (&AuthenticationError{}).Error()

	if "" == msg {
		t.Errorf("Error returned an empty message")
	}
}

func TestInvalidResponseError_Error(t *testing.T) {
	msg := (&InvalidResponseError{}).Error()

	if "" == msg {
		t.Errorf("Error returned an empty message")
	}
}
