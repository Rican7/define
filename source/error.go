package source

const emptyResultErrorMessage = "The source returned an empty result"

// EmptyResultError represents an error caused by an empty result
type EmptyResultError struct {
}

func (e *EmptyResultError) Error() string {
	return emptyResultErrorMessage
}
