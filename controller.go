package gweb

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/cache"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/tool/encryption"
	"runtime/debug"
	"time"

	"net/http"
	"reflect"
	"regexp"
	"strings"
)

var removeSeparatorRegexp = regexp.MustCompile("/+")

func fixPath(path string) string {
	return removeSeparatorRegexp.ReplaceAllString(path, "/")
}

type handler struct {
	call func(context *Context) Result
}

func (h *handler) Handle(context *Context) Result {
	return h.call(context)
}

type Context struct {
	Response   http.ResponseWriter //
	Request    *http.Request       //
	Session    *Session            //
	PathParams map[string]string   //
	RoutePath  string              //route 的路径，相当于router根目录,请求request path remove restful path
	Function   *Function           //
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

type IController interface {
	DefaultHandle(context *Context) Result
	NotFoundHandler(context *Context) Result
}

/*
MethodGet     = "GET"
MethodHead    = "HEAD"
MethodPost    = "POST"
MethodPut     = "PUT"
MethodPatch   = "PATCH" // RFC 5789
MethodDelete  = "DELETE"
MethodConnect = "CONNECT"
MethodOptions = "OPTIONS"
MethodTrace   = "TRACE"
*/

type IHandlerGet interface {
	IHandler
	HandleGet(context *Context) Result
}
type IHandlerPost interface {
	IHandler
	HandlePost(context *Context) Result
}
type IHandlerHead interface {
	IHandler
	HandleHead(context *Context) Result
}
type IHandlerPut interface {
	IHandler
	HandlePut(context *Context) Result
}
type IHandlerPatch interface {
	IHandler
	HandlePatch(context *Context) Result
}
type IHandlerDelete interface {
	IHandler
	HandleDelete(context *Context) Result
}
type IHandlerConnect interface {
	IHandler
	HandleConnect(context *Context) Result
}
type IHandlerOptions interface {
	IHandler
	HandleOptions(context *Context) Result
}
type IHandlerTrace interface {
	IHandler
	HandleTrace(context *Context) Result
}
type IHandler interface {
	Handle(context *Context) Result
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

type Function struct {
	RoutePath  string
	Handler    IHandler
	controller *Controller
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

		GLSESSIONID = encryption.CipherEncrypter(encryption.NewSecretKey(conf.Config.SecureKey), fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05")))
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
		}

		interceptor := function.controller.Interceptors.Get()
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
		panic(errors.New("function 无法获取 controller"))
	}

}

func NewFunction(RoutePath string, call IHandler) *Function {
	function := &Function{}
	function.RoutePath = RoutePath
	function.Handler = call
	return function
}

var AppRouter = mux.NewRouter()

var _ IController = &Controller{}

type Controller struct {
	RoutePath        string       //定义路由的路径
	Interceptors     Interceptors //
	ParentController *Controller  //
	Router           *mux.Router  //dir
	ViewSubDir       string       //
}

func (c *Controller) DefaultHandle(context *Context) Result {
	panic("implement me")
}

func (c *Controller) NotFoundHandler(context *Context) Result {

	http.NotFound(context.Response, context.Request)

	return &EmptyResult{}
}
func (c *Controller) NewController(controller IController, actionName string) IController {
	path := "/" + strings.Trim(actionName, "/") + "/"

	h := c.Router.HandleFunc("/"+strings.Trim(actionName, "/"), func(writer http.ResponseWriter, request *http.Request) {
		c := reflect.ValueOf(controller).Elem().FieldByName("Controller")

		function := &Function{
			RoutePath:  "*",
			Handler:    &handler{call: controller.DefaultHandle},
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
	return controller
}

func (c *Controller) AddHandler(routePath string, call IHandler) {
	function := NewFunction(routePath, call)
	if strings.EqualFold(function.RoutePath, "") {
		panic(errors.New("不允许有空的路由"))
		return
	}
	function.controller = c

	var methods []string
	if _, ok := call.(IHandlerGet); ok {
		methods = append(methods, http.MethodGet)
	}
	if _, ok := call.(IHandlerPost); ok {
		methods = append(methods, http.MethodPost)
	}
	if _, ok := call.(IHandlerHead); ok {
		methods = append(methods, http.MethodHead)
	}
	if _, ok := call.(IHandlerPut); ok {
		methods = append(methods, http.MethodPut)
	}
	if _, ok := call.(IHandlerPatch); ok {
		methods = append(methods, http.MethodPatch)
	}
	if _, ok := call.(IHandlerDelete); ok {
		methods = append(methods, http.MethodDelete)
	}
	if _, ok := call.(IHandlerConnect); ok {
		methods = append(methods, http.MethodConnect)
	}
	if _, ok := call.(IHandlerOptions); ok {
		methods = append(methods, http.MethodOptions)
	}
	if _, ok := call.(IHandlerTrace); ok {
		methods = append(methods, http.MethodTrace)
	}
	if len(methods) == 0 {
		methods = append(methods,
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		)
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

		switch context.Request.Method {
		case http.MethodGet:
			if handler, ok := f.Handler.(IHandlerGet); ok {
				result = handler.HandleGet(context)
			}
		case http.MethodHead:
			if handler, ok := f.Handler.(IHandlerHead); ok {
				result = handler.HandleHead(context)
			}
		case http.MethodPost:
			if handler, ok := f.Handler.(IHandlerPost); ok {
				result = handler.HandlePost(context)
			}
		case http.MethodPut:
			if handler, ok := f.Handler.(IHandlerPut); ok {
				result = handler.HandlePut(context)
			}
		case http.MethodPatch:
			if handler, ok := f.Handler.(IHandlerPatch); ok {
				result = handler.HandlePatch(context)
			}
		case http.MethodDelete:
			if handler, ok := f.Handler.(IHandlerDelete); ok {
				result = handler.HandleDelete(context)
			}
		case http.MethodConnect:
			if handler, ok := f.Handler.(IHandlerConnect); ok {
				result = handler.HandleConnect(context)
			}
		case http.MethodOptions:
			if handler, ok := f.Handler.(IHandlerOptions); ok {
				result = handler.HandleOptions(context)
			}
		case http.MethodTrace:
			if handler, ok := f.Handler.(IHandlerTrace); ok {
				result = handler.HandleTrace(context)
			}
		default:
			result = f.Handler.Handle(context)
		}

		if result == nil {
			result = f.Handler.Handle(context)
		}

		if result == nil {

			glog.Error(errors.New("Action:" + context.Request.URL.String() + "-> 返回视图类型为空"))
		}
	}

	return result
}

// NewController dirName
func NewController(controller IController, dirName, viewSubDir string) IController {
	path := "/" + strings.Trim(dirName, "/")

	var routePath string

	if strings.EqualFold(path, "/") == false {
		h := AppRouter.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
			c := reflect.ValueOf(controller).Elem().FieldByName("Controller")

			function := &Function{
				RoutePath:  path,
				Handler:    &handler{call: controller.DefaultHandle},
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
		Handler:    &handler{call: controller.DefaultHandle},
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

	return controller

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
			RoutePath:  "*",
			Handler:    &handler{call: controller.DefaultHandle},
			controller: cAddr.Interface().(*Controller),
		}

		function.ServeHTTP(writer, request)
	})
	glog.Panic(h.GetError())

	v := reflect.ValueOf(controller)
	RoutePathValue := v.Elem().FieldByName("RoutePath")
	RoutePathValue.SetString(path)
	return controller
}
