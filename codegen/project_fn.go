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
	"strings"
)

func (p *Project) loadFn() (err error) {
	for _, info := range p.Program.Created {
		for _, file := range info.Files {
			fnFile, has, fileErr := p.loadFnFile(file)
			if fileErr != nil {
				err = fmt.Errorf("load fn file failed, %v", fileErr)
				return
			}
			if has {
				p.Fns = append(p.Fns, fnFile)
			}
		}
	}

	return
}

func (p *Project) loadFnFile(f *ast.File) (fnFile FnFile, has bool, err error) {

	fileInfo := p.Program.Fset.File(f.Pos())
	if fileInfo == nil {
		err = fmt.Errorf("get %s file failed", f.Name.Name)
		return
	}
	filename := fileInfo.Name()
	// path
	fnFile.Path = filename
	// file doc
	fnFile.Doc = parseDoc(f.Doc.Text())

	// Package
	pkgPath, _, hasPkg := p.PackageOfFile(f)
	if !hasPkg {
		err = fmt.Errorf("read %s failed, package name is not founed", filename)
		return
	}
	fnFile.Package = pkgPath
	// imports
	fnFile.Imports = p.GetImports(f)

	// func
	fns := make([]Fn, 0, 1)

	if f.Decls != nil {
		for _, decl := range f.Decls {
			funcDecl, fnOk := decl.(*ast.FuncDecl)
			if !fnOk {
				continue
			}
			// not method
			if funcDecl.Recv != nil {
				continue
			}
			// not doc
			if funcDecl.Doc == nil {
				continue
			}
			fn := Fn{}
			// name
			fn.Name = funcDecl.Name.Name
			fn.Exported = ast.IsExported(fn.Name)
			// doc
			fn.Doc = parseDoc(funcDecl.Doc.Text())
			address := ""
			proxy := make([]string, 0, 1)
			openAPIContent := ""
			for _, line := range fn.Doc {
				if line == "" {
					continue
				}
				if strings.Contains(line, "@Fn") {
					address = strings.TrimSpace(line[strings.Index(line, "@Fn")+4:])
					continue
				}
				if strings.Contains(line, "@Proxy") {
					proxy = strings.Split(strings.TrimSpace(line[strings.Index(line, "@Proxy")+7:]), " ")
					continue
				}
				if strings.Contains(line, "@OpenAPI") {
					openAPIContent = strings.TrimSpace(line[strings.Index(line, "@OpenAPI")+9:])
					continue
				}
			}
			if address == "" {
				continue
			}
			fn.Address = address
			fn.Proxy = proxy
			fn.OpenAPI = openAPIContent

			fn.Imports = make(map[string]Import)
			// params
			params, paramsErr := p.parseFnParams(fnFile.Package, fnFile.Imports, fn.Imports, funcDecl.Type.Params)
			if paramsErr != nil {
				err = fmt.Errorf("%s:%s, %v", filename, fn.Name, paramsErr)
				return
			}
			fn.In = params

			// results
			results, resultsErr := p.parseFnResults(fnFile.Package, fnFile.Imports, fn.Imports, funcDecl.Type.Results)
			if resultsErr != nil {
				err = fmt.Errorf("%s:%s, %v", filename, fn.Name, resultsErr)
				return
			}
			fn.Out = results

			// fin
			fns = append(fns, fn)
		}
	}

	if len(fns) > 0 {
		fnFile.Functions = fns
		has = true
	}

	return
}
