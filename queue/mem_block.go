package queue

import (
	"encoding/hex"
	"fmt"
	"github.com/nbvghost/glog"
	"math/rand"
	"sync"
	"time"
)

type MemBlock struct {
	lock      sync.RWMutex
	die       bool
	pools     *MemQueue
	Input     chan interface{}
	LastInput int64
	Hash      string
}

func (p *MemBlock) generatorHash() string {
	dest := [8]byte{}
	if _, err := rand.Read(dest[:]); err != nil {
		glog.Panic(err)
	}
	p.Hash = time.Now().Format("20060102150405.999999999") + "." + hex.EncodeToString(dest[:])
	return p.Hash
}

func (p *MemBlock) IsDie() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.die
}
func (p *MemBlock) Die() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.die = true

}
func (p *MemBlock) DetectDie() {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			//poolsLen := len(pools)
			if ((time.Now().UnixNano() - p.LastInput) / 1000 / 1000) > int64(Params.PoolTimeOut) {
				if p.IsDie() {
					if len(p.Input) == 0 {
						close(p.Input)
						glog.Trace(fmt.Sprintf("删除Pool,Hash：%v  ChanLen:%v  ChanCap:%v  删除成功：%v", p.Hash, len(p.Input), cap(p.Input), p.pools.Remove(p)))
						return
					}

				} else {
					p.Die()
				}
			}
		}
	}
}
func (p *MemBlock) ReceiveMessageTo(c chan<- interface{}) {
	t := time.NewTicker(time.Second * 10)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			//poolsLen := len(pools)
			if ((time.Now().UnixNano() - p.LastInput) / 1000 / 1000) > int64(Params.PoolTimeOut) {
				if p.IsDie() {
					if len(p.Input) == 0 {
						close(p.Input)
						glog.Trace(fmt.Sprintf("删除Pool,Hash：%v  ChanLen:%v  ChanCap:%v  删除成功：%v", p.Hash, len(p.Input), cap(p.Input), p.pools.Remove(p)))
						return
					}

				} else {
					p.Die()
				}
			}

		case msg := <-p.Input:
			c <- msg
		}
	}
}
func (p *MemBlock) CanInput() bool {
	if len(p.Input) == cap(p.Input) || p.IsDie() {
		return false
	} else {
		return true
	}
}
