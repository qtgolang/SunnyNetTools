package main

import "C"
import (
	"encoding/json"
	"github.com/qtgolang/SunnyNet/src/protobuf"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
)

func _PbToJson(data []byte, skip int) string {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	var msg protobuf.Message
	msg.Unmarshal(data[skip:])
	b, e := json.Marshal(msg)
	if e != nil {
		return ""
	}
	PJson, _ := protobuf.ParseJson(string(b), "")
	s, _ := json.MarshalIndent(PJson, "", "\t")
	ss := string(s)
	ss = strings.ReplaceAll(ss, "\n", "\r\n")
	return ss
}
func Utf8ToGBK(src []byte) []byte {
	buf := make([]byte, len(src)*2)
	n, _, err := simplifiedchinese.GBK.NewEncoder().Transform(buf, src, true)
	if err != nil {
		return nil
	}
	return buf[:n]
}
