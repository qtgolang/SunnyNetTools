package Tools

import (
	"changeme/Service/Config"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"sort"
	"sync"
)

type ReplaceBody struct {
	App *SunnyNet.Sunny
}

var mapReplaceLock sync.Mutex

// CreateReplaceBody 创建替换
func (g *ReplaceBody) CreateReplaceBody() int {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapReplaceLock.Lock()
	defer mapReplaceLock.Unlock()
	Config.CertID++
	id := Config.CertID
	Config.Config.ReplaceRoles[id] = &Config.ReplaceBodyInfo{
		ID: id,
	}
	return id
}
func (g *ReplaceBody) ReplaceBodyUpdate(id int, Type, Source, Lod, New, Note, state string) any {
	mapReplaceLock.Lock()
	defer mapReplaceLock.Unlock()
	obj := Config.Config.ReplaceRoles[id]
	if obj == nil {
		return false
	}
	obj.Type = Type
	obj.Lod = Lod
	obj.New = New
	obj.Source = Source
	obj.Note = Note
	obj.State = state
	ok := obj.Parse()
	obj.Ok = ok
	Config.Config.Save()
	return obj
}

func (g *ReplaceBody) ReplaceBodyRemove(id int) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapReplaceLock.Lock()
	defer mapReplaceLock.Unlock()
	obj := Config.Config.ReplaceRoles[id]
	if obj == nil {
		return
	}
	delete(Config.Config.ReplaceRoles, id)
}

func (g *ReplaceBody) ReplaceBodyList() []Config.ReplaceBodyInfo {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapReplaceLock.Lock()
	defer mapReplaceLock.Unlock()
	var array []Config.ReplaceBodyInfo
	var keys []int
	for i := range Config.Config.ReplaceRoles {
		keys = append(keys, i)
	}
	sort.Ints(keys)
	for _, k := range keys {
		obj := Config.Config.ReplaceRoles[k]
		array = append(array, *obj)
	}
	return array
}
