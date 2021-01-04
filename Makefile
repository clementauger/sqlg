test:
	go test github.com/clementauger/sqlg/tpl
	go test github.com/clementauger/sqlg/example/first/sqlite
	go test github.com/clementauger/sqlg/example/first/oracle
	go test github.com/clementauger/sqlg/example/first/mssql
	go test github.com/clementauger/sqlg/example/first/mysql
	go test github.com/clementauger/sqlg/example/first/pg
	(cd example/first; go run .)
gen:
	(cd example/first; go generate -x -tags=sqlg .)
