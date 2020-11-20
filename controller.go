package gweb

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/cache"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool"
	"github.com/nbvghost/gweb/tool/encryption"
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
	RootPath   string //解析后的路径
	Data       map[string]interface{}
}

func (c *Context) Clone() Context {
	return Context{
		Response:   c.Response,
		Request:    c.Request,
		Session:    c.Session,
		PathParams: c.PathParams,
		RootPath:   c.RootPath,
		Data:       c.Data,
	}
}

type ActionFunction func(context *Context) Result

type Function struct {
	Methods        []HttpMethod
	RoutePath      string
	Function       ActionFunction
	baseController *BaseController
}

func (function *Function) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			glog.Trace(r)
			debug.PrintStack()

		}
	}()

	var session *Session

	cookie, err := r.Cookie("GLSESSIONID")
	var GLSESSIONID string
	if err != nil || strings.EqualFold(cookie.Value, "") {

		GLSESSIONID = tool.UUID()
		http.SetCookie(w, &http.Cookie{Name: "GLSESSIONID", Value: GLSESSIONID, Path: "/", MaxAge: int(30 * time.Minute)})
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

	jsonData := make(map[string]interface{})
	tool.JsonUnmarshal([]byte(conf.JsonText), &jsonData)

	var context = &Context{Response: w, Request: r, Session: session, Data: jsonData, PathParams: mux.Vars(r)}
	//todo:
	//context.RootPath = rootPath

	if function.baseController != nil {
		if function.baseController.Interceptors.Get() == nil {
			function.baseController.doAction(context, function).Apply(context)

		} else {
			isContinue, beforeResult := function.baseController.Interceptors.Get().ActionBefore(context)
			if isContinue == false {
				if beforeResult != nil {
					beforeResult.Apply(context)
				}
				return
			}

			serviceConfig := function.baseController.Interceptors.Get().ActionService(context)

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

			result := function.baseController.doAction(context, function)

			interceptorResult := function.baseController.Interceptors.Get().ActionAfter(context, result)
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

		(&BaseController{}).doAction(context, function).Apply(context)
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
}

type BaseController struct {
	RoutePath        string //定义路由的路径
	Interceptors     Interceptors
	ParentController *BaseController
	//sync.RWMutex
	Route *mux.Router
}

func (c *BaseController) Init() {

}
func (c *BaseController) GetPath() string {
	return c.RoutePath
}

/*func (c *BaseSubController) AddHandler(pattern string, function *Function) {
	c.Base.AddHandler("/"+c.SubPath+"/"+pattern, function)
}*/

var AppRouter = mux.NewRouter()

func NewHandler(function *Function) {
	if strings.EqualFold(function.RoutePath, "") {
		panic(errors.New("不允许有空的路由"))
		return
	}

	methods := make([]string, 0, len(function.Methods))
	for _, method := range function.Methods {
		methods = append(methods, string(method))
	}

	if len(methods) == 0 {
		methods = append(methods, string(MethodGet))
	}

	AppRouter.PathPrefix("/" + strings.Trim(function.RoutePath, "/") + "/").Handler(function)

}

func NewController(path string, controller IController) IController {
	path = "/" + strings.Trim(path, "/")
	router := AppRouter.PathPrefix(path).Subrouter()

	v := reflect.ValueOf(controller)

	RoutePathValue := v.Elem().FieldByName("RoutePath")
	RoutePathValue.SetString(path)

	RouteValue := v.Elem().FieldByName("Route")
	RouteValue.Set(reflect.ValueOf(router))
	controller.Init()
	return controller

}
func (c *BaseController) NewController(path string, controller IController) IController {

	path = "/" + strings.Trim(path, "/")

	router := c.Route.PathPrefix(path).Subrouter()

	v := reflect.ValueOf(controller)

	RoutePathValue := v.Elem().FieldByName("RoutePath")
	RoutePathValue.SetString(path)

	RouteValue := v.Elem().FieldByName("Route")
	RouteValue.Set(reflect.ValueOf(router))

	ParentControllerValue := v.Elem().FieldByName("ParentController")
	ParentControllerValue.Set(reflect.ValueOf(c))

	controller.Init()

	return controller

}

///func(context *Context) Result
func (c *BaseController) AddHandler(function *Function) {
	if strings.EqualFold(function.RoutePath, "") {
		panic(errors.New("不允许有空的路由"))
		return
	}
	function.baseController = c

	methods := make([]string, 0, len(function.Methods))
	for _, method := range function.Methods {
		methods = append(methods, string(method))
	}

	if len(methods) == 0 {
		methods = append(methods, string(MethodGet))
	}

	c.Route.Handle("/"+strings.TrimLeft(function.RoutePath, "/"), function).Methods(methods...)
}
func (c *BaseController) doAction(context *Context, f *Function) Result {
	glog.Debug(context.Request.Method, context.Request.URL)
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

func fixPath(path string) string {
	reg := regexp.MustCompile("(\\/)+")
	return reg.ReplaceAllString(path, "/")
}
