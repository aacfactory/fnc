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
	"bufio"
	"bytes"
	"io"
	"strings"
)

func parseDoc(text string) (doc []string) {
	if text == "" {
		return
	}
	buf := bufio.NewReader(bytes.NewBufferString(text))
	for {
		lineContent, _, readLineErr := buf.ReadLine()
		if readLineErr != nil {
			if readLineErr == io.EOF {
				break
			}
		}
		line := strings.TrimSpace(string(lineContent))
		if len(line) == 0 {
			continue
		}
		doc = append(doc, line)
	}
	return
}
