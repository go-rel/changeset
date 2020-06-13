package changeset

// Error struct.
type Error struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Code    int    `json:"code,omitempty"`
	Err     error  `json:"-"`
}

// Error prints error message.
func (e Error) Error() string {
	return e.Message
}

// Unwrap internal error.
func (e Error) Unwrap() error {
	return e.Err
}

// AddError adds an error to changeset.
//	ch := changeset.Cast(user, params, fields)
//	changeset.AddError(ch, "field", "error")
//	ch.Errors() // []errors.Error{{Field: "field", Message: "error"}}
func AddError(ch *Changeset, field string, message string) {
	ch.errors = append(ch.errors, Error{Message: message, Field: field})
}
