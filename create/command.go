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

package create

import (
	"github.com/aacfactory/errors"
	"github.com/urfave/cli/v2"
	"path/filepath"
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
		&cli.BoolFlag{
			Name:     "u",
			Usage:    "go get -u",
			Value:    false,
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
			err = errors.Warning("fnc: create fns project failed").WithCause(err).WithMeta("dir", projectDir)
			return
		}
		// check mod
		return
	},
}
