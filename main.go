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
	"github.com/aacfactory/fnc/project"
	"github.com/aacfactory/fnc/ssc"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
)

const (
	Name      = "FNC"
	Version   = "v1.10.10"
	Usage     = "see COMMANDS"
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
	app.Commands = []*cli.Command{
		codes.Command,
		project.Command,
		ssc.Command,
	}
	if err := app.RunContext(context.Background(), os.Args); err != nil {
		fmt.Println("fnc", "failed", err)
	}

}
