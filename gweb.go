package gweb

import (
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/conf"

	"net/http"
	_ "net/http/pprof"
)


func StartServer(serverMux *http.ServeMux,httpServer *http.Server, httpsServer *http.Server) {

	if serverMux==nil{
		return
	}
	if httpServer==nil && httpsServer==nil{
		return
	}

	static:=Static{}
	serverMux.HandleFunc("/file/up", static.fileUp)
	serverMux.HandleFunc("/file/load", static.fileLoad)
	serverMux.HandleFunc("/file/net/load", static.fileNetLoad)
	serverMux.HandleFunc("/file/temp/load", static.fileTempLoad)

	serverMux.Handle("/"+conf.Config.ResourcesDirName+"/", http.StripPrefix("/"+conf.Config.ResourcesDirName+"/", http.FileServer(http.Dir(conf.Config.ResourcesDir))))
	serverMux.Handle("/"+conf.Config.UploadDirName+"/", http.StripPrefix("/"+conf.Config.UploadDirName+"/", http.FileServer(http.Dir(conf.Config.UploadDir))))
	serverMux.Handle("/temp/", http.StripPrefix("/temp/", http.FileServer(http.Dir("temp"))))


	if httpServer==nil && httpsServer==nil {
		panic("选择http或https")
		return
	}

	if httpServer!=nil {

		if httpsServer==nil{

			glog.Trace("gweb start http at：" + httpServer.Addr)
			err := httpServer.ListenAndServe()
			panic(err)
		}else{
			go func() {

				glog.Trace("gweb start http at：" + httpServer.Addr)
				err := httpServer.ListenAndServe()
				panic(err)

			}()
		}
	}

	if httpsServer!=nil {

		glog.Trace("gweb start https at：" + httpsServer.Addr)
		err := httpsServer.ListenAndServeTLS(conf.Config.TLSCertFile, conf.Config.TLSKeyFile)
		panic(err)
	}

}
func main()  {

}