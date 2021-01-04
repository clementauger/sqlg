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

	// colunm and values printing
	"fields": func(converter runtime.CaseConverter, v interface{}, notSQLFields ...string) fieldAndValues {
		fields := fieldAndValues{}
		r := reflect.ValueOf(v)
		if r.Kind() == reflect.Ptr {
			r = r.Elem()
		} else if r.Kind() == reflect.Slice {
			if r.Len() == 0 {
				r = reflect.Zero(r.Type().Elem())
			} else {
				r = r.Index(0)
			}
		}

		if r.Kind() != reflect.Struct {
			r = r.Elem()
		}
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
	"cols": func(fields fieldAndValues) fieldPrinter {
		return fieldPrinter{kind: "col", fields: fields}
	},
	"val": func(collected *[]interface{}, placeholder string, v interface{}) fieldPrinter {
		return fieldPrinter{kind: "val", v: v, collected: collected, placeholder: placeholder}
	},
	"vals": func(collected *[]interface{}, placeholder string, v interface{}, fields fieldAndValues) fieldPrinter {
		return fieldPrinter{kind: "fieldval", v: v, fields: fields, collected: collected, placeholder: placeholder}
	},
	"update": func(collected *[]interface{}, placeholder string, v interface{}, fields fieldAndValues) fieldPrinter {
		return fieldPrinter{kind: "update", v: v, fields: fields, collected: collected, placeholder: placeholder}
	},
	// helpers to configure printing
	"prefix": func(prefix string, p fieldPrinter) fieldPrinter {
		p.alias = prefix
		return p
	},
	"placeholder": func(placeholder string, p fieldPrinter) fieldPrinter {
		p.placeholder = placeholder
		return p
	},
	"glue": func(glue string, p fieldPrinter) fieldPrinter {
		p.glue = glue
		return p
	},

	// raw print
	"raw": func(x interface{}) interface{} { return x },
}

type fieldAndValues []fieldAndValue

type fieldAndValue struct {
	Prop  string
	SQL   string
	Value interface{}
}

type fieldPrinter struct {
	kind        string
	fields      fieldAndValues
	v           interface{}
	alias       string
	glue        string
	placeholder string
	collected   *[]interface{}
}

func (c fieldPrinter) String() string {
	if c.kind == "col" {
		return c.printCols()
	}
	if c.kind == "fieldval" {
		return c.printFieldValues()
	}
	if c.kind == "update" {
		return c.printUpdate()
	}
	return c.printValue()
}

func (c fieldPrinter) printCols() string {
	var cols []string
	for _, f := range c.fields {
		s := f.SQL
		if c.alias != "" {
			s = c.alias + s
		}
		cols = append(cols, s)
	}

	glue := c.glue
	if c.glue == "" {
		glue = ","
	}
	return strings.TrimSuffix(strings.Join(cols, glue), glue)
}

func (c fieldPrinter) printValue() string {

	var placeholders []string

	placeholder := c.placeholder
	if c.placeholder == "" {
		placeholder = runtime.QuestionMark
	}

	var isslice bool
	var r reflect.Value
	if c.v != nil {
		r = reflect.ValueOf(c.v)
		if r.Kind() == reflect.Slice {
			isslice = true
		}
	}

	if isslice {
		for i := 0; i < r.Len(); i++ {
			rv := r.Index(i)

			var pl string
			if placeholder == runtime.Dollar {
				pl = fmt.Sprintf("$%v", len(*c.collected))
			} else if placeholder == runtime.Named {
				pl = fmt.Sprintf("@p%v", len(*c.collected))
			} else if placeholder == runtime.QuestionMark {
				pl = "?"
			}
			placeholders = append(placeholders, pl)
			*(c.collected) = append(*(c.collected), rv.Interface())

		}

	} else {
		var pl string
		if placeholder == runtime.Dollar {
			pl = fmt.Sprintf("$%v", len(*c.collected))
		} else if placeholder == runtime.Named {
			pl = fmt.Sprintf("@p%v", len(*c.collected))
		} else if placeholder == runtime.QuestionMark {
			pl = "?"
		}
		placeholders = append(placeholders, pl)
		*(c.collected) = append(*(c.collected), c.v)
	}

	glue := c.glue
	if c.glue == "" {
		glue = ","
	}

	return strings.TrimSuffix(strings.Join(placeholders, glue), glue)
}

