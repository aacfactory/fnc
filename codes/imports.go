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
	"go/ast"
	"strings"
)

func NewImports(file *ast.File) (imports Imports) {
	imports = make([]*Import, 0, 1)
	if file.Imports == nil || len(file.Imports) == 0 {
		return
	}
	for _, spec := range file.Imports {
		path := strings.ReplaceAll(spec.Path.Value, "\"", "")
		aliasSpec := spec.Name
		alias := ""
		name := path[strings.LastIndex(path, "/")+1:]
		if aliasSpec != nil && aliasSpec.Name != "" {
			alias = aliasSpec.Name
			name = alias
		}
		import0 := Import{
			Alias: alias,
			Name:  name,
			Path:  path,
		}
		imports = append(imports, &import0)
	}
	return
}

type Imports []*Import

func (s Imports) FindByName(name string) (v *Import, has bool) {
	for _, i := range s {
		if i.Name == name {
			v = i
			has = true
			return
		}
	}
	return
}

type Import struct {
	Name  string
	Alias string
	Path  string
}
