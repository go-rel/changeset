// Package changeset used to cast and validate data before saving it to the database.
package changeset

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/go-rel/changeset/params"
	"github.com/go-rel/rel"
)

// Changeset used to cast and validate data before saving it to the database.
type Changeset struct {
	errors        []error
	params        params.Params
	changes       map[string]interface{}
	values        map[string]interface{}
	types         map[string]reflect.Type
	constraints   Constraints
	zero          bool
	ignorePrimary bool
}

// Errors of changeset.
func (c Changeset) Errors() []error {
	return c.errors
}

// Error of changeset, returns the first error if any.
func (c Changeset) Error() error {
	if c.errors != nil {
		return c.errors[0]
	}

	return nil
}

// Get a change from changeset.
func (c Changeset) Get(field string) interface{} {
	return c.changes[field]
}

// Fetch a change or value from changeset.
func (c Changeset) Fetch(field string) interface{} {
	if change, ok := c.changes[field]; ok {
		return change
	}

	return c.values[field]
}

// Changes of changeset.
func (c Changeset) Changes() map[string]interface{} {
	return c.changes
}

// Values of changeset.
func (c Changeset) Values() map[string]interface{} {
	return c.values
}

// Types of changeset.
func (c Changeset) Types() map[string]reflect.Type {
	return c.types
}

// Constraints of changeset.
func (c Changeset) Constraints() Constraints {
	return c.constraints
}

// Apply mutation.
func (c *Changeset) Apply(doc *rel.Document, mut *rel.Mutation) {
	var (
		pField        = doc.PrimaryField()
		mutablePField = false
		now           = time.Now().Truncate(time.Second)
	)

	switch c.values[pField].(type) {
	case int:
	default:
		if c.changes[pField] != nil {
			mutablePField = true
		}
	}

	for field, value := range c.changes {
		switch v := value.(type) {
		case *Changeset:
			if mut.Cascade {
				c.applyAssocOne(doc, field, mut, v)
			}
		case []*Changeset:
			if mut.Cascade {
				c.applyAssocMany(doc, field, mut, v)
			}
		default:
			if pField != field && scannable(c.types[field]) { //if not PK - try to set
				c.set(doc, mut, field, v)
			} else if scannable(c.types[field]) { // if settable PK - check if not serial and was set manually
				if mutablePField && !c.ignorePrimary {
					c.set(doc, mut, field, v)
				}
			}

		}
	}

	// insert timestamp
	if doc.Flag(rel.HasCreatedAt) {
		if value, ok := doc.Value("created_at"); ok && value.(time.Time).IsZero() {
			c.set(doc, mut, "created_at", now)
		}
	}

	// update timestamp
	if doc.Flag(rel.HasUpdatedAt) {
		c.set(doc, mut, "updated_at", now)
	}

	// add error func
	if len(c.constraints) > 0 {
		mut.ErrorFunc = c.constraints.GetError
	}
}

func (c *Changeset) set(doc *rel.Document, mut *rel.Mutation, field string, value interface{}) {
	if doc.SetValue(field, value) {
		mut.Add(rel.Set(field, value))
	}
}

func (c *Changeset) applyAssocOne(doc *rel.Document, field string, mut *rel.Mutation, ch *Changeset) {
	var (
		assoc       = doc.Association(field)
		assocDoc, _ = assoc.Document()
		assocMut    = rel.Apply(assocDoc, ch)
	)

	mut.SetAssoc(field, assocMut)
}

func (c *Changeset) applyAssocMany(doc *rel.Document, field string, mut *rel.Mutation, chs []*Changeset) {
	var (
		assoc  = doc.Association(field)
		col, _ = assoc.Collection()
		muts   = make([]rel.Mutation, len(chs))
	)

	// Reset assoc.
	col.Reset()

	for i := range chs {
		muts[i] = rel.Apply(col.Add(), chs[i])
	}

	mut.SetAssoc(field, muts...)
}

var (
	rtTime    = reflect.TypeOf(time.Time{})
	rtScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
)

func scannable(rt reflect.Type) bool {
	var (
		kind = rt.Kind()
	)

	return !((kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Array) &&
		kind != reflect.Uint8 &&
		!rt.Implements(rtScanner) &&
		!rt.ConvertibleTo(rtTime))
}
