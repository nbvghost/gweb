package gweb

import (
	"encoding/json"
	"github.com/nbvghost/glog"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"time"

	"bytes"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool"
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
type SingleHostReverseProxyResult struct {
	Target *url.URL
}

func (r *SingleHostReverseProxyResult) Apply(context *Context) {

	rp:=httputil.NewSingleHostReverseProxy(r.Target)
	rp.ServeHTTP(context.Response,context.Request)

}
/**
类型/子类型	扩展名
application/envoy	evy
application/fractals	fif
application/futuresplash	spl
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

	if strings.EqualFold(path,"/"){
		if strings.EqualFold(conf.Config.DefaultPage,"")==false{
			path=path+conf.Config.DefaultPage
			var redirectToUrlResult=&RedirectToUrlResult{Url:path}
			redirectToUrlResult.Apply(context)
			return
		}

	}

	path = strings.TrimRight(path,"/")
	b, err := ioutil.ReadFile(conf.Config.ViewDir + path+conf.Config.ViewSuffix)
	if err != nil {
		//fmt.Println(context.Request.Header)





		var haveMIME = false
		b, err := ioutil.ReadFile(conf.Config.ViewDir + path)
		if err==nil{
			re, err := regexp.Compile("\\/([0-9a-zA-Z_]+)\\.([0-9a-zA-Z]+)$")
			glog.Error(err)

			if re.MatchString(path){
				Groups:=re.FindAllStringSubmatch(path, -1)
				//[[/fgsd_gffdgdf.txt fgsd_gffdgdf txt]]
				//{"ContentType": "text/html","Extension":"html"}
				Extension:=Groups[0][2]
				for index:= range conf.Config.ViewActionMapping{
					ce:=conf.Config.ViewActionMapping[index]
					if strings.EqualFold(ce.Extension,Extension){
						context.Response.Header().Set("Content-Type", ce.ContentType+"; charset=utf-8")
						//w.Header().Set("X-Content-Type-Options", "nosniff")
						context.Response.WriteHeader(http.StatusOK)
						context.Response.Write(b)
						haveMIME =true
						break
					}

				}
			}

		}
		if haveMIME==false{

			fi,err:=os.Stat(conf.Config.ViewDir + path)
			//log.Println(err)
			if err==nil && fi.IsDir(){

				path=path+"/"+conf.Config.DefaultPage
				var redirectToUrlResult=&RedirectToUrlResult{Url:path}
				redirectToUrlResult.Apply(context)

			}else{
				//没有找到路由，
				http.NotFound(context.Response, context.Request)
			}


		}

	} else {
		context.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
		context.Response.WriteHeader(http.StatusOK)
		t, err := template.New("default").Funcs(tool.FuncMap()).Parse(string(b))
		glog.Error(err)
		t.Execute(context.Response, nil)
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

//只映射已经定义的后缀模板文件
type HTMLResult struct {
	Name   string
	Params map[string]interface{}
}

func (r *HTMLResult) Apply(context *Context) {

	path := context.Request.URL.Path

	var b []byte
	var err error

	if strings.EqualFold(r.Name, ""){
		//html 只处理，已经定义后缀名的文件
		b, err = ioutil.ReadFile(fixPath(conf.Config.ViewDir + "/" + path+conf.Config.ViewSuffix))
	}else{
		b, err = ioutil.ReadFile(fixPath(conf.Config.ViewDir + "/" + r.Name + conf.Config.ViewSuffix))
	}
	if err != nil {
		//判断是否有默认页面
		//fmt.Println(fixPath(Config.ViewDir + "/" + path +"/"+ Config.DefaultPage))
		b, err = ioutil.ReadFile(fixPath(conf.Config.ViewDir + "/" + path + "/" + conf.Config.DefaultPage + conf.Config.ViewSuffix))
		if err != nil {
			(&ViewActionMappingResult{}).Apply(context)
			return
		}
	}

	/*if err != nil {

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
				(&ViewActionMappingResult{}).Apply(context)
				return
			}
		}

	}*/

	//t, err := template.New("default").Funcs(FuncMap).Parse(string(b))
	t := template.New("default").Funcs(tool.FuncMap())
	t, err = t.Parse(string(b))
	//template.Must(t.Parse(string(b)))
	if err != nil {
		log.Println(err)
		t, err = template.New("").Parse(err.Error())
	}

	data := make(map[string]interface{})
	data["session"] = context.Session.Attributes.GetMap()
	data["query"] = tool.QueryParams(context.Request.URL.Query())
	data["params"] = r.Params
	data["host"] = context.Request.Host
	data["time"] = time.Now().Unix() * 1000

	jsonData:=make(map[string]interface{})
	json.Unmarshal([]byte(conf.JsonText),&jsonData)

	data["data"] =jsonData
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

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	err = encoder.Encode(r.Data)
	//return buffer.Bytes(), err
	//b, err = json.Marshal(r.Data)
	b = buffer.Bytes()

	if err != nil {
		(&ErrorResult{Error: err}).Apply(context)
		return
	}

	context.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	context.Response.WriteHeader(http.StatusOK)
	//context.Response.Header().Add("Content-Type", "application/json")
	context.Response.Write(b)
}
type HtmlPlainResult struct {
	Data string
}

func (r *HtmlPlainResult) Apply(context *Context) {

	t := template.New("default").Funcs(tool.FuncMap())
	t, err := t.Parse(r.Data)
	//template.Must(t.Parse(string(b)))
	if err != nil {
		log.Println(err)
		t, err = template.New("").Parse(err.Error())
	}

	data := make(map[string]interface{})
	data["session"] = context.Session.Attributes.GetMap()
	data["query"] = tool.QueryParams(context.Request.URL.Query())
	data["host"] = context.Request.Host
	data["time"] = time.Now().Unix() * 1000
	context.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	context.Response.WriteHeader(http.StatusOK)
	t.Execute(context.Response, data)



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
	Data        []byte
	ContentType string //: image/png
}

func (r *ImageBytesResult) Apply(context *Context) {

	//context.Response.Header().Add()
	context.Response.Header().Set("Content-Type", r.ContentType)
	context.Response.Write(r.Data)

}
