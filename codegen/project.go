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
	"go/parser"
	"go/types"
	"golang.org/x/tools/go/loader"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Project struct {
	Path           string                     `json:"path,omitempty"`
	Module         Module                     `json:"module,omitempty"`
	Program        *loader.Program            `json:"-"`
	ImportPrograms map[string]*loader.Program `json:"-"`
	Fns            []FnFile                   `json:"fns,omitempty"`
	Structs        map[string]Struct          `json:"structs,omitempty"`
	tmpStructMap   map[string]bool
}

type requiredProject struct {
	name    string
	program *loader.Program
	err     error
}

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
			Error: func(err error) {},
		},
	}

	for pkg, filenames := range files {
		config.CreateFromFilenames(pkg, filenames...)
	}
	importPrograms := make(map[string]*loader.Program)

	if mod.Requires != nil && len(mod.Requires) > 0 {
		gopath, hasGopath := os.LookupEnv("GOPATH")
		if hasGopath {
			modPrefix := filepath.Join(gopath, "pkg", "mod")
			requiredCh := make(chan *requiredProject, 64)
			wg := &sync.WaitGroup{}
			for _, require := range mod.Requires {
				wg.Add(1)
				go func(wg *sync.WaitGroup, require Require, requiredCh chan *requiredProject) {
					defer wg.Done()
					importCWD := filepath.Join(modPrefix, fmt.Sprintf("%s@%s", require.Name, require.Version))
					if require.Replace != "" {
						replace := require.Replace
						if strings.Index(replace, ".") == 0 {
							replace = filepath.Join(path, require.Replace)
						}
						importCWD = replace
					}
					configOfImport := loader.Config{
						Cwd:         importCWD,
						ParserMode:  parser.ParseComments,
						AllowErrors: true,
						TypeChecker: types.Config{
							Error: func(err error) {},
						},
					}
					importModFiles := make(map[string][]string)

					importModFilesErr := loadPackageFiles(importCWD, require.Name, importModFiles)
					if importModFilesErr != nil {
						requiredCh <- &requiredProject{
							name:    require.Name,
							program: nil,
							err:     fmt.Errorf("fnc load project failed, %v", importModFilesErr),
						}
						return
					}

					if len(importModFiles) == 0 {
						requiredCh <- &requiredProject{
							name:    require.Name,
							program: nil,
							err:     fmt.Errorf("fnc load project failed, no go files loaded"),
						}
						return
					}
					for pkg, filenames := range importModFiles {
						configOfImport.CreateFromFilenames(pkg, filenames...)
					}
					programOfImport, loadImportErr := configOfImport.Load()
					if loadImportErr != nil {
						err = fmt.Errorf("fnc load import mod project failed, %v", loadImportErr)
						return
					}
					requiredCh <- &requiredProject{
						name:    require.Name,
						program: programOfImport,
						err:     fmt.Errorf("fnc load project failed, no go files loaded"),
					}
				}(wg, require, requiredCh)
			}
			wg.Wait()
			close(requiredCh)

			for {
				r, ok := <-requiredCh
				if !ok {
					break
				}
				importPrograms[r.name] = r.program
			}
		} else {
			Log().Warnf("fnc can not get GOOATH, so it will not parse imported in go.mod")

		}

	}

	program, loadErr := config.Load()
	if loadErr != nil {
		err = fmt.Errorf("fnc load project failed, %v", loadErr)
		return
	}

	p = &Project{
		Path:           path,
		Module:         mod,
		Program:        program,
		ImportPrograms: importPrograms,
		Fns:            make([]FnFile, 0, 1),
		Structs:        make(map[string]Struct),
		tmpStructMap:   make(map[string]bool),
	}

	loadFnErr := p.loadFn()
	if loadFnErr != nil {
		err = fmt.Errorf("fnc load project failed, %v", loadErr)
		return
	}

	for _, fnFile := range p.Fns {
		for _, fn := range fnFile.Functions {
			Log().Debugf("fnc load fn: [%s]%s", fnFile.Path, fn)
		}
	}

	return
}

func loadPackageFiles(root string, modName string, pfs map[string][]string) (err error) {

	if strings.Contains(root, "\\") {
		root = strings.ReplaceAll(root, "\\", "/")
	}

	if strings.Contains(root, "/.") {
		return
	}

	modIdx := strings.Index(root, modName)
	if modIdx < 0 {
		fmt.Println(root, modName, modIdx)
		err = fmt.Errorf("load project files failed, module name is not in path")
		return
	}
	name := root[modIdx:]

	walkErr := filepath.WalkDir(root, func(path string, dir fs.DirEntry, cause error) (err error) {
		if root == path {
			return
		}
		if strings.Contains(path, "\\") {
			path = strings.ReplaceAll(path, "\\", "/")
		}
		if strings.Contains(path, "/.git/") {
			return
		}
		if dir.IsDir() {
			subErr := loadPackageFiles(path, modName, pfs)
			if subErr != nil {
				err = subErr
				return
			}
			return
		}

		info, infoErr := dir.Info()
		if infoErr != nil {
			err = infoErr
			return
		}
		if info.IsDir() {
			return
		}
		if !strings.HasSuffix(info.Name(), ".go") {
			return
		}

		subPath := path[strings.Index(path, name)+len(name)+1:]
		if strings.Index(subPath, "/") > 0 {
			return
		}

		files, hasPkg := pfs[name]
		if !hasPkg {
			files = make([]string, 0, 1)
		}

		has := false
		for _, filePath := range files {
			if filePath == path {
				has = true
				break
			}
		}
		if has {
			return
		}
		files = append(files, path)
		pfs[name] = files
		return
	})

	if walkErr != nil {
		err = fmt.Errorf("load filename from %s failed, %v", root, walkErr)
		return
	}

	return
}
