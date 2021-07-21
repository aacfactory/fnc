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

package goparser

type Project struct {
	Path    string
	Module  Module
	Package Package
}

func (p *Project) Filter(fn func(file File) bool) (files []File) {
	files = p.filter(p.Package, fn)
	return
}

func (p *Project) filter(pkg Package, fn func(file File) bool) (files []File) {
	files = make([]File, 0, 1)
	if pkg.Files != nil {
		for _, file := range pkg.Files {
			if fn(file) {
				files = append(files, file)
			}
		}
	}
	if pkg.SubPackages != nil {
		for _, subPackage := range pkg.SubPackages {
			subFiles := p.filter(subPackage, fn)
			if subFiles != nil {
				files = append(files, subFiles...)
			}
		}
	}
	return
}

type Module struct {
	Name     string
	Version  string
	Requires []Require
}

type Require struct {
	Name    string
	Version string
}

type Package struct {
	Path        string
	Name        string
	Doc         string
	Files       []File
	SubPackages []Package
}

type File struct {
	Path       string
	Comment    string
	Package    string
	Imports    []Import
	Constants  []Const
	Variables  []Var
	Functions  []Func
	Interfaces []Interface
	Structs    []Struct
	AliasTypes []AliasType
}

type Import struct {
	Name      string
	Alias     string
	Anonymous bool
	Static    bool
}

type Const struct {
	Exported bool
	Comment  string
	Name     string
	Type     string
	Value    string
}

type Var struct {
	Exported bool
	Comment  string
	Name     string
	Type     Type
	Value    string
}

type AliasType struct {
	Exported bool
	Name     string
	Type     Type
}

type FuncReceiver struct {
	Name string
	Type Struct
}

type FuncItem struct {
	Name string
	Type Type
}

type Func struct {
	Exported bool
	Comment  string
	Name     string
	Receiver *FuncReceiver
	In       []FuncItem
	Out      []FuncItem
}

type Field struct {
	Exported bool
	Comment  string
	Name     string
	Type     Type
	Tags     []FieldTag
}

type FieldTag struct {
	Name   string
	Values []string
}

type Interface struct {
	Exported  bool
	Comment   string
	Name      string
	Functions []Func
}

type Struct struct {
	Exported bool
	Comment  string
	Name     string
	Fields   []Field
}

type Type struct {
	Kind       string // indent(int string ...), map, array, chan, struct, func
	Package    string
	Name       string
	Ptr        bool
	InnerTypes []Type
}
