package parse

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

func FromPackagePath(packagePath, engine string, args []string) (map[string]FileObjects, error) {
	if engine != "" {
		engine = "," + engine
	}
	var f bool
	for i, a := range args {
		if strings.HasPrefix(a, "-tags") {
			if !strings.Contains(a, "sqlg") {
				f = true
				args[i] = a + ",sqlg"
				break
			}
		}
	}
	if !f {
		args = append(args, "-tags=sqlg"+engine)
	}
	cfg := &packages.Config{
		// Mode:       packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.LoadTypes,
		Mode:       packages.NeedSyntax | packages.NeedImports | packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes,
		BuildFlags: args,
	}
	pkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load package: %v", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("failed to load packages")
	}

	out := map[string]FileObjects{}
	for _, pkg := range pkgs {
		objs, err := parsePackage(pkg.Types, pkg.Fset, pkg.Syntax)
		if err != nil {
			return nil, err
		}
		for f, o := range objs {
			x := out[f]
			x.PackagePath = o.PackagePath
			x.PackageName = o.PackageName
			x.filePath = o.filePath
			x.fileTags = o.fileTags
			x.Types = append(x.Types, o.Types...)
			x.Methods = append(x.Methods, o.Methods...)
			out[f] = x
		}
	}
	return out, nil
}

func FromString(pkgName, filePath, src string) (map[string]FileObjects, error) {
	fset := token.NewFileSet()
	var files []*ast.File
	for _, file := range []struct{ name, input string }{{filePath, src}} {
		f, err := parser.ParseFile(fset, file.name, file.input, 0)
		if err != nil {
			return nil, fmt.Errorf("parse failure: %v", err)
		}
		files = append(files, f)
	}

	// Type-check a package consisting of these files.
	// Type information for the imported "fmt" package
	// comes from $GOROOT/pkg/$GOOS_$GOOARCH/fmt.a.
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check(pkgName, fset, files, nil)
	if err != nil {
		return nil, fmt.Errorf("type check failed: %v", err)
	}

	return parsePackage(pkg, fset, files)
}

type userType struct {
	Name          string
	PackagePath   string
	PackageName   string
	FileName      string
	Methods       []userMethod
	Comment       string
	Tracer        *types.Var
	Logger        *types.Var
	CaseConverter *types.Var
}

type sqlMode string

var (
	modeQuery sqlMode = "query"
	modeExec  sqlMode = "exec"
)

type userMethod struct {
	FileName     string
	Name         string
	Prepared     bool
	InsertedID   string
	AffectedRows string
	Mode         sqlMode
	Query        string
	InParams     []userParam
	OutParams    []userParam
	Receiver     userParam
	Comment      string

	FinalErr string

	Tracer        *types.Var
	Logger        *types.Var
	CaseConverter *types.Var
}

type userParam struct {
	IsPtr            bool
	IsSlice          bool
	Name             string
	GoType           string // is the full go type such as []int or *Author, or Author
	PackagePath      string
	BasicKind        types.BasicKind
	BasicInfo        types.BasicInfo
	StructProperties []string // Property of the underlying struct, being a slice, a pointer or a value
	StructName       string   // Name of the underlying struct type, being a slice, a pointer or a value

	IsFunc bool
	Func   userFunc
}

type userFunc struct {
	InParams  []userParam
	OutParams []userParam
}

func (p userFunc) String() string {
	out := "func( "
	for _, pp := range p.InParams {
		out += pp.GoType + ","
	}
	out = strings.TrimSuffix(out, ",")
	out += ") ("
	for _, pp := range p.OutParams {
		out += pp.GoType + ","
	}
	out = strings.TrimSuffix(out, ",")
	out += ")"
	return out
}

type FileObjects struct {
	filePath    string
	fileTags    string
	PackageName string
	PackagePath string
	Types       []userType
	Methods     []userMethod
}

func (f FileObjects) Imports() map[string]string {
	out := map[string]string{}
	for _, m := range f.Methods {
		for _, p := range m.InParams {
			if p.PackagePath != "" && p.PackagePath != f.PackagePath {
				out[p.PackagePath] = ""
			}
		}
		for _, p := range m.OutParams {
			if p.PackagePath != "" && p.PackagePath != f.PackagePath {
				out[p.PackagePath] = ""
			}
		}
	}
	return out
}

