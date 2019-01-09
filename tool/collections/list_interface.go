package collections

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type ListInterface struct {
	Collection []interface{}
	SortFunc   func(i, j int) bool
}

func (list *ListInterface) Shift() interface{} {
	old := (*list).Collection
	n := len(old)
	if n == 0 {
		return ""
	}
	x := old[0:1][0]
	(*list).Collection = old[1:]
	return x
}
func (list *ListInterface) Pop() interface{} {
	old := (*list).Collection
	n := len(old)
	if n == 0 {
		return ""
	}
	x := old[n-1] //old[n-1].([]string)[0]
	(*list).Collection = old[0 : n-1]
	return x
}
func (list *ListInterface) String() string {
	var txt = ""
	for _, value := range list.Collection {
		txt = txt + "," + fmt.Sprint(value)
	}
	return txt
}
func (list *ListInterface) Append(elems interface{}) {

	list.Collection = append(list.Collection, elems)
}
func (list *ListInterface) Len() int {
	return len(list.Collection)
}
func (list *ListInterface) Less(i, j int) bool {
	return list.SortFunc(i, j)
}
func (list *ListInterface) Swap(i, j int) {
	var temp = list.Collection[i]
	list.Collection[i] = list.Collection[j]
	list.Collection[j] = temp
}
func (list *ListInterface) SortL() {
	list.SortFunc = func(i, j int) bool {
		return fmt.Sprint(list.Collection[i]) < fmt.Sprint(list.Collection[j])
	}

	sort.Sort(list)
}
func (list *ListInterface) Join(fix string) string {
	var txt = ""
	for _, value := range list.Collection {
		if strings.EqualFold(txt, "") {
			txt = fmt.Sprint(value)
		} else {
			txt = txt + fix + fmt.Sprint(value)
		}

	}
	return txt
}
func (list *ListInterface) SortH() {
	list.SortFunc = func(i, j int) bool {
		a := fmt.Sprint(list.Collection[i])
		b := fmt.Sprint(list.Collection[j])
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
	}
	sort.Sort(list)
}
