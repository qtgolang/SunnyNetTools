package Tools

import (
	"changeme/Service/Config"
	"encoding/base64"
	"fmt"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"sort"
	"sync"
	"time"
)

type ScriptLog struct {
	App *SunnyNet.Sunny
}
type LogInfo struct {
	Time string
	Info []string
}

var logLock sync.Mutex
var scriptLog = make(map[int]LogInfo)
var scriptLogId = 0

func (l *ScriptLog) ScriptLogInit() {
	go l.init()
}
func (l *ScriptLog) init() {
	for {
		time.Sleep(time.Millisecond * 500)
		logLock.Lock()
		if len(scriptLog) > 0 {
			var keys []int
			for i, _ := range scriptLog {
				keys = append(keys, i)
			}
			sort.Ints(keys)
			var array []LogInfo
			for _, k := range keys {
				array = append(array, scriptLog[k])
			}
			Config.AppList["Main"].EmitEvent("addScriptLog", array)
			scriptLog = make(map[int]LogInfo)
			scriptLogId = 0
		}
		logLock.Unlock()
	}
}
func (l *ScriptLog) PrintScriptLog(_ int, info ...any) {
	logLock.Lock()
	defer logLock.Unlock()
	var array []string
	for _, v := range info {
		if v == nil {
			array = append(array, "nil")
			continue
		}
		objArray, ok := v.([]any)
		if !ok {
			array = append(array, "nil")
			continue
		}
		for _, vv := range objArray {
			array = append(array, fmt.Sprintf("%v", vv))
		}
	}
	if len(array) > 0 {
		scriptLogId++
		scriptLog[scriptLogId] = LogInfo{Info: array, Time: time.Now().Format("2006-01-02 15:04:05")}
	}
}
func (l *ScriptLog) SaveScriptCode(_ int, code []byte) {
	Config.Config.ScriptCode = base64.StdEncoding.EncodeToString(code)
	Config.Config.Save()
}
