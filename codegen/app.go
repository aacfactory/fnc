package codegen

import (
	"fmt"
	"path/filepath"
)

func Generate(path string, plugins []string) (err error) {

	projectPath, projectPathErr := filepath.Abs(path)
	if projectPathErr != nil {
		err = fmt.Errorf("get abs filepath of %s failed, %v", path, projectPathErr)
		return
	}

	project, projectErr := LoadProject(projectPath)
	if projectErr != nil {
		Log().Errorf("load project %s failed, %v", path, projectErr)
		return
	}

	fnFiles, loadFnErr := LoadFn(project)
	if loadFnErr != nil {
		Log().Errorf("load @Fn in project %s failed, %v", path, loadFnErr)
		return
	}


	for _, fnFile := range fnFiles {
		Log().Debugf("fn file %s", fnFile.Path)
		for _, fn := range fnFile.Functions {
			Log().Debugf("fn %s", fn.Name)

		}
	}

return
}
