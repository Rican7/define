package source

import "testing"

func TestRemoveDiacritics(t *testing.T) {
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
			text: "façade",
			want: "facade",
		},
		"multiple marks": {
			text: "résumé",
			want: "resume",
		},
	} {
		t.Run(testName, func(t *testing.T) {
			if got := RemoveDiacritics(testData.text); got != testData.want {
				t.Errorf("RemoveDiacritics returned wrong value. Got %#v. Want %#v.", got, testData.want)
			}
		})
	}
}

func TestEqualFoldPlain(t *testing.T) {
	for testName, testData := range map[string]struct {
		s    string
		t    string
		want bool
	}{
		"empty": {
			s:    "",
			t:    "",
			want: true,
		},
		"no marks": {
			s:    "tree",
			t:    "tree",
			want: true,
		},
		"single mark": {
			s:    "façade",
			t:    "facade",
			want: true,
		},
		"multiple marks": {
			s:    "résumé",
			t:    "resume",
			want: true,
		},
		"multiple marks and capitalizations": {
			s:    "Résumé",
			t:    "resume",
			want: true,
		},
		"different": {
			s:    "test",
			t:    "resume",
			want: false,
		},
	} {
		t.Run(testName, func(t *testing.T) {
			if got := EqualFoldPlain(testData.s, testData.t); got != testData.want {
				t.Errorf("EqualFoldPlain returned wrong value. Got %#v. Want %#v.", got, testData.want)
			}
		})
	}
}