func (c fieldPrinter) printFieldValues() string {
	if c.v == nil {
		return ""
	}

	r := reflect.ValueOf(c.v)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	} else if r.Kind() == reflect.Slice {
		if r.Len() == 0 {
			r = reflect.Zero(r.Type().Elem())
		} else {
			r = r.Index(0)
		}
	}

	if r.Kind() != reflect.Struct {
		r = r.Elem()
	}

	var placeholders []string

	placeholder := c.placeholder
	if c.placeholder == "" {
		placeholder = runtime.QuestionMark
	}

	for _, field := range c.fields {
		if _, ok := r.Type().FieldByName(field.Prop); !ok {
			continue
		}
		f := r.FieldByName(field.Prop)
		var pl string
		if placeholder == runtime.Dollar {
			pl = fmt.Sprintf("$%v", len(*c.collected))
			*(c.collected) = append(*(c.collected), f.Interface())

		} else if placeholder == runtime.Named {
			n := fmt.Sprintf("p%v", len(*c.collected))
			pl = "@" + n
			*(c.collected) = append(*(c.collected), sql.NamedArg{Name: n, Value: f.Interface()})

		} else if placeholder == runtime.QuestionMark {
			pl = "?"
			*(c.collected) = append(*(c.collected), f.Interface())
		}
		placeholders = append(placeholders, pl)
	}

	glue := c.glue
	if c.glue == "" {
		glue = ","
	}

	return strings.TrimSuffix(strings.Join(placeholders, glue), glue)
}

