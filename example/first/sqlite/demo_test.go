//+build !sqlg

package sqlite_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/clementauger/sqlg/example/first/model"
	store "github.com/clementauger/sqlg/example/first/sqlite"
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

	mock.ExpectExec("INSERT INTO authors ( bio ) VALUES ( ? )").
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

func TestCreateAuthors(t *testing.T) {
	var store store.MyDatastore

	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO authors ( bio ) VALUES ( ? ) , ( ? )").
		WithArgs("bio", "bio2").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// now we execute our method
	ctx := context.Background()
	var a []model.Author
	a = append(a, model.Author{Bio: "bio"})
	a = append(a, model.Author{Bio: "bio2"})
	err = store.CreateAuthors2(ctx, db, a)
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateAuthor2(t *testing.T) {
	var store store.MyDatastore

	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO authors ( bio ) VALUES ( ? )").
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

func TestUpdateAuthor(t *testing.T) {
	var store store.MyDatastore

	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec(`UPDATE authors SET
        			bio = ?
        		 WHERE id = ?`).
		WithArgs("", 0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// now we execute our method
	ctx := context.Background()
	var a model.Author
	err = store.UpdateAuthor(ctx, db, a)
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
