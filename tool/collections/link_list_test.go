package collections

import (
	"net/url"
	"reflect"
	"testing"
)

func TestLinkList_Add(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		args   args
	}{
		{name:"fdsfds",args:args{key:"1",value:"1"}},
		{name:"fdsfds",args:args{key:"2",value:"2"}},
		{name:"fdsfds",args:args{key:"3",value:"3"}},
		{name:"fdsfds",args:args{key:"4",value:"4"}},
		{name:"fdsfds",args:args{key:"5",value:"5"}},
	}
	v := &LinkList{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Add(tt.args.key,tt.args.value)
		})
	}

	t.Log(v)
	v.SortDesc()
	t.Log(v)
}

func TestLinkList_Get(t *testing.T) {
	type fields struct {
		RootNode *KV
		Last     *KV
		Map      map[string]*KV
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *KV
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
		RootNode *KV
		Last     *KV
		Map      map[string]*KV
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
		RootNode *KV
		Last     *KV
		Map      map[string]*KV
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
		RootNode *KV
		Last     *KV
		Map      map[string]*KV
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