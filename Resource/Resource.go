package Resource

import (
	"embed"
	_ "embed"
	"encoding/base64"
	"encoding/json"
)

func Bs64ToBs(s string) []byte {
	b, _ := base64.StdEncoding.DecodeString(s)
	return b
}

//go:embed GoCode.txt
var GoCode []byte

//go:embed BuiltFunc.txt
var GoBuiltFuncCode []byte

//go:embed ScriptTemplate/_aes-des.goc
//go:embed ScriptTemplate/_forward.goc
//go:embed ScriptTemplate/_header.goc
//go:embed ScriptTemplate/_quickResponse.goc
//go:embed ScriptTemplate/_ReplaceUrlParameter.goc
//go:embed ScriptTemplate/_response.goc
//go:embed ScriptTemplate/_responseFile.goc
//go:embed ScriptTemplate/_rsa.goc
//go:embed ScriptTemplate/_pid.goc
//go:embed ScriptTemplate/annotation.txt
//go:embed ScriptTemplate/info.json
var scriptTemplate embed.FS

var ScriptTemplate []*ScriptTemplateInfo
var ScriptAnnotation string

type ScriptTemplateInfo struct {
	Name    string
	Explain string
	Data    string
}
type TemplateInfo struct {
	Name     string `json:"Name"`
	Explain  string `json:"Explain"`
	FileName string `json:"FileName"`
}

func init() {
	annotation, _ := scriptTemplate.ReadFile("ScriptTemplate/annotation.txt")
	ScriptAnnotation = string(annotation)

	scriptInfo, _ := scriptTemplate.ReadFile("ScriptTemplate/info.json")
	var info []TemplateInfo
	_ = json.Unmarshal(scriptInfo, &info)
	for _, entry := range info {
		infoCode := &ScriptTemplateInfo{}
		infoCode.Name = entry.Name
		infoCode.Explain = entry.Explain
		fileData, err1 := scriptTemplate.ReadFile("ScriptTemplate/" + entry.FileName)
		if err1 == nil {
			infoCode.Data = string(fileData) + ScriptAnnotation
			ScriptTemplate = append(ScriptTemplate, infoCode)
		}
	}
}
