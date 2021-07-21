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

func parseFnParams(project *Project, imports []Import, used map[string]Import, params *ast.FieldList) (p1 FuncItem, p2 FuncItem, err error) {
	if params == nil || len(params.List) != 2 {
		err = fmt.Errorf("parse params is invalied, must has two params, first is fns.FnContext, secend is a struct typed")
		return
	}
	//used = make(map[string]Import)
	p1, err = parseFnParam1(project, imports, used, params.List[0])
	if err != nil {
		return
	}
	p2, err = parseFnParam2(project, imports, used, params.List[1])
	if err != nil {
		return
	}
	return
}

func parseFnParam1(project *Project, imports []Import, used map[string]Import, param *ast.Field) (p FuncItem, err error) {
	if param == nil {
		err = fmt.Errorf("parse first param failed, it is nil")
		return
	}
	if param.Names == nil || len(param.Names) != 1 {
		err = fmt.Errorf("parse first param failed,first is fns.FnContext, secend is a struct typed")
		return
	}
	p.Name = param.Names[0].Name
	fpTypeExpr, fpTypeExprOk := param.Type.(*ast.SelectorExpr)
	if !fpTypeExprOk {
		err = fmt.Errorf("parse first param failed,first is fns.FnContext, secend is a struct typed")
		return
	}
	fpTypeStructName := fpTypeExpr.Sel.Name
	fpTypePkgExpr, fpTypePkgExprOk := fpTypeExpr.X.(*ast.Ident)
	if !fpTypePkgExprOk {
		err = fmt.Errorf("parse first param failed,first is fns.FnContext, secend is a struct typed")
		return
	}
	fpTypeStructPkg := fpTypePkgExpr.Name

	for _, i := range imports {
		if i.Name == fpTypeStructPkg {
			if i.Path != "github.com/aacfactory/fns" {
				err = fmt.Errorf("parse first param failed, first is fns.FnContext, secend is a struct typed")
				return
			}
			used[i.Name] = i
		}
	}

	if fpTypeStructName != "FnContext" {
		err = fmt.Errorf("parse first param failed, first is fns.FnContext, secend is a struct typed")
		return
	}

	p.Type = Type{
		IsInterface: true,
		Package: Import{
			Path: "github.com/aacfactory/fns",
			Name: fpTypeStructPkg,
		},
		Name:      "FnContext",
		InnerType: nil,
	}

	return
}

func parseFnParam2(project *Project, imports []Import, used map[string]Import, param *ast.Field) (p FuncItem, err error) {
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

		obj, hasObj := project.FindObject(import0.Path, structName)
		if !hasObj {
			err = fmt.Errorf("parse first param failed, first is fns.FnContext, secend is a struct typed")
			return
		}

		typ, typErr := parseTypeFromObject(project, obj)
		if typErr != nil {
			err = fmt.Errorf("parse first param failed, first is fns.FnContext, secend is a struct typed, %v", typErr)
			return
		}
		p.Type = typ

	case *ast.SelectorExpr:

	case *ast.Ident:
		expr := paramTypeExpr.(*ast.Ident)
		obj, hasObj := project.ObjectOf(expr)
		fmt.Println("ident:", hasObj, obj)

	default:
		fmt.Println("p2, not supported", reflect.TypeOf(paramTypeExpr))
		err = fmt.Errorf("parse first param failed, first is fns.FnContext, secend is a struct typed")
		return
	}

	return
}
