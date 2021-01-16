//+build !sqlg

// Code generated by sqlg DO NOT EDIT

package mysql

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/clementauger/sqlg/example/first/model"
	sqlg "github.com/clementauger/sqlg/runtime"
	"github.com/clementauger/sqlg/tpl"
	"text/template"
)

var queryTemplates410ea3 = map[string]*template.Template{
	"myDatastore__CreateAuthor": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`{{$fields := fields $.SQLGConverter .a "id"}}
		INSERT INTO authors ( {{$fields | cols}} ) VALUES ( {{$fields | vals $.SQLGValues $.SQLGFlavor .a}} )`,
	)),
	"myDatastore__CreateAuthor2": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`{{$fields := fields $.SQLGConverter .a "id"}}
		INSERT INTO authors
		( {{$fields | cols}} )
		VALUES
		( {{$fields | vals $.SQLGValues $.SQLGFlavor .a}} ) `,
	)),
	"myDatastore__CreateAuthors": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`INSERT INTO authors (bio)
		VALUES
		{{range $i, $a := .a}}
		 ( {{$a.Bio | val $.SQLGValues $.SQLGFlavor}} ) {{comma $i (len $.a)}}
		{{end}}
	`,
	)),
	"myDatastore__CreateAuthors2": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`{{$fields := fields $.SQLGConverter .a "id"}}
		INSERT INTO authors ( {{$fields | cols}} )
		VALUES
		{{range $i, $a := .a}}
		 ( {{$fields | vals $.SQLGValues $.SQLGFlavor $a}} ) {{comma $i (len $.a)}}
		{{end}}
	`,
	)),
	"myDatastore__CreateAuthors3": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`{{$fields := fields $.SQLGConverter .a "id"}}
		INSERT INTO authors ( {{$fields | cols}} )
		VALUES
		{{range $i, $a := .a}}
		 ( {{$fields | vals $.SQLGValues $.SQLGFlavor $a}} ) {{comma $i (len $.a)}}
		{{end}}
	`,
	)),
	"myDatastore__DeleteAuthor": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`DELETE FROM authors WHERE id={{.id | val $.SQLGValues $.SQLGFlavor}}`,
	)),
	"myDatastore__DeleteAuthor2": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`DELETE FROM authors WHERE id={{.id | val $.SQLGValues $.SQLGFlavor}}`,
	)),
	"myDatastore__GetAuthor": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`SELECT * FROM authors WHERE id={{.id | val $.SQLGValues $.SQLGFlavor}}`,
	)),
	"myDatastore__GetAuthor2": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`{{$fields := fields $.SQLGConverter .a}}
		SELECT {{$fields | cols}} FROM authors WHERE id={{.id | val $.SQLGValues $.SQLGFlavor}}`,
	)),
	"myDatastore__GetAuthor3": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`{{$fields := fields $.SQLGConverter .a}}
		SELECT {{$fields | cols | prefix "alias."}} FROM authors as alias WHERE alias.id={{.id | val $.SQLGValues $.SQLGFlavor}}`,
	)),
	"myDatastore__GetAuthorCount": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`SELECT *, COUNT(*) as count FROM authors WHERE id={{.id | val $.SQLGValues $.SQLGFlavor}}`,
	)),
	"myDatastore__GetAuthorsWihIterator": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`SELECT * FROM authors WHERE id={{.id | val $.SQLGValues $.SQLGFlavor}}`,
	)),
	"myDatastore__GetAuthorsWihNamedIterator": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`SELECT * FROM authors WHERE id={{.id | val $.SQLGValues $.SQLGFlavor}}`,
	)),
	"myDatastore__GetSomeAuthors": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`SELECT * FROM authors
		WHERE id IN ({{.ids | val $.SQLGValues $.SQLGFlavor}})
		GROUP BY {{.groupby | print}}
		ORDER BY {{.orderby | raw}}
		LIMIT {{.start | val $.SQLGValues $.SQLGFlavor}}, {{.end | val $.SQLGValues $.SQLGFlavor}}
		`,
	)),
	"myDatastore__GetSomeY": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`SELECT * FROM y`,
	)),
	"myDatastore__UpdateAuthor": template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
		`{{$fields := fields $.SQLGConverter .a "id"}}
		UPDATE authors SET
		 {{$fields | update $.SQLGValues $.SQLGFlavor .a}}
		 WHERE id = {{.a.ID | val $.SQLGValues $.SQLGFlavor}}`,
	)),
}

var rawQueries410ea3 = map[string]string{
	"myDatastore__GetAuthors": `SELECT * FROM authors`,
	"myDatastore__ProductUpdate": `UPDATE products SET price = price * 1.10
  WHERE price <= 99.99
  RETURNING name, price AS new_price`,
}

