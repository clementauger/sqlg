package parse

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"go/format"
	"go/types"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/clementauger/sqlg/runtime"
	"github.com/clementauger/sqlg/tpl"

	mssql "github.com/clementauger/sqlg/tpl/mssql"
	pg "github.com/clementauger/sqlg/tpl/pg"
)

func varName(name string, suffixes ...string) string {
	if len(suffixes) < 1 {
		return name
	}
	h := sha256.New()
	for _, s := range suffixes {
		fmt.Fprint(h, s)
	}
	s := fmt.Sprintf("%x", h.Sum(nil))
	return name + "" + s[:6]
}

func generateMethod(meth userMethod, engine, queryTemplates, rawQueries string) (string, error) {
	var out string

	// generate the method signature
	tName := strings.Title(meth.Receiver.StructName)
	exportedRcvType := strings.Replace(meth.Receiver.GoType, meth.Receiver.StructName, tName, -1)
	comment := strings.TrimSpace(meth.Comment)
	if comment != "" {
		comment = "// " + strings.Replace(comment, "\n", "\n// ", -1)
		out += comment + "\n"
	}
	out += fmt.Sprintf("func (%v %v) %v", meth.Receiver.Name, exportedRcvType, meth.Name)
	out += "(ctx context.Context, "
	if meth.Mode == modeQuery {
		out += "db sqlg.Querier,"
	} else {
		out += "db sqlg.Execer,"
	}
	for _, p := range meth.InParams {
		out += p.Name + " " + p.GoType + ","
	}
	out = strings.TrimSuffix(out, ",")
	out += ") ("
	for _, p := range meth.OutParams {
		if p.IsFunc {
			out += p.Name + " " + p.Func.OutParams[0].StructName + "Iterator,"
		} else {
			out += p.Name + " " + p.GoType + ","
		}
	}
	out = strings.TrimSuffix(out, ",")
	out += ") {\n"

	if meth.Query != "" {

		// generate the query using templates
		sqlQuery := varName("sqlQuery", meth.FileName)
		out += fmt.Sprintf("var %v string\n", sqlQuery)
		SQLGValues := varName("SQLGValues", meth.FileName)
		SQLGFlavor := varName("SQLGFlavor", meth.FileName)
		if len(meth.InParams) > 0 {
			out += fmt.Sprintf("%v := &[]interface{}{}\n", SQLGValues)
			// if byDir {
			out += fmt.Sprintf("%v := %q\n", SQLGFlavor, runtime.EngineToPlaceholder(engine))
			// } else {
			// 	out += fmt.Sprintf("%v := runtime.Current\n", SQLGFlavor)
			// }
			out += "{\n"
			query := varName("query", meth.FileName)
			// create a new buffer to store the template result
			out += fmt.Sprintf("var %v bytes.Buffer\n", query)
			templateInput := varName("templateInput", meth.FileName)
			// create the map of template parameters
			out += fmt.Sprintf("%v := map[string]interface{}{\n", templateInput)
			if meth.CaseConverter != nil {
				out += fmt.Sprintf("%q:%v.%v,\n", "SQLGConverter", meth.Receiver.Name, meth.CaseConverter.Name())
			} else {
				out += fmt.Sprintf("%q:sqlg.NilCaseConverter{},\n", "SQLGConverter")
			}
			out += fmt.Sprintf("%q:%v,\n", "SQLGValues", SQLGValues)
			out += fmt.Sprintf("%q:%v,\n", "SQLGFlavor", SQLGFlavor)
			for _, h := range meth.InParams {
				if h.Name != "_" {
					out += fmt.Sprintf("%q:%v,\n", h.Name, h.Name)
				}
			}
			for _, h := range meth.OutParams {
				if h.Name != "_" {
					out += fmt.Sprintf("%q:%v,\n", h.Name, h.Name)
				}
			}
			for _, y := range meth.TemplateParams {
				out += fmt.Sprintf("%v:%v,\n", y.Name, y.Expr)
			}
			out += "}\n"
			// execute the template query
			out += fmt.Sprintf("err = %v[%q].Execute(&%v, %v)\n",
				queryTemplates, meth.Receiver.StructName+"__"+meth.Name, query, templateInput)
			out += "if err != nil { return }\n"
			out += fmt.Sprintf("%v = %v.String()\n\n", sqlQuery, query)

			if meth.Logger != nil {
				out += fmt.Sprintf("%v.%v.Log(%q, %q, %v, (*%v)... )\n",
					meth.Receiver.Name, meth.Logger.Name(), meth.Receiver.PackagePath+"/"+meth.Receiver.StructName,
					meth.Name,
					sqlQuery, SQLGValues)
			}
			if meth.Tracer != nil {
				out += fmt.Sprintf("%v.%v.Begin(%q, %q, %v, (*%v)... )\n",
					meth.Receiver.Name, meth.Tracer.Name(), meth.Receiver.PackagePath+"/"+meth.Receiver.StructName,
					meth.Name,
					sqlQuery, SQLGValues)
				out += fmt.Sprintf("defer func(){\n")
				out += fmt.Sprintf("%v.%v.End(%q, %q, %v )\n",
					meth.Receiver.Name, meth.Tracer.Name(), meth.Receiver.PackagePath+"/"+meth.Receiver.StructName,
					meth.Name, "err")
				out += fmt.Sprintf("}()\n")
			}

			out += "}\n"
		} else {
			// get the query from non template raw queries
			out += fmt.Sprintf("%v = %v[%q]\n\n",
				sqlQuery, rawQueries, meth.Receiver.StructName+"__"+meth.Name)

			if meth.Logger != nil {
				out += fmt.Sprintf("%v.%v.Log(%q, %q, %v)\n",
					meth.Receiver.Name, meth.Logger.Name(), meth.Receiver.PackagePath+"/"+meth.Receiver.StructName,
					meth.Name,
					sqlQuery)
			}
			if meth.Tracer != nil {
				out += fmt.Sprintf("%v.%v.Begin(%q, %q, %v )\n",
					meth.Receiver.Name, meth.Tracer.Name(), meth.Receiver.PackagePath+"/"+meth.Receiver.StructName,
					meth.Name,
					sqlQuery)
				out += fmt.Sprintf("defer func(){\n")
				out += fmt.Sprintf("%v.%v.End(%q, %q, %v )\n",
					meth.Receiver.Name, meth.Tracer.Name(), meth.Receiver.PackagePath+"/"+meth.Receiver.StructName,
					meth.Name, "err")
				out += fmt.Sprintf("}()\n")
			}
		}
		out += "\n"
		// execute the query
		if meth.Mode == modeQuery {
			// create the rows variable
			rows := varName("rows", meth.FileName)
			out += fmt.Sprintf("var %v *sql.Rows\n", rows)
			if len(meth.InParams) > 0 {
				out += fmt.Sprintf("%v, err = db.QueryContext(ctx, %v, (*%v)...)\n", rows, sqlQuery, SQLGValues)
			} else {
				out += fmt.Sprintf("%v, err = db.QueryContext(ctx, %v)\n", rows, sqlQuery)
			}
			out += "if err != nil { return }\n"

			if meth.OutputsToIterator() {
				// out += fmt.Sprintf("var it %vIterator\n", meth.OutParams[0].Func.OutParams[0].StructName)
				out += fmt.Sprintf("it.rows = %v\n", rows)

			} else {

				out += fmt.Sprintf("for %v.Next() {\n", rows)
				// handle slice output
				if meth.OutputsToStructSlice() {
					// setup the new slice item
					item := varName("item", meth.FileName)
					if meth.OutParams[0].IsPtr {
						out += fmt.Sprintf(`%v := new(`, item)
					} else {
						out += fmt.Sprintf(`var %v `, item)
					}
					if meth.OutParams[0].PackagePath != meth.Receiver.PackagePath {
						out += fmt.Sprintf(`%v.`, filepath.Base(meth.OutParams[0].PackagePath))
					}
					out += fmt.Sprintf("%v\n", meth.OutParams[0].StructName)
					if meth.OutParams[0].IsPtr {
						out += ")"
					}
					// scan the row
					out += fmt.Sprintf("err = %v.Scan(", rows)
					for _, n := range meth.OutParams[0].StructProperties {
						out += fmt.Sprintf("&%v.%v,", item, n)
					}
					out = strings.TrimSuffix(out, ",")
					out += ")\n"
					out += "if err != nil { return }\n"
					// save the result
					out += fmt.Sprintf("%v = append(%v, %v)\n",
						meth.OutParams[0].Name, meth.OutParams[0].Name, item)

				} else if meth.OutputsToStruct() {
					// scan the row
					out += fmt.Sprintf("err = %v.Scan(", rows)
					for _, n := range meth.OutParams[0].StructProperties {
						out += fmt.Sprintf("&%v.%v,", meth.OutParams[0].Name, n)
					}
					out = strings.TrimSuffix(out, ",")
					out += ")\n"
					out += "if err != nil { return }\n"

				} else if meth.OutputsToBasic() {
					// scan the row
					out += fmt.Sprintf("err = %v.Scan(", rows)
					for _, n := range meth.OutParams[:len(meth.OutParams)-1] {
						out += fmt.Sprintf("&%v,", n.Name)
					}
					out = strings.TrimSuffix(out, ",")
					out += ")\n"
					out += "if err != nil { return }\n"
				}
				out += "}\n"
				// Close the rows
				out += fmt.Sprintf("if err = %v.Close(); err != nil { return }\n", rows)
				// Check for errors
				out += fmt.Sprintf("err = %v.Err()\n", rows)
			}
		} else if meth.Mode == modeExec {
			//-
			res := varName("res", meth.FileName)
			if meth.AffectedRows != "" || meth.InsertedID != "" {
				out += fmt.Sprintf("var %v sql.Result\n", res)
				if len(meth.InParams) > 0 {
					out += fmt.Sprintf("%v, err = db.ExecContext(ctx, %v, (*%v)...)\n",
						res, sqlQuery, SQLGValues)
				} else {
					out += fmt.Sprintf("%v, err = db.ExecContext(ctx, %v)\n",
						res, sqlQuery)
				}
				out += "if err != nil { return }\n"
				if meth.AffectedRows != "" {
					j := meth.Param(meth.AffectedRows)
					if j != nil && j.IsPtr {
						out += fmt.Sprintf("*%v,err=%v.RowsAffected()\n", meth.AffectedRows, res)
					} else {
						out += fmt.Sprintf("%v,err=%v.RowsAffected()\n", meth.AffectedRows, res)
					}
					out += "if err != nil { return }\n"
				}
				if meth.InsertedID != "" {
					j := meth.Param(meth.InsertedID)
					if j != nil && j.IsPtr {
						out += fmt.Sprintf("*%v,err=%v.LastInsertId()\n", meth.InsertedID, res)
					} else {
						out += fmt.Sprintf("%v,err=%v.LastInsertId()\n", meth.InsertedID, res)
					}
					out += "if err != nil { return }\n"
				}
			} else {
				if len(meth.InParams) > 0 {
					out += fmt.Sprintf("_,err=db.ExecContext(ctx, %v, (*%v)...)\n",
						sqlQuery, SQLGValues)
				} else {
					out += fmt.Sprintf("_,err=db.ExecContext(ctx, %v)\n",
						sqlQuery)
				}
				out += "if err != nil { return }\n"
			}
		}

	} else {
		if meth.FinalErr != "" && meth.FinalErr != "nil" {
			errOut := meth.OutParams[len(meth.OutParams)-1]
			out += fmt.Sprintf("%v = %v\n", errOut.Name, meth.FinalErr)
		}
	}
	out += "return\n"
	out += "}\n\n"
	return out, nil
}

