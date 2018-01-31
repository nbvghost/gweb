package gweb

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"server.local/gweb/tool"
	"reflect"

)

type Context struct {
	Response   http.ResponseWriter
	Request    *http.Request
	Session    *Session
	PathParams map[string]string
}

type Function struct {
	Function func(context *Context) Result
}

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
	RequestMapping map[string]*Function
	Context        Context
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
		if r := recover(); r != nil {
			panic("重复的path:" + path)
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
func (c *BaseController) AddHandler(pattern string, function *Function) {
	c.Lock()
	defer c.Unlock()
	if c.RequestMapping == nil {
		c.RequestMapping = make(map[string]*Function)
	}
	_pattern := c.Root +"/"+ pattern
	c.RequestMapping[delRepeatAll(_pattern, "/", "/")] = function
	//fmt.Println(c.RequestMapping)
}
func (c *BaseController) doAction(path string, context *Context) Result {

	var f *Function
	var result Result

	if strings.Contains(path, ":") == true {
		return &ErrorResult{errors.New("地址:(" + path + ")不允许包含有':'")}
	} else {
		if c.RequestMapping[path] != nil {

			//fmt.Println(path,path)
			f = c.RequestMapping[path]

		} else {

			c.Lock()
			for key, value := range c.RequestMapping {
				if su, params := matchURL(key, path); su {
					context.PathParams = params
					f = value
				}
			}
			c.Unlock()

			if f == nil {
				f = c.RequestMapping[c.Root+"*"]
				if f == nil {
					f = c.RequestMapping[c.Root]
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

	cookie, err := r.Cookie("GLSESSIONID")
	var session *Session
	var GLSESSIONID string
	if err != nil {

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

	var context = &Context{Response: w, Request: r, Session: session}

	bo := c.Interceptors.ExecuteAll(context)
	if bo == false {
		return
	}
	result := c.doAction(r.URL.Path, context)
	result.Apply(context)
}

func delRepeatAll(src string, repeat string, new string) string {
	reg := regexp.MustCompile("(" + repeat + "){2,}")
	return reg.ReplaceAllString(src, new)
}
func matchURL(RoutePath string, Path string) (bool, map[string]string) {
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
