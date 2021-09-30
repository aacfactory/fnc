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
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/loader"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func NewModule(modFilePath string) (mod *Module, err error) {
	p, readErr := ioutil.ReadFile(modFilePath)
	if readErr != nil {
		err = fmt.Errorf("fnc: read mod file failed, %s, %v", modFilePath, readErr)
		return
	}
	mf, parseErr := modfile.Parse(modFilePath, p, func(path, version string) (string, error) {
		return version, nil
	})
	if parseErr != nil {
		err = fmt.Errorf("fnc: parse mod file failed, %s, %v", modFilePath, parseErr)
		return
	}
	mod = &Module{}
	mod.GoVersion = mf.Go.Version
	mod.Name = mf.Module.Mod.Path
	mod.Path = fmt.Sprintf("%s/src/%s", gopath(), mod.Name)
	mod.Requires = make([]Require, 0, 1)
	if mf.Require != nil {
		for _, require := range mf.Require {
			r := Require{
				Name:    require.Mod.Path,
				Version: require.Mod.Version,
				Replace: "",
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

func (mod *Module) GetStruct(pkg string, pkgAlias string, name string) (v *Struct, has bool) {
	key := fmt.Sprintf("%s:%s:%s", pkg, pkgAlias, name)
	v, has = mod.Structs[key]
	if has {
		return
	}
	v, has = mod.findStruct(pkg, pkgAlias, name)
	if has {
		mod.SetStruct(v)
	}
	return
}

func (mod *Module) SetStruct(v *Struct) {
	key := fmt.Sprintf("%s:%s:%s", v.Package, v.PackageAlias, v.Name)
	mod.Structs[key] = v
	return
}

func (mod *Module) findStruct(pkg string, pkgAlias string, name string) (v *Struct, has bool) {
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
		v, has = mod.findStructInProgram(pkg, pkgAlias, name, mod.Root)
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
					panic(fmt.Sprintf("fnc: load dep(%s) failed, %s, %v", require.Name, require.Path(), depErr))
					return
				}
				mod.Deps[require.Name] = dep
				deps[require.Name] = dep
			}
		}
	}

	if len(deps) == 0 {
		panic(fmt.Sprintf("fnc: no %s require in go.mod", pkg))
		return
	}
	for _, dep := range deps {
		v, has = mod.findStructInProgram(pkg, pkgAlias, name, dep)
		if has {
			return
		}
	}
	return
}

func (mod *Module) findStructInProgram(pkgName string, pkgAlias string, name string, program *loader.Program) (v *Struct, has bool) {
	pkg := program.Package(pkgName)
	if pkg == nil {
		return
	}
	var spec *ast.TypeSpec
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
				break
			}
		}
		if spec != nil {
			break
		}
	}
	imports := getImports(srcFile)
	doc := spec.Doc.Text()
	v = &Struct{
		Package:      pkgName,
		PackageAlias: pkgAlias,
		Name:         name,
		Fields:       make([]*Field, 0, 1),
		Annotations:  getAnnotations(doc),
	}
	// todo
	fmt.Println(imports)


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
				s = fmt.Sprintf("%s/%s/v%d@%s", gopathModPath(), name, bv, r.ReplaceVersion)
			} else {
				s = fmt.Sprintf("%s/%s@%s", gopathModPath(), r.Replace, r.ReplaceVersion)
			}
		}
	} else {
		gto, name, bv := moduleVersion(r.Name)
		if gto {
			s = fmt.Sprintf("%s/%s/v%d@%s", gopathModPath(), name, bv, r.Version)
		} else {
			s = fmt.Sprintf("%s/%s@%s", gopathModPath(), r.Name, r.Version)
		}
	}
	return
}
