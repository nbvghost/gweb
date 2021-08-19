package gweb

import (
	"errors"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/tool/encryption"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
)

type FileDirType int32

var FileDirT FileDirType = 1
var FileDirTemp FileDirType = 1

func IsFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {

		return false
	} else {
		return true
	}
}
func WriteTempUrlNameFile(b []byte, Url string) string {
	var f *os.File

	fileName := encryption.Md5ByString(Url)
	fullPath := "temp/" + fileName

	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		err = os.MkdirAll("temp", os.ModePerm)
		glog.Error(err)
	}

	f, err := os.Create(fullPath) //创建文件
	glog.Error(err)
	defer f.Close()
	f.Write(b)
	f.Sync()

	return fileName //不包括temp dir
}

func WriteWithFile(file multipart.File, header *multipart.FileHeader, dynamicDirName string, dirType string) (error, string) {

	if strings.EqualFold(dynamicDirName, "") {
		return errors.New("dynamicDirName 不能为空"), ""
	}

	netPath := strings.Trim(dynamicDirName, "/") + "/" + strings.Trim(dirType, "/") + "/"
	filePath := strings.TrimRight(conf.Config.UploadDir, "/") + "/" + netPath

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

	fileName := header.Filename
	if fileInfo, err := os.Stat(filePath + fileName); err == nil {
		if fileInfo.IsDir() {

			return errors.New("目标是一个文件夹"), ""

		}

		f, err := os.OpenFile(filePath+fileName, os.O_RDONLY, os.ModePerm) //打开文件
		if glog.Error(err) {

			return err, ""
		}

		fBytes, err := ioutil.ReadAll(f)
		if glog.Error(err) {

			return err, ""
		}

		if strings.EqualFold(encryption.Md5ByBytes(fileBytes), encryption.Md5ByBytes(fBytes)) {
			return nil, netPath + fileName
		}

		names := strings.Split(header.Filename, ".")
		if len(names) == 0 {
			fileName = header.Filename + "_copy"
		} else if len(names) == 1 {
			fileName = names[0] + "_copy"
		} else {
			fileName = ""
			for i := 0; i < len(names)-1; i++ {

				if i == 0 {
					fileName = fileName + names[i]
				} else {
					fileName = fileName + "_" + names[i]
				}
			}
			fileName = fileName + "." + names[len(names)-1]
		}

	} else {
		fileName = header.Filename
	}

	f, err := os.OpenFile(filePath+fileName, os.O_RDWR|os.O_CREATE, os.ModePerm) //打开文件
	if glog.Error(err) {
		return err, ""
	} else {
		_, err = f.Write(fileBytes)
		glog.Error(err)
		glog.Error(f.Sync())
		glog.Error(f.Close())
	}
	return nil, netPath + fileName

}
func WriteFilePath(fileBytes []byte, dynamicDirName, dirType string, fileName string) (error, string) {

	filePath := strings.Trim(dynamicDirName, "/") + "/" + strings.Trim(dirType, "/") + "/"
	fileFullPath := strings.TrimRight(conf.Config.UploadDir, "/") + "/" + filePath
	if _, err := os.Stat(fileFullPath); os.IsNotExist(err) {
		err = os.MkdirAll(fileFullPath, os.ModePerm)
		glog.Error(err)
	}

	fileUri := fileFullPath + "/" + fileName

	f, err := os.OpenFile(fileUri, os.O_RDWR|os.O_CREATE, os.ModePerm) //打开文件
	if glog.Error(err) {
		return err, ""
	} else {
		f.Write(fileBytes)
		f.Sync()
		f.Close()
	}
	return nil, filePath + "/" + fileName
	//tool.CheckError(err)
	//fmt.Println(f)

}
func WriteFile(fileBytes []byte, ContentType string, dynamicDirName, dirType string) string {
	md5Name := encryption.Md5ByBytes(fileBytes)

	var f *os.File

	fileTypes := strings.Split(ContentType, "/")
	if len(fileTypes) != 2 {
		glog.Error(errors.New("ContentType 格式不正确"))
		return ""
	}

	contentTypeList := strings.Split(fileTypes[1], "-")
	fileType := contentTypeList[len(contentTypeList)-1]
	filePath := strings.Trim(dynamicDirName, "/") + "/" + strings.Trim(dirType, "/") + "/"
	fileFullPath := strings.TrimRight(conf.Config.UploadDir, "/") + "/" + filePath
	if _, err := os.Stat(fileFullPath); os.IsNotExist(err) {
		err = os.MkdirAll(fileFullPath, os.ModePerm)
		glog.Error(err)
	}

	fileName := fileFullPath + "/" + md5Name + "." + fileType

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		//不存在的文件
		f, err = os.Create(fileName) //创建文件
		glog.Error(err)
		defer f.Close()
		f.Write(fileBytes)
		f.Sync()

	}
	return filePath + "/" + md5Name + "." + fileType
}
func WriteTempFile(fileBytes []byte, ContentType string) string {

	md5Name := encryption.Md5ByBytes(fileBytes)
	var f *os.File

	fileTypes := strings.Split(ContentType, "/")
	if len(fileTypes) != 2 {
		glog.Error(errors.New("ContentType 格式不正确"))
		return ""
	}
	contentTypeList := strings.Split(fileTypes[1], "-")
	fileType := contentTypeList[len(contentTypeList)-1]

	fileName := md5Name + "." + fileType
	fullPath := "temp/" + fileName

	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		err = os.MkdirAll("temp", os.ModePerm)
		glog.Error(err)
	}

	f, err := os.Create(fullPath) //创建文件
	glog.Error(err)
	defer f.Close()
	f.Write(fileBytes)
	f.Sync()

	return fileName

}
