package goparser

import "testing"

func TestParseStructTags(t *testing.T) {
	v := `json:"id,omitempty" xml:"id" fn:"-" `
	tags, ok := ParseStructTags(v)
	t.Log(ok, tags)
}
