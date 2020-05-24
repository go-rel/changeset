package changeset

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/assert"
)

func TestConstraint_GetError(t *testing.T) {
	ch := &Changeset{}
	UniqueConstraint(ch, "slug")
	ForeignKeyConstraint(ch, "user_id", Name("user_id_ibfk1"), Exact(true))
	CheckConstraint(ch, "state")

	tests := []struct {
		name     string
		err      rel.ConstraintError
		expected error
	}{
		{
			name:     "unique",
			err:      rel.ConstraintError{Key: "slug_unique_index", Type: rel.UniqueConstraint},
			expected: NewError("slug has already been taken", "slug"),
		},
		{
			name:     "fk",
			err:      rel.ConstraintError{Key: "user_id_ibfk1", Type: rel.ForeignKeyConstraint},
			expected: NewError("does not exist", "user_id"),
		},
		{
			name:     "check",
			err:      rel.ConstraintError{Key: "state_check", Type: rel.CheckConstraint},
			expected: NewError("state is invalid", "state"),
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ch.Constraints().GetError(tt.err))
		})
	}
}
