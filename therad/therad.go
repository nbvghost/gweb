package therad

import "runtime/debug"

func NewCoroutine(run func(), fail func(v interface{}, stack []byte)) {

	go func() {
		defer func() {
			if r := recover(); r != nil {
				b := debug.Stack()
				fail(r, b)
			}
		}()
		run()
	}()
}


func NewCoroutineParams(run func(args []interface{}), fail func(v interface{}, stack []byte),params ...interface{}) {

	go func() {
		defer func() {
			if r := recover(); r != nil {
				b := debug.Stack()
				fail(r, b)
			}
		}()
		run(params)
	}()
}