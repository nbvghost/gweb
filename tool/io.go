package tool

import (
	"errors"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/gweb/tool/encryption"
	"io/ioutil"
	"mime/multipart"
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
func WriteWithFile(file multipart.File, header *multipart.FileHeader, dynamicDirName string) (error, string) {

	if strings.EqualFold(dynamicDirName, "") {
		return errors.New("dynamicDirName 不能为空"), ""
	}

	now := time.Now()
	filePath := conf.Config.UploadDir + strings.Trim(conf.Config.UploadDirName, "/") + "/" + dynamicDirName + "/" + strconv.Itoa(now.Year()) + "/" + strconv.Itoa(int(now.Month())) + "/"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(filePath, os.ModePerm)
		if glog.Error(err) {

			return err, ""
		}
	}

	fileBytes, err := ioutil.ReadAll(file)
	if glog.Error(err) {

		return err, ""
	}

	fileName := filePath + header.Filename
	if fileInfo, err := os.Stat(fileName); err == nil {
		if fileInfo.IsDir() {

			return errors.New("目标是一个文件夹"), ""

		}

		f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm) //打开文件
		if glog.Error(err) {

			return err, ""
		}

		fBytes, err := ioutil.ReadAll(f)
		if glog.Error(err) {

			return err, ""
		}

		if strings.EqualFold(encryption.Md5ByBytes(fileBytes), encryption.Md5ByBytes(fBytes)) {
			return nil, fileName
		}

		names := strings.Split(header.Filename, ".")
		if len(names) == 0 {
			fileName = filePath + header.Filename + "_copy"
		} else if len(names) == 1 {
			fileName = filePath + names[0] + "_copy"
		} else {
			fileName = filePath
			for i := 0; i < len(names)-1; i++ {

				if i == 0 {
					fileName = fileName + names[i]
				} else {
					fileName = fileName + "_" + names[i]
				}
			}
			fileName = fileName + "." + names[len(names)-1]
		}

	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm) //打开文件
	if glog.Error(err) {
		return err, ""
	} else {
		f.Write(fileBytes)
		f.Sync()
		f.Close()
	}
	return nil, fileName

}
func WriteFilePath(read []byte, subPath string, fileName string) (error, string) {

	filePath := conf.Config.UploadDir + "/" + subPath
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(filePath, os.ModePerm)
		glog.Error(err)
	}

	fileUri := filePath + "/" + fileName

	f, err := os.OpenFile(fileUri, os.O_RDWR|os.O_CREATE, os.ModePerm) //打开文件
	if glog.Error(err) {
		return err, ""
	} else {
		f.Write(read)
		f.Sync()
		f.Close()
	}
	return nil, fileUri
	//tool.CheckError(err)
	//fmt.Println(f)

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
