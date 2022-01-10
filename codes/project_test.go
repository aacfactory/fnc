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

package codes_test

import (
	"github.com/aacfactory/fnc/codes"
	"testing"
)

func TestNewProject(t *testing.T) {
	//x := `D:\studio\workspace\go\src\tkh.com\tkh`
	x := `D:\studio\workspace\go\src\github.com\aacfactory\fns-example\standalone`
	p, pErr := codes.NewProject(x, true)
	if pErr != nil {
		t.Error(pErr)
		return
	}
	scanErr := p.Scan()
	if scanErr != nil {
		t.Error(scanErr)
		return
	}
	//fmt.Println(fmt.Sprintf("%v", p))
	generateErr := p.Generate()
	if generateErr != nil {
		t.Error(generateErr)
		return
	}
}
