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
		//log.Println(t.Field(fi).Type.Kind())
		//log.Println(t.Field(fi).Type.Name(), "Name")
		//log.Println(t.Field(fi).Type)
		if subT := t.Field(fi).Type; subT.Kind() == reflect.Struct {
			subN := t.Field(fi).Type.Name()
			switch subN {
			case "Time":
				fields = append(fields, t.Field(fi).Name)
			default:
				sfs := GetAllFieldName(subT)
				fields = append(fields, sfs...)
			}
		} else {
			fields = append(fields, t.Field(fi).Name)
		}

		fi++
	}
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

		if vtarget.FieldByName(sourceFields[index]).IsValid(){
			///log.Println("sourceFields",sourceFields[index],vsource.FieldByName(sourceFields[index]),vtarget.FieldByName(sourceFields[index]))
			//log.Println("sourceFields",sourceFields[index],vsource.FieldByName(sourceFields[index]).CanInterface(),vtarget.FieldByName(sourceFields[index]).CanInterface())

			if vsource.FieldByName(sourceFields[index]).CanInterface() && vtarget.FieldByName(sourceFields[index]).CanInterface(){
				vsData := vsource.FieldByName(sourceFields[index]).Interface()
				vtData := vtarget.FieldByName(sourceFields[index]).Interface()
				//log.Println(sourceFields[index],vsData, vsData, reflect.DeepEqual(vsData, vtData))
				if reflect.DeepEqual(vsData, vtData) == false {
					changeMap[sourceFields[index]] = vsData
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

		if vtarget.FieldByName(sourceFields[index]).IsValid(){
			//log.Println("sourceFields",sourceFields[index],vsource.FieldByName(sourceFields[index]),vtarget.FieldByName(sourceFields[index]))
			//log.Println("sourceFields",sourceFields[index],vsource.FieldByName(sourceFields[index]).CanInterface(),vtarget.FieldByName(sourceFields[index]).CanInterface())

			if vsource.FieldByName(sourceFields[index]).CanInterface() && vtarget.FieldByName(sourceFields[index]).CanInterface(){
				vsData := vsource.FieldByName(sourceFields[index]).Interface()
				vtData := vtarget.FieldByName(sourceFields[index]).Interface()
				//log.Println(sourceFields[index],vsData, vsData, reflect.DeepEqual(vsData, vtData))
				if reflect.DeepEqual(vsData, vtData) == false {
					changeMap[sourceFields[index]] = vsData
				}



				vtarget.FieldByName(sourceFields[index]).Set(vsource.FieldByName(sourceFields[index]))

			}


		}



	}
	return changeMap

}