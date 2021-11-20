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
	"fmt"
	"github.com/aacfactory/fnc/codes"
	"testing"
)

func TestNewProject(t *testing.T) {
	x := `D:\studio\workspace\go\src\tkh.com\tkh`
	p, pErr := codes.NewProject(x)
	if pErr != nil {
		t.Error(pErr)
		return
	}
	services := p.Fns
	for k, service := range services {
		fmt.Println(fmt.Sprintf("**********%s**********", k))
		fmt.Println(fmt.Sprintf("name: %s", service.Name()))
		fmt.Println(fmt.Sprintf("title: %s", service.Title()))
		fmt.Println(fmt.Sprintf("description: %s", service.Description()))
		fmt.Println(fmt.Sprintf("internal: %v", service.Internal()))
		fmt.Println(fmt.Sprintf("fn: %d", len(service.Fns)))
		for f, fn := range service.Fns {
			fmt.Println(fmt.Sprintf("%s >>>", f))
			fmt.Println(fmt.Sprintf("name: %s", fn.Name()))
			fmt.Println(fmt.Sprintf("title: %s", fn.Title()))
			fmt.Println(fmt.Sprintf("description: %s", fn.Description()))

		}
	}

}
