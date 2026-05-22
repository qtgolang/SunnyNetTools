//go:build windows

package Tools

import (
	"bytes"
	"changeme/Service/Config"
	"crypto/md5"
	"encoding/base64"
	"github.com/orcastor/fico"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type AddCustomTools struct {
	App *SunnyNet.Sunny
}

var toolsLock sync.Mutex

func (g *AddCustomTools) CustomToolsAdd(filePath string) string {
	var w bytes.Buffer
	err := fico.PE2ICO(&w, filePath, fico.Config{Format: "png", Width: 128, Height: 128})
	if err != nil {
		panic(err)
	}
	fileNameWithExt := filepath.Base(filePath)
	fileName := strings.TrimSuffix(fileNameWithExt, filepath.Ext(fileNameWithExt))
	var obj Config.ToolsInfo
	md := md5.Sum([]byte(time.Now().String()))
	obj.ID = hex.EncodeToString(md[:])
	obj.Name = fileName
	obj.Icon = base64.StdEncoding.EncodeToString(w.Bytes())
	obj.File = filePath
	Config.Config.ToolsList = append(Config.Config.ToolsList, obj)
	Config.AppList["Main"].EmitEvent("addTools", fileName, obj)
	Config.Config.Save()
	return ""
}
func (g *AddCustomTools) CustomToolsDel(id string) bool {
	toolsLock.Lock()
	defer toolsLock.Unlock()
	for i, obj := range Config.Config.ToolsList {
		if obj.ID == id {
			Config.Config.ToolsList = append(Config.Config.ToolsList[:i], Config.Config.ToolsList[i+1:]...)
			Config.Config.Save()
			return true
		}
	}
	return false
}

func (g *AddCustomTools) CustomToolsList() []Config.ToolsInfo {
	toolsLock.Lock()
	defer toolsLock.Unlock()
	return Config.Config.ToolsList
}
func (g *AddCustomTools) SaveCustomTools(objInfo Config.ToolsInfo) bool {
	toolsLock.Lock()
	defer toolsLock.Unlock()
	for v, obj := range Config.Config.ToolsList {
		if obj.ID == objInfo.ID {
			//obj.Icon = objInfo.Icon
			obj.Args = objInfo.Args
			obj.File = objInfo.File
			obj.Name = objInfo.Name
			Config.Config.ToolsList[v] = obj
			Config.Config.Save()
			return true
		}
	}
	return false
}
func (g *AddCustomTools) ExecCustomTools(id string) string {
	toolsLock.Lock()
	defer toolsLock.Unlock()
	for _, obj := range Config.Config.ToolsList {
		if obj.ID == id {
			if Exists(obj.File) {
				args := strings.Fields(obj.Args)
				cmd := exec.Command(obj.File, args...)
				_ = cmd.Start()
				return ""
			}
			return "文件路径不存在!"
		}
	}
	return "找不到此选项"
}
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true

}
