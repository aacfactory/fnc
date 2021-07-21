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

package codegen_test

import (
	"fmt"
	"github.com/aacfactory/fnc/codegen"
	"go/format"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"testing"
)

func TestGenerate(t *testing.T) {

	path := `D:\studio\workspace\go\src\github.com\aacfactory\example`

	codegen.Generate(path, nil)

}

func TestFMT(t *testing.T) {

	file := "D:\\studio\\workspace\\go\\src\\github.com\\aacfactory\\example\\main.go"

	c, _ := ioutil.ReadFile(file)
	b, _ := imports.Process("", c, nil)
	d, _ := format.Source(b)
	fmt.Println(string(d))
}
