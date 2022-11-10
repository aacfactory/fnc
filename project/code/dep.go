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
	authComponents := make([]string, 0, 1)
	rbacComponents := make([]string, 0, 1)
	for _, dependency := range g.Settings.Dependencies {
		switch dependency.Name {
		case "authorizations":
			es := strings.Split(dependency.Kind, ":")
			switch es[0] {
			case "jwt":
				authComponents = append(authComponents, "jwt.Component()")
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/authorizations/encoding/jwt", "jwt"))
			default:
				break
			}
			switch es[1] {
			case "redis":
				authComponents = append(authComponents, "authredis.Component()")
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/authorizations/store/redis", "authredis"))
			case "postgres":
				authComponents = append(authComponents, "authpostgre.Component()")
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/authorizations/store/postgres", "authpostgre"))
			case "mysql":
				authComponents = append(authComponents, "authmysql.Component()")
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/authorizations/store/mysql", "authmysql"))
			case "dgraph":
				authComponents = append(authComponents, "authdgraph.Component()")
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/authorizations/store/dgraph", "authdgraph"))
			case "rgraph":
				authComponents = append(authComponents, "authrgraph.Component()")
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/authorizations/store/rgraph", "authrgraph"))
			default:
				break
			}
		case "rbac":
			switch dependency.Kind {
			case "postgres":
				rbacComponents = append(rbacComponents, "rbacpostgres.Component()")
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/rbac/store/postgres", "rbacpostgres"))
			case "mysql":
				rbacComponents = append(rbacComponents, "rbacmysql.Component()")
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/rbac/store/mysql", "rbacmysql"))
			case "dgraph":
				rbacComponents = append(rbacComponents, "rbacdgraph.Component()")
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/rbac/store/dgraph", "rbacdgraph"))
			case "rgraph":
				rbacComponents = append(rbacComponents, "rbacrgraph.Component()")
				file.AddImport(gcg.NewPackageWithAlias("github.com/aacfactory/fns-contrib/rbac/store/rgraph", "rbacrgraph"))
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
			fmt.Sprintf("authorizations.Service(%s),", strings.Join(authComponents, ", ")),
			gcg.NewPackage("github.com/aacfactory/fns/service/builtin/authorizations"),
		).Line()
	}

	// deploy permissions
	_, hasPermissions := g.Settings.FindDependency("rbac")
	if hasPermissions {
		has = true
		body.Tab().Tab().Token(
			fmt.Sprintf("rbac.Service(%s),", strings.Join(rbacComponents, ", ")),
			gcg.NewPackage("github.com/aacfactory/fns/service/builtin/rbac"),
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
