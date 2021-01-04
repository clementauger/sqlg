//+build !sqlg

// Code generated by sqlg DO NOT EDIT

package pg

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/clementauger/sqlg/example/first/model"
	sqlg "github.com/clementauger/sqlg/runtime"
	tpl "github.com/clementauger/sqlg/tpl/pg"
	"text/template"
)

var queryTemplatesd36938 = map[string]*template.Template{
	"myDatastore__CreateSomeValues": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`{{$fields := fields $.SQLGConverter .v "id"}}
		{{$vals := fields $.SQLGConverter .v "id" "v"}}
		INSERT INTO sometype ( {{$fields | cols}} )
		VALUES ( {{$vals | vals $.SQLGValues $.SQLGFlavor}}, {{.v | pqArray | val $.SQLGValues $.SQLGFlavor}} )`,
	)),
	"myDatastore__DeleteManyAuthors": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`DELETE FROM authors WHERE id ANY ( {{.ids | pqArray | val $.SQLGValues $.SQLGFlavor}}::int[] ) RETURNING *`,
	)),
}

var rawQueriesd36938 = map[string]string{
	"myDatastore__DeleteAuthors": `DELETE FROM authors WHERE bio = ''`,
}

func (m *MyDatastore) CreateSomeValues(ctx context.Context, db sqlg.Execer, v model.SomeType) (id int64, err error) {
	var sqlQueryd36938 string
	SQLGValuesd36938 := &[]interface{}{}
	SQLGFlavord36938 := "$n"
	{
		var queryd36938 bytes.Buffer
		templateInputd36938 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValuesd36938,
			"SQLGFlavor":    SQLGFlavord36938,
			"v":             v,
			"id":            id,
			"err":           err,
		}
		err = queryTemplatesd36938["myDatastore__CreateSomeValues"].Execute(&queryd36938, templateInputd36938)
		if err != nil {
			return
		}
		sqlQueryd36938 = queryd36938.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "CreateSomeValues", sqlQueryd36938, (*SQLGValuesd36938)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "CreateSomeValues", sqlQueryd36938, (*SQLGValuesd36938)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "CreateSomeValues", err)
		}()
	}

	var resd36938 sql.Result
	resd36938, err = db.ExecContext(ctx, sqlQueryd36938, (*SQLGValuesd36938)...)
	if err != nil {
		return
	}
	id, err = resd36938.LastInsertId()
	if err != nil {
		return
	}
	return
}

func (m MyDatastore) DeleteAuthors(ctx context.Context, db sqlg.Execer) (err error) {
	var sqlQueryd36938 string
	sqlQueryd36938 = rawQueriesd36938["myDatastore__DeleteAuthors"]

	m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteAuthors", sqlQueryd36938)
	m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteAuthors", sqlQueryd36938)
	defer func() {
		m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteAuthors", err)
	}()

	_, err = db.ExecContext(ctx, sqlQueryd36938)
	if err != nil {
		return
	}
	return
}

func (m MyDatastore) DeleteManyAuthors(ctx context.Context, db sqlg.Querier, ids []int) (ab []model.Author, err error) {
	var sqlQueryd36938 string
	SQLGValuesd36938 := &[]interface{}{}
	SQLGFlavord36938 := "$n"
	{
		var queryd36938 bytes.Buffer
		templateInputd36938 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValuesd36938,
			"SQLGFlavor":    SQLGFlavord36938,
			"ids":           ids,
			"ab":            ab,
			"err":           err,
		}
		err = queryTemplatesd36938["myDatastore__DeleteManyAuthors"].Execute(&queryd36938, templateInputd36938)
		if err != nil {
			return
		}
		sqlQueryd36938 = queryd36938.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteManyAuthors", sqlQueryd36938, (*SQLGValuesd36938)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteManyAuthors", sqlQueryd36938, (*SQLGValuesd36938)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteManyAuthors", err)
		}()
	}

	var rowsd36938 *sql.Rows
	rowsd36938, err = db.QueryContext(ctx, sqlQueryd36938, (*SQLGValuesd36938)...)
	if err != nil {
		return
	}
	for rowsd36938.Next() {
		var itemd36938 model.Author
		err = rowsd36938.Scan(&itemd36938.ID, &itemd36938.Bio)
		if err != nil {
			return
		}
		ab = append(ab, itemd36938)
	}
	if err = rowsd36938.Close(); err != nil {
		return
	}
	err = rowsd36938.Err()
	return
}
