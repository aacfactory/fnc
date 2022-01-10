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
	"github.com/aacfactory/cases"
	"github.com/aacfactory/gcg"
	"go/ast"
	"reflect"
)

func NewType(e ast.Expr, pkgPath string, imports Imports, mod *Module) (typ *Type, err error) {
	typ, err = tryDecodeCommonType(e, imports)
	if err != nil {
		return
	}
	if typ != nil {
		return
	}
	switch e.(type) {
	case *ast.Ident:
		expr := e.(*ast.Ident)
		if expr.Obj == nil {
			_struct, structErr := NewStruct(pkgPath, expr.Name, mod)
			if structErr != nil {
				err = fmt.Errorf("decode type of %s %s failed, expr type is %v, %v", pkgPath, e, reflect.TypeOf(e), structErr)
				return
			}
			typ = &Type{
				Kind:   "struct",
				Indent: "",
				Struct: _struct,
			}
		} else {
			// 同一个文件
			if expr.Obj.Kind != ast.Typ || expr.Obj.Decl == nil {

				err = fmt.Errorf("decode type of %s %s failed, expr type is %v", pkgPath, e, reflect.TypeOf(e))
				return
			}
			declExpr, declOk := expr.Obj.Decl.(*ast.TypeSpec)
			if !declOk {
				err = fmt.Errorf("decode type of %s %s failed, expr type is %v, decl type is %v", pkgPath, e, reflect.TypeOf(e), reflect.TypeOf(expr.Obj.Decl))
				return
			}
			declType := declExpr.Type
			switch declType.(type) {
			case *ast.StructType:
				_struct, structErr := NewStruct(pkgPath, expr.Obj.Name, mod)
				if structErr != nil {
					err = fmt.Errorf("decode type of %s %s failed, expr type is %v, %v", pkgPath, e, reflect.TypeOf(e), structErr)
					return
				}
				typ = &Type{
					Kind:   "struct",
					Indent: "",
					Struct: _struct,
				}
			case *ast.ArrayType:
				eltExpr := declType.(*ast.ArrayType).Elt
				eltType, eltErr := NewType(eltExpr, pkgPath, imports, mod)
				if eltErr != nil {
					err = fmt.Errorf("decode type of %s %s failed, expr type is %v, %v", pkgPath, e, reflect.TypeOf(e), eltErr)
					return
				}
				if eltType == nil {
					err = fmt.Errorf("decode type of %s %s failed, expr type is %v, not elt", pkgPath, e, reflect.TypeOf(e))
					return
				}
				typ = &Type{
					Kind:   "array",
					Indent: "",
					Struct: nil,
					X:      eltType,
					Y:      nil,
				}
			case *ast.Ident:
				identExpr := declType.(*ast.Ident)
				if expr.Obj != nil {
					err = fmt.Errorf("decode type of %s %s failed, expr type is %v", pkgPath, e, reflect.TypeOf(e))
					return
				}
				typ = &Type{
					Kind:   "builtin",
					Indent: identExpr.Name,
				}
			case *ast.MapType:
				mapExpr := declType.(*ast.MapType)
				keyExpr, keyOk := mapExpr.Key.(*ast.Ident)
				if !keyOk {
					err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
					return
				}
				if keyExpr.Obj != nil {
					err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
					return
				}
				valType, valErr := NewType(mapExpr.Value, pkgPath, imports, mod)
				if valErr != nil {
					err = fmt.Errorf("decode type of %s failed, expr type is %v, %v", e, reflect.TypeOf(e), valErr)
					return
				}
				if valType != nil {
					typ = &Type{
						Kind:   "map",
						Indent: "",
						Struct: nil,
						X: &Type{
							Kind:   "builtin",
							Indent: keyExpr.Name,
							Struct: nil,
							X:      nil,
							Y:      nil,
						},
						Y: valType,
					}
				}
			default:
				err = fmt.Errorf("decode type of %s %s failed, expr type is %v", pkgPath, e, reflect.TypeOf(e))
				return
			}

		}
	case *ast.SelectorExpr:
		expr := e.(*ast.SelectorExpr)
		structName := expr.Sel.Name
		ident, identOk := expr.X.(*ast.Ident)
		if !identOk {
			err = fmt.Errorf("parse struct %s/%s failed", pkgPath, structName)
			return
		}
		structPkgName := ident.Name
		import_, hasImport := imports.FindByName(structPkgName)
		if !hasImport {
			err = fmt.Errorf("parse struct %s/%s failed for no %s import found in file", pkgPath, structName, structPkgName)
			return
		}
		_struct, structErr := NewStruct(import_.Path, structName, mod)
		if structErr != nil {
			err = fmt.Errorf("decode type of %s %s failed, expr type is %v, %v", pkgPath, e, reflect.TypeOf(e), structErr)
			return
		}
		typ = &Type{
			Kind:   "struct",
			Indent: "",
			Import: import_,
			Struct: _struct,
			X:      nil,
			Y:      nil,
		}
	case *ast.StarExpr:
		expr := e.(*ast.StarExpr).X
		eltType, eltErr := NewType(expr, pkgPath, imports, mod)
		if eltErr != nil {
			err = eltErr
			return
		}
		if eltType == nil {
			err = fmt.Errorf("decode type of %s %s failed, expr type is %v", pkgPath, e, reflect.TypeOf(e))
			return
		}
		typ = &Type{
			Kind:   "star",
			Indent: "",
			Import: nil,
			Struct: nil,
			X:      eltType,
			Y:      nil,
		}
	case *ast.ArrayType:
		expr := e.(*ast.ArrayType).Elt
		eltType, eltErr := NewType(expr, pkgPath, imports, mod)
		if eltErr != nil {
			err = eltErr
			return
		}
		if eltType == nil {
			err = fmt.Errorf("decode type of %s %s failed, expr type is %v", pkgPath, e, reflect.TypeOf(e))
			return
		}
		typ = &Type{
			Kind:   "array",
			Indent: "",
			Struct: nil,
			X:      eltType,
			Y:      nil,
		}
	case *ast.MapType:
		expr := e.(*ast.MapType)
		keyExpr, keyOk := expr.Key.(*ast.Ident)
		if !keyOk {
			err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
			return
		}
		if keyExpr.Obj != nil {
			err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
			return
		}
		valType, valErr := NewType(expr.Value, pkgPath, imports, mod)
		if valErr != nil {
			err = valErr
			return
		}
		if valType != nil {
			typ = &Type{
				Kind:   "map",
				Indent: "",
				Struct: nil,
				X: &Type{
					Kind:   "builtin",
					Indent: keyExpr.Name,
					Struct: nil,
					X:      nil,
					Y:      nil,
				},
				Y: valType,
			}
		}
	default:
		err = fmt.Errorf("decode type of %s %s failed, expr type is %v", pkgPath, e, reflect.TypeOf(e))
		return
	}
	return
}

