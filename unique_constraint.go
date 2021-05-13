package changeset

import (
	"strings"

	"github.com/go-rel/rel"
)

// UniqueConstraintMessage is the default error message for UniqueConstraint.
var UniqueConstraintMessage = "{field} has already been taken"

// UniqueConstraint adds an unique constraint to changeset.
func UniqueConstraint(ch *Changeset, field string, opts ...Option) {
	options := Options{
		message: UniqueConstraintMessage,
		name:    field,
		exact:   false,
	}
	options.apply(opts)

	ch.constraints = append(ch.constraints, Constraint{
		Field:   field,
		Message: strings.Replace(options.message, "{field}", field, 1),
		Code:    options.code,
		Name:    options.name,
		Exact:   options.exact,
		Type:    rel.UniqueConstraint,
	})
}
