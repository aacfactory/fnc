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