func (f FileObjects) Queries(raw bool) ([]string, map[string]string) {
	out := map[string]string{}
	var keys []string
	for _, m := range f.Methods {
		if m.Query == "" {
			continue
		}
		if raw && len(m.InParams) < 1 {
			name := m.Receiver.StructName + "__" + m.Name
			out[name] = m.Query
			keys = append(keys, name)
		} else if !raw && len(m.InParams) > 0 {
			name := m.Receiver.StructName + "__" + m.Name
			out[name] = m.Query
			keys = append(keys, name)
		}
	}
	sort.Strings(keys)
	return keys, out
}

func validFunc(s string) bool {
	fns := []string{"InsertedID", "AffectedRows", "Prepared", "Exec", "Query"}
	for _, n := range fns {
		if s == n {
			return true
		}
	}
	return false
}

func typeFieldLookup(st *types.Struct, typSearch, namSearch string) *types.Var {
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if f.Embedded() && strings.HasSuffix(f.Type().String(), typSearch) {
			return f
		} else if !f.Embedded() && f.Name() == namSearch && strings.HasSuffix(f.Type().String(), typSearch) {
			return f
		}
	}
	return nil
}

func typeFieldLookupByName(st *types.Struct, namSearches ...string) *types.Var {
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if f.Embedded() {
			continue
		}
		for _, namSearch := range namSearches {
			if f.Name() == namSearch {
				return f
			}
		}
	}
	return nil
}

func lookupTypComments(fset *token.FileSet, syntax []*ast.File, start token.Pos, end token.Pos, typName string) string {
	file, ok := fileForPos(fset, syntax, start, end)
	if ok && file != nil {
		for _, decl := range file.Decls {
			if x, ok := decl.(*ast.GenDecl); ok {
				if len(x.Specs) > 0 {
					if y, ok := x.Specs[0].(*ast.TypeSpec); ok {
						if y.Name.Name == typName {
							return x.Doc.Text()
						}
					}
				}
			}
		}
	}
	return ""
}

