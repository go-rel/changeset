package changeset

// Error struct.
type Error struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Code    int    `json:"code,omitempty"`
	// Err     error  `json:"-"` TODO: Unwrap
}

// Error prints error message.
func (e Error) Error() string {
	return e.Message
}

// NewError creates an error with field and message.
func NewError(message string, field string) error {
	return NewErrorWithCode(message, field, 0)
}

// NewErrorWithCode creates an error with code.
func NewErrorWithCode(message string, field string, code int) error {
	return Error{message, field, code}
}

// AddError adds an error to changeset.
//	ch := changeset.Cast(user, params, fields)
//	changeset.AddError(ch, "field", "error")
//	ch.Errors() // []errors.Error{{Field: "field", Message: "error"}}
func AddError(ch *Changeset, field string, message string) {
	ch.errors = append(ch.errors, NewError(message, field))
}
