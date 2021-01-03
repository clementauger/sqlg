package runtime

type Tracer interface {
	Begin(typ, method, query string, args ...interface{})
	End(typ, method string, err error)
}

type NilTracer struct {
	Tracer
}

func (n *NilTracer) Configure(tracer Tracer) {
	n.Tracer = tracer
}
func (n NilTracer) Begin(typ, method, query string, args ...interface{}) {
	if n.Tracer != nil {
		n.Tracer.Begin(typ, method, query, args)
	}
}
func (n NilTracer) End(typ, method string, err error) {
	if n.Tracer != nil {
		n.Tracer.End(typ, method, err)
	}
}
