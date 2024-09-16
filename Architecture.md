# 综合云存储系统架构

## 1. 系统概览

本系统是一个分布式云存储平台，采用微服务架构，使用现代化的技术栈和工具链。系统的主要特点包括高可用性、可扩展性、安全性和易维护性。

## 2. 技术栈

### 2.1 前端
- React: 用于构建用户界面
- Redux: 状态管理
- React Router: 路由管理
- Axios: HTTP 客户端
- Ant Design: UI 组件库

### 2.2 后端
- Go: 主要编程语言
- gRPC: 服务间通信
- Gin: HTTP API 框架

### 2.3 数据存储
- TiDB: 分布式 SQL 数据库
- Ceph: 对象存储系统
- Redis: 缓存系统

### 2.4 消息队列
- RabbitMQ: 异步消息处理

### 2.5 API 网关
- Traefik: 动态路由和负载均衡

### 2.6 服务发现与配置
- Consul: 服务发现和配置管理

### 2.7 监控与日志
- Prometheus: 监控系统
- Grafana: 可视化仪表盘
- ELK Stack: 日志收集和分析

### 2.8 CI/CD
- GitLab: 版本控制和 CI/CD 管道
- Jenkins: 自动化构建和部署
- Docker: 容器化
- Kubernetes: 容器编排

## 3. 微服务架构

系统被拆分为以下微服务：

1. API 网关服务 (apigw)
2. 账户服务 (account)
3. 上传服务 (upload)
4. 下载服务 (download)
5. 传输服务 (transfer)
6. 文件元数据服务 (metadata)
7. 存储服务 (storage)
8. 共享服务 (share)
9. 通知服务 (notification)

## 4. 系统架构图

```
![Architect](https://github.com/user-attachments/assets/bb37f8c0-afee-4512-af42-d033eb94cd2a)

```

## 5. CI/CD 流程

1. 开发人员将代码提交到 GitLab
2. GitLab 触发 CI 管道，运行单元测试和代码质量检查
3. 通过测试后，Jenkins 从 GitLab 拉取代码
4. Jenkins 构建 Docker 镜像并推送到私有 Docker 仓库
5. Jenkins 更新 Kubernetes 部署配置
6. Kubernetes 拉取新镜像并更新相应的服务

## 6. 微服务详细说明

### 6.1 API 网关服务 (apigw)
- 职责：请求路由、负载均衡、认证授权、限流
- 技术：Traefik、Go
- API：/api/v1/*

### 6.2 账户服务 (account)
- 职责：用户注册、登录、信息管理
- 技术：Go、gRPC、TiDB
- API：
  - POST /api/v1/account/register
  - POST /api/v1/account/login
  - GET /api/v1/account/profile
  - PUT /api/v1/account/profile

### 6.3 上传服务 (upload)
- 职责：处理文件上传请求，文件分片
- 技术：Go、gRPC、Ceph
- API：
  - POST /api/v1/upload/init
  - POST /api/v1/upload/chunk
  - POST /api/v1/upload/complete

### 6.4 下载服务 (download)
- 职责：处理文件下载请求
- 技术：Go、gRPC、Ceph
- API：
  - GET /api/v1/download/:fileId

### 6.5 传输服务 (transfer)
- 职责：管理文件传输任务，断点续传
- 技术：Go、gRPC、Redis
- API：
  - POST /api/v1/transfer/pause/:taskId
  - POST /api/v1/transfer/resume/:taskId
  - GET /api/v1/transfer/status/:taskId

### 6.6 文件元数据服务 (metadata)
- 职责：管理文件元数据信息
- 技术：Go、gRPC、TiDB
- API：
  - GET /api/v1/metadata/:fileId
  - PUT /api/v1/metadata/:fileId

### 6.7 存储服务 (storage)
- 职责：与对象存储系统交互
- 技术：Go、gRPC、Ceph
- 内部服务，不对外暴露 API

### 6.8 共享服务 (share)
- 职责：管理文件共享和权限
- 技术：Go、gRPC、TiDB
- API：
  - POST /api/v1/share/create
  - GET /api/v1/share/:shareId
  - DELETE /api/v1/share/:shareId

### 6.9 通知服务 (notification)
- 职责：处理系统通知和消息
- 技术：Go、gRPC、RabbitMQ
- API：
  - GET /api/v1/notifications
  - POST /api/v1/notifications/read/:notificationId

## 7. 数据流

1. 文件上传流程：
   客户端 -> API网关 -> 上传服务 -> 存储服务 -> Ceph
                    -> 元数据服务 -> TiDB

2. 文件下载流程：
   客户端 -> API网关 -> 下载服务 -> 存储服务 -> Ceph
                    -> 元数据服务 -> TiDB

3. 文件共享流程：
   客户端 -> API网关 -> 共享服务 -> 元数据服务 -> TiDB
                    -> 通知服务 -> RabbitMQ

## 8. 安全考虑

1. 全站 HTTPS
2. JWT 进行身份验证
3. 文件加密存储
4. 细粒度的访问控制
5. 定期安全审计
6. DDoS 防护

## 9. 性能优化

1. 使用 CDN 加速文件下载
2. Redis 缓存热点数据
3. 文件分片上传和并发下载
4. 数据库读写分离
5. 使用消息队列处理异步任务
6. Kubernetes 自动扩缩容

## 10. 可靠性和可用性

1. 服务多副本部署
2. 数据多副本存储
3. 跨区域部署
4. 定期备份和恢复演练
5. 熔断和限流保护
6. 全面的监控和告警机制

这个综合架构设计涵盖了您所需的所有元素，包括前端技术、CI/CD 流程和细化的微服务结构。它提供了一个可扩展、高性能和可靠的分布式云存储系统框架。
