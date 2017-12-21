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
}
type HeroTemplate struct { // 英雄模板
	HeroName            string  `val:"名字"`
	SkillID             []int32 `val:"技能ID列表"`
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
		structType := vf.Type().Elem() // map value type

		keys, ret := GetAllRowIntKeys(tableName)
		if ret != GAMEDATA_OK {
			fmt.Println("GetAllRowIntKeys Error")
			return
		}

		for _, v := range keys {
			structEntry := reflect.New(structType).Elem()

			loadStructData(tableName, v, structEntry, structType)

			vf.SetMapIndex(reflect.ValueOf(int32(v)), structEntry)
		}
	}

	// Delete tables
	tables = make(map[string]*Table)
}

func loadStructData(tableName string, index int, sVal reflect.Value, sType reflect.Type) {

	for j := 0; j < sVal.NumField(); j++ {
		sef := sVal.Field(j)
		tsf := sType.Field(j)

		switch tsf.Type.Kind() {
		case reflect.Int, reflect.Int32:
			rowName := tsf.Tag.Get("val")
			strValue, ret := GetString(tableName, strconv.Itoa(index), rowName)
			if ret != GAMEDATA_OK {
				fmt.Printf("reflect.Int GetString(%v, %v, %v) Error \n", tableName, rowName, index)
				return
			}
			intValue := 0
			err := errors.New("parse int error")
			if strValue != "" {
				intValue, err = strconv.Atoi(strValue)
				if err != nil {
					fmt.Printf("reflect.Int strconv.Atoi(%v) Error %v \n", strValue, err)
					return
				}
			}
			sef.SetInt(int64(intValue))

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
