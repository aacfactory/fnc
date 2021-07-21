/*
 * Copyright 2021 Wang Min Xiang
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
							Log().Debugf("const %s, type is %v", name, type0)
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
							// todo parse value
							//value := valueSpec.Values[0].(*ast.BasicLit).Value
							comment := ""
							if valueSpec.Comment != nil {
								comment = valueSpec.Comment.Text()
							}
							type0, typeOk := ParseType(valueSpec.Type)
							if typeOk {
								file.Variables = append(file.Variables, Var{
									Exported: IsExported(name),
									Comment:  comment,
									Name:     name,
									Type:     type0,
									Value:    "",
								})
							}
							Log().Debugf("var %s, type is %v", name, type0)
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
						switch typeSpec.Type.(type) {
						case *ast.StarExpr:
							Log().Debugf("type %s, type is %v, it is not supported", name, reflect.TypeOf(typeSpec.Type))
						case *ast.Ident:
							Log().Debugf("type %s, type is %v, it is not supported", name, reflect.TypeOf(typeSpec.Type))
						case *ast.StructType:
							struct0, structOk := ParseStruct(typeSpec)
							if structOk {
								file.Structs = append(file.Structs, struct0)
							}
						case *ast.FuncType:
							// alias type
							Log().Debugf("type %s, type is %v, it is not supported", name, reflect.TypeOf(typeSpec.Type))
						case *ast.InterfaceType:
							Log().Debugf("type %s, type is %v, it is not supported", name, reflect.TypeOf(typeSpec.Type))
						case *ast.ArrayType:
							// alias type
							Log().Debugf("type %s, type is %v, it is not supported", name, reflect.TypeOf(typeSpec.Type))
						case *ast.MapType:
							// alias type
							Log().Debugf("type %s, type is %v, it is not supported", name, reflect.TypeOf(typeSpec.Type))
						case *ast.ChanType:
							// alias type
							Log().Debugf("type %s, type is %v, it is not supported", name, reflect.TypeOf(typeSpec.Type))

						default:
							Log().Debugf("type %s, type is %v, it is not supported", name, reflect.TypeOf(typeSpec.Type))
							return
						}
					}
				}
			case *ast.FuncDecl:
				funcDecl := decl.(*ast.FuncDecl)
				fn, fnOk := ParseFuncDecl(funcDecl)
				if fnOk {
					file.Functions = append(file.Functions, fn)
				}
			case *ast.BadDecl:
				// error

			}
		}

	}
	return
}
