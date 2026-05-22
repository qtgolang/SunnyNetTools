package Tools

import (
	"changeme/Service/Config"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/shirou/gopsutil/process"
	"time"
)

type Device struct {
	App  *SunnyNet.Sunny
	load bool
	list map[int32]string
	stop bool
}

func (g *Device) SetDeviceStopUpdate(stop bool) {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()
	g.stop = stop
	if g.load {
		g.update()
	}
}
func (g *Device) IsLoadDevice() bool {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()
	return g.load
}
func (g *Device) ProcessAny(open, StopNetwork bool) {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()
	g.App.ProcessALLName(open, StopNetwork)
}
func (g *Device) ProcessAddName(name string) {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()
	g.App.ProcessAddName(name)
}
func (g *Device) ProcessAddPid(pid int) {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()
	g.App.ProcessAddPid(pid)
}
func (g *Device) ProcessDelName(name string) {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()
	g.App.ProcessDelName(name)
}
func (g *Device) ProcessDelPid(pid int) {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()
	g.App.ProcessDelPid(pid)
}

// LoadDevice true表示使用NFAPI驱动 如果为false 表示使用Proxifier
func (g *Device) LoadDevice(mode int) bool {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()
	if g.load {
		return true
	}
	g.stop = false
	g.load = g.App.OpenDrive(mode)
	g.list = make(map[int32]string)
	if g.load {
		g.update()
	}
	return g.load
}

func (g *Device) update() {
	go func() {
		for {
			if g.stop {
				return
			}
			newList := make(map[int32]bool)
			arr, _ := process.Processes()
			for _, o := range arr {
				if g.list[o.Pid] == "" {
					n, _ := o.Name()
					if n != "" {
						g.list[o.Pid] = n
						Config.AppList["Main"].EmitEvent("DeviceUpdateProcessesList", o.Pid, n, false)
					}
				}
				newList[o.Pid] = true
			}
			for k, Name := range g.list {
				if !newList[k] {
					delete(g.list, k)
					Config.AppList["Main"].EmitEvent("DeviceUpdateProcessesList", k, Name, true)
				}
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
}

type ProcessesInfo struct {
	Pid  int32
	Name string
}

func (g *Device) GetAllProcessesList() []ProcessesInfo {
	var array []ProcessesInfo
	arr, _ := process.Processes()
	for _, v := range arr {
		n, _ := v.Name()
		if n != "" {
			p := v.Pid
			array = append(array, ProcessesInfo{Pid: p, Name: n})
		}
	}
	return array
}
