package gweb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/gweb/cache"
	"github.com/nbvghost/gweb/tool/encryption"

	"html/template"
	"net/http/httptest"

	"github.com/nbvghost/glog"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"time"

	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool"
)

var _ Result = (*ErrorResult)(nil)
var _ Result = (*SingleHostReverseProxyResult)(nil)
var _ Result = (*SingleHostForwardProxyResult)(nil)
var _ Result = (*ViewActionMappingResult)(nil)
var _ Result = (*ViewResult)(nil)
var _ Result = (*EmptyResult)(nil)
var _ Result = (*cacheHTMLResult)(nil)
var _ Result = (*HTMLResult)(nil)
var _ Result = (*JsonResult)(nil)
var _ Result = (*FileServerResult)(nil)
var _ Result = (*HtmlPlainResult)(nil)
var _ Result = (*TextResult)(nil)
var _ Result = (*JavaScriptResult)(nil)
var _ Result = (*XMLResult)(nil)
var _ Result = (*RedirectToUrlResult)(nil)
var _ Result = (*ImageResult)(nil)
var _ Result = (*ImageBytesResult)(nil)

type Result interface {
	Apply(context *Context)
}

type ErrorResult struct {
	Error error
}

func NewErrorResult(err error) *ErrorResult {
	return &ErrorResult{Error: err}
}

func (r *ErrorResult) Apply(context *Context) {
	if r.Error != nil {
		http.Error(context.Response, r.Error.Error(), http.StatusNotFound)
	} else {
		http.Error(context.Response, "error", http.StatusNotFound)
	}
}

type SingleHostReverseProxyResult struct {
	Target *url.URL
}

func (r *SingleHostReverseProxyResult) Apply(context *Context) {
	rp := httputil.NewSingleHostReverseProxy(r.Target)
	rp.ServeHTTP(context.Response, context.Request)
}

type SingleHostForwardProxyResult struct {
	Target *url.URL
}

func (r *SingleHostForwardProxyResult) Apply(context *Context) {

	transport := http.DefaultTransport

	// step 1Forward Proxy
	//outReq := new(http.Request)
	//*outReq = *context.Request // this only does shallow copies of maps
	//fmt.Printf("Received request %s %s %s\n", context.Request.Method, context.Request.URL.String(), r.Target.String())
	outReq, _ := http.NewRequest(context.Request.Method, r.Target.String(), context.Request.Body)

	if clientIP, _, err := net.SplitHostPort(context.Request.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// step 2
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		context.Response.WriteHeader(http.StatusBadGateway)

	} else {
		// step 3
		for key, value := range res.Header {
			for _, v := range value {
				context.Response.Header().Add(key, v)
			}
		}
		context.Response.WriteHeader(res.StatusCode)
		io.Copy(context.Response, res.Body)
		res.Body.Close()
	}

}

type MIME string

const (
	MultipartByteranges MIME = "multipart/byteranges"
	MultipartFormData   MIME = "multipart/form-data"

	AudioWave   MIME = "audio/wave"
	AudioWav    MIME = "audio/wav"
	AudioXWav   MIME = "audio/x-wav"
	AudioWPnWav MIME = "audio/x-pn-wav"
	AudioWebm   MIME = "audio/webm"
	AudioOgg    MIME = "audio/ogg"
	AudioMpeg   MIME = "audio/mpeg"

	VideoWebm MIME = "video/webm"
	VideoOgg  MIME = "video/ogg"
	VideoMp4  MIME = "video/mp4"

	ApplicationOgg  MIME = "application/ogg"
	ApplicationJson MIME = "application/json"

	ApplicationJavascript  MIME = "application/javascript"
	ApplicationEcmascript  MIME = "application/ecmascript"
	ApplicationOctetStream MIME = "application/octet-stream"

	ImageGif    MIME = "image/gif"
	ImageJpeg   MIME = "image/jpeg"
	ImagePng    MIME = "image/png"
	ImageSvgXml MIME = "image/svg+xml"

	TextCss   MIME = "text/css"
	TextHtml  MIME = "text/html"
	TextPlain MIME = "text/plain"
)

