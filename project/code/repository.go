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

package code

import (
	"fmt"
	"github.com/aacfactory/fnc/project/model"
	"github.com/aacfactory/gcg"
	"os"
	"path/filepath"
)

func createRepository(g model.Generator) (err error) {
	dir := filepath.Join(g.Path, "repository")
	mdErr := os.MkdirAll(dir, 0600)
	if mdErr != nil {
		err = fmt.Errorf("make %s dir failed, %v", dir, mdErr)
		return
	}
	helps := "see https://github.com/aacfactory/fns-contrib/tree/main/databases/sql"
	_, hasPostgres := g.Settings.FindDependency("postgres")
	if hasPostgres {
		helps = helps + "\n" + "see https://github.com/aacfactory/fns-contrib/tree/main/databases/postres"
	}
	_, hasMysql := g.Settings.FindDependency("mysql")
	if hasMysql {
		helps = helps + "\n" + "see https://github.com/aacfactory/fns-contrib/tree/main/databases/mysql"
	}
	_, hasDgraph := g.Settings.FindDependency("dgraph")
	if hasDgraph {
		helps = helps + "\n" + "see https://github.com/aacfactory/fns-contrib/tree/main/databases/dgraph"
	}
	_, hasRgraph := g.Settings.FindDependency("rgraph")
	if hasRgraph {
		helps = helps + "\n" + "see https://github.com/aacfactory/fns-contrib/tree/main/databases/rgraph"
	}
	file := gcg.NewFile("repository")
	file.FileComments(
		"Package repository",
		"write data repository here",
		helps,
	)
	writer := gcg.FileRender(filepath.Join(dir, "doc.go"), true)
	renderErr := file.Render(writer)
	if renderErr != nil {
		err = fmt.Errorf("fnc: generate repository/doc.go failed, %v", renderErr)
		return
	}
	closeFileErr := writer.Close()
	if closeFileErr != nil {
		err = fmt.Errorf("fnc: generate repository/doc.go failed, %v", closeFileErr)
		return
	}
	return
}
