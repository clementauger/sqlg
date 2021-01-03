//+build !sqlg

package pg_test

import (
	"context"
	"database/sql"
)

type db struct {
	gotQuery string
	gotArgs  []interface{}
}

func (d *db) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	d.gotQuery = query
	d.gotArgs = args
	return &sql.Rows{}, nil
}
func (d *db) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	d.gotQuery = query
	d.gotArgs = args
	return sqlResult{}, nil
}

type sqlResult struct{}

func (s sqlResult) LastInsertId() (int64, error) { return 0, nil }
func (s sqlResult) RowsAffected() (int64, error) { return 0, nil }
