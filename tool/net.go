package tool

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"strings"
)

func QueryParams(m url.Values) map[string]string {
	data := make(map[string]string)
	for key, value := range m {
		//util.Trace(key)
		//util.Trace(value)
		data[key] = value[0]
	}
	return data
}
func DownloadInternetImage(url string,UserAgent string,Referer string) string {

	client:=http.Client{}
	req,err:=http.NewRequest("GET",url,nil)
	if err != nil {
		return ""
	}
	//req.Header.Add("User-Agent","Mozilla/5.0 (Linux; Android 7.0; SLA-AL00 Build/HUAWEISLA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/6.2 TBS/044109 Mobile Safari/537.36 MicroMessenger/6.6.7.1321(0x26060739) NetType/WIFI Language/zh_CN")
	if !strings.EqualFold(UserAgent,""){
		req.Header.Add("User-Agent",UserAgent)
	}

	if !strings.EqualFold(Referer,""){
		req.Header.Add("Referer",Referer)
	}

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	return WriteFile(b, resp.Header.Get("Content-Type"))

}