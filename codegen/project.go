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
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/loader"
	"os"
	"path/filepath"
	"strings"
)

func LoadProject(path string) (p *Project, err error) {

	mod, modErr := loadModuleFile(path)
	if modErr != nil {
		err = fmt.Errorf("fnc load project failed, %v", modErr)
		return
	}

	files := make(map[string][]string)

	filesErr := loadPackageFiles(path, mod.Name, files)
	if filesErr != nil {
		err = fmt.Errorf("fnc load project failed, %v", filesErr)
		return
	}

	if len(files) == 0 {
		err = fmt.Errorf("fnc load project failed, no go files loaded")
		return
	}

	config := loader.Config{
		Cwd:         path,
		ParserMode:  parser.ParseComments,
		AllowErrors: true,
		TypeChecker: types.Config{
			Error: func(err error) {

			},
		},
	}

	for pkg, filenames := range files {
		config.CreateFromFilenames(pkg, filenames...)
	}

	importPrograms := make(map[string]*loader.Program)
	gopath, hasGopath := os.LookupEnv("GOPATH")
	if hasGopath {
		modPrefix := filepath.Join(gopath, "pkg", "mod")
		for _, require := range mod.Requires {
			importCWD := filepath.Join(modPrefix, fmt.Sprintf("%s@%s", require.Name, require.Version))
			if require.Replace != "" {
				replace := require.Replace
				if strings.Index(replace, ".") == 0 {
					replace = filepath.Join(path, require.Replace)
				}
				importCWD = replace
			}
			configOfImport := loader.Config{
				Cwd:         importCWD,
				ParserMode:  parser.ParseComments,
				AllowErrors: true,
				TypeChecker: types.Config{
					Error: func(err error) {

					},
				},
			}
			importModFiles := make(map[string][]string)

			importModFilesErr := loadPackageFiles(importCWD, require.Name, importModFiles)
			if importModFilesErr != nil {
				err = fmt.Errorf("fnc load project failed, %v", importModFilesErr)
				return
			}

			if len(importModFiles) == 0 {
				err = fmt.Errorf("fnc load project failed, no go files loaded")
				return
			}
			for pkg, filenames := range importModFiles {
				configOfImport.CreateFromFilenames(pkg, filenames...)
			}
			programOfImport, loadImportErr := configOfImport.Load()
			if loadImportErr != nil {
				err = fmt.Errorf("fnc load import mod project failed, %v", loadImportErr)
				return
			}
			importPrograms[require.Name] = programOfImport
		}
	} else {
		Log().Warnf("fnc can not get GOOATH, so it will not parse imported in go.mod")
	}

	program, loadErr := config.Load()
	if loadErr != nil {
		err = fmt.Errorf("fnc load project failed, %v", loadErr)
		return
	}

	p = &Project{
		Path:            path,
		Module:          mod,
		Program:         program,
		ImportPrograms:  importPrograms,
		Fns:             make([]FnFile, 0, 1),
		Structs:         make(map[string]Struct),
		structFindStack: NewStructFindStack(),
	}

	return
}

type Project struct {
	Path           string                     `json:"path,omitempty"`
	Module         Module                     `json:"module,omitempty"`
	Program        *loader.Program            `json:"-"`
	ImportPrograms map[string]*loader.Program `json:"-"`
	Fns            []FnFile                   `json:"fns,omitempty"`
	// Structs
	// key = package(path).StructName
	Structs         map[string]Struct `json:"structs,omitempty"`
	structFindStack *StructFindStack
}

func (p *Project) FindObject(pkgName string, name string) (obj types.Object, has bool) {
	pkg := p.Program.Package(pkgName)
	if pkg == nil {
		return
	}
	if pkg.Defs == nil {
		return
	}
	for ident, object := range pkg.Defs {
		if object == nil {
			continue
		}
		if ident.Name == name {
			if object.Type() == nil {
				continue
			}
			if object.Type().String() == fmt.Sprintf("%s.%s", pkgName, name) {
				has = true
				obj = object
				return
			}
		}
	}
	return
}

func (p *Project) TypeOf(expr ast.Expr) (typ types.Type, has bool) {
	if p.Program == nil {
		panic(fmt.Errorf("fnc get type of expr failed, program is not setup"))
		return
	}
	// created
	for _, info := range p.Program.Created {
		typ = info.TypeOf(expr)
		if typ != nil {
			has = true
			return
		}
	}
	// imported
	for _, info := range p.Program.Imported {
		typ = info.TypeOf(expr)
		if typ != nil {
			has = true
			return
		}
	}
	return
}

func (p *Project) ObjectOf(ident *ast.Ident) (obj types.Object, has bool) {
	if p.Program == nil {
		panic(fmt.Errorf("fnc get object of ident failed, program is not setup"))
		return
	}
	// created
	for _, info := range p.Program.Created {
		obj = info.ObjectOf(ident)
		if obj != nil {
			has = true
			return
		}
	}
	// imported
	for _, info := range p.Program.Imported {
		obj = info.ObjectOf(ident)
		if obj != nil {
			has = true
			return
		}
	}
	return
}

func (p *Project) FilepathOfFile(f *ast.File) (filePath string, has bool) {
	var fileInfo *token.File
	fileInfo = p.Program.Fset.File(f.Pos())
	if fileInfo == nil {
		for _, program := range p.ImportPrograms {
			fileInfo = program.Fset.File(f.Pos())
			if fileInfo != nil {
				break
			}
		}
	}
	if fileInfo == nil {
		return
	}
	filePath = fileInfo.Name()
	has = true
	return
}

func (p *Project) PackageOfFile(f *ast.File) (pkgPath string, pkgName string, has bool) {
	for _, info := range p.Program.AllPackages {
		for _, file := range info.Files {
			if file == f {
				pkgPath = info.Pkg.Path()
				pkgName = info.Pkg.Name()
				has = true
				return
			}
		}
	}
	for _, program := range p.ImportPrograms {
		for _, info := range program.AllPackages {
			for _, file := range info.Files {
				if file == f {
					pkgPath = info.Pkg.Path()
					pkgName = info.Pkg.Name()
					has = true
					return
				}
			}
		}
	}
	return
}
