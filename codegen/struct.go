package codegen

import (
	"go/ast"
	"strings"
)

func parseStruct(spec *ast.TypeSpec) (str Struct, ok bool) {
	specType, typeOk := spec.Type.(*ast.StructType)
	if !typeOk {
		return
	}

	if specType.Fields == nil || specType.Fields.NumFields() == 0 {
		return
	}

	fields := make([]Field, 0, 1)
	for _, field := range specType.Fields.List {
		// todo names
		for _, fieldNameIdent := range field.Names {
			if fieldNameIdent == nil {
				continue
			}
			fieldName := fieldNameIdent.Name

			field0 := Field{
				Exported: ast.IsExported(fieldName),
				Comment:  field.Doc.Text(),
				Name:     fieldName,
				Type:     Type{},
				Tags:     nil,
			}
			typ, typOk := parseType(field.Type)
			if typOk {
				field0.Type = typ
			}
			if field.Tag != nil {
				tags, tagsOk := parseStructTags(field.Tag.Value)
				if tagsOk {
					field0.Tags = tags
				}
			}

			fields = append(fields, field0)

		}
	}

	name := spec.Name.String()
	exported := ast.IsExported(name)

	str = Struct{
		Exported: exported,
		Comment:  spec.Comment.Text(),
		Name:     name,
		Fields:   fields,
	}

	ok = true
	return
}

func parseStructTags(v string) (tags []FieldTag, ok bool) {
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
			subTags, subTagsOk := parseStructTags(sub[r+2:])
			if subTagsOk {
				tags = append(tags, subTags...)
			}
		}
	}

	ok = true
	return
}

