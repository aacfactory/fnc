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

package codes

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/loader"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func gopath() (s string) {
	s, _ = os.LookupEnv("GOPATH")
	if s == "" {
		panic("fnc: get GOPATH failed")
	}
	return
}

func gopathModPath() (s string) {
	s = filepath.Join(gopath(), "pkg", "mod")
	return
}

func moduleVersion(modName string) (gto bool, name string, bv int) {
	idx := strings.LastIndex(modName, "/")
	if idx < 1 {
		name = modName
		return
	}
	left := modName[0:idx]
	right := modName[idx+1:]
	if strings.Index(right, "v") == 0 {
		vers := right[1:]
		n, toIntErr := strconv.Atoi(vers)
		if toIntErr != nil {
			name = modName
			return
		}
		if n > 1 {
			gto = true
			name = left
			bv = n
			return
		}
	} else {
		name = modName
		return
	}
	return
}

func loadProgram(pkg string, programDir string) (program *loader.Program, err error) {
	gfs, fsErr := NewFileSet(programDir)
	if fsErr != nil {
		err = fsErr
		return
	}
	config := loader.Config{
		Fset:       token.NewFileSet(),
		ParserMode: parser.ParseComments,
		TypeChecker: types.Config{
			Error: func(err error) {
			},
		},
		TypeCheckFuncBodies: nil,
		Build:               nil,
		Cwd:                 programDir,
		AllowErrors:         true,
		CreatePkgs:          make([]loader.PkgSpec, 0, 1),
		ImportPkgs:          nil,
		FindPackage:         nil,
		AfterTypeCheck:      nil,
	}
	goFiles := gfs.MapToPackageFiles()
	for gfp, gf := range goFiles {
		if len(gf) > 0 {
			config.CreateFromFilenames(gfp, gf...)

		}
	}
	program, err = config.Load()
	if err != nil {
		err = fmt.Errorf("fnc: load %s program failed, %v", pkg, err)
		return
	}
	return
}

func getImports(file *ast.File) (imports []Import) {
	if file.Imports == nil || len(file.Imports) == 0 {
		return
	}
	for _, spec := range file.Imports {
		name := strings.ReplaceAll(spec.Path.Value, "\"", "")
		aliasSpec := spec.Name
		alias := ""
		ident := name[strings.LastIndex(name, "/")+1:]
		if aliasSpec != nil && aliasSpec.Name != "" {
			alias = aliasSpec.Name
			ident = alias
		}
		import0 := Import{
			Alias: alias,
			Name:  name,
			Ident: ident,
		}
		imports = append(imports, import0)
	}
	return
}

func getStructFieldTag(tag string) (v map[string]string) {
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := tag[:i]
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := tag[:i+1]
		value, err := strconv.Unquote(qvalue)
		if err != nil {
			continue
		}
		v[name] = value
		tag = tag[i+1:]
	}
	return
}

func getAnnotations(doc string) (v map[string]string) {
	v = make(map[string]string)
	if doc == "" || !strings.Contains(doc, "@") {
		return
	}
	doc = strings.ReplaceAll(doc, "\r", "")
	docLines := strings.Split(doc, "\n")
	readBlock := false
	blockName := ""
	block := ""
	for _, line := range docLines {
		line = strings.TrimSpace(line)
		if readBlock {
			if line == "<<<" {
				if len(block) > 1 {
					block = block[1:]
				}
				v[blockName] = block
				readBlock = false
				blockName = ""
				block = ""
				continue
			}
			block = block + "\n" + line
			continue
		}
		if strings.Index(line, "@") != 0 {
			continue
		}
		idx := strings.Index(line, " ")
		if idx < 0 {
			v[line[1:]] = ""
			continue
		}
		key := line[1:idx]
		val := line[idx+1:]
		if val == ">>>" {
			readBlock = true
			blockName = key
			continue
		}
		v[key] = val
	}
	return
}

