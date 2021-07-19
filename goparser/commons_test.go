package goparser

import "testing"

func TestIsExported(t *testing.T)  {
	t.Log(IsExported("Abb"))
	t.Log(IsExported("abb"))
	t.Log(IsExported("_bb"))
}
