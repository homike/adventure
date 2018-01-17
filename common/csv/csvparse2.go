package csv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//LoadTemplates2 : Load Tempaltes Lick this:
/*
type Templates struct {
	HeroTemplates        map[int32]structs.HeroTemplate        `table:"hero"`
	GlobalData           structs.GlobalTemplate                `table:"GlobalData"`
}

type HeroTemplate struct { // 英雄模板
	HeroName            string  `val:"名字"`
	SkillID             []int32 `val:"技能ID列表"`
}
type GlobalTemplate struct {
	EmployReturnExp              int32 `val:"英雄.解雇返回经验需要的总经验值"`
	EmployReturnExpPer           int32 `val:"英雄.解雇返回经验比例"`
}
*/
func LoadTemplates2(templates interface{}) {
	// Load Tables
	LoadTables()

	va := reflect.ValueOf(templates).Elem()
	vt := va.Type()
	for i := 0; i < va.NumField(); i++ {
		vf := va.Field(i)
		vtf := vt.Field(i)

		tableName := vtf.Tag.Get("table")

		if vf.Type().Kind() == reflect.Struct { // key-value config
			loadKeyValueData(tableName, vf)

		} else { // map config
			structType := vf.Type().Elem() // map value type (HeroTemplate)

			keys, ret := GetAllRowIntKeys(tableName)
			if ret != GAMEDATA_OK {
				fmt.Println("GetAllRowIntKeys Error")
				return
			}

			for _, key := range keys {
				structEntry := reflect.New(structType).Elem() // new HeroTemplate

				loadStructData(tableName, key, structEntry, structType)

				vf.SetMapIndex(reflect.ValueOf(int32(key)), structEntry) // set map value
			}
		}
	}

	// Delete tables
	tables = make(map[string]*Table)
}

func loadKeyValueData(tableName string, vf reflect.Value) {
	for j := 0; j < vf.NumField(); j++ {
		vff := vf.Field(j)
		vftf := vf.Type().Field(j)

		key := vftf.Tag.Get("val")

		switch vftf.Type.Kind() {
		case reflect.Int8, reflect.Int, reflect.Int32, reflect.Int64:
			intValue, _ := GetInt(tableName, key, "Value")
			vff.SetInt(int64(intValue))

		case reflect.Slice: // Only for slic int
			strValue, ret := GetString(tableName, key, "Value")
			if ret != GAMEDATA_OK {
				fmt.Printf("reflect.Slice GetString(%v, %v) Error \n", tableName, key)
				return
			}

			if vftf.Type.Elem().Kind() == reflect.String {
				strSplits := strings.Split(strValue, ";")
				vff.Set(reflect.ValueOf(strSplits))

			} else {
				arrValue := []int32{}
				if strValue != "" {
					strSplits := strings.Split(strValue, ";")
					for _, v := range strSplits {
						intValue, err := strconv.Atoi(v)
						if err != nil {
							fmt.Printf("reflect.Slice strconv.Atoi(%v) Error %v \n", v, err)
							return
						}
						arrValue = append(arrValue, int32(intValue))
					}
				}
				vff.Set(reflect.ValueOf(arrValue))
			}
		}

	}
}

func loadStructData(tableName string, index int, sVal reflect.Value, sType reflect.Type) {

	for j := 0; j < sVal.NumField(); j++ {
		sef := sVal.Field(j)
		tsf := sType.Field(j)

		switch tsf.Type.Kind() {
		case reflect.Uint8, reflect.Uint64:
			intValue := uint64(0)
			rowName := tsf.Tag.Get("val")
			if rowName == "ID" {
				intValue = uint64(index)
			} else {
				strValue, ret := GetString(tableName, strconv.Itoa(index), rowName)
				if ret != GAMEDATA_OK {
					fmt.Printf("reflect.UInt GetString(%v, %v, %v) Error \n", tableName, rowName, index)
					return
				}
				err := errors.New("parse int error")
				if strValue != "" {
					intValue, err = strconv.ParseUint(strValue, 10, 64)
					if err != nil {
						fmt.Printf("reflect.UInt strconv.Atoi(%v) Error %v \n", strValue, err)
						return
					}
				}
			}
			sef.SetUint(intValue)

		case reflect.Int8, reflect.Int, reflect.Int32, reflect.Int64:
			intValue := int64(0)
			rowName := tsf.Tag.Get("val")
			if rowName == "ID" {
				intValue = int64(index)
			} else {
				strValue, ret := GetString(tableName, strconv.Itoa(index), rowName)
				if ret != GAMEDATA_OK {
					fmt.Printf("reflect.Int GetString(%v, %v, %v) Error \n", tableName, rowName, index)
					return
				}
				err := errors.New("parse int error")
				if strValue != "" {
					intValue, err = strconv.ParseInt(strValue, 10, 64)
					if err != nil {
						fmt.Printf("reflect.Int strconv.Atoi(%v) Error %v \n", strValue, err)
						return
					}
				}
			}
			sef.SetInt(intValue)

		case reflect.String:
			rowName := tsf.Tag.Get("val")
			strValue, ret := GetString(tableName, strconv.Itoa(index), rowName)
			if ret != GAMEDATA_OK {
				fmt.Printf("reflect.String GetString(%v, %v, %v) Error \n", tableName, index, rowName)
				return
			}
			sef.SetString(strValue)

		case reflect.Slice: // Only for slic int
			rowName := tsf.Tag.Get("val")
			strValue, ret := GetString(tableName, strconv.Itoa(index), rowName)
			if ret != GAMEDATA_OK {
				fmt.Printf("reflect.Slice GetString(%v, %v, %v) Error \n", tableName, rowName, index)
				return
			}
			arrValue := []int32{}
			if strValue != "" {
				strSplits := strings.Split(strValue, ";")
				for _, v := range strSplits {
					intValue, err := strconv.Atoi(v)
					if err != nil {
						fmt.Printf("reflect.Slice strconv.Atoi(%v) Error %v \n", v, err)
						return
					}
					arrValue = append(arrValue, int32(intValue))
				}
			}
			sef.Set(reflect.ValueOf(arrValue))
		}
	}
}
