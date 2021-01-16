//+build sqlg,pg

package main

import (
	"fmt"

	"github.com/clementauger/sqlg/example/first/model"
)

// CreateTable authors.
func (m myDatastore) CreateTable() (err error) {
	return fmt.Errorf("todo")
}

func (m myDatastore) DeleteManyAuthors(ids []int) (ab []model.Author, err error) {
	m.Query(`DELETE FROM authors WHERE id ANY ( {{.ids | pqArray}}::int[] ) RETURNING *`)
	return
}

func (m myDatastore) DeleteAuthors() (err error) {
	m.Exec(`DELETE FROM authors WHERE bio = ''`)
	return
}

func (m *myDatastore) CreateSomeValues(v model.SomeType) (id int64, err error) {
	m.Exec(`{{$fields := fields .v "id"}}
		{{$vals := fields .v "id" "v"}}
		INSERT INTO sometype ( {{$fields | cols}} )
		VALUES ( {{$vals | vals}}, {{.v | pqArray}} )`).InsertedID(id)
	return
}