// MyDatastore stores stuff.
type MyDatastore struct {
	Tracer    sqlg.NilTracer
	Logger    sqlg.NilLogger
	Converter sqlg.ToSnake
}

// MyDatastoreIface is an interface of MyDatastore
type MyDatastoreIface interface {
	CreateAuthor(ctx context.Context, db sqlg.Execer, a model.Author) (id int64, err error)
	CreateAuthor2(ctx context.Context, db sqlg.Execer, a model.Author) (id int64, err error)
	CreateAuthors(ctx context.Context, db sqlg.Execer, a []model.Author) (err error)
	CreateAuthors2(ctx context.Context, db sqlg.Execer, a []model.Author) (err error)
	CreateAuthors3(ctx context.Context, db sqlg.Execer, a []model.Author) (err error)
	DeleteAuthor(ctx context.Context, db sqlg.Execer, id int) (err error)
	DeleteAuthor2(ctx context.Context, db sqlg.Execer, id int) (count int64, err error)
	GetAuthor(ctx context.Context, db sqlg.Querier, id int) (a model.Author, err error)
	GetAuthor2(ctx context.Context, db sqlg.Querier, id int) (a model.Author, err error)
	GetAuthor3(ctx context.Context, db sqlg.Querier, id int) (a model.Author, err error)
	GetAuthorCount(ctx context.Context, db sqlg.Querier, id int) (a model.AuthorCount, err error)
	GetAuthors(ctx context.Context, db sqlg.Querier) (a []model.Author, err error)
	GetAuthorsWihIterator(ctx context.Context, db sqlg.Querier, id int) (it AuthorIterator, err error)
	GetAuthorsWihNamedIterator(ctx context.Context, db sqlg.Querier, id int) (it AuthorIterator, err error)
	GetSomeAuthors(ctx context.Context, db sqlg.Querier, ids []int, start int, end int, orderby string, groupby string) (param0 []model.Author, err error)
	GetSomeY(ctx context.Context, db sqlg.Querier, u model.Y) (param0 []model.Y, err error)
	ProductUpdate(ctx context.Context, db sqlg.Querier) (name string, price int, err error)
	UpdateAuthor(ctx context.Context, db sqlg.Execer, a model.Author) (err error)
	CreateSomeValues(ctx context.Context, db sqlg.Execer, v model.SomeType) (id int64, err error)
	CreateTable(ctx context.Context, db sqlg.Execer) (err error)
	DeleteAuthors(ctx context.Context, db sqlg.Execer) (err error)
	DeleteManyAuthors(ctx context.Context, db sqlg.Execer, ids []int) (_ []model.Author, err error)
}

// AuthorIterator is an iterator of Author
type AuthorIterator struct {
	rows  *sql.Rows
	err   error
	value model.Author
}

func (x AuthorIterator) Err() error {
	if x.err != nil {
		return x.err
	}
	return x.rows.Err()
}

func (x AuthorIterator) Close() error {
	return x.rows.Close()
}

func (x AuthorIterator) All() (ret []model.Author) {
	for x.Next() {
		ret = append(ret, x.Value())
	}
	return ret
}
func (it *AuthorIterator) Value() model.Author {
	return it.value
}
func (it *AuthorIterator) Next() bool {
	if !it.rows.Next() {
		it.rows.Close()
		return false
	}
	it.err = it.rows.Scan(&(it.value.ID), &(it.value.Bio))
	if it.err == nil {
		return true
	}
	return false
}

func (m *MyDatastore) CreateAuthor(ctx context.Context, db sqlg.Execer, a model.Author) (id int64, err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"a":             a,
			"id":            id,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__CreateAuthor"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthor", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthor", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthor", err)
		}()
	}

	var res410ea3 sql.Result
	res410ea3, err = db.ExecContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	id, err = res410ea3.LastInsertId()
	if err != nil {
		return
	}
	return
}

func (m *MyDatastore) CreateAuthor2(ctx context.Context, db sqlg.Execer, a model.Author) (id int64, err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"a":             a,
			"id":            id,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__CreateAuthor2"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthor2", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthor2", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthor2", err)
		}()
	}

	var res410ea3 sql.Result
	res410ea3, err = db.ExecContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	id, err = res410ea3.LastInsertId()
	if err != nil {
		return
	}
	return
}

func (m *MyDatastore) CreateAuthors(ctx context.Context, db sqlg.Execer, a []model.Author) (err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"a":             a,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__CreateAuthors"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthors", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthors", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthors", err)
		}()
	}

	_, err = db.ExecContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	return
}

