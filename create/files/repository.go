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
	"os"
	"path/filepath"
)

func NewRepositoryFile(dir string) (hooks *RepositoryFile, err error) {
	if !filepath.IsAbs(dir) {
		dir, err = filepath.Abs(dir)
		if err != nil {
			err = errors.Warning("forg: new repositories file failed").WithCause(err).WithMeta("dir", dir)
			return
		}
	}
	dir = filepath.Join(dir, "repositories")
	if !files.ExistFile(dir) {
		mdErr := os.MkdirAll(dir, 0600)
		if mdErr != nil {
			err = errors.Warning("forg: new repositories file failed").WithCause(mdErr).WithMeta("dir", dir)
			return
		}
	}
	hooks = &RepositoryFile{
		filename: filepath.ToSlash(filepath.Join(dir, "doc.go")),
	}
	return
}

type RepositoryFile struct {
	filename string
}

func (f *RepositoryFile) Name() (name string) {
	name = f.filename
	return
}

func (f *RepositoryFile) Write(ctx context.Context) (err error) {
	const (
		content = `// Package repositories
// read https://github.com/aacfactory/fns-contrib/tree/main/databases/sql for more details.
package hooks`
	)
	writeErr := os.WriteFile(f.filename, []byte(content), 0600)
	if writeErr != nil {
		err = errors.Warning("forg: repositories file write failed").WithCause(writeErr).WithMeta("filename", f.filename)
		return
	}
	return
}