func (c fieldPrinter) printUpdate() string {
	if c.v == nil {
		return ""
	}

	r := reflect.ValueOf(c.v)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	} else if r.Kind() == reflect.Slice {
		if r.Len() == 0 {
			r = reflect.Zero(r.Type().Elem())
		} else {
			r = r.Index(0)
		}
	}

	if r.Kind() != reflect.Struct {
		r = r.Elem()
	}

	var parts []string

	placeholder := c.placeholder
	if c.placeholder == "" {
		placeholder = runtime.QuestionMark
	}

	for _, field := range c.fields {
		if _, ok := r.Type().FieldByName(field.Prop); !ok {
			continue
		}
		f := r.FieldByName(field.Prop)
		*(c.collected) = append(*(c.collected), f.Interface())

		var pl string
		if placeholder == runtime.Dollar {
			pl = fmt.Sprintf("$%v", len(*c.collected))
		} else if placeholder == runtime.Named {
			pl = fmt.Sprintf("@p%v", len(*c.collected))
		} else if placeholder == runtime.QuestionMark {
			pl = "?"
		}
		colName := field.SQL
		if c.alias != "" {
			colName = c.alias + colName
		}
		parts = append(parts, fmt.Sprintf("%v = %v", colName, pl))
	}

	glue := c.glue
	if c.glue == "" {
		glue = ","
	}

	return strings.TrimSuffix(strings.Join(parts, glue), glue)
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

	var transformVals func(n parse.Node) bool
	transformVals = func(n parse.Node) bool {

		if x, ok := n.(*parse.PipeNode); ok {

			var shouldInclude bool
			visit(x, hasIdentifier(&shouldInclude, "vals"))

			if shouldInclude {
				collected := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGValues",
					},
				}
				placeholder := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGFlavor",
					},
				}
				t := append([]parse.Node{}, x.Cmds[0].Args[1:]...)
				x.Cmds[0].Args = append(x.Cmds[0].Args[:1], collected, placeholder)
				x.Cmds[0].Args = append(x.Cmds[0].Args, t...)
			}
		} else if x, ok := n.(*parse.CommandNode); ok {

			var shouldInclude bool
			visit(x, hasIdentifier(&shouldInclude, "vals"))

			if shouldInclude {
				collected := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGValues",
					},
				}
				placeholder := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGFlavor",
					},
				}
				t := append([]parse.Node{}, x.Args[1:]...)
				x.Args = append(x.Args[:1], collected, placeholder)
				x.Args = append(x.Args, t...)
			}
		}
		return true
	}
	visit(t.Tree.Root, transformVals)

	var transformVal func(n parse.Node) bool
	transformVal = func(n parse.Node) bool {

		if x, ok := n.(*parse.PipeNode); ok {

			var shouldInclude bool
			visit(x, hasIdentifier(&shouldInclude, "val"))

			if shouldInclude {
				collected := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGValues",
					},
				}
				placeholder := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGFlavor",
					},
				}
				t := append([]parse.Node{}, x.Cmds[0].Args[1:]...)
				x.Cmds[0].Args = append(x.Cmds[0].Args[:1], collected, placeholder)
				x.Cmds[0].Args = append(x.Cmds[0].Args, t...)
			}
		} else if x, ok := n.(*parse.CommandNode); ok {

			var shouldInclude bool
			visit(x, hasIdentifier(&shouldInclude, "val"))

			if shouldInclude {
				collected := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGValues",
					},
				}
				placeholder := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGFlavor",
					},
				}
				t := append([]parse.Node{}, x.Args[1:]...)
				x.Args = append(x.Args[:1], collected, placeholder)
				x.Args = append(x.Args, t...)
			}
		}
		return true
	}
	visit(t.Tree.Root, transformVal)

	var transformUpdate func(n parse.Node) bool
	transformUpdate = func(n parse.Node) bool {

		if x, ok := n.(*parse.PipeNode); ok {

			var shouldInclude bool
			visit(x, hasIdentifier(&shouldInclude, "update"))

			if shouldInclude {
				collected := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGValues",
					},
				}
				placeholder := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGFlavor",
					},
				}
				t := append([]parse.Node{}, x.Cmds[0].Args[1:]...)
				x.Cmds[0].Args = append(x.Cmds[0].Args[:1], collected, placeholder)
				x.Cmds[0].Args = append(x.Cmds[0].Args, t...)
			}
		} else if x, ok := n.(*parse.CommandNode); ok {

			var shouldInclude bool
			visit(x, hasIdentifier(&shouldInclude, "update"))

			if shouldInclude {
				collected := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGValues",
					},
				}
				placeholder := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGFlavor",
					},
				}
				t := append([]parse.Node{}, x.Args[1:]...)
				x.Args = append(x.Args[:1], collected, placeholder)
				x.Args = append(x.Args, t...)
			}
		}
		return true
	}
	visit(t.Tree.Root, transformUpdate)

	var transformFields func(n parse.Node) bool
	transformFields = func(n parse.Node) bool {

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

	var transformValuePrintings func(n parse.Node) bool
	transformValuePrintings = func(n parse.Node) bool {
		if x, ok := n.(*parse.ActionNode); ok {

			var hasIdent bool
			visit(x, hasIdentifier(&hasIdent, "comma", "cols", "raw", "print", "printf", "fields", "cols", "vals", "val", "update"))

			if !hasIdent {
				val := &parse.IdentifierNode{
					Ident: "val",
				}
				collected := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGValues",
					},
				}
				placeholder := &parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident: []string{
						"$", "SQLGFlavor",
					},
				}

				cmd1 := &parse.CommandNode{}
				cmd1.Args = append(cmd1.Args, val, collected, placeholder)

				x.Pipe.Cmds = append(x.Pipe.Cmds, cmd1)
			}

		}
		return true
	}
	visit(t.Tree.Root, transformValuePrintings)

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
