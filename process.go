package main

import (
	"changeme/CommAnd"
	"github.com/qtgolang/SunnyNet/src/JsCall"
	"github.com/qtgolang/SunnyNet/src/protobuf/JSON"
)

func processEvent(command string, args *JSON.SyJson) any {
	switch command {
	case "加载驱动":
		return app.App.StartProcess()
	case "枚举进程":
		p := CommAnd.EnumerateProcesses()
		return p
	case "进程驱动添加PID":
		gx := args.GetData("isSelected") == "true"
		app.App.ProcessCancelAll()
		if gx {
			app.App.ProcessAddPid(getInt(args.GetData("PID")))
		} else {
			app.App.ProcessDelPid(getInt(args.GetData("PID")))
		}
		return true
	case "进程驱动添加进程名":
		gx := args.GetData("isSet") == "true"
		Name := JsCall.ToGBK(args.GetData("Name"))
		if Name == "{OpenALL}" {
			app.App.ProcessALLName(gx)
			return true
		}
		if gx {
			app.App.ProcessAddName(Name)
		} else {
			app.App.ProcessDelName(Name)
		}
		return true
	default:
		return ReplaceRulesEvent(command, args)
	}
}
