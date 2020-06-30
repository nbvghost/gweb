package collections

import "sync"

type SafeStringMap struct {
	sync.RWMutex
	Map map[string]string
}

func (ssm *SafeStringMap) Put(key string, value string) {
	ssm.Lock()
	defer ssm.Unlock()
	if ssm.Map==nil{
		ssm.Map = make(map[string]string)
	}

	ssm.Map[key] = value

}
func (ssm *SafeStringMap) Del(k string) {
	ssm.Lock()
	defer ssm.Unlock()
	delete(ssm.Map, k)
	//db.NotifyAll(&db.Message{db.Socket_Type_2_STC,k})
}

func (ssm *SafeStringMap) Get(key string) (string,bool) {
	ssm.Lock()
	defer ssm.Unlock()

	v,have:=ssm.Map[key]

	return v,have

}
