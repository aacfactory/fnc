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
	"bytes"
	"fmt"
	"github.com/aacfactory/logs"
	"go/ast"
	"path/filepath"
	"sort"
	"strings"
)

func NewProject(projectDirPath string, debug bool) (p *Project, err error) {
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
	lvl := logs.ErrorLevel
	if debug {
		lvl = logs.DebugLevel
	}
	log, logErr := logs.New(logs.Name("fnc"), logs.Color(true), logs.WithLevel(lvl))
	if logErr != nil {
		err = fmt.Errorf("fnc: new project failed for create log failed, %v", logErr)
		return
	}
	p = &Project{
		log:      log,
		dir:      projectDirPath,
		mod:      mod,
		services: make(map[string]*Service),
	}

	return
}

type Project struct {
	log      logs.Logger
	dir      string
	mod      *Module
	services map[string]*Service
}

func (p *Project) Path() (v string) {
	v = p.mod.Path
	return
}

func (p *Project) Scan() (err error) {
	for _, info := range p.mod.CreatedPackageInfos() {
		fns := make([]*Fn, 0, 1)
		serviceName := ""
		serviceAnnotations := make(map[string]string)
		for _, file := range info.Files {
			if serviceName == "" {
				serviceName, serviceAnnotations = p.scanServiceDoc(file)
				if serviceName != "" {
					_, existService := p.services[serviceName]
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
			//if len(fns) > 0 {
			//	err = fmt.Errorf("fnc: scan fn failed for some fns has no named service")
			//	return
			//}
			continue
		}
		if len(fns) == 0 {
			continue
		}
		servicePath := info.Pkg.Path()
		if strings.Index(servicePath, "main/") == 0 {
			servicePath = strings.Replace(servicePath, "main/", "", 1)
		}
		service := &Service{
			DirPath:     filepath.Join(p.dir, servicePath),
			Package:     info.Pkg.Name(),
			Imports:     make([]*Import, 0, 1),
			fns:         make(map[string]*Fn),
			Annotations: serviceAnnotations,
		}
		for _, fn := range fns {
			addErr := service.AddFn(fn)
			if addErr != nil {
				err = addErr
				return
			}
		}
		p.services[serviceName] = service
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
	// filePkg
	pkgPath, pkgName, hasPkg := p.mod.GetPackageOfFile(file)
	if !hasPkg {
		err = fmt.Errorf("fnc: scan %s fn failed for get packge of file", file.Name)
		return
	}
	pkg := &Import{
		Name:  pkgName,
		Alias: "",
		Path:  pkgPath,
	}
	// imports
	imports := NewImports(file)

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
		ok, fn, fnErr := p.scanFn(pkg, imports, funcDecl)
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

func (p *Project) scanFn(filePkg *Import, imports Imports, decl *ast.FuncDecl) (ok bool, fn *Fn, err error) {
	// doc
	doc := decl.Doc.Text()
	annotations := getAnnotations(doc)
	if len(annotations) == 0 {
		return
	}
	name, hasName := annotations["fn"]
	if !hasName {
		return
	}
	defer func() {
		if err != nil {
			if p.log.ErrorEnabled() {
				p.log.Error().Message(fmt.Sprintf("fnc: scan %s failed", name))
			}
		} else {
			if p.log.DebugEnabled() {
				p.log.Debug().Message(fmt.Sprintf("fnc: scan %s succeed", name))
			}
		}
	}()
	// check params
	params := decl.Type.Params
	if params == nil || !(len(params.List) > 0 && len(params.List) < 3) {
		err = fmt.Errorf("fnc: scan %s fn failed for fn must has one or two params, first must be fns.Context, seconed can be value object struct", name)
		return
	}
	ctxType, ctxTypeErr := NewType(params.List[0].Type, filePkg.Path, imports, p.mod)
	if ctxTypeErr != nil || !ctxType.IsFnsContext() {
		err = fmt.Errorf("fnc: scan %s fn failed for fn must has one or two params, first must be fns.Context, seconed can be value object struct", name)
		return
	}
	var param *FnField
	if len(params.List) == 2 {
		paramName := params.List[1].Names[0].Name
		paramType, paramTypeErr := NewType(params.List[1].Type, filePkg.Path, imports, p.mod)
		if paramTypeErr != nil {
			err = fmt.Errorf("fnc: scan %s fn failed for fn must has one or two params, first must be fns.Context, seconed can be value object struct, %v", name, paramTypeErr)
			return
		}
		if !(paramType.IsArray() || paramType.IsStruct() || paramType.IsStar() || paramType.IsBuiltin()) {
			err = fmt.Errorf("fnc: scan %s fn failed for fn must has one or two params, first must be fns.Context, seconed can be value object struct", name)
			return
		}
		_, hasImport := paramType.GetImport()
		param = &FnField{
			InFile: !hasImport,
			Name:   paramName,
			Type:   paramType,
		}
	}

	// check results
	results := decl.Type.Results
	if results == nil || !(len(results.List) > 0 && len(results.List) < 3) {
		err = fmt.Errorf("fnc: scan %s fn failed for fn must has one or two results, first can be star struct or star struct array, last must be github.com/aacfactory/errors.CodeError", name)
		return
	}

	resultName := ""
	var resultTypeExpr ast.Expr
	var errTypeExpr ast.Expr
	if len(results.List) == 1 {
		errTypeExpr = results.List[0].Type
	}
	if len(results.List) == 2 {
		resultName = results.List[0].Names[0].Name
		resultTypeExpr = results.List[0].Type
		errTypeExpr = results.List[1].Type
	}
	errType, errTypeErr := NewType(errTypeExpr, filePkg.Path, imports, p.mod)
	if errTypeErr != nil || !errType.IsFnsCodeError() {
		err = fmt.Errorf("fnc: scan %s fn failed for fn must has one or two results, first can be star struct or star struct array, last must be github.com/aacfactory/errors.CodeError", name)
		return
	}
	var result *FnField
	if resultName != "" {
		resultType, resultTypeErr := NewType(resultTypeExpr, filePkg.Path, imports, p.mod)
		if resultTypeErr != nil {
			err = fmt.Errorf("fnc: scan %s fn failed for fn must has one or two results, first can be star struct or star struct array, last must be github.com/aacfactory/errors.CodeError, %v", name, resultTypeErr)
			return
		}
		if resultType == nil {
			err = fmt.Errorf("fnc: scan %s fn failed for fn must has one or two results, first can be star struct or star struct array, last must be github.com/aacfactory/errors.CodeError", name)
			return
		}
		if !(resultType.IsArray() || resultType.IsStar() || resultType.IsBuiltin() || resultType.IsMap()) {
			err = fmt.Errorf("fnc: scan %s fn failed for fn must has one or two results, first can be star struct or star struct array, last must be github.com/aacfactory/errors.CodeError", name)
			return
		}
		_, hasImport := resultType.GetImport()
		result = &FnField{
			InFile: !hasImport,
			Name:   resultName,
			Type:   resultType,
		}
	}

	// return
	fn = &Fn{
		FuncName:    decl.Name.Name,
		Param:       param,
		Result:      result,
		Annotations: annotations,
	}
	ok = true
	return
}

func (p *Project) Services() (v []*Service) {
	v = make([]*Service, 0, 1)
	for _, service := range p.services {
		v = append(v, service)
	}
	sort.Slice(v, func(i, j int) bool {
		return v[i].Name() < v[j].Name()
	})
	return
}

func (p *Project) Generate() (err error) {
	services := p.Services()
	if len(services) == 0 {
		return
	}
	for _, service := range services {
		err = service.generate()
		if err != nil {
			return
		}
	}
	return
}

func (p *Project) String() (s string) {
	b := bytes.NewBufferString("")
	b.WriteString(fmt.Sprintf("\nMod        : %s\n", p.mod.Name))
	b.WriteString(fmt.Sprintf("Service Num: %d\n", len(p.services)))
	for _, service := range p.Services() {
		b.WriteString(fmt.Sprintf("*************%s*************\n", service.Name()))
		b.WriteString(fmt.Sprintf("Package : %s\n", service.Package))
		b.WriteString(fmt.Sprintf("Title   : %s\n", service.Title()))
		b.WriteString(fmt.Sprintf("Internal: %v\n", service.Internal()))
		b.WriteString(fmt.Sprintf("Path    : %s\n", service.DirPath))

		b.WriteString(fmt.Sprintf("Fn Num  : %d\n", len(service.fns)))
		for _, fn := range service.Fns() {
			b.WriteString(fmt.Sprintf("\tFn: %s\n", fn.Name()))
			b.WriteString(fmt.Sprintf("\t\tTitle : %s\n", fn.Title()))
			b.WriteString(fmt.Sprintf("\t\tParam : %v\n", fn.Param))
			b.WriteString(fmt.Sprintf("\t\tResult: %v\n", fn.Result))
		}
	}
	s = b.String()
	return
}
