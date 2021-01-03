package pg

import (
	"github.com/clementauger/sqlg/tpl"
	"text/template"

	"github.com/lib/pq"
)

// FuncMap returns the CommonFuncMap plus those for postgresql.
func FuncMap() template.FuncMap {
	out := template.FuncMap{}
	for key, fn := range tpl.FuncMap() {
		out[key] = fn
	}
	out["pqArray"] = func(s interface{}) interface{} {
		return pq.Array(s)
	}
	return out
}
