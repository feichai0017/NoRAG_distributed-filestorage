package database

import (
	mydb "cloud_distributed_storage/database/mysql"
	"fmt"
)

// UserSignup: for user signup
func UserSignup(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user(`user_name`,`user_pwd`)" +
		"values(?,?)")
	if err != nil {
		fmt.Println("failed to insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("failed to insert, err:" + err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err && rf > 0 {
		return true
	}
	return false
}

// UserSignIn: for user login
func UserSignIn(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name = ? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Printf("usrname not found:" + username)
		return false
	}

	pRows, err := mydb.ParseRows(rows)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if len(pRows) > 0 && pRows[0]["user_pwd"].(string) == encpwd {
		return true
	}
	return false
}

// UpdateToken: update username and token in mysql
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token(`user_name`,`user_token`)values(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DBConn().Prepare(
		"select user_name, signup_at from tbl_user where user = ? limit 1")

	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
