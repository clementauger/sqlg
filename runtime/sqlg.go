package runtime

import (
	"context"
	"database/sql"

	"github.com/iancoleman/strcase"
)

type SQLg interface {
	WithParam(name string, value interface{}) SQLg
	Query(sql string)
	Exec(sql string) resulter
	// Insert(intoTable string, value interface{}, pkfields ...string) resulter
	// Update(intoTable string, value interface{}, pkfields ...string) resulter
}
type resulter interface {
	AffectedRows(dest interface{}) resulter
	InsertedID(dest interface{}) resulter
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type CaseConverter interface {
	Convert(s string) string
}

type NilCaseConverter struct {
	converter CaseConverter
}

func (n *NilCaseConverter) Configure(converter CaseConverter) {
	n.converter = converter
}
func (n NilCaseConverter) Convert(s string) string {
	if n.converter != nil {
		s = n.converter.Convert(s)
	}
	return s
}

type ToSnake struct{}

func (n ToSnake) Convert(s string) string {
	return strcase.ToSnake(s)
}

type ToCamel struct{}

func (n ToCamel) Convert(s string) string {
	return strcase.ToCamel(s)
}