type Type struct {
	Kind   string
	Indent string // built-in type name
	Import *Import
	Struct *Struct
	X      *Type // if map, x is key
	Y      *Type // if map, y is value
}

func (t *Type) ObjectKey() (v string) {
	if t.IsStruct() {
		v = t.Struct.ObjectKey()
		return
	}
	if t.IsStar() {
		v = t.X.ObjectKey()
		return
	}
	if t.IsArray() {
		v = t.X.ObjectKey() + "_array"
		atoms, _ := cases.Snake().Parse(v)
		v = cases.LowerCamel().Format(atoms)
		return
	}
	return
}

func (t *Type) generateObject() (code *gcg.Statement) {
	code = gcg.Statements()
	if t.IsBuiltin() {
		switch t.Indent {
		case "string":
			code.Token("fns.StringObjectDocument()")
		case "bool":
			code.Token("fns.BoolObjectDocument()")
		case "int", "int64", "uint", "uint32", "uint64":
			code.Token("fns.IntObjectDocument()")
		case "int32", "int8", "int16", "uint8", "uint16":
			code.Token("fns.Int32ObjectDocument()")
		case "float32":
			code.Token("fns.Float32ObjectDocument()")
		case "float64":
			code.Token("fns.Float64ObjectDocument()")
		default:
			code.Token("fns.JsonRawObjectDocument()")
		}
		return
	}
	if t.IsTime() || t.IsFnsJsonTime() {
		code.Token("fns.DateTimeObjectDocument()")
		return
	}
	if t.IsFnsJsonDate() {
		code.Token("fns.DateObjectDocument()")
		return
	}
	if t.IsFnsEmpty() {
		code.Token("fns.EmptyObjectDocument()")
		return
	}
	if t.IsFnsJsonObject() || t.IsFnsJsonArray() || t.IsFnsJsonRawMessage() {
		code.Token("fns.JsonRawObjectDocument()")
		return
	}
	if t.IsStruct() {
		code.Add(t.Struct.generateObject())
		return
	}
	if t.IsStar() {
		code.Add(t.X.Struct.generateObject())
		return
	}
	if t.IsArray() {
		code.Token("fns.ArrayObjectDocument(").Line()
		code.Token(fmt.Sprintf("\"\", \"\", \"\",")).Line()
		code.Add(t.X.generateObject()).Symbol(",").Line()
		code.Symbol(")")
		return
	}
	if t.IsMap() {
		code.Token("fns.MapObjectDocument(").Line()
		code.Token(fmt.Sprintf("\"\", \"\", \"\",")).Line()
		code.Add(t.Y.generateObject()).Symbol(",").Line()
		code.Symbol(")")
		return
	}
	return
}

