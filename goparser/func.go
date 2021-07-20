package goparser

import (
	"fmt"
	"go/ast"
)

func ParseFuncDecl(decl *ast.FuncDecl) (fn Func, ok bool) {

	fn.Name = decl.Name.Name
	fn.Exported = IsExported(fn.Name)
	fn.Comment = decl.Doc.Text()

	if decl.Recv != nil {
		// todo
		fn.Receiver = &FuncReceiver{
			Name: "x",
			Type: Struct{},
		}
	}

	// params
	if decl.Type == nil {
		return
	}

	params := decl.Type.Params
	if params != nil && params.NumFields() > 0 {
		for _, param := range params.List {
			if param.Names != nil {
				for _, name := range param.Names {
					p := FuncItem{}
					p.Name = name.Name
					typ, pTypeOk := ParseType(param.Type)
					fmt.Println(typ, pTypeOk)
				}
			}
		}
	}

	// results

	ok = true
	return
}
