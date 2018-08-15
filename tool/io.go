package tool

import (
	"time"
	"os"
	"strings"
	"strconv"
	"github.com/nbvghost/gweb/conf"
)

func WriteTempFile(b []byte, ContentType string) string {

	md5Name := Md5ByBytes(b)
	var f *os.File

	fileTypes := strings.Split(ContentType, "/")
	if len(fileTypes)==0||len(fileTypes)==1{

		return ""
	}
	fileType := fileTypes[1]
	fileType = strings.Split(fileType, "+")[0]
	filePath := string(md5Name[0:1])

	path := filePath + "/"
	fileName := md5Name + "." + fileType
	fullPath:="temp/"+path

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err = os.MkdirAll(fullPath, os.ModePerm)
		CheckError(err)
	}

	if _, err := os.Stat(fullPath+fileName); os.IsNotExist(err) {

		f, err = os.Create(fullPath+fileName) //创建文件
		CheckError(err)
		defer f.Close()
		f.Write(b)
		f.Sync()
	} else {
		//f, err = os.OpenFile(fileName, os.O_RDONLY, os.ModePerm) //打开文件
		//tool.CheckError(err)
		//fmt.Println(f)
	}
	return path+fileName

}
func WriteFile(b []byte, ContentType string) string {
	md5Name := Md5ByBytes(b)
	now := time.Now()
	var f *os.File

	fileType := strings.Split(ContentType, "/")[1]
	fileType = strings.Split(fileType, "+")[0]
	filePath := conf.Config.UploadDir+"/" + strconv.Itoa(now.Year()) + "/" + strconv.Itoa(int(now.Month())) + "/" + strconv.Itoa(now.Day()) + "/" + md5Name[0:2]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(filePath, os.ModePerm)
		CheckError(err)
	}

	fileName := filePath + "/" + md5Name + "." + fileType

	if _, err := os.Stat(fileName); os.IsNotExist(err) {

		f, err = os.Create(fileName) //创建文件
		CheckError(err)
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