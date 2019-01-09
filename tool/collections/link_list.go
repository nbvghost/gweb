package collections

import "net/url"

type KV struct {
	Key   string
	Value string
	Next  *KV
}
type LinkList struct {
	RootNode *KV
	Last     *KV
	Map      map[string]*KV
}

func (v *LinkList) Get(key string) *KV {
	if v.Map == nil {
		v.Map = make(map[string]*KV)
		return &KV{}
	}
	_, ok := v.Map[key]
	if ok {
		return v.Map[key]
	}
	return &KV{}
}

func (v *LinkList) Add(key, value string) {
	if v.Map == nil {
		v.Map = make(map[string]*KV)
	}
	item, ok := v.Map[key]
	if ok {
		item.Value = value
	} else {
		node := &KV{Value: value, Key: key}
		v.Map[key] = node
		if v.RootNode == nil {
			v.RootNode = node
			v.Last = node
		}
		if v.Last != nil {
			v.Last.Next = node
			v.Last = node
		}

	}
}
func (v *LinkList) GetMap() map[string]string {
	data := make(map[string]string)
	for key := range v.Map {
		item := v.Map[key]
		data[key] = item.Value
	}
	return data
}
func (v *LinkList) GetValue() url.Values {
	data := url.Values{}
	for key := range v.Map {
		item := v.Map[key]
		data.Set(key, item.Value)
	}
	return data
}