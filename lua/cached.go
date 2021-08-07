package lua

import (
	"reflect"
	"sync"
)

var udataCached = make(map[uintptr]interface{})
var udataMutex sync.Mutex

func CachedAdd(it interface{}) {
	udataMutex.Lock()
	defer udataMutex.Unlock()
	ptr := reflect.ValueOf(it).Pointer()
	if _, ok := udataCached[ptr]; !ok {
		udataCached[ptr] = it
	}
}

func CachedRemove(it interface{}) {
	udataMutex.Lock()
	defer udataMutex.Unlock()
	ptr := reflect.ValueOf(it).Pointer()
	delete(udataCached, ptr)
}

func CachedAddRemove(add interface{}, remove interface{}) {
	if remove == add {
		return
	}
	remove1 := reflect.ValueOf(remove).Pointer()
	add1 := reflect.ValueOf(add).Pointer()
	if remove1 == add1 {
		return
	}
	udataMutex.Lock()
	defer udataMutex.Unlock()

	delete(udataCached, remove1)
	udataCached[add1] = add
}

