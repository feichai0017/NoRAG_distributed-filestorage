package orm

import (
	mydb "cloud_distributed_storage/Backend/service/dbproxy/conn"
	"database/sql"
	"log"
	"time"
)

// GrantPermission 授予权限
func GrantPermission(roleName, userName, fileSha1 string, permRead, permWrite, permDelete, permShare bool, expireTime *time.Time) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(`
        INSERT INTO tbl_permission 
        (role_name, user_name, file_sha1, perm_read, perm_write, perm_delete, perm_share, expire_time) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
        perm_read = VALUES(perm_read),
        perm_write = VALUES(perm_write),
        perm_delete = VALUES(perm_delete),
        perm_share = VALUES(perm_share),
        expire_time = VALUES(expire_time)
    `)
	if err != nil {
		log.Println("Failed to prepare statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(roleName, userName, fileSha1, permRead, permWrite, permDelete, permShare, expireTime)
	if err != nil {
		log.Println("Failed to execute statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	res.Suc = true
	return
}

// RevokePermission 撤销权限
func RevokePermission(roleName, userName, fileSha1 string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("DELETE FROM tbl_permission WHERE role_name = ? AND user_name = ? AND file_sha1 = ?")
	if err != nil {
		log.Println("Failed to prepare statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(roleName, userName, fileSha1)
	if err != nil {
		log.Println("Failed to execute statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	res.Suc = true
	return
}

// CheckPermission 检查权限
func CheckPermission(userName, fileSha1 string) (res ExecResult) {
	query := `
        SELECT p.perm_read, p.perm_write, p.perm_delete, p.perm_share, p.expire_time
        FROM tbl_permission p
        LEFT JOIN tbl_user_role ur ON p.role_name = ur.role_name
        WHERE (p.user_name = ? OR ur.user_name = ?) AND p.file_sha1 = ?
        AND (p.expire_time IS NULL OR p.expire_time > NOW())
    `

	stmt, err := mydb.DBConn().Prepare(query)
	if err != nil {
		log.Println("Failed to prepare statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	var permRead, permWrite, permDelete, permShare bool
	var expireTime sql.NullTime

	err = stmt.QueryRow(userName, userName, fileSha1).Scan(&permRead, &permWrite, &permDelete, &permShare, &expireTime)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Suc = true
			res.Data = map[string]bool{
				"read":   false,
				"write":  false,
				"delete": false,
				"share":  false,
			}
		} else {
			log.Println("Failed to execute query, err:", err.Error())
			res.Suc = false
			res.Msg = err.Error()
		}
		return
	}

	res.Suc = true
	res.Data = map[string]bool{
		"read":   permRead,
		"write":  permWrite,
		"delete": permDelete,
		"share":  permShare,
	}
	return
}

// ListUserPermissions 列出用户的所有权限
func ListUserPermissions(userName string) (res ExecResult) {
	query := `
        SELECT p.file_sha1, f.file_name, p.perm_read, p.perm_write, p.perm_delete, p.perm_share, p.expire_time
        FROM tbl_permission p
        LEFT JOIN tbl_user_role ur ON p.role_name = ur.role_name
        LEFT JOIN tbl_file f ON p.file_sha1 = f.file_sha1
        WHERE p.user_name = ? OR ur.user_name = ?
        AND (p.expire_time IS NULL OR p.expire_time > NOW())
    `

	stmt, err := mydb.DBConn().Prepare(query)
	if err != nil {
		log.Println("Failed to prepare statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(userName, userName)
	if err != nil {
		log.Println("Failed to execute query, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer rows.Close()

	var permissions []map[string]interface{}
	for rows.Next() {
		var fileSha1, fileName string
		var permRead, permWrite, permDelete, permShare bool
		var expireTime sql.NullTime

		err := rows.Scan(&fileSha1, &fileName, &permRead, &permWrite, &permDelete, &permShare, &expireTime)
		if err != nil {
			log.Println("Failed to scan row, err:", err.Error())
			continue
		}

		perm := map[string]interface{}{
			"file_sha1":   fileSha1,
			"file_name":   fileName,
			"read":        permRead,
			"write":       permWrite,
			"delete":      permDelete,
			"share":       permShare,
			"expire_time": nil,
		}
		if expireTime.Valid {
			perm["expire_time"] = expireTime.Time
		}
		permissions = append(permissions, perm)
	}

	res.Suc = true
	res.Data = permissions
	return
}
