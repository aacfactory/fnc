package goparser

import (
	"go/ast"
	"go/token"
	"reflect"
)

func ParseType(expr ast.Expr) (typ Type, ok bool) {
	if expr == nil {
		return
	}
	switch expr.(type) {
	case *ast.BadExpr:
		return
	case *ast.Ident:
		ident := expr.(*ast.Ident)
		typ = Type{
			Kind:       "Indent",
			Package:    "",
			Name:       ident.String(),
			Ptr:        false,
			InnerTypes: nil,
		}
		ok = true
	case *ast.Ellipsis:
		// An Ellipsis node stands for the "..." type in a
		ellipsis := expr.(*ast.Ellipsis)
		subType, subOk := ParseType(ellipsis.Elt)
		if subOk {
			typ = Type{
				Kind:       "Ellipsis",
				Package:    "",
				Name:       "",
				Ptr:        false,
				InnerTypes: []Type{subType},
			}
		}
		ok = true
	case *ast.BasicLit:
		basicLit := expr.(*ast.BasicLit)
		name := ""
		switch basicLit.Kind {
		case token.INT:
			name = "int"
		case token.FLOAT:
			name = "float64"
		case token.IMAG:
			name = "imag"
		case token.CHAR:
			name = "char"
		case token.STRING:
			name = "string"
		}
		typ = Type{
			Kind:       "BasicLit",
			Package:    "",
			Name:       name,
			Ptr:        false,
			InnerTypes: nil,
		}
		ok = true
	case *ast.FuncLit:
		// todo in and out
		// funcLit := expr.(*ast.FuncLit)
		// funcLit.Type.Params
		typ = Type{
			Kind:       "FuncLit",
			Package:    "",
			Name:       "func",
			Ptr:        false,
			InnerTypes: nil,
		}
		ok = true
	case *ast.CompositeLit:
		Log().Debugf("parse type warned, %v type is CompositeLit, it is not supported", expr)
	case *ast.ParenExpr:
		Log().Debugf("parse type warned, %v type is ParenExpr, it is not supported", expr)
	case *ast.SelectorExpr:
		Log().Debugf("parse type warned, %v type is SelectorExpr, it is not supported", expr)
	case *ast.IndexExpr:
		indexExpr := expr.(*ast.IndexExpr)
		subType, subOk := ParseType(indexExpr.X)
		if subOk {
			typ = Type{
				Kind:       "Array",
				Package:    "",
				Name:       "",
				Ptr:        false,
				InnerTypes: []Type{subType},
			}
		}
		ok = true
	case *ast.SliceExpr:
		sliceExpr := expr.(*ast.SliceExpr)
		subType, subOk := ParseType(sliceExpr.X)
		if subOk {
			typ = Type{
				Kind:       "Array",
				Package:    "",
				Name:       "",
				Ptr:        false,
				InnerTypes: []Type{subType},
			}
		}
		ok = true
	case *ast.TypeAssertExpr:
		Log().Debugf("parse type warned, %v type is TypeAssertExpr, it is not supported", expr)
	case *ast.CallExpr:
		Log().Debugf("parse type warned, %v type is CallExpr, it is not supported", expr)
	case *ast.StarExpr:
		starExpr := expr.(*ast.StarExpr)
		subType, subOk := ParseType(starExpr.X)
		if subOk {
			typ = Type{
				Kind:       "Ptr",
				Package:    "",
				Name:       "",
				Ptr:        false,
				InnerTypes: []Type{subType},
			}
		}
		ok = true
	case *ast.UnaryExpr:
		Log().Debugf("parse type warned, %v type is UnaryExpr, it is not supported", expr)
	case *ast.BinaryExpr:
		Log().Debugf("parse type warned, %v type is BinaryExpr, it is not supported", expr)
	case *ast.KeyValueExpr:
		Log().Debugf("parse type warned, %v type is KeyValueExpr, it is not supported", expr)
	case *ast.ArrayType:
		arrayType := expr.(*ast.ArrayType)
		subType, subOk := ParseType(arrayType.Elt)
		if subOk {
			typ = Type{
				Kind:       "Array",
				Package:    "",
				Name:       "",
				Ptr:        false,
				InnerTypes: []Type{subType},
			}
		}
		ok = true
	case *ast.StructType:
		structType := expr.(*ast.StructType)
		subs := make([]Type, 0, 1)
		if structType.Fields != nil && structType.Fields.NumFields() > 0 {
			for _, field := range structType.Fields.List {
				subType, subOk := ParseType(field.Type)
				if subOk {
					subs = append(subs, subType)
				}
			}
		}
		typ = Type{
			Kind:       "Struct",
			Package:    "",
			Name:       "",
			Ptr:        false,
			InnerTypes: subs,
		}
		ok = true
	case *ast.FuncType:
		typ = Type{
			Kind:       "Func",
			Package:    "",
			Name:       "",
			Ptr:        false,
			InnerTypes: nil,
		}
		ok = true
	case *ast.InterfaceType:
		typ = Type{
			Kind:       "Interface",
			Package:    "",
			Name:       "",
			Ptr:        false,
			InnerTypes: nil,
		}
		ok = true
	case *ast.MapType:
		mapType := expr.(*ast.MapType)
		typ = Type{
			Kind:       "Interface",
			Package:    "",
			Name:       "",
			Ptr:        false,
			InnerTypes: nil,
		}
		keyType, keyTypeOk := ParseType(mapType.Key)
		valueType, valueTypeOk := ParseType(mapType.Value)
		if keyTypeOk && valueTypeOk {
			typ.InnerTypes = []Type{keyType, valueType}
		}
		ok = true
	case *ast.ChanType:
		Log().Debugf("parse type warned, %v is not supported", reflect.TypeOf(expr))
	default:
		Log().Debugf("parse type warned, %v is not supported", reflect.TypeOf(expr))
		return
	}
	return
}
