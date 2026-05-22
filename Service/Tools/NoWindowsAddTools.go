//go:build !windows

package Tools

import (
	"changeme/Service/Config"
	"github.com/qtgolang/SunnyNet/SunnyNet"
)

type AddCustomTools struct {
	App     *SunnyNet.Sunny
	AddFunc func(name string, obj Config.ToolsInfo)
}

func (g *AddCustomTools) CustomToolsList() []Config.ToolsInfo {
	return nil
}
func (g *AddCustomTools) CustomToolsAdd(filePath string) string {
	return "当前系统:不支持自定义添加"
}
func (g *AddCustomTools) ExecCustomTools(id string) string {
	return "当前系统:不支持"
}
func (g *AddCustomTools) SaveCustomTools(objInfo Config.ToolsInfo) bool {
	return false
}
func (g *AddCustomTools) CustomToolsDel(id string) bool {
	return false
}
