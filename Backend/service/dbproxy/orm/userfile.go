package orm

import (
	mydb "cloud_distributed_storage/Backend/database/mysql"
	"log"
	"time"
)

// OnUserFileUploadFinished: User file upload finished
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user_file (`user_name`, `file_sha1`, `file_name`, `file_size`, `status`) values (?, ?, ?, ?, 1)")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

// QueryUserFileMetas: Query user file meta info
func QueryUserFileMetas(username string, limit int) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1, file_name, file_size, upload_at, last_update from tbl_user_file where user_name = ? limit ?")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	var userFiles []TableUserFile
	for rows.Next() {
		ufile := TableUserFile{}
		err := rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		userFiles = append(userFiles, ufile)
	}

	res.Suc = true
	res.Data = userFiles
	return
}

// DeleteUserFile: Delete user file meta info
func DeleteUserFile(username, filehash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("update tbl_user_file set status=2 where user_name=? and file_sha1=? limit 1")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash)
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

// UpdateFileName: Rename file name
func UpdateFileName(username, filehash, filename string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("update tbl_user_file set file_name=? where user_name=? and file_sha1=? limit 1")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(filename, username, filehash)
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

// QueryUserFileMeta: Query user file meta info
func QueryUserFileMeta(username, filehash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1, file_name, file_size, upload_at, last_update from tbl_user_file where user_name = ? and file_sha1 = ? limit 1")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	row := stmt.QueryRow(username, filehash)
	ufile := TableUserFile{}
	err = row.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	res.Suc = true
	res.Data = ufile
	return
}
