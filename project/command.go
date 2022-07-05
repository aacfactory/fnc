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
	"github.com/aacfactory/fnc/project/model"
	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
)

var Command = &cli.Command{
	Name:        "create",
	Aliases:     nil,
	Usage:       "create .",
	Description: "create fns project",
	ArgsUsage:   "",
	Category:    "",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Required: false,
			Name:     "settings",
			Value:    "",
			Usage:    "project settings file",
		},
		&cli.GenericFlag{
			Name:     "u",
			Usage:    "go get -u",
			Required: false,
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		projectDir := strings.TrimSpace(ctx.Args().First())
		if projectDir == "" {
			projectDir = "."
		}
		projectDir, err = filepath.Abs(projectDir)
		if err != nil {
			err = fmt.Errorf("fnc: create failed for project path is invalid, %v", err)
			return
		}
		latest := ctx.IsSet("u")
		settingsPath := ctx.Path("settings")
		g := model.Generator{}
		if settingsPath != "" {
			settingsAbsPath, settingsAbsPathErr := filepath.Abs(settingsPath)
			if settingsAbsPathErr != nil {
				err = fmt.Errorf("fnc: create project failed, absolute representation of settings file failed, %v", settingsAbsPathErr)
				return
			}
			settingsContent, readErr := ioutil.ReadFile(settingsAbsPath)
			if readErr != nil {
				err = fmt.Errorf("fnc: create project failed, read settings file failed, %v", readErr)
				return
			}
			decodeErr := yaml.Unmarshal(settingsContent, &g)
			if decodeErr != nil {
				err = fmt.Errorf("fnc: create project failed, decode settings file failed, %v", decodeErr)
				return
			}
		} else {
			goVersion := runtime.Version()
			if len(strings.Split(goVersion, ".")) > 2 {
				goVersion = goVersion[0:strings.LastIndexByte(goVersion, '.')]
			}
			g.Module.Go = goVersion
			//
			askErr := doAsk(&g)
			if askErr != nil {
				err = fmt.Errorf("fnc: create project failed, %v", askErr)
				return
			}
		}
		g.Path = projectDir
		g.Module.Latest = latest
		err = create(g)
		return
	},
}
