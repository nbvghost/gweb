package tool

import (
	"log"
	"runtime"
)

func Trace(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	//util.Trace(funcName,file,line,ok)
	log.Println(file, line, v)
}

func CheckError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Println(file, line, err)
	}
}
