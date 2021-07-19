package goparser

import "go/ast"

func IsExported(name string) bool {
	return name[0] >= 'A' && name[0] <= 'Z'
}

func ParseType(expr ast.Expr) (typ Type) {

	return
}