func (t *Type) CodeString() (v string) {
	if t.IsBuiltin() {
		v = t.Indent
	} else if t.IsStruct() {
		import_, has := t.GetImport()
		if has {
			v = import_.Name + "." + t.Struct.Name
		} else {
			v = t.Struct.Name
		}
	} else if t.IsArray() {
		v = "[]" + t.X.CodeString()
	} else if t.IsStar() {
		v = "*" + t.X.CodeString()
	}
	return
}

func (t *Type) Annotations() (v map[string]string) {
	if t.IsBuiltin() {
		v = make(map[string]string)
		return
	}
	if t.IsStruct() {
		v = t.Struct.Annotations
		return
	}
	if t.IsStar() {
		v = t.X.Annotations()
		return
	}
	if t.IsArray() {
		v = t.X.Annotations()
		return
	}
	if t.IsMap() {
		v = t.Y.Annotations()
	}
	return
}

func (t *Type) GetImport() (v *Import, has bool) {
	if t.IsBuiltin() {
		return
	}
	if t.IsStruct() {
		has = t.Import != nil
		if has {
			v = t.Import
		}
		return
	}
	if t.IsStar() {
		v, has = t.X.GetImport()
		return
	}
	if t.IsArray() {
		v, has = t.X.GetImport()
		return
	}
	if t.IsMap() {
		v, has = t.Y.GetImport()
	}
	return
}

func (t *Type) HasStruct() bool {
	return t.Struct != nil
}

func (t *Type) IsBuiltin() bool {
	return t.Kind == "builtin"
}

func (t *Type) IsStruct() bool {
	return t.Kind == "struct"
}

func (t *Type) IsStar() bool {
	return t.Kind == "star"
}

func (t *Type) IsMap() bool {
	return t.Kind == "map"
}

func (t *Type) IsArray() bool {
	return t.Kind == "array"
}

func (t *Type) IsTime() bool {
	return t.IsStruct() && t.Indent == "time.Time"
}

func (t *Type) IsJsonRawMessage() bool {
	return t.IsStruct() && t.Indent == "json.RawMessage"
}

func (t *Type) IsFnsJsonRawMessage() bool {
	return t.IsStruct() && t.Indent == "github.com/aacfactory/json.RawMessage"
}

