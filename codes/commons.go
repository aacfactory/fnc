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
	"go/types"
	"golang.org/x/tools/go/loader"
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

func loadProgram(pkg string, programDir string) (program *loader.Program, err error) {
	gfs, fsErr := NewFileSet(programDir)
	if fsErr != nil {
		err = fsErr
		return
	}
	config := loader.Config{
		Fset:       token.NewFileSet(),
		ParserMode: parser.ParseComments,
		TypeChecker: types.Config{
			Error: func(err error) {
			},
		},
		TypeCheckFuncBodies: nil,
		Build:               nil,
		Cwd:                 programDir,
		AllowErrors:         true,
		CreatePkgs:          make([]loader.PkgSpec, 0, 1),
		ImportPkgs:          nil,
		FindPackage:         nil,
		AfterTypeCheck:      nil,
	}
	goFiles := gfs.MapToPackageFiles()
	for gfp, gf := range goFiles {
		if len(gf) > 0 {
			config.CreateFromFilenames(gfp, gf...)

		}
	}
	program, err = config.Load()
	if err != nil {
		err = fmt.Errorf("fnc: load %s program failed, %v", pkg, err)
		return
	}
	return
}

//func getImports(file *ast.File) (imports []*Import) {
//	if file.Imports == nil || len(file.Imports) == 0 {
//		return
//	}
//	for _, spec := range file.Imports {
//		path := strings.ReplaceAll(spec.Path.Value, "\"", "")
//		aliasSpec := spec.Name
//		alias := ""
//		name := path[strings.LastIndex(path, "/")+1:]
//		if aliasSpec != nil && aliasSpec.Name != "" {
//			alias = aliasSpec.Name
//			name = alias
//		}
//		import0 := Import{
//			Alias: alias,
//			Name:  name,
//			Path:  path,
//		}
//		imports = append(imports, &import0)
//	}
//	return
//}

func getStructFieldTag(tag string) (v map[string]string) {
	v = make(map[string]string)
	if tag[0] == '`' {
		tag = tag[1:]
	}
	if tag[len(tag)-1] == '`' {
		tag = tag[:len(tag)-1]
	}
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := tag[:i+1]
		tag = tag[i+1:]
		value, err := strconv.Unquote(qvalue)
		if err != nil {
			break
		}
		v[name] = value
	}
	return
}

func getAnnotations(doc string) (v map[string]string) {
	v = make(map[string]string)
	if doc == "" || !strings.Contains(doc, "@") {
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
