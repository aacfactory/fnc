package codegen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/types"
	"golang.org/x/tools/go/loader"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func LoadProject(path string) (p *Project, err error) {

	mod, modErr := loadModuleFile(path)
	if modErr != nil {
		err = fmt.Errorf("fnc load project failed, %v", modErr)
		return
	}

	files := make(map[string][]string)

	filesErr := loadPackageFiles(path, mod.Name, files)
	if filesErr != nil {
		err = fmt.Errorf("fnc load project failed, %v", filesErr)
		return
	}

	if len(files) == 0 {
		err = fmt.Errorf("fnc load project failed, no go files loaded")
		return
	}

	config := loader.Config{
		Cwd:         path,
		ParserMode:  parser.ParseComments,
		AllowErrors: true,
		TypeChecker: types.Config{
			Error: func(err error) {

			},
		},
	}

	for pkg, filenames := range files {
		config.CreateFromFilenames(pkg, filenames...)
	}

	program, loadErr := config.Load()
	if loadErr != nil {
		err = fmt.Errorf("fnc load project failed, %v", loadErr)
		return
	}

	p = &Project{
		Path:    path,
		Module:  mod,
		Program: program,
	}

	return
}

type Project struct {
	Path    string          `json:"path,omitempty"`
	Module  Module          `json:"module,omitempty"`
	Program *loader.Program `json:"-"`
	Fns     []FnFile        `json:"fns,omitempty"`
	// Structs
	// key = package(path).StructName
	Structs map[string]Struct `json:"structs,omitempty"`
}

func (p *Project) FindStruct(pkgPath string, name string) (str Struct, has bool) {
	str, has = p.Structs[fmt.Sprintf("%s.%s", pkgPath, name)]
	return
}

func (p *Project) PutStruct(pkgPath string, name string, str Struct) {
	_, has := p.Structs[fmt.Sprintf("%s.%s", pkgPath, name)]
	if has {
		return
	}
	p.Structs[fmt.Sprintf("%s.%s", pkgPath, name)] = str
	return
}

func (p *Project) FindObject(pkgName string, name string) (obj types.Object, has bool) {
	pkg := p.Program.Package(pkgName)
	if pkg == nil {
		return
	}
	if pkg.Defs == nil {
		return
	}
	for ident, object := range pkg.Defs {
		if object == nil {
			continue
		}
		if ident.Name == name {
			if object.Type() == nil {
				continue
			}
			if object.Type().String() == fmt.Sprintf("%s.%s", pkgName, name) {
				has = true
				obj = object
				return
			}
		}
	}
	return
}

func (p *Project) TypeOf(expr ast.Expr) (typ types.Type, has bool) {
	if p.Program == nil {
		panic(fmt.Errorf("fnc get type of expr failed, program is not setup"))
		return
	}
	// created
	for _, info := range p.Program.Created {
		typ = info.TypeOf(expr)
		if typ != nil {
			has = true
			return
		}
	}
	// imported
	for _, info := range p.Program.Imported {
		typ = info.TypeOf(expr)
		if typ != nil {
			has = true
			return
		}
	}
	return
}

func (p *Project) ObjectOf(ident *ast.Ident) (obj types.Object, has bool) {
	if p.Program == nil {
		panic(fmt.Errorf("fnc get object of ident failed, program is not setup"))
		return
	}
	// created
	for _, info := range p.Program.Created {
		obj = info.ObjectOf(ident)
		if obj != nil {
			has = true
			return
		}
	}
	// imported
	for _, info := range p.Program.Imported {
		obj = info.ObjectOf(ident)
		if obj != nil {
			has = true
			return
		}
	}
	return
}

func (p *Project) FilepathOfFile(f *ast.File) (filePath string, has bool) {
	fileInfo := p.Program.Fset.File(f.Pos())
	if fileInfo == nil {
		return
	}
	filePath = fileInfo.Name()
	has = true
	return
}

func (p *Project) PackageNameOfFile(f *ast.File) (name string, has bool) {
	pkgPos := p.Program.Fset.Position(f.Package)
	path := pkgPos.Filename
	if path == "" {
		return
	}
	lineNo := pkgPos.Line
	if lineNo < 1 {
		return
	}

	column := pkgPos.Column
	if column < 0 {
		return
	}

	content, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return
	}
	buf := bufio.NewReader(bytes.NewReader(content))
	lines := 0
	for {
		lines++
		lineContent, _, lineErr := buf.ReadLine()
		if lineNo == lines {
			items := strings.Split(string(lineContent), " ")
			name = items[column]
			has = name != ""
			return
		}
		if lineErr != nil {
			break
		}
	}

	return
}

func loadFilenames(path0 string, filenames map[string]string) (err error) {

	walkErr := filepath.Walk(path0, func(path string, info fs.FileInfo, cause error) (err error) {
		if strings.Contains(path, "/.git/") {
			return
		}
		if info.IsDir() {
			if path == path0 || strings.HasSuffix(path, ".git") {
				return
			}
			subErr := loadFilenames(path, filenames)
			if subErr != nil {
				err = subErr
				return
			}
			return
		}
		if !strings.HasSuffix(info.Name(), ".go") {
			return
		}
		if _, has := filenames[path]; !has {
			filenames[path] = info.Name()
		}

		return
	})

	if walkErr != nil {
		err = fmt.Errorf("load filename from %s failed, %v", path0, walkErr)
		return
	}

	return
}
