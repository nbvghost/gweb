package collections

import "sync"

type SafeMap struct {
	sync.RWMutex
	Map map[interface{}]interface{}
}

func (this *SafeMap) Put(key interface{}, value interface{}) {
	this.Lock()
	defer this.Unlock()
	this.Map[key] = value

}
func (this *SafeMap) DelAll() {
	this.Lock()
	defer this.Unlock()
	this.Map = make(map[interface{}]interface{})
}
func (this *SafeMap) Get(key interface{}) interface{} {
	this.Lock()
	defer this.Unlock()
	return this.Map[key]

}
