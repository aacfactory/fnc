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
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type FileSet struct {
	Package  string
	Files    []string
	Parent   *FileSet
	Children map[string]*FileSet
}

func (v *FileSet) getPackagePrefix() (x string) {
	if v.Parent == nil {
		return
	}
	x = v.Parent.Package
	p := v.Parent.getPackagePrefix()
	if p != "" {
		x = fmt.Sprintf("%s/%s", p, x)
	}
	return
}

func (v *FileSet) MapToPackageFiles() (x map[string][]string) {
	x = make(map[string][]string)
	if len(v.Files) == 0 {
		return
	}
	key := v.Package
	packagePrefix := v.getPackagePrefix()
	if packagePrefix != "" {
		key = fmt.Sprintf("%s/%s", packagePrefix, key)
	}
	x[key] = v.Files
	for _, child := range v.Children {
		y := child.MapToPackageFiles()
		for k, v := range y {
			x[k] = v
		}
	}
	return
}

func (v FileSet) String() string {
	return fmt.Sprintf("%+v", v.MapToPackageFiles())
}

func NewFileSet(dir string) (v *FileSet, err error) {
	dir = filepath.ToSlash(dir)
	if !filepath.IsAbs(dir) {
		absFilePath, absErr := filepath.Abs(dir)
		if absErr != nil {
			err = fmt.Errorf("fnc: new fileset failed for absolute representation of %s, %v", dir, absErr)
			return
		}
		dir = absFilePath
	}
	v = &FileSet{
		Package:  "",
		Parent:   nil,
		Files:    make([]string, 0, 1),
		Children: make(map[string]*FileSet),
	}
	fset := token.NewFileSet()
	walkErr := filepath.Walk(dir, func(path string, info fs.FileInfo, cause error) (err error) {
		if cause != nil {
			err = cause
			return
		}
		path = filepath.ToSlash(path)
		if path == dir {
			return
		}
		if strings.Contains(path, "/.") {
			return
		}
		if filepath.ToSlash(filepath.Dir(path)) != dir {
			return
		}
		if info.IsDir() {
			// children
			child, childErr := NewFileSet(path)
			if childErr != nil {
				err = childErr
				return
			}
			if len(child.Files) == 0 {
				child.Package = info.Name()
			}
			child.Parent = v
			v.Children[child.Package] = child
			return
		}
		if !strings.HasSuffix(path, ".go") {
			return
		}
		if strings.HasSuffix(path, "_test.go") {
			return
		}
		// pkg
		if v.Package == "" {
			gfs, readErr := ioutil.ReadFile(path)
			if readErr != nil {
				err = readErr
				return
			}
			astFile, parseErr := parser.ParseFile(fset, path, gfs, parser.ParseComments)
			if parseErr != nil {
				err = parseErr
				return
			}
			v.Package = astFile.Name.Name
		}
		v.Files = append(v.Files, path)
		return
	})
	if walkErr != nil {
		err = fmt.Errorf("fnc: new fileset failed, %v", walkErr)
		return
	}
	return
}
