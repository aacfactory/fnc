package goparser

import (
	"strings"
)

func ParseStructTags(v string) (tags []FieldTag, ok bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return
	}
	idx := strings.IndexByte(v, ':')
	if idx > 0 {
		name := v[0:idx]
		sub := v[idx+1:]
		l := strings.IndexByte(sub, '"')
		if l < 0 {
			return
		}
		r := strings.IndexByte(sub[l+1:], '"')
		if r < 0 {
			return
		}
		values := strings.Split(sub[l+1:r+1], ",")
		tag := FieldTag{
			Name:   name,
			Values: make([]string, 0, 1),
		}
		for _, value := range values {
			value = strings.TrimSpace(value)
			if value != "" {
				tag.Values = append(tag.Values, value)
			}
		}
		tags = append(tags, tag)
		if len(sub) > r+2 {
			subTags, subTagsOk := ParseStructTags(sub[r+2:])
			if subTagsOk {
				tags = append(tags, subTags...)
			}
		}
	}

	ok = true
	return
}
