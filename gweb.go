package gweb

import (
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func StartServer(HTTP, HTTPS bool) {
	if !HTTP && !HTTPS {
		tool.Trace("选择http或https")
		return
	}

	if HTTP {
		s := &http.Server{
			Addr: conf.Config.HttpPort,
			//Handler:        http.DefaultServeMux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			//MaxHeaderBytes: 1 << 20,
		}
		tool.Trace("http server：" + conf.Config.HttpPort)
		err := s.ListenAndServe()
		tool.CheckError(err)
	}

	if HTTPS {
		s := &http.Server{
			Addr: conf.Config.HttpsPort,
			//Handler:        http.DefaultServeMux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			//MaxHeaderBytes: 1 << 20,
		}
		tool.Trace("https server：" + conf.Config.HttpsPort)
		err := s.ListenAndServeTLS(conf.Config.TLSCertFile, conf.Config.TLSKeyFile)
		tool.CheckError(err)
	}

}
