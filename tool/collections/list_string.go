package collections

import (
	"math"
	"sort"
	"strings"
)

type ListString struct {
	Collection []string
	SortFunc   func(i, j int) bool
}

func (list *ListString) Shift() string {
	old := (*list).Collection
	n := len(old)
	if n == 0 {
		return ""
	}
	x := old[0:1][0]
	(*list).Collection = old[1:]
	return x
}
func (list *ListString) Pop() string {
	old := (*list).Collection
	n := len(old)
	if n == 0 {
		return ""
	}
	x := old[n-1] //old[n-1].([]string)[0]
	(*list).Collection = old[0 : n-1]
	return x
}
func (list *ListString) String() string {
	var txt = ""
	for _, value := range list.Collection {
		txt = txt + "," + value
	}
	return txt
}
func (list *ListString) Append(elems string) {

	list.Collection = append(list.Collection, elems)
}
func (list *ListString) Len() int {
	return len(list.Collection)
}
func (list *ListString) Less(i, j int) bool {
	return list.SortFunc(i, j)
}
func (list *ListString) Swap(i, j int) {
	var temp = list.Collection[i]
	list.Collection[i] = list.Collection[j]
	list.Collection[j] = temp
}
func (list *ListString) SortL() {
	list.SortFunc = func(i, j int) bool {
		return list.Collection[i] < list.Collection[j]
	}

	sort.Sort(list)
}
func (list *ListString) Join(fix string) string {
	var txt = ""
	for _, value := range list.Collection {
		if strings.EqualFold(txt, "") {
			txt = value
		} else {
			txt = txt + fix + (value)
		}

	}
	return txt
}
func (list *ListString) SortH() {
	list.SortFunc = func(i, j int) bool {
		a := list.Collection[i]
		b := list.Collection[j]
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
