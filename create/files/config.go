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

package files

import (
	"context"
	"github.com/aacfactory/errors"
	"github.com/aacfactory/forg/files"
	"github.com/goccy/go-yaml"
	"os"
	"path/filepath"
	"strings"
)

func NewConfigFile(kind string, dir string) (cf *ConfigFile, err error) {
	if !filepath.IsAbs(dir) {
		dir, err = filepath.Abs(dir)
		if err != nil {
			err = errors.Warning("forg: new config file failed").WithCause(err).WithMeta("dir", dir)
			return
		}
	}
	dir = filepath.Join(dir, "configs")
	if !files.ExistFile(dir) {
		mdErr := os.MkdirAll(dir, 0600)
		if mdErr != nil {
			err = errors.Warning("forg: new config file failed").WithCause(mdErr).WithMeta("dir", dir)
			return
		}
	}
	name := "fns.yaml"
	if kind != "" {
		kind = strings.TrimSpace(strings.ToLower(kind))
		switch kind {
		case "local":
			name = "fns-local.yaml"
			break
		case "dev":
			name = "fns-dev.yaml"
			break
		case "test":
			name = "fns-test.yaml"
			break
		case "prod":
			name = "fns-prod.yaml"
			break
		default:
			err = errors.Warning("").WithCause(errors.Warning("kind is invalid")).WithMeta("kind", kind)
		}
	}

	filename := filepath.ToSlash(filepath.Join(dir, name))
	cf = &ConfigFile{
		kind:     kind,
		filename: filename,
	}
	return
}

type ConfigFile struct {
	kind     string
	filename string
}

func (cf *ConfigFile) Name() (name string) {
	name = cf.filename
	return
}

func (cf *ConfigFile) Write(ctx context.Context) (err error) {
	config := Config{}
	switch cf.kind {
	case "local":
		config.Log = &LogConfig{
			Level:     "debug",
			Formatter: "console",
			Color:     true,
		}
		config.Runtime = &RuntimeConfig{
			MaxWorkers:                0,
			WorkerMaxIdleSeconds:      0,
			HandleTimeoutSeconds:      0,
			LocalSharedStoreCacheSize: "64MB",
			AutoMaxProcs:              AutoMaxProcsConfig{},
			SecretKey:                 "",
		}
		break
	case "dev":
		config.Log = &LogConfig{
			Level:     "info",
			Formatter: "json",
			Color:     false,
		}
		config.Runtime = &RuntimeConfig{
			MaxWorkers:                0,
			WorkerMaxIdleSeconds:      0,
			HandleTimeoutSeconds:      0,
			LocalSharedStoreCacheSize: "64MB",
			AutoMaxProcs: AutoMaxProcsConfig{
				Min: 2,
				Max: 0,
			},
			SecretKey: "",
		}
		break
	case "test":
		config.Log = &LogConfig{
			Level:     "warn",
			Formatter: "json",
			Color:     false,
		}
		config.Runtime = &RuntimeConfig{
			MaxWorkers:                0,
			WorkerMaxIdleSeconds:      0,
			HandleTimeoutSeconds:      0,
			LocalSharedStoreCacheSize: "256MB",
			AutoMaxProcs: AutoMaxProcsConfig{
				Min: 2,
				Max: 0,
			},
			SecretKey: "",
		}
		break
	case "prod":
		config.Log = &LogConfig{
			Level:     "error",
			Formatter: "json",
			Color:     false,
		}
		config.Runtime = &RuntimeConfig{
			MaxWorkers:                0,
			WorkerMaxIdleSeconds:      0,
			HandleTimeoutSeconds:      0,
			LocalSharedStoreCacheSize: "256MB",
			AutoMaxProcs: AutoMaxProcsConfig{
				Min: 4,
				Max: 0,
			},
			SecretKey: "",
		}
		break
	default:
		config.Http = &HttpConfig{
			Port: 18080,
		}
		break
	}
	p, encodeErr := yaml.Marshal(config)
	if encodeErr != nil {
		err = errors.Warning("forg: config file write failed").WithCause(encodeErr).WithMeta("filename", cf.filename)
		return
	}
	writeErr := os.WriteFile(cf.filename, p, 0600)
	if writeErr != nil {
		err = errors.Warning("forg: config file write failed").WithCause(writeErr).WithMeta("filename", cf.filename)
		return
	}
	return
}

type Config struct {
	Http    *HttpConfig    `json:"http" yaml:"http,omitempty"`
	Log     *LogConfig     `json:"log" yaml:"log,omitempty"`
	Runtime *RuntimeConfig `json:"runtime" yaml:"runtime,omitempty"`
}

type LogConfig struct {
	Level     string `json:"level" yaml:"level,omitempty"`
	Formatter string `json:"formatter" yaml:"formatter,omitempty"`
	Color     bool   `json:"color" yaml:"color,omitempty"`
}

type RuntimeConfig struct {
	MaxWorkers                int                `json:"maxWorkers" yaml:"maxWorkers,omitempty"`
	WorkerMaxIdleSeconds      int                `json:"workerMaxIdleSeconds" yaml:"workerMaxIdleSeconds,omitempty"`
	HandleTimeoutSeconds      int                `json:"handleTimeoutSeconds" yaml:"handleTimeoutSeconds,omitempty"`
	LocalSharedStoreCacheSize string             `json:"localSharedStoreCacheSize" yaml:"localSharedStoreCacheSize,omitempty"`
	AutoMaxProcs              AutoMaxProcsConfig `json:"autoMaxProcs" yaml:"autoMaxProcs,omitempty"`
	SecretKey                 string             `json:"secretKey" yaml:"secretKey,omitempty"`
}

type AutoMaxProcsConfig struct {
	Min int `json:"min" yaml:"min"`
	Max int `json:"max" yaml:"max"`
}

type HttpConfig struct {
	Port int `json:"port" yaml:"port,omitempty"`
}
