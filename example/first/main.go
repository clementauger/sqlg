//+build !sqlg

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/clementauger/sqlg/example/first/model"
	store "github.com/clementauger/sqlg/example/first/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

type defaultLogger struct{}

func (l defaultLogger) Log(typ, method, query string, args ...interface{}) {
	log.Printf("%v.%v: %q", typ, method, query)
}

func main() {
	os.Remove("sqlite-database.db")
	defer os.Remove("sqlite-database.db")

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	db, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer db.Close()
	createTable(db)
	store := store.MyDatastore{}
	store.Logger.Configure(defaultLogger{})

	ctx := context.Background()
	var id int64
	id, err = store.CreateAuthor(ctx, db, model.Author{Bio: "bio"})
	if err != nil {
		log.Fatalf("create author: %v", err)
	}
	fmt.Println("created author id #", id)
	authors, err := store.GetAuthors(ctx, db)
	if err != nil {
		log.Fatalf("get author: %v", err)
	}
	for _, a := range authors {
		fmt.Println(a)
	}
	it, err := store.GetAuthorsWihIterator(ctx, db, 1)
	if err != nil {
		log.Fatalf("get author: %v", err)
	}
	for it.Next() {
		fmt.Println(it.Value())
	}
	fmt.Println("err:", it.Err())

	it2, err := store.GetAuthorsWihIterator(ctx, db, 1)
	fmt.Println(it2.All())
	fmt.Println("err:", it.Err())

	a, err := store.GetAuthor2(ctx, db, 1)
	fmt.Println(a)
	fmt.Println("err:", err)

	a, err = store.GetAuthor3(ctx, db, 1)
	fmt.Println(a)
	fmt.Println("err:", err)

	a, err = store.GetAuthor3(ctx, db, 2)
	fmt.Println(a)
	fmt.Println("err:", err)

	authors, err = store.GetSomeAuthors(ctx, db, []int{0, 1}, 0, 10, "bio", "id")
	fmt.Println(authors)
	fmt.Println("err:", err)
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE authors (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		bio TEXT
	  );`

	statement, err := db.Prepare(createStudentTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("student table created")
}
