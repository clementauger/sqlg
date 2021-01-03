package tpl

import (
	"database/sql"
	"reflect"
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/clementauger/sqlg/runtime"

	"fmt"
)

type collectedValues interface{}

// CommonFuncMap is the set of function used to transform the templates.
var CommonFuncMap = map[string]interface{}{
	// for your convenience
	"comma": func(index, max int) string {
		if index < max-1 {
			return ","
		}
		return ""
	},

	// sqlg internals
	"collect": func(values *[]interface{}, placeholder string, s interface{}) interface{} {
		r := reflect.ValueOf(s)
		vs := *values
		vl := len(vs)
		if r.Kind() == reflect.Slice {
			for i := 0; i < r.Len(); i++ {
				if placeholder == runtime.Named {
					vs = append(vs, sql.Named(fmt.Sprintf("p%v", vl+i), r.Index(i).Interface()))
				} else {
					vs = append(vs, r.Index(i).Interface())
				}
			}
			*values = vs
			return make([]collectedValues, r.Len())
		}
		if placeholder == runtime.Named {
			vs = append(vs, sql.Named(fmt.Sprintf("p%v", vl), s))
		} else {
			vs = append(vs, s)
		}
		*values = vs
		return ""
	},
	"placeholder": func(values *[]interface{}, placeholder string, s interface{}) string {
		if cvalues, ok := s.([]collectedValues); ok {
			out := ""
			for i := range cvalues {
				if placeholder == runtime.Dollar {
					out += fmt.Sprintf("$%v,", len(*values)-len(cvalues)+i)
				} else if placeholder == runtime.Named {
					out += fmt.Sprintf("@p%v,", len(*values)-len(cvalues)+i)
				} else if placeholder == runtime.QuestionMark {
					out += "?,"
				} else {
					panic("no such placeholder")
				}
			}
			out = strings.TrimSuffix(out, ",")
			return out
		}
		if placeholder == runtime.Dollar {
			return fmt.Sprintf("$%v", len(*values))
		} else if placeholder == runtime.Named {
			return fmt.Sprintf("@p%v", len(*values))
		} else if placeholder == runtime.QuestionMark {
			return "?"
		}
		panic("no such placeholder")
	},
	"vals": func(converter runtime.CaseConverter, s interface{}, notFields ...string) (ret []interface{}) {
		r := reflect.ValueOf(s)
		if r.Kind() == reflect.Struct {
			for i := 0; i < r.NumField(); i++ {
				f := r.Type().Field(i)
				sf := f.Name
				if converter != nil {
					sf = converter.Convert(sf)
				}
				var ok bool = true
				for _, not := range notFields {
					if not == sf {
						ok = false
						break
					}
				}
				if !ok {
					continue
				}
				ret = append(ret, r.Field(i).Interface())
			}
		}
		return
	},

	// colunm printing
	"fields": func(converter runtime.CaseConverter, v interface{}, notSQLFields ...string) []fieldAndValue {
		fields := []fieldAndValue{}
		r := reflect.ValueOf(v)
		for i := 0; i < r.Type().NumField(); i++ {
			f := r.Type().Field(i)
			sf := f.Name
			if converter != nil {
				sf = converter.Convert(sf)
			}
			var ignore bool
			for _, not := range notSQLFields {
				if sf == not {
					ignore = true
					break
				}
			}
			if !ignore {
				fields = append(fields, fieldAndValue{
					Prop:  f.Name,
					SQL:   sf,
					Value: r.Field(i).Interface(),
				})
			}
		}
		return fields
	},
	"cols": func(v interface{}, notSQLFields ...string) colPrinting {
		return colPrinting{v: v, not: notSQLFields}
	},
	"convert": func(converter runtime.CaseConverter, p colPrinting) colPrinting {
		p.converter = converter
		return p
	},
	"prefix": func(prefix string, p colPrinting) colPrinting {
		p.alias = prefix
		return p
	},
	"glue": func(glue string, p colPrinting) colPrinting {
		p.glue = glue
		return p
	},

	// raw print
	"raw": func(x interface{}) interface{} { return x },
}

type fieldAndValue struct {
	Prop  string
	SQL   string
	Value interface{}
}

type colPrinting struct {
	v         interface{}
	not       []string
	alias     string
	glue      string
	converter runtime.CaseConverter
}

func (c colPrinting) String() string {
	if c.v == nil {
		return ""
	}
	var cols []string
	r := reflect.ValueOf(c.v)
	for i := 0; i < r.Type().NumField(); i++ {
		f := r.Type().Field(i)
		sf := f.Name
		if c.converter != nil {
			sf = c.converter.Convert(sf)
		}
		var ignore bool
		for _, not := range c.not {
			if sf == not {
				ignore = true
				break
			}
		}
		if !ignore {
			if c.alias != "" {
				sf = c.alias + sf
			}
			cols = append(cols, sf)
		}
	}
	return strings.TrimSuffix(strings.Join(cols, c.glue), c.glue)
}

// FuncMap returns the CommonFuncMap.
func FuncMap() template.FuncMap {
	return CommonFuncMap
}

