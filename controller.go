package gweb

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/cache"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/tool"
	"github.com/nbvghost/tool/encryption"

	"net/http"
	"reflect"
	"regexp"
	"runtime/debug"
	"strings"
	"time"
)

type Context struct {
	Response   http.ResponseWriter
	Request    *http.Request
	Session    *Session
	PathParams map[string]string
	RoutePath  string //route 的路径，相当于router根目录,请求request path remove restful path
	//Data       interface{}
	Function *Function
}

func (c *Context) Clone() Context {
	return Context{
		Response:   c.Response,
		Request:    c.Request,
		Session:    c.Session,
		PathParams: c.PathParams,
		RoutePath:  c.RoutePath,
		//Data:       c.Data,
	}
}

type ActionFunction func(context *Context) Result

type Function struct {
	Methods    []HttpMethod
	RoutePath  string
	Function   ActionFunction
	controller *Controller
}

/*func getRootPath() string {

}*/

func mapToPairs(m map[string]string) []string {
	var i int
	p := make([]string, len(m)*2)
	for k, v := range m {
		p[i] = k
		p[i+1] = v
		i += 2
	}
	return p
}

func (function *Function) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var startTime = time.Now().UnixNano()
	defer func() {
		if r := recover(); r != nil {
			glog.Trace(r)
			debug.PrintStack()
		}
		//context.Request.Method, context.Request.URL
		glog.Debug(fmt.Sprintf("%10v\t%10v\t%10v", r.Method, fmt.Sprintf("%vms", float64(time.Now().UnixNano()-startTime)/float64(time.Millisecond)), r.URL))
	}()

	var session *Session

	cookie, err := r.Cookie("GLSESSIONID")
	var GLSESSIONID string
	if err != nil || strings.EqualFold(cookie.Value, "") {

		GLSESSIONID = tool.UUID()
		http.SetCookie(w, &http.Cookie{Name: "GLSESSIONID", Value: GLSESSIONID, Path: "/", MaxAge: conf.Config.SessionExpires})
		session = &Session{Attributes: &Attributes{}, CreateTime: time.Now().Unix(), LastOperationTime: time.Now().Unix(), GLSESSIONID: GLSESSIONID}
		Sessions.AddSession(GLSESSIONID, session)

	} else {

		session = Sessions.GetSession(cookie.Value)
		if session == nil {
			session = &Session{Attributes: &Attributes{}, CreateTime: time.Now().Unix(), LastOperationTime: time.Now().Unix(), GLSESSIONID: cookie.Value}

			Sessions.AddSession(cookie.Value, session)
		}
		session.LastOperationTime = time.Now().Unix()
	}
	session.LastRequestURL = r.URL

	//c.Lock()
	w.Header().Add("Server-Name", conf.Config.Name)
	w.Header().Add("Server-Ver", conf.Config.Ver)
	//c.Unlock()

	pathParams := mux.Vars(r)
	var context = &Context{Response: w, Request: r, Session: session, PathParams: pathParams, Function: function}

	//context.RootPath=

	if function.controller != nil {

		if function.controller.Router == nil {
			context.RoutePath = function.controller.RoutePath
		} else {
			//todo:这个地方修改没有测试，原逻辑会导致路由被修改而异常
			context.RoutePath = function.controller.RoutePath

			varsRegexp := regexp.MustCompile("\\{(.*?)+\\}")
			if varsRegexp.MatchString(context.RoutePath) {
				context.RoutePath = varsRegexp.ReplaceAllStringFunc(context.RoutePath, func(s string) string {
					ss := s[1 : len(s)-1]
					return pathParams[ss]
				})
			}

			//url, err := function.controller.Router.//.Queries().URL(mapToPairs(mux.Vars(r))...)
			//url, err := function.controller.Router.Queries().URL(mapToPairs(mux.Vars(r))...)
			//if err != nil {
			//	context.RoutePath = "/" + strings.Trim(function.controller.RoutePath, "/") + "/"
			//} else {
			//	context.RoutePath = "/" + strings.Trim(url.Path, "/") + "/"
			//}
		}

		interceptor := function.controller.Interceptors.Get()

		//todo:不主动获取拦截器，必须指定
		/*controller := function.controller.ParentController
		for controller != nil {
			//interceptor=controller.Interceptors.Get(); interceptor != nil
			if interceptor == nil {
				interceptor = controller.Interceptors.Get()
			} else {
				break
			}

		}*/
		if interceptor == nil {
			function.controller.doAction(context, function).Apply(context)
		} else {
			isContinue, beforeResult := interceptor.ActionBefore(context)
			if isContinue == false {
				if beforeResult != nil {
					beforeResult.Apply(context)
				}
				return
			}

			serviceConfig := interceptor.ActionService(context)

			if serviceConfig.CacheConfig.EnableHTMLCache {

				var fullPath = context.Request.URL.Path
				if strings.EqualFold(context.Request.URL.RawQuery, "") == false {
					fullPath = fullPath + "?" + context.Request.URL.RawQuery
				}
				fullPathMd5 := encryption.Md5ByString(fullPath)
				cacheItem, err := cache.Read(fmt.Sprintf("cache/%v/%v", serviceConfig.CacheConfig.PrefixName, fullPathMd5))
				if err == nil {
					context.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
					context.Response.Write(cacheItem.Byte)
					return

				}

			}

			result := function.controller.doAction(context, function)

			interceptorResult := interceptor.ActionAfter(context, result)
			if interceptorResult == nil {
				interceptorResult = result
			}

			if serviceConfig.CacheConfig.EnableHTMLCache {

				if htmlResult, ok := interceptorResult.(*HTMLResult); ok {

					interceptorResult = &cacheHTMLResult{
						HTMLResult:  htmlResult,
						ServiceName: serviceConfig.CacheConfig.PrefixName,
					}
				}

			}

			interceptorResult.Apply(context)
		}

	} else {

		(&Controller{}).doAction(context, function).Apply(context)
		panic(errors.New("Function 无法获取 Controller"))
	}

}

