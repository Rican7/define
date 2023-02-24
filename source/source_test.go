// Copyright © 2023 Trevor N. Suarez (Rican7)

package source

import (
	"reflect"
	"testing"
)

func TestDictionaryResults_IsSortedForPrimaryResult(t *testing.T) {
	for testName, testData := range map[string]struct {
		results DictionaryResults
		word    string
		want    bool
	}{
		"nil": {
			results: nil,
			word:    "",
			want:    true,
		},
		"empty": {
			results: DictionaryResults{},
			word:    "",
			want:    true,
		},
		"single result and sorted": {
			results: DictionaryResults{{Word: "test"}},
			word:    "test",
			want:    true,
		},
		"multi results and sorted": {
			results: DictionaryResults{
				{Word: "test"},
				{Word: "nah"},
				{Word: "nope"},
				{Word: "def not"},
				{Word: "pls stop"},
			},
			word: "test",
			want: true,
		},
		"single results and no match": {
			results: DictionaryResults{{Word: "nah"}},
			word:    "test",
			want:    false,
		},
		"multi results and no match": {
			results: DictionaryResults{
				{Word: "nah"},
				{Word: "nope"},
				{Word: "def not"},
				{Word: "pls stop"},
			},
			word: "test",
			want: false,
		},
		"multi results and not sorted": {
			results: DictionaryResults{
				{Word: "nah"},
				{Word: "nope"},
				{Word: "test"},
				{Word: "def not"},
				{Word: "pls stop"},
			},
			word: "test",
			want: false,
		},
	} {
		t.Run(testName, func(t *testing.T) {
			if got := testData.results.IsSortedForPrimaryResult(testData.word); got != testData.want {
				t.Errorf("IsSortedForPrimaryResult returned wrong value. Got %#v. Want %#v.", got, testData.want)
			}
		})
	}
}

func TestDictionaryResults_SortForPrimaryResult(t *testing.T) {
	for testName, testData := range map[string]struct {
		results DictionaryResults
		word    string
		want    DictionaryResults
	}{
		"nil": {
			results: nil,
			word:    "",
			want:    nil,
		},
		"empty": {
			results: DictionaryResults{},
			word:    "",
			want:    DictionaryResults{},
		},
		"single result and sorted": {
			results: DictionaryResults{{Word: "test"}},
			word:    "test",
			want:    DictionaryResults{{Word: "test"}},
		},
		"multi results and sorted": {
			results: DictionaryResults{
				{Word: "test"},
				{Word: "nah"},
				{Word: "nope"},
				{Word: "def not"},
				{Word: "pls stop"},
			},
			word: "test",
			want: DictionaryResults{
				{Word: "test"},
				{Word: "nah"},
				{Word: "nope"},
				{Word: "def not"},
				{Word: "pls stop"},
			},
		},
		"single results and no match": {
			results: DictionaryResults{{Word: "nah"}},
			word:    "test",
			want:    DictionaryResults{{Word: "nah"}},
		},
		"multi results and no match": {
			results: DictionaryResults{
				{Word: "nah"},
				{Word: "nope"},
				{Word: "def not"},
				{Word: "pls stop"},
			},
			word: "test",
			want: DictionaryResults{
				{Word: "nah"},
				{Word: "nope"},
				{Word: "def not"},
				{Word: "pls stop"},
			},
		},
		"multi results and not sorted with match in middle": {
			results: DictionaryResults{
				{Word: "nah"},
				{Word: "nope"},
				{Word: "test"},
				{Word: "def not"},
				{Word: "pls stop"},
			},
			word: "test",
			want: DictionaryResults{
				{Word: "test"},
				{Word: "nah"},
				{Word: "nope"},
				{Word: "def not"},
				{Word: "pls stop"},
			},
		},
		"multi results and not sorted with match at end": {
			results: DictionaryResults{
				{Word: "nah"},
				{Word: "nope"},
				{Word: "def not"},
				{Word: "pls stop"},
				{Word: "test"},
			},
			word: "test",
			want: DictionaryResults{
				{Word: "test"},
				{Word: "nah"},
				{Word: "nope"},
				{Word: "def not"},
				{Word: "pls stop"},
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			testData.results.SortForPrimaryResult(testData.word)
			if !reflect.DeepEqual(testData.results, testData.want) {
				t.Errorf("SortForPrimaryResult returned wrong value. Sorted %#v. Want %#v.", testData.results, testData.want)
			}
		})
	}
}

func TestPronunciations_String(t *testing.T) {
	for testName, testData := range map[string]struct {
		pronunciations Pronunciations
		want           string
	}{
		"nil": {
			pronunciations: nil,
			want:           "",
		},
		"empty": {
			pronunciations: Pronunciations{},
			want:           "",
		},
		"one": {
			pronunciations: Pronunciations{"test-1"},
			want:           "/test-1/",
		},
		"two": {
			pronunciations: Pronunciations{"test-1", "test-2"},
			want:           "/test-1/ (/test-2/)",
		},
		"three": {
			pronunciations: Pronunciations{"test-1", "test-2", "test-3"},
			want:           "/test-1/ (/test-2/ /test-3/)",
		},
	} {
		t.Run(testName, func(t *testing.T) {
			if got := testData.pronunciations.String(); got != testData.want {
				t.Errorf("String returned wrong value. Got %#v. Want %#v.", got, testData.want)
			}
		})
	}
}

func TestPronunciation_String(t *testing.T) {
	for testName, testData := range map[string]struct {
		pronunciation Pronunciation
		want          string
	}{
		"empty": {
			pronunciation: Pronunciation(""),
			want:          "//",
		},
		"word": {
			pronunciation: Pronunciation("test-1"),
			want:          "/test-1/",
		},
	} {
		t.Run(testName, func(t *testing.T) {
			if got := testData.pronunciation.String(); got != testData.want {
				t.Errorf("String returned wrong value. Got %#v. Want %#v.", got, testData.want)
			}
		})
	}
}

func TestAttributedText_String(t *testing.T) {
	for testName, testData := range map[string]struct {
		attributedText AttributedText
		want           string
	}{
		"empty": {
			attributedText: AttributedText{},
			want:           "\"\"",
		},
		"text only": {
			attributedText: AttributedText{
				Text: "test",
			},
			want: "\"test\"",
		},
		"text and author": {
			attributedText: AttributedText{
				Text: "test",

				Attribution: Attribution{
					Author: "Mr. Testy",
				},
			},
			want: "\"test\" - Mr. Testy",
		},
		"text and source": {
			attributedText: AttributedText{
				Text: "test",

				Attribution: Attribution{
					Source: "WikiTest",
				},
			},
			want: "\"test\" (WikiTest)",
		},
		"text, author, and source": {
			attributedText: AttributedText{
				Text: "test",

				Attribution: Attribution{
					Author: "Mr. Testy",
					Source: "WikiTest",
				},
			},
			want: "\"test\" - Mr. Testy (WikiTest)",
		},
	} {
		t.Run(testName, func(t *testing.T) {
			if got := testData.attributedText.String(); got != testData.want {
				t.Errorf("String returned wrong value. Got %#v. Want %#v.", got, testData.want)
			}
		})
	}
}