func (m *MyDatastore) CreateAuthors2(ctx context.Context, db sqlg.Execer, a []model.Author) (err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"a":             a,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__CreateAuthors2"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthors2", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthors2", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthors2", err)
		}()
	}

	_, err = db.ExecContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	return
}

func (m *MyDatastore) CreateAuthors3(ctx context.Context, db sqlg.Execer, a []model.Author) (err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"a":             a,
			"err":           err,
			"b":             model.Author{},
		}
		err = queryTemplates410ea3["myDatastore__CreateAuthors3"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthors3", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthors3", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "CreateAuthors3", err)
		}()
	}

	_, err = db.ExecContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	return
}

func (m *MyDatastore) DeleteAuthor(ctx context.Context, db sqlg.Execer, id int) (err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"id":            id,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__DeleteAuthor"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteAuthor", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteAuthor", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteAuthor", err)
		}()
	}

	_, err = db.ExecContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	return
}

func (m *MyDatastore) DeleteAuthor2(ctx context.Context, db sqlg.Execer, id int) (count int64, err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"id":            id,
			"count":         count,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__DeleteAuthor2"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteAuthor2", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteAuthor2", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "DeleteAuthor2", err)
		}()
	}

	var res410ea3 sql.Result
	res410ea3, err = db.ExecContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	count, err = res410ea3.RowsAffected()
	if err != nil {
		return
	}
	return
}

// GetAuthor retrieves
// an Author by its ID.
func (m MyDatastore) GetAuthor(ctx context.Context, db sqlg.Querier, id int) (a model.Author, err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"id":            id,
			"a":             a,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__GetAuthor"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor", err)
		}()
	}

	var rows410ea3 *sql.Rows
	rows410ea3, err = db.QueryContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	for rows410ea3.Next() {
		err = rows410ea3.Scan(&a.ID, &a.Bio)
		if err != nil {
			return
		}
	}
	if err = rows410ea3.Close(); err != nil {
		return
	}
	err = rows410ea3.Err()
	return
}

func (m MyDatastore) GetAuthor2(ctx context.Context, db sqlg.Querier, id int) (a model.Author, err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"id":            id,
			"a":             a,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__GetAuthor2"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor2", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor2", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor2", err)
		}()
	}

	var rows410ea3 *sql.Rows
	rows410ea3, err = db.QueryContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	for rows410ea3.Next() {
		err = rows410ea3.Scan(&a.ID, &a.Bio)
		if err != nil {
			return
		}
	}
	if err = rows410ea3.Close(); err != nil {
		return
	}
	err = rows410ea3.Err()
	return
}

func (m MyDatastore) GetAuthor3(ctx context.Context, db sqlg.Querier, id int) (a model.Author, err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"id":            id,
			"a":             a,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__GetAuthor3"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor3", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor3", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor3", err)
		}()
	}

	var rows410ea3 *sql.Rows
	rows410ea3, err = db.QueryContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	for rows410ea3.Next() {
		err = rows410ea3.Scan(&a.ID, &a.Bio)
		if err != nil {
			return
		}
	}
	if err = rows410ea3.Close(); err != nil {
		return
	}
	err = rows410ea3.Err()
	return
}

// GetAuthorCount retrieves
// Author and count.
func (m MyDatastore) GetAuthorCount(ctx context.Context, db sqlg.Querier, id int) (a model.AuthorCount, err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"id":            id,
			"a":             a,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__GetAuthorCount"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthorCount", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthorCount", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthorCount", err)
		}()
	}

	var rows410ea3 *sql.Rows
	rows410ea3, err = db.QueryContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	for rows410ea3.Next() {
		err = rows410ea3.Scan(&a.ID, &a.Bio, &a.Count)
		if err != nil {
			return
		}
	}
	if err = rows410ea3.Close(); err != nil {
		return
	}
	err = rows410ea3.Err()
	return
}

func (m *MyDatastore) GetAuthors(ctx context.Context, db sqlg.Querier) (a []model.Author, err error) {
	var sqlQuery410ea3 string
	sqlQuery410ea3 = rawQueries410ea3["myDatastore__GetAuthors"]

	m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthors", sqlQuery410ea3)
	m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthors", sqlQuery410ea3)
	defer func() {
		m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthors", err)
	}()

	var rows410ea3 *sql.Rows
	rows410ea3, err = db.QueryContext(ctx, sqlQuery410ea3)
	if err != nil {
		return
	}
	for rows410ea3.Next() {
		var item410ea3 model.Author
		err = rows410ea3.Scan(&item410ea3.ID, &item410ea3.Bio)
		if err != nil {
			return
		}
		a = append(a, item410ea3)
	}
	if err = rows410ea3.Close(); err != nil {
		return
	}
	err = rows410ea3.Err()
	return
}

