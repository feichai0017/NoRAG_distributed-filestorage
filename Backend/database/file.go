package database

import (
	mydb "cloud_distributed_storage/database/mysql"
	"database/sql"
	"fmt"
)

// OnFileUploadFinished: file uploaded to database
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"INSERT IGNORE INTO tbl_file (`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) VALUES (?, ?, ?, ?, 1)",
	)
	if err != nil {
		fmt.Printf("failed to prepare statement,err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Printf("failed to execute the sql" + err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("file with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false

}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// GetFileMeta: get filemeta by filehash
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file where file_sha1=? and status=1",
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tfile, nil
}

// UpdateFileLocation: update file location
func UpdateFileLocation(filehash string, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_file set file_addr=? where file_sha1=? limit 1",
	)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileaddr, filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("file with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}
