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

import (
	"bufio"
	"fmt"
	"github.com/aacfactory/fnc/project/model"
	"os"
	"strings"
)

func doAsk(g *model.Generator) (err error) {
	// name
	name, nameErr := ask("enter project name: ")
	if nameErr != nil {
		err = nameErr
		return
	}
	if name == "" {
		err = fmt.Errorf("fnc: create project failed cause name is required")
		return
	}
	g.Settings.Name = name
	// mod name
	modName, modNameErr := ask("enter project module name: ")
	if modNameErr != nil {
		err = modNameErr
		return
	}
	if modName == "" {
		err = fmt.Errorf("fnc: create project failed cause module name is required")
		return
	}
	g.Module.Name = modName
	fmt.Println("modName", ">", modName)
	// http version
	needHttp3, needHttp3Err := ask("does project use http3 (y/n): ")
	if needHttp3Err != nil {
		err = needHttp3Err
		return
	}
	if strings.ToLower(needHttp3) == "y" {
		g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/http/http3")
		g.Settings.Dependencies = append(g.Settings.Dependencies, &model.Dependency{
			Name: "http",
			Kind: "http3",
		})
	}
	// cluster
	needCluster, needClusterErr := ask("does project use cluster (y/n): ")
	if needClusterErr != nil {
		err = needClusterErr
		return
	}
	if strings.ToLower(needCluster) == "y" {
		kindAsk := "cluster kind:\n\t1) default\n\t2) docker swarm\n\t3) kubernetes\nplease enter the number: "
		kindNo, kindNoErr := ask(kindAsk)
		if kindNoErr != nil {
			err = kindNoErr
			return
		}
		kind := ""
		switch kindNo {
		case "1":
			kind = "default"
		case "2":
			kind = "swarm"
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/cluster/swarm")
		case "3":
			kind = "kubernetes"
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/cluster/kubernetes")
		default:
			err = fmt.Errorf("fnc: please choose in list cqrs server kind")
			return
		}
		g.Settings.Dependencies = append(g.Settings.Dependencies, &model.Dependency{
			Name: "cluster",
			Kind: kind,
		})
	}
	// ddd
	needDDD, needDDDErr := ask("does project use CQRS(DDD) (y/n): ")
	if needDDDErr != nil {
		err = needDDDErr
		return
	}
	if strings.ToLower(needDDD) == "y" {
		kindAsk := "server kind:\n\t1) query\n\t2) command\n\t3) full\nplease enter the number: "
		kindNo, kindNoErr := ask(kindAsk)
		if kindNoErr != nil {
			err = kindNoErr
			return
		}
		kind := ""
		switch kindNo {
		case "1":
			kind = "query"
		case "2":
			kind = "command"
		case "3":
			kind = "full"
		default:
			err = fmt.Errorf("fnc: please choose in list cqrs server kind")
			return
		}
		g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/cqrs")

		if kind == "query" || kind == "full" {
			// query store
			//qsAsk := "query database kind:\n\t1) postgres\n\t2) mysql\n\t3) dgraph\n\t4) rgraph\nplease enter the number: "
			qsAsk := "query database kind:\n\t1) postgres\n\t2) mysql\nplease enter the number: "
			qsNo, qsNoErr := ask(qsAsk)
			if qsNoErr != nil {
				err = qsNoErr
				return
			}
			switch qsNo {
			case "1":
				g.Module.Requires = append(g.Module.Requires, "github.com/lib/pq")
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/databases/sql")
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/databases/postgres")
			case "2":
				g.Module.Requires = append(g.Module.Requires, "github.com/go-sql-driver/mysql")
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/databases/sql")
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/databases/mysql")
			case "3":
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/databases/dgraph")
			case "4":
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/databases/rgraph")
			default:
				err = fmt.Errorf("fnc: please choose in list query database kind")
				return
			}
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/cqrs")
		}
		g.Settings.Dependencies = append(g.Settings.Dependencies, &model.Dependency{
			Name: "cqrs",
			Kind: kind,
		})
	}

	// dep >>>
	// auth
	needAuthorizations, needAuthorizationsErr := ask("does project need authorizations (y/n): ")
	if needAuthorizationsErr != nil {
		err = needAuthorizationsErr
		return
	}
	if strings.ToLower(needAuthorizations) == "y" {
		encodingAsk := "authorizations encodings:\n\t1) default\n\t2) jwt\nplease enter the number: "
		authorizationsEncodingNo, authorizationsEncodingNoErr := ask(encodingAsk)
		if authorizationsEncodingNoErr != nil {
			err = authorizationsEncodingNoErr
			return
		}
		authorizationsEncoding := ""
		switch authorizationsEncodingNo {
		case "1":
			authorizationsEncoding = "default"
		case "2":
			authorizationsEncoding = "jwt"
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/authorizations/encoding/jwt")
		default:
			err = fmt.Errorf("fnc: please choose in list authorizations encoding")
			return
		}
		storeAsk := "authorizations stores:\n\t1) discard\n\t2) redis\n\t3) postgres\n\t4) mysql\nplease enter the number: "
		//storeAsk := "authorizations stores:\n\t1) discard\n\t2) redis\n\t3) postgres\n\t4) mysql\n\t5) dgraph\n\t6) rgraph\nplease enter the number: "
		authorizationsStoreNo, authorizationsStoreNoErr := ask(storeAsk)
		if authorizationsStoreNoErr != nil {
			err = authorizationsStoreNoErr
			return
		}
		authorizationsStore := ""
		switch authorizationsStoreNo {
		case "1":
			authorizationsStore = "discard"
		case "2":
			authorizationsStore = "redis"
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/authorizations/store/redis")
		case "3":
			authorizationsStore = "postgres"
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/authorizations/store/postgres")
		case "4":
			authorizationsStore = "mysql"
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/authorizations/store/mysql")
		case "5":
			authorizationsStore = "dgraph"
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/authorizations/store/dgraph")
		case "6":
			authorizationsStore = "rgraph"
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/authorizations/store/rgraph")
		default:
			err = fmt.Errorf("fnc: please choose in list authorizations store")
			return
		}
		g.Settings.Dependencies = append(g.Settings.Dependencies, &model.Dependency{
			Name: "authorizations",
			Kind: fmt.Sprintf("%s:%s", authorizationsEncoding, authorizationsStore),
		})
	}
	// permissions
	needPermissions, needPermissionsErr := ask("does project need permissions (y/n): ")
	if needPermissionsErr != nil {
		err = needPermissionsErr
		return
	}
	if strings.ToLower(needPermissions) == "y" {
		//storeAsk := "permissions stores:\n\t1) postgres\n\t2) mysql\n\t3) dgraph\n\t4) rgraph\nplease enter the number: "
		storeAsk := "permissions stores:\n\t1) postgres\n\t2) mysql\nplease enter the number: "
		storeNo, storeNoErr := ask(storeAsk)
		if storeNoErr != nil {
			err = storeNoErr
			return
		}
		store := ""
		switch storeNo {
		case "1":
			store = "postgres"
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/permissions/store/postgres")
		case "2":
			store = "mysql"
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/permissions/store/mysql")
		case "3":
			store = "dgraph"
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/permissions/store/dgraph")
		case "4":
			store = "rgraph"
			g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/permissions/store/rgraph")
		default:
			err = fmt.Errorf("fnc: please choose in list permissions policy store")
			return
		}
		g.Settings.Dependencies = append(g.Settings.Dependencies, &model.Dependency{
			Name: "permissions",
			Kind: store,
		})
	}
	if needDDD != "y" {
		// sql
		needSQL, needSQLErr := ask("does project need sql (y/n): ")
		if needSQLErr != nil {
			err = needSQLErr
			return
		}
		if strings.ToLower(needSQL) == "y" {
			sqlAsk := "sql type:\n\t1) postgres\n\t2) mysql\nplease enter the number: "
			sqlNo, sqlNoErr := ask(sqlAsk)
			if sqlNoErr != nil {
				err = sqlNoErr
				return
			}
			sqlType := ""
			switch sqlNo {
			case "1":
				sqlType = "postgres"
				g.Module.Requires = append(g.Module.Requires, "github.com/lib/pq")
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/databases/sql")
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/databases/postgres")
			case "2":
				sqlType = "mysql"
				g.Module.Requires = append(g.Module.Requires, "github.com/go-sql-driver/mysql")
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/databases/sql")
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/databases/mysql")
			default:
				err = fmt.Errorf("fnc: please choose in list sql type")
				return
			}
			g.Settings.Dependencies = append(g.Settings.Dependencies, &model.Dependency{
				Name: "sql",
				Kind: sqlType,
			})
		}
		// message queue
		needMQ, needMQErr := ask("does project need message queue (y/n): ")
		if needMQErr != nil {
			err = needMQErr
			return
		}
		if strings.ToLower(needMQ) == "y" {
			mqAsk := "message queue type:\n\t1) rabbitMQ\n\t2) kafka\n\t3) rocketMQ\n\t4) nats\nplease enter the number: "
			mqNo, mqNoErr := ask(mqAsk)
			if mqNoErr != nil {
				err = mqNoErr
				return
			}
			mqType := ""
			switch mqNo {
			case "1":
				mqType = "rabbitMQ"
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/message-queues/rabbit")
			case "2":
				mqType = "kafka"
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/message-queues/kafka")
			case "3":
				mqType = "rocketMQ"
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/message-queues/rocket")
			case "4":
				mqType = "nats"
				g.Module.Requires = append(g.Module.Requires, "github.com/aacfactory/fns-contrib/message-queues/nats")
			default:
				err = fmt.Errorf("fnc: please choose in list message queue type")
				return
			}
			g.Settings.Dependencies = append(g.Settings.Dependencies, &model.Dependency{
				Name: "mq",
				Kind: mqType,
			})
		}
	}
	// dep <<<
	return
}

func ask(question string) (answer string, err error) {
	in := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	line, lineErr := in.ReadString('\n')
	if lineErr != nil {
		err = fmt.Errorf("fnc: create project failed, ask failed, %v", lineErr)
		err = lineErr
		return
	}
	answer = strings.TrimSpace(line)
	return
}
