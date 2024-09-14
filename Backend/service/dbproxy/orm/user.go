package orm

import (
	mydb "cloud_distributed_storage/Backend/database/mysql"
	"database/sql"
	"log"
	"time"
)

func UserSignup(username, passwd, email, phone string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"INSERT INTO tbl_user (`user_name`, `user_pwd`, `email`, `phone`, `signup_at`) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd, email, phone, time.Now())
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		res.Suc = true
		return
	}

	res.Suc = false
	res.Msg = "Failed to insert user"
	return
}

func UserLogin(username string, encpwd string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("SELECT id, user_pwd FROM tbl_user WHERE user_name = ? LIMIT 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	var (
		userId int64
		dbPwd  string
	)
	err = stmt.QueryRow(username).Scan(&userId, &dbPwd)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Suc = false
			res.Msg = "User not found"
		} else {
			log.Println(err.Error())
			res.Suc = false
			res.Msg = err.Error()
		}
		return
	}

	if dbPwd == encpwd {
		res.Suc = true
		res.Data = userId
	} else {
		res.Suc = false
		res.Msg = "Incorrect password"
	}
	return
}

func UpdateToken(username, token string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"REPLACE INTO tbl_user_token (`user_name`, `user_token`) VALUES (?, ?)")
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

func GetUserInfo(username string) (res ExecResult) {
	user := TableUser{}
	stmt, err := mydb.DBConn().Prepare(
		"SELECT id, user_name, email, phone, email_validated, phone_validated, signup_at, last_active, profile, status " +
			"FROM tbl_user WHERE user_name = ? LIMIT 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone,
		&user.EmailValidated, &user.PhoneValidated, &user.SignupAt,
		&user.LastActive, &user.Profile, &user.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Suc = false
			res.Msg = "User not found"
		} else {
			log.Println(err.Error())
			res.Suc = false
			res.Msg = err.Error()
		}
		return
	}

	res.Suc = true
	res.Data = user
	return
}

func UserExist(username string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("SELECT 1 FROM tbl_user WHERE user_name = ? LIMIT 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	var exists int
	err = stmt.QueryRow(username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	res.Suc = true
	res.Data = exists == 1
	return
}
