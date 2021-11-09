package changeset

// Convert a struct as changeset, every field's value will be treated as changes. Returns a new changeset.
// PK changes in the changeset created with this function will be ignored
func Convert(data interface{}) *Changeset {
	ch := &Changeset{}
	ch.values = make(map[string]interface{})
	ch.changes, ch.types, _ = mapSchema(data, false)
	ch.ignorePrimary = true // set ignore primary to prevent implicit PK change
	return ch
}
