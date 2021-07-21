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

// ssss
// bvbbb
package goparser

import (
	"context"
	"time"
)

const (
	A S = 1 //A sss
)

var (
	bb string            = "bb" //bb xx
	mm map[string]string = map[string]string{"s": "b"}
)

type S int
type SP *int

type SS struct {
	Id   string `json:"id,omitempty" xml:"id"`
	Name *SS
	TT   time.Time
}

type ff func(ctx context.Context, ss *SS, bb SS, tt time.Time) (err error)

type tt interface {
	A()
}

type Id []string

type IA map[string]interface{}
