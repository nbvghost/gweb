package gweb

type InterceptorFlow string

const (
	InterceptorFlowBreak    InterceptorFlow = "BREAK"
	InterceptorFlowContinue InterceptorFlow = "CONTINUE"
)

type Interceptors struct {
	//lock *sync.Mutex
	//list []Interceptor
	interceptor []Interceptor
}

type CacheConfig struct {
	PrefixName      string //前缀名字，主要用于生成缓存目录的名字
	EnableHTMLCache bool   //是否启用html缓存功能
}
type ServiceConfig struct {
	CacheConfig CacheConfig
}

type Interceptor interface {
	ActionBefore(context *Context) (bool, Result)
	ActionService(context *Context) ServiceConfig
	ActionAfter(context *Context, result Result) Result
}

/*func (inter *Interceptors) Len() int {
	if inter == nil {
		return 0
	}
	if inter.list == nil {
		return 0
	}
	if len(inter.list) == 0 {
		return 0
	}
	return len(inter.list)
}*/
func (inter *Interceptors) ActionAfter(context *Context, result Result) Result {
	for i := range inter.interceptor {
		return inter.interceptor[i].ActionAfter(context, result)
	}
	return nil
}
func (inter *Interceptors) ActionService(context *Context) ServiceConfig {
	for i := range inter.interceptor {
		return inter.interceptor[i].ActionService(context)
	}
	return ServiceConfig{}
}
func (inter *Interceptors) ActionBefore(context *Context) (bool, Result) {
	for i := range inter.interceptor {
		canNext, r := inter.interceptor[i].ActionBefore(context)
		if !canNext {
			return canNext, r
		}
	}
	return true, nil
}
func (inter *Interceptors) AddInterceptor(value Interceptor) {
	inter.interceptor = append(inter.interceptor, value)
	/*
		if inter.list == nil {
			inter.list = make([]Interceptor, 0)

		}

		if inter.Contains(value) == false {
			inter.list = append(inter.list, value)
		} else {
			glog.Error(errors.New("已经存在"))
		}*/
}

//func (inter *Interceptors) ExecuteBeforeAll(c *BaseController, context *Context) (bool, Result) {
//
//	/*for _, value := range inter.list {
//		isContinue, interceptorResult := value.ExecuteBefore(context)
//		if isContinue == false {
//			return isContinue, interceptorResult
//		}
//	}*/
//	if inter.interceptor == nil {
//		return true, nil
//	}
//	return inter.interceptor.ExecuteBefore(context)
//}

//func (inter *Interceptors) ExecuteAfterAll(c *BaseController, context *Context, f *Function) Result {
//
//	if inter.interceptor == nil {
//		return nil
//	}
//
//	/*var interceptorFlow InterceptorFlow
//	var interceptorResult Result
//	for _, value := range inter.list {
//		interceptorFlow, interceptorResult = value.ExecuteAfter(context, f)
//		switch interceptorFlow {
//		case InterceptorFlowBreak:
//			if interceptorResult == nil {
//				return nil
//			} else {
//				return interceptorResult
//			}
//		case InterceptorFlowContinue:
//		default:
//			glog.Trace(fmt.Sprintf("未匹配的拦截器流转类型%v", interceptorFlow))
//		}
//
//	}
//	if interceptorResult == nil {
//		return nil
//	} else {
//		return interceptorResult
//	}*/
//}

/*func (inter *Interceptors) Contains(interceptor Interceptor) bool {
	have := false
	for _, value := range inter.list {
		if interceptor == value {
			have = true
			break
		}
	}
	return have

}
*/
