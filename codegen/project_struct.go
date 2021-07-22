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

package codegen

import (
	"fmt"
	"go/ast"
)

func (p *Project) FindStruct(pkgPath string, name string) (fullName string, str Struct, has bool) {
	fullName = fmt.Sprintf("%s.%s", pkgPath, name)
	str, has = p.Structs[fullName]
	if has {
		return
	}
	if pkgPath == "time" {
		fullName = "time.Time"
		has = true
		return
	}
	if pkgPath == "encoding/json" {
		if name == "RawMessage" {
			fullName = "encoding/json.RawMessage"
			has = true
		}
		return
	}
	processing := p.tmpStructMap[fullName]
	if processing {
		has = true
		return
	}

	p.tmpStructMap[fullName] = true
	defer delete(p.tmpStructMap, fullName)

	spec, doc, file, hasSpec := p.getStructSpecType(pkgPath, name)
	if !hasSpec {
		return
	}
	str0, ok := p.getStructFromTypeSpec(file, spec, doc)
	if !ok {
		return
	}
	str = str0
	has = true
	p.Structs[fullName] = str0
	Log().Debugf("fnc load struct type: %v", str)
	return
}

func (p *Project) getStructSpecType(pkgPath string, name string) (spec *ast.TypeSpec, doc string, fileOfSpec *ast.File, has bool) {
	pkg := p.Program.Package(pkgPath)
	if pkg == nil {
		for _, program := range p.ImportPrograms {
			pkg = program.Package(pkgPath)
			if pkg != nil {
				break
			}
		}
	}
	if pkg == nil {
		return
	}
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			genDecl, genDeclOk := decl.(*ast.GenDecl)
			if !genDeclOk {
				continue
			}
			typeSpec, typeSpecOk := genDecl.Specs[0].(*ast.TypeSpec)
			if !typeSpecOk {
				continue
			}
			if typeSpec.Name.Name == name {
				doc = genDecl.Doc.Text()
				fileOfSpec = file
				spec = typeSpec
				has = true
				return
			}
		}
	}
	return
}

func (p *Project) getStructFromTypeSpec(file *ast.File, spec *ast.TypeSpec, doc string) (str Struct, ok bool) {
	pkgPath, pkgName, hasPkgPath := p.PackageOfFile(file)
	if !hasPkgPath {
		return
	}

	expr, isStructType := spec.Type.(*ast.StructType)
	if !isStructType {
		return
	}

	imports := p.GetImports(file)
	str.Name = spec.Name.Name
	str.Package = Import{
		Path: pkgPath,
		Name: pkgName,
	}
	str.Exported = ast.IsExported(str.Name)
	str.Doc = parseDoc(doc)

	if expr.Fields != nil && expr.Fields.NumFields() > 0 {
		for i := 0; i < expr.Fields.NumFields(); i++ {
			field := expr.Fields.List[i]
			if field.Names == nil || len(field.Names) != 1 {
				switch field.Type.(type) {
				case *ast.Ident:
					joinExpr := field.Type.(*ast.Ident)
					_, joined, hasJoined := p.FindStruct(pkgPath, joinExpr.Name)
					if hasJoined {
						for _, f := range joined.Fields {
							str.PutField(f)
						}
					}
				case *ast.SelectorExpr:
					joinExpr := field.Type.(*ast.SelectorExpr)
					structName := joinExpr.Sel.Name
					ident, identOk := joinExpr.X.(*ast.Ident)
					if !identOk {
						continue
					}

					fieldPkgName := ident.Name
					fieldPkgPath := ""
					for _, import0 := range imports {
						if import0.Name == fieldPkgName {
							fieldPkgPath = import0.Path
							continue
						}
					}
					if fieldPkgPath == "" {
						continue
					}

					_, fieldStruct, defined := p.FindStruct(fieldPkgPath, structName)
					if defined {
						for _, f := range fieldStruct.Fields {
							str.PutField(f)
						}
					}
				default:
					continue
				}
			} else {
				fd, fieldOk := p.parseStructField(file, field)
				if fieldOk {
					str.PutField(fd)
				}
			}
		}
	}
	ok = true
	return
}
