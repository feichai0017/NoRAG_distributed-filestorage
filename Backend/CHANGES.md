# 系统更新日志

## 数据库结构变更

1. 修改 `tbl_file` 表：
    - 添加 `owner_id` 字段
    - 原因：支持文件所有权管理，实现多用户系统

2. 修改 `tbl_user` 表：
    - 添加 `email`, `phone`, `email_validated`, `phone_validated`, `profile` 字段
    - 原因：增强用户信息管理，支持邮箱和手机号验证

3. 新增 `tbl_role` 表：
    - 包含 `id`, `role_name`, `description`, `create_at`, `update_at` 字段
    - 原因：实现基于角色的访问控制（RBAC）

4. 新增 `tbl_user_role` 表：
    - 包含 `id`, `user_id`, `role_id`, `create_at` 字段
    - 原因：建立用户和角色之间的多对多关系

5. 新增 `tbl_permission` 表：
    - 包含 `id`, `role_id`, `user_id`, `file_id`, `perm_read`, `perm_write`, `perm_delete`, `perm_share`, `expire_time`, `create_at`, `update_at` 字段
    - 原因：实现细粒度的文件访问权限控制

6. 修改 `tbl_user_file` 表：
    - 使用 `user_id` 和 `file_id` 替代原有的 `username` 和 `file_sha1`
    - 原因：提高查询效率，更好地与其他表关联

## ORM 更新

1. 更新 `file.go`：
    - 修改 `OnFileUploadFinished` 函数，添加 `owner_id` 参数
    - 更新 `GetFileMeta` 和 `GetFileMetaList` 函数以包含新字段
    - 原因：适应新的文件表结构，支持文件所有权

2. 更新 `user.go`：
    - 修改 `UserSignup` 函数，添加 email 和 phone 参数
    - 更新 `GetUserInfo` 函数以返回更多用户信息
    - 原因：支持更丰富的用户信息管理

3. 更新 `user_file.go`：
    - 修改所有函数使用 `user_id` 和 `file_id` 而不是 `username` 和 `file_sha1`
    - 原因：适应新的用户文件表结构

4. 新增 `role.go`：
    - 添加 `CreateRole`, `GetRoleInfo`, `UpdateRole`, `DeleteRole`, `ListRoles` 等函数
    - 原因：支持角色管理功能

5. 新增 `permission.go`：
    - 添加 `GrantPermission`, `RevokePermission`, `CheckPermission`, `ListUserPermissions` 等函数
    - 原因：实现文件访问权限控制

## DBProxy 客户端更新

1. 更新 `client.go`：
    - 修改 `FileMeta` 结构体，添加 `OwnerID` 字段
    - 更新 `TableFileToFileMeta` 函数以适应新的 `TableFile` 结构
    - 修改 `UserSignup` 函数，添加 email 和 phone 参数
    - 更新 `QueryUserFileMeta` 和 `QueryUserFileMetas` 函数，使用 `userID` 替代 `username`
    - 添加新的角色和权限相关函数
    - 原因：适应新的 ORM 结构和 RBAC 模型需求

## 总体更新原因

1. 实现基于角色的访问控制（RBAC）：
    - 增加了角色和权限管理，提高了系统的安全性和灵活性
    - 允许更精细的文件访问控制

2. 增强用户管理：
    - 添加了更多用户信息字段，支持邮箱和手机号验证
    - 改善了用户体验和账户安全

3. 优化数据结构：
    - 使用整数ID替代字符串作为主键，提高查询效率
    - 建立了更清晰的表间关系，方便数据管理和查询

4. 改进文件管理：
    - 引入文件所有权概念，支持多用户系统
    - 实现更细粒度的文件权限控制

5. 提高系统可扩展性：
    - 新的结构设计使得系统更容易扩展新功能
    - RBAC 模型为未来添加更复杂的权限规则奠定了基础

这些更改collectively提高了系统的功能性、安全性、可扩展性和性能，为future的发展打下了坚实的基础。