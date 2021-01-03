//+build !sqlg

// Code generated by sqlg DO NOT EDIT

package oracle

import (
	"bytes"
	"context"
	"fmt"
	"github.com/clementauger/sqlg/example/first/model"
	sqlg "github.com/clementauger/sqlg/runtime"
	"github.com/clementauger/sqlg/tpl"
	"text/template"
)

var queryTemplates12fb19 = map[string]*template.Template{
	"myDatastore__DeleteManyAuthors": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`DELETE FROM authors WHERE id IN ( {{.ids | collect $.SQLGValues $.SQLGFlavor | placeholder $.SQLGValues $.SQLGFlavor}} )`,
	)),
}

var rawQueries12fb19 = map[string]string{
	"myDatastore__DeleteAuthors": `DELETE FROM authors WHERE bio = ''`,
}

func (m *MyDatastore) CreateSomeValues(ctx context.Context, db sqlg.Execer, v model.SomeType) (id int64, err error) {
	err = fmt.Errorf("unsupported")
	return
}

func (m *MyDatastore) DeleteAuthors(ctx context.Context, db sqlg.Execer) (err error) {
	var sqlQuery12fb19 string
	sqlQuery12fb19 = rawQueries12fb19["myDatastore__DeleteAuthors"]

	m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteAuthors", sqlQuery12fb19)
	m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteAuthors", sqlQuery12fb19)
	defer func() {
		m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteAuthors", err)
	}()

	_, err = db.ExecContext(ctx, sqlQuery12fb19)
	if err != nil {
		return
	}
	return
}

func (m *MyDatastore) DeleteManyAuthors(ctx context.Context, db sqlg.Execer, ids []int) (_ []model.Author, err error) {
	var sqlQuery12fb19 string
	SQLGValues12fb19 := &[]interface{}{}
	SQLGFlavor12fb19 := "?"
	{
		var query12fb19 bytes.Buffer
		templateInput12fb19 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues12fb19,
			"SQLGFlavor":    SQLGFlavor12fb19,
			"ids":           ids,
			"err":           err,
		}
		err = queryTemplates12fb19["myDatastore__DeleteManyAuthors"].Execute(&query12fb19, templateInput12fb19)
		if err != nil {
			return
		}
		sqlQuery12fb19 = query12fb19.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteManyAuthors", sqlQuery12fb19, (*SQLGValues12fb19)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteManyAuthors", sqlQuery12fb19, (*SQLGValues12fb19)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteManyAuthors", err)
		}()
	}

	_, err = db.ExecContext(ctx, sqlQuery12fb19, (*SQLGValues12fb19)...)
	if err != nil {
		return
	}
	return
}
