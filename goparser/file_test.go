package goparser

import (
	"testing"
)

func TestParseFile(t *testing.T) {
	EnableLog(true)
	path := "D:\\studio\\workspace\\go\\src\\github.com\\aacfactory\\fnc\\goparser\\file_case_test.go"
	file, err := ParseFile(path)
	t.Log(err, file)
}
