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

package conf

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/aacfactory/fnc/project/model"
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Create(g model.Generator) (err error) {
	if g.Settings.Oas.Title == "" {
		g.Settings.Oas.Title = g.Settings.Name
	}
	var jwtKeyPair *rsa.PrivateKey
	app := Config{
		Name: g.Settings.Name,
		Runtime: &Runtime{
			MaxWorkers:           256 * 1024,
			WorkerMaxIdleSeconds: 10,
			HandleTimeoutSeconds: 10,
		},
		Log: &LogConfig{
			Level:     "info",
			Formatter: "console",
			Color:     false,
		},
		OAS: &OAS{
			Title:       g.Settings.Oas.Title,
			Description: g.Settings.Oas.Description,
			Terms:       g.Settings.Oas.Terms,
			Contact: func() *OASContact {
				if g.Settings.Oas.Contact == nil {
					return nil
				}
				return &OASContact{
					Name:  g.Settings.Oas.Contact.Name,
					Url:   g.Settings.Oas.Contact.Url,
					Email: g.Settings.Oas.Contact.Email,
				}
			}(),
			License: func() *OASLicense {
				if g.Settings.Oas.License == nil {
					return nil
				}
				return &OASLicense{
					Name: g.Settings.Oas.License.Name,
					Url:  g.Settings.Oas.License.Url,
				}
			}(),
			Servers: func() []OASServer {
				if g.Settings.Oas.Servers == nil || len(g.Settings.Oas.Servers) == 0 {
					return nil
				}
				ss := make([]OASServer, 0, 1)
				for _, server := range g.Settings.Oas.Servers {
					ss = append(ss, OASServer{
						URL:         server.URL,
						Description: server.Description,
					})
				}
				return ss
			}(),
		},
		Server: &Server{
			Port:         8800,
			Cors:         nil,
			TLS:          nil,
			Options:      nil,
			Interceptors: nil,
		},
		Cluster: nil,
		Authorizations: func() *Auth {
			dep, has := g.Settings.FindDependency("authorizations")
			if !has {
				return nil
			}
			es := strings.Split(dep.Kind, ":")
			encoding := es[0]
			if encoding == "jwt" {
				jwtKeyPair, _ = rsa.GenerateKey(rand.Reader, 4096)
				return &Auth{
					Encoding: &AuthJWTEncoding{
						Method:      "RS512",
						SK:          "",
						PublicKey:   "./config/jwt/pub.pem",
						PrivateKey:  "./config/jwt/pri.pem",
						Issuer:      "fns",
						Audience:    nil,
						Expirations: "720h0m0s",
					},
				}
			}
			return nil
		}(),
	}
	appConf, _ := yaml.Marshal(app)
	local := Config{
		Runtime: nil,
		Log: &LogConfig{
			Level:     "debug",
			Formatter: "console",
			Color:     true,
		},
		Server: &Server{
			Port:         8800,
			Cors:         nil,
			TLS:          nil,
			Options:      nil,
			Interceptors: nil,
		},
		Cluster: nil,
	}
	cluster, hasCluster := g.Settings.FindDependency("cluster")
	if hasCluster {
		local.Cluster = &Cluster{
			DevMode:           true,
			NodesProxyAddress: "proxyHost",
			Kind:              cluster.Kind,
			Client: ClusterClient{
				MaxIdleConnSeconds:    60,
				MaxConnsPerHost:       0,
				MaxIdleConnsPerHost:   0,
				RequestTimeoutSeconds: 10,
			},
			Options: make(map[string]interface{}),
		}
		if cluster.Kind == "members" {
			local.Cluster.Options[cluster.Kind] = []string{"ip:port"}
		}
		if cluster.Kind == "swarm" {
			local.Cluster.Options[cluster.Kind] = map[string]interface{}{
				"fromEnv": false,
				"host":    "master ip:port",
				"certDir": "./config/swarm",
				"labels":  []string{"FNS-SERVICE"},
			}
			_ = os.MkdirAll(filepath.Join(g.Path, "config", "swarm"), 0666)
		}
		if cluster.Kind == "kubernetes" {
			local.Cluster.Options[cluster.Kind] = map[string]interface{}{
				"inCluster":      false,
				"kubeConfigPath": "./config/kube",
				"namespace":      "fns",
				"labels":         []string{"FNS-SERVICE"},
			}
			_ = os.MkdirAll(filepath.Join(g.Path, "config", "kube"), 0666)
		}
	}
	sql, hasSQL := g.Settings.FindDependency("sql")
	if hasSQL {
		driver := ""
		if sql.Kind == "postgres" {
			driver = "postgres"
		} else if sql.Kind == "mysql" {
			driver = "mysql"
		}
		local.SQL = &SQL{
			Driver:           driver,
			MasterSlaverMode: false,
			DSN:              []string{"dsn"},
			MaxIdles:         2,
			MaxOpens:         4,
			EnableDebugLog:   false,
			GTMCleanUpSecond: 60,
			Isolation:        2,
		}
	}
	localConf, _ := yaml.Marshal(local)
	dev := Config{
		Runtime: nil,
		Log: &LogConfig{
			Level:     "debug",
			Formatter: "json",
			Color:     false,
		},
		Server: &Server{
			Port:         8800,
			Cors:         nil,
			TLS:          nil,
			Options:      nil,
			Interceptors: nil,
		},
		Cluster: nil,
	}
	if hasCluster {
		dev.Cluster = &Cluster{
			DevMode:           true,
			NodesProxyAddress: "proxyHost",
			Kind:              cluster.Kind,
			Client: ClusterClient{
				MaxIdleConnSeconds:    60,
				MaxConnsPerHost:       0,
				MaxIdleConnsPerHost:   0,
				RequestTimeoutSeconds: 10,
			},
			Options: make(map[string]interface{}),
		}
		if cluster.Kind == "members" {
			dev.Cluster.Options[cluster.Kind] = []string{"ip:port"}
		}
		if cluster.Kind == "swarm" {
			dev.Cluster.Options[cluster.Kind] = map[string]interface{}{
				"fromEnv": false,
				"host":    "master ip:port",
				"certDir": "./config/swarm",
				"labels":  []string{"FNS-SERVICE"},
			}
		}
		if cluster.Kind == "kubernetes" {
			dev.Cluster.Options[cluster.Kind] = map[string]interface{}{
				"inCluster":      true,
				"kubeConfigPath": "~/.kube",
				"namespace":      "fns",
				"labels":         []string{"FNS-SERVICE"},
			}
		}
	}
	if hasSQL {
		driver := ""
		if sql.Kind == "postgres" {
			driver = "postgres"
		} else if sql.Kind == "mysql" {
			driver = "mysql"
		}
		dev.SQL = &SQL{
			Driver:           driver,
			MasterSlaverMode: false,
			DSN:              []string{"dsn"},
			MaxIdles:         8,
			MaxOpens:         64,
			EnableDebugLog:   false,
			GTMCleanUpSecond: 60,
			Isolation:        2,
		}
	}
	devConf, _ := yaml.Marshal(dev)
	qa := Config{
		Runtime: nil,
		Log: &LogConfig{
			Level:     "warn",
			Formatter: "json",
			Color:     false,
		},
		Server: &Server{
			Port:         8800,
			Cors:         nil,
			TLS:          nil,
			Options:      nil,
			Interceptors: nil,
		},
		Cluster: nil,
	}
	if hasCluster {
		qa.Cluster = &Cluster{
			DevMode:           false,
			NodesProxyAddress: "",
			Kind:              cluster.Kind,
			Client: ClusterClient{
				MaxIdleConnSeconds:    60,
				MaxConnsPerHost:       0,
				MaxIdleConnsPerHost:   0,
				RequestTimeoutSeconds: 10,
			},
			Options: make(map[string]interface{}),
		}
		if cluster.Kind == "members" {
			qa.Cluster.Options[cluster.Kind] = []string{"ip:port"}
		}
		if cluster.Kind == "swarm" {
			qa.Cluster.Options[cluster.Kind] = map[string]interface{}{
				"fromEnv": false,
				"host":    "master ip:port",
				"certDir": "./config/swarm",
				"labels":  []string{"FNS-SERVICE"},
			}
		}
		if cluster.Kind == "kubernetes" {
			qa.Cluster.Options[cluster.Kind] = map[string]interface{}{
				"inCluster":      true,
				"kubeConfigPath": "~/.kube",
				"namespace":      "fns",
				"labels":         []string{"FNS-SERVICE"},
			}
		}
	}
	if hasSQL {
		driver := ""
		if sql.Kind == "postgres" {
			driver = "postgres"
		} else if sql.Kind == "mysql" {
			driver = "mysql"
		}
		qa.SQL = &SQL{
			Driver:           driver,
			MasterSlaverMode: false,
			DSN:              []string{"dsn"},
			MaxIdles:         8,
			MaxOpens:         64,
			EnableDebugLog:   false,
			GTMCleanUpSecond: 60,
			Isolation:        2,
		}
	}
	qaConf, _ := yaml.Marshal(qa)
	prod := Config{
		Runtime: nil,
		Log: &LogConfig{
			Level:     "error",
			Formatter: "json",
			Color:     false,
		},
		Server: &Server{
			Port:         8800,
			Cors:         nil,
			TLS:          nil,
			Options:      nil,
			Interceptors: nil,
		},
		Cluster: nil,
	}
	if hasCluster {
		prod.Cluster = &Cluster{
			DevMode:           false,
			NodesProxyAddress: "",
			Kind:              cluster.Kind,
			Client: ClusterClient{
				MaxIdleConnSeconds:    60,
				MaxConnsPerHost:       0,
				MaxIdleConnsPerHost:   0,
				RequestTimeoutSeconds: 10,
			},
			Options: make(map[string]interface{}),
		}
		if cluster.Kind == "members" {
			prod.Cluster.Options[cluster.Kind] = []string{"ip:port"}
		}
		if cluster.Kind == "swarm" {
			prod.Cluster.Options[cluster.Kind] = map[string]interface{}{
				"fromEnv": false,
				"host":    "master ip:port",
				"certDir": "./config/swarm",
				"labels":  []string{"FNS-SERVICE"},
			}
		}
		if cluster.Kind == "kubernetes" {
			prod.Cluster.Options[cluster.Kind] = map[string]interface{}{
				"inCluster":      true,
				"kubeConfigPath": "~/.kube",
				"namespace":      "fns",
				"labels":         []string{"FNS-SERVICE"},
			}
		}
	}
	if hasSQL {
		driver := ""
		if sql.Kind == "postgres" {
			driver = "postgres"
		} else if sql.Kind == "mysql" {
			driver = "mysql"
		}
		prod.SQL = &SQL{
			Driver:           driver,
			MasterSlaverMode: false,
			DSN:              []string{"dsn"},
			MaxIdles:         8,
			MaxOpens:         64,
			EnableDebugLog:   false,
			GTMCleanUpSecond: 60,
			Isolation:        2,
		}
	}
	prodConf, _ := yaml.Marshal(prod)

	configDir := filepath.Join(g.Path, "config")
	mdErr := os.MkdirAll(configDir, 600)
	if mdErr != nil {
		err = fmt.Errorf("fnc: create project failed at config mkdir, %v", mdErr)
		return
	}
	err = ioutil.WriteFile(filepath.Join(configDir, "fns.yaml"), appConf, 0600)
	if err != nil {
		err = fmt.Errorf("fnc: create project failed at write config file, %v", err)
		return
	}
	err = ioutil.WriteFile(filepath.Join(configDir, "fns-local.yaml"), localConf, 0600)
	if err != nil {
		err = fmt.Errorf("fnc: create project failed at write config file, %v", err)
		return
	}
	err = ioutil.WriteFile(filepath.Join(configDir, "fns-dev.yaml"), devConf, 0600)
	if err != nil {
		err = fmt.Errorf("fnc: create project failed at write config file, %v", err)
		return
	}
	err = ioutil.WriteFile(filepath.Join(configDir, "fns-qa.yaml"), qaConf, 0600)
	if err != nil {
		err = fmt.Errorf("fnc: create project failed at write config file, %v", err)
		return
	}
	err = ioutil.WriteFile(filepath.Join(configDir, "fns-prod.yaml"), prodConf, 0600)
	if err != nil {
		err = fmt.Errorf("fnc: create project failed at write config file, %v", err)
		return
	}
	if jwtKeyPair != nil {
		jwtDir := filepath.Join(configDir, "jwt")
		err = os.MkdirAll(jwtDir, 600)
		if err != nil {
			err = fmt.Errorf("fnc: create project failed at jwt config mkdir, %v", err)
			return
		}
		priPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(jwtKeyPair),
		})
		err = ioutil.WriteFile(filepath.Join(jwtDir, "pri.pem"), priPEM, 0600)
		if err != nil {
			err = fmt.Errorf("fnc: create project failed at write jwt rsa private key file, %v", err)
			return
		}
		pub := jwtKeyPair.Public().(*rsa.PublicKey)
		pubPem, pubErr := x509.MarshalPKIXPublicKey(pub)
		if pubErr != nil {
			err = fmt.Errorf("fnc: create project failed at write jwt rsa publich key file, %v", pubErr)
			return
		}
		pubPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubPem,
		})
		err = ioutil.WriteFile(filepath.Join(jwtDir, "pub.pem"), pubPEM, 0600)
		if err != nil {
			err = fmt.Errorf("fnc: create project failed at write jwt rsa public key file, %v", err)
			return
		}
	}
	return
}
