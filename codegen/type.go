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

import "encoding/json"

type Type struct {
	IsContext   bool   `json:"isContext,omitempty"`
	IsBasic     bool   `json:"isBasic,omitempty"`
	IsStruct    bool   `json:"isStruct,omitempty"`
	IsInterface bool   `json:"isInterface,omitempty"`
	IsPtr       bool   `json:"isPtr,omitempty"`
	IsArray     bool   `json:"isArray,omitempty"`
	IsMap       bool   `json:"isMap,omitempty"`
	IsErr       bool   `json:"isErr,omitempty"`
	Name        string `json:"name,omitempty"`
	InnerTypes  []Type `json:"innerType,omitempty"`
}

func (t Type) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}
