//+build sqlg
//+build mysql oracle sqlite

package main

import (
	"fmt"
	"github.com/clementauger/sqlg/example/first/model"
)

func (m *myDatastore) DeleteManyAuthors(ids []int) (_ []model.Author, err error) {
	m.Exec(`DELETE FROM authors WHERE id IN ( {{.ids}} )`)
	return
}

func (m *myDatastore) DeleteAuthors() (err error) {
	m.Exec(`DELETE FROM authors WHERE bio = ''`)
	return
}

func (m *myDatastore) CreateSomeValues(v model.SomeType) (id int64, err error) {
	return id, fmt.Errorf("unsupported")
}