func generateQueryTemplates(f FileObjects, engine, queryTemplates, rawQueries string) (string, error) {
	var out string
	funcMap := tpl.FuncMap()
	if engine == "mssql" {
		funcMap = mssql.FuncMap()
	} else if engine == "pg" {
		funcMap = pg.FuncMap()
	}
	out += fmt.Sprintf("var %v = map[string]*template.Template{\n", queryTemplates)
	keys, queries := f.Queries(false)
	for _, key := range keys {
		query := queries[key]
		tr, err := tpl.Transform(query, funcMap)
		if err != nil {
			return "", fmt.Errorf(
				"failed to transform %v.%v.%v query %q: %v",
				f.filePath, strings.Split(key, "__")[0], strings.Split(key, "__")[1], query, err)
		}
		out += fmt.Sprintf(
			`%q:template.Must(template.New("").Funcs(tpl.FuncMap()).Parse(
							%v,
							)),
							`,
			key, tr)
	}
	out += "}\n\n"
	// write the map of raw queries
	out += fmt.Sprintf("var %v = map[string]string{\n", rawQueries)
	keys, queries = f.Queries(true)
	for _, key := range keys {
		query := queries[key]
		out += fmt.Sprintf("%q:%v,\n", key, query)
	}
	out += "}\n\n"

	return out, nil
}

