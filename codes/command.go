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
	"github.com/aacfactory/fnc/commons"
	"github.com/urfave/cli/v2"
	"path/filepath"
	"strings"
	"time"
)

var Command = &cli.Command{
	Name:        "codes",
	Aliases:     nil,
	Usage:       "codes {project path}",
	Description: "scan fns project and generate fn codes",
	ArgsUsage:   "",
	Category:    "",
	Action: func(ctx *cli.Context) (err error) {
		//generating := commons.NewLoading("generating ...", 500*time.Millisecond)
		//generating.Show()
		//defer generating.Close()

		projectDir := strings.TrimSpace(ctx.Args().First())
		if projectDir == "" {
			projectDir = "."
		}
		projectDir, err = filepath.Abs(projectDir)
		if err != nil {
			//loading.Close()
			err = fmt.Errorf("fnc: codes failed for project path is invalid, %v", err)
			return
		}
		scanning := commons.NewLoading("scanning ...", 500*time.Millisecond)
		scanning.Show()
		debug := ctx.Bool("debug")
		p, pErr := NewProject(projectDir, debug)
		if pErr != nil {
			scanning.Close()
			err = pErr
			return
		}

		scanErr := p.Scan()
		scanning.Close()
		if scanErr != nil {
			err = scanErr
			return
		}
		generating := commons.NewLoading("generating ...", 500*time.Millisecond)
		generating.Show()
		generateErr := p.Generate()
		generating.Close()
		if generateErr != nil {
			err = generateErr
			return
		}
		fmt.Println("generated!!!")
		return
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "debug",
			EnvVars:  []string{"FNC_DEBUG"},
			Usage:    "print debug infos",
			Required: false,
		},
	},
}
