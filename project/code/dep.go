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

func createDep(g model.Generator) (err error) {
	file := gcg.NewFileWithoutNote("modules")
	file.AddImport(gcg.NewPackage("github.com/aacfactory/fns/service"))
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

	fn := gcg.Func()
	fn.Name("dependencies")
	fn.AddResult("services", gcg.Token("[]service.Service"))
	body := gcg.Statements()
	has := false
	body.Tab().Token("services = append(").Line()
	body.Tab().Tab().Token("services,").Line()
	// sql
	_, hasSQL := g.Settings.FindDependency("sql")
	if hasSQL {
		has = true
		body.Tab().Tab().Token("sql.Service(),", gcg.NewPackage("github.com/aacfactory/fns-contrib/databases/sql")).Line()
	}
	// deploy auth
	_, hasAuth := g.Settings.FindDependency("authorizations")
	if hasAuth {
		has = true
		body.Tab().Tab().Token(
			fmt.Sprintf("authorizations.Service(),"),
			gcg.NewPackage("github.com/aacfactory/fns/service/builtin/authorizations"),
		).Line()
	}

	// deploy permissions
	_, hasPermissions := g.Settings.FindDependency("permissions")
	if hasPermissions {
		has = true
		body.Tab().Tab().Token(
			fmt.Sprintf("permissions.Service(),"),
			gcg.NewPackage("github.com/aacfactory/fns/service/builtin/permissions"),
		).Line()
	}

	// mq
	mq, hasMQ := g.Settings.FindDependency("mq")
	if hasMQ {
		has = true
		body.Tab().Tab().Token(
			fmt.Sprintf("%s.Service(),", mq.Kind),
			gcg.NewPackage(fmt.Sprintf("github.com/aacfactory/fns-contrib/message-queues/%s", mq.Kind)),
		).Line()
	}

	// cqrs
	_, hasCqrs := g.Settings.FindDependency("cqrs")
	if hasCqrs {
		has = true
		body.Tab().Tab().Token(
			"cqrs.Service(),",
			gcg.NewPackage("github.com/aacfactory/fns-contrib/cqrs"),
		).Line()
	}
	body.Tab().Token(")").Line()
	body.Return()

	if !has {
		body = gcg.Statements()
		body.Return()
	}

	fn.Body(body)
	file.AddCode(fn.Build())

	writer := gcg.FileRender(filepath.Join(g.Path, "modules", "dependencies.go"), true)
	renderErr := file.Render(writer)
	if renderErr != nil {
		err = fmt.Errorf("fnc: generate modules/dependencies.go failed, %v", renderErr)
		return
	}
	closeFileErr := writer.Close()
	if closeFileErr != nil {
		err = fmt.Errorf("fnc: generate modules/dependencies.go failed, %v", closeFileErr)
		return
	}
	return
}
