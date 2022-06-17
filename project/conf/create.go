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
	var jwtKeyPair *rsa.PrivateKey
	app := Config{
		Name:    g.Settings.Name,
		Service: nil,
		Log: &LogConfig{
			Level:     "info",
			Formatter: "console",
			Color:     true,
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
		Service: nil,
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
	localConf, _ := yaml.Marshal(local)
	dev := Config{
		Service: nil,
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
	devConf, _ := yaml.Marshal(dev)
	qa := Config{
		Service: nil,
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
	qaConf, _ := yaml.Marshal(qa)
	prod := Config{
		Service: nil,
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
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(jwtKeyPair),
		})
		err = ioutil.WriteFile(filepath.Join(jwtDir, "pri.pem"), priPEM, 0600)
		if err != nil {
			err = fmt.Errorf("fnc: create project failed at write jwt rsa private key file, %v", err)
			return
		}
		pubPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "BEGIN PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&jwtKeyPair.PublicKey),
		})
		err = ioutil.WriteFile(filepath.Join(jwtDir, "pub.pem"), pubPEM, 0600)
		if err != nil {
			err = fmt.Errorf("fnc: create project failed at write jwt rsa public key file, %v", err)
			return
		}
	}
	return
}
