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

type Namespace struct {
	DirPath string
	Package string
	Imports []Import
	Fns     map[string]*Fn
}

func (n *Namespace) AddImport(v Import) {
	if v.Alias == "_" {
		return
	}
	if v.Name == n.Package {
		return
	}
	added := false
	for _, import0 := range n.Imports {
		if import0.Name == v.Name {
			added = true
			break
		}
	}
	if added {
		return
	}
	n.Imports = append(n.Imports, v)
}

type Import struct {
	Name  string
	Alias string
	Ident string
}

type FnField struct {
	Name string
	Kind string
	Type *Struct
}

type Struct struct {
	Package      string
	PackageAlias string
	Name         string
	Fields       []*Field
	Annotations  map[string]string
}

func (s Struct) IsTime() (ok bool) {
	if s.Package == "time" && s.Name == "Time" {
		ok = true
	}
	return
}

func (s Struct) IsJson() (ok bool) {
	if s.Name == "RawMessage" {
		ok = true
	}
	if s.Package == "github.com/aacfactory/json" && s.Name == "Object" {
		ok = true
	}
	return
}

func (s Struct) IsJsonArray() (ok bool) {
	if s.Package == "github.com/aacfactory/json" && s.Name == "Array" {
		ok = true
	}
	return
}

type Field struct {
	Exported    bool
	Name        string
	Tag         string
	Kind        string
	Type        *Struct
	Annotations map[string]string
}

type Fn struct {
	Name        string
	FuncName    string
	Param       *FnField
	Result      *FnField
	Annotations map[string]string
}
