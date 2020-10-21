package utils

import "sync"

var gvlock *sync.Mutex
var gvs map[string]interface{}

func init() {
	gvlock = &sync.Mutex{}
	gvs = map[string]interface{}{}
}

func SetVar(key string, value interface{}) {
	gvlock.Lock()
	defer gvlock.Unlock()
	gvs[key] = value
}

func GetVar(key string) (interface{}, bool) {
	gvlock.Lock()
	defer gvlock.Unlock()
	val, found := gvs[key]
	return val, found
}

func RemoveVar(key string) {
	gvlock.Lock()
	defer gvlock.Unlock()
	delete(gvs, key)
}
