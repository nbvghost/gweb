package conf

type ViewActionMapping struct {
	ContentType string `json:"ContentType"`
	Extension   string `json:"Extension"`
}

var JsonText = ""
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
