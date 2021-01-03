package tpl_test

import (
	"github.com/clementauger/sqlg/tpl"
	"github.com/clementauger/sqlg/tpl/pg"
	"testing"
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
	{{.y | pqArray}}
`,
			out: `text
	{{range $i, $a := .a}}
	 ( {{$a.Bio | collect $.SQLGValues | placeholder $.SQLGFlavor}} ) {{comma $i (len $a)}}
	{{end}}
	{{.y | pqArray | collect $.SQLGValues | placeholder $.SQLGFlavor}}
`,
		},
	}

	for _, test := range table {
		got, err := tpl.Transform(test.src, pg.FuncMap())
		if err != test.err {
			t.Fatalf("got unexpected error %v wanted %v", err, test.err)
		}
		if got != test.out {
			t.Fatalf("got unexpected output %q wanted %q", got, test.out)
		}
	}
}
