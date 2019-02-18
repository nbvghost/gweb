package gweb

import "sync"

type Attributes struct {
	Map sync.Map
}

func (att *Attributes) Put(key string, value interface{}) {
	//att.Lock()
	//att.Map[key] = value
	att.Map.Store(key,value)
	//defer att.Unlock()
}
func (att *Attributes) Get(key string) interface{} {
	//att.RLock()
	//defer att.RUnlock()
	v,_:=att.Map.Load(key)
	return v
}
func (att *Attributes) Delete(key string) {
	//att.RLock()
	//defer att.RUnlock()
	//delete(att.Map, key)
	att.Map.Delete(key)
}
