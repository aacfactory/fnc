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
	"github.com/aacfactory/gcg"
	"strings"
)

type FnField struct {
	InFile bool
	Name   string
	Type   *Type
}

func (x *FnField) key() (v string) {
	v = x.Type.ObjectKey()
	return
}

func (x *FnField) generateObjectDocument() (code gcg.Code) {
	code = x.Type.generateObject()
	return
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

func (f *Fn) HasValidate() (v bool) {
	v = f.Annotations["validate"] == "true"
	return
}

func (f *Fn) HasAuthorization() (v bool) {
	v = f.Annotations["authorization"] == "true"
	return
}

func (f *Fn) HasDeprecated() (v bool) {
	v = f.Annotations["deprecated"] == "true"
	return
}

func (f *Fn) HasPermission() (v bool) {
	v = f.Annotations["permission"] == "true"
	return
}

func (f *Fn) IsInternal() (v bool) {
	v = f.Annotations["internal"] == "true"
	return
}

func (f *Fn) HasTx() (kind string, opts []string, has bool) {
	kind = f.Annotations["transactional"]
	has = kind != ""
	if has {
		optIdx := strings.Index(kind, ":")
		if optIdx > 0 {
			// @transactional sql:1s
			kind = kind[0:optIdx]
			opt := kind[optIdx+1:]
			opts = strings.Split(opt, ",")
		}
	}
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
	v = "_" + cases.LowerCamel().Format(atoms)
	return
}

func (f *Fn) NameToProxyName() (v string, err error) {
	atoms, parseErr := cases.Camel().Parse(f.FuncName)
	if parseErr != nil {
		err = fmt.Errorf("fn name is invalid, %s is not camel format", f.Name())
		return
	}
	v = cases.Camel().Format(atoms)
	return
}
