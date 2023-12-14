package main

import (
	"context"
	"encoding/json"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/protobuf/JSON"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strconv"
)

// App struct
type App struct {
	ctx context.Context
	App *SunnyNet.Sunny
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

type Command struct {
	Command string
	Args    any
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Do(arg any) any {
	defer func() {
		recover()
	}()
	as := arg.(map[string]any)
	if as == nil {
		return nil
	}
	_t := as["Command"]
	if _t == nil {
		return nil
	}
	command := _t.(string)
	Args := as["Args"]
	if Args == nil {
		return event(command, nil)
	}
	bs, _ := json.Marshal(Args)
	_json := string(bs)
	sj := JSON.NewSyJson()
	sj.Parse(_json)
	return event(command, sj)
}
func (a *App) CallDo(arg ...any) {
	runtime.EventsEmit(a.ctx, "Do", arg...)
}
func getFloat64(arg any) float64 {
	switch v := arg.(type) {
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	case int16:
		return float64(v)
	case int8:
		return float64(v)
	case uint:
		return float64(v)
	case uint64:
		return float64(v)
	case uint32:
		return float64(v)
	case uint16:
		return float64(v)
	case uint8:
		return float64(v)
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0
		}
		return f
	case float32:
		return float64(v)
	case float64:
		return v
	default:
		return 0
	}
}
func getInt64(arg any) int64 {
	switch v := arg.(type) {
	case int:
		return int64(v)
	case int64:
		return int64(v)
	case int32:
		return int64(v)
	case int16:
		return int64(v)
	case int8:
		return int64(v)
	case uint:
		return int64(v)
	case uint64:
		return int64(v)
	case uint32:
		return int64(v)
	case uint16:
		return int64(v)
	case uint8:
		return int64(v)
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0
		}
		return i
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	default:
		return 0
	}
}
func getInt(arg any) int {
	return int(getInt64(arg))
}
