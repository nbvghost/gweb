package gweb

import (
	"encoding/gob"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/conf"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"strings"
	"time"
)

func initGo() {

	go func() {
		if conf.Config.SessionExpires > 0 {
			for {
				Sessions.Range(func(key, value interface{}) bool {
					session := value.(*Session)
					if time.Now().Unix()-session.LastOperationTime >= int64(conf.Config.SessionExpires) {
						Sessions.DeleteSession(key.(string))
					}
					return true
				})
				time.Sleep(time.Second)
			}
		}

	}()

	go func() {

		//err:=os.RemoveAll("temp")
		//tool.CheckError(err)
		for {

			fileList, err := ioutil.ReadDir("temp")
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			for _, v := range fileList {
				//fmt.Println(k,v)
				//file, err := ioutil.ReadFile("temp" + "/" + v.Name())

				file, err := os.Stat("temp" + "/" + v.Name())
				if time.Now().Unix() > file.ModTime().Add(time.Minute*3).Unix() {
					err = os.Remove("temp" + "/" + v.Name() + "/" + file.Name())
					glog.Error(err)
				}

			}
			time.Sleep(time.Second)
		}

	}()
}
func init() {
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})
	//testing.Init()
	//flag.StringVar(&gwebJson, "gweb", "gweb.json", "-gweb 指定gweb.json的位置")

	initGo()

}

/*type Static struct {
}*/

/*func (static Static) fileNetLoad(writer http.ResponseWriter, request *http.Request) {
	url := request.URL.Query().Get("url")
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		writer.Write([]byte{})
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		writer.Write([]byte{})
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		writer.Write([]byte{})
		return
	}
	writer.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	writer.Write(b)
}*/

/*func (static Static) FileLoad(ctx *Context) Result {
	path := ctx.Request.URL.Query().Get("path")
	urldd, err := url.Parse(path)
	if glog.Error(err) || (strings.EqualFold(urldd.Scheme, "") && strings.EqualFold(urldd.Host, "")) {
		//dir, _ := filepath.Split(path)
		http.ServeFile(ctx.Response, ctx.Request, strings.Trim(conf.Config.UploadDir, "/")+"/"+strings.Trim(conf.Config.UploadDirName, "/")+"/"+path)
		return &EmptyResult{}
	}
	return &RedirectToUrlResult{Url: path}
}*/
/*func (static Static) FileTempLoad(ctx *Context) Result {
	path := ctx.Request.URL.Query().Get("path")
	http.ServeFile(ctx.Response, ctx.Request, "temp/"+path)
	return &EmptyResult{}
}
*/
func FileLoadAction(ctx *Context) Result {
	path := ctx.Request.URL.Query().Get("path")
	urldd, err := url.Parse(path)
	if glog.Error(err) || (strings.EqualFold(urldd.Scheme, "") && strings.EqualFold(urldd.Host, "")) {
		//dir, _ := filepath.Split(path)
		http.ServeFile(ctx.Response, ctx.Request, strings.TrimRight(conf.Config.UploadDir, "/")+"/"+strings.TrimLeft(path, "/"))
		return &EmptyResult{}
	}
	return &RedirectToUrlResult{Url: path}
}
func FileTempLoadAction(ctx *Context) Result {
	path := ctx.Request.URL.Query().Get("path")
	http.ServeFile(ctx.Response, ctx.Request, "temp/"+strings.TrimLeft(path, "/"))
	return &EmptyResult{}
}
