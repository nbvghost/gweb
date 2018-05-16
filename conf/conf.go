package conf
type mConfiger struct {
	ViewDir           string `json:"ViewDir"`
	ResourcesDir      string `json:"ResourcesDir"`
	DefaultPage       string `json:"DefaultPage"`
	JsonDataPath      string `json:"JsonDataPath"`
	HttpPort          string `json:"HttpPort"`
	HttpsPort         string `json:"HttpsPort"`
	Domain         string `json:"Domain"`
	ViewSuffix        string `json:"ViewSuffix"`
	ViewActionMapping bool   `json:"ViewActionMapping"`
	TLSCertFile       string `json:"TLSCertFile"`
	TLSKeyFile        string `json:"TLSKeyFile"`
	DBUrl        string `json:"DBUrl"`
}

var JsonData = make(map[string]interface{})
var Config mConfiger
