package changeset

import (
	"reflect"
	"strings"
)

// PutChangeErrorMessage is the default error message for PutChange.
var PutChangeErrorMessage = "{field} is invalid"

// PutChange to changeset.
func PutChange(ch *Changeset, field string, value interface{}, opts ...Option) {
	options := Options{
		message: PutChangeErrorMessage,
	}
	options.apply(opts)

	if typ, exist := ch.doc.Type(field); exist {
		if value != (interface{})(nil) {
			rt := reflect.TypeOf(value)
			if rt.Kind() == reflect.Ptr {
				rt = rt.Elem()
			}

			if rt.ConvertibleTo(typ) {
				ch.changes[field] = value
				return
			}
		} else {
			ch.changes[field] = reflect.Zero(reflect.PtrTo(typ)).Interface()
			return
		}
	}

	msg := strings.Replace(options.message, "{field}", field, 1)
	AddError(ch, field, msg)
}
