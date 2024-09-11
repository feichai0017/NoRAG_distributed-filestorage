package config

import "fmt"

var (
	MySQLSource = "root:119742@tcp(localhost:3301)/fileserver?charset=utf8"
)

func UpdateDBHost(host string) {
	MySQLSource = fmt.Sprintf("root:119742@tcp(%s:3301)/fileserver?charset=utf8", host)
	fmt.Println("Updated MySQLSource:", MySQLSource) // Debug log
}
