package gweb

import (
	"errors"
	"github.com/nbvghost/gweb/tool"
	"sync"
)

type Interceptors struct {
	lock  *sync.Mutex
	list []Interceptor
}
type Interceptor interface {
	Execute(Context *Context)(bool,Result)
}
func (inter *Interceptors) Add(value Interceptor) {
	if inter.lock==nil{
		inter.lock =&sync.Mutex{}
	}

	if inter.list == nil {
		inter.list = make([]Interceptor, 0)
	}

	if inter.Contains(value) == false {
		inter.list = append(inter.list, value)
	} else {
		tool.CheckError(errors.New("已经存在"))
	}
}
func (inter *Interceptors) ExecuteAll(c *BaseController) (bool,Result) {
	if inter.lock==nil{
		inter.lock =&sync.Mutex{}
	}
	inter.lock.Lock()
	defer inter.lock.Unlock()
	for _, value := range inter.list {
		//Execute(Session *Session,Request *http.Request) Result

		bo,result:= value.Execute(c.Context)
		if bo == false {
			return false,result
		}
	}
	return true,nil
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
