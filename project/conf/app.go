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
	"encoding/json"
)

type Config struct {
	Name           string     `yaml:"name,omitempty"`
	Runtime        *Runtime   `yaml:"runtime,omitempty"`
	Log            *LogConfig `yaml:"log,omitempty"`
	OAS            *OAS       `yaml:"oas,omitempty"`
	Server         *Server    `yaml:"server,omitempty"`
	Cluster        *Cluster   `yaml:"cluster,omitempty"`
	Authorizations *Auth      `yaml:"authorizations,omitempty"`
	SQL            *SQL       `yaml:"sql,omitempty"`
}

type Cluster struct {
	DevMode           bool                   `yaml:"devMode"`
	NodesProxyAddress string                 `yaml:"nodesProxyAddress"`
	Kind              string                 `yaml:"kind,omitempty"`
	Client            ClusterClient          `yaml:"client,omitempty"`
	Options           map[string]interface{} `yaml:"options,omitempty"`
}

type ClusterClient struct {
	MaxIdleConnSeconds    int `yaml:"maxIdleConnSeconds,omitempty"`
	MaxConnsPerHost       int `yaml:"maxConnsPerHost,omitempty"`
	MaxIdleConnsPerHost   int `yaml:"maxIdleConnsPerHost,omitempty"`
	RequestTimeoutSeconds int `yaml:"requestTimeoutSeconds,omitempty"`
}

type LogConfig struct {
	Level     string `yaml:"level,omitempty"`
	Formatter string `yaml:"formatter,omitempty"`
	Color     bool   `yaml:"color,omitempty"`
}

type OASContact struct {
	Name  string `yaml:"name,omitempty"`
	Url   string `yaml:"url,omitempty"`
	Email string `yaml:"email,omitempty"`
}

type OASLicense struct {
	Name string `yaml:"name,omitempty"`
	Url  string `yaml:"url,omitempty"`
}

type OASServer struct {
	URL         string `yaml:"url,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type OAS struct {
	Title       string      `yaml:"title,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Terms       string      `yaml:"terms,omitempty"`
	Contact     *OASContact `yaml:"contact,omitempty"`
	License     *OASLicense `yaml:"license,omitempty"`
	Servers     []OASServer `yaml:"servers,omitempty"`
}

type Cors struct {
	AllowedOrigins   []string `yaml:"allowedOrigins,omitempty"`
	AllowedHeaders   []string `yaml:"allowedHeaders,omitempty"`
	ExposedHeaders   []string `yaml:"exposedHeaders,omitempty"`
	AllowCredentials bool     `yaml:"allowCredentials,omitempty"`
	MaxAge           int      `yaml:"maxAge,omitempty"`
}

type TLS struct {
	// Kind
	// ACME
	// SSC(SELF-SIGN-CERT)
	// DEFAULT
	Kind    string          `yaml:"kind,omitempty"`
	Options json.RawMessage `yaml:"options,omitempty"`
}

type Server struct {
	Port         int                        `yaml:"port,omitempty"`
	Cors         *Cors                      `yaml:"cors,omitempty"`
	TLS          *TLS                       `yaml:"tls,omitempty"`
	Options      json.RawMessage            `yaml:"options,omitempty"`
	Interceptors map[string]json.RawMessage `yaml:"interceptors,omitempty"`
}

type Runtime struct {
	MaxWorkers           int `yaml:"maxWorkers"`
	WorkerMaxIdleSeconds int `yaml:"workerMaxIdleSeconds"`
	HandleTimeoutSeconds int `yaml:"handleTimeoutSeconds"`
}

type SQL struct {
	Driver           string   `yaml:"driver"`
	MasterSlaverMode bool     `yaml:"masterSlaverMode"`
	DSN              []string `yaml:"dsn"`
	MaxIdles         int      `yaml:"maxIdles"`
	MaxOpens         int      `yaml:"maxOpens"`
	EnableDebugLog   bool     `yaml:"enableDebugLog"`
	GTMCleanUpSecond int      `yaml:"gtmCleanUpSecond"`
	Isolation        int      `yaml:"isolation"`
}
