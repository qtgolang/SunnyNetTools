package Tools

import (
	"changeme/Service/Config"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"sort"
	"sync"
)

type ReplaceHost struct {
	App *SunnyNet.Sunny
}

var mapReplaceHostLock sync.Mutex

func (g *ReplaceBody) CreateReplaceHost() int {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapReplaceHostLock.Lock()
	defer mapReplaceHostLock.Unlock()
	Config.CertID++
	id := Config.CertID
	Config.Config.ReplaceHost[id] = &Config.ReplaceHostInfo{
		ID: id,
	}
	return id
}
func (g *ReplaceBody) ReplaceHostUpdate(id int, Lod, New, Note string) bool {
	mapReplaceHostLock.Lock()
	defer mapReplaceHostLock.Unlock()
	obj := Config.Config.ReplaceHost[id]
	if obj == nil {
		return false
	}
	obj.LodInfo.Host = Config.InvalidHost
	obj.NewInfo.Host = Config.InvalidHost
	obj.LodInfo.Parse(Lod, 0)
	obj.NewInfo.Parse(New, 0)
	
	if obj.LodInfo.Port == 0 {
		obj.Lod = obj.LodInfo.Host
	} else {
		obj.Lod = obj.LodInfo.String()
	}

	if obj.NewInfo.Port == 0 {
		obj.New = obj.NewInfo.Host
	} else {
		obj.New = obj.NewInfo.String()
	}

	obj.Note = Note
	Config.Config.Save()
	return true
}

func (g *ReplaceBody) ReplaceHostRemove(id int) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapReplaceHostLock.Lock()
	defer mapReplaceHostLock.Unlock()
	obj := Config.Config.ReplaceHost[id]
	if obj == nil {
		return
	}
	delete(Config.Config.ReplaceHost, id)
	Config.Config.Save()
}

func (g *ReplaceBody) ReplaceHostList() []Config.ReplaceHostInfo {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapReplaceLock.Lock()
	defer mapReplaceLock.Unlock()
	var keys []int
	for i := range Config.Config.ReplaceHost {
		keys = append(keys, i)
	}
	sort.Ints(keys)
	var array []Config.ReplaceHostInfo
	for _, k := range keys {
		obj := Config.Config.ReplaceHost[k]
		array = append(array, *obj)
	}
	return array
}
