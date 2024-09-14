package orm

import (
	mydb "cloud_distributed_storage/Backend/database/mysql"
	"database/sql"
	"log"
	"time"
)

// CreateRole 创建新角色
func CreateRole(roleName string, description string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"INSERT INTO tbl_role (role_name, description, create_at) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(roleName, description, time.Now())
	if err != nil {
		log.Println("Failed to execute statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	res.Suc = true
	return
}

// GetRoleInfo 获取角色信息
func GetRoleInfo(roleID int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"SELECT id, role_name, description, create_at, update_at FROM tbl_role WHERE id = ?")
	if err != nil {
		log.Println("Failed to prepare statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	var role TableRole
	err = stmt.QueryRow(roleID).Scan(&role.ID, &role.RoleName, &role.Description, &role.CreateAt, &role.UpdateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Suc = false
			res.Msg = "Role not found"
		} else {
			log.Println("Failed to execute query, err:", err.Error())
			res.Suc = false
			res.Msg = err.Error()
		}
		return
	}

	res.Suc = true
	res.Data = role
	return
}

// UpdateRole 更新角色信息
func UpdateRole(roleID int64, roleName, description string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"UPDATE tbl_role SET role_name = ?, description = ?, update_at = ? WHERE id = ?")
	if err != nil {
		log.Println("Failed to prepare statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(roleName, description, time.Now(), roleID)
	if err != nil {
		log.Println("Failed to execute statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	res.Suc = true
	return
}

// DeleteRole 删除角色
func DeleteRole(roleID int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("DELETE FROM tbl_role WHERE id = ?")
	if err != nil {
		log.Println("Failed to prepare statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(roleID)
	if err != nil {
		log.Println("Failed to execute statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	res.Suc = true
	return
}

// ListRoles 列出所有角色
func ListRoles() (res ExecResult) {
	rows, err := mydb.DBConn().Query("SELECT id, role_name, description, create_at, update_at FROM tbl_role")
	if err != nil {
		log.Println("Failed to execute query, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer rows.Close()

	var roles []TableRole
	for rows.Next() {
		var role TableRole
		err := rows.Scan(&role.ID, &role.RoleName, &role.Description, &role.CreateAt, &role.UpdateAt)
		if err != nil {
			log.Println("Failed to scan row, err:", err.Error())
			continue
		}
		roles = append(roles, role)
	}

	res.Suc = true
	res.Data = roles
	return
}

// AssignRoleToUser 为用户分配角色
func AssignRoleToUser(userID, roleID int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"INSERT INTO tbl_user_role (user_id, role_id, create_at) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, roleID, time.Now())
	if err != nil {
		log.Println("Failed to execute statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	res.Suc = true
	return
}

// RemoveRoleFromUser 从用户移除角色
func RemoveRoleFromUser(userID, roleID int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("DELETE FROM tbl_user_role WHERE user_id = ? AND role_id = ?")
	if err != nil {
		log.Println("Failed to prepare statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, roleID)
	if err != nil {
		log.Println("Failed to execute statement, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	res.Suc = true
	return
}

// GetUserRoles 获取用户的所有角色
func GetUserRoles(userID int64) (res ExecResult) {
	rows, err := mydb.DBConn().Query(`
        SELECT r.id, r.role_name, r.description, r.create_at, r.update_at
        FROM tbl_role r
        INNER JOIN tbl_user_role ur ON r.id = ur.role_id
        WHERE ur.user_id = ?
    `, userID)
	if err != nil {
		log.Println("Failed to execute query, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer rows.Close()

	var roles []TableRole
	for rows.Next() {
		var role TableRole
		err := rows.Scan(&role.ID, &role.RoleName, &role.Description, &role.CreateAt, &role.UpdateAt)
		if err != nil {
			log.Println("Failed to scan row, err:", err.Error())
			continue
		}
		roles = append(roles, role)
	}

	res.Suc = true
	res.Data = roles
	return
}

// GetRoleUsers 获取拥有特定角色的所有用户
func GetRoleUsers(roleID int64) (res ExecResult) {
	rows, err := mydb.DBConn().Query(`
        SELECT u.id, u.user_name, u.email, u.phone, u.signup_at, u.last_active, u.status
        FROM tbl_user u
        INNER JOIN tbl_user_role ur ON u.id = ur.user_id
        WHERE ur.role_id = ?
    `, roleID)
	if err != nil {
		log.Println("Failed to execute query, err:", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer rows.Close()

	var users []TableUser
	for rows.Next() {
		var user TableUser
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.SignupAt, &user.LastActive, &user.Status)
		if err != nil {
			log.Println("Failed to scan row, err:", err.Error())
			continue
		}
		users = append(users, user)
	}

	res.Suc = true
	res.Data = users
	return
}
