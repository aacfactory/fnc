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

package main

import (
	"context"
	"fmt"
	"github.com/aacfactory/fnc/codes"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
)

const (
	Name      = "FNC"
	Version   = "v1.0.0"
	Usage     = "fnc -p={fns project dir} "
	Copyright = `Copyright 2021 Wang Min Xiang

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`
)

func main() {
	runtime.GOMAXPROCS(0)
	app := cli.NewApp()
	app.Name = Name
	app.Usage = Usage
	app.Version = Version
	app.Authors = []*cli.Author{
		{
			Name:  "Wang Min Xiang",
			Email: "wangminxiang@aacfactory.co",
		},
	}
	app.Copyright = Copyright
	app.Flags = []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "plugin",
			EnvVars:  []string{"FNC_PLUGIN"},
			Usage:    "add some plugins such as http",
			Required: false,
		},
		&cli.PathFlag{
			Name:     "project",
			Aliases:  []string{"p"},
			EnvVars:  []string{"FNC_PROJECT"},
			Usage:    "go project path",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "debug",
			EnvVars:  []string{"FNC_DEBUG"},
			Usage:    "print debug infos",
			Required: false,
		},
	}
	app.Action = func(c *cli.Context) (err error) {
		projectDir := c.Path("project")
		debug := c.Bool("debug")
		p, pErr := codes.NewProject(projectDir, debug)
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
		// todo plugins
		// plugins := c.StringSlice("plugin")
		return
	}

	if err := app.RunContext(context.Background(), os.Args); err != nil {
		fmt.Println("fnc", "failed", err)
	}

}
