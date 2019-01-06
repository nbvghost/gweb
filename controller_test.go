package gweb

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestGETMethod(t *testing.T) {
	type args struct {
		RoutePath string
		call      func(context *Context) Result
	}
	tests := []struct {
		name string
		args args
		want function
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GETMethod(tt.args.RoutePath, tt.args.call); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GETMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOPTMethod(t *testing.T) {
	type args struct {
		RoutePath string
		call      func(context *Context) Result
	}
	tests := []struct {
		name string
		args args
		want function
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OPTMethod(tt.args.RoutePath, tt.args.call); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OPTMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHEAMethod(t *testing.T) {
	type args struct {
		RoutePath string
		call      func(context *Context) Result
	}
	tests := []struct {
		name string
		args args
		want function
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HEAMethod(tt.args.RoutePath, tt.args.call); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HEAMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPOSMethod(t *testing.T) {
	type args struct {
		RoutePath string
		call      func(context *Context) Result
	}
	tests := []struct {
		name string
		args args
		want function
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := POSMethod(tt.args.RoutePath, tt.args.call); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("POSMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPUTMethod(t *testing.T) {
	type args struct {
		RoutePath string
		call      func(context *Context) Result
	}
	tests := []struct {
		name string
		args args
		want function
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PUTMethod(tt.args.RoutePath, tt.args.call); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PUTMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDELMethod(t *testing.T) {
	type args struct {
		RoutePath string
		call      func(context *Context) Result
	}
	tests := []struct {
		name string
		args args
		want function
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DELMethod(tt.args.RoutePath, tt.args.call); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DELMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTRAMethod(t *testing.T) {
	type args struct {
		RoutePath string
		call      func(context *Context) Result
	}
	tests := []struct {
		name string
		args args
		want function
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TRAMethod(tt.args.RoutePath, tt.args.call); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TRAMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCONMethod(t *testing.T) {
	type args struct {
		RoutePath string
		call      func(context *Context) Result
	}
	tests := []struct {
		name string
		args args
		want function
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CONMethod(tt.args.RoutePath, tt.args.call); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CONMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestALLMethod(t *testing.T) {
	type args struct {
		RoutePath string
		call      func(context *Context) Result
	}
	tests := []struct {
		name string
		args args
		want function
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ALLMethod(tt.args.RoutePath, tt.args.call); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ALLMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseController_NewController(t *testing.T) {
	type args struct {
		path string
		ic   IController
	}
	tests := []struct {
		name string
		c    *BaseController
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.NewController(tt.args.path, tt.args.ic)
		})
	}
}

func TestBaseController_AddSubController(t *testing.T) {
	type args struct {
		path  string
		isubc IController
	}
	tests := []struct {
		name string
		c    *BaseController
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.AddSubController(tt.args.path, tt.args.isubc)
		})
	}
}

func TestBaseController_AddHandler(t *testing.T) {
	type args struct {
		_function function
	}
	tests := []struct {
		name string
		c    *BaseController
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.AddHandler(tt.args._function)
		})
	}
}



func TestBaseController_ServeHTTP(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		c    *BaseController
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}

func Test_delRepeatAll(t *testing.T) {
	type args struct {
		src string
		new string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := delRepeatAll(tt.args.src, tt.args.new); got != tt.want {
				t.Errorf("delRepeatAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_getPathParams(b *testing.B) {


	for i:=0;i<b.N;i++{
		//got, got1 := getPathParams("/platfdfsdsfdsorm/{dsfdsfds}顺{ChannelID}吴{Dfdfd}清{sdfdsf}/1303/game","/platfdfsdsfdsorm/sdfdsfdsf顺45435435dsf吴dsf43543543dfsgdfs清sdfsdfds/1303/game")
		//got, got1 := getPathParams("{dsfdsfds}顺{ChannelID}吴{Dfdfd}清{sdfdsf}","sdfdsfdsf顺45435435dsf吴dsf43543543dfsgdfs清sdfsdfds")
		//got, got1 := getPathParams("fsdg{dsfdsfds}sdafsda","fsdgfdafsdafsdafsda")
		fmt.Println(getPathParams("/{sfsdfds}/{dfds}","/雷克萨反对sdfdssdfdsffdsfsdfsf/sad呆困运输成本吵过架基材"))
		//   20000	     67015 ns/op
		//   10000	    106824 ns/op
		//fmt.Println(got,got1)
		//fmt.Println("--------------------------")
	}
}
func Test_validateRoutePath(t *testing.T) {

	fmt.Println(validateRoutePath("/{ds}{ds}/"))

}
func Test_getPathParams(t *testing.T) {
	type args struct {
		RoutePath string
		Path      string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 map[string]string
	}{
		//{name:"Test_getPathParams",args:args{RoutePath:"/sdfsd/{dfdsfs}_dsfdsf/{DFdfd}/dfdf_{sdfdsfsdf}/{dfdsfddd}-dfds{fd}-{jk}/f{dfd}",Path:"/sdfsd/dfd5f4ds_dsfdsf/sdf/dfdf_sd/dfsdsfds-dfdsfdf-dfdf/fsdfsd"},want:true},
		{name:"Test_getPathParams",args:args{RoutePath:"dsfds/dd{dsfsd}",Path:"dsfds/ddsdafsdaf"},want:true},
	}	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getPathParams(tt.args.RoutePath, tt.args.Path)
			if got != tt.want {
				t.Errorf("getPathParams() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getPathParams() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_fixPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fixPath(tt.args.path); got != tt.want {
				t.Errorf("fixPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
