package gweb

import (
	"errors"
	"fmt"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/cache"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool/encryption"
	"net/http"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/nbvghost/gweb/tool"

	"runtime/debug"
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
	Method    string
	RoutePath string
	Function  ActionFunction
}

func GETMethod(RoutePath string, call ActionFunction) Function {
	var _function Function
	_function.Method = http.MethodGet
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func OPTMethod(RoutePath string, call ActionFunction) Function {
	var _function Function
	_function.Method = http.MethodOptions
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func HEAMethod(RoutePath string, call ActionFunction) Function {
	var _function Function
	_function.Method = http.MethodHead
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func POSMethod(RoutePath string, call ActionFunction) Function {
	var _function Function
	_function.Method = http.MethodPost
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func PUTMethod(RoutePath string, call ActionFunction) Function {
	var _function Function
	_function.Method = http.MethodPut
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func DELMethod(RoutePath string, call ActionFunction) Function {
	var _function Function
	_function.Method = http.MethodDelete
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func TRAMethod(RoutePath string, call ActionFunction) Function {
	var _function Function
	_function.Method = http.MethodTrace
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func CONMethod(RoutePath string, call ActionFunction) Function {
	var _function Function
	_function.Method = http.MethodConnect
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func ALLMethod(RoutePath string, call ActionFunction) Function {
	var _function Function
	_function.Method = "ALL"
	_function.RoutePath = RoutePath
	_function.Function = call
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
	Init()
	ServeHTTP(w http.ResponseWriter, r *http.Request, rootPath string)
	addRequestMapping(key string, f *Function) *ListMapping
}

/*type ISubController interface {
	Apply(parant *BaseSubController)
}
type BaseSubController struct {
	sync.RWMutex
	Base    *BaseController
	SubPath string
}*/
type ListMapping struct {
	_list []*Mapping
	sync.RWMutex
}
type Mapping struct {
	Key string
	F   *Function
}

func (lm *ListMapping) Range(call func(index int, e *Mapping) bool) {

	if lm == nil {
		return
	}

	for index := range lm._list {

		co := call(index, lm._list[index])
		if co == false {
			break
		}

	}

}
func (lm *ListMapping) GetByKey(Key string) *Mapping {
	if lm == nil {
		return nil
	}
	for index, value := range lm._list {
		if strings.EqualFold(value.Key, Key) {
			return lm._list[index]
		}
	}
	return nil

}
func (lm *ListMapping) Add(e *Mapping) {
	lm.Lock()
	defer lm.Unlock()
	if has := lm.GetByKey(e.Key); has != nil {
		panic(errors.New("不允许添加相同的路由:" + has.Key))
	}

	if lm._list == nil {
		lm._list = make([]*Mapping, 0)
	}

	lm._list = append(lm._list, e)

	sort.SliceStable(lm._list, func(i, j int) bool {
		e := lm._list[i]
		_e := lm._list[j]

		eRs := strings.Split(e.Key, "/")
		_eRs := strings.Split(_e.Key, "/")

		if len(eRs) > len(_eRs) {

			return true
		} else {
			return false
		}
	})
}

var controllerMap = make(map[string]interface{})

type BaseController struct {
	//Context          *Context
	RequestMapping   *ListMapping //map[string]*function
	RoutePath        string       //定义路由的路径
	Interceptors     Interceptors
	ParentController *BaseController
	//sync.RWMutex
}

func (c *BaseController) Init() {

}

func (c *BaseController) addRequestMapping(key string, f *Function) *ListMapping {
	//c.Lock()
	//defer c.Unlock()
	//c.RequestMapping[key] =f
	if c.RequestMapping == nil {
		c.RequestMapping = &ListMapping{}
	}
	c.RequestMapping.Add(&Mapping{Key: key, F: f})
	return c.RequestMapping
}

/*func (c *BaseSubController) AddHandler(pattern string, function *Function) {
	c.Base.AddHandler("/"+c.SubPath+"/"+pattern, function)
}*/

func (c *BaseController) NewController(path string, controller IController) {

	if strings.Contains(path, "//") {
		panic(errors.New("重复的//"))
		return
	}

	if strings.EqualFold(path, "/") || strings.EqualFold(path, "") {
		path = "/"
	} else {
		path = "/" + strings.Trim(path, "/") + "/"
	}
	c.RoutePath = path
	//path = fixPath(path)
	/*if !strings.EqualFold(path[len(path)-1:], "/") {

		path = path + "/"

	}*/
	if validateRoutePath(path) == false {
		return
	}
	controller.Init()
	//http.Handle(path, c)

	pathList := strings.Split(strings.Trim(path, "/"), "/")

	var lastItem map[string]interface{} = controllerMap
	for index := range pathList {

		if _, ok := lastItem[pathList[index]]; !ok {

			if index == len(pathList)-1 {
				//最后一项
				lastItem[pathList[index]] = map[string]interface{}{"": controller}
			} else {
				lastItem[pathList[index]] = make(map[string]interface{})
				lastItem = lastItem[pathList[index]].(map[string]interface{})
			}

		} else {
			if index == len(pathList)-1 {
				lastItem[pathList[index]].(map[string]interface{})[""] = controller
				//panic(errors.New("重复的路由："+pathList[index]))

			} else {
				lastItem = lastItem[pathList[index]].(map[string]interface{})
			}
		}

	}

}
func (c *BaseController) AddSubController(path string, isubc IController) {
	//subbc := &BaseController{}
	//subbc.Base = c
	//subbc.SubPath = path

	if strings.EqualFold(path, "/") || strings.EqualFold(path, "") {
		panic(errors.New("路由地址为*或空，请使用ALLMethod方法，创建function"))
		//panic(errors.New("不允许有空的路由"))
		return
	} else {
		path = strings.Trim(path, "/") + "/"
	}

	if strings.EqualFold(c.RoutePath, "/") {
		path = c.RoutePath + path
	} else {
		path = c.RoutePath + path
	}

	/*path = fixPath(c.Root + "/" + path)
	if !strings.EqualFold(path[len(path)-1:], "/") {
		path = path + "/"
	}*/

	value := reflect.Indirect(reflect.ValueOf(isubc))
	//fmt.Println(value.Interface())

	RootField := value.FieldByName("RoutePath")

	//fmt.Println(RootField)
	//fmt.Println("----")
	if RootField.Kind() == reflect.String {
		if RootField.CanSet() {
			RootField.SetString(path)
		}
	}
	//fmt.Println(isubc)
	//fmt.Println("----")
	if validateRoutePath(path) == false {
		return
	}

	isubc.Init()

	key := "Get," + path

	//log.Println(key)

	/*if c.RequestMapping[key] != nil {
		glog.Trace(key, "已经存在，将被替换成新的方法")
	}*/
	var _function Function
	_function.Method = "Get"
	_function.RoutePath = path
	_function.Function = func(context *Context) Result {

		return &ViewActionMappingResult{}
	}

	//c.RequestMapping[key] = &_function
	subMapping := isubc.addRequestMapping(key, &_function)
	subMapping.Range(func(index int, e *Mapping) bool {

		//isubc.addRequestMapping(e.Key, e.F)
		return true
	})

	//http.Handle(path, isubc)

	pathList := strings.Split(strings.TrimLeft(strings.TrimRight(path, "/"), "/"), "/")

	var lastItem map[string]interface{} = controllerMap
	for index := range pathList {

		if _, ok := lastItem[pathList[index]]; !ok {

			if index == len(pathList)-1 {
				//最后一项
				lastItem[pathList[index]] = map[string]interface{}{"": isubc}
			} else {
				lastItem[pathList[index]] = make(map[string]interface{})
				lastItem = lastItem[pathList[index]].(map[string]interface{})
			}

		} else {
			if index == len(pathList)-1 {
				//panic(errors.New("重复的路由"))
				lastItem[pathList[index]].(map[string]interface{})[""] = isubc
			} else {
				lastItem = lastItem[pathList[index]].(map[string]interface{})
			}
		}

	}

}

///func(context *Context) Result
func (c *BaseController) AddHandler(_function Function) {
	if strings.EqualFold(_function.RoutePath, "") {
		panic(errors.New("不允许有空的路由"))
		return
	}

	/*c.Lock()
	defer c.Unlock()
	if c.RequestMapping == nil {
		c.RequestMapping = make(map[string]*function)
	}*/

	if strings.EqualFold(_function.RoutePath, "*") || strings.EqualFold(_function.RoutePath, "") {
		if !strings.EqualFold(_function.Method, "ALL") {
			//panic("路由地址为*或空，请使用ALLMethod方法，创建function")

		}
		panic(errors.New("路由地址为*或空，请使用ALLMethod方法，创建function"))
		//panic(errors.New("不允许有空的路由"))
		return
	}

	var _pattern = ""

	//_function.RoutePath = strings.Trim(_function.RoutePath,"/")
	_function.RoutePath = strings.TrimLeft(_function.RoutePath, "/")

	_pattern = c.RoutePath + _function.RoutePath

	if validateRoutePath(_pattern) == false {
		return
	}
	key := _function.Method + "," + _pattern

	//log.Println(key)

	/*if c.RequestMapping[key] != nil {
		glog.Trace(key, "已经存在，将被替换成新的方法")
	}*/
	c.addRequestMapping(key, &_function)
	//c.RequestMapping[key] = &_function
	//fmt.Println(c.RequestMapping)
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
func (c *BaseController) pathParams(Method, Path string) (*Function, map[string]string) {

	var f *Function
	var p map[string]string

	if c.RequestMapping.GetByKey("ALL,"+Path) != nil {

		//fmt.Println(path,path)
		return c.RequestMapping.GetByKey("ALL," + Path).F, map[string]string{}

	} else if c.RequestMapping.GetByKey(Method+","+Path) != nil {
		return c.RequestMapping.GetByKey(Method + "," + Path).F, map[string]string{}

	} else {
		//地址包括参数的方法

		c.RequestMapping.Range(func(index int, e *Mapping) bool {

			keys := strings.Split(e.Key, ",") //[Method,Path]
			if su, params := getPathParams(string(keys[1]), Path); su {
				if strings.EqualFold(keys[0], "ALL") {
					p = params
					f = e.F
					return false
				} else if strings.EqualFold(string(keys[0]), Method) {
					p = params
					f = e.F
					return false
				}

			}

			return true
		})

		//是否有对应的路由
		/*if f == nil {
			f = c.RequestMapping["ALL,"+c.Root+"*"]
			if f == nil {
				f = c.RequestMapping["ALL,"+c.Root]
			}
		}*/
	}

	return f, p
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Debug(r.Method, r.URL)

	path, _ := filepath.Split(r.URL.Path)

	pathList := strings.Split(strings.Trim(path, "/"), "/")

	RootPathList := make([]string, 0)

	var lastItem = controllerMap
	for index := range pathList {

		if value, ok := lastItem[pathList[index]]; ok {

			lastItem = value.(map[string]interface{})
			RootPathList = append(RootPathList, pathList[index])
		} else {

			for k, _ := range lastItem {

				re, err := regexp.Compile("\\{(.*?)+\\}")
				glog.Error(err)
				if re.MatchString(k) {

					lastItem = lastItem[k].(map[string]interface{})
					RootPathList = append(RootPathList, pathList[index])
				}
			}

		}

	}

	controller, ok := lastItem[""].(IController)
	if ok == false {
		controller, ok = lastItem[""].(map[string]interface{})[""].(IController)
		if ok == false {
			controller = &BaseController{}
		}
	}
	controller.ServeHTTP(w, r, "/"+strings.Join(RootPathList, "/")+"/")

}
func (c *BaseController) ServeHTTP(w http.ResponseWriter, r *http.Request, rootPath string) {

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
	var context = &Context{Response: w, Request: r, Session: session, Data: jsonData}
	context.RootPath = rootPath

	Method := context.Request.Method

	var f *Function
	f, context.PathParams = c.pathParams(Method, context.Request.URL.Path)

	if c.Interceptors.Get() == nil {
		c.doAction(context, f).Apply(context)

	} else {
		isContinue, beforeResult := c.Interceptors.Get().ActionBefore(context)
		if isContinue == false {
			if beforeResult != nil {
				beforeResult.Apply(context)
			}
			return
		}

		serviceName := c.Interceptors.Get().ActionBeforeServiceName(context)

		if strings.EqualFold(serviceName, "") == false {
			var fullPath = context.Request.URL.Path
			if strings.EqualFold(context.Request.URL.RawQuery, "") == false {
				fullPath = fullPath + "?" + context.Request.URL.RawQuery
			}

			fullPathMd5 := encryption.Md5ByString(fullPath)
			cacheItem, err := cache.Read(fmt.Sprintf("cache/%v/%v", serviceName, fullPathMd5))
			if err == nil {
				context.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
				context.Response.Write(cacheItem.Byte)
				return

			}
		}

		result := c.doAction(context, f)

		interceptorResult := c.Interceptors.Get().ActionAfter(context, result)
		if interceptorResult == nil {
			interceptorResult = result
		}

		interceptorResult.Apply(context)
	}

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

func delRepeatAll(src string, new string) string {
	reg := regexp.MustCompile("(\\/)+")
	return reg.ReplaceAllString(src, new)
}
func validateRoutePath(RoutePath string) bool {
	re, err := regexp.Compile("^[0-9a-zA-Z_\\/\\{\\}\\.]+$")
	glog.Error(err)

	if re.MatchString(RoutePath) == false && strings.EqualFold(RoutePath, "") == false {
		//panic("路径:" + RoutePath + ":不允许含有0-9a-zA-Z/{}之外的字符")
		panic(errors.New("路径:" + RoutePath + ":不允许含有0-9a-zA-Z/{}之外的字符"))
		return false
	}
	routePaths := strings.Split(RoutePath, "/")

	rea, err := regexp.Compile("\\{[0-9a-zA-Z_]+\\}")
	glog.Error(err)
	reb, err := regexp.Compile("^\\{[0-9a-zA-Z_]+\\}$")
	glog.Error(err)

	for index := range routePaths {

		if strings.Count(routePaths[index], "{") != strings.Count(routePaths[index], "}") {
			//panic("路径:" + RoutePath + ":{或}个数不匹配")
			panic(errors.New("路径:" + RoutePath + ":{或}个数不匹配"))
			return false
		}

		if rea.MatchString(routePaths[index]) {
			if reb.MatchString(routePaths[index]) {
				continue
			} else {
				//panic("路径:" + RoutePath + "中" + routePaths[index] + "，只有一个{paramName}参数形式")
				panic(errors.New("路径:" + RoutePath + "中" + routePaths[index] + "，只有一个{paramName}参数形式"))
				return false
			}
		}
	}

	return true
}

/*
RoutePath 定义的路由
Path 用户路由
*/
func getPathParams(RoutePath string, Path string) (bool, map[string]string) {
	result := make(map[string]string)
	_RoutePath := delRepeatAll(RoutePath, "/")
	_Path := delRepeatAll(Path, "/")

	mRoutePaths := strings.Split(_RoutePath, "/")
	mPaths := strings.Split(_Path, "/")

	tr := RoutePath[len(RoutePath)-1:]
	isDirPath := strings.EqualFold(tr, "/")

	//两个目录级别要一样。
	if len(mRoutePaths) != len(mPaths) && isDirPath == false {
		return false, result
	}
	if len(mRoutePaths) > len(mPaths) {
		return false, result
	}

	re, err := regexp.Compile("\\{(.*?)+\\}")
	glog.Error(err)

	for index := range mRoutePaths {

		haveParams := re.MatchString(mRoutePaths[index])
		if haveParams {
			//有参数

			//获取地址参数
			Submatchs := re.FindAllStringSubmatch(mRoutePaths[index], -1)
			dfd := re.Split(mRoutePaths[index], -1) //不是参数的文本
			//fmt.Println("--------",mRoutePaths[index],dfd,Submatchs)

			subPath := mPaths[index]
			//var ars []string

			//顺45435435dsf吴dsf43543543dfsgdfs清sdfdsfds
			var kindex = 0
			var pIndex = 0
			for subIndex := range dfd {
				kindex = strings.Index(subPath, string(dfd[subIndex]))

				value := string(subPath[0:kindex])
				//fmt.Println("keywork",value)
				if !strings.EqualFold(value, "") {
					item := Submatchs[pIndex]
					result[item[1]] = value
					pIndex++
				}

				subPath = string(subPath[kindex+len(dfd[subIndex]):])
				//fmt.Println("++++++++++",string(dfd[subIndex]))
				//fmt.Println("//////////",subPath)

			}

			if len(subPath)-1 >= kindex {
				value := string(subPath[kindex:])
				//fmt.Println("keywork",value)
				if !strings.EqualFold(value, "") {
					item := Submatchs[pIndex]
					result[item[1]] = value
					pIndex++
				}
			}

		} else {
			//没有参数
			if !strings.EqualFold(mRoutePaths[index], mPaths[index]) {
				if isDirPath {
					if len(mRoutePaths) == index+1 {
						return true, result
					}

				}

				//return true, pathData
				return false, result
			}
		}
	}

	//fmt.Println(result)

	return true, result

}

func fixPath(path string) string {
	_path := delRepeatAll(path, "/")
	/*if strings.EqualFold(string(_path[0]),"/"){
		_path =string(_path[1:])
	}*/
	return _path
}
