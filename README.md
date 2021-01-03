# sqlg

`sqlg` is a command line utility to generate boilerplate sql instructions when
using go `sql` package.

It parses `.go` file which contains such kind of functions

```go
//+build sqlg

//go:generate go run github.com/clementauger/sqlg/cmd/sqlg -clean sqlite
//go:generate go run github.com/clementauger/sqlg/cmd/sqlg -engine sqlite

// myDatastore stores stuff.
type myDatastore struct {
	sqlg.SQLg
}

// GetAuthor retrieves an Author by its ID.
func (m myDatastore) GetAuthor(id int) (a model.Author, err error) {
	m.Query(`SELECT * FROM authors WHERE id={{.id}}`)
	return
}
```

And generates the corresponding go code using `go generate -x -tags=sqlg .`

```go
// GetAuthor retrieves an Author by its ID.
func (m MyDatastore) GetAuthor(ctx context.Context, db sqlg.Querier, id int) (a model.Author, err error) {
	var sqlQuery86120a string
	SQLGValues86120a := &[]interface{}{}
	SQLGFlavor86120a := "?"
	{
		var query86120a bytes.Buffer
		templateInput86120a := map[string]interface{}{
			"SQLGConverter": m.Converter,
			"SQLGValues":    SQLGValues86120a,
			"SQLGFlavor":    SQLGFlavor86120a,
			"id":            id,
			"a":             a,
			"err":           err,
		}
		err = queryTemplates86120a["myDatastore__GetAuthor"].Execute(&query86120a, templateInput86120a)
		if err != nil {
			return
		}
		sqlQuery86120a = query86120a.String()

		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor", sqlQuery86120a, (*SQLGValues86120a)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor", sqlQuery86120a, (*SQLGValues86120a)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor", err)
		}()
	}

	var rows86120a *sql.Rows
	rows86120a, err = db.QueryContext(ctx, sqlQuery86120a, (*SQLGValues86120a)...)
	if err != nil {
		return
	}
	for rows86120a.Next() {
		err = rows86120a.Scan(&a.ID, &a.Bio)
		if err != nil {
			return
		}
	}
	if err = rows86120a.Close(); err != nil {
		return
	}
	err = rows86120a.Err()
	return
}
```

The sql query becomes

```sql
SELECT * FROM authors WHERE id=?
```

It attempts to provide support around the work of writing maintenable sql queries with maximum control.

Queries are parsed as `go/templates` at runtime to generate an appropriate sql queries using some helpers.

```go
func (m *myDatastore) CreateAuthor(a model.Author) (id int64, err error) {
	m.Exec(`INSERT INTO authors ( {{cols .a "id"}} ) VALUES ( {{vals .a "id"}} )`).InsertedID(id)
	return
}
```

It becomes

```sql
INSERT INTO authors ( bio ) VALUES ( ? )
```


When a query is templated, it recieves input and output function parameters as a `map[string]interface{}`
  so you can use those values to build the sql output.

It comes with functions like `cols(someStructValue interface{}, notFields... string)` to map a struct properties list into their corrsponding sql columns.
  Conversevely `{{vals .a "id"}}` is an helper to list values of a struct properties list except some fields.
  Those values are recorded to be passed as query arguments when invoking `db.Query` or `db.Exec` methods.
  Those values are printed with their corresponding placeholder syntax within the query.

It tries to be useful with some helpers like `comma(index, max)` `prefix(string, colPrinting)`

```go
func (m *myDatastore) CreateAuthors(a []model.Author) (err error) {
	m.Exec(`INSERT INTO authors as alias ({{cols .a "id" | prefix "alias"}})
		VALUES
		{{range $i, $a := .a}}
		 ( {{$a.Bio}} ) {{comma $i $a}}
		{{end}}
	`)
	return
}
```

When you need special syntax per engine, use build tags

This is a query written for `postgresql`

```go
// file: demo.go

//+build sqlg,pg

package main

import "github.com/clementauger/sqlg/example/first/model"

func (m myDatastore) DeleteManyAuthors(ids []int) (ab []model.Author, err error) {
	m.Query(`DELETE FROM authors WHERE id ANY ( {{.ids | pqArray}}::int[] ) RETURNING *`)
	return
}
```

Postgresql is provided a special function `pqArray` to emit a `pq.Array` value.

When using `mssql` engine, values are automatically converted to `sql.NamedArg` and the indexed placeholder are generated at runtime.

so the query become

```sql
INSERT INTO authors ( id,bio ) VALUES ( @p0,@p1 )
```

Some specific functions are available to cast the value to an `sql.Valuer` with template helpers such

```go
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
```

Other engines are supported (`mysql oracle sqlite`).

in below examples the same function signature is used for different engines that can t provide exactly the same functionality.

```go
// file: demo_others.go

//+build sqlg
//+build mysql oracle sqlite

package main

import (
	"fmt"
	"github.com/clementauger/sqlg/example/first/model"
)

func (m *myDatastore) DeleteManyAuthors(ids []int) (_ []model.Author, err error) {
	m.Exec(`DELETE FROM authors WHERE id IN ( {{.ids}} )`)
	return
}
```

When needed, you can mark a method as unsupported

```go
//+build sqlg
//+build mysql oracle sqlite

func (m *myDatastore) CreateSomeValues(v model.SomeType) (id int64, err error) {
	return id, fmt.Errorf("unsupported")
}
```

Only types embedding an `sql.SQLg` interface will be processed as query handlers.

```go
// myDatastore stores stuff.
type myDatastore struct {
	sqlg.SQLg
}
```

You can enable tracing or logging by adding properties to your struct definition


```go
// myDatastore stores stuff.
type myDatastore struct {
	sqlg.SQLg
	Tracer    sqlg.NilTracer
	Logger    sqlg.NilLogger
	Converter sqlg.ToSnake
}
```

it adds corresponding expressions at runtime

```go
		m.Logger.Log("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor", sqlQuery86120a, (*SQLGValues86120a)...)
		m.Tracer.Begin("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor", sqlQuery86120a, (*SQLGValues86120a)...)
		defer func() {
			m.Tracer.End("github.com/clementauger/sqlg/example/first/myDatastore", "GetAuthor", err)
		}()
```

If can also generates and return query iterators, useful when working with large data sets

```go
func (m myDatastore) GetAuthorsWihIterator(id int) (it func() (model.Author, error), err error) {
	m.Query(`SELECT * FROM authors WHERE id={{.id}}`)
	return
}
```

The generated iterator exposes the method `Value() model.Author `, use it like so

```go
it, err := store.GetAuthorsWihIterator(ctx, db, 1)
if err != nil {
	log.Fatalf("get author: %v", err)
}
for it.Next() {
	fmt.Println(it.Value())
}
fmt.Println("err:", it.Err())
```

You write test using the neat http://github.com/DATA-DOG/go-sqlmock api

```go

import (
	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateAuthorWithMock(t *testing.T) {
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
```

Consider to put your models into a dedicated folder

Check the [example](https://github.com/clementauger/sqlg/tree/main/example/first) folder.

# Install

```sh
go get github.com/clementauger/sqlg/...
```
