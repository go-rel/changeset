// Package changeset used to cast and validate data before saving it to the database.
package changeset

import (
	"strings"

	"github.com/Fs02/rel"
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
func (constraints Constraints) GetError(err rel.ConstraintError) error {
	for _, c := range constraints {
		if c.Type == err.Type {
			if c.Exact && c.Name != err.Key {
				continue
			}

			if !c.Exact && !strings.Contains(err.Key, c.Name) {
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
