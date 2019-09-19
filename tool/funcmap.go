package tool

import (
	"bufio"
	"bytes"
	"github.com/nbvghost/glog"
	"html/template"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/nbvghost/gweb/conf"
	"time"
)

func FuncMap() template.FuncMap {

	return template.FuncMap{
		"IncludeHTML":     includeHTML,
		"Split":           splitFunc,
		"FromJSONToMap":   fromJSONToMap,
		"FromJSONToArray": fromJSONToArray,
		"CipherDecrypter": cipherDecrypter,
		"CipherEncrypter": cipherEncrypter,
		"Int2String":      int2String,
		"Uint2String":     uint2String,
		"Float2String":    float2String,
		"ToJSON":          toJSON,
		"DateTimeFormat":  DateTimeFormat,
		"HTML":            HTML,
		"UrlQueryEncode":  urlQueryEncode,
		"DigitAdd":        digitAdd,
		"DigitSub":        digitSub,
		"DigitMul":        digitMul,
		"DigitDiv":        digitDiv,
		"DigitMod":        digitMod,
	}
}
func digitAdd(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a+_b, 'f', prec, 64), 64)
	return f
}
func digitSub(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a-_b, 'f', prec, 64), 64)
	return f
}
func digitMul(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a*_b, 'f', prec, 64), 64)
	return f
}
func digitDiv(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	//f, _ := strconv.ParseFloat(strconv.FormatFloat(_a/_b, 'f', prec, 64), 64)
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a/_b, 'f', prec, 64), 64)
	return f
}
func digitMod(a, b interface{}) uint64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()

	///f, _ := strconv.ParseFloat(strconv.FormatFloat(_a%_b, 'f', prec, 64), 64)
	return uint64(_a) % uint64(_b)

}
func urlQueryEncode(source map[string]string) template.URL {
	//fmt.Println(source)
	v := &url.Values{}
	for key := range source {
		v.Set(key, source[key])
	}
	return template.URL(v.Encode())
}
func HTML(source string) template.HTML {
	//fmt.Println(source)
	return template.HTML(source)
}
func DateTimeFormat(source time.Time, format string) string {
	//fmt.Println(source)
	//fmt.Println(format)
	return source.Format(format)
}
func toJSON(source interface{}) string {
	b, err := JsonMarshal(source)
	glog.Error(err)
	return string(b)
}
func int2String(source interface{}) string {

	return strconv.FormatInt((source.(int64)), 10)
}
func uint2String(source interface{}) string {

	return strconv.FormatUint((source.(uint64)), 10)
}
func float2String(source interface{}) string {
	return strconv.FormatFloat((source.(float64)), 'f', -1, 64)
}
func cipherDecrypter(source string) string {

	str := CipherDecrypter(public_PassWord, source)
	return str
}
func cipherEncrypter(source string) string {
	str := CipherEncrypter(public_PassWord, source)
	return str
}
func fromJSONToMap(source string) map[string]interface{} {
	d := make(map[string]interface{})
	err := JsonUnmarshal([]byte(source), &d)
	glog.Error(err)
	return d
}
func fromJSONToArray(source string) []interface{} {
	d := make([]interface{}, 0)
	err := JsonUnmarshal([]byte(source), &d)
	glog.Error(err)
	return d
}
func splitFunc(source string, sep string) []string {

	return strings.Split(source, sep)
}
func includeHTML(url string, params interface{}) template.HTML {
	//util.Trace(params)
	//paramsMap := make(map[string]interface{})

	b := bytes.NewBuffer(make([]byte, 0))
	ww := bufio.NewWriter(b)

	t, err := template.ParseFiles(conf.Config.ViewDir + "/" + url)
	if os.IsNotExist(err) {
		ww.WriteString("IncludeHTML:not found path in:" + url)
		t = template.New("static")
	} else {
		t.Execute(ww, params)
	}

	//checkError(err, "read from file template")

	ww.Flush()
	//template.JSEscape()
	//template.HTMLEscapeString()

	//	util.Trace(string(b.Bytes()))
	///return string(b.Bytes());
	return template.HTML(string(b.Bytes()))
}