/**
类型/子类型	扩展名
application/json						json
application/envoy						evy
application/fractals					fif
application/futuresplash				spl
application/hta	hta
application/internet-property-stream	acx
application/mac-binhex40	hqx
application/msword	doc
application/msword	dot
application/octet-stream	*
application/octet-stream	bin
application/octet-stream	class
application/octet-stream	dms
application/octet-stream	exe
application/octet-stream	lha
application/octet-stream	lzh
application/oda	oda
application/olescript	axs
application/pdf	pdf
application/pics-rules	prf
application/pkcs10	p10
application/pkix-crl	crl
application/postscript	ai
application/postscript	eps
application/postscript	ps
application/rtf	rtf
application/set-payment-initiation	setpay
application/set-registration-initiation	setreg
application/vnd.ms-excel	xla
application/vnd.ms-excel	xlc
application/vnd.ms-excel	xlm
application/vnd.ms-excel	xls
application/vnd.ms-excel	xlt
application/vnd.ms-excel	xlw
application/vnd.ms-outlook	msg
application/vnd.ms-pkicertstore	sst
application/vnd.ms-pkiseccat	cat
application/vnd.ms-pkistl	stl
application/vnd.ms-powerpoint	pot
application/vnd.ms-powerpoint	pps
application/vnd.ms-powerpoint	ppt
application/vnd.ms-project	mpp
application/vnd.ms-works	wcm
application/vnd.ms-works	wdb
application/vnd.ms-works	wks
application/vnd.ms-works	wps
application/winhlp	hlp
application/x-bcpio	bcpio
application/x-cdf	cdf
application/x-compress	z
application/x-compressed	tgz
application/x-cpio	cpio
application/x-csh	csh
application/x-director	dcr
application/x-director	dir
application/x-director	dxr
application/x-dvi	dvi
application/x-gtar	gtar
application/x-gzip	gz
application/x-hdf	hdf
application/x-internet-signup	ins
application/x-internet-signup	isp
application/x-iphone	iii
application/x-javascript	js
application/x-latex	latex
application/x-msaccess	mdb
application/x-mscardfile	crd
application/x-msclip	clp
application/x-msdownload	dll
application/x-msmediaview	m13
application/x-msmediaview	m14
application/x-msmediaview	mvb
application/x-msmetafile	wmf
application/x-msmoney	mny
application/x-mspublisher	pub
application/x-msschedule	scd
application/x-msterminal	trm
application/x-mswrite	wri
application/x-netcdf	cdf
application/x-netcdf	nc
application/x-perfmon	pma
application/x-perfmon	pmc
application/x-perfmon	pml
application/x-perfmon	pmr
application/x-perfmon	pmw
application/x-pkcs12	p12
application/x-pkcs12	pfx
application/x-pkcs7-certificates	p7b
application/x-pkcs7-certificates	spc
application/x-pkcs7-certreqresp	p7r
application/x-pkcs7-mime	p7c
application/x-pkcs7-mime	p7m
application/x-pkcs7-signature	p7s
application/x-sh	sh
application/x-shar	shar
application/x-shockwave-flash	swf
application/x-stuffit	sit
application/x-sv4cpio	sv4cpio
application/x-sv4crc	sv4crc
application/x-tar	tar
application/x-tcl	tcl
application/x-tex	tex
application/x-texinfo	texi
application/x-texinfo	texinfo
application/x-troff	roff
application/x-troff	t
application/x-troff	tr
application/x-troff-man	man
application/x-troff-me	me
application/x-troff-ms	ms
application/x-ustar	ustar
application/x-wais-source	src
application/x-x509-ca-cert	cer
application/x-x509-ca-cert	crt
application/x-x509-ca-cert	der
application/ynd.ms-pkipko	pko
application/zip	zip
audio/basic	au
audio/basic	snd
audio/mid	mid
audio/mid	rmi
audio/mpeg	mp3
audio/x-aiff	aif
audio/x-aiff	aifc
audio/x-aiff	aiff
audio/x-mpegurl	m3u
audio/x-pn-realaudio	ra
audio/x-pn-realaudio	ram
audio/x-wav	wav
image/bmp	bmp
image/cis-cod	cod
image/gif	gif
image/ief	ief
image/jpeg	jpe
image/jpeg	jpeg
image/jpeg	jpg
image/pipeg	jfif
image/svg+xml	svg
image/tiff	tif
image/tiff	tiff
image/x-cmu-raster	ras
image/x-cmx	cmx
image/x-icon	ico
image/x-portable-anymap	pnm
image/x-portable-bitmap	pbm
image/x-portable-graymap	pgm
image/x-portable-pixmap	ppm
image/x-rgb	rgb
image/x-xbitmap	xbm
image/x-xpixmap	xpm
image/x-xwindowdump	xwd
message/rfc822	mht
message/rfc822	mhtml
message/rfc822	nws
text/css	css
text/h323	323
text/html	htm
text/html	html
text/html	stm
text/iuls	uls
text/plain	bas
text/plain	c
text/plain	h
text/plain	txt
text/richtext	rtx
text/scriptlet	sct
text/tab-separated-values	tsv
text/webviewhtml	htt
text/x-component	htc
text/x-setext	etx
text/x-vcard	vcf
video/mpeg	mp2
video/mpeg	mpa
video/mpeg	mpe
video/mpeg	mpeg
video/mpeg	mpg
video/mpeg	mpv2
video/quicktime	mov
video/quicktime	qt
video/x-la-asf	lsf
video/x-la-asf	lsx
video/x-ms-asf	asf
video/x-ms-asf	asr
video/x-ms-asf	asx
video/x-msvideo	avi
video/x-sgi-movie	movie
x-world/x-vrml	flr
x-world/x-vrml	vrml
x-world/x-vrml	wrl
x-world/x-vrml	wrz
x-world/x-vrml	xaf
x-world/x-vrml	xof

*/
type ViewActionMappingResult struct {
}

