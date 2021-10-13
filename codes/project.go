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
		mod: mod,
		Fns: make(map[string]*Namespace),
	}
	err = p.scan()
	return
}

type Project struct {
	mod *Module
	Fns map[string]*Namespace
}

func (p *Project) Path() (v string) {
	v = p.mod.Path
	return
}

func (p *Project) scan() (err error) {
	for _, info := range p.mod.CreatedPackageInfos() {
		fmt.Println("info", info.Pkg.Name(), info.Pkg.Path(), len(info.Files))
		//for _, file := range info.Files {
		//	fns, scanFnsErr := p.scanFile(file)
		//	if scanFnsErr != nil {
		//		err = scanFnsErr
		//		return
		//	}
		//	if fns == nil || len(fns) == 0 {
		//		continue
		//	}
		//}
	}
	return
}

func (p *Project) scanFile(file *ast.File) (fns []*Fn, err error) {
	if file.Decls == nil {
		return
	}

	return
}