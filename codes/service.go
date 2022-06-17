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

type Component struct {
	Name   string
	Loader string
	Struct string
}

type Service struct {
	DirPath     string
	Package     string
	Imports     Imports
	fns         map[string]*Fn
	Annotations map[string]string
	Components  map[string]*Component
}

func (svc *Service) Internal() (v bool) {
	x, has := svc.Annotations["internal"]
	if has {
		v = x == "true"
	}
	return
}

func (svc *Service) Name() (v string) {
	v = svc.Annotations["service"]
	return
}

func (svc *Service) Title() (v string) {
	v = svc.Annotations["title"]
	return
}

func (svc *Service) Description() (v string) {
	v = svc.Annotations["description"]
	return
}

func (svc *Service) AddComponent(c *Component) (err error) {
	name := c.Name
	_, exist := svc.Components[name]
	if exist {
		err = fmt.Errorf("fnc: %s component in %s is duplicated", name, svc.Name())
		return
	}
	svc.Components[name] = c
	return
}

func (svc *Service) AddFn(fn *Fn) (err error) {
	name := fn.Name()
	_, exist := svc.fns[name]
	if exist {
		err = fmt.Errorf("fnc: %s fn in %s is duplicated", name, svc.Name())
		return
	}
	svc.fns[name] = fn
	if fn.Param != nil && !fn.Param.InFile {
		v, has := fn.Param.Type.GetImport()
		if has {
			svc.addImport(v)
		}
	}
	if fn.Result != nil && !fn.Result.InFile {
		v, has := fn.Result.Type.GetImport()
		if has {
			svc.addImport(v)
		}
	}
	return
}

func (svc *Service) addImport(v *Import) {
	if v.Alias == "_" {
		return
	}
	if v.Name == svc.Package {
		return
	}
	added := false
	for _, import0 := range svc.Imports {
		if import0.Path == v.Path {
			added = true
			break
		}
	}
	if added {
		return
	}
	svc.Imports = append(svc.Imports, v)
}

func (svc *Service) Fns() (v []*Fn) {
	v = make([]*Fn, 0, 1)
	for _, fn := range svc.fns {
		v = append(v, fn)
	}
	sort.Slice(v, func(i, j int) bool {
		return v[i].Name() < v[j].Name()
	})
	return
}

func (svc *Service) generate() (err error) {
	genErr := svc.generateFile()
	if genErr != nil {
		err = genErr
		return
	}
	return
}
