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

package tc_test

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"testing"
)

func TestCmd(t *testing.T) {
	app := cli.NewApp()

	app.Commands = []*cli.Command{
		{
			Name: "create",
			Action: func(ctx *cli.Context) (err error) {
				arg := ctx.Args().First()
				fmt.Println(arg)
				return
			},
		},
	}
	app.Run(os.Args)

}
