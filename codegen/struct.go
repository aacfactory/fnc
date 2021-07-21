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
