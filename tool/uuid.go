package tool

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/nbvghost/glog"
	"io"
	"math"
	rand2 "math/rand"
	"net"
	"strconv"
	"time"
)
var hardwareAddrs=GetHardwareAddrs()
func GetHardwareAddrs() string {
	hardwareAddrs:=""
	interfaces, sdfsdfsd := net.Interfaces()
	glog.Error(sdfsdfsd)
	for _, inter := range interfaces {
		//fmt.Println(inter.Name)
		mac := inter.HardwareAddr //获取本机MAC地址
		//fmt.Println("MAC = ", mac)
		hardwareAddrs = hardwareAddrs + mac.String()
	}

	return hardwareAddrs
}
func UUID() string {
	rander := rand.Reader
	uuid := make([]byte, 64)
	if _, err := io.ReadFull(rander, uuid); err != nil {
		panic(err.Error()) // rand should never fail
	}
	t := base64.URLEncoding.EncodeToString(uuid)
	_ranSource := rand2.New(rand2.NewSource(time.Now().UnixNano()))
	t = hardwareAddrs + t + strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.FormatInt(_ranSource.Int63n(math.MaxInt64), 10)
	//fmt.Println(t)
	h := md5.New()
	h.Write([]byte(t))
	return hex.EncodeToString(h.Sum(nil))
	//return uuid
}
