package Resource

import (
	_ "embed"
	"encoding/base64"
)

func Bs64ToBs(s string) []byte {
	b, _ := base64.StdEncoding.DecodeString(s)
	return b
}

//go:embed GoCode.txt
var GoCode []byte

//go:embed BuiltFunc.txt
var GoBuiltFuncCode []byte
