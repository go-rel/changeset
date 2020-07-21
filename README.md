# changeset

[![GoDoc](https://godoc.org/github.com/Fs02/changeset?status.svg)](https://godoc.org/github.com/Fs02/changeset) 
[![Build Status](https://travis-ci.org/Fs02/changeset.svg?branch=master)](https://travis-ci.org/Fs02/changeset) 
[![Go Report Card](https://goreportcard.com/badge/github.com/Fs02/changeset)](https://goreportcard.com/report/github.com/Fs02/changeset)
[![Gocover](https://gocover.io/_badge/github.com/Fs02/changeset)](https://gocover.io/github.com/Fs02/changeset)

Changeset mutator for [REL](https://github.com/Fs02/rel). Changesets allow filtering, casting, validation and definition of constraints when manipulating structs.

## Install

```bash
go get github.com/Fs02/changeset
```

## Example

```golang
package main

import (
	"time"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/mysql"
	"github.com/Fs02/changeset"
	"github.com/Fs02/changeset/params"
)

type Product struct {
	ID        int
	Name      string
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ChangeProduct prepares data before database operation.
// Such as casting value to appropriate types and perform validations.
func ChangeProduct(product interface{}, params params.Params) *changeset.Changeset {
	ch := changeset.Cast(product, params, []string{"name", "price"})
	changeset.ValidateRequired(ch, []string{"name", "price"})
	changeset.ValidateMin(ch, "price", 100)
	return ch
}

func main() {
    // initialize mysql adapter.
    adapter, _ := mysql.Open(dsn)
    defer adapter.Close()

    // initialize rel's repo.
    repo := rel.New(adapter)

	var product Product

	// Inserting Products.
	// Changeset is used when creating or updating your data.
	ch := ChangeProduct(product, params.Map{
		"name":  "shampoo",
		"price": 1000,
	})

	if ch.Error() != nil {
		// handle error
	}

	// Changeset can also be created directly from json string.
	jsonch := ChangeProduct(product, params.ParseJSON(`{
		"name":  "soap",
		"price": 2000,
	}`))

	// Create products with changeset and return the result to &product,
	if err := repo.Insert(context.TODO(), &product, ch); err != nil {
		// handle error
	}
}
```

## License

Released under the [MIT License](https://github.com/Fs02/changeset/blob/master/LICENSE)
