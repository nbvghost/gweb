package collections

import (
	"fmt"
	"net/url"
)

type KV struct {
	Key   string
	Value string
	Previous  *KV
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
func (v *LinkList) SortAscFunc(a *KV,b *KV) int{
	if a ==nil || b==nil{
		return 0
	}
	if b.Next==nil{
		return 0
	}
	if a.Key < b.Key{

		return 1

	}else if a.Key > b.Key{

		return -1

	}else{

		if a.Value < b.Value{
			return 1

		}else if a.Value > b.Value{
			return -1
		}else{
			//==0
			return 0

		}

	}
}
func (v *LinkList) SortDescFunc(a *KV,b *KV) int{
	if a ==nil || b==nil{
		return 0
	}
	if b.Next==nil{
		return 0
	}
	if a.Key < b.Key{

		return -1

	}else if a.Key > b.Key{

		return 1

	}else{

		if a.Value < b.Value{
			return -1

		}else if a.Value > b.Value{
			return 1
		}else{
			//==0
			return 0

		}

	}
}
//todo:还没有写完
func (v *LinkList) Sort(sortDescFunc func(a *KV,b *KV)int){
	//desc -1
	//asc 1
	//eque 0
	v.Last=nil
	nextNode:=v.RootNode

	for nextNode!=nil{
		var swpNode *KV
		if nextNode.Next!=nil{
			sortType:=sortDescFunc(nextNode,nextNode.Next)
			if sortType ==-1{
				isRootNode:=false
				if nextNode.Previous==nil{
					isRootNode=true
				}

				previousNode:=nextNode.Previous

				upNode:=nextNode.Next
				upNodeChild:=nextNode.Next.Next

				upNode.Previous=previousNode//
				nextNode.Previous=upNode
				upNode.Next=nextNode

				upNodeChild.Previous=upNode.Next
				upNode.Next.Next=upNodeChild


				if isRootNode{
					upNode.Previous=nil
					v.RootNode=upNode
				}else{
					previousNode.Next=upNode
				}
				swpNode =upNode

			}else if sortType ==1{

			}

			if swpNode==nil{
				nextNode=nextNode.Next
			}else{
				if swpNode.Previous==nil{
					nextNode=v.RootNode
				}else{
					nextNode=swpNode.Previous
				}
			}
		}else{

			newLastNode:=nextNode.Previous
			if nextNode==v.Last{
				break
			}
			v.Last = newLastNode

			lastNode:=nextNode
			lastNode.Previous.Next=nil
			lastNode.Previous=nil

			v.RootNode.Previous=lastNode

			lastNode.Next=v.RootNode
			v.RootNode=lastNode

			nextNode=v.RootNode

		}

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

func (v *LinkList) Add(key, value string) *KV {
	if v.Map == nil {
		v.Map = make(map[string]*KV)
	}
	//item, ok := v.Map[key]

		node := &KV{Value: value, Key: key}
		v.Map[key] = node
		if v.RootNode == nil {
			v.RootNode = node
			v.Last = node
		}else if v.Last != nil {
			node.Previous=v.Last
			v.Last.Next = node
			v.Last = node
		}
		return node
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