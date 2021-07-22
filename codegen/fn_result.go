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
)

func (p *Project) parseFnResults(pkgPath string, imports []Import, used map[string]Import, result *ast.FieldList) (results []FuncItem, err error) {
	if result == nil || len(result.List) < 2 || len(result.List) > 2 {
		err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
		return
	}
	if result.NumFields() == 1 {
		rErr, parseErr := parseFnResultError(imports, used, result.List[0])
		if parseErr != nil {
			err = parseErr
			return
		}
		results = append(results, rErr)
		return
	}

	r1, parseR1Err := parseFnResult(p, pkgPath, imports, used, result.List[0])
	if parseR1Err != nil {
		err = parseR1Err
		return
	}
	rErr, parseRErrorErr := parseFnResultError(imports, used, result.List[1])
	if parseRErrorErr != nil {
		err = parseRErrorErr
		return
	}
	results = append(results, r1, rErr)

	return
}

func parseFnResult(project *Project, pkgPath string, imports []Import, used map[string]Import, result *ast.Field) (p FuncItem, err error) {
	if result == nil {
		err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
		return
	}
	if result.Names == nil || len(result.Names) != 1 {
		err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
		return
	}
	p.Name = result.Names[0].Name
	resultTypeExpr := result.Type
	switch result.Type.(type) {
	case *ast.Ident:
		// 同文件
		expr := resultTypeExpr.(*ast.Ident)
		if expr.Obj == nil {
			err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
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
	case *ast.SelectorExpr:
		// 值对象
		expr := resultTypeExpr.(*ast.SelectorExpr)

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
	case *ast.StarExpr:
		expr := resultTypeExpr.(*ast.StarExpr)
		xExpr, selectorExprOk := expr.X.(*ast.SelectorExpr)
		if !selectorExprOk {
			err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
			return
		}
		structName := xExpr.Sel.Name
		packageIdent, packageIdentOk := xExpr.X.(*ast.Ident)
		if !packageIdentOk {
			err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
			return
		}
		packageName := packageIdent.Name

		import0, hasImport := used[packageName]
		if !hasImport {
			import0, hasImport = FindImport(imports, packageName)
			if !hasImport {
				err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
				return
			}
			used[packageName] = import0
		}

		fullName, _, hasStruct := project.FindStruct(import0.Path, structName)
		if !hasStruct {
			err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
			return
		}
		p.Type = Type{
			IsPtr:    true,
			IsStruct: true,
			Name:     fullName,
		}
	case *ast.ArrayType:
		expr := resultTypeExpr.(*ast.ArrayType)
		typ, typeOk := project.parseStructFieldArrayType(pkgPath, imports, expr.Elt)
		if !typeOk {
			err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
			return
		}
		p.Type = typ
	default:
		err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
		return
	}

	return
}

func parseFnResultError(imports []Import, used map[string]Import, result *ast.Field) (p FuncItem, err error) {
	if result == nil {
		err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
		return
	}
	if result.Names == nil || len(result.Names) != 1 {
		err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
		return
	}
	p.Name = result.Names[0].Name
	switch result.Type.(type) {
	case *ast.Ident:
		expr := result.Type.(*ast.Ident)
		if expr.Name == "error" {
			p.Type = Type{
				IsBasic: true,
				IsErr:   true,
				Name:    "error",
			}
			return
		}
		err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
		return
	case *ast.SelectorExpr:
		expr := result.Type.(*ast.SelectorExpr)
		x, xOk := expr.X.(*ast.Ident)
		if !xOk {
			err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
			return
		}
		if expr.Sel.Name == "CodeError" && x.Name == "errors" {
			p.Type = Type{
				IsErr: true,
				Name:  "github.com/aacfactory/errors.CodeError",
			}
			return
		}
		err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
		return
	default:
		err = fmt.Errorf("parse results is invalied, must has [1,2) params, first can be array or struct, secend must be error")
		return
	}
	return
}
