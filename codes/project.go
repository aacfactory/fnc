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
	"go/ast"
	"path/filepath"
)

const (
	useAuthPackage       = ""
	usePermissionPackage = ""
	useTxSQLPackage      = ""
	useCachePackage      = ""
)

// golang.org/x/tools/imports
// go/format
func NewProject(projectDirPath string) (p *Project, err error) {
	projectDirPath = filepath.ToSlash(projectDirPath)
	if !filepath.IsAbs(projectDirPath) {
		absFilePath, absErr := filepath.Abs(projectDirPath)
		if absErr != nil {
			err = fmt.Errorf("fnc: new project failed for absolute representation of %s, %v", projectDirPath, absErr)
			return
		}
		projectDirPath = absFilePath
	}
	mod, modErr := NewModule(filepath.Join(projectDirPath, "go.mod"))
	if modErr != nil {
		err = modErr
		return
	}
	p = &Project{
		dir: projectDirPath,
		mod: mod,
		Fns: make(map[string]*Service),
	}
	err = p.scan()
	return
}

type Project struct {
	dir string
	mod *Module
	Fns map[string]*Service
}

func (p *Project) Path() (v string) {
	v = p.mod.Path
	return
}

func (p *Project) scan() (err error) {
	for _, info := range p.mod.CreatedPackageInfos() {
		fns := make([]*Fn, 0, 1)
		serviceName := ""
		serviceAnnotations := make(map[string]string)
		for _, file := range info.Files {
			if serviceName == "" {
				serviceName, serviceAnnotations = p.scanServiceDoc(file)
				if serviceName != "" {
					_, existService := p.Fns[serviceName]
					if existService {
						err = fmt.Errorf("fnc: scan fn failed for %s is duplicated", serviceName)
					}
					continue
				}
			}
			fnsOfFile, scanFnErr := p.scanFile(file)
			if scanFnErr != nil {
				err = scanFnErr
				return
			}
			if fnsOfFile != nil && len(fnsOfFile) > 0 {
				fns = append(fns, fnsOfFile...)
			}
		}
		if serviceName == "" {
			if len(fns) > 0 {
				err = fmt.Errorf("fnc: scan fn failed for some fns has no named service")
				return
			}
			continue
		}
		if len(fns) == 0 {
			continue
		}
		service := &Service{
			DirPath:     filepath.Join(p.dir, info.Pkg.Path()),
			Package:     info.Pkg.Name(),
			Imports:     make([]Import, 0, 1),
			Fns:         make(map[string]*Fn),
			Annotations: serviceAnnotations,
		}
		for _, fn := range fns {
			addErr := service.AddFn(fn)
			if addErr != nil {
				err = addErr
				return
			}
		}
		p.Fns[serviceName] = service
	}
	return
}

func (p *Project) scanServiceDoc(file *ast.File) (serviceName string, annotations map[string]string) {
	doc := file.Doc.Text()
	if doc == "" {
		return
	}
	annotations = getAnnotations(doc)
	if len(annotations) == 0 {
		return
	}
	name, hasName := annotations["service"]
	if !hasName {
		return
	}
	serviceName = name
	return
}

func (p *Project) scanFile(file *ast.File) (fns []*Fn, err error) {
	if file.Decls == nil {
		return
	}
	fns = make([]*Fn, 0, 1)
	for _, decl := range file.Decls {
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
		ok, fn, fnErr := p.scanFn(file, funcDecl)
		if fnErr != nil {
			err = fnErr
			return
		}
		if ok {
			fns = append(fns, fn)
		}
	}
	return
}

func (p *Project) scanFn(file *ast.File, decl *ast.FuncDecl) (ok bool, fn *Fn, err error) {
	if file.Decls == nil {
		return
	}
	doc := file.Doc.Text()
	annotations := getAnnotations(doc)
	if len(annotations) == 0 {
		return
	}
	name, hasName := annotations["fn"]
	if !hasName {
		return
	}
	if decl.Recv != nil {
		err = fmt.Errorf("fnc: scan %s fn failed for fn can has recv", name)
		return
	}
	// imports

	// check params and results
	params := decl.Type.Params
	if params == nil || (len(params.List) > 0 && len(params.List) < 3) {
		err = fmt.Errorf("fnc: scan %s fn failed for fn must has one or two params, first must be *fns.Context, seconed can be star struct", name)
		return
	}
	ctxTypeExpr, ctxTypeExprOk := params.List[0].Type.(*ast.SelectorExpr)

	results := decl.Type.Results
	if results == nil || (len(results.List) > 0 && len(results.List) < 3) {
		err = fmt.Errorf("fnc: scan %s fn failed for fn must has one or two results, first can be star struct, last must be github.com/aacfactory/errors.CodeError", name)
		return
	}
	// return
	fn = &Fn{
		FuncName:    decl.Name.Name,
		Param:       nil,
		Result:      nil,
		Annotations: annotations,
	}
	ok = true
	return
}