type HttpMethod string

const (
	MethodGet     HttpMethod = "GET"
	MethodHead    HttpMethod = "HEAD"
	MethodPost    HttpMethod = "POST"
	MethodPut     HttpMethod = "PUT"
	MethodPatch   HttpMethod = "PATCH" // RFC 5789
	MethodDelete  HttpMethod = "DELETE"
	MethodConnect HttpMethod = "CONNECT"
	MethodOptions HttpMethod = "OPTIONS"
	MethodTrace   HttpMethod = "TRACE"
)

func NewFunction(RoutePath string, call ActionFunction, args ...HttpMethod) *Function {
	function := &Function{}
	function.Methods = args
	function.RoutePath = RoutePath
	function.Function = call
	return function
}

type IController interface {
	Init()
	DefaultHandle(context *Context) Result
	NotFoundHandler(context *Context) Result
}

var AppRouter = mux.NewRouter()

//var _ IController = (*Controller)(nil)

type Controller struct {
	RoutePath        string //定义路由的路径
	Interceptors     Interceptors
	ParentController *Controller
	Router           *mux.Router //dir
	ViewSubDir       string
}

func (c *Controller) DefaultHandle(context *Context) Result {
	panic("implement me")
}

func (c *Controller) NotFoundHandler(context *Context) Result {

	http.NotFound(context.Response, context.Request)

	return &EmptyResult{}
}

func NewStaticController(controller IController, actionName string) IController {

	path := "/" + strings.Trim(actionName, "/") + "/"
	if strings.EqualFold(actionName, "/") || strings.EqualFold(actionName, "") {
		path = "/"
	} else {
		path = "/" + strings.Trim(actionName, "/") + "/"
	}
	route := AppRouter.PathPrefix(path)
	glog.Panic(route.GetError())
	h := route.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		c := reflect.ValueOf(controller).Elem().FieldByName("Controller")
		cAddr := c.Addr()

		function := &Function{
			Methods:    []HttpMethod{MethodGet},
			RoutePath:  "*",
			Function:   controller.DefaultHandle,
			controller: cAddr.Interface().(*Controller),
		}

		function.ServeHTTP(writer, request)
	})
	glog.Panic(h.GetError())

	v := reflect.ValueOf(controller)
	RoutePathValue := v.Elem().FieldByName("RoutePath")
	RoutePathValue.SetString(path)

	//RouteValue := v.Elem().FieldByName("Router")
	//RouteValue.Set(reflect.ValueOf(route.Subrouter()))

	controller.Init()

	return controller
}

type NotFoundHandler struct {
	function *Function
}

func (h *NotFoundHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if h.function.controller.RoutePath != "/" {
		h.function.RoutePath = "/" + strings.ReplaceAll(request.URL.Path, h.function.controller.RoutePath, "")
	}
	h.function.ServeHTTP(writer, request)
}

