package gweb

import (
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/conf"
	"net/http"
	_ "net/http/pprof"
	"time"
)


func StartServer(HTTP, HTTPS bool) {

	static:=Static{}
	http.HandleFunc("/file/up", static.fileUp)
	http.HandleFunc("/file/load", static.fileLoad)
	http.HandleFunc("/file/net/load", static.fileNetLoad)
	http.HandleFunc("/file/temp/load", static.fileTempLoad)

	http.Handle("/"+conf.Config.ResourcesDirName+"/", http.StripPrefix("/"+conf.Config.ResourcesDirName+"/", http.FileServer(http.Dir(conf.Config.ResourcesDir))))
	http.Handle("/"+conf.Config.UploadDirName+"/", http.StripPrefix("/"+conf.Config.UploadDirName+"/", http.FileServer(http.Dir(conf.Config.UploadDir))))
	http.Handle("/temp/", http.StripPrefix("/temp/", http.FileServer(http.Dir("temp"))))


	if !HTTP && !HTTPS {
		glog.Trace("选择http或https")
		return
	}

	if HTTP {


		if HTTPS==false{
			s := &http.Server{
				Addr: conf.Config.HttpPort,
				//Handler:        http.DefaultServeMux,
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
				//MaxHeaderBytes: 1 << 20,
			}
			glog.Trace("http server：" + conf.Config.HttpPort)
			err := s.ListenAndServe()
			glog.Error(err)
		}else{
			go func() {
				s := &http.Server{
					Addr: conf.Config.HttpPort,
					//Handler:        http.DefaultServeMux,
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
					//MaxHeaderBytes: 1 << 20,
				}
				glog.Trace("http server：" + conf.Config.HttpPort)
				err := s.ListenAndServe()
				glog.Error(err)
			}()
		}
	}

	if HTTPS {
		s := &http.Server{
			Addr: conf.Config.HttpsPort,
			//Handler:        http.DefaultServeMux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			//MaxHeaderBytes: 1 << 20,
		}
		glog.Trace("https server：" + conf.Config.HttpsPort)
		err := s.ListenAndServeTLS(conf.Config.TLSCertFile, conf.Config.TLSKeyFile)
		glog.Error(err)
	}

}
func main()  {

}