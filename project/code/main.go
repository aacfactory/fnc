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

package code

import (
	"fmt"
	"github.com/aacfactory/fnc/project/model"
	"github.com/aacfactory/gcg"
	"path/filepath"
	"strings"
)

func createMain(g model.Generator) (err error) {
	file := gcg.NewFileWithoutNote("main")
	// imports
	file.AddImport(gcg.NewPackage("fmt"))
	file.AddImport(gcg.NewPackage("github.com/aacfactory/fns"))
	for _, dependency := range g.Settings.Dependencies {
		switch dependency.Name {
		case "authorizations":
			es := strings.Split(dependency.Kind, ":")
			switch es[0] {
			case "jwt":
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/authorizations/encoding/jwt", "_"))
			default:
				break
			}
			switch es[1] {
			case "redis":
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/authorizations/store/redis", "_"))
			case "postgres":
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/authorizations/store/postgres", "_"))
			case "mysql":
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/authorizations/store/mysql", "_"))
			case "dgraph":
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/authorizations/store/dgraph", "_"))
			case "rgraph":
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/authorizations/store/rgraph", "_"))
			default:
				break
			}
		case "permissions":
			switch dependency.Kind {
			case "postgres":
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/permissions/store/postgres", "_"))
			case "mysql":
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/permissions/store/mysql", "_"))
			case "dgraph":
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/permissions/store/dgraph", "_"))
			case "rgraph":
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/permissions/store/rgraph", "_"))
			default:
				break
			}
		case "cluster":
			switch dependency.Kind {
			case "swarm":
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/cluster/swarm", "_"))
			case "kubernetes":
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/cluster/kubernetes", "_"))
			default:
				break
			}
		case "sql":
			switch dependency.Kind {
			case "postgres":
				file.AddImport(gcg.NewPackageWithAlias("github.com/lib/pq", "_"))
			case "mysql":
				file.AddImport(gcg.NewPackageWithAlias("github.com/go-sql-driver/mysql", "_"))
			default:
				break
			}
		default:
			break
		}
	}

	// version
	vars := gcg.Vars()
	version := gcg.Var("Version", gcg.String())
	version.Value(gcg.Literal("v0.0.1"))
	vars.Add(version)
	file.AddCode(vars.Build())

	// go:generate
	file.AddCode(gcg.Line().Token("//go:generate fnc codes ."))
	// main
	mainFn := gcg.Func()
	mainFn.Name("main")
	mainBody := gcg.Statements()
	mainBody.Tab().Comments("set system environment to make config be active, e.g.: export FNS-ACTIVE=local")
	mainBody.Tab().Token("app := fns.New(").Line()
	httpOpt, hasHttp := g.Settings.FindDependency("http")
	if hasHttp {
		switch httpOpt.Kind {
		case "http3":
			mainBody.Tab().Tab().Token("fns.Server(http3.Server()),", gcg.NewPackage("github.com/aacfactory/fns-contrib/http/http3")).Line()
			mainBody.Tab().Tab().Token("fns.Server(http3.ClientBuilder),").Line()
		default:
			break
		}
	}
	mainBody.Tab().Tab().Token("fns.Version(Version),").Line()
	mainBody.Tab().Token(")").Line().Line()

	// deploy
	mainBody.Tab().Token("deployErr := app.Deploy(").Line()
	// deploy sql
	_, hasSQL := g.Settings.FindDependency("sql")
	if hasSQL {
		mainBody.Tab().Tab().Token("sql.Service(),", gcg.NewPackage("github.com/aacfactory/fns-contrib/databases/sql")).Line()
	}
	// deploy auth
	_, hasAuth := g.Settings.FindDependency("authorizations")
	if hasAuth {
		mainBody.Tab().Tab().Token(
			fmt.Sprintf("authorizations.Service(),"),
			gcg.NewPackage("github.com/aacfactory/fns/service/builtin/authorizations"),
		).Line()
	}

	// deploy permissions
	_, hasPermissions := g.Settings.FindDependency("permissions")
	if hasPermissions {

		mainBody.Tab().Tab().Token(
			fmt.Sprintf("permissions.Service(),"),
			gcg.NewPackage("github.com/aacfactory/fns/service/builtin/permissions"),
		).Line()
	}
	// deploy cqrs
	_, hasCqrs := g.Settings.FindDependency("cqrs")
	if hasCqrs {
		mainBody.Tab().Tab().Token(
			"cqrs.Service(),",
			gcg.NewPackage("github.com/aacfactory/fns-contrib/cqrs"),
		).Line()
	}
	// deploy examples
	mainBody.Tab().Tab().Token(
		"examples.Service(),",
		gcg.NewPackage(strings.Join([]string{g.Module.Name, "modules", "examples"}, "/")),
	).Line()

	mainBody.Tab().Token(")").Line()
	mainBody.Tab().Token("if deployErr != nil {").Line()
	mainBody.Tab().Tab().Token("app.Log().Error().Caller().Message(fmt.Sprintf(\"%+v\", deployErr))").Line()
	mainBody.Tab().Token("}").Line().Line()

	// run
	mainBody.Tab().Token("if runErr := app.Run(); runErr != nil {").Line()
	mainBody.Tab().Tab().Token("app.Log().Error().Caller().Message(fmt.Sprintf(\"%+v\", runErr))").Line()
	mainBody.Tab().Token("}").Line()
	// debug log
	mainBody.Tab().Token("if app.Log().DebugEnabled() {").Line()
	mainBody.Tab().Tab().Token("app.Log().Debug().Caller().Message(\"running...\")").Line()
	mainBody.Tab().Token("}").Line().Line()
	// sync
	mainBody.Tab().Token("if syncErr := app.Sync(); syncErr != nil {").Line()
	mainBody.Tab().Tab().Token("app.Log().Error().Caller().Message(fmt.Sprintf(\"%+v\", syncErr))").Line()
	mainBody.Tab().Token("}").Line()
	// debug log
	mainBody.Tab().Token("if app.Log().DebugEnabled() {").Line()
	mainBody.Tab().Tab().Token("app.Log().Debug().Message(\"stopped!!!\")").Line()
	mainBody.Tab().Token("}").Line().Line()

	mainFn.Body(mainBody)
	file.AddCode(mainFn.Build())

	writer := gcg.FileRender(filepath.Join(g.Path, "main.go"), true)
	renderErr := file.Render(writer)
	if renderErr != nil {
		err = fmt.Errorf("fnc: generate main.go failed, %v", renderErr)
		return
	}
	closeFileErr := writer.Close()
	if closeFileErr != nil {
		err = fmt.Errorf("fnc: generate main.go failed, %v", closeFileErr)
		return
	}
	return
}
