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
	"path/filepath"
)

func createServices(g model.Generator) (err error) {
	file := gcg.NewFile("modules")
	file.AddImport(gcg.NewPackage("github.com/aacfactory/fns/service"))

	fn := gcg.Func()
	fn.Name("services")
	fn.AddResult("services", gcg.Token("[]service.Service"))
	body := gcg.Statements()
	body.Tab().Token("services = append(").Line()
	body.Tab().Tab().Token("services,").Line()
	body.Tab().Tab().Token("dependencies()...,").Line()
	body.Tab().Token(")").Line()
	body.Tab().Token("services = append(").Line()
	body.Tab().Tab().Token("services,").Line()
	body.Tab().Tab().Token("examples.Service(),", gcg.NewPackage(fmt.Sprintf("%s/modules/examples", g.Module.Name))).Line()
	body.Tab().Token(")").Line()
	body.Return()

	fn.Body(body)
	file.AddCode(fn.Build())

	writer := gcg.FileRender(filepath.Join(g.Path, "modules", "services.go"), true)
	renderErr := file.Render(writer)
	if renderErr != nil {
		err = fmt.Errorf("fnc: generate modules/services.go failed, %v", renderErr)
		return
	}
	closeFileErr := writer.Close()
	if closeFileErr != nil {
		err = fmt.Errorf("fnc: generate modules/services.go failed, %v", closeFileErr)
		return
	}
	return
}