// Transform a given template to add the sequence
// '| collect $.SQLGValues | placeholder $.SQLGFlavor'
// to each ActionNode which does not contain a 'comma' identifier.
// 'collect' records values into .SQLGValues to pass them
// to the Querier. 'placeholder' identifier emits a placeholder
// for each value to apply to the query.
func Transform(src string, funcs template.FuncMap) (string, error) {

	t, err := template.New("").Funcs(funcs).Parse(src)
	if err != nil {
		return "", err
	}

	hasIdentifier := func(out *bool, idents ...string) func(n parse.Node) bool {
		return func(n parse.Node) bool {
			if x, ok := n.(*parse.IdentifierNode); ok {
				for _, i := range idents {
					if x.String() == i {
						*out = true
						break
					}
				}
			}
			return true
		}
	}

	var transformValuePrintings func(n parse.Node) bool
	transformValuePrintings = func(n parse.Node) bool {
		if x, ok := n.(*parse.ActionNode); ok {

			var hasIdent bool
			visit(x, hasIdentifier(&hasIdent, "comma", "cols", "raw", "print", "printf", "fields"))

			if !hasIdent {
				collect := &parse.IdentifierNode{
					Ident: "collect",
				}
				values := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGValues",
					},
				}
				placeholder := &parse.IdentifierNode{
					Ident: "placeholder",
				}
				flavor := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGFlavor",
					},
				}

				cmd1 := &parse.CommandNode{}
				cmd1.Args = append(cmd1.Args, collect, values, flavor)
				cmd2 := &parse.CommandNode{}
				cmd2.Args = append(cmd2.Args, placeholder, values, flavor)

				x.Pipe.Cmds = append(x.Pipe.Cmds, cmd1, cmd2)
			}

		}
		return true
	}
	visit(t.Tree.Root, transformValuePrintings)

	var transformColPrintings func(n parse.Node) bool
	transformColPrintings = func(n parse.Node) bool {
		if x, ok := n.(*parse.ActionNode); ok {

			var shouldInclude bool
			visit(x, hasIdentifier(&shouldInclude, "cols"))

			if shouldInclude {

				var hasGlue bool
				visit(x, hasIdentifier(&hasGlue, "glue"))

				convert := &parse.IdentifierNode{
					Ident: "convert",
				}
				converter := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGConverter",
					},
				}
				cmd1 := &parse.CommandNode{}
				cmd1.Args = append(cmd1.Args, convert, converter)
				x.Pipe.Cmds = append(x.Pipe.Cmds, cmd1)

				if !hasGlue {
					glue := &parse.IdentifierNode{
						Ident: "glue",
					}
					comma := &parse.StringNode{
						Quoted: fmt.Sprintf("%q", ","),
						Text:   ",",
					}
					cmd2 := &parse.CommandNode{}
					cmd2.Args = append(cmd2.Args, glue, comma)
					x.Pipe.Cmds = append(x.Pipe.Cmds, cmd2)
				}

			}
		}
		return true
	}
	visit(t.Tree.Root, transformColPrintings)

	var transformValsPrintings func(n parse.Node) bool
	transformValsPrintings = func(n parse.Node) bool {
		if x, ok := n.(*parse.ActionNode); ok {

			var shouldInclude bool
			visit(x, hasIdentifier(&shouldInclude, "vals"))

			if shouldInclude {

				converter := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGConverter",
					},
				}
				t := append([]parse.Node{}, x.Pipe.Cmds[0].Args[1:]...)
				x.Pipe.Cmds[0].Args = append(x.Pipe.Cmds[0].Args[:1], converter)
				x.Pipe.Cmds[0].Args = append(x.Pipe.Cmds[0].Args, t...)

			}
		}
		return true
	}
	visit(t.Tree.Root, transformValsPrintings)

	var transformFields func(n parse.Node) bool
	transformFields = func(n parse.Node) bool {
		// log.Printf("%T %v\n", n, n)

		if x, ok := n.(*parse.PipeNode); ok {

			var shouldInclude bool
			visit(x, hasIdentifier(&shouldInclude, "fields"))

			if shouldInclude {
				converter := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGConverter",
					},
				}
				t := append([]parse.Node{}, x.Cmds[0].Args[1:]...)
				x.Cmds[0].Args = append(x.Cmds[0].Args[:1], converter)
				x.Cmds[0].Args = append(x.Cmds[0].Args, t...)
			}
		} else if x, ok := n.(*parse.CommandNode); ok {

			var shouldInclude bool
			visit(x, hasIdentifier(&shouldInclude, "fields"))

			if shouldInclude {
				converter := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGConverter",
					},
				}
				t := append([]parse.Node{}, x.Args[1:]...)
				x.Args = append(x.Args[:1], converter)
				x.Args = append(x.Args, t...)
			}
		}
		return true
	}
	visit(t.Tree.Root, transformFields)

	return t.Tree.Root.String(), nil
}

func visit(n parse.Node, fn func(parse.Node) bool) bool {
	if n == nil {
		return true
	}
	if !fn(n) {
		return false
	}
	if l, ok := n.(*parse.ListNode); ok {
		for _, nn := range l.Nodes {
			if !visit(nn, fn) {
				continue
			}
		}
	}
	if l, ok := n.(*parse.RangeNode); ok {
		visit(l.BranchNode.Pipe, fn)
		if l.BranchNode.List != nil {
			visit(l.BranchNode.List, fn)
		}
		if l.BranchNode.ElseList != nil {
			visit(l.BranchNode.ElseList, fn)
		}
	}
	if l, ok := n.(*parse.ActionNode); ok {
		for _, c := range l.Pipe.Decl {
			visit(c, fn)
		}
		for _, c := range l.Pipe.Cmds {
			if visit(c, fn) {
				for _, a := range c.Args {
					visit(a, fn)
				}
			}
		}
	}
	if l, ok := n.(*parse.CommandNode); ok {
		for _, a := range l.Args {
			visit(a, fn)
		}
	}
	if l, ok := n.(*parse.PipeNode); ok {
		for _, a := range l.Decl {
			visit(a, fn)
		}
		for _, a := range l.Cmds {
			visit(a, fn)
		}
	}
	return true
}
