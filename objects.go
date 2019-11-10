package objects

import (
	"sync"
)

type ObjectId int32

var manager struct {
	sync.Mutex
	objs     map[ObjectId]interface{}
	freeObjs map[ObjectId]interface{}
	next     ObjectId
}

func init() {
	manager.Lock()
	defer manager.Unlock()

	manager.objs = make(map[ObjectId]interface{})
	manager.freeObjs = make(map[ObjectId]interface{})
	manager.next = 1
}

func NewObjectId(obj interface{}) ObjectId {
	manager.Lock()
	defer manager.Unlock()

	var id ObjectId

	if len(manager.freeObjs) == 0 {
		id = manager.next
		manager.next++
	} else {
		for obj := range manager.freeObjs {
			id = obj
			break
		}
		delete(manager.freeObjs, id)
	}

	manager.objs[id] = obj
	return id
}

func Reset() {
	manager.Lock()
	defer manager.Unlock()

	manager.objs = make(map[ObjectId]interface{})
	manager.freeObjs = make(map[ObjectId]interface{})
	manager.next = 1
}

func (id ObjectId) IsNil() bool {
	return id == 0
}

func (id ObjectId) Get() interface{} {
	manager.Lock()
	defer manager.Unlock()

	return manager.objs[id]
}

func (id *ObjectId) Free() interface{} {
	manager.Lock()
	defer manager.Unlock()

	obj := manager.objs[*id]
	delete(manager.objs, *id)

	manager.freeObjs[*id] = nil

	return obj
}
