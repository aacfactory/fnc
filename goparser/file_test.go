package goparser

import (
	"testing"
)

func TestParseFile(t *testing.T) {
	path := "D:\\studio\\workspace\\go\\src\\github.com\\aacfactory\\fnc\\goparser\\file_case_test.go"
	file, err := ParseFile(path)
	t.Log(err, file)
}
