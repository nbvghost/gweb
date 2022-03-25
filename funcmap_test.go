package gweb

import (
	"bytes"
	"log"
	"net/http"
	"testing"
	"text/template"
)

type testFunc struct {
}

func (m *testFunc) Call(ctx *Context) IFuncResult {

	return NewMapFuncResult(map[string]interface{}{"tt": 55})

}

func TestRegisterRenderFunction(t *testing.T) {
	type args struct {
		funcName string
		function interface{}
		args     []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		/*{name: "TestRegisterRenderFunction", args: args{
			funcName: "TestFunc",
			function: func(fm *Context, a int, b int, c int, s int) int {

				return 222
			},
			args: []interface{}{454, 22, 4545, 55},
		}},
		{name: "TestRegisterRenderFunction", args: args{
			funcName: "TestFunc1",
			function: func(fm *Context) int {

				return 222
			},
			args: []interface{}{},
		}},*/
		{name: "TestRegisterRenderFunctionMap", args: args{
			funcName: "TestFunc",
			function: func(fm *Context) map[string]interface{} {

				return map[string]interface{}{"dsfds": "sdfsdf"}
			},
			args: []interface{}{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterFunction("test", &testFunc{})

			r, _ := http.NewRequest("GET", "http://a.b.c/test/a/b/c/index", nil)
			funcMap := NewFuncMap(&Context{RoutePath: "fdsfdsf", Request: r})

			/*argsList := make([]reflect.Value, 0)
			for _, arg := range tt.args.args {
				argsList = append(argsList, reflect.ValueOf(arg))
			}

			fd := reflect.ValueOf(funcMap[tt.args.funcName]).Call(argsList)
			log.Println(fd)*/

			templ, err := template.New("dfds").Funcs(template.FuncMap(funcMap)).Parse("@{{$k:=TestFunc}}/{{$k.tt}}@")
			log.Println(err)
			buffer := bytes.NewBuffer([]byte{})
			log.Println(templ.Execute(buffer, nil))
			log.Println(string(buffer.Bytes()))

			//fds := (reflect.ValueOf(funcMap["HTML"])).Call([]reflect.Value{reflect.ValueOf("dsfdsfsdf")}) //reflect.ValueOf(funcMap["HTML"])
			//log.Println(fds)

		})
	}
}
