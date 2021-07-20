package codegen

func Generate(path string, plugins []string) {

	project, projectErr := LoadProject(path)
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

}
