#!/bin/bash

# 设置工作目录
WORK_DIR="/usr/local/Distributed_system/cloud_distributed_storage/Backend/"
cd "$WORK_DIR" || exit

# 创建日志目录
LOG_DIR="/tmp/log/filestore-server"
mkdir -p "$LOG_DIR"

# 停止进程
stop_process() {
    echo "正在停止 $1 服务..."
    pid=$(pgrep -f "service/$1")
    if [[ $pid != '' ]]; then
        kill $pid
        sleep 2  # 给进程一些时间来关闭
        if ! pgrep -f "service/$1" > /dev/null; then
            echo -e "\033[32m已关闭\033[0m $1"
            return 0
        else
            echo -e "\033[31m关闭失败\033[0m $1"
            return 1
        fi
    else
        echo -e "\033[33m服务未运行\033[0m $1"
        return 0
    fi
}

# 服务列表
services=(
    "apigw"
    "account"
    "transfer"
    "download"
    "upload"
    "dbproxy"
)

# 停止所有服务
for service in "${services[@]}"; do
    stop_process "$service"
done

echo "所有微服务已停止."