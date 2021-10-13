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
	"go/token"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/loader"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func NewModule(modFile string) (mod *Module, err error) {
	modFile = filepath.ToSlash(modFile)
	if !filepath.IsAbs(modFile) {
		modFileAbs, absErr := filepath.Abs(modFile)
		if absErr != nil {
			err = fmt.Errorf("fnc: new module failed for absolute representation of %s, %v", modFile, absErr)
			return
		}
		modFile = modFileAbs
	}
	modProjectPath, _ := filepath.Split(modFile)
	modProjectPath = modProjectPath[0 : len(modProjectPath)-1]
	p, readErr := ioutil.ReadFile(modFile)
	if readErr != nil {
		err = fmt.Errorf("fnc: read mod file failed, %s, %v", modFile, readErr)
		return
	}
	mf, parseErr := modfile.Parse(modFile, p, func(path, version string) (string, error) {
		return version, nil
	})
	if parseErr != nil {
		err = fmt.Errorf("fnc: parse mod file failed, %s, %v", modFile, parseErr)
		return
	}
	mod = &Module{}
	mod.GoVersion = mf.Go.Version
	mod.Name = mf.Module.Mod.Path
	mod.Path = modProjectPath
	mod.Requires = make([]Require, 0, 1)
	if mf.Require != nil {
		for _, require := range mf.Require {
			r := Require{
				Name:           require.Mod.Path,
				Version:        require.Mod.Version,
				Replace:        "",
				ReplaceVersion: "",
				Indirect:       require.Indirect,
			}
			mod.Requires = append(mod.Requires, r)
		}
	}
	if mf.Replace != nil {
		for _, replace := range mf.Replace {
			for i, require := range mod.Requires {
				if replace.Old.Path == require.Name {
					require.Replace = replace.New.Path
					require.ReplaceVersion = replace.New.Version
					mod.Requires[i] = require
				}
			}
		}
	}
	root, rootErr := loadProgram(mod.Name, mod.Path)
	if rootErr != nil {
		err = rootErr
		return
	}
	mod.Root = root
	mod.Deps = make(map[string]*loader.Program)
	mod.Structs = make(map[string]*Struct)
	return
}

type Module struct {
	Name      string
	Path      string
	GoVersion string
	Requires  []Require
	Root      *loader.Program
	Deps      map[string]*loader.Program
	Structs   map[string]*Struct
}

func (mod *Module) CreatedPackageInfos() (infos []*loader.PackageInfo) {
	infos = mod.Root.Created
	return
}

func (mod *Module) FileInfo(f *ast.File) (v *token.File) {
	v = mod.Root.Fset.File(f.Pos())
	return
}

func (mod *Module) GetPackageOfFile(f *ast.File) (name string, ident string, has bool) {
	for _, info := range mod.Root.Created {
		for _, file := range info.Files {
			if file == f {
				name = info.Pkg.Path()
				ident = info.Pkg.Name()
				has = true
				return
			}
		}
	}
	return
}

func (mod *Module) GetStruct(pkg string, pkgAlias string, name string) (v *Struct, has bool, err error) {
	key := fmt.Sprintf("%s:%s:%s", pkg, pkgAlias, name)
	v, has = mod.Structs[key]
	if has {
		return
	}
	v, has, err = mod.findStruct(pkg, pkgAlias, name)

	return
}

func (mod *Module) SetStruct(v *Struct) {
	key := fmt.Sprintf("%s:%s:%s", v.Package, v.PackageAlias, v.Name)
	mod.Structs[key] = v
	return
}

func (mod *Module) findStruct(pkg string, pkgAlias string, name string) (v *Struct, has bool, err error) {
	if pkg == "time" {
		if name == "Time" {
			v = &Struct{
				Package:      "time",
				PackageAlias: pkgAlias,
				Name:         "Time",
				Fields:       nil,
				Annotations:  nil,
			}
			has = true
			return
		}
	}
	if pkg == "encoding/json" {
		if name == "RawMessage" {
			v = &Struct{
				Package:      "encoding/json",
				PackageAlias: pkgAlias,
				Name:         "RawMessage",
				Fields:       nil,
				Annotations:  nil,
			}
			has = true
			return
		}
	}
	if pkg == "github.com/aacfactory/fns" {
		if name == "Empty" {
			v = &Struct{
				Package:      "github.com/aacfactory/fns",
				PackageAlias: pkgAlias,
				Name:         "Empty",
				Fields:       nil,
				Annotations:  nil,
			}
			has = true
			return
		}
	}
	if pkg == "github.com/aacfactory/json" {
		if name == "Date" {
			v = &Struct{
				Package:      "github.com/aacfactory/json",
				PackageAlias: pkgAlias,
				Name:         "Date",
				Fields:       nil,
				Annotations:  nil,
			}
			has = true
			return
		}
		if name == "RawMessage" {
			v = &Struct{
				Package:      "github.com/aacfactory/json",
				PackageAlias: pkgAlias,
				Name:         "RawMessage",
				Fields:       nil,
				Annotations:  nil,
			}
			has = true
			return
		}
		if name == "Object" {
			v = &Struct{
				Package:      "github.com/aacfactory/json",
				PackageAlias: pkgAlias,
				Name:         "Object",
				Fields:       nil,
				Annotations:  nil,
			}
			has = true
			return
		}
		if name == "Array" {
			v = &Struct{
				Package:      "github.com/aacfactory/json",
				PackageAlias: pkgAlias,
				Name:         "Array",
				Fields:       nil,
				Annotations:  nil,
			}
			has = true
			return
		}
	}
	if strings.Contains(pkg, mod.Name) {
		v, has, err = mod.findStructInProgram(pkg, pkgAlias, name, mod.Root)
		return
	}
	deps := make(map[string]*loader.Program)
	for s, program := range mod.Deps {
		if strings.Contains(pkg, s) {
			deps[s] = program
		}
	}
	if len(deps) == 0 {
		for _, require := range mod.Requires {
			if strings.Contains(pkg, require.Name) {
				dep, depErr := loadProgram(require.Name, require.Path())
				if depErr != nil {
					err = fmt.Errorf("fnc: load dep(%s) failed, %s, %v", require.Name, require.Path(), depErr)
					return
				}
				mod.Deps[require.Name] = dep
				deps[require.Name] = dep
			}
		}
	}

	if len(deps) == 0 {
		err = fmt.Errorf("fnc: no %s require in go.mod", pkg)
		return
	}
	for _, dep := range deps {
		v, has, err = mod.findStructInProgram(pkg, pkgAlias, name, dep)
		if has {
			return
		}
	}
	return
}

