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

import "github.com/urfave/cli/v2"

var Command = &cli.Command{
	Name:        "create",
	Aliases:     nil,
	Usage:       "create .",
	Description: "create fns project",
	ArgsUsage:   "",
	Category:    "",
	Action: func(context *cli.Context) (err error) {

		return
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Required: true,
			Name:     "mod",
			Value:    "",
			Usage:    "go module name",
			Aliases:  []string{"m"},
		},
		&cli.PathFlag{
			Required: true,
			Name:     "path",
			Value:    "",
			Usage:    "project home path",
			Aliases:  []string{"p"},
		},
	},
}