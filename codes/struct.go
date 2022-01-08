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

package codes

import (
	"fmt"
	"go/ast"
	"strings"
)

func NewStruct(pkg string, name string, mod *Module) (s *Struct, err error) {
	s0, has := mod.GetStruct(pkg, name)
	if has {
		s = s0
		return
	}
	mods := make(map[string]*Module)
	if strings.Index(pkg, mod.Name) == 0 {
		pkg = strings.Replace(pkg, mod.Name, "main", 1)
	}
	if strings.Index(pkg, "main") == 0 {
		mods[mod.Name] = mod
	} else {
		deps, depsErr := mod.FindDeps(pkg)
		if depsErr != nil {
			err = fmt.Errorf("decode %s.%s struct failed for %v", pkg, name, depsErr)
			return
		}
		for _, dep := range deps {
			mods[dep.Name] = dep
		}
	}
	if len(mods) == 0 {
		err = fmt.Errorf("decode %s.%s struct failed for not found", pkg, name)
		return
	}
	for _, mod0 := range mods {
		s1, found, foundErr := findStructInProgram(pkg, name, mod0)
		if foundErr != nil {
			err = foundErr
			return
		}
		if found {
			s = s1
			return
		}
	}
	err = fmt.Errorf("decode %s.%s struct failed for not found", pkg, name)
	return
}

func findStructInProgram(pkgPath string, name string, mod *Module) (v *Struct, has bool, err error) {
	pkg := mod.Program.Package(pkgPath)
	if pkg == nil {
		return
	}
	var srcFile *ast.File
	var spec *ast.TypeSpec
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
				srcFile = file
				spec = typeSpec
				break
			}
		}
		if spec != nil {
			break
		}
	}
	structType, structTypeOk := spec.Type.(*ast.StructType)
	if !structTypeOk {
		return
	}
	imports := NewImports(srcFile)
	doc := ""
	if spec.Doc != nil {
		doc = spec.Doc.Text()
	}
	v = &Struct{
		Package:     pkgPath,
		Name:        name,
		Fields:      make([]*Field, 0, 1),
		Annotations: getAnnotations(doc),
	}
	mod.SetStruct(v)
	if structType.Fields != nil && structType.Fields.NumFields() > 0 {
		for i := 0; i < structType.Fields.NumFields(); i++ {
			fieldSpec := structType.Fields.List[i]
			fieldName := fieldSpec.Names[0].Name
			fieldNameExported := ast.IsExported(fieldName)
			if !fieldNameExported {
				continue
			}
			tags := make(map[string]string)
			if fieldSpec.Tag != nil {
				tags = getStructFieldTag(fieldSpec.Tag.Value)
			} else {
				tags = make(map[string]string)
			}
			if jsonTag, hasTag := tags["json"]; hasTag && jsonTag == "-" {
				continue
			}
			annotations := make(map[string]string)
			fieldDoc := ""
			if fieldSpec.Doc != nil {
				annotations = getAnnotations(fieldDoc)
			}
			// kind && type
			typ, typeErr := NewType(fieldSpec.Type, pkgPath, imports, mod)
			if typeErr != nil {
				err = fmt.Errorf("parse %s/%s.%s failed for %v", pkgPath, name, fieldName, typeErr)
				return
			}
			field := &Field{
				Name:        fieldName,
				Tag:         tags,
				Type:        typ,
				Annotations: annotations,
			}
			v.Fields = append(v.Fields, field)
		}
	}
	has = true
	return
}

type Struct struct {
	Package     string
	Name        string
	Fields      []*Field
	Annotations map[string]string
}

func (s Struct) Key() (key string) {
	key = s.Package + "." + s.Name
	return
}

type Field struct {
	Name        string
	Tag         map[string]string
	Type        *Type
	Annotations map[string]string
}

func (x *Field) Title() (title string) {
	v, has := x.Annotations["title"]
	if has {
		title = v
		return
	}
	title = x.Type.Struct.Key()
	return
}

func (x *Field) Description() (description string) {
	v, has := x.Annotations["description"]
	if has {
		description = v
		return
	}
	return
}
