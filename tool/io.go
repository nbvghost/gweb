package tool

import (
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool/encryption"

	"os"
	"strconv"
	"strings"
	"time"
)

type FileDirType int32

var FileDirT FileDirType = 1
var FileDirTemp FileDirType = 1

func CreateFile(FilePath, FileName string) *os.File {

	return nil
}
func IsFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {

		return false
	} else {
		return true
	}
}
func WriteTempUrlNameFile(b []byte, Url string) string {

	md5Name := encryption.Md5ByString(Url)
	var f *os.File

	filePath := string(md5Name[0:1])

	path := filePath + "/"
	fileName := md5Name
	fullPath := "temp/" + path

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err = os.MkdirAll(fullPath, os.ModePerm)
		glog.Error(err)
	}

	if _, err := os.Stat(fullPath + fileName); os.IsNotExist(err) {

		f, err = os.Create(fullPath + fileName) //创建文件
		glog.Error(err)
		defer f.Close()
		f.Write(b)
		f.Sync()
	} else {
		return path + fileName
	}
	return path + fileName
}
func WriteTempFile(b []byte, ContentType string) string {

	md5Name := encryption.Md5ByBytes(b)
	var f *os.File

	fileTypes := strings.Split(ContentType, "/")
	if len(fileTypes) == 0 || len(fileTypes) == 1 {

		return ""
	}
	fileType := fileTypes[1]
	fileType = strings.Split(fileType, "+")[0]
	filePath := string(md5Name[0:1])

	path := filePath + "/"
	fileName := md5Name + "." + fileType
	fullPath := "temp/" + path

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err = os.MkdirAll(fullPath, os.ModePerm)
		glog.Error(err)
	}

	if _, err := os.Stat(fullPath + fileName); os.IsNotExist(err) {

		f, err = os.Create(fullPath + fileName) //创建文件
		glog.Error(err)
		defer f.Close()
		f.Write(b)
		f.Sync()
	} else {
		//f, err = os.OpenFile(fileName, os.O_RDONLY, os.ModePerm) //打开文件
		//tool.CheckError(err)
		//fmt.Println(f)
	}
	return path + fileName

}
func WriteFile(b []byte, ContentType string) string {
	md5Name := encryption.Md5ByBytes(b)
	now := time.Now()
	var f *os.File

	fileType := strings.Split(ContentType, "/")[1]
	fileType = strings.Split(fileType, "+")[0]
	filePath := conf.Config.UploadDir + "/" + strconv.Itoa(now.Year()) + "/" + strconv.Itoa(int(now.Month())) + "/" + strconv.Itoa(now.Day()) + "/" + md5Name[0:2]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(filePath, os.ModePerm)
		glog.Error(err)
	}

	fileName := filePath + "/" + md5Name + "." + fileType

	if _, err := os.Stat(fileName); os.IsNotExist(err) {

		f, err = os.Create(fileName) //创建文件
		glog.Error(err)
		defer f.Close()
		f.Write(b)
		f.Sync()

	} else {
		//f, err = os.OpenFile(fileName, os.O_RDONLY, os.ModePerm) //打开文件
		//tool.CheckError(err)
		//fmt.Println(f)
	}
	return fileName
}
