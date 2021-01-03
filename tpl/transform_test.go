package tpl_test

import (
	"testing"

	"github.com/clementauger/sqlg/tpl"
	"github.com/clementauger/sqlg/tpl/pg"
)

func TestTransform(t *testing.T) {
	type input struct {
		src string
		out string
		err error
	}
	table := []input{
		input{
			src: `text
	{{range $i, $a := .a}}
	 ( {{$a.Bio}} ) {{comma $i (len $a) }}
	{{end}}
	{{.y | pqArray}}`,
			out: `text
	{{range $i, $a := .a}}
	 ( {{$a.Bio | collect $.SQLGValues $.SQLGFlavor | placeholder $.SQLGValues $.SQLGFlavor}} ) {{comma $i (len $a)}}
	{{end}}
	{{.y | pqArray | collect $.SQLGValues $.SQLGFlavor | placeholder $.SQLGValues $.SQLGFlavor}}`,
		},
		input{
			src: `UPDATE authors SET
{{$fields := fields .a "id"}}
{{range $i, $field := $fields}}
	{{$field.SQL | print}} = {{$field.Value}} {{comma $i (len $fields) }}
{{end}}
WHERE id = {{.a.id}}`,
			out: `UPDATE authors SET
{{$fields := fields $.SQLGConverter .a "id"}}
{{range $i, $field := $fields}}
	{{$field.SQL | print}} = {{$field.Value | collect $.SQLGValues $.SQLGFlavor | placeholder $.SQLGValues $.SQLGFlavor}} {{comma $i (len $fields)}}
{{end}}
WHERE id = {{.a.id | collect $.SQLGValues $.SQLGFlavor | placeholder $.SQLGValues $.SQLGFlavor}}`,
		},
	}

	for _, test := range table {
		got, err := tpl.Transform(test.src, pg.FuncMap())
		if err != test.err {
			t.Fatalf("got unexpected error %v wanted %v", err, test.err)
		}
		if got != test.out {
			t.Fatalf("got\n%v\n\nwanted\n%v", got, test.out)
		}
	}
}
