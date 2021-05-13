package changeset

import (
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestConstraint_GetError(t *testing.T) {
	ch := &Changeset{}
	UniqueConstraint(ch, "slug")
	ForeignKeyConstraint(ch, "user_id", Name("user_id_ibfk1"), Exact(true))
	CheckConstraint(ch, "state")

	tests := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "unique",
			err:      rel.ConstraintError{Key: "slug_unique_index", Type: rel.UniqueConstraint},
			expected: Error{Message: "slug has already been taken", Field: "slug", Err: rel.ConstraintError{Key: "slug_unique_index", Type: rel.UniqueConstraint}},
		},
		{
			name:     "fk",
			err:      rel.ConstraintError{Key: "user_id_ibfk1", Type: rel.ForeignKeyConstraint},
			expected: Error{Message: "does not exist", Field: "user_id", Err: rel.ConstraintError{Key: "user_id_ibfk1", Type: rel.ForeignKeyConstraint}},
		},
		{
			name:     "check",
			err:      rel.ConstraintError{Key: "state_check", Type: rel.CheckConstraint},
			expected: Error{Message: "state is invalid", Field: "state", Err: rel.ConstraintError{Key: "state_check", Type: rel.CheckConstraint}},
		},
		{
			name:     "undefined unique",
			err:      rel.ConstraintError{Key: "other_unique_index", Type: rel.UniqueConstraint},
			expected: rel.ConstraintError{Key: "other_unique_index", Type: rel.UniqueConstraint},
		},
		{
			name:     "undefined fk",
			err:      rel.ConstraintError{Key: "other_id_ibfk1", Type: rel.ForeignKeyConstraint},
			expected: rel.ConstraintError{Key: "other_id_ibfk1", Type: rel.ForeignKeyConstraint},
		},
		{
			name:     "other error",
			err:      rel.NotFoundError{},
			expected: rel.NotFoundError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ch.Constraints().GetError(tt.err))
		})
	}
}
