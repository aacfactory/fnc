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

package codes

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func gopath() (s string) {
	s, _ = os.LookupEnv("GOPATH")
	if s == "" {
		panic("fnc: get GOPATH failed")
	}
	return
}

func gopathModPath() (s string) {
	s = fmt.Sprintf("%s/pkg/mod", gopath())
	return
}

func moduleVersion(modName string) (gto bool, name string, bv int) {
	idx := strings.LastIndex(modName, "/")
	if idx < 1 {
		name = modName
		return
	}
	left := modName[0:idx]
	right := modName[idx+1:]
	if strings.Index(right, "v") == 0 {
		vers := right[1:]
		n, toIntErr := strconv.Atoi(vers)
		if toIntErr != nil {
			name = modName
			return
		}
		if n > 1 {
			gto = true
			name = left
			bv = n
			return
		}
	} else {
		name = modName
		return
	}
	return
}