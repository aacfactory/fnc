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
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"path/filepath"
)

func NewModule(modFilePath string) (mod Module, err error) {
	p, readErr := ioutil.ReadFile(modFilePath)
	if readErr != nil {
		err = fmt.Errorf("fnc: read mod file failed, %s, %v", modFilePath, readErr)
		return
	}
	mf, parseErr := modfile.Parse(modFilePath, p, func(path, version string) (string, error) {
		return version, nil
	})
	if parseErr != nil {
		err = fmt.Errorf("fnc: parse mod file failed, %s, %v", modFilePath, parseErr)
		return
	}
	mod.GoVersion = mf.Go.Version
	mod.Name = mf.Module.Mod.Path
	mod.Path = fmt.Sprintf("%s/src/%s", gopath(), mod.Name)
	mod.Requires = make([]Require, 0, 1)
	if mf.Require != nil {
		for _, require := range mf.Require {
			r := Require{
				Name:    require.Mod.Path,
				Version: require.Mod.Version,
				Replace: "",
			}
			mod.Requires = append(mod.Requires, r)
		}
	}
	if mf.Replace != nil {
		for _, replace := range mf.Replace {
			for i, require := range mod.Requires {
				if replace.Old.Path == require.Name {
					require.Replace = replace.New.Path
					require.ReplaceVersion = replace.New.Version
					mod.Requires[i] = require
				}
			}
		}
	}

	return
}

type Module struct {
	Name      string
	Path      string
	GoVersion string
	Requires  []Require
}

func (mod Module) String() (s string) {
	s = mod.Name + " " + mod.Path + "\n"
	s = s + mod.GoVersion + "\n"
	for _, require := range mod.Requires {
		s = s + require.Name + " " + require.Version
		if require.Replace != "" {
			s = s + " => " + require.Replace
		}
		s = s + " [" + require.Path() + "]"
		s = s + "\n"
	}
	return
}

type Require struct {
	Name           string
	Version        string
	Replace        string
	ReplaceVersion string
}

func (r Require) Path() (s string) {
	if r.Replace != "" {
		if r.ReplaceVersion != "" {
			f, err := filepath.Abs(r.Replace)
			if err != nil {
				panic(fmt.Sprintf("fnc: get require file path failed, %s, %v", r.Replace, err))
			}
			s = f
		} else {
			gto, name, bv := moduleVersion(r.Replace)
			if gto {
				s = fmt.Sprintf("%s/%s/v%d@%s", gopathModPath(), name, bv, r.ReplaceVersion)
			} else {
				s = fmt.Sprintf("%s/%s@%s", gopathModPath(), r.Replace, r.ReplaceVersion)
			}
		}
	} else {
		gto, name, bv := moduleVersion(r.Name)
		if gto {
			s = fmt.Sprintf("%s/%s/v%d@%s", gopathModPath(), name, bv, r.Version)
		} else {
			s = fmt.Sprintf("%s/%s@%s", gopathModPath(), r.Name, r.Version)
		}
	}
	return
}
