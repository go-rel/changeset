// Package changeset used to cast and validate data before saving it to the database.
package changeset

import (
	"github.com/Fs02/form/params"
	"github.com/Fs02/rel"
)

// Changeset used to cast and validate data before saving it to the database.
type Changeset struct {
	doc         *rel.Document
	params      params.Params
	changes     map[string]interface{}
	errors      []error
	constraints Constraints
	zero        bool
}

// Errors of changeset.
func (changeset *Changeset) Errors() []error {
	return changeset.errors
}

// Error of changeset, returns the first error if any.
func (changeset *Changeset) Error() error {
	if changeset.errors != nil {
		return changeset.errors[0]
	}
	return nil
}

// Get a change from changeset.
func (changeset *Changeset) Get(field string) interface{} {
	return changeset.changes[field]
}

// Fetch a change or value from changeset.
func (changeset *Changeset) Fetch(field string) interface{} {
	if change, ok := changeset.changes[field]; ok {
		return change
	}

	val, _ := changeset.doc.Value(field)
	return val
}

// Changes of changeset.
func (changeset *Changeset) Changes() map[string]interface{} {
	return changeset.changes
}