func parsePackage(pkg *types.Package, fset *token.FileSet, syntax []*ast.File) (map[string]FileObjects, error) {
	ret := map[string]FileObjects{}
	// qual := types.RelativeTo(pkg.Types)
	scope := pkg.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)

		if obj != nil && obj.Type() != nil && obj.Type().Underlying() != nil {
			// obj is types.Named,
			// obj.Type() is obj.TypeName
			// obj.Type().Underlyng exposes the struct.Types which
			// gives access to the desired types.Var.Embedded
			st, ok := obj.Type().Underlying().(*types.Struct)
			if !ok {
				continue
			}

			hasSQLG := typeFieldLookup(st, "runtime.SQLg", "") != nil

			if !hasSQLG {
				continue
			}

			var u userType
			u.PackagePath = pkg.Path()
			u.PackageName = pkg.Name()
			u.FileName = fset.File(obj.Pos()).Name()
			u.Name = obj.Name()

			// todo: make use of type infrence to not rely on type string and field name comparisons
			tracer := typeFieldLookupByName(st, "tracer", "Tracer")
			logger := typeFieldLookupByName(st, "logger", "Logger")
			caseConverter := typeFieldLookupByName(st, "converter", "Converter")

			u.Tracer = tracer
			u.Logger = logger
			u.CaseConverter = caseConverter

			// lookup for type comment.
			u.Comment = lookupTypComments(fset, syntax, obj.Pos(), obj.Pos(), u.Name)

			// fmt.Printf("\t%s\n", types.ObjectString(obj, qual))

			// lookup for type methods.
			if _, ok := obj.(*types.TypeName); ok {
				for _, meth := range typeutil.IntuitiveMethodSet(obj.Type(), nil) {
					if !meth.Obj().Exported() {
						continue // skip unexported names
					}
					if len(meth.Index()) > 1 {
						continue // Dont process methods that are from embedded types
					}

					// grab the receiver variable name of the method
					sig := meth.Obj().Type().(*types.Signature)
					receiverName := sig.Recv().Name()

					// Lokup for the ast tree of the function body
					var fnDecl *ast.FuncDecl
					path, _ := pathEnclosingInterval(fset, syntax, meth.Obj().Pos(), meth.Obj().Pos())
					for _, n := range path {
						if x, ok := n.(*ast.FuncDecl); ok {
							fnDecl = x
							break
						}
					}
					if fnDecl == nil {
						return ret, fmt.Errorf("failed to find function body of %q", meth.Obj().Name())
					}

					opts := browseBodyFunc(fnDecl, receiverName)

					methFile := fset.File(meth.Obj().Pos()).Name()
					var um userMethod
					um.Name = meth.Obj().Name()
					um.FileName = methFile
					um.Comment = fnDecl.Doc.Text()

					um.Tracer = u.Tracer
					um.Logger = u.Logger
					um.CaseConverter = u.CaseConverter

					um.Receiver.StructName = u.Name
					um.Receiver.Name = receiverName
					um.Receiver.PackagePath = pkg.Path()
					if _, ok := sig.Recv().Type().(*types.Pointer); ok {
						um.Receiver.IsPtr = true
						um.Receiver.GoType = "*"
					}
					um.Receiver.GoType += obj.Name()

					if q, ok := opts["Query"]; ok && len(q) > 0 {
						um.Mode = modeQuery
						um.Query = q[0]
					} else if q, ok := opts["Exec"]; ok && len(q) > 0 {
						um.Mode = modeExec
						um.Query = q[0]
					}
					if q, ok := opts["Prepared"]; ok && len(q) > 0 {
						um.Prepared = q[0] == "true"
					}
					if q, ok := opts["InsertedID"]; ok && len(q) > 0 {
						um.InsertedID = q[0]
					}
					if q, ok := opts["AffectedRows"]; ok && len(q) > 0 {
						um.AffectedRows = q[0]
					}

					if _, ok := opts["Query"]; !ok {
						errExpr := lookupFinalErrExpr(fnDecl)
						um.FinalErr = errExpr
					}

					var err error
					um.InParams, err = tuplesToParams(obj.Pkg().Path(), sig.Params(), true)
					if err != nil {
						return nil, err
					}
					um.OutParams, err = tuplesToParams(obj.Pkg().Path(), sig.Results(), true)
					if err != nil {
						return nil, err
					}

					err = checkMethod(u, um)
					if err != nil {
						return nil, err
					}

					buildTags, ok := lookupForFileBuildags(fset, syntax, meth.Obj())
					if !ok {
						return nil, fmt.Errorf("failed to find ast file for method %v.%v", obj.Pkg().Path(), meth.Obj().Name())
					}
					x := ret[methFile]
					x.Methods = append(x.Methods, um)
					x.fileTags = buildTags
					x.filePath = methFile
					x.PackageName = u.PackageName
					x.PackagePath = u.PackagePath
					ret[methFile] = x

				}
			}

			buildTags, ok := lookupForFileBuildags(fset, syntax, obj)
			if !ok {
				return nil, fmt.Errorf("failed to find ast file for type %v.%v", obj.Pkg().Path(), obj.Name())
			}
			x := ret[u.FileName]
			x.Types = append(x.Types, u)
			x.fileTags = buildTags
			x.filePath = u.FileName
			x.PackageName = u.PackageName
			x.PackagePath = u.PackagePath
			ret[u.FileName] = x
		}

	}
	return ret, nil
}

func lookupForFileBuildags(fset *token.FileSet, files []*ast.File, obj types.Object) (string, bool) {
	fileSyntax, ok := fileForPos(fset, files, obj.Pos(), obj.Pos())
	if !ok {
		return "", false
	}
	if len(fileSyntax.Comments) < 1 {
		return "", true
	}
	text := fileSyntax.Comments[0].Text()
	if strings.HasPrefix(text, "+build") {
		text = strings.TrimPrefix(text, "+build")
	}
	text = strings.TrimSpace(text)

	return text, true
}

// Browse body function top statements
// looking for calls to the sqlg runtime functions.
// record func calls and their arguments into a map
// of [function]=> {arguments...}
func browseBodyFunc(fnDecl *ast.FuncDecl, receiverName string) map[string][]string {
	opts := map[string][]string{}
	for _, n := range fnDecl.Body.List {
		if x, ok := n.(*ast.ExprStmt); ok {
			var currentFunc string
			// in k.Z() ast.Inspect will first present the Ident of k, then Z.
			// So we look for chain calls that starts with an ident.Value == receiverName,
			// If the first ident found is not the receiver, return false, and stop browsing this branch.
			// If first ident is the receiver name, identify function calls and their arguments
			// to record them into the opts map.
			// Unless validIdent is true, dont process the node and keep browsing the branch.
			var foundIdent bool
			var validIdent bool
			ast.Inspect(x, func(n ast.Node) bool {
				if x, ok := n.(*ast.Ident); ok {
					if !foundIdent {
						foundIdent = true
						if x.Name == receiverName {
							validIdent = true
							return true
						}
					}
					if !validIdent {
						return false
					}
				}
				if validIdent {
					switch x := n.(type) {
					case *ast.Ident:
						if validFunc(x.Name) {
							currentFunc = x.Name
						} else {
							opts[currentFunc] = append(opts[currentFunc], x.Name)
						}
					case *ast.BasicLit:
						opts[currentFunc] = append(opts[currentFunc], x.Value)
					}
				}
				return true
			})
		}
	}
	return opts
}
func lookupFinalErrExpr(fnDecl *ast.FuncDecl) string {
	var out string
	for _, n := range fnDecl.Body.List {
		if x, ok := n.(*ast.ReturnStmt); ok {
			if len(x.Results) > 0 {
				fset := token.NewFileSet()
				var b bytes.Buffer
				format.Node(&b, fset, x.Results[len(x.Results)-1])
				out = b.String()
			} else {
				out = ""
			}
		}
	}
	return out
}