func (m MyDatastore) GetAuthorsWihIterator(ctx context.Context, db sqlg.Querier, id int) (it AuthorIterator, err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"id":            id,
			"it":            it,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__GetAuthorsWihIterator"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthorsWihIterator", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthorsWihIterator", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthorsWihIterator", err)
		}()
	}

	var rows410ea3 *sql.Rows
	rows410ea3, err = db.QueryContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	it.rows = rows410ea3
	return
}

func (m MyDatastore) GetAuthorsWihNamedIterator(ctx context.Context, db sqlg.Querier, id int) (it AuthorIterator, err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"id":            id,
			"it":            it,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__GetAuthorsWihNamedIterator"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthorsWihNamedIterator", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthorsWihNamedIterator", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthorsWihNamedIterator", err)
		}()
	}

	var rows410ea3 *sql.Rows
	rows410ea3, err = db.QueryContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	it.rows = rows410ea3
	return
}

func (m *MyDatastore) GetSomeAuthors(ctx context.Context, db sqlg.Querier, ids []int, start int, end int, orderby string, groupby string) (param0 []model.Author, err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"ids":           ids,
			"start":         start,
			"end":           end,
			"orderby":       orderby,
			"groupby":       groupby,
			"param0":        param0,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__GetSomeAuthors"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "GetSomeAuthors", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "GetSomeAuthors", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "GetSomeAuthors", err)
		}()
	}

	var rows410ea3 *sql.Rows
	rows410ea3, err = db.QueryContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	for rows410ea3.Next() {
		var item410ea3 model.Author
		err = rows410ea3.Scan(&item410ea3.ID, &item410ea3.Bio)
		if err != nil {
			return
		}
		param0 = append(param0, item410ea3)
	}
	if err = rows410ea3.Close(); err != nil {
		return
	}
	err = rows410ea3.Err()
	return
}

func (m *MyDatastore) GetSomeY(ctx context.Context, db sqlg.Querier, u model.Y) (param0 []model.Y, err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"u":             u,
			"param0":        param0,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__GetSomeY"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "GetSomeY", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "GetSomeY", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "GetSomeY", err)
		}()
	}

	var rows410ea3 *sql.Rows
	rows410ea3, err = db.QueryContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	for rows410ea3.Next() {
		var item410ea3 model.Y
		err = rows410ea3.Scan(&item410ea3.W)
		if err != nil {
			return
		}
		param0 = append(param0, item410ea3)
	}
	if err = rows410ea3.Close(); err != nil {
		return
	}
	err = rows410ea3.Err()
	return
}

func (m *MyDatastore) ProductUpdate(ctx context.Context, db sqlg.Querier) (name string, price int, err error) {
	var sqlQuery410ea3 string
	sqlQuery410ea3 = rawQueries410ea3["myDatastore__ProductUpdate"]

	m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "ProductUpdate", sqlQuery410ea3)
	m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "ProductUpdate", sqlQuery410ea3)
	defer func() {
		m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "ProductUpdate", err)
	}()

	var rows410ea3 *sql.Rows
	rows410ea3, err = db.QueryContext(ctx, sqlQuery410ea3)
	if err != nil {
		return
	}
	for rows410ea3.Next() {
		err = rows410ea3.Scan(&name, &price)
		if err != nil {
			return
		}
	}
	if err = rows410ea3.Close(); err != nil {
		return
	}
	err = rows410ea3.Err()
	return
}

func (m *MyDatastore) UpdateAuthor(ctx context.Context, db sqlg.Execer, a model.Author) (err error) {
	var sqlQuery410ea3 string
	SQLGValues410ea3 := &[]interface{}{}
	SQLGFlavor410ea3 := "?"
	{
		var query410ea3 bytes.Buffer
		templateInput410ea3 := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues410ea3,
			"SQLGFlavor":    SQLGFlavor410ea3,
			"a":             a,
			"err":           err,
		}
		err = queryTemplates410ea3["myDatastore__UpdateAuthor"].Execute(&query410ea3, templateInput410ea3)
		if err != nil {
			return
		}
		sqlQuery410ea3 = query410ea3.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "UpdateAuthor", sqlQuery410ea3, (*SQLGValues410ea3)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "UpdateAuthor", sqlQuery410ea3, (*SQLGValues410ea3)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "UpdateAuthor", err)
		}()
	}

	_, err = db.ExecContext(ctx, sqlQuery410ea3, (*SQLGValues410ea3)...)
	if err != nil {
		return
	}
	return
}
