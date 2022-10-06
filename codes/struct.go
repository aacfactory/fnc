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
	"github.com/aacfactory/cases"
	"github.com/aacfactory/gcg"
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
		modRootName := mod.Name[strings.LastIndex(mod.Name, "/")+1:]
		depPkgPath := pkgPath[strings.LastIndex(pkgPath, modRootName):]
		for _, info := range mod.Program.Created {
			created := info.String()
			if strings.Index(created, "main") == 0 {
				created = mod.Name + created[4:]
			}
			if depPkgPath == created || pkgPath == created {
				v, has, err = findStructInProgram(info.String(), name, mod)
				if has {
					return
				}
			}
		}
		if !has {
			return
		}
		return
	}
	var srcFile *ast.File
	var spec *ast.TypeSpec
	doc := ""
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
				if genDecl != nil {
					doc = genDecl.Doc.Text()
				}
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
	v = &Struct{
		Package:     pkgPath,
		Name:        name,
		Fields:      make([]*Field, 0, 1),
		Annotations: getAnnotations(doc),
		Ref:         false,
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
			if fieldSpec.Doc != nil {
				annotations = getAnnotations(fieldSpec.Doc.Text())
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
	Ref         bool
}

func (s *Struct) Title() (title string) {
	v, has := s.Annotations["title"]
	if has {
		title = v
		return
	}
	return
}

func (s *Struct) Description() (description string) {
	v, has := s.Annotations["description"]
	if has {
		description = v
		return
	}
	return
}

func (s *Struct) Key() (key string) {
	key = s.Package + "." + s.Name
	return
}

func (s *Struct) ObjectKey() (v string) {
	v = strings.ReplaceAll(s.Key(), "/", "_")
	v = strings.ReplaceAll(v, ".", "_")
	atoms, _ := cases.Snake().Parse(v)
	v = cases.LowerCamel().Format(atoms)
	return
}

func (s *Struct) generateObjectRef() (code *gcg.Statement) {
	code = gcg.Statements()
	code.Token(fmt.Sprintf("documents.RefStruct(\"%s\", \"%s\")", s.Package, s.Name))
	return
}

func (s *Struct) generateObject(servicePkg string) (code *gcg.Statement) {
	code = gcg.Statements()
	if s.Ref {
		code.Token(fmt.Sprintf("documents.RefStruct(\"%s\", \"%s\")", s.Package, s.Name))
		return
	}
	if s.Package != servicePkg && s.Package != "main/repository" {
		code.Token(fmt.Sprintf("documents.RefStruct(\"%s\", \"%s\")", s.Package, s.Name))
		return
	}
	s.Ref = true
	code.Token(fmt.Sprintf("documents.Struct(\"%s\", \"%s\", \"%s\", \"%s\")", s.Package, s.Name, s.Title(), s.Description()))
	if s.Fields != nil && len(s.Fields) > 0 {
		code.Token(".").Line()
		i := 0
		for _, field := range s.Fields {
			i++
			code.Add(field.generateObject(servicePkg))
			if i != len(s.Fields) {
				code.Token(".").Line()
			}
		}
	}
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

func (x *Field) Validation() (required bool, validation string) {
	v, has := x.Tag["validate"]
	if has {
		required = strings.Contains(v, "required")
		message := x.Tag["message"]
		validation = fmt.Sprintf("validate:\"%s\" message:\"%s\"", v, message)
		return
	}
	return
}

func (x *Field) Deprecated() (deprecated bool) {
	v, has := x.Annotations["deprecated"]
	if has {
		deprecated = v == "true"
		return
	}
	return
}

func (x *Field) Enum() (enum string) {
	v, has := x.Annotations["enum"]
	if has {
		enum = v
		return
	}
	return
}

func (x *Field) generateObject(servicePKG string) (code *gcg.Statement) {
	code = gcg.Statements()
	key, hasKey := x.Tag["json"]
	if !hasKey {
		key = x.Name
	} else {
		keyItems := strings.Split(key, ",")
		if keyItems[0] == "-" {
			return
		}
		key = keyItems[0]
	}
	if key == "" {
		return
	}

	fc := x.Type.generateObject(servicePKG)
	if x.Title() != "" {
		fc.Token(fmt.Sprintf(".SetTitle(\"%s\")", x.Title()))
	}
	if x.Description() != "" {
		fc.Token(fmt.Sprintf(".SetDescription(\"%s\")", x.Description()))
	}
	require, validation := x.Validation()
	if require {
		fc.Token(fmt.Sprintf(".AsRequired(`%s`)", validation))
	}
	if x.Deprecated() {
		fc.Token(fmt.Sprintf(".AsDeprecated()"))
	}
	// AddEnum
	enum := x.Enum()
	if enum != "" && x.Type.IsBuiltin() {
		enums := ""
		enumItems := strings.Split(enum, ",")
		if x.Type.Indent == "string" {
			for _, item := range enumItems {
				enums = enums + "," + `"` + item + `"`
			}
			enums = enums[1:]
		} else {
			for _, item := range enumItems {
				enums = "," + item
			}
			enums = enums[1:]
		}
		fc.Token(fmt.Sprintf(".AddEnum(%s)", enums))
	}
	code.Token(fmt.Sprintf("AddProperty(\"%s\", ", key)).Line().Add(fc).Token(",").Line().Token(")")
	return
}
