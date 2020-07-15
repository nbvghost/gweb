package tool

import (
	"log"
	"reflect"
)

/*
map的深度复制

只支持值类型为：

Array,Slice,元数据类型，不支持struct,map
*/
func DeepCopyMap(source interface{}) (target interface{}) {

	sourceType := reflect.TypeOf(source)

	if sourceType.Kind() != reflect.Map {
		log.Panic("source is not Map Kind")
	}

	sourceValue := reflect.ValueOf(source)
	targetValue := reflect.MakeMap(sourceType)

	mapKey := sourceValue.MapKeys()
	for index := range mapKey {
		elem := sourceValue.MapIndex(mapKey[index]).Elem()
		kind := elem.Kind()

		mapKeyType := reflect.New(mapKey[index].Type())
		mapKeySet := mapKeyType.Elem()
		mapKeySet.Set(mapKey[index])

		switch kind {

		case reflect.Bool,
			reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Float32,
			reflect.Float64,
			reflect.Complex64,
			reflect.Complex128,
			reflect.String:
			targetValue.SetMapIndex(mapKeySet, sourceValue.MapIndex(mapKey[index]).Elem())
		case reflect.Array:
			newSliceValue := sourceValue.MapIndex(mapKey[index])
			newSliceLen := newSliceValue.Elem().Len()
			newSlice := reflect.New(reflect.ArrayOf(newSliceLen, reflect.TypeOf(newSliceValue.Interface()).Elem()))
			reflect.Copy(newSlice.Elem(), newSliceValue.Elem())
			targetValue.SetMapIndex(mapKeySet, newSlice.Elem())
		case reflect.Slice:
			newSliceValue := sourceValue.MapIndex(mapKey[index])
			newSliceLen := newSliceValue.Elem().Len()
			newSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(newSliceValue.Interface()).Elem()), newSliceLen, newSliceLen)
			reflect.Copy(newSlice, newSliceValue.Elem())
			targetValue.SetMapIndex(mapKeySet, newSlice)
		default:
			log.Panic(kind.String() + " kind not support")
		}

	}

	return targetValue.Interface()

}
