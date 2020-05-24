package changeset

import (
	"reflect"
	"strings"
)

// PutDefaultErrorMessage is the default error message for PutDefault.
var PutDefaultErrorMessage = "{field} is invalid"

// PutDefault to changeset.
func PutDefault(ch *Changeset, field string, value interface{}, opts ...Option) {
	options := Options{
		message: PutDefaultErrorMessage,
	}
	options.apply(opts)

	if typ, exist := ch.doc.Type(field); exist {
		rt := reflect.TypeOf(value)
		if rt.ConvertibleTo(typ) {
			existingValue, _ := ch.doc.Value(field)

			if (ch.params == nil || !ch.params.Exists(field)) && // no input
				ch.changes[field] == nil && // no change
				isZero(existingValue) { // existing value is zero value
				ch.changes[field] = value
			}
			return
		}
	}

	msg := strings.Replace(options.message, "{field}", field, 1)
	AddError(ch, field, msg)
}