func checkMethod(u userType, um userMethod) error {
	// Do some checks
	// 1. The method must return a trailing error of type error
	if um.OutParams[len(um.OutParams)-1].GoType != "error" {
		return fmt.Errorf(
			"invalid method signature %v.%v.%v, it must return a trailing error of type errors.error",
			u.PackagePath, u.Name, um.Name)
	}
	// 2. all return parameters are named
	for i, o := range um.OutParams {
		if o.Name == "" {
			o.Name = fmt.Sprintf("param%v", i)
		}
		um.OutParams[i] = o
	}
	// 3. make sure the error param is named err
	if um.OutParams[len(um.OutParams)-1].Name != "err" {
		x := um.OutParams[len(um.OutParams)-1]
		x.Name = "err"
		um.OutParams[len(um.OutParams)-1] = x
	}
	// 4. make sure names are uniq
	j := map[string]bool{}
	for _, o := range um.OutParams {
		_, ok := j[o.Name]
		if ok {
			return fmt.Errorf("non uniq result parameter (%v %v) in %v.%v.%v",
				o.Name, o.GoType, u.PackagePath, u.Name, um.Name)
		}
		j[o.Name] = true
	}
	// 5. make sure out params, except the trailing err, are either
	// a struct, or a slice of struct, or
	// a bunch of basic
	if um.OutputsToIterator() {
	} else if um.OutputsToStruct() {
	} else if um.OutputsToStructSlice() {
	} else if um.OutputsToBasic() {
	} else if um.OutputsOnlyError() {
	} else {
		return fmt.Errorf("ambiguous result parameters in %v.%v.%v",
			u.PackagePath, u.Name, um.Name)
	}
	// 6. input params can not be named
	// error ctx db
	for _, i := range um.InParams {
		if i.Name == "err" {
			return fmt.Errorf("invalid input parameter (%v %v) in %v.%v.%v: %v",
				i.Name, i.GoType, u.PackagePath, u.Name, um.Name, "it cannot be named error")
		}
		if i.Name == "db" {
			return fmt.Errorf("invalid input parameter (%v %v) in %v.%v.%v: %v",
				i.Name, i.GoType, u.PackagePath, u.Name, um.Name, "it cannot be named db")
		}
		if i.Name == "ctx" {
			return fmt.Errorf("invalid input parameter (%v %v) in %v.%v.%v: %v",
				i.Name, i.GoType, u.PackagePath, u.Name, um.Name, "it cannot be named ctx")
		}
	}
	// 7. Output struct types must have at least one property
	// if this is a struct
	for _, i := range um.OutParams[:len(um.OutParams)-1] {
		if i.BasicKind == types.Invalid &&
			len(i.StructProperties) < 1 && !i.IsFunc {
			return fmt.Errorf("invalid result type (%v %v) in %v.%v.%v: it must have at least one property",
				i.Name, i.StructName, u.PackagePath, u.Name, um.Name)
		}
	}
	// 8. LastInsertID and AffectedRows must be int64
	if um.InsertedID != "" {
		if p := um.Param(um.InsertedID); p != nil {
			if p.BasicKind != types.Int64 {
				return fmt.Errorf("invalid type (%v %v) in %v.%v.%v: it mustbe int64",
					p.Name, p.GoType, u.PackagePath, u.Name, um.Name)
			}
		}
	}
	if um.AffectedRows != "" {
		if p := um.Param(um.AffectedRows); p != nil {
			if p.BasicKind != types.Int64 {
				return fmt.Errorf("invalid type (%v %v) in %v.%v.%v: it mustbe int64",
					p.Name, p.GoType, u.PackagePath, u.Name, um.Name)
			}
		}
	}
	// 9. Output iterator types must have zero input parameters,
	// one trailing error and a least one non error output parameter.
	for _, i := range um.OutParams[:len(um.OutParams)-1] {
		if i.IsFunc {
			if len(i.Func.InParams) > 0 {
				return fmt.Errorf("invalid iterator type %v in %v.%v.%v: it must have zero input parameter",
					i.GoType, u.PackagePath, u.Name, um.Name)
			} else if len(i.Func.OutParams) < 1 {
				return fmt.Errorf("invalid iterator type %v in %v.%v.%v: it must have at least two output parameters",
					i.GoType, u.PackagePath, u.Name, um.Name)
			} else if i.Func.OutParams[len(i.Func.OutParams)-1].GoType != "error" {
				return fmt.Errorf("invalid iterator type %v in %v.%v.%v: it must have a trailing error output parameter",
					i.GoType, u.PackagePath, u.Name, um.Name)
			} else if len(i.Func.OutParams) < 2 {
				return fmt.Errorf("invalid iterator type %v in %v.%v.%v: it must have at least one output data parameter",
					i.GoType, u.PackagePath, u.Name, um.Name)
			}
		}
	}
	return nil
}

