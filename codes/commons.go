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
	"go/parser"
	"go/types"
	"golang.org/x/tools/go/loader"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func gopath() (s string) {
	s, _ = os.LookupEnv("GOPATH")
	if s == "" {
		panic("fnc: get GOPATH failed")
	}
	return
}

func gopathModPath() (s string) {
	s = filepath.Join(gopath(), "pkg", "mod")
	return
}

func moduleVersion(modName string) (gto bool, name string, bv int) {
	idx := strings.LastIndex(modName, "/")
	if idx < 1 {
		name = modName
		return
	}
	left := modName[0:idx]
	right := modName[idx+1:]
	if strings.Index(right, "v") == 0 {
		vers := right[1:]
		n, toIntErr := strconv.Atoi(vers)
		if toIntErr != nil {
			name = modName
			return
		}
		if n > 1 {
			gto = true
			name = left
			bv = n
			return
		}
	} else {
		name = modName
		return
	}
	return
}

func getGoFiles(dirPath string) (files []string, err error) {
	if strings.Contains(dirPath, "\\") {
		dirPath = strings.ReplaceAll(dirPath, "\\", "/")
	}
	if strings.Contains(dirPath, "/.") {
		return
	}
	walkErr := filepath.Walk(dirPath, func(path string, info fs.FileInfo, cause error) (err error) {
		if cause != nil {
			err = cause
			return
		}
		if path == dirPath {
			return
		}
		if info.IsDir() {
			return
		}
		if !strings.HasSuffix(path, ".go") {
			return
		}
		files = append(files, path)
		return
	})
	if walkErr != nil {
		err = fmt.Errorf("fnc: get go files in %s failed, %v", dirPath, walkErr)
		return
	}
	return
}

func loadProgram(pkg string, programDir string) (program *loader.Program, err error) {
	config := loader.Config{
		Cwd:         programDir,
		ParserMode:  parser.ParseComments,
		AllowErrors: true,
		TypeChecker: types.Config{
			Error: func(err error) {},
		},
	}
	goFiles, getGoFilesErr := getGoFiles(programDir)
	if getGoFilesErr != nil {
		err = getGoFilesErr
		return
	}
	config.CreateFromFilenames(pkg, goFiles...)
	program, err = config.Load()
	if err != nil {
		err = fmt.Errorf("fnc: load %s program failed, %v", pkg, err)
		return
	}
	return
}

func getImports(file *ast.File) (imports []Import) {
	if file.Imports == nil || len(file.Imports) == 0 {
		return
	}
	for _, spec := range file.Imports {
		name := strings.ReplaceAll(spec.Path.Value, "\"", "")
		aliasSpec := spec.Name
		alias := ""
		ident := name[strings.LastIndex(name, "/")+1:]
		if aliasSpec != nil && aliasSpec.Name != "" {
			alias = aliasSpec.Name
			ident = alias
		}
		import0 := Import{
			Alias: alias,
			Name:  name,
			Ident: ident,
		}
		imports = append(imports, import0)
	}
	return
}

func getAnnotations(doc string) (v map[string]string) {
	v = make(map[string]string)
	if doc == "" {
		return
	}
	doc = strings.ReplaceAll(doc, "\r", "")
	docLines := strings.Split(doc, "\n")
	readBlock := false
	blockName := ""
	block := ""
	for _, line := range docLines {
		line = strings.TrimSpace(line)
		if readBlock {
			if line == "<<<" {
				if len(block) > 1 {
					block = block[1:]
				}
				v[blockName] = block
				readBlock = false
				blockName = ""
				block = ""
				continue
			}
			block = block + "\n" + line
			continue
		}
		if strings.Index(line, "@") != 0 {
			continue
		}
		idx := strings.Index(line, " ")
		if idx < 0 {
			v[line[1:]] = ""
			continue
		}
		key := line[1:idx]
		val := line[idx+1:]
		if val == ">>>" {
			readBlock = true
			blockName = key
			continue
		}
		v[key] = val
	}
	return
}
