package changeset

import (
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestCheckConstraint(t *testing.T) {
	ch := &Changeset{}
	assert.Nil(t, ch.Constraints())

	CheckConstraint(ch, "field1")
	assert.Equal(t, 1, len(ch.Constraints()))
	assert.Equal(t, rel.CheckConstraint, ch.Constraints()[0].Type)
}
