package codegen_test

import (
	"github.com/aacfactory/fnc/codegen"
	"testing"
)

func TestLoadProject(t *testing.T) {
	path := `D:\studio\workspace\go\src\github.com\aacfactory\example`
	project, projectErr := codegen.LoadProject(path)
	if projectErr != nil {
		t.Errorf("load project failed, %v", projectErr)
		return
	}

	t.Log("path:", project.Path)
	t.Log("mod:", project.Module.Name, project.Module.Version, project.Module.Requires)
	t.Log("files:")

}