func generateTypeInterface(fileObjects map[string]FileObjects, pkgPath, typName string) string {
	var out string
	// write the interface
	tName := strings.Title(typName)
	ifaceName := tName + "Iface"
	out += fmt.Sprintf("// %v is an interface of %v\n", ifaceName, tName)
	out += fmt.Sprintf("type %v interface{\n", ifaceName)
	files := []string{}
	for fpath := range fileObjects {
		files = append(files, fpath)
	}
	sort.Strings(files)
	for _, fpath := range files {
		f := fileObjects[fpath]
		if f.PackagePath != pkgPath {
			continue
		}
		for _, meth := range f.Methods {
			if meth.Receiver.StructName != typName {
				continue
			}
			out += fmt.Sprintf(" %v", meth.Name)
			out += "(ctx context.Context, "
			if meth.Mode == modeQuery {
				out += "db sqlg.Querier,"
			} else {
				out += "db sqlg.Execer,"
			}
			for _, p := range meth.InParams {
				out += p.Name + " " + p.GoType + ","
			}
			out = strings.TrimSuffix(out, ",")
			out += ") ("
			for _, p := range meth.OutParams {
				if p.IsFunc {
					out += p.Name + " " + p.Func.OutParams[0].StructName + "Iterator,"
				} else {
					out += p.Name + " " + p.GoType + ","
				}
			}
			out = strings.TrimSuffix(out, ",")
			out += ")\n"
		}
	}
	out += "}\n"
	return out
}

