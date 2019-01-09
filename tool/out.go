package tool

import (
	"github.com/nbvghost/gweb/conf"
	"log"
	"runtime"
)

func Trace(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if conf.Config.Debug {
		//util.Trace(funcName,file,line,ok)
		for _, va := range v {
			if va != nil {
				log.Println(file, line, va)
			}
		}
	}

}

func CheckError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Println(file, line, err)
	}
}
