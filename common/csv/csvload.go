package csv

import (
	"adventure/advserver/config"
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	GAMEDATA_OK = iota

	GAMEDATA_FORMAT_ERROR
	GAMEDATA_NOT_EXISTS
)

type Record struct {
	Fields map[string]string
}

type Table struct {
	Records map[string]*Record
}

var (
	tables map[string]*Table
	lock   sync.RWMutex
)

func LoadTables() {
	lock.Lock()
	defer lock.Unlock()

	timenow := time.Now()

	tables = make(map[string]*Table)

	config := config.GetConfig()

	pattern := config.GameData + "/*.csv"

	//log.Printf("Loading GameData From", pattern)
	files, err := filepath.Glob(pattern)

	if err != nil {
		log.Printf("failed to Glob the GameData %v", err)
		panic(err)
	}

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			log.Printf("error opening file", err)
			continue
		}

		err = parse(file)
		if err != nil {
			file.Close()
			panic(err)
		}
		file.Close()
	}

	log.Printf(fmt.Sprintf("[ %v, CSV(s) Loaded, take %v secs] \n", len(tables), time.Now().Sub(timenow).Seconds()))
}

func parse(file *os.File) error {

	// csv 读取器
	csv_reader := csv.NewReader(file)
	records, err := csv_reader.ReadAll()
	if err != nil {
		log.Printf("cannot parse csv file.", file.Name(), err)
		return errors.New(fmt.Sprint("cannot parse csv file.", file.Name(), err))
	}

	// 是否为空档
	if len(records) == 0 {
		log.Printf("csv file is empty", file.Name())
		return errors.New(fmt.Sprint("csv file is empty", file.Name()))
	}

	// 处理表名
	fi, err := file.Stat()
	if err != nil {
		log.Printf("cannot stat the file", file.Name())
		return errors.New(fmt.Sprint("cannot stat the file", file.Name()))
	}
	tblname := strings.TrimSuffix(fi.Name(), path.Ext(file.Name()))
	//

	// 记录数据, 第一行为表头，因此从第二行开始
	for line := 1; line < len(records); line++ {
		for field := 1; field < len(records[line]); field++ { // 每条记录的第一个字段作为行索引
			set(tblname, records[line][0], records[0][field], records[line][field])
		}
	}
	//log.Printf("Config file:%s loading completed", tblname)

	return nil
}

//---------------------------------------------------------- Set Field value
func set(tblname string, rowname string, fieldname string, value string) {
	//lock.Lock()
	//defer lock.Unlock()

	tbl := tables[tblname]

	if tbl == nil {
		tbl = &Table{}
		tbl.Records = make(map[string]*Record)
		tables[tblname] = tbl
	}

	rec := tbl.Records[rowname]
	if rec == nil {
		rec = &Record{}
		rec.Fields = make(map[string]string)
		tbl.Records[rowname] = rec
	}

	rec.Fields[fieldname] = value
}

//---------------------------------------------------------- Get Field value

func getSubSringNoKey(tblname string, rowname string) (string, int) {
	lock.RLock()
	defer lock.RUnlock()

	tbl, ok := tables[tblname]
	if !ok {
		log.Printf(fmt.Sprint("table ", tblname, " not exists!"))
		return "", GAMEDATA_NOT_EXISTS
	}

	for _, rec := range tbl.Records {

		dirtyWord, ok := rec.Fields["word"]
		if !ok {
			log.Printf(fmt.Sprint("table ", tblname, " field word not exists!"))
			return "", GAMEDATA_NOT_EXISTS
		}

		isSub := strings.Contains(rowname, dirtyWord)
		if isSub == true {
			log.Printf(fmt.Sprint("table ", tblname, " dirtyWord ", dirtyWord, " exists a substring for name ", rowname))
			return rowname, GAMEDATA_OK
		}
	}
	return rowname, GAMEDATA_NOT_EXISTS
}

