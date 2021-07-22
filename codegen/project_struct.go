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
	"container/list"
	"fmt"
	"go/ast"
	"strings"
)

type StructFindStack struct {
	structs *list.List
}

func NewStructFindStack() *StructFindStack {
	structs := list.New()
	return &StructFindStack{structs}
}

func (stack *StructFindStack) Push(fullName string) {
	stack.structs.PushBack(fullName)
}

func (stack *StructFindStack) Pop() (fullName string) {
	e := stack.structs.Back()
	if e != nil {
		stack.structs.Remove(e)
		return e.Value.(string)
	}
	return ""
}

func (stack *StructFindStack) Peak() (fullName string) {
	e := stack.structs.Back()
	if e != nil {
		return e.Value.(string)
	}

	return ""
}

func (stack *StructFindStack) Len() int {
	return stack.structs.Len()
}

func (stack *StructFindStack) Empty() bool {
	return stack.structs.Len() == 0
}

func (p *Project) FindStructSpecType(pkgPath string, name string) (spec *ast.TypeSpec, fileOfSpec *ast.File, has bool) {
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
				fileOfSpec = file
				spec = typeSpec
				has = true
				return
			}
		}
	}
	return
}

func (p *Project) FindStruct(pkgPath string, name string) (str Struct, has bool) {
	key := fmt.Sprintf("%s.%s", pkgPath, name)
	str, has = p.Structs[key]
	if has {
		return
	}
	if pkgPath == "time" {
		str = Struct{
			Exported: true,
			Doc:      nil,
			Package: Import{
				Path: "time",
				Name: "time",
			},
			Name:   name,
			Fields: nil,
		}
		has = true
		return
	}
	if pkgPath == "encoding/json" {
		if name == "RawMessage" {
			str = Struct{
				Exported: true,
				Doc:      nil,
				Package: Import{
					Path: "encoding/json",
					Name: "json",
				},
				Name:   "RawMessage",
				Fields: nil,
			}
			has = true
		}
		return
	}
	prev := p.structFindStack.Peak()
	if prev == key {
		return
	}
	p.structFindStack.Push(key)
	defer p.structFindStack.Pop()
	spec, file, hasSpec := p.FindStructSpecType(pkgPath, name)
	if !hasSpec {
		return
	}
	str0, ok := p.parseStructFromTypeSpec(file, spec)
	if !ok {
		return
	}
	str = str0
	has = true
	p.Structs[key] = str
	Log().Debugf("fnc load struct type: %s %s", str.Package.Path, str.Name)
	return
}

func (p *Project) parseStructFromTypeSpec(file *ast.File, spec *ast.TypeSpec) (str Struct, ok bool) {
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
	if spec.Doc != nil {
		comments := spec.Doc.List
		if comments != nil {
			for _, comment := range comments {
				content := strings.TrimSpace(comment.Text)
				if content != "" {
					str.Doc = append(str.Doc, content)
				}
			}
		}
	}
	if expr.Fields != nil && expr.Fields.NumFields() > 0 {
		for i := 0; i < expr.Fields.NumFields(); i++ {
			field := expr.Fields.List[i]
			if field.Names == nil || len(field.Names) != 1 {
				switch field.Type.(type) {
				case *ast.Ident:
					joinExpr := field.Type.(*ast.Ident)
					joined, hasJoined := p.FindStruct(pkgPath, joinExpr.Name)
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

					fieldStruct, defined := p.FindStruct(fieldPkgPath, structName)
					if defined {
						for _, f := range fieldStruct.Fields {
							str.PutField(f)
						}
					}
				default:
					continue
				}
			} else {
				fd, fieldOk := p.parseStructFieldFromStructType(file, field)
				if fieldOk {
					str.PutField(fd)
				}
			}
		}
	}
	ok = true
	return
}
