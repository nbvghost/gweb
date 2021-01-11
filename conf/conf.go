package conf

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sync"
)

type ViewActionMapping struct {
	ContentType string `json:"ContentType"`
	Extension   string `json:"Extension"`
}

type JsonDataEntity struct {
	m        map[string]interface{}
	gobBuf   []byte
	jsonData []byte
	sync.RWMutex
}

func (j *JsonDataEntity) New(jsonData []byte) error {
	defer j.Unlock()
	j.Lock()

	j.jsonData = jsonData
	err := json.Unmarshal(jsonData, &j.m)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	err = gob.NewEncoder(buf).Encode(j.m)
	if err != nil {
		return err
	}
	j.gobBuf = buf.Bytes()
	return nil
}
func (j *JsonDataEntity) Parse(v interface{}) {
	json.Unmarshal(j.jsonData, v)
}
func (j *JsonDataEntity) Stringify() string {
	return string(j.jsonData)
}
func (j *JsonDataEntity) CopyMap() map[string]interface{} {
	copyMap := make(map[string]interface{})
	buf := bytes.NewBuffer(j.gobBuf)
	gob.NewDecoder(buf).Decode(&copyMap)
	return copyMap
}

var JsonData = &JsonDataEntity{}

var Config = &struct {
	ViewDir           string              `json:"ViewDir"`           //
	UploadDir         string              `json:"UploadDir"`         //
	ResourcesDir      string              `json:"ResourcesDir"`      //
	ResourcesDirName  string              `json:"ResourcesDirName"`  //
	UploadDirName     string              `json:"UploadDirName"`     //
	LogDir            string              `json:"LogDir"`            //
	DefaultPage       string              `json:"DefaultPage"`       //
	JsonDataPath      string              `json:"JsonDataPath"`      //
	HttpPort          string              `json:"HttpPort"`          //
	HttpsPort         string              `json:"HttpsPort"`         //
	Debug             bool                `json:"Debug"`             //
	Domain            string              `json:"Domain"`            //
	ViewSuffix        string              `json:"ViewSuffix"`        //动态模板，默认文件类型
	ViewActionMapping []ViewActionMapping `json:"ViewActionMapping"` //可以进行映射的文件类型
	TLSCertFile       string              `json:"TLSCertFile"`       //
	TLSKeyFile        string              `json:"TLSKeyFile"`        //
	DBUrl             string              `json:"DBUrl"`             //
	LogServer         string              `json:"LogServer"`         //
	Name              string              `json:"Name"`              //
	Ver               string              `json:"Ver"`               //
	SessionExpires    int                 `json:"SessionExpires"`    //
	SecureKey         string              `json:"SecureKey"`         //
}{}
