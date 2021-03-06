package gweb

import (
	"fmt"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"strings"
	"time"
)

var gweb = "gweb.json"

func initGo() {

	readDataFile := func() error {
		mJsonData, err := ioutil.ReadFile(conf.Config.JsonDataPath)
		if err != nil {
			return err //glog.Trace("当前未使用data.json 文件")
		} else {
			//fmt.Printf("当前data.json数据：\n%v\n", string(mJsonData))
			conf.JsonText = string(mJsonData)
			return nil
			//err = json.Unmarshal(mJsonData, &conf.JsonData)
			//glog.Error(err)
			//return err
		}
	}
	err := readDataFile()
	if err != nil {
		if strings.EqualFold(conf.JsonText, "") {
			glog.Trace("当前未使用data.json 文件")
		}
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			readDataFile()
			//time.Sleep(time.Second * 3)
		}
		ticker.Stop()
	}()

	go func() {
		if conf.Config.SessionExpires > 0 {
			for {
				Sessions.Range(func(key, value interface{}) bool {
					session := value.(*Session)
					if time.Now().Unix()-session.LastOperationTime >= conf.Config.SessionExpires {
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
				fileNodes, err := ioutil.ReadDir("temp" + "/" + v.Name())
				for _, file := range fileNodes {
					if time.Now().Unix() > file.ModTime().Add(time.Minute*3).Unix() {
						err = os.Remove("temp" + "/" + v.Name() + "/" + file.Name())
						glog.Error(err)
					}
				}
				if len(fileNodes) <= 0 {
					err = os.Remove("temp" + "/" + v.Name())
					glog.Error(err)
				}

			}
			time.Sleep(time.Second)
		}

	}()
}
func init() {
	//testing.Init()
	//flag.StringVar(&gwebJson, "gweb", "gweb.json", "-gweb 指定gweb.json的位置")
	LoadConfig(gweb)
	initGo()
}

//todo:暂时使用 LoadConfig 参数来指定 gweb.json 文件,后面改用 os.Args
func LoadConfig(gwebFile string) {

	content, err := ioutil.ReadFile(gwebFile)
	if err != nil {
		glog.Trace("缺少配制文件：gweb.json")
		glog.Trace("使用默认配制：")

		conf.Config.ViewDir = "view"
		conf.Config.ResourcesDir = "resources"
		conf.Config.ResourcesDirName = "resources"

		conf.Config.UploadDir = "upload"
		conf.Config.UploadDirName = "upload"

		conf.Config.DefaultPage = "index"
		conf.Config.HttpPort = ":80"
		conf.Config.HttpsPort = ":443"
		conf.Config.SessionExpires = 1800
		conf.Config.Domain = "localhost"
		conf.Config.JsonDataPath = "data.json"
		conf.Config.ViewSuffix = ".html"
		conf.Config.ViewActionMapping = []conf.ViewActionMapping{}
		conf.Config.DBUrl = ""

	} else {
		err = tool.JsonUnmarshal(content, &conf.Config)
		glog.Error(err)
	}

	if strings.EqualFold(conf.Config.ResourcesDir, "") {
		conf.Config.ResourcesDir = "resources"
	}
	if strings.EqualFold(conf.Config.ResourcesDirName, "") {
		conf.Config.ResourcesDirName = "resources"
	}

	if strings.EqualFold(conf.Config.UploadDir, "") {
		conf.Config.UploadDir = "upload"
	}
	if strings.EqualFold(conf.Config.UploadDirName, "") {
		conf.Config.UploadDirName = "upload"
	}
	if strings.EqualFold(conf.Config.Name, "") {
		conf.Config.Name = "default"
	}
	if strings.EqualFold(conf.Config.Ver, "") {
		conf.Config.Ver = "0.0.0"
	}

	conf.Config.ViewDir = strings.Trim(conf.Config.ViewDir, "/")
	conf.Config.ResourcesDir = strings.Trim(conf.Config.ResourcesDir, "/")
	conf.Config.ResourcesDirName = strings.Trim(conf.Config.ResourcesDirName, "/")
	conf.Config.UploadDir = strings.Trim(conf.Config.UploadDir, "/")
	conf.Config.UploadDirName = strings.Trim(conf.Config.UploadDirName, "/")
	conf.Config.DefaultPage = strings.Trim(conf.Config.DefaultPage, "/")

	dt, _ := tool.JsonMarshal(conf.Config)
	//tool.Trace("当前配制信息：" + string(dt))
	glog.Debug(fmt.Sprintf("当前配制信息：\n%v\n", string(dt)))

}
func FileUploadAction(context *Context, dynamicDirName string) {

	context.Request.ParseForm()
	File, FileHeader, err := context.Request.FormFile("file")
	if glog.Error(err) {
		result := make(map[string]interface{})
		result["Success"] = false
		result["Message"] = err
		result["Path"] = ""
		result["Url"] = ""
		rb, _ := tool.JsonMarshal(result)
		context.Response.Write(rb)
		return
	}
	defer File.Close()

	err, fileName := tool.WriteWithFile(File, FileHeader, dynamicDirName)
	if glog.Error(err) {
		result := make(map[string]interface{})
		result["Success"] = false
		result["Message"] = err
		result["Path"] = ""
		result["Url"] = ""
		rb, _ := tool.JsonMarshal(result)
		context.Response.Write(rb)
	} else {
		result := make(map[string]interface{})
		result["Success"] = true
		result["Message"] = "OK"
		result["Path"] = fileName
		result["Url"] = "//" + conf.Config.Domain + "/file/load?path=" + fileName
		rb, _ := tool.JsonMarshal(result)
		context.Response.Write(rb)
	}

}

type Static struct {
}

func (static Static) fileNetLoad(writer http.ResponseWriter, request *http.Request) {
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
	//return WriteFile(b, resp.Header.Get("Content-Type"))
	writer.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	writer.Write(b)
}
func (static Static) fileLoad(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Query().Get("path")

	urldd, err := url.Parse(path)
	if glog.Error(err) {
		writer.Write([]byte(path))
		return
	}
	if strings.EqualFold(urldd.Scheme, "") && strings.EqualFold(urldd.Host, "") {
		http.Redirect(writer, request, "/"+path, http.StatusFound)
	} else {
		http.Redirect(writer, request, path, http.StatusFound)
	}
}
func (static Static) fileTempLoad(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Query().Get("path")
	//fmt.Println(util.GetHost(context))
	//return &gweb.ImageResult{FilePath: path}
	//return &gweb.RedirectToUrlResult{Url:"/file/"}
	//tempFiles[path]=time.Now().Unix()
	http.Redirect(writer, request, "/temp/"+path, http.StatusFound)

}
