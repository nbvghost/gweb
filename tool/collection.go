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

func (this *List) String() string {
	var txt = ""
	for _, value := range this.Collection {
		txt = txt + "," + value.(string)
	}
	return txt
}
func (this *List) Append(elems interface{}) {

	this.Collection = append(this.Collection, elems)
}
func (this *List) Len() int {
	return len(this.Collection)
}
func (this *List) Less(i, j int) bool {
	return this.SortFunc(i, j)
}
func (this *List) Swap(i, j int) {
	var temp interface{} = this.Collection[i]
	this.Collection[i] = this.Collection[j]
	this.Collection[j] = temp
}
func (this *List) SortL() {
	this.SortFunc = func(i, j int) bool {
		return this.Collection[i].(string) < this.Collection[j].(string)
	}

	sort.Sort(this)
}
func (this *List) Join(fix string) string {
	var txt = ""
	for _, value := range this.Collection {
		if strings.EqualFold(txt, "") {
			txt = (value).(string)
		} else {
			txt = txt + fix + (value).(string)
		}

	}
	return txt
}
func (this *List) SortH() {
	this.SortFunc = func(i, j int) bool {
		a := this.Collection[i].(string)
		b := this.Collection[j].(string)
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
	sort.Sort(this)
}
