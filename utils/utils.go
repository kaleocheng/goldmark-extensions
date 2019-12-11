package utils

import (
	"strings"
)

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func Argvs(source []byte) []string {
	return Map(strings.Split(string(source), ","), func(s string) string {
		s = strings.Trim(s, `" `)
		return s
	})
}
