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
	"go/ast"
	"strings"
)

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

func (p *Project) GetImports(file *ast.File) (imports []Import) {

	if file.Imports == nil || len(file.Imports) == 0 {
		return
	}

	for _, spec := range file.Imports {
		path := strings.ReplaceAll(spec.Path.Value, "\"", "")
		alias := spec.Name
		name := path[strings.LastIndex(path, "/")+1:]
		if alias != nil && alias.Name != "" {
			name = alias.Name
		}
		import0 := Import{
			Path: path,
			Name: name,
		}
		imports = append(imports, import0)
	}

	return
}