func (r *ViewActionMappingResult) Apply(context *Context) {

	path := context.Request.URL.Path

	if strings.EqualFold(path, "/") {
		if strings.EqualFold(conf.Config.DefaultPage, "") == false {
			path = path + conf.Config.DefaultPage
			var redirectToUrlResult = &RedirectToUrlResult{Url: path}
			redirectToUrlResult.Apply(context)
			return
		}

	}

	path = strings.TrimRight(path, "/")
	//b, err := ioutil.ReadFile(conf.Config.ViewDir + path + conf.Config.ViewSuffix)
	b, err := cache.Read(conf.Config.ViewDir + path + conf.Config.ViewSuffix)
	if err != nil {
		//不存在
		//fmt.Println(context.Request.Header)

		var haveMIME = false
		//b, err := ioutil.ReadFile(conf.Config.ViewDir + path)
		b, err := cache.Read(conf.Config.ViewDir + path)
		if err == nil {
			re, err := regexp.Compile("\\/([0-9a-zA-Z_]+)\\.([0-9a-zA-Z]+)$")
			glog.Error(err)

			if re.MatchString(path) {
				Groups := re.FindAllStringSubmatch(path, -1)
				//[[/fgsd_gffdgdf.txt fgsd_gffdgdf txt]]
				//{"ContentType": "text/html","Extension":"html"}
				Extension := Groups[0][2]
				for index := range conf.Config.ViewActionMapping {
					ce := conf.Config.ViewActionMapping[index]
					if strings.EqualFold(ce.Extension, Extension) {
						context.Response.Header().Set("Content-Type", ce.ContentType+"; charset=utf-8")
						//w.Header().Set("X-Content-Type-Options", "nosniff")
						context.Response.WriteHeader(http.StatusOK)
						context.Response.Write(b.Byte)
						haveMIME = true
						break
					}

				}
			}

		}
		if haveMIME == false {

			fi, err := os.Stat(conf.Config.ViewDir + path)
			//log.Println(err)
			if err == nil && fi.IsDir() {

				path = path + "/" + conf.Config.DefaultPage
				var redirectToUrlResult = &RedirectToUrlResult{Url: path}
				redirectToUrlResult.Apply(context)

			} else {
				//没有找到路由，
				http.NotFound(context.Response, context.Request)
			}

		}

	} else {
		context.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
		context.Response.WriteHeader(http.StatusOK)
		t, err := template.New("default").Funcs(NewFuncMap(context)).Parse(string(b.Byte))
		glog.Error(err)

		data := make(map[string]interface{})
		data["session"] = context.Session.Attributes.GetMap()
		data["query"] = tool.QueryParams(context.Request.URL.Query())

		glog.Error(t.Execute(context.Response, data))
	}

}

