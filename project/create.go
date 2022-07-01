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
	"github.com/aacfactory/fnc/commons"
	"github.com/aacfactory/fnc/project/code"
	"github.com/aacfactory/fnc/project/conf"
	"github.com/aacfactory/fnc/project/model"
	"os/exec"
	"time"
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
	tidyErr := goModTidy()
	if tidyErr != nil {
		err = tidyErr
		return
	}
	fmt.Println("\nfnc: create project succeed")
	return
}

func goModInit(name string) (err error) {
	loading := commons.NewLoading("go mod init", 500*time.Millisecond)
	loading.Show()
	fin := make(chan error, 1)
	go func(fin chan error) {
		modInitCmd := exec.Command("go", "mod", "init", name)
		modInitErr := modInitCmd.Run()
		if modInitErr != nil {
			fin <- fmt.Errorf("fnc: create project failed at go mod init, %v", modInitErr)
			return
		}
		close(fin)
	}(fin)
	cmdErr, ok := <-fin
	loading.Close()
	if !ok {
		return
	}
	err = cmdErr
	return
}

func goModTidy() (err error) {
	loading := commons.NewLoading("go mod tidy", 500*time.Millisecond)
	loading.Show()
	fin := make(chan error, 1)
	go func(fin chan error) {
		modTidyCmd := exec.Command("go", "mod", "tidy")
		modTidyCmdErr := modTidyCmd.Run()
		if modTidyCmdErr != nil {
			fin <- fmt.Errorf("fnc: create project failed at go mod tidy, %v", modTidyCmdErr)
			return
		}
		close(fin)
	}(fin)
	cmdErr, ok := <-fin
	loading.Close()
	if !ok {
		return
	}
	err = cmdErr
	return
}

func goGetRequires(requires []string) (err error) {
	loading := commons.NewLoading("go get github.com/aacfactory/fns", 500*time.Millisecond)
	loading.Show()
	getFnsFin := make(chan error, 1)
	go func(fin chan error) {
		getFns := exec.Command("go", "get", "github.com/aacfactory/fns")
		getFnsErr := getFns.Run()
		if getFnsErr != nil {
			fin <- fmt.Errorf("fnc: create project failed at go get github.com/aacfactory/fns, %v", getFnsErr)
			return
		}
		close(fin)
	}(getFnsFin)
	cmdErr := <-getFnsFin
	loading.Close()
	if cmdErr != nil {
		err = cmdErr
		return
	}

	// requires
	if requires == nil || len(requires) == 0 {
		return
	}
	for _, require := range requires {
		rloading := commons.NewLoading(fmt.Sprintf("go get %s", require), 500*time.Millisecond)
		rloading.Show()
		getReqFin := make(chan error, 1)
		go func(fin chan error) {
			getReq := exec.Command("go", "get", require)
			getReqErr := getReq.Run()
			if getReqErr != nil {
				fin <- fmt.Errorf("fnc: create project failed at go get %s, %v", require, getReqErr)
				return
			}
			close(fin)
		}(getReqFin)
		reqErr := <-getReqFin
		rloading.Close()
		if reqErr != nil {
			err = reqErr
			return
		}
	}
	return
}
