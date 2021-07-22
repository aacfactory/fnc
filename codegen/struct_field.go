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
	"go/ast"
	"reflect"
)

func parseStructField(field *ast.Field) (sf Field, ok bool) {
	sf = Field{
		Exported: false,
		Doc:      nil,
		Name:     "",
		Type:     Type{},
		Tags:     nil,
	}
	// doc
	if field.Doc != nil && field.Doc.List != nil {
		doc := make([]string, 0, 1)
		for _, comment := range field.Doc.List {
			doc = append(doc, comment.Text)
		}
		sf.Doc = doc
	}
	// name
	name := field.Names[0].Name
	sf.Name = name
	sf.Exported = ast.IsExported(sf.Name)
	// tags
	if field.Tag != nil {
		tags, tagsOk := NewFieldTags(field.Tag.Value)
		if tagsOk {
			sf.Tags = tags
		}
	}
	// type
	switch field.Type.(type) {
	case *ast.Ident:

	case *ast.BasicLit:

	case *ast.IndexExpr:

	case *ast.ArrayType:

	case *ast.StarExpr:
		expr := field.Type.(*ast.StarExpr)
		fmt.Println(reflect.TypeOf(expr.X))

	case *ast.SelectorExpr:
		// struct
		expr := field.Type.(*ast.SelectorExpr)

		// package = expr.X.(*ast.Ident).Name
		// struct = expr.Sel.Name
		fmt.Println("selector:", expr.Sel.Name, expr.Sel.Obj, reflect.TypeOf(expr.X), expr.X.(*ast.Ident).Name)

	default:
		Log().Debugf("parse %s field failed, %v not supported", name, reflect.TypeOf(field.Type))
	}
	fmt.Println("field", name, reflect.TypeOf(field.Type))
	return
}
