package orm

import (
	mydb "cloud_distributed_storage/Backend/database/mysql"
	"database/sql"
	"github.com/go-acme/lego/v4/log"
)

func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_file (`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) values (?, ?, ?, ?, 1)")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
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
	return

}

// GetFileMeta: Get file meta info
func GetFileMeta(filehash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1, file_name, file_size, file_addr from tbl_file where file_sha1 = ? and status = 1 limit 1")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No records found")
			res.Suc = false
			res.Data = nil
			return
		} else {
			log.Println("Failed to execute statement, err: ", err.Error())
			res.Suc = false
			res.Msg = err.Error()
			return
		}
	}
	res.Suc = true
	res.Data = tfile
	return
}

// GetFileMetaList: Get file meta info list
func GetFileMetaList(limit int) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1, file_name, file_size, file_addr from tbl_file where status = 1 limit ?")
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

	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	var tfiles []TableFile
	for i := 0; i < len(values) && rows.Next(); i++ {
		tfile := TableFile{}
		err := rows.Scan(&tfile.FileHash, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr)
		if err != nil {
			log.Println("Failed to scan row, err: ", err.Error())
			break
		}
		tfiles = append(tfiles, tfile)
	}
	res.Suc = true
	res.Data = tfiles
	return
}

// UpdateFileLocation: Update file location
func UpdateFileLocation(filehash string, fileaddr string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("update tbl_file set file_addr = ? where file_sha1 = ? limit 1")
	if err != nil {
		log.Println("Failed to prepare statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileaddr, filehash)
	if err != nil {
		log.Println("Failed to execute statement, err: ", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			log.Printf("File with hash: %s not found", filehash)
			res.Suc = false
			res.Msg = "File not found"
			return
		}
		res.Suc = true
		return
	}
	res.Suc = false
	res.Msg = err.Error()
	return
}
