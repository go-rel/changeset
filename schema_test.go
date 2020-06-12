package changeset

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type CustomSchema struct {
	UUID  string
	Price int
}

func (c CustomSchema) PrimaryKey() (string, interface{}) {
	return "_uuid", c.UUID
}

func (c CustomSchema) Fields() map[string]int {
	return map[string]int{
		"_uuid":  0,
		"_price": 1,
	}
}

func (c CustomSchema) Types() []reflect.Type {
	return []reflect.Type{reflect.TypeOf(""), reflect.TypeOf(0)}
}

func (c CustomSchema) Values() []interface{} {
	return []interface{}{c.UUID, c.Price}
}

func TestInferFields(t *testing.T) {
	var (
		record = struct {
			A string
			B *int
			C []byte     `db:",primary"`
			D bool       `db:"D"`
			E []*float64 `db:"-"`
		}{}
		rt       = reflectTypePtr(record)
		expected = map[string]int{
			"a": 0,
			"b": 1,
			"c": 2,
			"D": 3,
		}
	)

	_, cached := fieldsCache.Load(rt)
	assert.False(t, cached)

	assert.Equal(t, expected, inferFields(record))

	_, cached = fieldsCache.Load(rt)
	assert.True(t, cached)

	assert.Equal(t, expected, inferFields(&record))
}

func TestInferFields_usingInterface(t *testing.T) {
	var (
		record   = CustomSchema{}
		rt       = reflectTypePtr(record)
		expected = record.Fields()
	)

	_, cached := fieldsCache.Load(rt)
	assert.False(t, cached)

	assert.Equal(t, expected, inferFields(record))

	_, cached = fieldsCache.Load(rt)
	assert.False(t, cached)
}

func TestInferTypes(t *testing.T) {
	type userDefined float64

	var (
		record = struct {
			A string
			B *int
			C []byte
			D bool
			E []*float64
			F userDefined
			G time.Time
		}{}
		rt       = reflectTypePtr(record)
		expected = []reflect.Type{
			reflect.TypeOf(""),
			reflect.TypeOf(0),
			reflect.TypeOf([]byte{}),
			reflect.TypeOf(false),
			reflect.TypeOf([]float64{}),
			reflect.TypeOf(userDefined(0)),
			reflect.TypeOf(time.Time{}),
		}
	)

	_, cached := typesCache.Load(rt)
	assert.False(t, cached)

	assert.Equal(t, expected, inferTypes(record))

	_, cached = typesCache.Load(rt)
	assert.True(t, cached)

	assert.Equal(t, expected, inferTypes(&record))
}

func TestInferTypes_usingInterface(t *testing.T) {
	var (
		record   = CustomSchema{}
		rt       = reflectTypePtr(record)
		expected = record.Types()
	)

	_, cached := typesCache.Load(rt)
	assert.False(t, cached)

	assert.Equal(t, expected, inferTypes(record))

	_, cached = typesCache.Load(rt)
	assert.False(t, cached)
}

func TestInferValues(t *testing.T) {
	var (
		address = "address"
		record  = struct {
			ID      int
			Name    string
			Skip    bool `db:"-"`
			Number  float64
			Address *string
			Data    []byte
		}{
			ID:      1,
			Name:    "name",
			Number:  10.5,
			Address: &address,
			Data:    []byte("data"),
		}
		expected = []interface{}{1, "name", 10.5, address, []byte("data")}
	)

	assert.Equal(t, expected, inferValues(record))
	assert.Equal(t, expected, inferValues(&record))
}

func TestInferValues_usingInterface(t *testing.T) {
	var (
		record = CustomSchema{
			UUID:  "abc123",
			Price: 100,
		}
		expected = []interface{}{"abc123", 100}
	)

	assert.Equal(t, expected, inferValues(record))
	assert.Equal(t, expected, inferValues(&record))
}