/*type NotFindResult struct {
}

func (r *NotFindResult) Apply(context *Context) {

	path := context.Request.URL.Path
	b, err := ioutil.ReadFile(fixPath(conf.Config.ViewDir + "/" + path))
	if err != nil {
		//没有找到路由，

		http.NotFound(context.Response, context.Request)




	} else {
		t, err := template.New("default").Funcs(tool.FuncMap()).Parse(string(b))
		tool.CheckError(err)
		t.Execute(context.Response, nil)
	}

}*/

//不做处理，返回原 Response
type ViewResult struct {
}

func (r *ViewResult) Apply(context *Context) {

}

type EmptyResult struct {
}

func (r *EmptyResult) Apply(context *Context) {

}

//只映射已经定义的后缀模板文件，并生成html缓存文件
type cacheHTMLResult struct {
	*HTMLResult
	ServiceName string
}

func (r *cacheHTMLResult) Apply(context *Context) {
	if r.ServiceName == "" {
		NewErrorResult(errors.New("CacheHTMLResult 结果，必须指定ServiceName值")).Apply(context)
		return
	}
	responseRecorder := httptest.NewRecorder()
	copyContext := context.Clone()
	copyContext.Response = responseRecorder
	r.HTMLResult.Apply(&copyContext)

	context.Response.WriteHeader(responseRecorder.Code)
	for key := range responseRecorder.Header() {
		context.Response.Header().Set(key, responseRecorder.Header().Get(key))
	}
	dataByte, err := ioutil.ReadAll(responseRecorder.Body)
	if glog.Error(err) {
		NewErrorResult(err).Apply(context)
		return
	}

	//path, filename := filepath.Split(context.Request.URL.Path)
	//path := context.Request.URL.Path
	var fullPath = context.Request.URL.Path
	if strings.EqualFold(context.Request.URL.RawQuery, "") == false {
		fullPath = fullPath + "?" + context.Request.URL.RawQuery
	}

	fullPathMd5 := encryption.Md5ByString(fullPath)

	cacheDir := fmt.Sprintf("cache/%v", r.ServiceName)
	cacheFile := cacheDir + "/" + fullPathMd5

	if tool.IsFileExist(cacheDir) == false {
		glog.Error(os.MkdirAll(cacheDir, os.ModePerm))
	}

	rp := regexp.MustCompile(`\s{2,}`)
	dataByte = rp.ReplaceAll(dataByte, []byte(" "))

	rp = regexp.MustCompile(`[\r\n]`)
	dataByte = rp.ReplaceAll(dataByte, []byte{})

	glog.Error(ioutil.WriteFile(cacheFile, dataByte, os.ModePerm))
	context.Response.Write(dataByte)
}

//只映射已经定义的后缀模板文件
type HTMLResult struct {
	Name       string
	StatusCode int
	Params     map[string]interface{}
	Template   []string //读取当前目录下 template 文件夹下的模板
}

