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
)

type FnFile struct {
	Path      string   `json:"path,omitempty"`
	Doc       []string `json:"doc,omitempty"`
	Package   string   `json:"package,omitempty"`
	Imports   []Import `json:"imports,omitempty"`
	Functions []Fn     `json:"functions,omitempty"`
}

func FindImport(imports []Import, name string) (result Import, has bool) {
	for _, i := range imports {
		if i.Name == name {
			has = true
			result = i
			return
		}
	}
	return
}

type Import struct {
	Path string `json:"path,omitempty"`
	Name string `json:"name,omitempty"`
}

func (i *Import) Anonymous() bool {
	return i.Name == "_"
}

func (i *Import) Static() bool {
	return i.Name == "."
}

type FuncItem struct {
	Name string `json:"name,omitempty"`
	Type Type   `json:"type,omitempty"`
}

type Fn struct {
	Exported bool `json:"exported,omitempty"`
	// key - import name
	Imports map[string]Import `json:"imports,omitempty"`
	Doc     []string          `json:"doc,omitempty"`
	Address string            `json:"address,omitempty"`
	Proxy   []string          `json:"proxy,omitempty"`
	OpenAPI string            `json:"openApi,omitempty"`
	Name    string            `json:"name,omitempty"`
	In      []FuncItem        `json:"in,omitempty"`
	Out     []FuncItem        `json:"out,omitempty"`
}

func (fn Fn) String() string {
	return fmt.Sprintf("[%s][%s][in(%d)][out(%d)]", fn.Address, fn.Name, len(fn.In), len(fn.Out))
}
