#!/bin/bash

# 设置工作目录
WORK_DIR="/Users/guochengsong/Documents/GitHub/cloud_distributed_file-system/Backend/"
cd "$WORK_DIR" || exit

# 创建日志目录
LOG_DIR="/tmp/log/filestore-server"
mkdir -p "$LOG_DIR"

# 检查进程是否运行
check_process() {
    sleep 2
    if pgrep -f "service/$1" > /dev/null; then
        echo -e "\033[32m已启动\033[0m $1"
        return 0
    else
        echo -e "\033[31m启动失败\033[0m $1"
        return 1
    fi
}

# 启动服务
start_service() {
    echo "正在启动 $1 服务..."
    go run "./Backend/service/$1/main.go" >> "$LOG_DIR/$1.log" 2>&1 &
    check_process "$1"
}

# 服务列表
services=(
    "dbproxy"
    "upload"
    "download"
    "transfer"
    "account"
    "apigw"
)

# 启动所有服务
for service in "${services[@]}"; do
    start_service "$service"
done

echo '所有微服务启动完毕.'