func (mod *Module) findStructInProgram(pkgName string, pkgAlias string, name string, program *loader.Program) (v *Struct, has bool, err error) {
	pkg := program.Package(pkgName)
	if pkg == nil {
		return
	}
	var spec *ast.TypeSpec
	var structType *ast.StructType
	var srcFile *ast.File
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
				structType0, structTypeOk := typeSpec.Type.(*ast.StructType)
				if structTypeOk {
					structType = structType0
				}
				break
			}
		}
		if spec != nil {
			break
		}
	}
	if structType == nil {
		return
	}
	imports := getImports(srcFile)
	doc := ""
	if spec.Doc != nil {
		doc = spec.Doc.Text()
	}
	v = &Struct{
		Package:      pkgName,
		PackageAlias: pkgAlias,
		Name:         name,
		Fields:       make([]*Field, 0, 1),
		Annotations:  getAnnotations(doc),
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
			field := &Field{
				Name:        fieldName,
				Tag:         nil,
				Type:        nil,
				Annotations: nil,
			}
			if fieldSpec.Tag != nil {
				field.Tag = getStructFieldTag(fieldSpec.Tag.Value)
			} else {
				field.Tag = make(map[string]string)
			}
			fieldDoc := ""
			if fieldSpec.Doc != nil {
				field.Annotations = getAnnotations(fieldDoc)
			}
			// kind && type
			typ, typeErr := newTypeFromSpec(mod, imports, pkgName, pkgAlias, fieldSpec.Type)
			if typeErr != nil {
				err = fmt.Errorf("fnc: parse %s/%s.%s failed for %v", pkgName, name, fieldName, typeErr)
				return
			}
			field.Type = typ

			v.Fields = append(v.Fields, field)
		}
	}
	has = true
	return
}

func (mod Module) String() (s string) {
	s = mod.Name + " " + mod.Path + "\n"
	s = s + mod.GoVersion + "\n"
	for _, require := range mod.Requires {
		s = s + require.Name + " " + require.Version
		if require.Replace != "" {
			s = s + " => " + require.Replace
		}
		s = s + " [" + require.Path() + "]"
		s = s + "\n"
	}
	return
}

type Require struct {
	Name           string
	Version        string
	Replace        string
	ReplaceVersion string
	Indirect       bool
}

func (r Require) Path() (s string) {
	if r.Replace != "" {
		if r.ReplaceVersion != "" {
			f, err := filepath.Abs(r.Replace)
			if err != nil {
				panic(fmt.Sprintf("fnc: get require file path failed, %s, %v", r.Replace, err))
			}
			s = f
		} else {
			gto, name, bv := moduleVersion(r.Replace)
			if gto {
				s = filepath.Join(gopathModPath(), name, fmt.Sprintf("v%d@%s", bv, r.ReplaceVersion))
			} else {
				s = filepath.Join(gopathModPath(), fmt.Sprintf("%s@%s", r.Replace, r.ReplaceVersion))
			}
		}
	} else {
		gto, name, bv := moduleVersion(r.Name)
		if gto {
			s = filepath.Join(gopathModPath(), name, fmt.Sprintf("v%d@%s", bv, r.Version))
		} else {
			s = filepath.Join(gopathModPath(), fmt.Sprintf("%s@%s", r.Name, r.Version))
		}
	}
	return
}

const (
	StructTypeKind               = "struct"
	ArrayTypeKind                = "array"
	MapTypeKind                  = "map"
	StringType                   = "string"
	BoolTypeKind                 = "bool"
	IntTypeKind                  = "int"
	Int8TypeKind                 = "int8"
	Int16TypeKind                = "int16"
	Int32TypeKind                = "int32"
	Int64TypeKind                = "int64"
	UIntTypeKind                 = "uint"
	UInt8TypeKind                = "uint8"
	UInt16TypeKind               = "uint16"
	UInt32TypeKind               = "uint32"
	UInt64TypeKind               = "uint64"
	Float32TypeKind              = "float32"
	Float64TypeKind              = "float64"
	TimeTypeKind                 = "time"
	FnsDateTypeKind              = "fns/date"
	JsonRawMessageTypeKind       = "json"
	FnsJsonRawMessageTypeKind    = "fns/json"
	FnsJsonObjectMessageTypeKind = "fns/json+object"
	FnsJsonArrayMessageTypeKind  = "fns/json+array"
)

type Type struct {
	Kind   string
	Indent string // build-in name, map key name
	Struct *Struct
	X      *Type
}