// /
// dirName
func NewController(controller IController, dirName, viewSubDir string) IController {
	path := "/" + strings.Trim(dirName, "/")

	var routePath string

	if strings.EqualFold(path, "/") == false {
		h := AppRouter.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
			c := reflect.ValueOf(controller).Elem().FieldByName("Controller")

			function := &Function{
				Methods:    []HttpMethod{MethodGet},
				RoutePath:  path,
				Function:   controller.DefaultHandle,
				controller: c.Addr().Interface().(*Controller),
			}

			function.ServeHTTP(writer, request)
		})
		glog.Panic(h.GetError())
		routePath = path + "/"
	} else {
		routePath = "/"
	}

	route := AppRouter.PathPrefix(routePath)
	glog.Panic(route.GetError())
	router := route.Subrouter()

	c := reflect.ValueOf(controller).Elem().FieldByName("Controller")
	function := &Function{
		Methods:    []HttpMethod{MethodGet},
		Function:   controller.NotFoundHandler,
		controller: c.Addr().Interface().(*Controller),
	}

	router.NotFoundHandler = &NotFoundHandler{function: function}

	v := reflect.ValueOf(controller)
	RoutePathValue := v.Elem().FieldByName("RoutePath")
	RoutePathValue.SetString(routePath)

	RouterValue := v.Elem().FieldByName("Router")
	RouterValue.Set(reflect.ValueOf(router))

	ViewSubDirValue := v.Elem().FieldByName("ViewSubDir")
	ViewSubDirValue.Set(reflect.ValueOf(strings.Trim(viewSubDir, "/")))

	controller.Init()
	return controller

}

func (c *Controller) NewController(controller IController, actionName string) IController {

	path := "/" + strings.Trim(actionName, "/") + "/"

	h := c.Router.HandleFunc("/"+strings.Trim(actionName, "/"), func(writer http.ResponseWriter, request *http.Request) {
		c := reflect.ValueOf(controller).Elem().FieldByName("Controller")

		function := &Function{
			Methods:    []HttpMethod{MethodGet},
			RoutePath:  "*",
			Function:   controller.DefaultHandle,
			controller: c.Addr().Interface().(*Controller),
		}

		function.ServeHTTP(writer, request)
	})
	glog.Panic(h.GetError())

	route := c.Router.PathPrefix(path)
	glog.Panic(route.GetError())
	router := route.Subrouter()

	v := reflect.ValueOf(controller)

	RoutePathValue := v.Elem().FieldByName("RoutePath")
	RoutePathValue.SetString(c.RoutePath + strings.Trim(path, "/") + "/")

	RouteValue := v.Elem().FieldByName("Router")
	RouteValue.Set(reflect.ValueOf(router))

	ParentControllerValue := v.Elem().FieldByName("ParentController")
	ParentControllerValue.Set(reflect.ValueOf(c))

	ViewSubDirValue := v.Elem().FieldByName("ViewSubDir")
	ViewSubDirValue.Set(reflect.ValueOf(c.ViewSubDir))

	controller.Init()

	return controller

}

///func(context *Context) Result
func (c *Controller) AddHandler(function *Function) {
	if strings.EqualFold(function.RoutePath, "") {
		panic(errors.New("不允许有空的路由"))
		return
	}
	function.controller = c

	methods := make([]string, 0, len(function.Methods))
	for _, method := range function.Methods {
		methods = append(methods, string(method))
	}

	if len(methods) == 0 {
		methods = append(methods, string(MethodGet))
	}

	h := c.Router.Handle("/"+strings.TrimLeft(function.RoutePath, "/"), function)
	glog.Panic(h.GetError())
	h.Methods(methods...)

}

func (c *Controller) AddStaticHandler(function *Function) {
	if strings.EqualFold(function.RoutePath, "") {
		panic(errors.New("不允许有空的路由"))
		return
	}
	function.controller = c

	p := c.Router.PathPrefix("/" + strings.Trim(function.RoutePath, "/") + "/")
	glog.Panic(p.GetError())
	h := p.Handler(function)
	glog.Panic(h.GetError())
	h.Methods(http.MethodGet)
}

func (c *Controller) doAction(context *Context, f *Function) Result {

	var result Result
	if f == nil {
		result = &ViewActionMappingResult{}
	} else {
		result = f.Function(context)
		if result == nil {
			glog.Error(errors.New("Action:" + context.Request.URL.String() + "-> 返回视图类型为空"))
		}
	}

	return result
}

var removeSeparatorRegexp = regexp.MustCompile("(\\/)+")

func fixPath(path string) string {
	return removeSeparatorRegexp.ReplaceAllString(path, "/")
}
