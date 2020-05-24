package changeset

import (
	"strings"

	"github.com/Fs02/rel"
)

// PutAssocErrorMessage is the default error message for PutAssoc.
var PutAssocErrorMessage = "{field} is invalid"

// PutAssoc to changeset.
func PutAssoc(ch *Changeset, field string, value interface{}, opts ...Option) {
	options := Options{
		message: PutAssocErrorMessage,
	}
	options.apply(opts)

	ok := false // TODO: need adjustment from rel
	if assoc := ch.doc.Association(field); ok {
		if assoc.Type() == rel.HasMany {
			if _, ok := value.([]*Changeset); ok {
				ch.changes[field] = value
				return
			}
		} else {
			if _, ok := value.(*Changeset); ok {
				ch.changes[field] = value
				return
			}
		}
	}

	msg := strings.Replace(options.message, "{field}", field, 1)
	AddError(ch, field, msg)
}
