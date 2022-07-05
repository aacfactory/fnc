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
	program, programErr := loadProgram(mod.Name, mod.Path)
	if programErr != nil {
		err = programErr
		return
	}
	mod.Program = program
	mod.Deps = make(map[string]*Module)
	mod.Structs = make(map[string]*Struct)
	return
}

type Module struct {
	Name      string
	Path      string
	GoVersion string
	Requires  []Require
	Program   *loader.Program
	Deps      map[string]*Module
	Structs   map[string]*Struct
}

func (mod *Module) CreatedPackageInfos() (infos []*loader.PackageInfo) {
	infos = mod.Program.Created
	return
}

func (mod *Module) FileInfo(f *ast.File) (v *token.File) {
	v = mod.Program.Fset.File(f.Pos())
	return
}

func (mod *Module) GetPackageOfFile(f *ast.File) (name string, ident string, has bool) {
	for _, info := range mod.Program.Created {
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

func (mod *Module) FindDeps(pkgPath string) (deps map[string]*Module, err error) {
	deps = make(map[string]*Module)
	for s, program := range mod.Deps {
		if pkgPath == s || strings.Contains(pkgPath, s+"/") {
			deps[s] = program
		}
	}
	if len(deps) == 0 {
		for _, require := range mod.Requires {
			if pkgPath == require.Name || strings.Contains(pkgPath, require.Name+"/") {
				dep, depErr := NewModule(filepath.Join(require.Path(), "go.mod"))
				if depErr != nil {
					err = fmt.Errorf("load dep(%s) failed, %s, %v", require.Name, require.Path(), depErr)
					return
				}
				mod.Deps[require.Name] = dep
				deps[require.Name] = dep
			}
		}
	}
	if len(deps) == 0 {
		err = fmt.Errorf("no %s require in go.mod", pkgPath)
		return
	}
	return
}

func (mod *Module) GetStruct(pkg string, name string) (v *Struct, has bool) {
	key := fmt.Sprintf("%s:%s", pkg, name)
	v, has = mod.Structs[key]
	return
}

func (mod *Module) SetStruct(v *Struct) {
	mod.Structs[v.Key()] = v
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
