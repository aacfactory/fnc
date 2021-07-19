package main

import (
	"context"
	"fmt"
	"github.com/aacfactory/fnc/codegen"
	"github.com/urfave/cli/v2"
	"os"
)

const (
	Name      = "FNC"
	Version   = "v1.0.0"
	Usage     = "fnc --plugin={some plugins} "
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
			Aliases:  []string{"p"},
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
		plugins := c.StringSlice("plugin")
		project := c.Path("project")
		debug := c.Bool("debug")

		codegen.CreateLog(debug)

		codegen.Generate(project, plugins)

		return
	}

	if err := app.RunContext(context.Background(), os.Args); err != nil {
		fmt.Println("fnc", "failed", err)
	}

}
