package csv

import (
	"encoding/json"
	"fmt"
	"nuanv3/shared/commondata"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func LoadTemplates(templates interface{}) (errs []error) {
	// Load Tables
	LoadTables()

	va := reflect.ValueOf(templates).Elem()
	vt := va.Type()
	for i := 0; i < va.NumField(); i++ {
		vf := va.Field(i)
		vtf := vt.Field(i)

		for j := 0; j < vf.NumField(); j++ {

			cd := newCsvData()
			err := cd.getCsvData(vtf, vf.Type().Field(j))
			if err != nil {
				errs = append(errs, err)
				continue
			}

			err = cd.setActivityConf(vf.Field(j))
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return
}

type CsvData struct {
	AID       string // activity id
	Version   uint8  // activity version
	Table     string // csv file name
	Key       string // key of csv
	Val       string // value of csv
	Type      string // type of csv
	IsVersion bool   // table with version or not
}

func newCsvData() *CsvData {
	return &CsvData{
		AID:       "",
		Version:   0,
		Table:     "",
		Key:       "",
		Val:       "",
		Type:      "",
		IsVersion: true,
	}
}

type (
	String  func(i ...interface{}) (string, error)
	Int8    func(i ...interface{}) (int8, error)
	Int16   func(i ...interface{}) (int16, error)
	Int32   func(i ...interface{}) (int32, error)
	Int     func(i ...interface{}) (int, error)
	Int64   func(i ...interface{}) (int64, error)
	Uint8   func(i ...interface{}) (uint8, error)
	Uint16  func(i ...interface{}) (uint16, error)
	Uint32  func(i ...interface{}) (uint32, error)
	Uint    func(i ...interface{}) (uint, error)
	Uint64  func(i ...interface{}) (uint64, error)
	Float32 func(i ...interface{}) (float32, error)
	Float64 func(i ...interface{}) (float64, error)
	Bool    func(i ...interface{}) (bool, error)
	Rewards func(i ...interface{}) ([]commondata.Reward, error)

	Keys   func() ([]string, error)
	Counts func() (int32, error)

	Day    func() time.Duration // 24 * Hour
	Hour   func() time.Duration
	Minute func() time.Duration
	Second func() time.Duration
	Time   func() time.Time

	Exist func(i interface{}) bool
)

//After getCsvData,then can set activity config correctly.
func (c *CsvData) setActivityConf(rv reflect.Value) error {
	if !rv.CanSet() {
		return fmt.Errorf("%v can not set", rv)
	}
	switch rv.Kind() {
	case reflect.String:
		r, err := c.getString()
		if err != nil {
			return err
		}
		rv.SetString(r)

	case reflect.Int8, reflect.Int32, reflect.Int16, reflect.Int, reflect.Int64:
		r, err := c.getInt64()
		if err != nil {
			return err
		}
		rv.SetInt(r)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
		r, err := c.getUint64()
		if err != nil {
			return err
		}
		rv.SetUint(r)
	case reflect.Float32, reflect.Float64:
		r, err := c.getFloat64()
		if err != nil {
			return err
		}
		rv.SetFloat(r)
	case reflect.Bool:
		r, err := c.getBool()
		if err != nil {
			return err
		}
		rv.SetBool(r)
	case reflect.Slice:
		vs, err := c.getSlice()
		if err != nil {
			return err
		}

		rv.Set(reflect.MakeSlice(rv.Type(), len(vs), len(vs)))

		for i := 0; i < rv.Len(); i++ {
			vsKind := reflect.ValueOf(vs[i]).Kind()

			if vsKind == reflect.Slice || vsKind == reflect.Struct {
				js, err := json.Marshal(vs[i])
				if err != nil {
					return err
				}

				err = setValue(rv.Index(i), string(js))
				if err != nil {
					return err
				}
			} else {
				err := setValue(rv.Index(i), fmt.Sprint(vs[i]))
				if err != nil {
					return err
				}
			}
		}
	case reflect.Struct:
		vm, err := c.getMap()
		if err != nil {
			return err
		}

		for k, v := range vm {
			delete(vm, k)
			vm[strings.ToUpper(strings.Replace(k, "_", "", -1))] = v
		}

		for i := 0; i < rv.NumField(); i++ {
			upperName := strings.ToUpper(rv.Type().Field(i).Name)

			v, ok := vm[upperName]
			if !ok {
				continue
			}

			err := setValue(rv.Field(i), fmt.Sprint(v))
			if err != nil {
				return err
			}
		}
	case reflect.Func:
		switch c.Type {
		case "String":
			rv.Set(reflect.ValueOf(c.getString))
		case "Int8":
			rv.Set(reflect.ValueOf(c.getInt8))
		case "Int16":
			rv.Set(reflect.ValueOf(c.getInt16))
		case "Int32":
			rv.Set(reflect.ValueOf(c.getInt32))
		case "Int":
			rv.Set(reflect.ValueOf(c.getInt))
		case "Int64":
			rv.Set(reflect.ValueOf(c.getInt64))
		case "Uint8":
			rv.Set(reflect.ValueOf(c.getUint8))
		case "Uint16":
			rv.Set(reflect.ValueOf(c.getUint16))
		case "Uint32":
			rv.Set(reflect.ValueOf(c.getUint32))
		case "Uint":
			rv.Set(reflect.ValueOf(c.getUint))
		case "Uint64":
			rv.Set(reflect.ValueOf(c.getUint64))
		case "Float32":
			rv.Set(reflect.ValueOf(c.getFloat32))
		case "Float64":
			rv.Set(reflect.ValueOf(c.getFloat64))
		case "Bool":
			rv.Set(reflect.ValueOf(c.getBool))
		case "Keys":
			rv.Set(reflect.ValueOf(c.getKeys))
		case "Counts":
			rv.Set(reflect.ValueOf(c.getCounts))
		case "Exist":
			rv.Set(reflect.ValueOf(c.getExist))
		case "Day":
			return c.getDay(rv)
		case "Hour":
			return c.getHour(rv)
		case "Minute":
			return c.getMinute(rv)
		case "Second":
			return c.getSecond(rv)
		case "Time":
			return c.getTime(rv)
		case "Rewards":
			rv.Set(reflect.ValueOf(c.getRewards))
		}
	}

	return nil
}

func (c *CsvData) getCsvData(aid, base reflect.StructField) error {
	err := c.getAID(aid.Tag)
	if err != nil {
		return err
	}
	err = c.getBase(base.Tag)
	if err != nil {
		return err
	}
	err = c.getVer()
	if err != nil {
		return err
	}
	c.Type = base.Type.Name()
	return nil
}

//Get Table,Key and Val of CsvData
func (c *CsvData) getBase(rs reflect.StructTag) error {

	table := rs.Get("table")
	if table == "" {
		return fmt.Errorf("Tag:%v table not set", rs)
	}

	key := rs.Get("key")
	val := rs.Get("val")
	if key == "" && val == "" {
		//return fmt.Errorf("Tag:%v key and val not set", rs)
	}

	if val == "" {
		val = "val"
	}

	if rs.Get("version") == "false" {
		c.IsVersion = false
	}

	c.Table = table
	c.Key = key
	c.Val = val
	return nil
}

//Get AID of CsvData.
func (c *CsvData) getAID(rs reflect.StructTag) error {
	// aid := rs.Get("id")
	// if aid == "" {
	// 	return fmt.Errorf("Tag:%v id not set", rs)
	// }
	c.AID = "1"
	return nil
}

//Get Version of CsvData.
func (c *CsvData) getVer() error {
	// CZXDO: 除去Version
	c.Version = 0
	c.IsVersion = false

	return nil
}

//Get table with version.
func (c *CsvData) getTableWithVer() string {
	if c.Version == 0 || !c.IsVersion {
		return c.Table
	}
	return fmt.Sprintf("%s%d", c.Table, c.Version)
}

//Get table all keys
func (c *CsvData) getKeys() ([]string, error) {
	keys, ret := GetAllRowKeys(c.getTableWithVer())
	if ret != GAMEDATA_OK {
		return keys, fmt.Errorf("GetKeys(%s) failed:%d", c.getTableWithVer(), ret)
	}
	return keys, nil
}

//Get table count of keys
func (c *CsvData) getCounts() (int32, error) {
	count := Count(c.getTableWithVer())
	if count == 0 {
		return count, fmt.Errorf("GetCounts(%s) failed:not found any key", c.getTableWithVer())
	}
	return count, nil
}

//Get exist of keys
func (c *CsvData) getExist(i interface{}) bool {
	ret := KeyExists(c.getTableWithVer(), fmt.Sprint(i))
	if ret == GAMEDATA_OK {
		return true
	}
	return false
}

func (c *CsvData) getDay(rv reflect.Value) error {
	// check exist and format valid
	r, err := c.getInt64()
	if err != nil {
		return err
	}

	rv.Set(reflect.ValueOf(func() time.Duration {
		return time.Duration(r) * time.Hour * 24
	}))

	return nil
}

func (c *CsvData) getHour(rv reflect.Value) error {
	// check exist and format valid
	r, err := c.getInt64()
	if err != nil {
		return err
	}

	rv.Set(reflect.ValueOf(func() time.Duration {
		return time.Duration(r) * time.Hour
	}))

	return nil
}

func (c *CsvData) getMinute(rv reflect.Value) error {
	// check exist and format valid
	r, err := c.getInt64()
	if err != nil {
		return err
	}

	rv.Set(reflect.ValueOf(func() time.Duration {
		return time.Duration(r) * time.Minute
	}))

	return nil
}

func (c *CsvData) getSecond(rv reflect.Value) error {
	// check exist and format valid
	r, err := c.getInt64()
	if err != nil {
		return err
	}

	rv.Set(reflect.ValueOf(func() time.Duration {
		return time.Duration(r) * time.Second
	}))

	return nil
}

func (c *CsvData) getTime(rv reflect.Value) error {
	// check exist and format valid
	/*
		r, err := c.getString()
		if err != nil {
			return err
		}

		t, err := time.ParseInLocation(su.InternalDateFormat, r, time.Local)
		if err != nil {
			logger.Error("time.ParseInLocation(%s, %s) error(%s)", su.InternalDateFormat, r)
			return err
		}

		rv.Set(reflect.ValueOf(func() time.Time {
			return t
		}))
	*/
	return nil
}

//Get string value by CsvData and it's base and version data.
func (c *CsvData) getStringBase(k []interface{}) (string, error) {
	key := c.Key
	if len(k) > 1 {
		return "", fmt.Errorf("In parameters number:%v more than limited number:1", len(k))
	} else if len(k) == 1 {
		key = fmt.Sprint(k[0])
	}

	s, ret := GetString(c.getTableWithVer(), key, c.Val)
	if ret != GAMEDATA_OK {
		if c.Version == 0 && c.IsVersion {
			//log.("CsvData table:%s key:%s value:%s ,fail to find data and change ver 0 to 1 try loading csv data again", c.Table, c.Key, c.Val)
			c.Version = 1
			s, ret := GetString(c.getTableWithVer(), key, c.Val)
			if ret != GAMEDATA_OK {
				return s, fmt.Errorf("GetInt(%s, %s, %s) failed:%d", c.getTableWithVer(), c.Key, c.Val, ret)
			}
			return s, nil
		} else {
			return s, fmt.Errorf("GetInt(%s, %s, %s) failed:%d", c.getTableWithVer(), c.Key, c.Val, ret)
		}
	}
	return s, nil
}

/* * * * * * * * * * * * * * * * * * * * * *
Get different types value by getStringBase.
* * * * * * * * * * * * * * * * * * * * * */

func (c *CsvData) getString(k ...interface{}) (string, error) {
	s, err := c.getStringBase(k)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (c *CsvData) getUint8(k ...interface{}) (uint8, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return 0, err
	}
	u, err := strconv.ParseUint(r, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(u), nil
}

func (c *CsvData) getUint16(k ...interface{}) (uint16, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return 0, err
	}
	u, err := strconv.ParseUint(r, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(u), nil
}

func (c *CsvData) getUint32(k ...interface{}) (uint32, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return 0, err
	}
	u, err := strconv.ParseUint(r, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(u), nil
}

func (c *CsvData) getUint(k ...interface{}) (uint, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return 0, err
	}
	u, err := strconv.ParseUint(r, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(u), nil
}

func (c *CsvData) getUint64(k ...interface{}) (uint64, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return 0, err
	}
	u, err := strconv.ParseUint(r, 10, 64)
	if err != nil {
		return 0, err
	}
	return u, nil
}

func (c *CsvData) getInt8(k ...interface{}) (int8, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(r, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(i), nil
}

func (c *CsvData) getInt16(k ...interface{}) (int16, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(r, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(i), nil
}

func (c *CsvData) getInt32(k ...interface{}) (int32, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(r, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

func (c *CsvData) getInt(k ...interface{}) (int, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func (c *CsvData) getInt64(k ...interface{}) (int64, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (c *CsvData) getFloat32(k ...interface{}) (float32, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(r, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func (c *CsvData) getFloat64(k ...interface{}) (float64, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(r, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func (c *CsvData) getBool(k ...interface{}) (bool, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return false, err
	}
	f, err := strconv.ParseBool(r)
	if err != nil {
		return f, err
	}
	return f, nil
}

func (c *CsvData) getSlice(k ...interface{}) ([]interface{}, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return nil, err
	}
	var s []interface{}
	err = json.Unmarshal([]byte(r), &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (c *CsvData) getMap(k ...interface{}) (map[string]interface{}, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal([]byte(r), &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

func (c *CsvData) getRewards(k ...interface{}) ([]commondata.Reward, error) {
	r, err := c.getStringBase(k)
	if err != nil {
		return nil, err
	}

	var rewards []commondata.Reward

	err = json.Unmarshal([]byte(r), &rewards)
	if err != nil {
		return nil, err
	}

	return rewards, nil
}

func setValue(rv reflect.Value, val string) error {
	if !rv.CanSet() {
		return nil
	}
	switch rv.Kind() {
	case reflect.String:
		rv.SetString(val)
	case reflect.Int8, reflect.Int32, reflect.Int16, reflect.Int, reflect.Int64:
		vi, err := valueInt64(val)
		if err != nil {
			return err
		}
		rv.SetInt(vi)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
		vu, err := valueUint64(val)
		if err != nil {
			return err
		}
		rv.SetUint(vu)
	case reflect.Float32, reflect.Float64:
		vf, err := valueFloat64(val)
		if err != nil {
			return err
		}
		rv.SetFloat(vf)
	case reflect.Bool:
		vb, err := valueBool(val)
		if err != nil {
			return err
		}
		rv.SetBool(vb)
	case reflect.Slice:
		vs, err := valueSlice(val)
		if err != nil {
			return err
		}

		rv.Set(reflect.MakeSlice(rv.Type(), len(vs), len(vs)))

		for i := 0; i < rv.Len(); i++ {
			vsKind := reflect.ValueOf(vs[i]).Kind()

			if vsKind == reflect.Slice || vsKind == reflect.Struct {
				js, err := json.Marshal(vs[i])
				if err != nil {
					return err
				}

				err = setValue(rv.Index(i), string(js))
				if err != nil {
					return err
				}
			} else {
				err := setValue(rv.Index(i), fmt.Sprint(vs[i]))
				if err != nil {
					return err
				}
			}
		}
	case reflect.Struct:
		vm, err := valueMap(val)
		if err != nil {
			return err
		}

		for k, v := range vm {
			delete(vm, k)
			vm[strings.ToUpper(strings.Replace(k, "_", "", -1))] = v
		}

		for i := 0; i < rv.NumField(); i++ {
			upperName := strings.ToUpper(rv.Type().Field(i).Name)

			v, ok := vm[upperName]
			if !ok {
				continue
			}

			err := setValue(rv.Field(i), fmt.Sprint(v))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func valueUint64(val string) (uint64, error) {
	u, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return u, err
	}
	return u, nil
}

func valueInt64(val string) (int64, error) {
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return i, err
	}
	return i, nil
}

func valueFloat64(val string) (float64, error) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return f, err
	}
	return f, nil
}

func valueBool(val string) (bool, error) {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return b, err
	}
	return b, nil
}

func valueSlice(val string) ([]interface{}, error) {
	var s []interface{}
	err := json.Unmarshal([]byte(val), &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func valueMap(val string) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(val), &m)
	if err != nil {
		return m, err
	}
	return m, nil
}
