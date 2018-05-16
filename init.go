package gweb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"fmt"
	"github.com/nbvghost/gweb/tool"
	"github.com/nbvghost/gweb/conf"
)

func init() {
	//fmt.Println(fixPath("/fg/fg/sdf/gd/fg/dsg/sd/fg/sd////sdf/g/sd/g/sd/g////sgdf/g/////sg//ds"))
	content, err := ioutil.ReadFile("gweb.json")
	if err != nil {
		tool.Trace("缺少配制文件：gweb.json")
		tool.Trace("使用默认配制：")
		conf.Config.ViewDir = "view"
		conf.Config.ResourcesDir = "resources"
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



	http.Handle("/"+conf.Config.ResourcesDir+"/", http.StripPrefix("/"+conf.Config.ResourcesDir+"/", http.FileServer(http.Dir(conf.Config.ResourcesDir))))
}
