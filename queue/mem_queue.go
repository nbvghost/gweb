package queue

import (
	"errors"
	"fmt"
	"github.com/nbvghost/glog"
	"strings"
	"sync"
	"time"
)

type messageOrder string

const DisOrder messageOrder = "DisOrder"
const Order messageOrder = "Order"

type MemQueue struct {
	list       []*MemBlock
	msgChan    chan interface{}
	PrintTime  time.Time
	Order      messageOrder
	maxPoolNum int
	poolNum    int

	inputTotalNum   uint
	outTotalNum     uint
	processTotalNum uint

	locker *sync.RWMutex

	totalNumLocker *sync.RWMutex

	once *sync.RWMutex

	workingPool *MemBlock

	workingReadPool *MemBlock
}

func NewPools(Order messageOrder) *MemQueue {
	pt := map[string]interface{}{
		"Name":                  "task.Pools",
		"PoolSize":              Params.PoolSize,
		"PoolTimeOut":           Params.PoolTimeOut,
		"MaxProcessMessageNum":  Params.MaxProcessMessageNum,
		"MaxWaitCollectMessage": Params.MaxWaitCollectMessage,
		"MaxPoolNum":            Params.MaxPoolNum,
	}

	if Params.PoolSize <= 0 {
		panic(errors.New("task.Params.PoolSize,不能有零值"))
	}

	if Params.PoolTimeOut <= 0 {
		panic(errors.New("task.Params.PoolTimeOut,不能有零值"))
	}

	if Params.MaxProcessMessageNum <= 0 {
		panic(errors.New("task.Params.MaxProcessMessageNum,不能有零值"))
	}

	if Params.MaxWaitCollectMessage <= 0 {
		panic(errors.New("task.Params.MaxWaitCollectMessage,不能有零值"))
	}

	if Params.MaxPoolNum < 0 {
		panic(errors.New("task.Params.MaxPoolNum,不能为负数"))
	}

	glog.Trace(pt)

	p := &MemQueue{msgChan: make(chan interface{}, Params.MaxProcessMessageNum), Order: Order,
		totalNumLocker: &sync.RWMutex{},
		once:           &sync.RWMutex{},
		locker:         &sync.RWMutex{}}

	return p

}

func (p *MemQueue) listenPools() <-chan interface{} {
	if p.Order == DisOrder {
		return p.msgChan
	} else if p.Order == Order {
		p.workingReadPool = p.getCanReadPool()
		if p.workingReadPool == nil {
			return nil
		} else {
			return p.workingReadPool.Input
		}

	} else {
		glog.Panic(errors.New(fmt.Sprintf("Pools 没有匹配到处理规则：%v", p.Order)))
		return nil
	}

}

func (p *MemQueue) getCanReadPool() *MemBlock {
	if len(p.list) > 1 {
		pool := p.list[0]
		if pool.CanInput() == false && len(pool.Input) == 0 {
			p.Remove(pool)
			return p.list[0]
		}
		return pool
	} else if len(p.list) == 1 {
		return p.list[0]
	} else {
		return nil
	}

}
func (p *MemQueue) Len() int {

	return len(p.list)
}

func (p *MemQueue) createPool() *MemBlock {

	if len(p.list) >= Params.MaxPoolNum && Params.MaxPoolNum != 0 {
		return p.list[len(p.list)-1]
	}

	pool := &MemBlock{
		Input: make(chan interface{}, Params.PoolSize), LastInput: time.Now().UnixNano(), pools: p,
	}
	pool.generatorHash()

	if p.Order == DisOrder {
		go pool.ReceiveMessageTo(p.msgChan)

	} else if p.Order == Order {
		//顺序通道，写的时候是往队列中最后一个写，如果最一个满了，会新建一个pool,些时前一个pool 就可以关闭，因为pool里的数据被读完，不会在往里写了
		if len(p.list) > 0 {
			p.list[len(p.list)-1].Die()
			//close(p.list[len(p.list)-1].Input)
		}
	}

	p.list = append(p.list, pool)

	if len(p.list) > p.maxPoolNum {
		p.maxPoolNum = len(p.list)
	}

	p.poolNum = len(p.list)

	return pool
}

