package tool

import (
	"math"
	"sort"
	"strings"
)

type List struct {
	Collection []interface{}
	SortFunc   func(i, j int) bool
}

func (list *List) Shift() interface{} {
	old := (*list).Collection
	n := len(old)
	if n == 0 {
		return nil
	}
	x := old[0:1][0]
	(*list).Collection = old[1:]
	return x
}
func (list *List) Pop() interface{} {
	old := (*list).Collection
	n := len(old)
	if n == 0 {
		return nil
	}
	x := old[n-1].([]interface{})[0]
	(*list).Collection = old[0 : n-1]
	return x
}
func (list *List) String() string {
	var txt = ""
	for _, value := range list.Collection {
		txt = txt + "," + value.(string)
	}
	return txt
}
func (list *List) Append(elems interface{}) {

	list.Collection = append(list.Collection, elems)
}
func (list *List) Len() int {
	return len(list.Collection)
}
func (list *List) Less(i, j int) bool {
	return list.SortFunc(i, j)
}
func (list *List) Swap(i, j int) {
	var temp = list.Collection[i]
	list.Collection[i] = list.Collection[j]
	list.Collection[j] = temp
}
func (list *List) SortL() {
	list.SortFunc = func(i, j int) bool {
		return list.Collection[i].(string) < list.Collection[j].(string)
	}

	sort.Sort(list)
}
func (list *List) Join(fix string) string {
	var txt = ""
	for _, value := range list.Collection {
		if strings.EqualFold(txt, "") {
			txt = (value).(string)
		} else {
			txt = txt + fix + (value).(string)
		}

	}
	return txt
}
func (list *List) SortH() {
	list.SortFunc = func(i, j int) bool {
		a := list.Collection[i].(string)
		b := list.Collection[j].(string)
		if a[0] > b[0] {
			return true
		} else if a[0] == b[0] {
			le := int(math.Max(float64(len(a)), float64(len(b))))
			for oo := 0; oo < le; oo++ {
				var aa byte = 0
				if oo > len(a)-1 {
					aa = 0
				} else {
					aa = a[oo]
				}

				var bb byte = 0
				if oo > len(b)-1 {
					bb = 0
				} else {
					bb = b[oo]
				}
				if aa > bb {
					return true
				} else if aa < bb {
					return false
				}
			}
			return false

		} else {
			return false
		}
		//return this.Collection[i].(string) > this.Collection[j].(string)
	}
	sort.Sort(list)
}
