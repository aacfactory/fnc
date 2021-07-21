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

package goparser

import (
	"fmt"
	"go/ast"
)

func ParseFuncDecl(decl *ast.FuncDecl) (fn Func, ok bool) {

	fn.Name = decl.Name.Name
	fn.Exported = IsExported(fn.Name)
	fn.Comment = decl.Doc.Text()

	if decl.Recv != nil {
		// todo
		fn.Receiver = &FuncReceiver{
			Name: "x",
			Type: Struct{},
		}
	}

	// params
	if decl.Type == nil {
		return
	}

	params := decl.Type.Params
	if params != nil && params.NumFields() > 0 {
		for _, param := range params.List {
			if param.Names != nil {
				for _, name := range param.Names {
					p := FuncItem{}
					p.Name = name.Name
					typ, pTypeOk := ParseType(param.Type)
					fmt.Println(typ, pTypeOk)
				}
			}
		}
	}

	// results

	ok = true
	return
}
