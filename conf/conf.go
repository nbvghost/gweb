package conf

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"sync"
)

type ViewActionMapping struct {
	ContentType string `json:"ContentType"`
	Extension   string `json:"Extension"`
}

type JsonDataEntity struct {
	Type   reflect.Type
	GobBuf []byte
	sync.RWMutex
}

/*func (j *JsonDataEntity) New(jsonData []byte) error {
	defer j.Unlock()
	j.Lock()

	j.jsonData = jsonData
	err := json.Unmarshal(jsonData, &j.m)
	if err != nil {
		return err
	}

	return nil
}*/
/*func (j *JsonDataEntity) Parse(v interface{}) {
	json.Unmarshal(j.jsonData, v)
}*/
/*func (j *JsonDataEntity) Stringify() string {
	return string(j.jsonData)
}*/
func (j *JsonDataEntity) Copy() interface{} {
	copy := reflect.New(j.Type)
	buf := bytes.NewBuffer(j.GobBuf)
	gob.NewDecoder(buf).Decode(copy.Interface())
	return copy.Interface()
}

var JsonData = &JsonDataEntity{}

func LoadData(target interface{}) error {
	gob.Register(target)

	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		JsonData.Type = v.Type().Elem()
	} else {
		JsonData.Type = v.Type()
	}

	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(target)
	if err != nil {
		return err
	}
	JsonData.GobBuf = buf.Bytes()
	return nil
}

var Config = &struct {
	//ResourcesDir      string              `json:"ResourcesDir"`      //
	//ResourcesDirName  string              `json:"ResourcesDirName"`  //
	//UploadDirName     string              `json:"UploadDirName"`     //
	//HttpPort          string              `json:"HttpPort"`          //
	//HttpsPort         string              `json:"HttpsPort"`         //
	//LogDir       string `json:"LogDir"`       //
	//TLSCertFile       string              `json:"TLSCertFile"`       //
	//TLSKeyFile        string              `json:"TLSKeyFile"`        //
	//DBUrl             string              `json:"DBUrl"`             //
	//LogServer      string `json:"LogServer"`      //
	//Domain            string              `json:"Domain"`            //
	ViewDir     string `json:"ViewDir"`     //
	UploadDir   string `json:"UploadDir"`   //
	DefaultPage string `json:"DefaultPage"` //
	//JsonDataPath      string              `json:"JsonDataPath"`      //
	Debug             bool                `json:"Debug"`             //
	ViewSuffix        string              `json:"ViewSuffix"`        //动态模板，默认文件类型
	ViewActionMapping []ViewActionMapping `json:"ViewActionMapping"` //可以进行映射的文件类型
	Name              string              `json:"Name"`              //
	Ver               string              `json:"Ver"`               //
	SessionExpires    int                 `json:"SessionExpires"`    //

}{
	ViewDir:     "view",
	UploadDir:   "upload",
	DefaultPage: "index",
	//JsonDataPath:   "config.json",
	ViewSuffix:     ".html",
	Name:           "gweb",
	Ver:            "1.3",
	SessionExpires: 1800,
}