func (p *MemQueue) readMessage(num int) []interface{} {
	msgs := make([]interface{}, 0)
	defer func() {
		p.OutMany(uint(len(msgs)))
		p.printStat()
	}()

	if num == 0 {
		return msgs
	}
	t := time.NewTicker(time.Duration(Params.MaxWaitCollectMessage) * time.Millisecond)
	defer t.Stop()

	for {

		select {
		case <-t.C:
			p.printStat()
			//如果有收集的消息的话，在超时后返回，没有的话，继续收集
			if len(msgs) > 0 {
				return msgs
			} else {
				continue
			}

		case msg, isOpen := <-p.listenPools():
			if isOpen && msg != nil {
				msgs = append(msgs, msg)
				if len(msgs) >= num {
					return msgs
				}
			} else {
				glog.Trace(map[string]interface{}{"Message": "GetMessage取到空的对象", "MsgIsNil": msg == nil, "IsOpen": isOpen})
				time.Sleep(time.Second)

			}

		}

	}
}
func (p *MemQueue) GetMessage(num int) []interface{} {
	p.once.Lock()
	defer p.once.Unlock()
	return p.readMessage(num)

}
func (p *MemQueue) Remove(target *MemBlock) bool {
	p.locker.Lock()
	defer p.locker.Unlock()

	for i := 0; i < len(p.list); i++ {

		if p.list[i] == target && target.IsDie() {
			close(p.list[i].Input)
			p.list = append(p.list[:i], p.list[i+1:]...)
			return true
		}
	}
	return false
}
func (p *MemQueue) Push(messages ...interface{}) {

	p.push(messages...)

	p.InputMany(uint(len(messages)))
}
func (p *MemQueue) push(messages ...interface{}) {
	if p.workingPool == nil {
		p.workingPool = p.getAbleWritePool()
	}

	for index := range messages {

		ticker := time.NewTicker(time.Millisecond * 1000)
	writeTimeout:
		for {
			select {
			case <-ticker.C:
				oldHash := p.workingPool.Hash
				newHash := ""
				glog.Debug(fmt.Sprintf("缓冲区已满:%v", oldHash))
				if len(p.workingPool.Input) == cap(p.workingPool.Input) {
					p.workingPool = p.getAbleWritePool()
					newHash = p.workingPool.Hash
				}

				if strings.EqualFold(oldHash, newHash) {
					glog.Debug(fmt.Sprintf("缓冲池已满，缓冲区数量：%v，最大缓冲区数量：%v", len(p.list), Params.MaxPoolNum))
				} else {
					glog.Debug(fmt.Sprintf("新建缓冲区:%v", p.workingPool.Hash))
				}

			case p.workingPool.Input <- messages[index]:
				break writeTimeout

			}
		}
		ticker.Stop()

	}

	p.workingPool.LastInput = time.Now().UnixNano()

}
func (p *MemQueue) getAbleWritePool() *MemBlock {

	var ablePool *MemBlock
	length := len(p.list)
	if p.Order == DisOrder {

		for i := length - 1; i >= 0; i-- {
			if p.list[i].CanInput() {
				ablePool = p.list[i]
				return ablePool
			}
		}

		ablePool = p.createPool()

	} else if p.Order == Order {
		p.locker.Lock()
		defer p.locker.Unlock()

		if len(p.list) > 0 {
			_ablePool := p.list[len(p.list)-1]
			if _ablePool.CanInput() {
				ablePool = _ablePool
			}
		}

		if ablePool == nil {
			ablePool = p.createPool()
		}

	}
	return ablePool
}

func (p *MemQueue) OutMany(num uint) {
	p.totalNumLocker.Lock()
	defer p.totalNumLocker.Unlock()
	p.outTotalNum = p.outTotalNum + num
}
func (p *MemQueue) InputMany(num uint) {
	p.totalNumLocker.Lock()
	defer p.totalNumLocker.Unlock()
	p.inputTotalNum = p.inputTotalNum + num
}
func (p *MemQueue) ProcessOne() {
	p.totalNumLocker.Lock()
	defer p.totalNumLocker.Unlock()
	p.processTotalNum++
}
func (p *MemQueue) IsEmpty() bool {

	for i := 0; i < len(p.list); i++ {
		if len(p.list[i].Input) > 0 {
			return false
		}
	}

	p.printStat()
	if p.outTotalNum != p.inputTotalNum || p.outTotalNum != p.processTotalNum {
		return false
	}

	return true

}
func (p *MemQueue) printStat() {
	now := time.Now()
	if now.Sub(p.PrintTime) > time.Second*10 {
		p.poolNum = len(p.list)
		glog.Trace(fmt.Sprintf("MaxPoolNum:%v   PoolNum:%v  InputTotalNum:%v   OutTotalNum:%v   ProcessTotalNum:%v", p.maxPoolNum, p.poolNum, p.inputTotalNum, p.outTotalNum, p.processTotalNum))
		p.PrintTime = now
	}
}
