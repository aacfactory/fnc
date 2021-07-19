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
	f, parseErr := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if parseErr != nil {
		err = fmt.Errorf("parse file %s failed, %v", path, parseErr)
		return
	}
	// path
	file.Path = path
	// file doc
	file.Comment = f.Doc.Text()
	// Package
	file.Package = f.Name.String()

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
		file.Imports = imports
	}

	// decl
	if f.Decls != nil {
		for _, decl := range f.Decls {
			switch decl.(type) {
			case *ast.GenDecl:
				gen := decl.(*ast.GenDecl)
				specs := gen.Specs
				if gen.Tok == token.CONST {
					for _, spec := range specs {
						constSpec, ok := spec.(*ast.ValueSpec)
						if !ok {
							err = fmt.Errorf("scan %s failed, get const gen decl[%d, %d], but can not convert to *ast.ValueSpec", path, gen.Pos(), gen.End())
							return
						}
						size := len(constSpec.Names)
						for i := 0; i < size; i++ {
							name := constSpec.Names[i].Name
							value := constSpec.Values[0].(*ast.BasicLit).Value
							comment := ""
							if constSpec.Comment != nil {
								comment = constSpec.Comment.Text()
							}
							type0 := ""
							if constSpec.Type != nil {
								type0 = constSpec.Type.(*ast.Ident).String()
							}
							file.Constants = append(file.Constants, Const{
								Exported: IsExported(name),
								Comment:  comment,
								Name:     name,
								Type:     type0,
								Value:    value,
							})
						}
					}
				} else if gen.Tok == token.VAR {
					for _, spec := range specs {
						valueSpec, ok := spec.(*ast.ValueSpec)
						if !ok {
							err = fmt.Errorf("scan %s failed, get var gen decl[%d, %d], but can not convert to *ast.ValueSpec", path, gen.Pos(), gen.End())
							return
						}
						size := len(valueSpec.Names)
						for i := 0; i < size; i++ {
							name := valueSpec.Names[i].Name
							value := valueSpec.Values[0].(*ast.BasicLit).Value
							comment := ""
							if valueSpec.Comment != nil {
								comment = valueSpec.Comment.Text()
							}
							type0 := ParseType(valueSpec.Type)
							file.Variables = append(file.Variables, Var{
								Exported: IsExported(name),
								Comment:  comment,
								Name:     name,
								Type:     type0,
								Value:    value,
							})
						}
					}
				} else if gen.Tok == token.TYPE {
					for _, spec := range specs {
						typeSpec, ok := spec.(*ast.TypeSpec)
						if !ok {
							err = fmt.Errorf("scan %s failed, get type gen decl[%d, %d], but can not convert to *ast.TypeSpec", path, gen.Pos(), gen.End())
							return
						}
						name := typeSpec.Name.String()
						fmt.Println("type", name, reflect.TypeOf(typeSpec.Type), typeSpec)
						switch typeSpec.Type.(type) {
						case *ast.StarExpr:
							// alias type
							ptrExpr := typeSpec.Type.(*ast.StarExpr)
							aliasType, isIdent := ptrExpr.X.(*ast.Ident)
							if isIdent {
								file.AliasTypes = append(file.AliasTypes, AliasType{
									Exported: aliasType.IsExported(),
									Name:     name,
									Type:     Type{
										Kind:       "indent",
										Package:    "",
										Name:       aliasType.String(),
										Ptr:        false,
										InnerTypes: nil,
									},
								})
							}
						case *ast.Ident:
							// alias type
							aliasType := typeSpec.Type.(*ast.Ident)
							file.AliasTypes = append(file.AliasTypes, AliasType{
								Exported: aliasType.IsExported(),
								Name:     name,
								Type:     Type{
									Kind:       "indent",
									Package:    "",
									Name:       aliasType.String(),
									Ptr:        false,
									InnerTypes: nil,
								},
							})
						case *ast.StructType:

						case *ast.FuncType:
							// alias type

						case *ast.InterfaceType:

						case *ast.ArrayType:
							// alias type

						case *ast.MapType:
							// alias type

						case *ast.ChanType:
							// alias type

						default:
							err = fmt.Errorf("scan %s failed, get type gen decl[%d, %d], but it is bad", path, gen.Pos(), gen.End())
							return
						}
					}
				}
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
