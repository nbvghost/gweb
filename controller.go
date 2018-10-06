package gweb

import (
	"errors"
	"github.com/nbvghost/gweb/conf"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/nbvghost/gweb/tool"

	"log"
	"runtime"
	"runtime/debug"
	"time"
)

type Context struct {
	Response   http.ResponseWriter
	Request    *http.Request
	Session    *Session
	PathParams map[string]string
	Data map[string]interface{}
}

type function struct {
	Method string
	RoutePath string
	Function func(context *Context) Result
}

func GETMethod(RoutePath string,call func(context *Context) Result) function  {
	var _function function
	_function.Method="GET"
	_function.RoutePath=RoutePath
	_function.Function =call
	return _function
}
func OPTMethod(RoutePath string,call func(context *Context) Result) function  {
	var _function function
	_function.Method="OPTIONS"
	_function.RoutePath=RoutePath
	_function.Function =call
	return _function
}
func HEAMethod(RoutePath string,call func(context *Context) Result) function  {
	var _function function
	_function.Method="HEAD"
	_function.RoutePath=RoutePath
	_function.Function =call
	return _function
}
func POSMethod(RoutePath string,call func(context *Context) Result) function  {
	var _function function
	_function.Method="POST"
	_function.RoutePath=RoutePath
	_function.Function =call
	return _function
}
func PUTMethod(RoutePath string,call func(context *Context) Result) function  {
	var _function function
	_function.Method="PUT"
	_function.RoutePath=RoutePath
	_function.Function =call
	return _function
}
func DELMethod(RoutePath string,call func(context *Context) Result) function  {
	var _function function
	_function.Method="DELETE"
	_function.RoutePath=RoutePath
	_function.Function =call
	return _function
}
func TRAMethod(RoutePath string,call func(context *Context) Result) function  {
	var _function function
	_function.Method="TRACE"
	_function.RoutePath=RoutePath
	_function.Function =call
	return _function
}
func CONMethod(RoutePath string,call func(context *Context) Result) function  {
	var _function function
	_function.Method="CONNECT"
	_function.RoutePath=RoutePath
	_function.Function =call
	return _function
}
func ALLMethod(RoutePath string,call func(context *Context) Result) function  {
	var _function function
	_function.Method="ALL"
	_function.RoutePath=RoutePath
	_function.Function =call
	return _function
}
/*"OPTIONS"                ; Section 9.2
| "GET"                    ; Section 9.3
| "HEAD"                   ; Section 9.4
| "POST"                   ; Section 9.5
| "PUT"                    ; Section 9.6
| "DELETE"                 ; Section 9.7
| "TRACE"                  ; Section 9.8
| "CONNECT"                ; Section 9.9*/
type IController interface {
	Apply()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
/*type ISubController interface {
	Apply(parant *BaseSubController)
}
type BaseSubController struct {
	sync.RWMutex
	Base    *BaseController
	SubPath string
}*/
type BaseController struct {
	sync.RWMutex
	RequestMapping map[string]*function
	Context        *Context
	Root           string
	Interceptors   Interceptors
	ParentController *BaseController
}


/*func (c *BaseSubController) AddHandler(pattern string, function *Function) {
	c.Base.AddHandler("/"+c.SubPath+"/"+pattern, function)
}*/
func (c *BaseController) NewController(path string, ic IController) {
	c.Root = path
	defer func() {
		if r := recover(); r != nil{
			_, file, line, _ := runtime.Caller(1)
			log.Println(file, line, r)
		}
	}()

	path = fixPath(path)

	if !strings.EqualFold(path[len(path)-1:], "/") {

		path = path + "/"

	}
	ic.Apply()
	http.Handle(path,ic)

}
func (c *BaseController) AddSubController(path string, isubc IController) {
	//subbc := &BaseController{}
	//subbc.Base = c
	//subbc.SubPath = path

	path = fixPath(c.Root+"/"+path)

	if !strings.EqualFold(path[len(path)-1:], "/") {

		path = path + "/"

	}

	value:=reflect.Indirect(reflect.ValueOf(isubc))
	//fmt.Println(value.Interface())

	RootField:=value.FieldByName("Root")

	//fmt.Println(RootField)
	//fmt.Println("----")
	if RootField.Kind()==reflect.String{
		if RootField.CanSet(){
			RootField.SetString(path)
		}
	}
	//fmt.Println(isubc)
	//fmt.Println("----")

	isubc.Apply()
	http.Handle(path,isubc)
}
///func(context *Context) Result
func (c *BaseController) AddHandler(_function function) {
	c.Lock()
	defer c.Unlock()
	if c.RequestMapping == nil {
		c.RequestMapping = make(map[string]*function)
	}
	if strings.EqualFold(_function.RoutePath,"*") || strings.EqualFold(_function.RoutePath,""){
		if !strings.EqualFold(_function.Method,"ALL"){
			panic("路由地址为*或空，请使用ALLMethod方法，创建function")
		}
	}

	_pattern := c.Root +"/"+ _function.RoutePath
	key:=_function.Method+","+delRepeatAll(_pattern, "/", "/")

	if c.RequestMapping[key]!=nil{
		tool.Trace(key,"已经存在，将被替换成新的方法")
	}
	c.RequestMapping[key] = &_function
}
//func (c *BaseController) AddHandler(pattern string, function *Function) {
//	c.Lock()
//	defer c.Unlock()
//	if c.RequestMapping == nil {
//		c.RequestMapping = make(map[string]*Function)
//	}
//	_pattern := c.Root +"/"+ pattern
//	c.RequestMapping[delRepeatAll(_pattern, "/", "/")] = function
//}
func (c *BaseController) doAction(path string, context *Context) Result {

	var f *function
	var result Result
	Method := context.Request.Method

	if strings.Contains(path, ":") == true  ||strings.Contains(path, ",") == true {
		return &ErrorResult{errors.New("地址:(" + path + ")不允许包含有':'")}
	} else {
		if c.RequestMapping["ALL,"+path] != nil {

			//fmt.Println(path,path)
			f = c.RequestMapping["ALL,"+path]

		}else if c.RequestMapping[Method+","+path] != nil {
			f = c.RequestMapping[Method+","+path]

		}else {
			//地址包括参数的方法
			c.Lock()
			for key, value := range c.RequestMapping {
				keys:=strings.Split(key,",")//[Method,Path]
				if su, params := getPathParams(string(keys[1]), path); su {
					if strings.EqualFold(keys[0],"ALL"){
						context.PathParams = params
						f = value
						break
					}else if strings.EqualFold(string(keys[0]),Method){
						context.PathParams = params
						f = value
						break
					}

				}
			}
			c.Unlock()

			if f == nil {
				f = c.RequestMapping["ALL,"+c.Root+"*"]
				if f == nil {
					f = c.RequestMapping["ALL,"+c.Root]
				}
			}
		}
	}

	if f == nil {
		result = &NotFindResult{}
	} else {
		result = f.Function(context)
		if result == nil {
			tool.CheckError(errors.New("Action:" + path + "-> 返回视图类型为空"))
		}
	}

	return result
}
func (c *BaseController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			tool.Trace(r)
			debug.PrintStack()

		}
	}()

	var session *Session

	cookie, err := r.Cookie("GLSESSIONID")
	var GLSESSIONID string
	if err != nil || strings.EqualFold(cookie.Value,"") {

		GLSESSIONID = tool.UUID()
		http.SetCookie(w, &http.Cookie{Name: "GLSESSIONID", Value: GLSESSIONID, Path: "/"})
		session = &Session{Attributes: &Attributes{Map: make(map[string]interface{})}, CreateTime: time.Now().Unix(), Operation: time.Now().Unix(), ActionTime: time.Now().Unix(), GLSESSIONID: GLSESSIONID}
		Sessions.addSession(GLSESSIONID, session)

	} else {

		session = Sessions.GetSession(cookie.Value)
		if session == nil {
			session = &Session{Attributes: &Attributes{Map: make(map[string]interface{})}, CreateTime: time.Now().Unix(), Operation: time.Now().Unix(), ActionTime: time.Now().Unix(), GLSESSIONID: cookie.Value}

			Sessions.addSession(cookie.Value, session)
		}
		session.ActionTime = time.Now().Unix()
	}
	session.LastRequestURL = r.URL

	var context = &Context{Response: w, Request: r, Session: session,Data:conf.JsonData}
	c.Context = context
	bo,result := c.Interceptors.ExecuteAll(c)
	if bo == false {
		if result!=nil{
			result.Apply(context)
		}
		return
	}
	result = c.doAction(r.URL.Path, context)
	result.Apply(context)
}

func delRepeatAll(src string, repeat string, new string) string {
	reg := regexp.MustCompile("(" + repeat + "){2,}")
	return reg.ReplaceAllString(src, new)
}
func getPathParams(RoutePath string, Path string) (bool, map[string]string) {
	_RoutePath := delRepeatAll(RoutePath, "/", "/")
	_Path := delRepeatAll(Path, "/", "/")

	mRoutePaths := strings.Split(_RoutePath, "/")
	mPaths := strings.Split(_Path, "/")

	mapData := make(map[string]string)

	if len(mRoutePaths) != len(mPaths) {
		return false, nil
	}

	for index, value := range mRoutePaths {

		if strings.Contains(value, ":") {
			mapData[value[1:]] = mPaths[index]
		} else {
			if strings.EqualFold(mRoutePaths[index], mPaths[index]) {

			} else {
				return false, nil
			}
		}

	}

	return true, mapData
}

func fixPath(path string) string {
	_path:=delRepeatAll(path, "/", "/")

	/*if strings.EqualFold(string(_path[0]),"/"){
		_path =string(_path[1:])
	}*/

	return _path
}