func (t *Type) IsFnsJsonObject() bool {
	if t.IsStar() {
		return t.X.Indent == "github.com/aacfactory/json.Object"
	}
	return false
}

func (t *Type) IsFnsEmpty() bool {
	if t.IsStruct() && t.Indent == "github.com/aacfactory/fns.Empty" {

	}
	if t.IsStar() {
		return t.X.Indent == "github.com/aacfactory/fns.Empty"
	}
	return false
}

func (t *Type) IsFnsJsonArray() bool {
	if t.IsStar() {
		return t.X.Indent == "github.com/aacfactory/json.Array"
	}
	return false
}

func (t *Type) IsFnsJsonDate() bool {
	return t.IsStruct() && t.Indent == "github.com/aacfactory/json.Date"
}

func (t *Type) IsFnsJsonTime() bool {
	return t.IsStruct() && t.Indent == "github.com/aacfactory/json.Time"
}

func (t *Type) IsFnsContext() bool {
	return t.IsStruct() && t.Indent == "github.com/aacfactory/fns.Context"
}

func (t *Type) IsFnsCodeError() bool {
	return t.IsStruct() && t.Indent == "github.com/aacfactory/errors.CodeError"
}

func tryDecodeCommonType(e ast.Expr, imports Imports) (typ *Type, err error) {
	switch e.(type) {
	case *ast.Ident:
		expr := e.(*ast.Ident)
		if expr.Obj == nil {
			if isBuiltin(expr.Name) {
				// 基础内置类型
				typ = &Type{
					Kind:   "builtin",
					Indent: expr.Name,
				}
			}
		}
	case *ast.SelectorExpr:
		expr := e.(*ast.SelectorExpr)
		importExpr, ok := expr.X.(*ast.Ident)
		if !ok {
			err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
			return
		}
		importName := importExpr.Name
		_import, has := imports.FindByName(importName)
		if !has {
			err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
			return
		}
		structName := expr.Sel.Name
		switch _import.Path {
		case "time":
			if structName != "Time" {
				err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
				return
			}
			typ = &Type{
				Kind:   "struct",
				Indent: "time.Time",
				Import: _import,
				Struct: &Struct{
					Package:     "time",
					Name:        "Time",
					Fields:      nil,
					Annotations: nil,
				},
				X: nil,
				Y: nil,
			}
		case "json":
			if structName != "RawMessage" {
				err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
				return
			}
			typ = &Type{
				Kind:   "struct",
				Indent: "json.RawMessage",
				Import: _import,
				Struct: &Struct{
					Package:     "json",
					Name:        "RawMessage",
					Fields:      nil,
					Annotations: nil,
				},
				X: nil,
				Y: nil,
			}
		case "github.com/aacfactory/json":
			switch structName {
			case "RawMessage":
				typ = &Type{
					Kind:   "struct",
					Indent: "github.com/aacfactory/json.RawMessage",
					Import: _import,
					Struct: &Struct{
						Package:     "github.com/aacfactory/json",
						Name:        "RawMessage",
						Fields:      nil,
						Annotations: nil,
					},
					X: nil,
					Y: nil,
				}
			case "Time":
				typ = &Type{
					Kind:   "struct",
					Indent: "github.com/aacfactory/json.Time",
					Import: _import,
					Struct: &Struct{
						Package:     "github.com/aacfactory/json",
						Name:        "Time",
						Fields:      nil,
						Annotations: nil,
					},
					X: nil,
					Y: nil,
				}
			case "Date":
				typ = &Type{
					Kind:   "struct",
					Indent: "github.com/aacfactory/json.Date",
					Import: _import,
					Struct: &Struct{
						Package:     "github.com/aacfactory/json",
						Name:        "Date",
						Fields:      nil,
						Annotations: nil,
					},
					X: nil,
					Y: nil,
				}
			case "Object":
				typ = &Type{
					Kind:   "struct",
					Indent: "github.com/aacfactory/json.Object",
					Import: _import,
					Struct: &Struct{
						Package:     "github.com/aacfactory/json",
						Name:        "Object",
						Fields:      nil,
						Annotations: nil,
					},
					X: nil,
					Y: nil,
				}
			case "Array":
				typ = &Type{
					Kind:   "struct",
					Indent: "github.com/aacfactory/json.Array",
					Import: _import,
					Struct: &Struct{
						Package:     "github.com/aacfactory/json",
						Name:        "Array",
						Fields:      nil,
						Annotations: nil,
					},
					X: nil,
					Y: nil,
				}
			default:
				err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
				return
			}
		case "github.com/aacfactory/errors":
			if structName != "CodeError" {
				err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
				return
			}
			typ = &Type{
				Kind:   "struct",
				Indent: "github.com/aacfactory/errors.CodeError",
				Import: _import,
				Struct: &Struct{
					Package:     "github.com/aacfactory/errors",
					Name:        "CodeError",
					Fields:      nil,
					Annotations: nil,
				},
				X: nil,
				Y: nil,
			}
		case "github.com/aacfactory/fns":
			if structName == "Context" {
				typ = &Type{
					Kind:   "struct",
					Indent: "github.com/aacfactory/fns.Context",
					Import: _import,
					Struct: &Struct{
						Package:     "github.com/aacfactory/fns",
						Name:        "Context",
						Fields:      nil,
						Annotations: nil,
					},
					X: nil,
					Y: nil,
				}
			} else if structName == "Empty" {
				typ = &Type{
					Kind:   "struct",
					Indent: "github.com/aacfactory/fns.Empty",
					Import: _import,
					Struct: &Struct{
						Package:     "github.com/aacfactory/fns",
						Name:        "Empty",
						Fields:      nil,
						Annotations: nil,
					},
					X: nil,
					Y: nil,
				}
			} else {
				err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
				return
			}
		default:

		}
	case *ast.StarExpr:
		starExpr := e.(*ast.StarExpr).X
		eltType, eltErr := tryDecodeCommonType(starExpr, imports)
		if eltErr != nil {
			err = eltErr
			return
		}
		if eltType != nil {
			typ = &Type{
				Kind:   "star",
				Indent: "",
				Struct: nil,
				X:      eltType,
				Y:      nil,
			}
		}
	case *ast.ArrayType:
		expr := e.(*ast.ArrayType).Elt
		eltType, eltErr := tryDecodeCommonType(expr, imports)
		if eltType != nil {
			err = eltErr
			return
		}
		if eltType != nil {
			typ = &Type{
				Kind:   "array",
				Indent: "",
				Struct: nil,
				X:      eltType,
				Y:      nil,
			}
		}
	case *ast.MapType:
		expr := e.(*ast.MapType)
		keyExpr, keyOk := expr.Key.(*ast.Ident)
		if !keyOk {
			err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
			return
		}
		if keyExpr.Obj != nil {
			err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
			return
		}
		valType, valErr := tryDecodeCommonType(expr.Value, imports)
		if valErr != nil {
			err = valErr
			return
		}
		if valType != nil {
			typ = &Type{
				Kind:   "map",
				Indent: "",
				Struct: nil,
				X: &Type{
					Kind:   "builtin",
					Indent: keyExpr.Name,
					Struct: nil,
					X:      nil,
					Y:      nil,
				},
				Y: valType,
			}
		}
	default:
		err = fmt.Errorf("decode type of %s failed, expr type is %v", e, reflect.TypeOf(e))
	}
	return
}

func isBuiltin(v string) (ok bool) {
	switch v {
	case "string":
		ok = true
	case "int", "int8", "int16", "int32", "int64":
		ok = true
	case "float32", "float64":
		ok = true
	case "bool":
		ok = true
	case "uint", "uint8", "uint16", "uint32", "uint64":
		ok = true
	case "complex64", "complex128":
		ok = true
	default:
		ok = false
	}
	return
}
