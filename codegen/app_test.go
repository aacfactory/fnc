package codegen_test

import (
	"github.com/aacfactory/fnc/codegen"
	"testing"
)

func TestGenerate(t *testing.T) {

	path := `D:\studio\workspace\go\src\github.com\aacfactory\example`

	codegen.Generate(path , nil)

}
