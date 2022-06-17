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
	Name           string
	Service        *Service
	Log            *LogConfig
	OAS            *OAS
	Server         *Server
	Cluster        *Cluster
	Authorizations *Auth
}

type Cluster struct {
	Kind    string
	Client  ClusterClient
	Options json.RawMessage
}

type ClusterClient struct {
	MaxIdleConnSeconds    int
	MaxConnsPerHost       int
	MaxIdleConnsPerHost   int
	RequestTimeoutSeconds int
}

type LogConfig struct {
	Level     string
	Formatter string
	Color     bool
}

type OASContact struct {
	Name  string
	Url   string
	Email string
}

type OASLicense struct {
	Name string
	Url  string
}

type OASServer struct {
	URL         string
	Description string
}

type OAS struct {
	Title       string
	Description string
	Terms       string
	Contact     *OASContact
	License     *OASLicense
	Servers     []OASServer
}

type Cors struct {
	AllowedOrigins   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

type TLS struct {
	// Kind
	// ACME
	// SSC(SELF-SIGN-CERT)
	// DEFAULT
	Kind    string
	Options json.RawMessage
}

type Server struct {
	Port         int
	Cors         *Cors
	TLS          *TLS
	Options      json.RawMessage
	Interceptors map[string]json.RawMessage
}

type Service struct {
	MaxWorkers           int
	WorkerMaxIdleSeconds int
	HandleTimeoutSeconds int
}
