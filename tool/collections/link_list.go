package collections

import (
	"fmt"
	"net/url"
)

type KV struct {
	Key   string
	Value string
	Next  *KV
}

func (kv *KV)String()string  {
	return fmt.Sprintf("{Key:%v,Value:%v,Next:%v}",kv.Key,kv.Value,kv.Next)
}
type LinkList struct {
	RootNode *KV
	Last     *KV
	Map      map[string]*KV
}

//todo:还没有写完
func (v *LinkList) SortDesc() {

	//desc -1
	//asc 1
	//eque 0
	sortFunc:= func(a *KV,b *KV) int{
		if a ==nil || b==nil{
			return 0
		}
		if b.Next==nil{
			return 0
		}



		if a.Key < b.Next.Key{

			return -1

		}else if a.Key > b.Next.Key{

			return 1

		}else{

			if a.Value < b.Next.Value{
				return -1

			}else if a.Key > b.Next.Key{
				return 1
			}else{
				//==0
				return 0

			}

		}

	}
	if sortFunc(nil,nil)>0{

	}

	readNode := v.RootNode
	for readNode != nil{

		if readNode.Next==nil{
			break
		}


		//sortType:=sortFunc()


		if readNode.Key < readNode.Next.Key{

			Key   :=readNode.Key
			Value :=readNode.Value

			readNode.Key=readNode.Next.Key
			readNode.Value=readNode.Next.Value

			readNode.Next.Key=Key
			readNode.Next.Value=Value






		}else if readNode.Key > readNode.Next.Key{



		}else{
			//==0
			if readNode.Value < readNode.Next.Value{

				Key   :=readNode.Key
				Value :=readNode.Value

				readNode.Key=readNode.Next.Key
				readNode.Value=readNode.Next.Value

				readNode.Next.Key=Key
				readNode.Next.Value=Value

			}else if readNode.Key > readNode.Next.Key{

			}else{
				//==0

			}

		}

		readNode = readNode.Next

	}
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
		}else if v.Last != nil {
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