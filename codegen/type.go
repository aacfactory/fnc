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
	"fmt"
	"go/types"
	"reflect"
)

func parseTypeFromObject(project *Project, obj types.Object) (typ Type, err error) {

	fmt.Println(obj.Pkg().Name(), obj.Pkg().Path(), obj.Name(), reflect.TypeOf(obj.Type()))
	typ.Name = obj.Name()
	typ.Package = Import{
		Path: obj.Pkg().Path(),
		Name: obj.Pkg().Name(),
	}

	objType := obj.Type()
	switch objType.(type) {
	case *types.Named:
		namedType := objType.(*types.Named)
		underlying := namedType.Underlying()
		if underlying == nil {
			err = fmt.Errorf("parse object(%v) failed", obj)
			return
		}
		structType, structTypeOk := underlying.(*types.Struct)
		if !structTypeOk {
			err = fmt.Errorf("parse object(%v) failed, it is not a struct", obj)
			return
		}
		str, structTypeErr := parseTypeStruct(project, structType)
		if structTypeErr != nil {
			err = fmt.Errorf("parse object(%v) failed, %v", obj, structTypeErr)
			return
		}
		typ.IsStruct = true
		typ.Struct = &str
	case *types.Pointer:

	case *types.Basic:

	case *types.Map:

	case *types.Struct:
		structType, structTypeOk := objType.(*types.Struct)
		if !structTypeOk {
			err = fmt.Errorf("parse object(%v) failed, it is not a struct", obj)
			return
		}
		str, structTypeErr := parseTypeStruct(project, structType)
		if structTypeErr != nil {
			err = fmt.Errorf("parse object(%v) failed, %v", obj, structTypeErr)
			return
		}
		typ.IsStruct = true
		typ.Struct = &str
	case *types.Array:

	case *types.Slice:

	default:

	}
	return
}