func getSubSring(tblname string, rowname string) (string, int) {
	lock.RLock()
	defer lock.RUnlock()

	tbl, ok := tables[tblname]
	if !ok {
		log.Printf(fmt.Sprint("table ", tblname, " not exists!"))
		return "", GAMEDATA_NOT_EXISTS
	}

	_, ok = tbl.Records[rowname]
	if ok {
		log.Printf(fmt.Sprint("table ", tblname, " key ", rowname, " exists!"))
		return rowname, GAMEDATA_OK
	}
	for _, rec := range tbl.Records {

		dirtyWord, ok := rec.Fields["word"]
		if !ok {
			log.Printf(fmt.Sprint("table ", tblname, " field word not exists!"))
			return "", GAMEDATA_NOT_EXISTS
		}

		isSub := strings.Contains(rowname, dirtyWord)
		if isSub == true {
			log.Printf(fmt.Sprint("table ", tblname, " dirtyWord ", dirtyWord, " exists a substring for name ", rowname))
			return rowname, GAMEDATA_OK
		}
	}
	return rowname, GAMEDATA_NOT_EXISTS
}

func GetAllRowKeys(tblname string) ([]string, int) {
	lock.RLock()
	defer lock.RUnlock()

	rowKeys := []string{}

	tbl, ok := tables[tblname]
	if !ok {
		log.Printf(fmt.Sprint("table ", tblname, " not exists!"))
		return rowKeys, GAMEDATA_NOT_EXISTS
	}

	for key, _ := range tbl.Records {
		rowKeys = append(rowKeys, key)
	}

	return rowKeys, GAMEDATA_OK
}

func GetAllRowValues(tblname string, fieldname string) ([]string, int) {
	lock.RLock()
	defer lock.RUnlock()

	rowValues := []string{}

	tbl, ok := tables[tblname]
	if !ok {
		log.Printf(fmt.Sprint("table ", tblname, " not exists!"))
		return rowValues, GAMEDATA_NOT_EXISTS
	}

	for _, rec := range tbl.Records {
		rows, ok := rec.Fields[fieldname]
		if !ok {
			log.Printf(fmt.Sprint("table ", tblname, " field ", fieldname, " not exists!"))
			return rowValues, GAMEDATA_NOT_EXISTS
		}

		rowValues = append(rowValues, rows)
	}

	return rowValues, GAMEDATA_OK
}

func GetRowValues(tblname string, rowname string) ([]string, int) {
	lock.RLock()
	defer lock.RUnlock()

	rowValues := []string{}

	tbl, ok := tables[tblname]
	if !ok {
		log.Printf(fmt.Sprint("table ", tblname, " not exists!"))
		return rowValues, GAMEDATA_NOT_EXISTS
	}

	for _, field := range tbl.Records[rowname].Fields {
		rowValues = append(rowValues, field)
	}

	return rowValues, GAMEDATA_OK
}

func get(tblname string, rowname string, fieldname string) (string, int) {
	lock.RLock()
	defer lock.RUnlock()

	tbl, ok := tables[tblname]
	if !ok {
		log.Printf(fmt.Sprint("table ", tblname, " not exists!"))
		return "", GAMEDATA_NOT_EXISTS
	}

	rec, ok := tbl.Records[rowname]
	if !ok {
		log.Printf(fmt.Sprint("table ", tblname, " row ", rowname, " not exists!"))
		return "", GAMEDATA_NOT_EXISTS
	}

	value, ok := rec.Fields[fieldname]
	if !ok {
		log.Printf(fmt.Sprint("table ", tblname, " field ", fieldname, " not exists!"))
		return "", GAMEDATA_NOT_EXISTS
	}
	return value, GAMEDATA_OK
}

//---------------------------------------------------------- Get Field value as Integer
func GetInt(tblname string, rowname string, fieldname string) (int32, int) {
	val, ret := get(tblname, rowname, fieldname)
	if ret != GAMEDATA_OK {
		return 0, ret
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		log.Printf(fmt.Sprintf("cannot parse integer from gamedata %v %v %v %v\n", tblname, rowname, fieldname, err))
		return 0, GAMEDATA_FORMAT_ERROR
	}

	return int32(v), GAMEDATA_OK
}

func KeyExists(tblname string, rowname string) int {
	lock.RLock()
	defer lock.RUnlock()

	tbl, ok := tables[tblname]
	if !ok {
		return GAMEDATA_NOT_EXISTS
	}

	_, ok = tbl.Records[rowname]
	if !ok {
		return GAMEDATA_NOT_EXISTS
	}

	return GAMEDATA_OK
}

