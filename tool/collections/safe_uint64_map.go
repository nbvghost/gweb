package collections

import "sync"

type SafeUint64StringMap struct {
	sync.RWMutex
	Map map[uint64]string
}

func (ssm *SafeUint64StringMap) Put(key uint64, value string) {
	ssm.Lock()
	defer ssm.Unlock()
	if ssm.Map==nil{
		ssm.Map = make(map[uint64]string)
	}

	ssm.Map[key] = value

}
func (ssm *SafeUint64StringMap) Del(k uint64) {
	ssm.Lock()
	defer ssm.Unlock()
	delete(ssm.Map, k)
	//db.NotifyAll(&db.Message{db.Socket_Type_2_STC,k})
}

func (ssm *SafeUint64StringMap) Get(key uint64) (string,bool) {
	ssm.Lock()
	defer ssm.Unlock()

	v,have:=ssm.Map[key]

	return v,have

}