func generateIterator(meth userMethod) string {
	typName := meth.OutParams[0].Func.OutParams[0].StructName
	var out string
	// write the interface
	itName := typName + "Iterator"
	out += fmt.Sprintf("// %v is an iterator of %v\n", itName, typName)
	out += fmt.Sprintf(`type %v struct{
		rows *sql.Rows
		err error
		value %v
	}
	`, itName, meth.OutParams[0].Func.OutParams[0].GoType)

	out += fmt.Sprintf(`
	func(x %v) Err() error {
		if x.err!=nil{
			return x.err
		}
		return x.rows.Err()
	}
	`, itName)

	out += fmt.Sprintf(`
	func(x %v) Close() error {
		return x.rows.Close()
	}
	`, itName)

	out += fmt.Sprintf(`
	func(x %v) All() (ret []%v) {
		for x.Next() {
			ret = append(ret,x.Value())
		}
		return ret
	}
	`, itName, meth.OutParams[0].Func.OutParams[0].GoType)

	out += fmt.Sprintf(`func(it *%v) Value() (%v) {
		`, itName, meth.OutParams[0].Func.OutParams[0].GoType)
	out += "return it.value\n"
	out += "}\n"

	out += fmt.Sprintf(`func(it *%v) Next() ( bool) {
		`, itName)
	out += fmt.Sprintf(`if !it.rows.Next(){
			it.rows.Close()
			return false
		}
		`)

	if meth.OutParams[0].Func.OutParams[0].IsPtr {
		// always alloc to prevent beginner problems.
		out += fmt.Sprintf(`it.value = new(`)
		if meth.OutParams[0].Func.OutParams[0].PackagePath != meth.Receiver.PackagePath {
			out += fmt.Sprintf(`%v.`, filepath.Base(meth.OutParams[0].Func.OutParams[0].PackagePath))
		}
		out += fmt.Sprintf("%v\n", meth.OutParams[0].Func.OutParams[0].StructName)
		out += ")"
	}

	out += fmt.Sprintf(`it.err =it.rows.Scan(`)
	for _, n := range meth.OutParams[0].Func.OutParams[0].StructProperties {
		out += fmt.Sprintf("&(it.value.%v),", n)
	}
	out = strings.TrimSuffix(out, ",")
	out += fmt.Sprintf(")\n")
	out += fmt.Sprintf("if it.err == nil { return true  }\n")
	out += fmt.Sprintf("return false\n")
	out += "}\n"
	return out
}

