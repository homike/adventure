package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	AdvDB *sql.DB
)

func init() {
	mysqlCfg := "root:123456@(127.0.0.1:3306)/adventuremain?parseTime=true&loc=Local&charset=utf8"
	var err error
	if AdvDB, err = sql.Open("mysql", mysqlCfg); err != nil {
		//logger.Error("sql.Open(\"mysql\", %s) failed (%v)", mysqlCfg, err)
		fmt.Println("sql.Open(mysql) error: ", err)
		return
	}

	AdvDB.SetMaxIdleConns(100 / 4)
	AdvDB.SetMaxOpenConns(100)
}
