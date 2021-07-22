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
	"strings"
)

type Struct struct {
	Exported bool     `json:"exported,omitempty"`
	Doc      []string `json:"doc,omitempty"`
	Package  Import   `json:"package,omitempty"`
	Name     string   `json:"name,omitempty"`
	Fields   []Field  `json:"fields,omitempty"`
}

func (s Struct) String() string {
	return fmt.Sprintf("[%v][%s.%s][%d]", s.Exported, s.Package.Path, s.Name, len(s.Fields))
}

func (s *Struct) PutField(field Field) {
	if s.Fields == nil {
		s.Fields = make([]Field, 0, 1)
	}
	has := false
	for _, f := range s.Fields {
		if f.Name == field.Name {
			has = true
			break
		}
	}
	if !has {
		s.Fields = append(s.Fields, field)
	}
}

type Field struct {
	Exported bool       `json:"exported,omitempty"`
	Doc      []string   `json:"doc,omitempty"`
	Name     string     `json:"name,omitempty"`
	Type     Type       `json:"type,omitempty"`
	Tags     []FieldTag `json:"tags,omitempty"`
}

func NewFieldTags(v string) (tags []FieldTag, ok bool) {
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
			subTags, subTagsOk := NewFieldTags(sub[r+2:])
			if subTagsOk {
				tags = append(tags, subTags...)
			}
		}
	}

	ok = true
	return
}

type FieldTag struct {
	Name   string   `json:"name,omitempty"`
	Values []string `json:"values,omitempty"`
}
