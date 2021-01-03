//+build !sqlg

package oracle_test

import (
	"context"
	"github.com/clementauger/sqlg/example/first/model"
	store "github.com/clementauger/sqlg/example/first/oracle"
	"testing"
)

func TestCreateAuthor(t *testing.T) {
	var store store.MyDatastore

	ctx := context.Background()
	db := &db{}
	var a model.Author
	store.CreateAuthor(ctx, db, a)
	wantQuery := `INSERT INTO authors ( id,bio ) VALUES ( ?,? )`
	if db.gotQuery != wantQuery {
		t.Fatalf("invalid query\nwanted=%q\ngot   =%q", wantQuery, db.gotQuery)
	}
	gotValue := db.gotArgs[0]
	if wantedValue, gotOk := db.gotArgs[0].(int); !gotOk {
		t.Fatalf("invalid argument type at index %v\nwanted=%T\ngot   =%T", 0, wantedValue, gotValue)
	}
	gotValue = db.gotArgs[1]
	if wantedValue, gotOk := db.gotArgs[1].(string); !gotOk {
		t.Fatalf("invalid argument type at index %v\nwanted=%T\ngot   =%T", 1, wantedValue, gotValue)
	}
}
