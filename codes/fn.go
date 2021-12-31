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
	"sort"
)

type Service struct {
	DirPath     string
	Package     string
	Imports     Imports
	fns         map[string]*Fn
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
	_, exist := n.fns[name]
	if exist {
		err = fmt.Errorf("fnc: %s fn in %s is duplicated", name, n.Name())
		return
	}
	n.fns[name] = fn
	if fn.Param != nil && !fn.Param.InFile {
		v, has := fn.Param.Type.GetImport()
		if has {
			n.addImport(v)
		}
	}
	if fn.Result != nil && !fn.Result.InFile {
		v, has := fn.Result.Type.GetImport()
		if has {
			n.addImport(v)
		}
	}
	return
}

func (n *Service) addImport(v *Import) {
	if v.Alias == "_" {
		return
	}
	if v.Name == n.Package {
		return
	}
	added := false
	for _, import0 := range n.Imports {
		if import0.Path == v.Path {
			added = true
			break
		}
	}
	if added {
		return
	}
	n.Imports = append(n.Imports, v)
}

func (n *Service) Fns() (v []*Fn) {
	v = make([]*Fn, 0, 1)
	for _, fn := range n.fns {
		v = append(v, fn)
	}
	sort.Slice(v, func(i, j int) bool {
		return v[i].Name() < v[j].Name()
	})
	return
}

func (n *Service) generate() (err error) {
	// todo
	return
}

type FnField struct {
	InFile bool
	Name   string
	Type   *Type
}

func (x *FnField) String() (s string) {
	s = fmt.Sprintf("%s %s", x.Name, x.Title())
	return
}

func (x *FnField) Title() (title string) {
	if x.Type.Struct != nil {
		v, has := x.Type.Annotations()["title"]
		if has {
			title = v
		} else {
			title = x.Type.Struct.Key()
		}
	} else {
		title = x.Name
	}
	return
}

func (x *FnField) Description() (description string) {
	if x.Type.Struct != nil {
		v, has := x.Type.Annotations()["description"]
		if has {
			description = v
			return
		}
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
