package tool

import (
	"encoding/json"
	"sync"
)

var json_operation_locker sync.RWMutex
func JsonMarshal(v interface{}) ([]byte, error)  {
	json_operation_locker.Lock()
	defer json_operation_locker.Unlock()
	return json.Marshal(v)
}
func JsonUnmarshal(data []byte, v interface{}) error  {
	json_operation_locker.Lock()
	defer json_operation_locker.Unlock()
	return json.Unmarshal(data,v)
}