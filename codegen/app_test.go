package codegen_test

import (
	"fmt"
	"github.com/aacfactory/fnc/codegen"
	"go/format"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"testing"
)

func TestGenerate(t *testing.T) {

	path := `D:\studio\workspace\go\src\github.com\aacfactory\example`

	codegen.Generate(path, nil)

}

func TestFMT(t *testing.T) {

	file := "D:\\studio\\workspace\\go\\src\\github.com\\aacfactory\\example\\main.go"

	c, _ := ioutil.ReadFile(file)
	b, _ := imports.Process("", c, nil)
	d, _ := format.Source(b)
	fmt.Println(string(d))
}
