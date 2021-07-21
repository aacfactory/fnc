package codegen

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

func parseStructField(field *ast.Field) (sf Field, ok bool) {
	sf = Field{
		Exported: false,
		Doc:      nil,
		Name:     "",
		Type:     Type{},
		Tags:     nil,
	}
	// doc
	if field.Doc != nil && field.Doc.List != nil {
		doc := make([]string, 0, 1)
		for _, comment := range field.Doc.List {
			doc = append(doc, comment.Text)
		}
		sf.Doc = doc
	}
	// name
	name := field.Names[0].Name
	sf.Name = name
	sf.Exported = ast.IsExported(sf.Name)
	// tags
	if field.Tag != nil {
		tags, tagsOk := parseStructFieldTags(field.Tag.Value)
		if tagsOk {
			sf.Tags = tags
		}
	}
	// type
	switch field.Type.(type) {
	case *ast.Ident:

	case *ast.BasicLit:

	case *ast.IndexExpr:

	case *ast.ArrayType:

	case *ast.StarExpr:
		expr := field.Type.(*ast.StarExpr)
		fmt.Println(reflect.TypeOf(expr.X))

	case *ast.SelectorExpr:
		// struct
		expr := field.Type.(*ast.SelectorExpr)

		// package = expr.X.(*ast.Ident).Name
		// struct = expr.Sel.Name
		fmt.Println("selector:", expr.Sel.Name, expr.Sel.Obj, reflect.TypeOf(expr.X),  expr.X.(*ast.Ident).Name)

	default:
		Log().Debugf("parse %s field failed, %v not supported", name, reflect.TypeOf(field.Type))
	}
	fmt.Println("field", name, reflect.TypeOf(field.Type))
	return
}

func parseStructFieldTags(v string) (tags []FieldTag, ok bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return
	}
	idx := strings.IndexByte(v, ':')
	if idx > 0 {
		name := v[0:idx]
		sub := v[idx+1:]
		l := strings.IndexByte(sub, '"')
		if l < 0 {
			return
		}
		r := strings.IndexByte(sub[l+1:], '"')
		if r < 0 {
			return
		}
		values := strings.Split(sub[l+1:r+1], ",")
		tag := FieldTag{
			Name:   name,
			Values: make([]string, 0, 1),
		}
		for _, value := range values {
			value = strings.TrimSpace(value)
			if value != "" {
				tag.Values = append(tag.Values, value)
			}
		}
		tags = append(tags, tag)
		if len(sub) > r+2 {
			subTags, subTagsOk := parseStructFieldTags(sub[r+2:])
			if subTagsOk {
				tags = append(tags, subTags...)
			}
		}
	}

	ok = true
	return
}
