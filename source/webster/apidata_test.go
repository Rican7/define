package webster

import (
	"reflect"
	"testing"
)

func TestCleanHeadword(t *testing.T) {
	for testName, testData := range map[string]struct {
		text string
		want string
	}{
		"empty": {
			text: "",
			want: "",
		},
		"no marks": {
			text: "tree",
			want: "tree",
		},
		"single mark": {
			text: "re*fuse",
			want: "refuse",
		},
		"multiple marks": {
			text: "vo*lu*mi*nous",
			want: "voluminous",
		},
	} {
		t.Run(testName, func(t *testing.T) {
			if got := cleanHeadword(testData.text); got != testData.want {
				t.Errorf("cleanHeadword returned wrong value. Got %#v. Want %#v.", got, testData.want)
			}
		})
	}
}

func TestCleanTextOfTokens(t *testing.T) {
	for testName, testData := range map[string]struct {
		text string
		want string
	}{
		"empty": {
			text: "",
			want: "",
		},
		"empty contents": {
			text: "{wi}{/wi}",
			want: "",
		},
		"unbalanced": {
			text: "{bc}{sx|test||} ",
			want: "test ",
		},
		"attributes": {
			text: "{bc}testing a {bc}{sx|test|test:2|}",
			want: "testing a test",
		},
		"complex": {
			text: "{bc}test {a_link|test} or test {bc}{sx|test||}",
			want: "test test or test test",
		},
	} {
		t.Run(testName, func(t *testing.T) {
			if got := cleanTextOfTokens(testData.text); got != testData.want {
				t.Errorf("cleanTextOfTokens returned wrong value. Got %#v. Want %#v.", got, testData.want)
			}
		})
	}
}

func TestParseSenseNumber(t *testing.T) {
	for testName, testData := range map[string]struct {
		toParse any
		want    *apiSenseNumber
	}{
		"nil": {
			toParse: nil,
			want:    nil,
		},
		"numeral": {
			toParse: "1",
			want: &apiSenseNumber{
				number: 1,
				letter: "",
				sub:    "",
			},
		},
		"letter": {
			toParse: "a",
			want: &apiSenseNumber{
				number: 0,
				letter: "a",
				sub:    "",
			},
		},
		"sub": {
			toParse: "(1)",
			want: &apiSenseNumber{
				number: 0,
				letter: "",
				sub:    "(1)",
			},
		},
		"numeral and letter": {
			toParse: "2 a",
			want: &apiSenseNumber{
				number: 2,
				letter: "a",
				sub:    "",
			},
		},
		"numeral and letter and sub": {
			toParse: "2 a (1)",
			want: &apiSenseNumber{
				number: 2,
				letter: "a",
				sub:    "(1)",
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			if got := parseSenseNumber(testData.toParse); !reflect.DeepEqual(got, testData.want) {
				t.Errorf("parseSenseNumber returned wrong value. Got %#v. Want %#v.", got, testData.want)
			}
		})
	}
}
