//+build !sqlg

package pg_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/clementauger/sqlg/example/first/model"
	store "github.com/clementauger/sqlg/example/first/pg"
	"github.com/lib/pq"
)

func TestCreateAuthor(t *testing.T) {
	var store store.MyDatastore

	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO authors ( bio ) VALUES ( $0 )").
		WithArgs("").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// now we execute our method
	ctx := context.Background()
	var a model.Author
	_, err = store.CreateAuthor(ctx, db, a)
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInQuery(t *testing.T) {
	var store store.MyDatastore

	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT * FROM authors
        		WHERE id IN ($0,$1,$2)
        		GROUP BY bio
        		ORDER BY bio
        		LIMIT $3, $4`).
		WithArgs(0, 1, 2, 0, 5).
		WillReturnRows(mock.NewRows([]string{"", "bio"}))

	// now we execute our method
	ctx := context.Background()
	_, err = store.GetSomeAuthors(ctx, db, []int{0, 1, 2}, 0, 5, "bio", "bio")
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAnyQuery(t *testing.T) {
	var store store.MyDatastore

	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(`DELETE FROM authors WHERE id ANY ( $0::int[] ) RETURNING *`).
		WithArgs(pq.Array([]int{0, 1, 2})).
		WillReturnRows(mock.NewRows([]string{"", "bio"}))

	// now we execute our method
	ctx := context.Background()
	_, err = store.DeleteManyAuthors(ctx, db, []int{0, 1, 2})
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
