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

package code

import (
	"fmt"
	"github.com/aacfactory/fnc/commons"
	"github.com/aacfactory/fnc/project/model"
	"os/exec"
	"time"
)

func Create(g model.Generator) (err error) {
	repositoryErr := createRepository(g)
	if repositoryErr != nil {
		err = repositoryErr
		return
	}
	examplesErr := createExamples(g)
	if examplesErr != nil {
		err = examplesErr
		return
	}
	mainErr := createMain(g)
	if mainErr != nil {
		err = mainErr
		return
	}

	loading := commons.NewLoading("fnc codes", 500*time.Millisecond)
	loading.Show()
	fin := make(chan error, 1)
	go func(fin chan error) {
		codesCmd := exec.Command("fnc", "codes", g.Path)
		codesCmdErr := codesCmd.Run()
		if codesCmdErr != nil {
			fin <- fmt.Errorf("fnc: create project failed at fnc codes %s, %v", g.Path, codesCmdErr)
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
