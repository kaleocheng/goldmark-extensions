package utils

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
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

func GetJSON(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
