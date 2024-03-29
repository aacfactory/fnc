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
	"github.com/aacfactory/errors"
	"github.com/aacfactory/forg"
	"github.com/urfave/cli/v2"
	"path/filepath"
	"strings"
)

var Command = &cli.Command{
	Name: "codes",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "debug",
			EnvVars:  []string{"FNC_DEBUG"},
			Usage:    "print debug infos",
			Required: false,
		},
		&cli.StringFlag{
			Name:      "work",
			Aliases:   []string{"w"},
			Usage:     "set workspace file path",
			Required:  false,
			EnvVars:   []string{"FNC_WORK"},
			TakesFile: false,
		},
	},
	Aliases:     nil,
	Usage:       "fnc codes {project path}",
	Description: "scan fns project and generate fn codes",
	ArgsUsage:   "",
	Category:    "",
	Action: func(ctx *cli.Context) (err error) {
		debug := ctx.Bool("debug")
		projectDir := strings.TrimSpace(ctx.Args().First())
		if projectDir == "" {
			projectDir = "."
		}
		if !filepath.IsAbs(projectDir) {
			projectDir, err = filepath.Abs(projectDir)
			if err != nil {
				err = errors.Warning("fnc: codes failed").WithCause(err).WithMeta("dir", projectDir)
				return
			}
		}
		projectDir = filepath.ToSlash(projectDir)
		work := ctx.String("work")
		var project *forg.Project
		if work != "" {
			project, err = forg.Load(projectDir, forg.WithWorkspace(work))
		} else {
			project, err = forg.Load(projectDir)
		}
		if err != nil {
			err = errors.Warning("fnc: codes failed").WithCause(err)
			return
		}
		process, codingErr := project.Coding(ctx.Context)
		if codingErr != nil {
			err = errors.Warning("fnc: codes failed").WithCause(codingErr)
			return
		}
		results := process.Start(ctx.Context)
		for {
			result, ok := <-results
			if !ok {
				if debug {
					fmt.Println("fnc: codes finished")
				}
				break
			}
			if debug {
				fmt.Println(result, "->", fmt.Sprintf("[%d/%d]", result.UnitNo, result.UnitNum), result.Data)
			}
			if result.Error != nil {
				fmt.Println(fmt.Sprintf("%+v", result.Error))
			}
		}
		return
	},
}
