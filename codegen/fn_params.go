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

func parseFnParams(project *Project, pkgPath string, imports []Import, used map[string]Import, params *ast.FieldList) (p []FuncItem, err error) {
	if params == nil || len(params.List) < 2 || len(params.List) > 2 {
		err = fmt.Errorf("parse params is invalied, must has two params, first is fns.FnContext, secend maybe a struct typed")
		return
	}
	p1, parseErr := parseFnParam1(imports, used, params.List[0])
	if parseErr != nil {
		return
	}
	p = append(p, p1)
	if len(params.List) == 2 {
		p2, parse2Err := parseFnParam2(project, pkgPath, imports, used, params.List[1])
		if parse2Err != nil {
			return
		}
		p = append(p, p2)
	}
	return
}

func parseFnParam1(imports []Import, used map[string]Import, param *ast.Field) (p FuncItem, err error) {
	if param == nil {
		err = fmt.Errorf("parse first param failed, first must be fns.FnContext")
		return
	}
	if param.Names == nil || len(param.Names) != 1 {
		err = fmt.Errorf("parse first param failed, first is fns.FnContext")
		return
	}
	p.Name = param.Names[0].Name
	fpTypeExpr, fpTypeExprOk := param.Type.(*ast.SelectorExpr)
	if !fpTypeExprOk {
		err = fmt.Errorf("parse first param failed, first is fns.FnContext")
		return
	}
	fpTypeStructName := fpTypeExpr.Sel.Name
	fpTypePkgExpr, fpTypePkgExprOk := fpTypeExpr.X.(*ast.Ident)
	if !fpTypePkgExprOk {
		err = fmt.Errorf("parse first param failed, first is fns.FnContext")
		return
	}
	fpTypeStructPkg := fpTypePkgExpr.Name

	for _, i := range imports {
		if i.Name == fpTypeStructPkg {
			if i.Path != "github.com/aacfactory/fns" {
				err = fmt.Errorf("parse first param failed, first is fns.FnContext")
				return
			}
			used[i.Name] = i
		}
	}

	if fpTypeStructName != "FnContext" {
		err = fmt.Errorf("parse first param failed, first is fns.FnContext")
		return
	}

	p.Type = Type{
		IsContext: true,
		Name:      "github.com/aacfactory/fns.FnContext",
	}

	return
}

func parseFnParam2(project *Project, pkgPath string, imports []Import, used map[string]Import, param *ast.Field) (p FuncItem, err error) {
	if param == nil {
		err = fmt.Errorf("parse second param failed, it is nil")
		return
	}
	if param.Names == nil || len(param.Names) != 1 {
		err = fmt.Errorf("parse first param failed, first is fns.FnContext, secend is a struct typed")
		return
	}

	p.Name = param.Names[0].Name

	paramTypeExpr := param.Type

	switch paramTypeExpr.(type) {
	case *ast.StarExpr:
		expr := paramTypeExpr.(*ast.StarExpr)
		xExpr, selectorExprOk := expr.X.(*ast.SelectorExpr)
		if !selectorExprOk {
			err = fmt.Errorf("parse first param failed, first is fns.FnContext, secend is a struct typed")
			return
		}
		structName := xExpr.Sel.Name
		packageIdent, packageIdentOk := xExpr.X.(*ast.Ident)
		if !packageIdentOk {
			err = fmt.Errorf("parse first param failed, first is fns.FnContext, secend is a struct typed")
			return
		}
		packageName := packageIdent.Name

		import0, hasImport := used[packageName]
		if !hasImport {
			import0, hasImport = FindImport(imports, packageName)
			if !hasImport {
				err = fmt.Errorf("parse first param failed, first is fns.FnContext, secend is a struct typed")
				return
			}
			used[packageName] = import0
		}

		fullName, _, hasStruct := project.FindStruct(import0.Path, structName)
		if !hasStruct {
			err = fmt.Errorf("parse first param failed, first is fns.FnContext, secend is a struct typed")
			return
		}
		p.Type = Type{
			IsPtr: true,
			Name:  fullName,
		}

	case *ast.SelectorExpr:
		// 值对象
		expr := paramTypeExpr.(*ast.SelectorExpr)

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

		fullName, _, defined := project.FindStruct(fieldPkgPath, structName)
		if !defined {
			return
		}
		p.Type = Type{
			IsStruct: true,
			Name:     fullName,
		}
	case *ast.Ident:
		// 同文件
		expr := paramTypeExpr.(*ast.Ident)
		if expr.Obj == nil {
			err = fmt.Errorf("parse first param failed, first is fns.FnContext, secend is a struct typed")
			return
		}
		fullName, _, loaded := project.FindStruct(pkgPath, expr.Name)
		if !loaded {
			return
		}

		p.Type = Type{
			IsStruct: true,
			Name:     fullName,
		}
	default:
		fmt.Println("p2, not supported", reflect.TypeOf(paramTypeExpr))
		err = fmt.Errorf("parse first param failed, first is fns.FnContext, secend is a struct typed")
		return
	}

	return
}
