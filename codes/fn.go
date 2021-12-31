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
	"github.com/aacfactory/cases"
)

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

func (f *Fn) NameToConstName() (v string, err error) {
	atoms, parseErr := cases.Snake().Parse(f.Name())
	if parseErr != nil {
		err = fmt.Errorf("fn name is invalid, %s is not snake format", f.Name())
		return
	}
	atoms = append(atoms, "fn")
	v = cases.Camel().Format(atoms)
	return
}
