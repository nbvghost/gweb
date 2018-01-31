package gweb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

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
		conf.Config.JsonDataPath = "data.json"
		conf.Config.ViewSuffix = ".html"		
		conf.Config.ViewActionMapping = true

	} else {
		err = json.Unmarshal(content, &conf.Config)
		tool.CheckError(err)
	}

	dt, _ := json.Marshal(conf.Config)
	tool.Trace("当前配制信息：" + string(dt))

	go func() {
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
