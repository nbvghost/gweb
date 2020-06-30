package tool

import (
	"github.com/nbvghost/glog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)
//Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1
//Mozilla/5.0 (Linux; Android 8.0; Pixel 2 Build/OPD3.170816.012) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Mobile Safari/537.36
//Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Safari/537.36
//Mozilla/5.0 (iPad; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.34 (KHTML, like Gecko) Version/11.0 Mobile/15A5341f Safari/604.1
func GetDeviceName(UserAgent string) string {
	UserAgent = strings.ToLower(UserAgent)
	//fmt.Println(UserAgent)
	if strings.Contains(UserAgent, "iphone") {
		//FrameworkHttp.OutHtmlFileWithPath(context,"game/web/ssc.html")
		return "iphone"
	}else if strings.Contains(UserAgent, "android") {
		//FrameworkHttp.OutHtmlFileWithPath(context,"game/web/ssc.html")
		return "android"
	}else if strings.Contains(UserAgent, "ipad") {
		//FrameworkHttp.OutHtmlFileWithPath(context,"game/web/ssc.html")
		return "ipad"
	}else if strings.Contains(UserAgent, "windows") {
		//FrameworkHttp.OutHtmlFileWithPath(context,"game/web/ssc.html")
		return "windows"
	}else{
		return "other"
	}
}
func GetIP(request *http.Request) string {
	//fmt.Println(context.Request)
	//fmt.Println(context.Request.Header.Get("X-Forwarded-For"))
	//fmt.Println(context.Request.RemoteAddr)
	//Ali-Cdn-Real-Ip


		//_IP := context.Request.Header.Get("X-Forwarded-For")

	IP := strings.Split(request.Header.Get("X-Forwarded-For"), ",")[0]
	if strings.EqualFold(IP, "") {
		text := request.RemoteAddr
		if strings.Contains(text, "::") {
			IP = "0.0.0.0"
		} else {
			IP = strings.Split(text, ":")[0]
		}
	}else{
		IP = strings.Split(IP, ":")[0]
	}

	return IP
}
func QueryParams(m url.Values) map[string]string {
	data := make(map[string]string)
	for key, value := range m {
		//util.Trace(key)
		//util.Trace(value)
		data[key] = value[0]
	}
	return data
}
/*func RequestByHeader(url string, UserAgent string, Referer string) ([]byte,error) {

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil,err
	}
	//req.Header.Add("User-Agent","Mozilla/5.0 (Linux; Android 7.0; SLA-AL00 Build/HUAWEISLA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/6.2 TBS/044109 Mobile Safari/537.36 MicroMessenger/6.6.7.1321(0x26060739) NetType/WIFI Language/zh_CN")
	if !strings.EqualFold(UserAgent, "") {
		req.Header.Add("User-Agent", UserAgent)
	}

	if !strings.EqualFold(Referer, "") {
		req.Header.Add("Referer", Referer)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	return b,err

}*/
func RequestByHeader(url string, UserAgent string, Referer string) (error,*http.Response,[]byte) {

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err,nil,nil
	}
	//req.Header.Add("User-Agent","Mozilla/5.0 (Linux; Android 7.0; SLA-AL00 Build/HUAWEISLA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/6.2 TBS/044109 Mobile Safari/537.36 MicroMessenger/6.6.7.1321(0x26060739) NetType/WIFI Language/zh_CN")
	if !strings.EqualFold(UserAgent, "") {
		req.Header.Add("User-Agent", UserAgent)
	}

	if !strings.EqualFold(Referer, "") {
		req.Header.Add("Referer", Referer)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err,nil,nil
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	return err,resp,b

}
func DownloadInternetImageTemp(url string, UserAgent string, Referer string) string {



	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	//req.Header.Add("User-Agent","Mozilla/5.0 (Linux; Android 7.0; SLA-AL00 Build/HUAWEISLA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/6.2 TBS/044109 Mobile Safari/537.36 MicroMessenger/6.6.7.1321(0x26060739) NetType/WIFI Language/zh_CN")
	if !strings.EqualFold(UserAgent, "") {
		req.Header.Add("User-Agent", UserAgent)
	}

	if !strings.EqualFold(Referer, "") {
		req.Header.Add("Referer", Referer)
	}

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	glog.Error(err)
	return WriteTempFile(b, resp.Header.Get("Content-Type"))

}
func DownloadInternetImage(url string, UserAgent string, Referer string) string {

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	//req.Header.Add("User-Agent","Mozilla/5.0 (Linux; Android 7.0; SLA-AL00 Build/HUAWEISLA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/6.2 TBS/044109 Mobile Safari/537.36 MicroMessenger/6.6.7.1321(0x26060739) NetType/WIFI Language/zh_CN")
	if !strings.EqualFold(UserAgent, "") {
		req.Header.Add("User-Agent", UserAgent)
	}

	if !strings.EqualFold(Referer, "") {
		req.Header.Add("Referer", Referer)
	}

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	glog.Error(err)
	return WriteFile(b, resp.Header.Get("Content-Type"))

}
