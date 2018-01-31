package gweb

import "sync"

type Attributes struct {
	sync.RWMutex
	Map map[string]interface{}
}

func (att *Attributes) Put(key string, value interface{}) {
	att.Lock()
	att.Map[key] = value
	defer att.Unlock()
}
func (att *Attributes) Get(key string) interface{} {
	att.RLock()
	defer att.RUnlock()
	return att.Map[key]
}
func (att *Attributes) Delete(key string) {
	att.RLock()
	defer att.RUnlock()
	delete(att.Map, key)
}
