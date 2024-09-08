package orm

import (
	mydb "cloud_distributed_storage/Backend/database/mysql"
	"log"
)

// UserSignup: Insert user info into database
func UserSignup(username string, passwd string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user (`user_name`, `user_pwd`) values (?, ?)")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			log.Printf("User: %s has been registered before", username)
		}
		res.Suc = true
		return
	}
	res.Suc = false
	res.Msg = "Failed to insert user"
	return
}

// UserLogin: Check user info in database
func UserLogin(username string, encpwd string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name = ? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	} else if rows == nil {
		log.Println("Username not found: " + username)
		res.Suc = false
		res.Msg = "Username not found"
		return
	}

	pRows, _ := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		res.Suc = true
		res.Data = true
		return
	}
	res.Suc = false
	res.Msg = "Incorrect password"
	return
}

// UpdateToken: Update user token
func UpdateToken(username string, token string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("replace into tbl_user_token (`user_name`, `user_token`) values (?, ?)")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

// GetUserInfo: Get user info
func GetUserInfo(username string) (res ExecResult) {
	user := TableUser{}
	stmt, err := mydb.DBConn().Prepare("select user_name, signup_at from tbl_user where user_name = ? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	res.Data = user
	return
}

// UserExist: Check if user exists
func UserExist(username string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("select 1 from tbl_user where user_name = ? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	} else if rows == nil {
		log.Println("Username not found: " + username)
		res.Suc = false
		res.Msg = "Username not found"
		return
	}

	res.Suc = true
	res.Data = map[string]bool{
		"exists": rows.Next(),
	}
	return
}
