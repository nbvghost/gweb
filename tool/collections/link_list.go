package collections

import (
    "fmt"
    "log"
    "net/url"
)

type Node struct {
    Key      string
    Value    string
    Previous *Node
    Next     *Node
}

func (node *Node) String() string {
    return fmt.Sprintf("{Key:%v,Value:%v,Next:%v}", node.Key, node.Value, node.Next)
}

type LinkList struct {
    RootNode *Node
    Last     *Node
    Map      map[string]*Node
}

func (v *LinkList) SortAscFunc(a *Node, b *Node) int {
    if a == nil || b == nil {
        return 0
    }
    if a.Key < b.Key {

        return 1

    } else if a.Key > b.Key {

        return -1

    } else {

        if a.Value < b.Value {
            return 1

        } else if a.Value > b.Value {
            return -1
        } else {
            //==0
            return 0

        }

    }
}
func (v *LinkList) SortDescFunc(a *Node, b *Node) int {
    if a == nil || b == nil {
        return 0
    }
    if a.Key < b.Key {

        return -1

    } else if a.Key > b.Key {

        return 1

    } else {

        if a.Value < b.Value {
            return -1

        } else if a.Value > b.Value {
            return 1
        } else {
            //==0
            return 0

        }

    }
}

func (v *LinkList)swap(a *Node, b *Node) (previous *Node)  {
    previous=a.Previous
    if previous==nil{
        next:=b.Next

        v.RootNode =b
        v.RootNode.Previous=nil

        v.RootNode.Next=a
        v.RootNode.Next.Previous = b
        v.RootNode.Next.Next=next

        if next!=nil{
            next.Previous=a
        }


        previous =b
    }else{


        next:=b.Next


        a.Previous=b

        b.Previous =previous
        b.Next =a
        previous.Next=b


        a.Next =next
        if a.Next!=nil{
            a.Next.Previous=a
        }

        previous = b.Previous
        //previous.Next.Previous=a


    }

    return previous
}
//todo:还没有写完
func (v *LinkList) Sort(sortFunc func(a *Node, b *Node) int) {
    //desc -1
    //asc 1
    //eque 0
    v.Last = nil
    nextNode := v.RootNode

     p := 0
    for nextNode != nil {
        p++

        if nextNode.Next != nil {

            //[0,2,3,5,6,4,8,9,60,1,23]

            sortType := sortFunc(nextNode, nextNode.Next)
            if sortType!=0{
                //nextNode=v.swap(nextNode,nextNode.Next)
                if sortType==-1{
                    nextNode=v.swap(nextNode,nextNode.Next)
                }else{
                    nextNode =nextNode.Next
                }

            }else{
                nextNode =nextNode.Next
            }


        } else {

            log.Println("count",p)
            v.Last = nextNode
            break



        }

    }

}
func (v *LinkList) Get(key string) *Node {
    if v.Map == nil {
        v.Map = make(map[string]*Node)
        return &Node{}
    }
    _, ok := v.Map[key]
    if ok {
        return v.Map[key]
    }
    return &Node{}
}

func (v *LinkList) Add(key, value string) *Node {
    if v.Map == nil {
        v.Map = make(map[string]*Node)
    }
    //item, ok := v.Map[key]

    node := &Node{Value: value, Key: key}
    v.Map[key] = node
    if v.RootNode == nil {
        v.RootNode = node
        v.Last = node
    } else if v.Last != nil {
        node.Previous = v.Last
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
