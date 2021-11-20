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

import "fmt"

type Service struct {
	DirPath     string
	Package     string
	Imports     []Import
	Fns         map[string]*Fn
	Annotations map[string]string
}

func (n *Service) Internal() (v bool) {
	x, has := n.Annotations["internal"]
	if has {
		v = x == "true"
	}
	return
}

func (n *Service) Name() (v string) {
	v = n.Annotations["service"]
	return
}

func (n *Service) Title() (v string) {
	v = n.Annotations["title"]
	return
}

func (n *Service) Description() (v string) {
	v = n.Annotations["description"]
	return
}

func (n *Service) AddFn(fn *Fn) (err error) {
	name := fn.Name()
	_, exist := n.Fns[name]
	if exist {
		err = fmt.Errorf("fnc: %s fn in %s is duplicated", name, n.Name())
		return
	}
	n.Fns[name] = fn
	if fn.Param != nil && fn.Param.Import != nil {
		n.addImport(*fn.Param.Import)
	}
	if fn.Result != nil && fn.Result.Import != nil {
		n.addImport(*fn.Result.Import)
	}
	return
}

func (n *Service) addImport(v Import) {
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
	Name    string
	IsArray bool
	Star    bool
	Import  *Import
	Struct  *Struct
}

func (x *FnField) Title() (title string) {
	v, has := x.Struct.Annotations["title"]
	if has {
		title = v
		return
	}
	title = x.Struct.Key()
	return
}

func (x *FnField) Description() (description string) {
	v, has := x.Struct.Annotations["description"]
	if has {
		description = v
		return
	}
	return
}

type Struct struct {
	Package      string
	PackageAlias string
	Name         string
	Fields       []*Field
	Annotations  map[string]string
}

func (s Struct) Key() (key string) {
	key = s.Package + "." + s.Name
	return
}

func (s Struct) IsTime() (ok bool) {
	if s.Package == "time" && s.Name == "Time" {
		ok = true
	}
	return
}

func (s Struct) IsDate() (ok bool) {
	if s.Package == "github.com/aacfactory/json" && s.Name == "Date" {
		ok = true
	}
	return
}

func (s Struct) IsJson() (ok bool) {
	if s.Name == "RawMessage" {
		ok = true
	}
	return
}

func (s Struct) IsJsonObject() (ok bool) {
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

type Fn struct {
	FuncName    string
	Param       *FnField
	Result      *FnField
	Annotations map[string]string
}

func (f *Fn) Name() (v string) {
	v = f.Annotations["fn"]
	return
}

func (f *Fn) Title() (v string) {
	v = f.Annotations["title"]
	return
}

func (f *Fn) Description() (v string) {
	v = f.Annotations["description"]
	return
}

func (f *Fn) HasParam() (v bool) {
	v = f.Param != nil
	return
}

func (f *Fn) HasResult() (v bool) {
	v = f.Result != nil
	return
}
