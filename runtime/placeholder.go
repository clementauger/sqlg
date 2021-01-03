package runtime

const (
	QuestionMark = "?"
	Dollar       = "$n"
	Named        = "@n"
)

const (
	MsSQL      = "mssql"
	MySQL      = "mysql"
	PostgreSQL = "pg"
	SQLite     = "sqlite"
	Oracle     = "oracle"
)

func EngineToPlaceholder(engine string) string {
	switch engine {
	case MySQL, SQLite, Oracle:
		return QuestionMark

	case PostgreSQL:
		return Dollar

	case MsSQL:
		return Named
	}
	return ""
}
