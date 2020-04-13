package gweb

import (
	"errors"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/conf"
	"net/http"
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
	Data       map[string]interface{}
}

type function struct {
	Method    string
	RoutePath string
	Function  func(context *Context) Result
}

func GETMethod(RoutePath string, call func(context *Context) Result) function {
	var _function function
	_function.Method = "GET"
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func OPTMethod(RoutePath string, call func(context *Context) Result) function {
	var _function function
	_function.Method = "OPTIONS"
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func HEAMethod(RoutePath string, call func(context *Context) Result) function {
	var _function function
	_function.Method = "HEAD"
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func POSMethod(RoutePath string, call func(context *Context) Result) function {
	var _function function
	_function.Method = "POST"
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func PUTMethod(RoutePath string, call func(context *Context) Result) function {
	var _function function
	_function.Method = "PUT"
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func DELMethod(RoutePath string, call func(context *Context) Result) function {
	var _function function
	_function.Method = "DELETE"
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func TRAMethod(RoutePath string, call func(context *Context) Result) function {
	var _function function
	_function.Method = "TRACE"
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func CONMethod(RoutePath string, call func(context *Context) Result) function {
	var _function function
	_function.Method = "CONNECT"
	_function.RoutePath = RoutePath
	_function.Function = call
	return _function
}
func ALLMethod(RoutePath string, call func(context *Context) Result) function {
	var _function function
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
	Apply()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	addRequestMapping(key string,f *function)*ListMapping
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
	F *function
}

func (lm *ListMapping)Range(call func(index int,e *Mapping)bool) {

	if lm==nil{
		return
	}

	for index:=range lm._list{

		co:=call(index,lm._list[index])
		if co==false{
			break
		}

	}

}
func (lm *ListMapping)GetByKey(Key string) *Mapping {
	if lm==nil{
		return nil
	}
	for index,value:=range lm._list{
		if strings.EqualFold(value.Key,Key){
			return lm._list[index]
		}
	}
	return nil

}
func (lm *ListMapping)Add(e *Mapping)  {
	lm.Lock()
	defer lm.Unlock()
	if lm.GetByKey(e.Key)!=nil{
		panic(errors.New("不允许添加相同的路由"))
	}


	if lm._list==nil{
		lm._list =make([]*Mapping,0)
	}

	lm._list=append(lm._list,e)

	sort.SliceStable(lm._list, func(i, j int) bool {
		e:=lm._list[i]
		_e:=lm._list[j]

		eRs:=strings.Split(e.Key,"/")
		_eRs:=strings.Split(_e.Key,"/")

		if len(eRs)>len(_eRs){

			return true
		}else{
			return false
		}
	})
}
type BaseController struct {
	RequestMapping   *ListMapping //map[string]*function
	Context          *Context
	Root             string
	Interceptors     Interceptors
	ParentController *BaseController
	sync.RWMutex
}
func (c *BaseController) addRequestMapping(key string,f *function)*ListMapping{
	//c.Lock()
	//defer c.Unlock()
	//c.RequestMapping[key] =f
	if c.RequestMapping==nil{
		c.RequestMapping = &ListMapping{}
	}
	c.RequestMapping.Add(&Mapping{Key:key,F:f})
	return c.RequestMapping
}
/*func (c *BaseSubController) AddHandler(pattern string, function *Function) {
	c.Base.AddHandler("/"+c.SubPath+"/"+pattern, function)
}*/

func (c *BaseController) NewController(path string, ic IController) {

	defer func() {
		if r := recover(); r != nil {
			//_, file, line, _ := runtime.Caller(1)
			//log.Println(file, line, r)
			glog.Trace(r)
			debug.PrintStack()
		}
	}()

	if strings.EqualFold(path,"/") || strings.EqualFold(path,""){
		path="/"
	}else{
		path="/"+strings.Trim(path,"/")+"/"
	}
	c.Root = path
	//path = fixPath(path)
	/*if !strings.EqualFold(path[len(path)-1:], "/") {

		path = path + "/"

	}*/
	if validateRoutePath(path) == false {
		return
	}
	ic.Apply()
	http.Handle(path, ic)


}
func (c *BaseController) AddSubController(path string, isubc IController) {
	//subbc := &BaseController{}
	//subbc.Base = c
	//subbc.SubPath = path



	if strings.EqualFold(path,"/") || strings.EqualFold(path,""){
		panic(errors.New("路由地址为*或空，请使用ALLMethod方法，创建function"))
		//panic(errors.New("不允许有空的路由"))
		return
	}else{
		path=strings.Trim(path,"/")+"/"
	}


	if strings.EqualFold(c.Root,"/"){
		path=c.Root+path
	}else{
		path = c.Root + path
	}

	/*path = fixPath(c.Root + "/" + path)
	if !strings.EqualFold(path[len(path)-1:], "/") {
		path = path + "/"
	}*/

	value := reflect.Indirect(reflect.ValueOf(isubc))
	//fmt.Println(value.Interface())

	RootField := value.FieldByName("Root")


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

	isubc.Apply()

	key := "Get," + path

	//log.Println(key)

	/*if c.RequestMapping[key] != nil {
		glog.Trace(key, "已经存在，将被替换成新的方法")
	}*/
	var _function function
	_function.Method = "Get"
	_function.RoutePath = path
	_function.Function = func(context *Context) Result {

		return &ViewActionMappingResult{}
	}


	//c.RequestMapping[key] = &_function
	subMapping:=isubc.addRequestMapping(key,&_function)
	subMapping.Range(func(index int, e *Mapping) bool {

		c.addRequestMapping(e.Key,e.F)
		return true
	})




	//http.Handle(path, isubc)


}

///func(context *Context) Result
func (c *BaseController) AddHandler(_function function) {
	if strings.EqualFold(_function.RoutePath,""){
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

	var _pattern =""

	//_function.RoutePath = strings.Trim(_function.RoutePath,"/")
	_function.RoutePath = strings.TrimLeft(_function.RoutePath,"/")

	_pattern = c.Root +  _function.RoutePath



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
func (c *BaseController) pathParams(Method,Path string) (*function,map[string]string) {

	var f *function
	var p map[string]string

	if c.RequestMapping.GetByKey("ALL,"+Path) != nil {

		//fmt.Println(path,path)
		return c.RequestMapping.GetByKey("ALL,"+Path).F, map[string]string{}

	} else if c.RequestMapping.GetByKey(Method+","+Path) != nil {
		return c.RequestMapping.GetByKey(Method+","+Path).F, map[string]string{}

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

	return f,p
}
func (c *BaseController) doAction(context *Context,f *function) Result {
	//path := strings.TrimRight(context.Request.URL.Path, "/")
	//path := context.Request.URL.Path
	//rowUrl := context.Request.URL.String()
	glog.Debug(context.Request.Method,context.Request.URL)

	//var f *function
	var result Result

	//Method := context.Request.Method
	//f,context.PathParams = c.pathParams(Method,path)

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
func (c *BaseController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
		session = &Session{Attributes: &Attributes{}, CreateTime: time.Now().Unix(), LastOperationTime: time.Now().Unix(),  GLSESSIONID: GLSESSIONID}
		Sessions.AddSession(GLSESSIONID, session)

	} else {

		session = Sessions.GetSession(cookie.Value)
		if session == nil {
			session = &Session{Attributes: &Attributes{}, CreateTime: time.Now().Unix(), LastOperationTime: time.Now().Unix(),  GLSESSIONID: cookie.Value}

			Sessions.AddSession(cookie.Value, session)
		}
		session.LastOperationTime = time.Now().Unix()
	}
	session.LastRequestURL = r.URL


	c.Lock()
	w.Header().Add("Server-Name",conf.Config.Name)
	w.Header().Add("Server-Ver",conf.Config.Ver)
	c.Unlock()


	jsonData:=make(map[string]interface{})
	tool.JsonUnmarshal([]byte(conf.JsonText),&jsonData)
	var context = &Context{Response: w, Request: r, Session: session, Data: jsonData}
	c.Context = context

	Method := context.Request.Method

	var f *function
	f,context.PathParams = c.pathParams(Method,context.Request.URL.Path)

	if c.Interceptors.Len()>0{
		bo, executeResult := c.Interceptors.ExecuteAll(c)
		if bo == false {
			if executeResult != nil {
				executeResult.Apply(context)
			}
			return
		}
	}

	result := c.doAction(context,f)
	result.Apply(context)
}

func delRepeatAll(src string, new string) string {
	reg := regexp.MustCompile("(\\/)+")
	return reg.ReplaceAllString(src, new)
}
func validateRoutePath(RoutePath string) bool {
	re, err := regexp.Compile("^[0-9a-zA-Z_\\/\\{\\}\\.]+$")
	glog.Error(err)

	if re.MatchString(RoutePath) == false && strings.EqualFold(RoutePath,"")==false {
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

	tr:=RoutePath[len(RoutePath)-1:]
	isDirPath:=strings.EqualFold(tr,"/")

	//两个目录级别要一样。
	if len(mRoutePaths) != len(mPaths) && isDirPath==false{
		return false, result
	}
	if len(mRoutePaths) > len(mPaths){
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
				if isDirPath{
					if len(mRoutePaths)==index+1{
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

	/*fmt.Println("----------")






	//获取地址参数
	Submatchs:=re.FindAllStringSubmatch(_RoutePath,-1)
	//SubmatchIndexs:=re.FindAllStringSubmatchIndex(_RoutePath,-1)
	if len(Submatchs)==0{
		return false, pathData
	}

	//fmt.Println(Submatchs,"Submatchs")


	paths:=re.Split(_RoutePath,-1)
	//fmt.Println(paths)

	//lastEndIndex:=0
	//"sdfsd/dfd5f4ds_dsfdsf/sdf/dfdf_sd/dfsdsfds-dfdsfdf-dfdf/f"
	//"sdfsd/{dfdsfs}_dsfdsf/{DFdfd}/dfdf_{sdfdsfsdf}/{dfdsfddd}-dfds{fd}-{jk}/f"
	varNameIndex:=0
	for index:=range paths {
		//_Path=paths[index]
		dfd:=strings.Index(_Path,paths[index])
		if dfd>0{

			//fmt.Println("键=值",string(Submatchs[varNameIndex][1])+"="+string(_Path[0:dfd]))
			pathData[string(Submatchs[varNameIndex][1])]=string(_Path[0:dfd])
			varNameIndex++
		}else if dfd<0{
			return false, pathData
		}

		_Path=string(_Path[dfd+len(paths[index]):])

	}
	if strings.EqualFold(_Path,"")==false{
		//fmt.Println("---------Path------------------")
		//fmt.Println(_Path)
		//fmt.Println(varNameIndex)
		//fmt.Println(Submatchs)
		//varNameIndex++
		//fmt.Println("键=值",string(Submatchs[varNameIndex][1])+"="+string(_Path))
		if varNameIndex<=len(Submatchs)-1{
			pathData[string(Submatchs[varNameIndex][1])]=_Path
		}
	}
	return true, pathData*/
}

func fixPath(path string) string {
	_path := delRepeatAll(path, "/")
	/*if strings.EqualFold(string(_path[0]),"/"){
		_path =string(_path[1:])
	}*/
	return _path
}
