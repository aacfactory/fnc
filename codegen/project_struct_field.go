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
	"go/ast"
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
		}
		ok = true
	case *ast.ArrayType:
		expr := field.Type.(*ast.ArrayType)
		typ, typeOk := parseStructArrayFieldType(p, pkgPath, imports, expr.Elt)
		if typeOk {
			f.Type = typ
			ok = true
		}
	case *ast.MapType:
		// todo
		ok = false
	default:
		// Log().Debugf("fnc parse field warned, %v is not supported", reflect.TypeOf(field.Type))
	}

	return
}

func parseStructArrayFieldType(p *Project, pkgPath string, imports []Import, expr ast.Expr) (typ Type, ok bool) {
	switch expr.(type) {
	case *ast.Ident:
		xExpr := expr.(*ast.Ident)
		if xExpr.Obj != nil {
			// 同一个文件
			str, loaded := p.FindStruct(pkgPath, xExpr.Name)
			if !loaded {
				return
			}
			typ = Type{
				IsArray: true,
				InnerType: &Type{
					IsStruct: true,
					Struct:   &str,
				},
			}
		} else {
			// 基础内置类型
			typ = Type{
				IsArray: true,
				InnerType: &Type{
					IsBasic: true,
					Name:    xExpr.Name,
				},
			}
		}
		ok = true
	case *ast.SelectorExpr:
		xExpr := expr.(*ast.SelectorExpr)

		structName := xExpr.Sel.Name
		ident, identOk := xExpr.X.(*ast.Ident)
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
		typ = Type{
			IsArray: true,
			InnerType: &Type{
				IsStruct:  true,
				Struct:    &fieldStruct,
				InnerType: nil,
			},
		}
		ok = true
	case *ast.StarExpr:
		xExpr := expr.(*ast.StarExpr)
		switch xExpr.X.(type) {
		case *ast.Ident:
			ident := xExpr.X.(*ast.Ident)
			str, loaded := p.FindStruct(pkgPath, ident.Name)
			if !loaded {
				return
			}
			typ = Type{
				IsArray: true,
				InnerType: &Type{
					IsPtr:  true,
					Struct: &str,
				},
			}
			ok = true
		case *ast.SelectorExpr:
			sExpr := xExpr.X.(*ast.SelectorExpr)
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
			typ = Type{
				IsArray: true,
				InnerType: &Type{
					IsPtr:  true,
					Struct: &fieldStruct,
				},
			}
			ok = true
		}
	case *ast.ArrayType:
		xExpr := expr.(*ast.ArrayType)
		typ0, subOk := parseStructArrayFieldType(p, pkgPath, imports, xExpr.Elt)
		if subOk {
			typ = Type{
				IsArray:   true,
				InnerType: &typ0,
			}
			ok = true
		}
	}

	return
}
