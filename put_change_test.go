package changeset

import (
	"reflect"
	"testing"

	"github.com/go-rel/changeset/params"
	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestPutChange(t *testing.T) {
	ch := &Changeset{
		changes: make(map[string]interface{}),
		values: map[string]interface{}{
			"field1": 0,
		},
		types: map[string]reflect.Type{
			"field1": reflect.TypeOf(0),
		},
	}

	assert.Nil(t, ch.Error())
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 0, len(ch.Changes()))

	// normal put changes
	PutChange(ch, "field1", 10)
	assert.Nil(t, ch.Error())
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 1, len(ch.Changes()))
	assert.Equal(t, 10, ch.Changes()["field1"])

	// put changes not valid and not allowed to error
	PutChange(ch, "field1", "10")
	assert.NotNil(t, ch.Error())
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, 1, len(ch.Errors()))
	assert.Equal(t, "field1 is invalid", ch.Error().Error())
	assert.Equal(t, 1, len(ch.Changes()))
	assert.Equal(t, 10, ch.Changes()["field1"])
}

func TestPutChange_ptr(t *testing.T) {
	var (
		nullable = false
		a        struct {
			ID       int
			Nullable *bool
		}
	)

	ch := Cast(a, params.Map{}, []string{})
	PutChange(ch, "nullable", &nullable)

	assert.Nil(t, ch.Error())
	assert.Equal(t, &nullable, ch.Changes()["nullable"])

	assert.NotPanics(t, func() {
		rel.Apply(rel.NewDocument(&a), ch)
		assert.Equal(t, &nullable, a.Nullable)
	})
}

func TestPutChange_nil(t *testing.T) {
	var a struct {
		ID       int
		Nullable *bool
	}

	ch := Cast(a, params.Map{}, []string{})
	PutChange(ch, "nullable", nil)

	assert.Nil(t, ch.Error())
	assert.Equal(t, nil, ch.Changes()["nullable"])

	assert.NotPanics(t, func() {
		rel.Apply(rel.NewDocument(&a), ch)
		assert.Nil(t, a.Nullable)
	})
}

func TestPutChange_typedNil(t *testing.T) {
	var a struct {
		ID       int
		Nullable *bool
	}

	ch := Cast(a, params.Map{}, []string{})
	PutChange(ch, "nullable", (*bool)(nil))

	assert.Nil(t, ch.Error())
	assert.Equal(t, nil, ch.Changes()["nullable"])

	assert.NotPanics(t, func() {
		rel.Apply(rel.NewDocument(&a), ch)
		assert.Nil(t, a.Nullable)
	})
}
