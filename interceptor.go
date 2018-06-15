package gweb

import (
	"errors"
	"github.com/nbvghost/gweb/tool"
	"sync"
)

type Interceptors struct {
	sync.Mutex
	list []Interceptor
	DisableManagerSession bool
}
type Interceptor interface {
	Execute(context *Context) bool
}
func (inter *Interceptors) Add(value Interceptor) {

	if inter.list == nil {
		inter.list = make([]Interceptor, 0)
	}

	if inter.Contains(value) == false {
		inter.list = append(inter.list, value)
	} else {
		tool.CheckError(errors.New("已经存在"))
	}
}
func (inter *Interceptors) ExecuteAll(c *BaseController) bool {
	for _, value := range inter.list {

		//fmt.Println(c.Context.Request.URL.Path)
		//fmt.Println(c.Root)
		inter.Lock()
		bo := value.Execute(c.Context)
		inter.Unlock()
		if bo == false {
			return false
		}
		/*ikey := strings.Split(key, "*")[0]
		if strings.Contains(path, ikey) {
			return true, value
		}*/

	}
	return true
}

func (inter *Interceptors) Contains(interceptor Interceptor) bool {
	have := false
	for _, value := range inter.list {
		if interceptor == value {
			have = true
			break
		}
	}
	return have

}
