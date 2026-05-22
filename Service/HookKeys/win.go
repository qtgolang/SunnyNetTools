//go:build windows
// +build windows

package HookKeys

import (
	"encoding/json"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

func KeyPressed(vk int) bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
	return ret&0x8000 != 0
}

func init() {
	go func() {
		for {
			shift := KeyPressed(0x10)
			ctrl := KeyPressed(0x11)
			alt := KeyPressed(0x12)
			win := KeyPressed(0x5B)
			for id, v := range _keys {
				lock.Lock()
				is := isEditKeyDown
				lock.Unlock()
				if !is && (v.Up == false && v.Alt == alt && v.Ctrl == ctrl && v.Shift == shift && v.Meta == win) {
					if KeyPressed(v.Code) {
						v.Fun(id)
						v.Up = true
					}
				} else {
					v.Up = false
				}
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()
}

var isEditKeyDown = false
var _keys = make(map[string]registerKey)
var lock sync.Mutex

func UpdateEditKeyDown(i bool) {
	lock.Lock()
	defer lock.Unlock()
	isEditKeyDown = i
}
func RegisterKeys(list string, fun func(id string)) {
	var listArray []registerKey
	_ = json.Unmarshal([]byte(list), &listArray)
	for _, v := range listArray {
		if strings.HasPrefix(v.ID, "system_") {
			continue
		}
		if v.Meta {
			continue
		}
		if v.Code == 0 {
			continue
		}
		lock.Lock()
		delete(_keys, v.ID)
		v.Fun = fun
		_keys[v.ID] = v
		lock.Unlock()
	}
}

type registerKey struct {
	Name  string          `json:"Name"`
	Ctrl  bool            `json:"ctrlKey"`
	Alt   bool            `json:"altKey"`
	Shift bool            `json:"shiftKey"`
	Meta  bool            `json:"metaKey"`
	Key   string          `json:"key"`
	Code  int             `json:"keyCode"`
	Value string          `json:"value"`
	ID    string          `json:"ID"`
	Fun   func(id string) `json:"-"`
	Up    bool            `json:"up"`
}
