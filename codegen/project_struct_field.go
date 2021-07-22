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
	"strings"
)

func (p *Project) parseStructFieldFromStructType(file *ast.File, field *ast.Field) (f Field, ok bool) {

	pkgPath, _, hasPkgPath := p.PackageOfFile(file)
	if !hasPkgPath {
		return
	}
	imports := p.GetImports(file)

	f.Name = field.Names[0].Name
	f.Exported = ast.IsExported(f.Name)
	if field.Doc != nil && field.Doc.List != nil {
		for _, line := range field.Doc.List {
			comment := strings.TrimSpace(line.Text)
			if comment != "" {
				f.Doc = append(f.Doc, comment)
			}
		}
	}
	if field.Tag != nil {
		tags, hasTags := NewFieldTags(field.Tag.Value)
		if hasTags {
			f.Tags = tags
		}
	}

	switch field.Type.(type) {
	case *ast.Ident:
		expr := field.Type.(*ast.Ident)
		if expr.Obj != nil {
			// 同一个文件
			str, loaded := p.FindStruct(pkgPath, expr.Name)
			if !loaded {
				return
			}
			f.Type = Type{
				IsStruct: true,
				Struct:   &str,
			}
		} else {
			// 基础内置类型
			f.Type = Type{
				IsBasic: true,
				Name:    expr.Name,
			}
		}
		ok = true
	case *ast.SelectorExpr:
		// 值对象
		expr := field.Type.(*ast.SelectorExpr)

		structName := expr.Sel.Name
		ident, identOk := expr.X.(*ast.Ident)
		if !identOk {
			return
		}

		fieldPkgName := ident.Name
		fieldPkgPath := ""
		for _, import0 := range imports {
			if import0.Name == fieldPkgName {
				fieldPkgPath = import0.Path
				break
			}
		}
		if fieldPkgPath == "" {
			return
		}

		fieldStruct, defined := p.FindStruct(fieldPkgPath, structName)
		if !defined {
			return
		}
		f.Type = Type{
			IsStruct:  true,
			Struct:    &fieldStruct,
			InnerType: nil,
		}

		ok = true
	case *ast.StarExpr:
		// 指针
		expr := field.Type.(*ast.StarExpr)

		fmt.Println("field:", field, f.Name, "StarExpr", expr.X, reflect.TypeOf(expr.X))
		switch expr.X.(type) {
		case *ast.Ident:
			ident := expr.X.(*ast.Ident)
			str, loaded := p.FindStruct(pkgPath, ident.Name)
			if !loaded {
				return
			}
			f.Type = Type{
				IsPtr:  true,
				Struct: &str,
			}
		case *ast.SelectorExpr:
			sExpr := expr.X.(*ast.SelectorExpr)
			fmt.Println("field:", field, f.Name, "StarExpr.SelectorExpr", sExpr.Sel.Name, sExpr.Sel.Obj, reflect.TypeOf(sExpr.X))
			structName := sExpr.Sel.Name
			ident, identOk := sExpr.X.(*ast.Ident)
			if !identOk {
				return
			}

			fieldPkgName := ident.Name
			fieldPkgPath := ""
			for _, import0 := range imports {
				if import0.Name == fieldPkgName {
					fieldPkgPath = import0.Path
					break
				}
			}
			if fieldPkgPath == "" {
				return
			}

			fieldStruct, defined := p.FindStruct(fieldPkgPath, structName)
			if !defined {
				return
			}
			f.Type = Type{
				IsPtr:  true,
				Struct: &fieldStruct,
			}
			// todo
			fmt.Println("field:", field, f.Name, "StarExpr.SelectorExpr", sExpr.Sel.Name, f.Type)

		}

		ok = true
	case *ast.ArrayType:
		// todo
		expr := field.Type.(*ast.ArrayType)
		fmt.Println("field:", f.Name, "ArrayType", expr.Elt, reflect.TypeOf(expr.Elt))

		ok = true
	case *ast.MapType:
		// todo
		expr := field.Type.(*ast.MapType)
		fmt.Println("field:", f.Name, "MapType", expr.Key, reflect.TypeOf(expr.Key), expr.Value, reflect.TypeOf(expr.Value))

		ok = true
	default:
		// Log().Debugf("fnc parse field warned, %v is not supported", reflect.TypeOf(field.Type))
	}

	return
}
