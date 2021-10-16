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
		Fns: make(map[string]*Namespace),
	}
	err = p.scan()
	return
}

type Project struct {
	dir string
	mod *Module
	Fns map[string]*Namespace
}

func (p *Project) Path() (v string) {
	v = p.mod.Path
	return
}

func (p *Project) scan() (err error) {
	for _, info := range p.mod.CreatedPackageInfos() {
		fns := make([]*Fn, 0, 1)
		ns := ""
		nsAnnotations := make(map[string]string)
		for _, file := range info.Files {
			if ns == "" {
				ns, nsAnnotations = p.scanNamespaceDoc(file)
				if ns != "" {
					_, existNS := p.Fns[ns]
					if existNS {
						err = fmt.Errorf("fnc: scan fn failed for %s is duplicated", ns)
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
		if ns == "" {
			if len(fns) > 0 {
				err = fmt.Errorf("fnc: scan fn failed for some fns has no named namespace")
				return
			}
			continue
		}
		if len(fns) == 0 {
			continue
		}
		namespace := &Namespace{
			DirPath:     filepath.Join(p.dir, info.Pkg.Path()),
			Package:     info.Pkg.Name(),
			Imports:     make([]Import, 0, 1),
			Fns:         make(map[string]*Fn),
			Annotations: nsAnnotations,
		}
		for _, fn := range fns {
			addErr := namespace.AddFn(fn)
			if addErr != nil {
				err = addErr
				return
			}
		}
		p.Fns[ns] = namespace
	}
	return
}

func (p *Project) scanNamespaceDoc(file *ast.File) (ns string, annotations map[string]string) {
	doc := file.Doc.Text()
	if doc == "" {
		return
	}
	annotations = getAnnotations(doc)
	if len(annotations) == 0 {
		return
	}
	namespace, hasNS := annotations["namespace"]
	if !hasNS {
		return
	}
	ns = namespace
	return
}

func (p *Project) scanFile(file *ast.File) (fns []*Fn, err error) {
	if file.Decls == nil {
		return
	}

	return
}