func (r *HTMLResult) Apply(context *Context) {
	path, filename := filepath.Split(context.Request.URL.Path)

	var b *cache.CacheFileItem
	var err error

	viewSubDir := context.Function.controller.ViewSubDir
	if strings.EqualFold(viewSubDir, "") == false {
		viewSubDir = viewSubDir + "/"
	}

	if strings.EqualFold(r.Name, "") {
		//html 只处理，已经定义后缀名的文件
		//b, err = ioutil.ReadFile(fixPath(conf.Config.ViewDir + "/" + path + conf.Config.ViewSuffix))
		b, err = cache.Read(fixPath(conf.Config.ViewDir + "/" + viewSubDir + path + "/" + filename + conf.Config.ViewSuffix))
	} else {
		//b, err = ioutil.ReadFile(fixPath(conf.Config.ViewDir + "/" + r.Name + conf.Config.ViewSuffix))
		b, err = cache.Read(fixPath(conf.Config.ViewDir + "/" + viewSubDir + context.RoutePath + "/" + r.Name + conf.Config.ViewSuffix))
	}
	if err != nil {
		//判断是否有默认页面
		//fmt.Println(fixPath(Config.ViewDir + "/" + path +"/"+ Config.DefaultPage))
		//b, err = ioutil.ReadFile(fixPath(conf.Config.ViewDir + "/" + path + "/" + conf.Config.DefaultPage + conf.Config.ViewSuffix))
		b, err = cache.Read(fixPath(conf.Config.ViewDir + "/" + viewSubDir + path + "/" + conf.Config.DefaultPage + conf.Config.ViewSuffix))
		if err != nil {
			(&ViewActionMappingResult{}).Apply(context)
			return
		}
	}

	//t, err := template.New("default").Funcs(FuncMap).Parse(string(b))
	t := template.New("HTMLResult").Funcs(NewFuncMap(context))

	for index := range r.Template {

		tpath := r.Template[index]
		var err error
		var tt *template.Template
		if strings.Contains(tpath, "*") {
			tt, err = t.ParseGlob(conf.Config.ViewDir + "/" + viewSubDir + tpath)
		} else {
			tt, err = t.ParseFiles(conf.Config.ViewDir + "/" + viewSubDir + tpath)
		}

		if err != nil {
			glog.Trace(err)
		} else {
			t = tt
		}
	}

	t, err = t.Parse(string(b.Byte))
	//template.Must(t.Parse(string(b)))
	if glog.Error(err) {
		t, err = template.New("HTMLResult").Parse(err.Error())
	}
	data := createPageParams(context, r.Params)
	context.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.StatusCode == 0 {
		context.Response.WriteHeader(http.StatusOK)
	} else {
		context.Response.WriteHeader(r.StatusCode)
	}

	glog.Error(t.Execute(context.Response, data))
}
func createPageParams(context *Context, Params map[string]interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	data["session"] = context.Session.Attributes.GetMap()
	data["query"] = tool.QueryParams(context.Request.URL.Query())
	data["params"] = Params
	data["debug"] = conf.Config.Debug
	data["host"] = context.Request.Host
	data["time"] = time.Now().Unix() * 1000
	data["rootPath"] = context.RoutePath

	jsonData := make(map[string]interface{})
	json.Unmarshal([]byte(conf.JsonText), &jsonData)

	data["data"] = jsonData

	return data
	//context.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	//context.Response.WriteHeader(http.StatusOK)
	//t.Execute(context.Response, data)
}

type JsonResult struct {
	Data interface{}
	///sync.RWMutex
}

/*func (r *JsonResult)encodeJson() (error,[]byte)  {
	r.Lock()
	defer r.Unlock()
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(r.Data)
	return err,buffer.Bytes()
}*/
func (r *JsonResult) Apply(context *Context) {
	var b []byte
	var err error

	b, err = json.Marshal(r.Data)
	if err != nil {
		(&ErrorResult{Error: err}).Apply(context)
		return
	}
	//return buffer.Bytes(), err
	//b, err = json.Marshal(r.Data)
	//b = buffer.Bytes()

	context.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	context.Response.WriteHeader(http.StatusOK)
	//context.Response.Header().Add("Content-Type", "application/json")
	context.Response.Write(b)
}

