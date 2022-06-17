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

package project

import (
	"fmt"
	"github.com/aacfactory/fnc/project/code"
	"github.com/aacfactory/fnc/project/conf"
	"github.com/aacfactory/fnc/project/model"
	"os/exec"
)

func create(g model.Generator) (err error) {
	// go mod init
	modInitErr := goModInit(g.Module.Name)
	if modInitErr != nil {
		err = modInitErr
		return
	}
	// go get
	getErr := goGetRequires(g.Module.Requires)
	if getErr != nil {
		err = getErr
		return
	}
	// config
	configErr := conf.Create(g)
	if configErr != nil {
		err = configErr
		return
	}
	codeErr := code.Create(g)
	if codeErr != nil {
		err = codeErr
		return
	}
	return
}

func goModInit(name string) (err error) {
	modInitCmd := exec.Command("go", "mod", "init", name)
	modInitErr := modInitCmd.Run()
	if modInitErr != nil {
		err = fmt.Errorf("fnc: create project failed at go mod init, %v", modInitErr)
		return
	}
	return
}

func goGetRequires(requires []string) (err error) {
	// fns
	getFns := exec.Command("go", "get", "github.com/aacfactory/fns")
	getFnsErr := getFns.Run()
	if getFnsErr != nil {
		err = fmt.Errorf("fnc: create project failed at go get github.com/aacfactory/fns, %v", getFnsErr)
		return
	}
	// requires
	if requires == nil || len(requires) == 0 {
		return
	}
	for _, require := range requires {
		get := exec.Command("go", "get", require)
		getErr := get.Run()
		if getErr != nil {
			err = fmt.Errorf("fnc: create project failed at go get %s, %v", require, getErr)
			return
		}
	}
	return
}
