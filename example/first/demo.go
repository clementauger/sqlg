//+build sqlg

//go:generate go run github.com/clementauger/sqlg/cmd/sqlg -clean sqlite pg mssql oracle mysql
//go:generate go run github.com/clementauger/sqlg/cmd/sqlg -engine sqlite
//go:generate go run github.com/clementauger/sqlg/cmd/sqlg -engine pg
//go:generate go run github.com/clementauger/sqlg/cmd/sqlg -engine mssql
//go:generate go run github.com/clementauger/sqlg/cmd/sqlg -engine oracle
//go:generate go run github.com/clementauger/sqlg/cmd/sqlg -engine mysql

package main

import (
	"fmt"
	"github.com/clementauger/sqlg/example/first/model"
	sqlg "github.com/clementauger/sqlg/runtime"
)

// myDatastore stores stuff.
type myDatastore struct {
	sqlg.SQLg
	Tracer    sqlg.NilTracer
	Logger    sqlg.NilLogger
	Converter sqlg.ToSnake
}

// GetAuthor retrieves
// an Author by its ID.
func (m myDatastore) GetAuthor(id int) (a model.Author, err error) {
	m.Query(`SELECT * FROM authors WHERE id={{.id}}`)
	return
}

func (m myDatastore) GetAuthorsWihIterator(id int) (it func() (model.Author, error), err error) {
	m.Query(`SELECT * FROM authors WHERE id={{.id}}`)
	return
}

type authorIterator func() (model.Author, error)

func (m myDatastore) GetAuthorsWihNamedIterator(id int) (it authorIterator, err error) {
	m.Query(`SELECT * FROM authors WHERE id={{.id}}`)
	return
}

func (m myDatastore) GetAuthor2(id int) (a model.Author, err error) {
	m.Query(`SELECT {{cols .a "id"}} FROM authors WHERE id={{.id}}`)
	// m.Query(`SELECT {{cols a "id" | convert .SQLGConverter | glue ","}} FROM authors WHERE id={{.id}}`)
	return
}
func (m myDatastore) GetAuthor3(id int) (a model.Author, err error) {
	m.Query(`SELECT {{cols .a "id" | prefix "alias."}} FROM authors as alias WHERE alias.id={{.id}}`)
	// m.Query(`SELECT {{cols a "id" | prefix "alias" | convert .SQLGConverter | glue ","}} FROM authors WHERE id={{.id}}`)
	return
}

func (m *myDatastore) GetAuthors() (a []model.Author, err error) {
	m.Query(`SELECT * FROM authors`)
	return
}

func (m *myDatastore) GetSomeAuthors(u model.Author, start, end int, orderby, groupby string) ([]model.Author, error) {
	m.Query(`SELECT * FROM authors
		GROUP BY {{.groupby | print}}
		ORDER BY {{.orderby | raw}}
		LIMIT {{.start }}, {{.end }}
		`)
	return nil, nil
}

func (m *myDatastore) GetSomeY(u model.Y) ([]model.Y, error) {
	var k string
	m.Query(`SELECT * FROM y`)
	fmt.Println(k)
	return nil, nil
}

func (m *myDatastore) DeleteAuthor(id int) error {
	m.Exec(`DELETE FROM authors WHERE id={{.id}}`)
	return nil
}

func (m *myDatastore) DeleteAuthor2(id int) (count int64, err error) {
	m.Exec(`DELETE FROM authors WHERE id={{.id}}`).AffectedRows(count)
	return
}

func (m *myDatastore) CreateAuthor(a model.Author) (id int64, err error) {
	m.Exec(`INSERT INTO authors ( {{cols .a "id"}} ) VALUES ( {{vals .a "id"}} )`).InsertedID(id)
	return
}

func (m *myDatastore) ProductUpdate() (name string, price int, err error) {
	m.Query(`UPDATE products SET price = price * 1.10
  WHERE price <= 99.99
  RETURNING name, price AS new_price`)
	return
}

func (m *myDatastore) CreateAuthors(a []model.Author) (err error) {
	m.Exec(`INSERT INTO authors (bio)
		VALUES
		{{range $i, $a := .a}}
		 ( {{$a.Bio}} ) {{comma $i $a}}
		{{end}}
	`)
	return
}
