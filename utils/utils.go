package utils

import (
	"encoding/json"
)

func Argvs(source []byte, argv *[]string) error {
	var argvBytes []byte
	argvBytes = append(argvBytes, []byte("[")...)
	argvBytes = append(argvBytes, source...)
	argvBytes = append(argvBytes, []byte("]")...)
	return json.Unmarshal(argvBytes, argv)
}
