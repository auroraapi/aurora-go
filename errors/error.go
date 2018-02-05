package errors

// Error is a generic error that is returned when something non-API related
// goes wrong
type Error struct {
	// Code is the specific error code (for debugging purposes)
	Code string `json:"code,omitempty"`

	// Message is a descriptive message of the error, why it occurred, how to resolve, etc.
	Message string `json:"message,omitempty"`

	// Info is an optional field describing in detail the error for debugging purposes.
	Info string `json:"-"`
}

func (e Error) Error() string {
	return e.Message
}
