package queue

import (
	"fmt"
	"github.com/nbvghost/glog"
	"testing"
	"time"
)

func init() {

}
func TestPools_Remove(t *testing.T) {
	PoolTool := NewPools(Order)
	PoolTool.createPool()
	PoolTool.createPool()
	PoolTool.createPool()
	PoolTool.createPool()

	type args struct {
		target *MemBlock
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "TestPools_Remove", args: args{target: PoolTool.list[3]}},
		{name: "TestPools_Remove", args: args{target: PoolTool.list[2]}},
		{name: "TestPools_Remove", args: args{target: PoolTool.list[0]}},
		{name: "TestPools_Remove", args: args{target: PoolTool.list[1]}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//p.Remove(tt.args.target)
		})
	}
	for {
		glog.Trace(PoolTool.list)
		time.Sleep(time.Second)
	}

}

func TestPool_generatorHash(t *testing.T) {

	tests := []struct {
		name string
		want string
	}{
		{name: "TestPool_generatorHash"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &MemBlock{}
			if got := p.generatorHash(); got != tt.want {
				t.Errorf("generatorHash() = %v, want %v", got, tt.want)
			}
			t.Log(p.Hash)
			t.Log(time.Now().Format("20060102150405.999999999"))
		})
	}
}
func BenchmarkPools_Push(b *testing.B) {
	p := NewPools(Order)
	fmt.Println(b.N)
	for i := 0; i < b.N; i++ {
		p.Push(125)
	}
}
func TestPools_Push(t *testing.T) {

	Params.MaxWaitCollectMessage = 1000
	Params.MaxPoolNum = 10

	p := NewPools(Order)
	//p := NewPools(Order)

	n := 0

	go func() {
		for {

			p.Push(n)

			n++
			if n > 50 {
				// n=0
				//time.Sleep(time.Second*20)
			}
			//time.Sleep(time.Millisecond * 1)

		}
	}()

	go func() {

		for {
			df := p.GetMessage(10)
			fmt.Println(len(df))
			time.Sleep(time.Millisecond * 1000)
		}

	}()

	for {
		p.printStat()
		time.Sleep(time.Second)
	}
}
