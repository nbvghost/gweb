package gweb

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool/number"
	"html/template"
	"log"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*var FunctionMap = template.FuncMap{
	"IncludeHTML":     includeHTML,
	"Split":           splitFunc,
	"FromJSONToMap":   fromJSONToMap,
	"FromJSONToArray": fromJSONToArray,
	"ParseFloat":      parseFloat,
	"ParseInt":        parseInt,
	"ToJSON":          toJSON,
	"DateTimeFormat":  dateTimeFormat,
	"HTML":            html,
	"UrlQueryEncode":  urlQueryEncode,
	"DigitAdd":        digitAdd,
	"DigitSub":        digitSub,
	"DigitMul":        digitMul,
	"DigitDiv":        digitDiv,
	"MakeArray":       makeArray,
	"DigitMod":        digitMod,
	//"CipherDecrypter": cipherDecrypter,
	//"CipherEncrypter": cipherEncrypter,
}*/

/*func FuncMap() template.FuncMap {

	return FunctionMap
}*/

var regFuncMap = make(map[string]interface{})

func RegisterRenderFunction(funcName string, function interface{}) error {

	if _, ok := regFuncMap[funcName]; ok {
		return errors.New(fmt.Sprintf("%v函数已经存在", funcName))
	}

	v := reflect.ValueOf(function)
	if v.Kind() != reflect.Func || v.Kind() == reflect.Ptr {
		return errors.New("function 必需是函数")
	}

	functionType := v.Type()

	functionNumIn := functionType.NumIn()
	if functionNumIn < 1 {
		return errors.New("function 参数个数必须1个或以上")
	}

	if strings.Contains(functionType.In(0).String(), "gweb.Context") == false {
		return errors.New("function 第一个参数必须是gweb.Context")
	}

	functionNumOut := functionType.NumOut()
	if functionNumOut != 1 {
		return errors.New("function 要有一个返回值")
	}

	regFuncMap[funcName] = function

	return nil
}

type FuncObject struct {
	funcMap template.FuncMap
	c       *Context
}

func NewFuncMap(context *Context) template.FuncMap {
	fm := &FuncObject{}
	fm.c = context
	fm.funcMap = make(template.FuncMap)
	fm.funcMap["IncludeHTML"] = fm.includeHTML
	fm.funcMap["Split"] = fm.splitFunc
	fm.funcMap["FromJSONToMap"] = fm.fromJSONToMap
	fm.funcMap["FromJSONToArray"] = fm.fromJSONToArray
	fm.funcMap["ParseFloat"] = fm.parseFloat
	fm.funcMap["ParseInt"] = fm.parseInt
	fm.funcMap["ToJSON"] = fm.toJSON
	fm.funcMap["DateTimeFormat"] = fm.dateTimeFormat
	fm.funcMap["HTML"] = fm.html
	fm.funcMap["UrlQueryEncode"] = fm.urlQueryEncode
	fm.funcMap["DigitAdd"] = fm.digitAdd
	fm.funcMap["DigitSub"] = fm.digitSub
	fm.funcMap["DigitMul"] = fm.digitMul
	fm.funcMap["DigitDiv"] = fm.digitDiv
	fm.funcMap["MakeArray"] = fm.makeArray
	fm.funcMap["DigitMod"] = fm.digitMod
	fm.funcMap["Test"] = fm.test

	for funcName := range regFuncMap {

		func(funcName string) {
			function := regFuncMap[funcName]

			v := reflect.ValueOf(function)
			functionType := v.Type()

			functionNumIn := functionType.NumIn()
			functionNumOut := functionType.NumOut()
			//---
			argsIn := make([]reflect.Type, 0)
			for i := 1; i < functionNumIn; i++ {
				argsIn = append(argsIn, functionType.In(i))
			}
			argsOut := make([]reflect.Type, 0)
			for i := 0; i < functionNumOut; i++ {
				argsOut = append(argsOut, functionType.Out(i))
			}

			//reflect.FuncOf(args)

			makeFuncType := reflect.FuncOf(argsIn, argsOut, false)

			backCallFunc := reflect.MakeFunc(makeFuncType, func(args []reflect.Value) (results []reflect.Value) {
				backCallFuncArgs := make([]reflect.Value, 0)
				backCallFuncArgs = append(backCallFuncArgs, reflect.ValueOf(fm.c))
				backCallFuncArgs = append(backCallFuncArgs, args...)

				resultArgs := v.Call(backCallFuncArgs)

				v := resultArgs[0]
				log.Println(v)
				log.Println("out", reflect.ValueOf(map[string]interface{}{"dfds": 154}).Interface())
				log.Println("out", reflect.Indirect(reflect.ValueOf(map[string]interface{}{"dfds": 154})).Interface())
				return []reflect.Value{v}
			})
			fm.funcMap[funcName] = backCallFunc.Interface()
		}(funcName)

	}

	return fm.funcMap
}

func (fo *FuncObject) test() map[string]interface{} {

	return map[string]interface{}{"fdsfds": 4545}
}
func (fo *FuncObject) digitAdd(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a+_b, 'f', prec, 64), 64)
	return f
}
func (fo *FuncObject) digitSub(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a-_b, 'f', prec, 64), 64)
	return f
}
func (fo *FuncObject) digitMul(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a*_b, 'f', prec, 64), 64)
	return f
}
func (fo *FuncObject) digitDiv(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	//f, _ := strconv.ParseFloat(strconv.FormatFloat(_a/_b, 'f', prec, 64), 64)
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a/_b, 'f', prec, 64), 64)
	return f
}
func (fo *FuncObject) digitMod(a, b interface{}) uint64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()

	///f, _ := strconv.ParseFloat(strconv.FormatFloat(_a%_b, 'f', prec, 64), 64)
	return uint64(_a) % uint64(_b)

}
func (fo *FuncObject) makeArray(len int) []int {

	return make([]int, len)
}
func (fo *FuncObject) urlQueryEncode(source map[string]string) template.URL {
	//fmt.Println(source)
	v := &url.Values{}
	for key := range source {
		v.Set(key, source[key])
	}
	return template.URL(v.Encode())
}
func (fo *FuncObject) html(source string) template.HTML {
	//fmt.Println(source)
	return template.HTML(source)
}
func (fo *FuncObject) dateTimeFormat(source time.Time, format string) string {
	//fmt.Println(source)
	//fmt.Println(format)
	return source.Format(format)
}
func (fo *FuncObject) toJSON(source interface{}) string {
	b, err := json.Marshal(source)
	glog.Error(err)
	return string(b)
}
func (fo *FuncObject) parseInt(source interface{}) int {

	return number.ParseInt(source)
}

func (fo *FuncObject) parseFloat(source interface{}) float64 {
	return number.ParseFloat(source)
}

/*func cipherDecrypter(source string) string {

	str := encryption.CipherDecrypter(encryption.public_PassWord, source)
	return str
}
func cipherEncrypter(source string) string {
	str := encryption.CipherEncrypter(encryption.public_PassWord, source)
	return str
}*/
func (fo *FuncObject) fromJSONToMap(source string) map[string]interface{} {
	d := make(map[string]interface{})
	err := json.Unmarshal([]byte(source), &d)
	glog.Error(err)
	return d
}
func (fo *FuncObject) fromJSONToArray(source string) []interface{} {
	d := make([]interface{}, 0)
	err := json.Unmarshal([]byte(source), &d)
	glog.Error(err)
	return d
}
func (fo *FuncObject) splitFunc(source string, sep string) []string {

	return strings.Split(source, sep)
}
func (fo *FuncObject) includeHTML(url string, params interface{}) template.HTML {
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
