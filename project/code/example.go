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
	"os/exec"
	"path/filepath"
)

func createExamples(g model.Generator) (err error) {
	dir := filepath.Join(g.Path, "modules", "examples")
	mdErr := os.MkdirAll(dir, 0600)
	if mdErr != nil {
		err = fmt.Errorf("make %s dir failed, %v", dir, mdErr)
		return
	}
	err = createExamplesDoc(dir)
	if err != nil {
		return
	}
	err = createExamplesFn(dir)
	if err != nil {
		return
	}
	codesCmd := exec.Command("fnc", "codes", g.Path)
	codesCmdErr := codesCmd.Run()
	if codesCmdErr != nil {
		err = fmt.Errorf("fnc: create project failed at fnc codes %s, %v", g.Path, codesCmdErr)
		return
	}
	return
}

func createExamplesDoc(root string) (err error) {
	file := gcg.NewFile("examples")
	file.FileComments(
		"Package examples",
		"@service users",
		"@title 案例",
		"@description 案例服务",
		"@internal false",
	)
	writer := gcg.FileRender(filepath.Join(root, "doc.go"), true)
	renderErr := file.Render(writer)
	if renderErr != nil {
		err = fmt.Errorf("fnc: generate examples/doc.go failed, %v", renderErr)
		return
	}
	closeFileErr := writer.Close()
	if closeFileErr != nil {
		err = fmt.Errorf("fnc: generate examples/doc.go failed, %v", closeFileErr)
		return
	}
	return
}

func createExamplesFn(root string) (err error) {
	file := gcg.NewFile("examples")

	file.AddImport(gcg.NewPackage("context"))
	file.AddImport(gcg.NewPackage("github.com/aacfactory/errors"))

	argument := gcg.Struct()
	argumentNameField := gcg.StructField("Name")
	argumentNameField.Comments("@title 名称", "@description 名称")
	argumentNameField.Type(gcg.Ident("string"))
	argumentNameField.Tag("json", "name")
	argumentNameField.Tag("validate", "required")
	argumentNameField.Tag("message", "name is required")
	argument.AddField(argumentNameField)
	file.AddCode(gcg.Type(
		"HelloArgument",
		argument.Build(),
		"@title Hello参数",
		"@description Hello参数",
	))

	result := gcg.Struct()
	resultNameField := gcg.StructField("Name")
	resultNameField.Comments("@title 名称", "@description 名称")
	resultNameField.Type(gcg.Ident("string"))
	resultNameField.Tag("json", "name")
	result.AddField(resultNameField)
	file.AddCode(gcg.Type(
		"HelloResult",
		result.Build(),
		"@title Hello返回值",
		"@description Hello返回值",
	))

	fn := gcg.Func()
	fn.Name("hello")
	fn.Comments(
		"@fn hello",
		"@validate true",
		"@authorization false",
		"@permission false",
		"@internal false",
		"@title 您好",
		"@description >>>",
		"您好",
		"----------",
		"errors:",
		"* examples_hello_failed",
		"<<<",
	)
	fn.AddParam("ctx", gcg.QualifiedIdent(gcg.NewPackage("context"), "Context"))
	fn.AddParam("argument", gcg.Ident("HelloArgument"))
	fn.AddResult("result", gcg.Star().Ident("HelloResult"))
	fn.AddResult("err", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/errors"), "CodeError"))

	body := gcg.Statements()
	body.Tab().Token("if argument.Name == \"error\" {").Line()
	body.Tab().Tab().Token("err = errors.ServiceError(\"hello failed\").WithMeta(\"code\", \"examples_hello_failed\")").Line()
	body.Tab().Token("}").Line()
	body.Tab().Token("result = &HelloResult{").Line()
	body.Tab().Tab().Token("Name: fmt.Sprintf(\"hello %s!\", argument.Name),", gcg.NewPackage("fmt")).Line()
	body.Tab().Token("}").Line()
	body.Tab().Return()
	fn.Body(body)
	file.AddCode(fn.Build())

	writer := gcg.FileRender(filepath.Join(root, "hello.go"), true)
	renderErr := file.Render(writer)
	if renderErr != nil {
		err = fmt.Errorf("fnc: generate examples/hello.go failed, %v", renderErr)
		return
	}
	closeFileErr := writer.Close()
	if closeFileErr != nil {
		err = fmt.Errorf("fnc: generate examples/hello.go failed, %v", closeFileErr)
		return
	}
	return
}
