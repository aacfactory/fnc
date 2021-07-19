package goparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
)

func ParseFile(path string) (file File, err error) {
	fs := token.NewFileSet()
	f, parseErr := parser.ParseFile(fs, path, nil, parser.ParseComments | parser.AllErrors)
	if parseErr != nil {
		err = fmt.Errorf("parse file %s failed, %v", path, parseErr)
		return
	}
	// path
	file.Path = path
	// file doc
	file.Comment =  f.Doc.Text()
	// Package
	fmt.Println("doc", f.Doc.Text())
	fmt.Println("pkg", f.Name)

	// imports
	if f.Imports != nil {
		for _, spec := range f.Imports {
			// name == _ -> ano , name == . -> Static,
			fmt.Println("import", spec.Name, spec.Path.Value)
		}
	}

	// decl
	if f.Decls != nil {
		for _, decl := range f.Decls {
			fmt.Println("decl", reflect.TypeOf(decl), decl)
			switch decl.(type) {
			case *ast.GenDecl:
				// import const var type
				// d := decl.(*ast.GenDecl)
				// d.Tok == ?
				// d.Specs
			case *ast.FuncDecl:

			case *ast.BadDecl:
				// error

			}
		}

	}
	for _, object := range f.Scope.Objects {
		fmt.Println("object", object)
	}
	return
}
