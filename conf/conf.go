package conf
type mConfiger struct {
	ViewDir           string `json:"ViewDir"`
	UploadDir      string `json:"UploadDir"`
	ResourcesDir      string `json:"ResourcesDir"`
	ResourcesDirName      string `json:"ResourcesDirName"`
	UploadDirName      string `json:"UploadDirName"`
	DefaultPage       string `json:"DefaultPage"`
	JsonDataPath      string `json:"JsonDataPath"`
	HttpPort          string `json:"HttpPort"`
	HttpsPort         string `json:"HttpsPort"`
	Debug         bool `json:"Debug"`
	Domain         string `json:"Domain"`
	ViewSuffix        string `json:"ViewSuffix"`
	ViewActionMapping bool   `json:"ViewActionMapping"`
	TLSCertFile       string `json:"TLSCertFile"`
	TLSKeyFile        string `json:"TLSKeyFile"`
	DBUrl        string `json:"DBUrl"`
}

var JsonData = make(map[string]interface{})
var Config mConfiger