func newTypeFromSpec(mod *Module, imports []Import, pkgName string, pkgAlias string, target ast.Expr) (typ *Type, err error) {
	switch target.(type) {
	case *ast.Ident:
		expr := target.(*ast.Ident)
		if expr.Obj != nil {
			// 同一个文件
			structType, has0, getErr := mod.GetStruct(pkgName, pkgAlias, expr.Obj.Name)
			if getErr != nil {
				err = getErr
				return
			}
			if !has0 {
				err = fmt.Errorf("can not get struct %s/%s", pkgName, expr.Obj.Name)
				return
			}
			typ = &Type{
				Kind:   StructTypeKind,
				Indent: "",
				Struct: structType,
				X:      nil,
			}
		} else {
			// 基础内置类型
			typ = &Type{
				Kind:   strings.ToLower(expr.Name),
				Indent: expr.Name,
				Struct: nil,
				X:      nil,
			}
		}
	case *ast.SelectorExpr:
		expr := target.(*ast.SelectorExpr)
		structName := expr.Sel.Name
		ident, identOk := expr.X.(*ast.Ident)
		if !identOk {
			err = fmt.Errorf("parse struct %s/%s failed", pkgName, structName)
			return
		}
		targetPkgName := ident.Name
		var targetPkg *Import
		for _, import0 := range imports {
			if import0.Ident == targetPkgName {
				targetPkg = &import0
				break
			}
		}
		if targetPkg == nil {
			err = fmt.Errorf("parse struct %s/%s failed", pkgName, structName)
			return
		}
		structType, has0, getErr := mod.GetStruct(targetPkg.Name, targetPkg.Alias, structName)
		if getErr != nil {
			err = getErr
			return
		}
		if !has0 {
			err = fmt.Errorf("fnc: can not get struct %s/%s", pkgName, structName)
			return
		}
		typ = &Type{
			Kind:   StructTypeKind,
			Indent: structName,
			Struct: structType,
			X:      nil,
		}
	case *ast.StarExpr:
		expr := target.(*ast.StarExpr)
		switch expr.X.(type) {
		case *ast.Ident:
			ident := expr.X.(*ast.Ident)
			structName := ident.Name
			structType, has0, getErr := mod.GetStruct(pkgName, pkgAlias, structName)
			if getErr != nil {
				err = getErr
				return
			}
			if !has0 {
				err = fmt.Errorf("can not get struct %s/%s", pkgName, structName)
				return
			}
			typ = &Type{
				Kind:   StructTypeKind,
				Indent: structName,
				Struct: structType,
				X:      nil,
			}
		case *ast.SelectorExpr:
			sExpr := expr.X.(*ast.SelectorExpr)
			structName := sExpr.Sel.Name
			ident, identOk := sExpr.X.(*ast.Ident)
			if !identOk {
				return
			}

			targetPkgName := ident.Name
			var targetPkg *Import
			for _, import0 := range imports {
				if import0.Ident == targetPkgName {
					targetPkg = &import0
					break
				}
			}
			if targetPkg == nil {
				err = fmt.Errorf("parse struct %s/%s failed", pkgName, structName)
				return
			}

			structType, has0, getErr := mod.GetStruct(targetPkg.Name, targetPkg.Alias, structName)
			if getErr != nil {
				err = getErr
				return
			}
			if !has0 {
				err = fmt.Errorf("can not get struct %s/%s", pkgName, structName)
				return
			}
			typ = &Type{
				Kind:   StructTypeKind,
				Indent: structName,
				Struct: structType,
				X:      nil,
			}
		}
	case *ast.SliceExpr, *ast.ArrayType:
		var expr ast.Expr
		sExpr, sOk := target.(*ast.SliceExpr)
		if sOk {
			expr = sExpr.X
		} else {
			aExpr := target.(*ast.ArrayType)
			expr = aExpr.Elt
		}
		x, xErr := newTypeFromSpec(mod, imports, pkgName, pkgAlias, expr)
		if xErr != nil {
			err = xErr
			return
		}
		typ = &Type{
			Kind:   ArrayTypeKind,
			Indent: "",
			Struct: nil,
			X:      x,
		}
	case *ast.MapType:
		expr := target.(*ast.MapType)
		keyExpr, keyOk := expr.Key.(*ast.Ident)
		if !keyOk {
			err = fmt.Errorf("parse map failed key is not build-in type")
		}
		valueExpr := expr.Value
		x, xErr := newTypeFromSpec(mod, imports, pkgName, pkgAlias, valueExpr)
		if xErr != nil {
			err = xErr
			return
		}
		typ = &Type{
			Kind:   MapTypeKind,
			Indent: keyExpr.Name,
			Struct: nil,
			X:      x,
		}
	default:
		err = fmt.Errorf("parse failed for %s is not supported", reflect.TypeOf(target).String())
	}
	return
}
