package orm

import (
	mydb "cloud_distributed_storage/Backend/service/dbproxy/conn"
	"log"
	"time"
)

// OnUserFileUploadFinished 当用户文件上传完成时调用
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

// QueryUserFileMetas 查询用户文件元信息
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

// DeleteUserFile 删除用户文件（软删除）
func DeleteUserFile(username string, filehash string) (res ExecResult) {
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

// RenameFileName 重命名用户文件
func RenameFileName(username, filehash, filename string) (res ExecResult) {
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

// QueryUserFileMeta 查询单个用户文件元信息
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

// RestoreUserFile 恢复已删除的用户文件
func RestoreUserFile(userName string, fileHash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(`
        UPDATE tbl_user_file
        SET status = 0, last_update = ?
        WHERE user_name = ? AND file_sha1 = ? AND status = 2
    `)
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now(), userName, fileHash)
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	res.Suc = true
	return
}

// QueryUserFilesByStatus 根据状态查询用户文件
func QueryUserFilesByStatus(userName string, status int, limit int) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(`
        SELECT file_sha1, file_name, file_size, upload_at, last_update
        FROM tbl_user_file
        WHERE user_name = ? AND status = ?
        ORDER BY last_update DESC
        LIMIT ?
    `)
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(userName, status, limit)
	if err != nil {
		log.Println("Failed to execute query, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer rows.Close()

	var userFiles []TableUserFile
	for rows.Next() {
		ufile := TableUserFile{}
		err := rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			log.Println("Failed to scan row, err: ", err.Error())
			continue
		}
		userFiles = append(userFiles, ufile)
	}

	res.Suc = true
	res.Data = userFiles
	return
}