func Generate(fileObjects map[string]FileObjects, engine string) (map[string]string, error) {

	out := map[string]string{}

	var finalErr error

	for _, f := range fileObjects {

		var fContent string
		// write the map of query templates
		queryTemplates := varName("queryTemplates", f.filePath)
		rawQueries := varName("rawQueries", f.filePath)

		q, err := generateQueryTemplates(f, engine, queryTemplates, rawQueries)
		if err != nil {
			return nil, err
		}
		fContent += q

		// write the struct statement
		for _, ut := range f.Types {
			tName := strings.Title(ut.Name)
			comment := strings.TrimSpace(ut.Comment)
			if strings.HasPrefix(comment, ut.Name) {
				comment = strings.TrimPrefix(comment, ut.Name)
				comment = tName + comment
			}
			if comment != "" {
				comment = "// " + strings.Replace(comment, "\n", "\n// ", -1)
				fContent += comment + "\n"
			}
			fContent += fmt.Sprintf("type %v struct{\n", tName)
			if ut.Tracer != nil {
				typName := filepath.Base(ut.Tracer.Type().String())
				pkgPath := ut.Tracer.Type().(*types.Named).Obj().Pkg().Path()
				if pkgPath == "github.com/clementauger/sqlg/runtime" {
					typName = strings.TrimPrefix(typName, "runtime.")
					typName = "sqlg." + typName
				} else if pkgPath == ut.PackagePath {
					typName = strings.TrimPrefix(typName, filepath.Base(pkgPath)+".")
				}
				fContent += fmt.Sprintf("%v %v\n", ut.Tracer.Name(), typName)
			}
			if ut.Logger != nil {
				typName := filepath.Base(ut.Logger.Type().String())
				pkgPath := ut.Logger.Type().(*types.Named).Obj().Pkg().Path()
				if pkgPath == "github.com/clementauger/sqlg/runtime" {
					typName = strings.TrimPrefix(typName, "runtime.")
					typName = "sqlg." + typName
				} else if pkgPath == ut.PackagePath {
					typName = strings.TrimPrefix(typName, filepath.Base(pkgPath)+".")
				}
				fContent += fmt.Sprintf("%v %v\n", ut.Logger.Name(), typName)
			}
			if ut.CaseConverter != nil {
				typName := filepath.Base(ut.CaseConverter.Type().String())
				pkgPath := ut.CaseConverter.Type().(*types.Named).Obj().Pkg().Path()
				if pkgPath == "github.com/clementauger/sqlg/runtime" {
					typName = strings.TrimPrefix(typName, "runtime.")
					typName = "sqlg." + typName
				} else if pkgPath == ut.PackagePath {
					typName = strings.TrimPrefix(typName, filepath.Base(pkgPath)+".")
				}
				fContent += fmt.Sprintf("%v %v\n", ut.CaseConverter.Name(), typName)
			}
			fContent += fmt.Sprintln("}")

			// write the interface
			fContent += generateTypeInterface(fileObjects, ut.PackagePath, ut.Name)
			fContent += fmt.Sprintln()
		}

		// generate iterators
		iterators := map[string]userMethod{}
		for _, meth := range f.Methods {
			if len(meth.OutParams) > 0 && meth.OutParams[0].IsFunc {
				iterators[meth.OutParams[0].Func.OutParams[0].StructName] = meth
			}
		}
		for _, meth := range iterators {
			fContent += generateIterator(meth) + "\n"
		}

		// generate methods
		for _, meth := range f.Methods {
			k, err := generateMethod(meth, engine, queryTemplates, rawQueries)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to generate method %v.%v.%v query %q: %v",
					f.filePath, meth.Receiver.StructName, meth.Name, meth.Query, err)
			}
			fContent += k + "\n"
		}

		var noComment string
		sc := bufio.NewScanner(strings.NewReader(fContent))
		for sc.Scan() {
			line := sc.Text()
			if strings.HasPrefix(line, "// ") {
				continue
			}
			noComment += line + "\n"
		}

		// generate import statements
		imports := f.Imports()
		imports["bytes"] = ""
		if strings.Contains(noComment, "sql.") {
			imports["database/sql"] = ""
		}
		if strings.Contains(noComment, "fmt.") {
			imports["fmt"] = ""
		}
		imports["text/template"] = ""
		if engine == runtime.MsSQL {
			imports["github.com/clementauger/sqlg/tpl/mssql"] = "tpl"
		} else if engine == runtime.PostgreSQL {
			imports["github.com/clementauger/sqlg/tpl/pg"] = "tpl"
		} else {
			imports["github.com/clementauger/sqlg/tpl"] = ""
		}
		imports["github.com/clementauger/sqlg/runtime"] = "sqlg"
		imports["context"] = ""

		// import tracer, logger, converter packages

		for _, ut := range f.Types {
			var pkgPath string
			if ut.Tracer != nil {
				pkgPath = ut.Tracer.Type().(*types.Named).Obj().Pkg().Path()
			} else if ut.Logger != nil {
				pkgPath = ut.Logger.Type().(*types.Named).Obj().Pkg().Path()
			} else if ut.CaseConverter != nil {
				pkgPath = ut.CaseConverter.Type().(*types.Named).Obj().Pkg().Path()
			}
			if _, ok := imports[pkgPath]; !ok && pkgPath != ut.PackagePath {
				imports[pkgPath] = ""
			}
		}

		importsText := "import(\n"
		for i, n := range imports {
			if n != "" {
				importsText += fmt.Sprintf("%v ", n)
			}
			importsText += fmt.Sprintf("%q\n", i)
		}
		importsText += ")\n"

		// generate the build constraints
		// tags := f.fileTags
		// tags = strings.Replace(tags, "sqlg,", "", -1)
		// tags = strings.Replace(tags, "sqlg ", "", -1)
		// tags = strings.Replace(tags, "sqlg", "", -1)
		// tags = strings.TrimSpace(tags)

		// write headers
		headers := "//+build !sqlg\n"
		// if !byDir && tags != "" {
		// 	headers += fmt.Sprintf("//+build %v\n", tags)
		// }
		headers += "\n"
		headers += "// Code generated by sqlg DO NOT EDIT\n\n"
		headers += "\n"
		headers += fmt.Sprintf("package %v\n", engine)
		headers += "\n"
		headers += importsText + "\n"

		fContent = headers + "\n" + fContent

		outfile := strings.TrimSuffix(f.filePath, ".go")
		outfile = fmt.Sprintf("%v_gen.go", outfile)
		// if byDir {
		outdir := filepath.Dir(outfile)
		outfile = filepath.Base(outfile)
		outfile = filepath.Join(outdir, engine, outfile)
		// }

		t, err := format.Source([]byte(fContent))
		if err != nil {
			log.Println(outfile, ":", err)
			log.Println(fContent)
			finalErr = fmt.Errorf("source formating failure: %v", err)
		} else {
			fContent = string(t)
		}
		out[outfile] = fContent
	}

	return out, finalErr
}
