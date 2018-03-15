package source

import "fmt"

const (
	emptyResultErrorMessage         = "The source returned an empty result"
	errorMessageForWordSuffixFormat = " for word: %q"
)

// EmptyResultError represents an error caused by an empty result
type EmptyResultError struct {
	Word string
}

func (e *EmptyResultError) Error() string {
	msg := emptyResultErrorMessage

	if "" != e.Word {
		msg = msg + fmt.Sprintf(errorMessageForWordSuffixFormat, e.Word)
	}

	return msg
}
