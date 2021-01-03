package runtime

type Logger interface {
	Log(typ, method, query string, args ...interface{})
}

type NilLogger struct {
	Logger
}

func (n *NilLogger) Configure(logger Logger) {
	n.Logger = logger
}
func (n NilLogger) Log(typ, method, query string, args ...interface{}) {
	if n.Logger != nil {
		n.Logger.Log(typ, method, query)
	}
}