func CheckDirtyKey(tblname string, rowname string) int {
	_, ret := getSubSringNoKey(tblname, rowname)
	if ret != GAMEDATA_OK {
		return ret
	}

	return GAMEDATA_OK
}

func CheckKey(tblname string, rowname string) int {
	_, ret := getSubSring(tblname, rowname)
	if ret != GAMEDATA_OK {
		return ret
	}

	return GAMEDATA_OK
}

func GetUint(tblname string, rowname string, fieldname string) (uint, int) {
	val, ret := get(tblname, rowname, fieldname)
	if ret != GAMEDATA_OK {
		return 0, ret
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		log.Printf(fmt.Sprintf("cannot parse unsigned integer from gamedata %v %v %v %v %v, 0 set as default", tblname, rowname, fieldname, val, err))
		return 0, GAMEDATA_OK
	}

	return uint(v), GAMEDATA_OK
}

func GetMaxKeyID(tblname string) (int, int) {
	IDs, ret := GetAllRowKeys(tblname)
	if ret != GAMEDATA_OK {
		log.Printf("GetAllRowKeys failed")
		return 0, GAMEDATA_OK
	}
	var max = 0

	for _, id := range IDs {
		v, err := strconv.Atoi(id)
		if err != nil {
			log.Printf(fmt.Sprintf("cannot parse unsigned integer from gamedata %v %v %v , 0 set as default", tblname, id, err))
			return 0, GAMEDATA_OK
		}

		if v > max {
			max = v
		}
	}

	return max, GAMEDATA_OK
}

//---------------------------------------------------------- Get Field value as Float
func GetFloat(tblname string, rowname string, fieldname string) (float64, int) {
	val, ret := get(tblname, rowname, fieldname)
	if ret != GAMEDATA_OK {
		return 0.0, ret
	}

	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Printf(fmt.Sprintf("cannot parse float from gamedata %v %v %v %v, 0.0 set as default", tblname, rowname, fieldname, err))
		return 0.0, GAMEDATA_OK
	}

	return f, GAMEDATA_OK
}

//---------------------------------------------------------- Get Field value as string
func GetString(tblname string, rowname string, fieldname string) (string, int) {
	return get(tblname, rowname, fieldname)
}

//---------------------------------------------------------- Get Row Count
func Count(tblname string) int32 {
	lock.RLock()
	defer lock.RUnlock()

	tbl := tables[tblname]

	if tbl == nil {
		return 0
	}

	return int32(len(tbl.Records))
}

//---------------------------------------------------------- Test Field Exists
func IsFieldExists(tblname string, fieldname string) bool {
	lock.RLock()
	defer lock.RUnlock()

	tbl := tables[tblname]

	if tbl == nil {
		return false
	}

	key := ""
	// get one record key, the first column
	for k, _ := range tbl.Records {
		key = k
		break
	}

	rec, ok := tbl.Records[key]
	if !ok {
		return false
	}

	_, ok = rec.Fields[fieldname]
	if !ok {
		return false
	}

	return true
}

func GetKeyByRowName(tblname, rowname, fieldname string) (string, int) {
	lock.RLock()
	defer lock.RUnlock()

	tbl, ok := tables[tblname]
	if !ok {
		log.Printf(fmt.Sprint("table ", tblname, " not exists!"))
		return "", GAMEDATA_NOT_EXISTS
	}

	for k, v := range tbl.Records {
		matchV, ok := v.Fields[fieldname]
		if !ok {
			break
		}
		if matchV == rowname {
			return k, GAMEDATA_OK
		}
	}
	log.Printf(fmt.Sprint("table ", tblname, " field ", fieldname, " row ", rowname, " not exists!"))
	return "", GAMEDATA_NOT_EXISTS
}

//---------------------------------------------------------- Load JSON From GameData Directory
func LoadJSON(filename string) ([]byte, error) {
	config := config.GetConfig()
	prefix := config.GameData
	path := prefix + "/" + filename
	return ioutil.ReadFile(path)
}
