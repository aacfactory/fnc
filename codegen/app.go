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
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

func Generate(path string, plugins []string) (err error) {

	projectPath, projectPathErr := filepath.Abs(path)
	if projectPathErr != nil {
		err = fmt.Errorf("get abs filepath of %s failed, %v", projectPath, projectPathErr)
		return
	}

	projectPath = filepath.Clean(projectPath)
	projectPath = strings.ReplaceAll(projectPath, "\\", "/")

	project, projectErr := LoadProject(projectPath)
	if projectErr != nil {
		Log().Errorf("load project %s failed, %v", projectPath, projectErr)
		return
	}

	projectContent, encodeErr := json.Marshal(project)
	if encodeErr != nil {
		err = encodeErr
		return
	}

	fmt.Println(string(projectContent))

	// todo write
	return
}
