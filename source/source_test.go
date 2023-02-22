// Copyright © 2023 Trevor N. Suarez (Rican7)

package source

import "testing"

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
