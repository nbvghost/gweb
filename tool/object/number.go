package object

import (
	"errors"
	"fmt"
	"github.com/nbvghost/glog"

	"math"

	"strconv"
)

func ParseFloat(value interface{}) float64 {
	if value == nil {
		return 0
	}
	switch value.(type) {
	case int:
		return float64(value.(int))
	case int8:
		return float64(value.(int8))
	case int32:
		return float64(value.(int32))
	case int64:
		return float64(value.(int64))
	case uint:
		return float64(value.(uint))
	case uint8:
		return float64(value.(uint8))
	case uint32:
		return float64(value.(uint32))
	case uint64:
		return float64(value.(uint64))
	case float32:
		return float64(value.(float32))
	case float64:
		return value.(float64)
	case string:
		numberStr := value.(string)
		numb, err := strconv.ParseFloat(numberStr, 64)
		if err != nil {
			glog.Debug(err.Error())
		}
		return numb
	default:
		glog.Error(errors.New("未支持的数据类型：" + fmt.Sprintf("%v", value)))
	}

	return 0
}
func ParseString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch value.(type) {
	case int:
		return strconv.Itoa(value.(int))
	case int8:
		return strconv.Itoa(int(value.(int8)))
	case int32:
		return strconv.Itoa(int(value.(int32)))
	case int64:
		return strconv.Itoa(int(value.(int64)))
	case uint:
		return strconv.Itoa(int(value.(uint)))
	case uint8:
		return strconv.Itoa(int(value.(uint8)))
	case uint32:
		return strconv.Itoa(int(value.(uint32)))
	case uint64:
		u := value.(uint64)
		if u > math.MaxInt64 {
			return strconv.FormatUint(math.MaxInt64, 10)
		}
		return strconv.Itoa(int(u))
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case string:
		numberStr := value.(string)
		return numberStr
	default:
		glog.Error(errors.New("未支持的数据类型：" + fmt.Sprintf("%v", value)))
	}

	return ""
}
func ParseInt(value interface{}) int {
	if value == nil {
		return 0
	}
	switch value.(type) {
	case int:
		return value.(int)
	case int8:
		return int(value.(int8))
	case int32:
		return int(value.(int32))
	case int64:
		return int(value.(int64))
	case uint:
		return int(value.(uint))
	case uint8:
		return int(value.(uint8))
	case uint32:
		return int(value.(uint32))
	case uint64:
		u := value.(uint64)
		if u > math.MaxInt64 {
			return math.MaxInt64
		}
		return int(u)
	case float32:
		return int(value.(float32))
	case float64:
		u := value.(float64)
		if u > math.MaxInt64 {
			return math.MaxInt64
		}
		return int(u)
	case string:
		numberStr := value.(string)
		numb, err := strconv.Atoi(numberStr)
		if err != nil {
			glog.Debug(err.Error())
		}
		return numb
	default:
		glog.Error(errors.New("未支持的数据类型：" + fmt.Sprintf("%v", value)))
	}
	return 0
}