func tuplesToParams(pkgPath string, list *types.Tuple, allowFunc bool) (ret []userParam, err error) {
	for i := 0; i < list.Len(); i++ {
		o := list.At(i)
		var p userParam
		p.Name = o.Name()
		p, err = typeToParam(pkgPath, o.Type(), p, allowFunc)
		if err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	return ret, nil
}

func (u userMethod) Param(name string) *userParam {
	for _, p := range u.InParams {
		if p.Name == name {
			return &p
		}
	}
	for _, p := range u.OutParams {
		if p.Name == name {
			return &p
		}
	}
	return nil
}
func (u userMethod) OutputsToStruct() bool {
	if len(u.OutParams) != 2 {
		return false
	}
	if u.OutParams[0].Name == "_" {
		return false
	}
	if u.OutParams[0].BasicKind != types.Invalid {
		return false
	}
	return true
}
func (u userMethod) OutputsToIterator() bool {
	if len(u.OutParams) < 1 {
		return false
	}
	return u.OutParams[0].IsFunc
}
func (u userMethod) OutputsToStructSlice() bool {
	if len(u.OutParams) != 2 {
		return false
	}
	if u.OutParams[0].Name == "_" {
		return false
	}
	if u.OutParams[0].BasicKind != types.Invalid {
		return false
	}
	if !u.OutParams[0].IsSlice {
		return false
	}
	return true
}
func (u userMethod) OutputsToBasic() bool {
	if len(u.OutParams) < 2 {
		return false
	}
	for _, p := range u.OutParams[:len(u.OutParams)-1] {
		if p.Name == "_" {
			continue
		}
		if p.BasicKind == types.Invalid {
			return false
		}
	}
	return true
}
func (u userMethod) OutputsOnlyError() bool {
	if len(u.OutParams) != 1 {
		return false
	}
	if u.OutParams[len(u.OutParams)-1].GoType != "error" {
		return false
	}
	return true
}

func typeToParam(pkgPath string, ty types.Type, p userParam, allowFunc bool) (userParam, error) {

	if s, ok := ty.(*types.Pointer); ok {
		p.IsPtr = true
		p.GoType = "*"
		if ss, ok := s.Elem().(*types.Named); ok {
			if ss.Obj().Name() != "error" {
				p.PackagePath = ss.Obj().Pkg().Path()
				if ss.Obj().Pkg().Path() != pkgPath {
					p.GoType += ss.Obj().Pkg().Name() + "."
				}
				if stru, ok := ss.Underlying().(*types.Struct); ok {
					stru.NumFields()
					for i := 0; i < stru.NumFields(); i++ {
						f := stru.Field(i)
						if f.Anonymous() {
							continue
						}
						if !f.Exported() {
							continue
						}
						p.StructProperties = append(p.StructProperties, f.Name())
					}
				}
			}
			p.GoType += ss.Obj().Name()
			p.StructName = ss.Obj().Name()
		} else {
			return p, fmt.Errorf("invalid signature variable type %s", ty.String())
		}

	} else if s, ok := ty.(*types.Named); ok {
		if s.Obj().Name() != "error" {
			if stru, ok := s.Underlying().(*types.Struct); ok {
				p.PackagePath = s.Obj().Pkg().Path()
				if s.Obj().Pkg().Path() != pkgPath {
					p.GoType = s.Obj().Pkg().Name() + "."
				}
				for i := 0; i < stru.NumFields(); i++ {
					f := stru.Field(i)
					if f.Anonymous() {
						continue
					}
					if !f.Exported() {
						continue
					}
					p.StructProperties = append(p.StructProperties, f.Name())
				}
			} else if s, ok := ty.(*types.Named); ok {
				if ss, ok := s.Underlying().(*types.Signature); ok {
					if !allowFunc {
						return p, fmt.Errorf("recursive function definition not allowed")
					}
					p.IsFunc = true
					p.GoType = ss.String()
					var err error
					p.Func.InParams, err = tuplesToParams(pkgPath, ss.Params(), false)
					if err != nil {
						return p, fmt.Errorf("invalid signature in %v: %v", ss.String(), err)
					}
					p.Func.OutParams, err = tuplesToParams(pkgPath, ss.Results(), false)
					if err != nil {
						return p, fmt.Errorf("invalid signature in %v: %v", ss.String(), err)
					}
					p.GoType = p.Func.String()
				} else {
					return p, fmt.Errorf("unsupported type %s", ty.String())
				}

			} else {
				return p, fmt.Errorf("unsupported type %s", ty.String())
			}
		}
		p.GoType += s.Obj().Name()
		p.StructName = s.Obj().Name()

	} else if s, ok := ty.(*types.Slice); ok {
		p.GoType = "[]"
		p.IsSlice = true
		if ss, ok := s.Elem().(*types.Named); ok {
			if ss.Obj().Name() != "error" {
				p.PackagePath = ss.Obj().Pkg().Path()
				if ss.Obj().Pkg().Path() != pkgPath {
					p.GoType += ss.Obj().Pkg().Name() + "."
				}
				if stru, ok := ss.Underlying().(*types.Struct); ok {
					stru.NumFields()
					for i := 0; i < stru.NumFields(); i++ {
						f := stru.Field(i)
						if f.Anonymous() {
							continue
						}
						if !f.Exported() {
							continue
						}
						p.StructProperties = append(p.StructProperties, f.Name())
					}
				} else {
					return p, fmt.Errorf("invalid signature variable type %s", ty.String())
				}
			}
			p.GoType += ss.Obj().Name()
			p.StructName = ss.Obj().Name()

		} else if ss, ok := s.Elem().(*types.Basic); ok {
			p.GoType += ss.Name()
			p.BasicKind = ss.Kind()
			p.BasicInfo = ss.Info()

		} else {
			return p, fmt.Errorf("invalid signature variable type %s", ty.String())
		}

	} else if s, ok := ty.(*types.Basic); ok {
		p.GoType = s.Name()
		p.BasicKind = s.Kind()
		p.BasicInfo = s.Info()

	} else if s, ok := ty.(*types.Signature); ok {
		if !allowFunc {
			return p, fmt.Errorf("recursive function definition not allowed")
		}
		p.IsFunc = true
		var err error
		p.Func.InParams, err = tuplesToParams(pkgPath, s.Params(), false)
		if err != nil {
			return p, fmt.Errorf("invalid signature in %v: %v", s.String(), err)
		}
		p.Func.OutParams, err = tuplesToParams(pkgPath, s.Results(), false)
		if err != nil {
			return p, fmt.Errorf("invalid signature in %v: %v", s.String(), err)
		}
		p.GoType = p.Func.String()

	} else {
		return p, fmt.Errorf("invalid signature variable type %s", ty.String())
	}

	return p, nil
}

func pathEnclosingInterval(fset *token.FileSet, files []*ast.File, start, end token.Pos) (path []ast.Node, exact bool) {
	for _, f := range files {
		if !tokenFileContainsPos(fset.File(f.Pos()), start) {
			continue
		}
		if path, exact := astutil.PathEnclosingInterval(f, start, end); path != nil {
			return path, exact
		}
	}
	return nil, false
}

func fileForPos(fset *token.FileSet, files []*ast.File, start, end token.Pos) (file *ast.File, exact bool) {
	for _, f := range files {
		if !tokenFileContainsPos(fset.File(f.Pos()), start) {
			continue
		}
		if path, _ := astutil.PathEnclosingInterval(f, start, end); path != nil {
			return f, true
		}
	}
	return nil, false
}

func tokenFileContainsPos(f *token.File, pos token.Pos) bool {
	p := int(pos)
	base := f.Base()
	return base <= p && p < base+f.Size()
}
