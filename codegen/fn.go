package codegen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"strings"
)

func LoadFn(project *Project) (fnFiles []FnFile, err error) {

	for _, filename := range project.Filenames {
		fmt.Println(filename)
		fnFile, has, fileErr := loadFnFile(project, filename)
		if fileErr != nil {
			err = fmt.Errorf("load fn file failed, %v", fileErr)
			return
		}
		if has {
			fnFiles = append(fnFiles, fnFile)
		}
	}

	return
}

func loadFnFile(project *Project, filename string) (fnFile FnFile, has bool, err error) {

	f, parseErr := parser.ParseFile(project.Program.Fset, filename, nil, parser.ParseComments)
	if parseErr != nil {
		err = fmt.Errorf("parse file %s failed, %v", filename, parseErr)
		return
	}
	// path
	fnFile.Path = filename
	// file doc
	fnFile.Comment = f.Doc.Text()
	// Package
	fnFile.Package = f.Name.String()
	// imports
	if f.Imports != nil {
		imports := make([]Import, 0, 1)
		for _, spec := range f.Imports {
			name := spec.Path.Value
			alias := spec.Name.String()
			anonymous := "_" == alias
			static := "." == alias
			import0 := Import{
				Name:      name,
				Alias:     alias,
				Anonymous: anonymous,
				Static:    static,
			}
			imports = append(imports, import0)
		}
		fnFile.Imports = imports
	}

	// func
	fns := make([]Fn, 0, 1)

	if f.Decls != nil {
		for _, decl := range f.Decls {
			funcDecl, fnOk := decl.(*ast.FuncDecl)
			if !fnOk {
				continue
			}
			// not method
			if funcDecl.Recv != nil {
				continue
			}
			// not doc
			if funcDecl.Doc == nil {
				continue
			}
			comments := funcDecl.Doc.List
			if comments == nil || len(comments) == 0 {
				continue
			}
			address := ""
			proxy := make([]string, 0, 1)
			openAPIContent := ""
			for _, comment := range comments {
				if comment == nil {
					continue
				}
				line := comment.Text
				if strings.Contains(line, "@Fn") {
					address = strings.TrimSpace(line[strings.Index(line, "@Fn")+4:])
					continue
				}
				if strings.Contains(line, "@Proxy") {
					proxy = strings.Split(strings.TrimSpace(line[strings.Index(line, "@Proxy")+7:]), " ")
					continue
				}
				if strings.Contains(line, "@OpenAPI") {
					openAPIContent = strings.TrimSpace(line[strings.Index(line, "@OpenAPI")+9:])
					continue
				}
			}
			if address == "" {
				continue
			}
			fn := Fn{}
			fn.Address = address
			fn.Proxy = proxy
			fn.OpenAPI = openAPIContent
			fn.Name = funcDecl.Name.Name
			fn.Exported = ast.IsExported(fn.Name)
			fn.Comment = funcDecl.Doc.Text()
			// in
			params := funcDecl.Type.Params
			if params == nil || len(params.List) != 2 {
				err = fmt.Errorf("get fn %s %s, but params is invalied, must has two params, first is fns.FnContext, secend is a struct typed", filename, fn.Name)
				return
			}
			// fp 1
			fp := params.List[0]
			if fp.Names == nil || len(fp.Names) != 1 {
				err = fmt.Errorf("get fn %s %s, but params is invalied, must has two params, first is fns.FnContext, secend is a struct typed", filename, fn.Name)
				return
			}
			fpName := fp.Names[0].Name
			fpTypeExpr, fpTypeExprOk := fp.Type.(*ast.SelectorExpr)
			if !fpTypeExprOk {
				err = fmt.Errorf("get fn %s %s, but params is invalied, must has two params, first is fns.FnContext, secend is a struct typed", filename, fn.Name)
				return
			}
			fpTypeStructName := fpTypeExpr.Sel.Name
			fpTypePkgExpr, fpTypePkgExprOk := fpTypeExpr.X.(*ast.Ident)
			if !fpTypePkgExprOk {
				err = fmt.Errorf("get fn %s %s, but params is invalied, must has two params, first is fns.FnContext, secend is a struct typed", filename, fn.Name)
				return
			}
			fpTypeStructPkg := fpTypePkgExpr.Name

			if !(fpTypeStructPkg == "fns" && fpTypeStructName == "FnContext") {
				err = fmt.Errorf("get fn %s %s, but params is invalied, must has two params, first is fns.FnContext, secend is a struct typed", filename, fn.Name)
				return
			}

			fpItem := FuncItem{
				Name: fpName,
				Type: Type{
					Kind:      "interface",
					Package:   "fns",
					Name:      "FnContext",
					InnerType: nil,
				},
			}
			fn.In = append(fn.In, fpItem)
			// fp 2
			sp := params.List[1]
			if sp.Names == nil || len(sp.Names) != 1 {
				err = fmt.Errorf("get fn %s %s, but params is invalied, must has two params, first is fns.FnContext, secend is a struct typed", filename, fn.Name)
				return
			}
			spName := sp.Names[0].Name

			spItem := FuncItem{
				Name: spName,
				Type: Type{},
			}
			switch sp.Type.(type) {
			case *ast.StarExpr:

			case *ast.SelectorExpr:

			case *ast.Ident:
				ident := sp.Type.(*ast.Ident)
				if ident.Obj == nil || ident.Obj.Kind != ast.Typ {
					err = fmt.Errorf("get fn %s %s, but params is invalied, must has two params, first is fns.FnContext, secend is a struct typed", filename, fn.Name)
					return
				}
				spec, specOk := ident.Obj.Decl.(*ast.TypeSpec)
				if !specOk {
					err = fmt.Errorf("get fn %s %s, but params is invalied, must has two params, first is fns.FnContext, secend is a struct typed", filename, fn.Name)
					return
				}
				str, parsed := parseStruct(spec)
				if !parsed {
					err = fmt.Errorf("get fn %s %s, but params is invalied, must has two params, first is fns.FnContext, secend is a struct typed", filename, fn.Name)
					return
				}
				spItem.Type.Name = ident.Obj.Name
				spItem.Type.Struct = &str
				spItem.Type.Kind = "struct"
				fmt.Println(ident.Name, ident.Obj.Name, ident.Obj.Kind, str)
			default:
				err = fmt.Errorf("get fn %s %s, but params is invalied, must has two params, first is fns.FnContext, secend is a struct typed", filename, fn.Name)
				return
			}
			//spTypeExpr, spTypeExprOk := sp.Type.(*ast.SelectorExpr)
			//if !spTypeExprOk {
			//	err = fmt.Errorf("get fn %s %s, but params is invalied, must has two params, first is fns.FnContext, secend is a struct typed", filename, fn.Name)
			//	return
			//}

			//fmt.Println(fn.Name, spName, reflect.TypeOf(sp.Type))

			// out
		}
	}

	if len(fns) > 0 {
		fnFile.Functions = fns
		has = true
	}

	return
}

type FnFile struct {
	Path      string
	Comment   string
	Package   string
	Imports   []Import
	Functions []Fn
}

type Import struct {
	Name      string
	Alias     string
	Anonymous bool
	Static    bool
}

type FuncItem struct {
	Name string
	Type Type
}

type Fn struct {
	Exported bool
	Comment  string
	Address  string
	Proxy    []string
	OpenAPI  string
	Name     string
	In       []FuncItem
	Out      []FuncItem
}

type Type struct {
	Kind      string // struct ptr array error basic
	Package   string
	Name      string
	Struct    *Struct
	InnerType *Type
}

type Struct struct {
	Exported bool
	Comment  string
	Name     string
	Fields   []Field
}

type Field struct {
	Exported bool
	Comment  string
	Name     string
	Type     Type
	Tags     []FieldTag
}

type FieldTag struct {
	Name   string
	Values []string
}
