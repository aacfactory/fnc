package codegen

import (
	"fmt"
	"go/ast"
	"strings"
)

func LoadFn(project *Project) (fnFiles []FnFile, err error) {

	for _, info := range project.Program.Created {
		for _, file := range info.Files {
			fnFile, has, fileErr := loadFnFile(project, file)
			if fileErr != nil {
				err = fmt.Errorf("load fn file failed, %v", fileErr)
				return
			}
			if has {
				fnFiles = append(fnFiles, fnFile)
			}
		}
	}

	return
}

func loadFnFile(project *Project, f *ast.File) (fnFile FnFile, has bool, err error) {

	fileInfo := project.Program.Fset.File(f.Pos())
	if fileInfo == nil {
		err = fmt.Errorf("get %s file failed", f.Name.Name)
		return
	}
	filename := fileInfo.Name()
	// path
	fnFile.Path = filename
	// file doc
	fnFile.Doc = make([]string, 0, 1)
	if f.Doc != nil && f.Doc.List != nil {
		for _, line := range f.Doc.List {
			fnFile.Doc = append(fnFile.Doc, line.Text)
		}
	}
	// Package
	pkgName, hasPkgName := project.PackageNameOfFile(f)
	if !hasPkgName {
		err = fmt.Errorf("read %s failed, package name is not founed", filename)
		return
	}
	fnFile.Package = pkgName
	// imports
	if f.Imports != nil {
		imports := make([]Import, 0, 1)
		for _, spec := range f.Imports {
			path := strings.ReplaceAll(spec.Path.Value, "\"", "")
			alias := spec.Name
			name := path[strings.LastIndex(path, "/")+1:]
			if alias != nil && alias.Name != "" {
				name = alias.Name
			}
			import0 := Import{
				Path: path,
				Name: name,
			}
			imports = append(imports, import0)
		}
		fnFile.Imports = imports
	}

	// func
	fns := make([]Fn, 0, 1)

	if f.Decls != nil {
		for _, decl := range f.Decls {
			funcDecl, fnOk := decl.(*ast.FuncDecl)
			if !fnOk {
				continue
			}
			// not method
			if funcDecl.Recv != nil {
				continue
			}
			// not doc
			if funcDecl.Doc == nil {
				continue
			}
			fn := Fn{}
			// name
			fn.Name = funcDecl.Name.Name
			fn.Exported = ast.IsExported(fn.Name)
			// doc
			comments := funcDecl.Doc.List
			if comments == nil || len(comments) == 0 {
				continue
			}
			address := ""
			proxy := make([]string, 0, 1)
			docs := make([]string, 0, 1)
			openAPIContent := ""
			for _, comment := range comments {
				if comment == nil {
					continue
				}
				line := comment.Text
				docs = append(docs, line)
				if strings.Contains(line, "@Fn") {
					address = strings.TrimSpace(line[strings.Index(line, "@Fn")+4:])
					continue
				}
				if strings.Contains(line, "@Proxy") {
					proxy = strings.Split(strings.TrimSpace(line[strings.Index(line, "@Proxy")+7:]), " ")
					continue
				}
				if strings.Contains(line, "@OpenAPI") {
					openAPIContent = strings.TrimSpace(line[strings.Index(line, "@OpenAPI")+9:])
					continue
				}
			}
			if address == "" {
				continue
			}
			fn.Address = address
			fn.Proxy = proxy
			fn.OpenAPI = openAPIContent

			fn.Doc = docs
			fn.Imports = make(map[string]Import)
			// params
			p1, p2, paramsErr := parseFnParams(project, fnFile.Imports, fn.Imports, funcDecl.Type.Params)
			if paramsErr != nil {
				err = fmt.Errorf("%s:%s, %v", filename, fn.Name, paramsErr)
				return
			}
			fn.In = append(fn.In, p1, p2)

			// results

			// fin
			fns = append(fns, fn)
		}
	}

	if len(fns) > 0 {
		fnFile.Functions = fns
		has = true
	}

	return
}

type FnFile struct {
	Path        string `json:"path,omitempty"`
	Doc         []string `json:"doc,omitempty"`
	Package     string `json:"package,omitempty"`
	Imports     []Import `json:"imports,omitempty"`
	Functions   []Fn `json:"functions,omitempty"`
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
	Type Type `json:"type,omitempty"`
}

type Fn struct {
	Exported bool             `json:"exported,omitempty"`
	// key - import name
	Imports map[string]Import `json:"imports,omitempty"`
	Doc      []string `json:"doc,omitempty"`
	Address  string `json:"address,omitempty"`
	Proxy    []string `json:"proxy,omitempty"`
	OpenAPI  string `json:"openApi,omitempty"`
	Name     string `json:"name,omitempty"`
	In       []FuncItem `json:"in,omitempty"`
	Out      []FuncItem `json:"out,omitempty"`
}

type Type struct {
	IsBasic     bool `json:"isBasic,omitempty"`
	IsStruct    bool `json:"isStruct,omitempty"`
	IsInterface bool `json:"isInterface,omitempty"`
	IsPtr       bool `json:"isPtr,omitempty"`
	IsArray     bool `json:"isArray,omitempty"`
	IsMap       bool `json:"isMap,omitempty"`
	IsErr       bool `json:"isErr,omitempty"`
	Package     Import `json:"package,omitempty"`
	Name        string `json:"name,omitempty"`
	Struct      *Struct `json:"struct,omitempty"`
	InnerType   *Type `json:"innerType,omitempty"`
}

type Struct struct {
	Exported bool `json:"exported,omitempty"`
	Doc      []string `json:"doc,omitempty"`
	Name     string `json:"name,omitempty"`
	Fields   []Field `json:"fields,omitempty"`
}

type Field struct {
	Exported bool `json:"exported,omitempty"`
	Doc      []string `json:"doc,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     Type `json:"type,omitempty"`
	Tags     []FieldTag `json:"tags,omitempty"`
}

type FieldTag struct {
	Name   string `json:"name,omitempty"`
	Values []string `json:"values,omitempty"`
}
