package collections

import (
	"log"
	"math/rand"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func BenchmarkLinkList_Add(b *testing.B) {
	type args struct {
		key   string
		value string
	}

	type Da struct {
		name   string
		args   args
	}
	tests := make([]Da,0)
	for i:=0;i<100;i++{
		tests=append(tests,Da{name:"fdsfds",args:args{key:strconv.FormatInt(rand.Int63n(10),10),value:strconv.FormatInt(rand.Int63n(99),10)}})
	}
	v := &LinkList{}
	for i:=0;i<b.N;i++ {
		v.Sort(v.SortDescFunc)
	}
}
func TestLinkList_Add(t *testing.T) {
	type args struct {
		key   string
		value string
	}

	type Da struct {
		name   string
		args   args
	}
	tests := make([]Da,0)
	for i:=0;i<100;i++{
		tests=append(tests,Da{name:"fdsfds",args:args{key:strconv.FormatInt(rand.Int63n(10),10),value:strconv.FormatInt(rand.Int63n(99),10)}})
	}
	v := &LinkList{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Add(tt.args.key,tt.args.value)
		})
	}

	t.Log(v)
	v.Sort(v.SortDescFunc)
	t.Log(v.RootNode)
	v.Sort(v.SortAscFunc)
	t.Log(v.RootNode)
}

func TestLinkList_Get(t *testing.T) {
	type fields struct {
		RootNode *Node
		Last     *Node
		Map      map[string]*Node
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &LinkList{
				RootNode: tt.fields.RootNode,
				Last:     tt.fields.Last,
				Map:      tt.fields.Map,
			}
			if got := v.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkList_GetMap(t *testing.T) {
	type fields struct {
		RootNode *Node
		Last     *Node
		Map      map[string]*Node
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &LinkList{
				RootNode: tt.fields.RootNode,
				Last:     tt.fields.Last,
				Map:      tt.fields.Map,
			}
			if got := v.GetMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkList_GetValue(t *testing.T) {
	type fields struct {
		RootNode *Node
		Last     *Node
		Map      map[string]*Node
	}
	tests := []struct {
		name   string
		fields fields
		want   url.Values
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &LinkList{
				RootNode: tt.fields.RootNode,
				Last:     tt.fields.Last,
				Map:      tt.fields.Map,
			}
			if got := v.GetValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkList_SortDesc(t *testing.T) {
	type fields struct {
		RootNode *Node
		Last     *Node
		Map      map[string]*Node
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	v := &LinkList{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
	t.Log(v)
}

func TestLinkList_Sort(t *testing.T) {

	v := &LinkList{}
	v.Add("b","2")



	abc:="abcdefghijklmnopqrstuvwxyz0123456789"

	r:=rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=0;i< 1000;i++{


		v.Add(
			string(abc[r.Intn(len(abc))])+string(abc[r.Intn(len(abc))])+string(abc[r.Intn(len(abc))])+string(abc[r.Intn(len(abc))])+string(abc[r.Intn(len(abc))]),
			string(abc[r.Intn(len(abc))])+string(abc[r.Intn(len(abc))])+string(abc[r.Intn(len(abc))])+string(abc[r.Intn(len(abc))])+string(abc[r.Intn(len(abc))]),

			)

	}

	type args struct {
		sortDescFunc func(a *Node, b *Node) int
	}
	tests := []struct {
		name   string
		args   args
	}{
		{ name: "TestLinkList_Sort_SortDescFunc",args: args{sortDescFunc:v.SortDescFunc}},
		{ name: "TestLinkList_Sort_SortAscFunc",args: args{sortDescFunc:v.SortAscFunc}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			v.Sort(tt.args.sortDescFunc)

			node:=v.RootNode
			for node!=nil{
				log.Println(node.Key,node.Value)
				node = node.Next
			}




		})
	}
}