package changeset

import (
	"testing"
	"time"

	"github.com/Fs02/changeset/params"
	"github.com/Fs02/rel"
	"github.com/stretchr/testify/assert"
)

type Status string

type User struct {
	ID           int
	Name         string
	Age          int
	Transactions []Transaction `ref:"id" fk:"user_id"`
	Address      Address
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Transaction struct {
	ID      int
	Item    string
	Status  Status
	BuyerID int  `db:"user_id"`
	Buyer   User `ref:"user_id" fk:"id"`
}

type Notes string

func (n Notes) Equal(other interface{}) bool {
	if o, ok := other.(Notes); ok {
		return n == o
	}

	return false
}

type Address struct {
	ID        int
	UserID    *int
	User      *User
	Street    string
	Notes     Notes
	Flagged   *bool
	DeletedAt *time.Time
}

type Owner struct {
	User   *User
	UserID *int
}

func TestChangeset(t *testing.T) {
	ch := Changeset{}
	assert.Nil(t, ch.Errors())
	assert.Nil(t, ch.Error())
	assert.Nil(t, ch.Changes())
	assert.Nil(t, ch.Values())
	assert.Nil(t, ch.Types())
	assert.Nil(t, ch.Constraints())
}

func TestChangeset_Get(t *testing.T) {
	ch := Changeset{
		changes: map[string]interface{}{
			"a": 2,
		},
	}

	assert.Equal(t, 2, ch.Get("a"))
	assert.Equal(t, nil, ch.Get("b"))
	assert.Equal(t, 1, len(ch.changes))
}

func TestChangeset_Fetch(t *testing.T) {
	ch := Changeset{
		changes: map[string]interface{}{
			"a": 1,
		},
		values: map[string]interface{}{
			"b": 2,
		},
	}

	assert.Equal(t, 1, ch.Fetch("a"))
	assert.Equal(t, 2, ch.Fetch("b"))
	assert.Equal(t, nil, ch.Fetch("c"))
	assert.Equal(t, 1, len(ch.changes))
	assert.Equal(t, 1, len(ch.values))
}

func TestChangesetApply(t *testing.T) {
	var (
		user    User
		now     = time.Now().Truncate(time.Second)
		flagged = true
		doc     = rel.NewDocument(&user)
		input   = params.Map{
			"name": "Luffy",
			"age":  20,
			"transactions": []params.Map{
				{
					"item":   "Sword",
					"status": "pending",
				},
				{
					"item":   "Shield",
					"status": "paid",
				},
			},
			"address": params.Map{
				"street":  "Grove Street",
				"notes":   "Brown fox jumps",
				"flagged": true,
			},
		}
		userMutation = rel.Apply(rel.NewDocument(&User{}),
			rel.Set("name", "Luffy"),
			rel.Set("age", 20),
			rel.Set("created_at", now),
			rel.Set("updated_at", now),
		)
		transaction1Mutation = rel.Apply(rel.NewDocument(&Transaction{}),
			rel.Set("item", "Sword"),
			rel.Set("status", Status("pending")),
		)
		transaction2Mutation = rel.Apply(rel.NewDocument(&Transaction{}),
			rel.Set("item", "Shield"),
			rel.Set("status", Status("paid")),
		)
		addressMutation = rel.Apply(rel.NewDocument(&Address{}),
			rel.Set("street", "Grove Street"),
			rel.Set("notes", Notes("Brown fox jumps")),
			rel.Set("flagged", true),
		)
	)

	ch := Cast(user, input, []string{"name", "age"})
	CastAssoc(ch, "transactions", func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"item", "status"})
		return ch
	})
	CastAssoc(ch, "address", func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"street", "notes", "flagged"})
		return ch
	})

	userMutation.SetAssoc("transactions", transaction1Mutation, transaction2Mutation)
	userMutation.SetAssoc("address", addressMutation)

	assert.Nil(t, ch.Error())
	assert.Equal(t, userMutation, rel.Apply(doc, ch))
	assert.Equal(t, User{
		Name:      "Luffy",
		Age:       20,
		CreatedAt: now,
		UpdatedAt: now,
		Transactions: []Transaction{
			{Item: "Sword", Status: "pending"},
			{Item: "Shield", Status: "paid"},
		},
		Address: Address{
			Street:  "Grove Street",
			Notes:   "Brown fox jumps",
			Flagged: &flagged,
		},
	}, user)
}

func TestChangesetApply_constraint(t *testing.T) {
	var (
		user  User
		now   = time.Now().Truncate(time.Second)
		doc   = rel.NewDocument(&user)
		input = params.Map{
			"name": "Luffy",
			"age":  20,
		}
		userMutation = rel.Apply(rel.NewDocument(&User{}),
			rel.Set("name", "Luffy"),
			rel.Set("age", 20),
			rel.Set("created_at", now),
			rel.Set("updated_at", now),
		)
	)

	ch := Cast(user, input, []string{"name", "age"})
	UniqueConstraint(ch, "name")
	mut := rel.Apply(doc, ch)

	assert.Nil(t, ch.Error())
	assert.Equal(t, userMutation.Mutates, mut.Mutates)
	assert.NotNil(t, mut.ErrorFunc)
	assert.Equal(t, User{
		Name:      "Luffy",
		Age:       20,
		CreatedAt: now,
		UpdatedAt: now,
	}, user)
}
