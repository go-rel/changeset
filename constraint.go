// Package changeset used to cast and validate data before saving it to the database.
package changeset

import (
	"strings"

	"github.com/go-rel/rel"
)

// Constraint defines information to infer constraint error.
type Constraint struct {
	Field   string
	Message string
	Code    int
	Name    string
	Exact   bool
	Type    rel.ConstraintType
}

// Constraints is slice of Constraint
type Constraints []Constraint

// GetError converts error based on constraints.
// If the original error is constraint error, and it's defined in the constraint list, then it'll be updated with constraint's message.
// If the original error is constraint error but not defined in the constraint list, it'll be converted to unexpected error.
// else it'll not modify the error.
func (constraints Constraints) GetError(err error) error {
	cerr, ok := err.(rel.ConstraintError)
	if !ok {
		return err
	}

	for _, c := range constraints {
		if c.Type == cerr.Type {
			if c.Exact && c.Name != cerr.Key {
				continue
			}

			if !c.Exact && !strings.Contains(cerr.Key, c.Name) {
				continue
			}

			return Error{
				Message: c.Message,
				Field:   c.Field,
				Code:    c.Code,
				Err:     err,
			}
		}
	}

	return err
}
