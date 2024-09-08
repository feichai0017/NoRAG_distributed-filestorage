package mysql

import (
	cfg "cloud_distributed_storage/Backend/service/dbproxy/config"
	"database/sql"
	"fmt"
	"log"
	"os"
)

var db *sql.DB

func InitDBConn() {
	fmt.Println("MySQLSource:", cfg.MySQLSource)
	db, _ = sql.Open("mysql", cfg.MySQLSource)
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Printf("Failed to connect to MySQL, err: %s", err.Error())
		os.Exit(1)
	}
}

// DBConn 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}

// ParseRows 将查询结果解析为 map
func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	var result []map[string]interface{}
	for rows.Next() {
		_ = rows.Scan(scanArgs...)
		entry := make(map[string]interface{})
		for i, col := range values {
			if col != nil {
				entry[columns[i]] = string(col.([]byte))
			}
		}
		result = append(result, entry)
	}
	return result
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
