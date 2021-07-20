package goparser

import (
	"go/ast"
)

func ParseStruct(spec *ast.TypeSpec) (str Struct, ok bool) {

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
				Exported: IsExported(fieldName),
				Comment:  field.Comment.Text(),
				Name:     fieldName,
				Type:     Type{},
				Tags:     nil,
			}
			typ, typOk := ParseType(field.Type)
			if typOk {
				field0.Type = typ
			}
			if field.Tag != nil {
				tags, tagsOk := ParseStructTags(field.Tag.Value)
				if tagsOk {
					field0.Tags = tags
				}
			}

			fields = append(fields, field0)

		}
	}

	name := spec.Name.String()
	exported := IsExported(name)

	str = Struct{
		Exported: exported,
		Comment:  spec.Comment.Text(),
		Name:     name,
		Fields:   fields,
	}

	ok = true
	return
}
