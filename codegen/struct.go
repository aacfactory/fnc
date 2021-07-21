package codegen

import (
	"go/ast"
	"go/types"
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
		if field.Names == nil || len(field.Names) != 1 {
			return
		}
		sf, sfOk := parseStructField(field)
		if sfOk {
			fields = append(fields, sf)
		}
	}

	name := spec.Name.String()
	exported := ast.IsExported(name)
	doc := make([]string, 0, 1)
	if spec.Doc != nil && spec.Doc.List != nil {
		for _, comment := range spec.Doc.List {
			doc = append(doc, comment.Text)
		}
	}
	str = Struct{
		Exported: exported,
		Doc:      doc,
		Name:     name,
		Fields:   fields,
	}

	ok = true
	return
}

func loadStructFromType(typ types.Type) (s *Struct, ok bool) {
	st, typeOk := typ.(*types.Struct)
	if !typeOk {
		return
	}


	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		field.Type()
	}

	return
}