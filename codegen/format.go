package codegen

import (
	"go/format"
	"golang.org/x/tools/imports"
)

func GoFormat(src []byte) (result []byte, err error) {
	imported, importErr := imports.Process("", src, nil)
	if importErr != nil {
		err = importErr
		return
	}
	formatted, fmtErr := format.Source(imported)
	if fmtErr != nil {
		err = fmtErr
		return
	}
	result = formatted
	return
}
