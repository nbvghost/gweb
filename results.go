package gweb

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"time"

	"server.local/gweb/conf"
	"server.local/gweb/tool"
)

type Result interface {
	Apply(context *Context)
}

type ErrorResult struct {
	Error error
}

func (r *ErrorResult) Apply(context *Context) {
	http.Error(context.Response, r.Error.Error(), http.StatusNotFound)

}

type NotFindResult struct {
}

func (r *NotFindResult) Apply(context *Context) {

	path := context.Request.URL.Path
	b, err := ioutil.ReadFile(fixPath(conf.Config.ViewDir + "/" + path))
	if err != nil {
		http.NotFound(context.Response, context.Request)
	} else {
		t, err := template.New("default").Funcs(tool.FuncMap()).Parse(string(b))
		tool.CheckError(err)
		t.Execute(context.Response, nil)
	}

}

//不做处理，返回原 Response
type ViewResult struct {
}

func (r *ViewResult) Apply(context *Context) {

}

type HTMLResult struct {
	Name   string
	Params map[string]interface{}
}

func (r *HTMLResult) Apply(context *Context) {

	path := context.Request.URL.Path

	var b []byte
	var err error
	if strings.EqualFold(r.Name, "") {
		b, err = ioutil.ReadFile(fixPath(conf.Config.ViewDir + "/" + path + conf.Config.ViewSuffix))
	} else {
		b, err = ioutil.ReadFile(fixPath(conf.Config.ViewDir + "/" + r.Name + conf.Config.ViewSuffix))
	}

	if err != nil {

		//判断是否有默认页面
		//fmt.Println(fixPath(Config.ViewDir + "/" + path +"/"+ Config.DefaultPage))
		b, err = ioutil.ReadFile(fixPath(conf.Config.ViewDir + "/" + path + "/" + conf.Config.DefaultPage + conf.Config.ViewSuffix))
		if err != nil {
			(&NotFindResult{}).Apply(context)
			return
		}
	}

	//t, err := template.New("default").Funcs(FuncMap).Parse(string(b))
	t := template.New("default").Funcs(tool.FuncMap())
	t, err = t.Parse(string(b))
	//template.Must(t.Parse(string(b)))
	if err != nil {
		log.Println(err)
		t, err = template.New("").Parse(err.Error())
	}

	data := make(map[string]interface{})
	data["session"] = context.Session.Attributes.Map
	data["query"] = tool.QueryParams(context.Request.URL.Query())
	data["params"] = r.Params
	data["host"] = context.Request.Host
	data["time"] = time.Now().Unix() * 1000
	data["data"] = conf.JsonData[path]
	//context.Response
	context.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	context.Response.WriteHeader(http.StatusOK)

	t.Execute(context.Response, data)
}

type JsonResult struct {
	Data interface{}
}

func (r *JsonResult) Apply(context *Context) {
	var b []byte
	var err error

	b, err = json.Marshal(r.Data)

	if err != nil {
		(&ErrorResult{Error: err}).Apply(context)
		return
	}

	context.Response.WriteHeader(http.StatusOK)
	context.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	context.Response.Write(b)
}

type TextResult struct {
	Data string
}

func (r *TextResult) Apply(context *Context) {
	context.Response.WriteHeader(http.StatusOK)
	context.Response.Header().Set("Content-Type", "text/plain;charset=utf-8")
	context.Response.Write([]byte(r.Data))
}

type RedirectToUrlResult struct {
	Url string
}

func (r *RedirectToUrlResult) Apply(context *Context) {
	//context.Response.Header().Set("Location", r.Url)
	//context.Response.WriteHeader(http.StatusFound)
	//context.Response.Header().Set("Content-Type", "")
	http.Redirect(context.Response, context.Request, r.Url, http.StatusFound)
}

type ImageResult struct {
	FilePath string
}

func (r *ImageResult) Apply(context *Context) {

	file, err := os.Open(r.FilePath)
	if err != nil {
		return
	}
	defer file.Close()

	buff, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	context.Response.Write(buff)

	//context.Response.Header().Set("Location", r.Url)
	//context.Response.WriteHeader(http.StatusFound)
	//context.Response.Header().Set("Content-Type", "")

}

type ImageBytesResult struct {
	Data []byte
}

func (r *ImageBytesResult) Apply(context *Context) {

	context.Response.Write(r.Data)

}
