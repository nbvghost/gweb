package thread

import (
    "os"
    "os/signal"
    "runtime/debug"
    "syscall"
    "time"
)

type Option struct {
    ReRun    bool
    Duration time.Duration
    Error    interface{}
    Stack    []byte
    Index    int
}

func NewCoroutine(run func(option *Option), fail func(option *Option), option *Option) {
    if option.Duration <= 0 {
        option.Duration = time.Second
    }
    go func() {
        defer func() {
            if r := recover(); r != nil {
                b := debug.Stack()
                option.Error = r
                option.Stack = b
                fail(option)
            }

            if option.ReRun {

                time.Sleep(option.Duration)
                NewCoroutine(run, fail, option)
            }

        }()
        run(option)
    }()
}

func NewCoroutineParams(run func(args []interface{}), fail func(v interface{}, stack []byte), params ...interface{}) {

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

func ListeningSignal(signals ...os.Signal) chan os.Signal {
    //syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM
    sn := make(chan os.Signal)
    //_sn := make(chan os.Signal)
    if len(signals) > 0 {
        signal.Notify(sn, signals...)

    } else {
        signal.Notify(sn, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
    }
    /*go func(sn,_sn chan os.Signal) {

      }(sn,_sn)*/

    //defer signal.Stop(sn)
    return sn
}
