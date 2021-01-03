package mssql

import (
	"github.com/clementauger/sqlg/tpl"
	"text/template"
	"time"

	mssql "github.com/denisenkom/go-mssqldb"
)

// FuncMap returns the CommonFuncMap plus those for mssql.
func FuncMap() template.FuncMap {
	out := template.FuncMap{}
	for key, fn := range tpl.FuncMap() {
		out[key] = fn
	}
	out["datetime"] = func(s time.Time) interface{} {
		return mssql.DateTime1(s)
	}
	out["datetimeoffset"] = func(s time.Time) interface{} {
		return mssql.DateTimeOffset(s)
	}
	out["nvarcharmax"] = func(s string) interface{} {
		return mssql.NVarCharMax(s)
	}
	out["varcharmax"] = func(s string) interface{} {
		return mssql.VarCharMax(s)
	}
	out["varchar"] = func(s string) interface{} {
		return mssql.VarChar(s)
	}
	return out
}
