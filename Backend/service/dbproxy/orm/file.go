package orm

import (
	mydb "cloud_distributed_storage/Backend/database/mysql"
	"database/sql"
	"log"
	"time"
)

func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string, ownerID int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"INSERT INTO tbl_file (`file_sha1`, `file_name`, `file_size`, `file_addr`, `owner_id`, `status`, `create_at`) " +
			"VALUES (?, ?, ?, ?, ?, 1, ?) ON DUPLICATE KEY UPDATE `owner_id`=?, `update_at`=?")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	nowTime := time.Now()
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr, ownerID, nowTime, ownerID, nowTime)
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			log.Printf("File with hash: %s has been uploaded before", filehash)
		}
		res.Suc = true
		return
	}

	res.Suc = false
	res.Msg = "Failed to upload file"
	return
}

func GetFileMeta(filehash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"SELECT id, file_sha1, file_name, file_size, file_addr, owner_id, create_at, update_at, status " +
			"FROM tbl_file WHERE file_sha1 = ? AND status = 1 LIMIT 1")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(
		&tfile.ID, &tfile.FileHash, &tfile.FileName, &tfile.FileSize,
		&tfile.FileAddr, &tfile.OwnerID, &tfile.CreateAt, &tfile.UpdateAt, &tfile.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Suc = false
			res.Msg = "File not found"
		} else {
			log.Println("Failed to execute statement, err: ", err.Error())
			res.Suc = false
			res.Msg = err.Error()
		}
		return
	}

	res.Suc = true
	res.Data = tfile
	return
}

func GetFileMetaList(limit int) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"SELECT id, file_sha1, file_name, file_size, file_addr, owner_id, create_at, update_at, status " +
			"FROM tbl_file WHERE status = 1 LIMIT ?")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit)
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer rows.Close()

	var tfiles []TableFile
	for rows.Next() {
		tfile := TableFile{}
		err := rows.Scan(
			&tfile.ID, &tfile.FileHash, &tfile.FileName, &tfile.FileSize,
			&tfile.FileAddr, &tfile.OwnerID, &tfile.CreateAt, &tfile.UpdateAt, &tfile.Status)
		if err != nil {
			log.Println("Failed to scan row, err: ", err.Error())
			continue
		}
		tfiles = append(tfiles, tfile)
	}

	res.Suc = true
	res.Data = tfiles
	return
}

func UpdateFileLocation(filehash string, fileaddr string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"UPDATE tbl_file SET file_addr = ?, update_at = ? WHERE file_sha1 = ? AND status = 1")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileaddr, time.Now(), filehash)
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			res.Suc = false
			res.Msg = "File not found or already updated"
		} else {
			res.Suc = true
		}
		return
	}

	res.Suc = false
	res.Msg = "Failed to update file location"
	return
}
