package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:119742@tcp(127.0.0.1:3301)/fileserver?charset=utf8")
	if err != nil {
		fmt.Printf("database faile to connected")
	}
	db.SetMaxOpenConns(1000)
	err = db.Ping()
	if err != nil {
		fmt.Printf("faile to connect mysql,err:" + err.Error())
		os.Exit(1)
	}
}

// DBConn: return database object
func DBConn() *sql.DB {
	return db
}

// ParseRows: parse *sql.Rows into a slice of maps
func ParseRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 准备容器来接收解析后的结果
	var result []map[string]interface{}

	for rows.Next() {
		// 创建一个临时的容器，用于存储每一行的数据
		columnPointers := make([]interface{}, len(columns))
		columnValues := make([]interface{}, len(columns))
		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}

		// 将行数据扫描到 columnPointers 中
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// 将数据存入 map 中，以列名为键
		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			val := columnValues[i]

			// 将 []byte 类型转换为字符串
			if b, ok := val.([]byte); ok {
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = val
			}
		}

		// 将 map 添加到结果切片中
		result = append(result, rowMap)
	}

	// 检查循环中的错误
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
