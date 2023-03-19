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

func NewHooksFile(dir string) (hooks *HooksFile, err error) {
	if !filepath.IsAbs(dir) {
		dir, err = filepath.Abs(dir)
		if err != nil {
			err = errors.Warning("forg: new hooks file failed").WithCause(err).WithMeta("dir", dir)
			return
		}
	}
	dir = filepath.Join(dir, "hooks")
	if !files.ExistFile(dir) {
		mdErr := os.MkdirAll(dir, 0600)
		if mdErr != nil {
			err = errors.Warning("forg: new hooks file failed").WithCause(mdErr).WithMeta("dir", dir)
			return
		}
	}
	hooks = &HooksFile{
		filename: filepath.ToSlash(filepath.Join(dir, "doc.go")),
	}
	return
}

type HooksFile struct {
	filename string
}

func (hooks *HooksFile) Name() (name string) {
	name = hooks.filename
	return
}

func (hooks *HooksFile) Write(ctx context.Context) (err error) {
	const (
		content = `// Package hooks
// read https://github.com/aacfactory/fns/blob/main/docs/hooks.md for more details.
package hooks`
	)
	writeErr := os.WriteFile(hooks.filename, []byte(content), 0600)
	if writeErr != nil {
		err = errors.Warning("forg: hooks file write failed").WithCause(writeErr).WithMeta("filename", hooks.filename)
		return
	}
	return
}