type FileServerResult struct {
	Prefix string
	Dir    string
}

// http.StripPrefix
func (fs *FileServerResult) Apply(context *Context) {
	//dir, _ := filepath.Split(context.Request.URL.Path)
	//log.Println(dir, fileName)
	//http.FileServer(http.Dir(conf.Config.ViewDir)+"/"+fs.Dir).ServeHTTP(context.Response, context.Request)
	//http.StripPrefix(conf.Config.ResourcesDir, http.FileServer(http.Dir(fs.Dir))).ServeHTTP(context.Response, context.Request)
	//http.StripPrefix(fs.StripPrefix, http.FileServer(http.Dir(dir))).ServeHTTP(context.Response, context.Request)
	//http.StripPrefix(fs.StripPrefix, http.FileServer(http.Dir(context.Request.URL.Path))).ServeHTTP(context.Response, context.Request)
	//http.StripPrefix(dir, http.FileServer(http.Dir("resources"))).ServeHTTP(context.Response, context.Request)
	//http.StripPrefix(dir, http.FileServer(fs.Dir+"/"+http.Dir(dir))).ServeHTTP(context.Response, context.Request)

	//http.StripPrefix("/resources/", http.FileServer(http.Dir(conf.Config.ResourcesDir+"/resources"))).ServeHTTP(context.Response, context.Request)
	//http.StripPrefix("/web/", http.FileServer(http.Dir(conf.Config.ViewDir+"/web/"))).ServeHTTP(context.Response, context.Request)

	//http.StripPrefix(fs.StripPrefix, http.FileServer(http.Dir(fs.Dir+fs.StripPrefix))).ServeHTTP(context.Response, context.Request)

	http.StripPrefix(fs.Prefix, http.FileServer(http.Dir(fs.Dir))).ServeHTTP(context.Response, context.Request)

}

type HtmlPlainResult struct {
	Data   string
	Params map[string]interface{}
}

func (r *HtmlPlainResult) Apply(context *Context) {

	t := template.New("HtmlPlainResult").Funcs(NewFuncMap(context))
	t, err := t.Parse(r.Data)
	//template.Must(t.Parse(string(b)))
	if err != nil {
		log.Println(err)
		t, err = template.New("HtmlPlainResult").Parse(err.Error())
	}

	data := createPageParams(context, r.Params)
	context.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	context.Response.WriteHeader(http.StatusOK)
	glog.Error(t.Execute(context.Response, data))

	//context.Response.Header().Set("Content-Type", "text/xml; charset=utf-8")
	//context.Response.WriteHeader(http.StatusOK)
	//context.Response.Write([]byte(r.Data))
}

type TextResult struct {
	Data string
}

func (r *TextResult) Apply(context *Context) {

	context.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	context.Response.WriteHeader(http.StatusOK)
	context.Response.Write([]byte(r.Data))
}

type JavaScriptResult struct {
	Data string
}

func (r *JavaScriptResult) Apply(context *Context) {

	context.Response.Header().Add("Content-Type", "application/javascript; charset=utf-8")
	context.Response.WriteHeader(http.StatusOK)
	context.Response.Write([]byte(r.Data))

}

type XMLResult struct {
	Data string
}

func (r *XMLResult) Apply(context *Context) {
	context.Response.Header().Set("Content-Type", "text/xml; charset=utf-8")
	context.Response.WriteHeader(http.StatusOK)
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

	file, err := cache.Read(r.FilePath)
	if err != nil {
		return
	}

	context.Response.Write(file.Byte)

	//context.Response.Header().Set("Location", r.Url)
	//context.Response.WriteHeader(http.StatusFound)
	//context.Response.Header().Set("Content-Type", "")

}

type ImageBytesResult struct {
	Data        []byte
	ContentType string //: image/png
}

func (r *ImageBytesResult) Apply(context *Context) {

	//context.Response.Header().Add()
	context.Response.Header().Set("Content-Type", r.ContentType)
	context.Response.Write(r.Data)

}
