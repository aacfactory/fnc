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
	"github.com/urfave/cli/v2"
	"path/filepath"
	"strings"
)

var Command = &cli.Command{
	Name:        "codes",
	Aliases:     nil,
	Usage:       "codes {project path}",
	Description: "scan fns project and generate fn codes",
	ArgsUsage:   "",
	Category:    "",
	Action: func(ctx *cli.Context) (err error) {
		projectDir := strings.TrimSpace(ctx.Args().First())
		if projectDir == "" {
			err = fmt.Errorf("fnc: codes failed for project path is undefined")
			return
		}
		projectDir, err = filepath.Abs(projectDir)
		if err != nil {
			err = fmt.Errorf("fnc: codes failed for project path is invalid, %v", err)
			return
		}
		debug := ctx.Bool("debug")
		p, pErr := NewProject(projectDir, debug)
		if pErr != nil {
			err = pErr
			return
		}
		scanErr := p.Scan()
		if scanErr != nil {
			err = scanErr
			return
		}
		generateErr := p.Generate()
		if generateErr != nil {
			err = generateErr
			return
		}

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
