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

package codes

import (
	"fmt"
	"github.com/aacfactory/gcg"
	"path/filepath"
)

func (svc *Service) generateFile() (err error) {

	fns := svc.Fns()
	if len(fns) == 0 {
		return
	}

	file := gcg.NewFile(svc.Package)

	// imports
	for _, _import := range svc.Imports {
		if _import.Alias != "" {
			file.AddImport(gcg.NewPackageWithAlias(_import.Path, _import.Alias))
		} else {
			file.AddImport(gcg.NewPackage(_import.Path))
		}
	}
	var code gcg.Code = nil
	var codeErr error = nil

	code, codeErr = svc.generateFileConst(fns)
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceProxy(fns)
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileService()
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceNamespace()
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceInternal()
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceBuild()
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceDescription()
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceHandle(fns)
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceClose()
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceInvokers(fns)
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	// write
	filename := filepath.Join(svc.DirPath, "fns_service_gen.go")
	w := gcg.FileRender(filename, true)
	wErr := file.Render(w)
	if wErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), wErr)
		return
	}
	closeFileErr := w.Close()
	if closeFileErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), closeFileErr)
		return
	}

	return
}

func (svc *Service) generateFileConst(fns []*Fn) (code gcg.Code, err error) {
	svc.Fns()
	return
}

func (svc *Service) generateFileService() (code gcg.Code, err error) {

	return
}

func (svc *Service) generateFileServiceNamespace() (code gcg.Code, err error) {

	return
}

func (svc *Service) generateFileServiceInternal() (code gcg.Code, err error) {

	return
}

func (svc *Service) generateFileServiceBuild() (code gcg.Code, err error) {

	return
}

func (svc *Service) generateFileServiceDescription() (code gcg.Code, err error) {

	return
}

func (svc *Service) generateFileServiceHandle(fns []*Fn) (code gcg.Code, err error) {

	return
}

func (svc *Service) generateFileServiceClose() (code gcg.Code, err error) {

	return
}

func (svc *Service) generateFileServiceInvokers(fns []*Fn) (code gcg.Code, err error) {

	return
}

func (svc *Service) generateFileServiceProxy(fns []*Fn) (code gcg.Code, err error) {

	return
}
