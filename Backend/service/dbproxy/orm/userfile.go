package orm

import (
	mydb "cloud_distributed_storage/Backend/database/mysql"
	"database/sql"
	"log"
	"time"
)

// OnUserFileUploadFinished 当用户文件上传完成时调用
func OnUserFileUploadFinished(userID int64, fileID int64, fileName string, fileSize int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"INSERT INTO tbl_user_file (`user_id`, `file_id`, `file_name`, `file_size`, `upload_at`) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, fileID, fileName, fileSize, time.Now())
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
func QueryUserFileMetas(userID int64, limit int) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(`
        SELECT uf.id, uf.file_name, uf.file_size, uf.upload_at, uf.last_update, f.file_sha1
        FROM tbl_user_file uf
        INNER JOIN tbl_file f ON uf.file_id = f.id
        WHERE uf.user_id = ? AND uf.status = 0
        ORDER BY uf.upload_at DESC
        LIMIT ?
    `)
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID, limit)
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
		err := rows.Scan(&ufile.ID, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated, &ufile.FileHash)
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

// DeleteUserFile 删除用户文件（软删除）
func DeleteUserFile(userID int64, fileHash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(`
        UPDATE tbl_user_file uf
        INNER JOIN tbl_file f ON uf.file_id = f.id
        SET uf.status = 2, uf.last_update = ?
        WHERE uf.user_id = ? AND f.file_sha1 = ? AND uf.status = 0
    `)
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now(), userID, fileHash)
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
func RenameFileName(userID int64, fileHash, newFileName string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(`
        UPDATE tbl_user_file uf
        INNER JOIN tbl_file f ON uf.file_id = f.id
        SET uf.file_name = ?, uf.last_update = ?
        WHERE uf.user_id = ? AND f.file_sha1 = ? AND uf.status = 0
    `)
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(newFileName, time.Now(), userID, fileHash)
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
func QueryUserFileMeta(userID int64, fileHash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(`
        SELECT uf.id, uf.file_name, uf.file_size, uf.upload_at, uf.last_update, f.file_sha1
        FROM tbl_user_file uf
        INNER JOIN tbl_file f ON uf.file_id = f.id
        WHERE uf.user_id = ? AND f.file_sha1 = ? AND uf.status = 0
        LIMIT 1
    `)
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ufile := TableUserFile{}
	err = stmt.QueryRow(userID, fileHash).Scan(
		&ufile.ID, &ufile.FileName, &ufile.FileSize,
		&ufile.UploadAt, &ufile.LastUpdated, &ufile.FileHash)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Suc = false
			res.Msg = "File not found"
		} else {
			log.Println("Failed to execute query, err: ", err.Error())
			res.Suc = false
			res.Msg = err.Error()
		}
		return
	}

	res.Suc = true
	res.Data = ufile
	return
}

// RestoreUserFile 恢复已删除的用户文件
func RestoreUserFile(userID int64, fileHash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(`
        UPDATE tbl_user_file uf
        INNER JOIN tbl_file f ON uf.file_id = f.id
        SET uf.status = 0, uf.last_update = ?
        WHERE uf.user_id = ? AND f.file_sha1 = ? AND uf.status = 2
    `)
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now(), userID, fileHash)
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
func QueryUserFilesByStatus(userID int64, status int, limit int) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(`
        SELECT uf.id, uf.file_name, uf.file_size, uf.upload_at, uf.last_update, f.file_sha1
        FROM tbl_user_file uf
        INNER JOIN tbl_file f ON uf.file_id = f.id
        WHERE uf.user_id = ? AND uf.status = ?
        ORDER BY uf.last_update DESC
        LIMIT ?
    `)
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID, status, limit)
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
		err := rows.Scan(&ufile.ID, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated, &ufile.FileHash)
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
