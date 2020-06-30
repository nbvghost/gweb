package tool

import (
	"reflect"
)

func GetAllFieldName(t reflect.Type) []string {
	fields := make([]string, 0)
	var fieldNum = t.NumField()
	var fi = 0
	for {
		if fi > (fieldNum - 1) {
			break
		}
		field:=t.Field(fi)
		fieldType:=field.Type
		//log.Println(fieldType.Kind())
		//log.Println(fieldType.Name(), "Name")
		//log.Println(fieldType)
		if subT := fieldType; subT.Kind() == reflect.Struct {

			subN := fieldType.Name()
			switch subN {
			case "Time":
				fields = append(fields, field.Name)
			default:
				sfs := GetAllFieldName(subT)
				//fmt.Println(sfs)
				fields = append(fields, sfs...)
			}
		} else {
			fields = append(fields, field.Name)
		}

		fi++
	}
	//fmt.Println(fields)
	return fields
}
func FindChange(source interface{}, target interface{}) map[string]interface{} {
	//第一步,先将结构体转化为map方便后续遍历
	tsource := reflect.TypeOf(source).Elem()

	vsource := reflect.ValueOf(source).Elem()
	vtarget := reflect.ValueOf(target).Elem()
	//ttarget := reflect.TypeOf(target).Elem()
	sourceFields := GetAllFieldName(tsource)
	//bmap := getSubField(ttarget)

	changeMap := make(map[string]interface{})


	//开始遍历A结构体的字段
	for index := range sourceFields {

		//vsource.FieldByName(sourceFields[index]).Interface()
		field:=sourceFields[index]
		vtV:=vtarget.FieldByName(field)
		vsV:=vsource.FieldByName(field)

		if vtV.IsValid(){
			///log.Println("sourceFields",sourceFields[index],vsource.FieldByName(sourceFields[index]),vtarget.FieldByName(sourceFields[index]))
			//log.Println("sourceFields",sourceFields[index],vsource.FieldByName(sourceFields[index]).CanInterface(),vtarget.FieldByName(sourceFields[index]).CanInterface())
			if vsV.CanInterface() && vtV.CanInterface(){
				vsData := vsV.Interface()
				vtData := vtV.Interface()
				//log.Println(sourceFields[index],vsData, vsData, reflect.DeepEqual(vsData, vtData))
				if reflect.DeepEqual(vsData, vtData) == false {
					changeMap[field] = vsData
				}
			}
		}
	}
	return changeMap
}
func CopyAndChange(source interface{}, target interface{}) map[string]interface{} {

	//第一步,先将结构体转化为map方便后续遍历
	tsource := reflect.TypeOf(source).Elem()

	vsource := reflect.ValueOf(source).Elem()
	vtarget := reflect.ValueOf(target).Elem()
	//ttarget := reflect.TypeOf(target).Elem()
	sourceFields := GetAllFieldName(tsource)
	//bmap := getSubField(ttarget)

	changeMap := make(map[string]interface{})


	//开始遍历A结构体的字段
	for index := range sourceFields {
		//vsource.FieldByName(sourceFields[index]).Interface()
		field:=sourceFields[index]
		vtV:=vtarget.FieldByName(field)
		vsV:=vsource.FieldByName(field)

		if vtV.IsValid(){
			//log.Println("sourceFields",sourceFields[index],vsource.FieldByName(sourceFields[index]),vtarget.FieldByName(sourceFields[index]))
			//log.Println("sourceFields",sourceFields[index],vsource.FieldByName(sourceFields[index]).CanInterface(),vtarget.FieldByName(sourceFields[index]).CanInterface())

			if vsV.CanInterface() && vtV.CanInterface(){
				vsData := vsV.Interface()
				vtData := vtV.Interface()
				//log.Println(sourceFields[index],vsData, vsData, reflect.DeepEqual(vsData, vtData))
				if reflect.DeepEqual(vsData, vtData) == false {
					changeMap[field] = vsData
				}
				//vtarget.FieldByName(field).Set(vsource.FieldByName(sourceFields[index]))
				vtV.Set(vsV)
			}
		}
	}
	return changeMap

}