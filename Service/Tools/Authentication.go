package Tools

import (
	"changeme/Service/Config"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"sort"
	"sync"
)

type Authentication struct {
	App *SunnyNet.Sunny
}

var authenticationLock sync.Mutex

func (g *ReplaceBody) CreateAuthentication() int {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	authenticationLock.Lock()
	defer authenticationLock.Unlock()
	Config.CertID++
	id := Config.CertID
	Config.Config.Authentication[id] = &Config.AuthenticationInfo{
		ID: id,
	}
	return id
}
func (g *ReplaceBody) AuthenticationUpdate(id int, User, Pass string) bool {
	authenticationLock.Lock()
	defer authenticationLock.Unlock()
	obj := Config.Config.Authentication[id]
	if obj == nil {
		return false
	}
	if obj.User != "" {
		g.App.Socket5DelUser(obj.User)
	}
	obj.User = User
	obj.Pass = Pass
	g.App.Socket5AddUser(User, Pass)
	Config.Config.Save()
	return true
}
func (g *ReplaceBody) AuthenticationRemove(id int) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	authenticationLock.Lock()
	defer authenticationLock.Unlock()
	obj := Config.Config.Authentication[id]
	if obj == nil {
		return
	}
	delete(Config.Config.Authentication, id)
}
func (g *ReplaceBody) AuthenticationList() []Config.AuthenticationInfo {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	authenticationLock.Lock()
	defer authenticationLock.Unlock()
	var keys []int
	for i := range Config.Config.Authentication {
		keys = append(keys, i)
	}
	sort.Ints(keys)
	var array []Config.AuthenticationInfo
	for _, k := range keys {
		obj := Config.Config.Authentication[k]
		array = append(array, *obj)
	}
	return array
}
