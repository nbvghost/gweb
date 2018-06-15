package gweb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"fmt"
	"github.com/nbvghost/gweb/tool"
	"github.com/nbvghost/gweb/conf"
	"os"
)

var tempFiles = make(map[string]int64)
func init() {
	//fmt.Println(fixPath("/fg/fg/sdf/gd/fg/dsg/sd/fg/sd////sdf/g/sd/g/sd/g////sgdf/g/////sg//ds"))
	content, err := ioutil.ReadFile("gweb.json")
	if err != nil {
		tool.Trace("缺少配制文件：gweb.json")
		tool.Trace("使用默认配制：")
		conf.Config.ViewDir = "view"
		conf.Config.UploadDir = "upload"
		conf.Config.ResourcesDir = "resources"
		conf.Config.ResourcesDirName = "resources"
		conf.Config.UploadDirName = "upload"
		conf.Config.DefaultPage = "index"
		conf.Config.HttpPort = ":80"
		conf.Config.HttpsPort = ":443"
		conf.Config.Domain ="localhost"
		conf.Config.JsonDataPath = "data.json"
		conf.Config.ViewSuffix = ".html"		
		conf.Config.ViewActionMapping = true
		conf.Config.DBUrl = ""

	} else {
		err = json.Unmarshal(content, &conf.Config)
		tool.CheckError(err)
	}

	dt, _ := json.Marshal(conf.Config)
	//tool.Trace("当前配制信息：" + string(dt))
	fmt.Printf("当前配制信息：\n%v\n",string(dt))

	go func() {
		mJsonData, err := ioutil.ReadFile(conf.Config.JsonDataPath)
		tool.CheckError(err)
		fmt.Printf("当前data.json数据：\n%v\n",string(mJsonData))

		for {
			mJsonData, err := ioutil.ReadFile(conf.Config.JsonDataPath)
			if err == nil {
				err = json.Unmarshal(mJsonData, &conf.JsonData)
				tool.CheckError(err)
			}
			time.Sleep(time.Second * 3)
		}
	}()


	go func() {

		err:=os.RemoveAll("temp")
		tool.CheckError(err)
		for{
			for k,v:=range tempFiles{
				if time.Now().Unix()>time.Unix(v,0).Add(time.Minute*3).Unix(){
					delete(tempFiles,k)

					err=os.Remove(k)
					tool.CheckError(err)
				}

			}
			time.Sleep(time.Second)
		}

	}()

	http.HandleFunc("/file/up", fileUp)
	http.HandleFunc("/file/load", fileLoad)
	http.HandleFunc("/file/temp/load", fileTempLoad)


	http.Handle("/"+conf.Config.ResourcesDirName+"/", http.StripPrefix("/"+conf.Config.ResourcesDirName+"/", http.FileServer(http.Dir(conf.Config.ResourcesDir))))
	http.Handle("/"+conf.Config.UploadDirName+"/", http.StripPrefix("/"+conf.Config.UploadDirName+"/", http.FileServer(http.Dir(conf.Config.UploadDir))))
	http.Handle("/temp/temp/", http.StripPrefix("/temp/temp/", http.FileServer(http.Dir("temp"))))
}
func fileUp(writer http.ResponseWriter, request *http.Request)  {

	request.ParseForm()
	File, FileHeader, err := request.FormFile("file")
	tool.CheckError(err)
	b, err := ioutil.ReadAll(File)
	tool.CheckError(err)
	defer File.Close()

	fileName := tool.WriteFile(b, FileHeader.Header.Get("Content-Type"))
	//base64Data := "data:" + FileHeader.Header.Get("Content-Type") + ";base64," + base64.StdEncoding.EncodeToString(b)
	result:=make(map[string]interface{})
	result["Success"]=true
	result["Message"]="OK"
	result["Data"]=fileName
	rb,_:=json.Marshal(result)
	writer.Write(rb)
	//framework.WriteJSON(context, &framework.ActionStatus{true, "oK", base64Data})
	//return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "ok", Data: fileName}}
}
func fileLoad(writer http.ResponseWriter, request *http.Request)  {
	path := request.URL.Query().Get("path")
	//fmt.Println(util.GetHost(context))
	//return &gweb.ImageResult{FilePath: path}
	//return &gweb.RedirectToUrlResult{Url:"/file/"}
	http.Redirect(writer, request,"/"+path, http.StatusFound)
}
func fileTempLoad(writer http.ResponseWriter, request *http.Request)  {
	path := request.URL.Query().Get("path")
	//fmt.Println(util.GetHost(context))
	//return &gweb.ImageResult{FilePath: path}
	//return &gweb.RedirectToUrlResult{Url:"/file/"}
	tempFiles[path]=time.Now().Unix()
	http.Redirect(writer, request,"/temp/"+path, http.StatusFound)






}