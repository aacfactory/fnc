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
	Type     string
	Value    string
}

type Func struct {
	Exported bool
	Comment  string
	Name     string
	Receiver string
	In       []Field
	Out      []Field
}

type Field struct {
	Exported bool
	Comment  string
	Name     string
	Type     string
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
	Exported  bool
	Comment   string
	Name      string
	Fields    []Field
	Functions []Func
}
