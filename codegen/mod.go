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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Module struct {
	Name     string
	Version  string
	Requires []Require
}

type Require struct {
	Name    string
	Version string
	Replace string
}

func loadModuleFile(path string) (mod Module, err error) {
	modFilePath := filepath.Join(path, "go.mod")
	content, readErr := ioutil.ReadFile(modFilePath)
	if readErr != nil {
		err = fmt.Errorf("read go.mod failed, %v", readErr)
		return
	}
	if content == nil || len(content) == 0 {
		err = fmt.Errorf("read go.mod failed, no content")
		return
	}

	buf := bufio.NewReader(bytes.NewReader(content))
	replaces := make(map[string]string)
	reqFlag := false
	for {
		lineContent, _, readLineErr := buf.ReadLine()
		if readLineErr != nil {
			if readLineErr == io.EOF {
				break
			}
			err = fmt.Errorf("reade go.mod failed, %v", readLineErr)
		}
		line := strings.TrimSpace(string(lineContent))
		if len(line) == 0 {
			continue
		}
		// name
		if idx := strings.Index(line, "module "); idx == 0 {
			mod.Name = strings.TrimSpace(line[7+idx:])
			continue
		}
		// version
		if idx := strings.Index(line, "go "); idx == 0 {
			mod.Version = strings.TrimSpace(line[3+idx:])
			continue
		}
		// require
		if strings.Index(line, "require (") == 0 {
			reqFlag = true
			continue
		}
		if reqFlag {
			if strings.Index(line, ")") == 0 {
				reqFlag = false
				continue
			}
			if strings.Index(line, "\t") == 0 {
				line = line[2:]
			}
			items := strings.Split(line, " ")
			mod.Requires = append(mod.Requires, Require{
				Name:    items[0],
				Version: items[1],
			})
		}
		// replace
		if strings.Index(line, "replace") == 0 {
			items := strings.Split(line, " ")
			key := fmt.Sprintf("%s@%s", items[1], items[2])
			replace := items[4]
			replaces[key] = replace
			continue
		}

	}

	if mod.Name == "" {
		err = fmt.Errorf("load go.mod failed, name is missing")
		return
	}

	for key, replace := range replaces {
		name := strings.Split(key, "@")[0]
		version := strings.Split(key, "@")[1]
		for i, require := range mod.Requires {
			if require.Name == name && require.Version == version {
				require.Replace = replace
				mod.Requires[i] = require
			}
		}
	}

	return
}
