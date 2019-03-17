package tool

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/nbvghost/glog"
	"html/template"
	"os"
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
	}
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
	b, err := json.Marshal(source)
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
	err := json.Unmarshal([]byte(source), &d)
	glog.Error(err)
	return d
}
func fromJSONToArray(source string) []interface{} {
	d := make([]interface{}, 0)
	err := json.Unmarshal([]byte(source), &d)
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
