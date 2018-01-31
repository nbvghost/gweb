package tool

import (
	"net/url"
)

func QueryParams(m url.Values) map[string]string {
	data := make(map[string]string)
	for key, value := range m {
		//util.Trace(key)
		//util.Trace(value)
		data[key] = value[0]
	}
	return data
}
