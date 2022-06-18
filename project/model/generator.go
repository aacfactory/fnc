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

package model

type Generator struct {
	Path     string
	Module   Module
	Settings Settings
}

type Module struct {
	Name     string
	Go       string
	Requires []string
}

type Contact struct {
	Name  string
	Url   string
	Email string
}

type License struct {
	Name string
	Url  string
}

type Server struct {
	URL         string
	Description string
}

type OAS struct {
	Title       string
	Description string
	Terms       string
	Contact     *Contact
	License     *License
	Servers     []Server
}

type Dependency struct {
	Name string
	Kind string
}

type Settings struct {
	Name         string
	Oas          OAS
	Http         string
	Cluster      string
	Dependencies []*Dependency
}

func (s Settings) FindDependency(name string) (dep *Dependency, has bool) {
	for _, dependency := range s.Dependencies {
		if dependency.Name == name {
			dep = dependency
			has = true
			return
		}
	}
	return
}
