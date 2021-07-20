package codegen

import (
	"go/ast"
	"go/token"
	"reflect"
)

func parseType(expr ast.Expr) (typ Type, ok bool) {
	if expr == nil {
		return
	}
	switch expr.(type) {
	case *ast.BadExpr:
		return
	case *ast.Ident:
		ident := expr.(*ast.Ident)
		typ = Type{
			Kind:      "Indent",
			Package:   "",
			Name:      ident.String(),
			Struct:    nil,
			InnerType: nil,
		}
		ok = true
	case *ast.Ellipsis:

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
			Kind:       "basic",
			Package:    "",
			Name:       name,
			Struct:    nil,
			InnerType: nil,
		}
		ok = true
	case *ast.FuncLit:
		Log().Debugf("parse type warned, %v type is FuncLit, it is not supported", expr)
	case *ast.CompositeLit:
		Log().Debugf("parse type warned, %v type is CompositeLit, it is not supported", expr)
	case *ast.ParenExpr:
		Log().Debugf("parse type warned, %v type is ParenExpr, it is not supported", expr)
	case *ast.SelectorExpr:
		Log().Debugf("parse type warned, %v type is SelectorExpr, it is not supported", expr)
	case *ast.IndexExpr:
		indexExpr := expr.(*ast.IndexExpr)
		subType, subOk := parseType(indexExpr.X)
		if subOk {
			typ = Type{
				Kind:       "array",
				Package:    "",
				Name:       "",
				Struct:    nil,
				InnerType: &subType,
			}
		}
		ok = true
	case *ast.SliceExpr:
		sliceExpr := expr.(*ast.SliceExpr)
		subType, subOk := parseType(sliceExpr.X)
		if subOk {
			typ = Type{
				Kind:       "array",
				Package:    "",
				Name:       "",
				Struct:    nil,
				InnerType: &subType,
			}
		}
		ok = true
	case *ast.TypeAssertExpr:
		Log().Debugf("parse type warned, %v type is TypeAssertExpr, it is not supported", expr)
	case *ast.CallExpr:
		Log().Debugf("parse type warned, %v type is CallExpr, it is not supported", expr)
	case *ast.StarExpr:
		starExpr := expr.(*ast.StarExpr)
		subType, subOk := parseType(starExpr.X)
		if subOk {
			typ = Type{
				Kind:       "ptr",
				Package:    "",
				Name:       "",
				Struct:    nil,
				InnerType: &subType,
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
		subType, subOk := parseType(arrayType.Elt)
		if subOk {
			typ = Type{
				Kind:       "array",
				Package:    "",
				Name:       "",
				Struct:    nil,
				InnerType: &subType,
			}
		}
		ok = true
	case *ast.StructType:
		structType := expr.(*ast.StructType)
		subs := make([]Type, 0, 1)
		if structType.Fields != nil && structType.Fields.NumFields() > 0 {
			for _, field := range structType.Fields.List {
				subType, subOk := parseType(field.Type)
				if subOk {
					subs = append(subs, subType)
				}
			}
		}
		typ = Type{
			Kind:       "struct",
			Package:    "",
			Name:       "",
			Struct:    nil,
			InnerType: nil,
		}
		ok = true
	case *ast.FuncType:
		Log().Debugf("parse type warned, %v type is FuncType, it is not supported", expr)
	case *ast.InterfaceType:
		Log().Debugf("parse type warned, %v type is InterfaceType, it is not supported", expr)
	case *ast.MapType:
		Log().Debugf("parse type warned, %v type is MapType, it is not supported", expr)
	case *ast.ChanType:
		Log().Debugf("parse type warned, %v is not supported", reflect.TypeOf(expr))
	default:
		Log().Debugf("parse type warned, %v is not supported", reflect.TypeOf(expr))
		return
	}
	return
}
