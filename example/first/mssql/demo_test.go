//+build !sqlg

package mssql_test

import (
	"context"
	"database/sql"
	"github.com/clementauger/sqlg/example/first/model"
	store "github.com/clementauger/sqlg/example/first/mssql"
	"testing"
)

func TestCreateAuthor(t *testing.T) {
	var store store.MyDatastore

	ctx := context.Background()
	db := &db{}
	var a model.Author
	store.CreateAuthor(ctx, db, a)
	wantQuery := `INSERT INTO authors ( id,bio ) VALUES ( @p0,@p1 )`
	if db.gotQuery != wantQuery {
		t.Fatalf("invalid query\nwanted=%q\ngot   =%q", wantQuery, db.gotQuery)
	}
	for i, a := range db.gotArgs {
		if x, ok := a.(sql.NamedArg); !ok {
			t.Fatalf("invalid argument value at index %v\nwanted=%T\ngot   =%T", i, x, a)
		}
	}
	wantName := "p0"
	gotName := db.gotArgs[0].(sql.NamedArg).Name
	if gotName != wantName {
		t.Fatalf("invalid argument name at index %v\nwanted=%q\ngot   =%q", 0, wantName, gotName)
	}
}
