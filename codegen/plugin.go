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
	"os"
	"os/exec"
)

type GenPlugin struct {
	Command string
	Debug   bool
}

func (p *GenPlugin) Execute(arg string) (err error) {

	cmd := exec.Command(p.Command, "--debug", fmt.Sprintf("%v", p.Debug), "-s", arg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	runErr := cmd.Start()
	if runErr != nil {
		err = fmt.Errorf("fnc call plugin %s failed, %v", p.Command, runErr)
		return
	}
	waitErr := cmd.Wait()
	if waitErr != nil {
		err = fmt.Errorf("fnc wait plugin %s executing failed, %v", p.Command, waitErr)
		return
	}

	return
}
