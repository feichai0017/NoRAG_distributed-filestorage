#!/bin/bash

# 设置工作目录
WORK_DIR="/usr/local/Distributed_system/cloud_distributed_storage/Backend/"
cd "$WORK_DIR" || exit

# 创建日志目录
LOG_DIR="/tmp/log/filestore-server"
mkdir -p "$LOG_DIR"

# 检查进程是否运行
check_process() {
    sleep 5  # 增加等待时间，确保服务有足够时间启动
    if pgrep -f "service/bin/$1" > /dev/null; then
        echo -e "\033[32m已启动\033[0m $1"
        return 0
    else
        echo -e "\033[31m启动失败\033[0m $1"
        return 1
    fi
}

# 编译service可执行文件
build_service() {
    echo "正在编译 $1 服务..."
    go build -o service/bin/$1 service/$1/main.go
    resbin=$(ls service/bin/ | grep $1)
    echo -e "\033[32m编译完成:\033[0m service/bin/$resbin"
}

# 启动服务
start_service() {
    echo "正在启动 $1 服务..."
    nohup ./service/bin/$1 --registry=consul >> "$LOG_DIR/$1.log" 2>&1 &
    if ! check_process "$1"; then
        echo "服务 $1 启动失败，退出脚本。"
        exit 1
    fi
    sleep 5  # 在启动下一个服务之前稍作等待
}

# 服务列表
services=(
    "apigw"
    "dbproxy"
    "upload"
    "download"
    "transfer"
    "account"
)

# 确保 bin 目录存在
mkdir -p service/bin/ && rm -f service/bin/*

# 编译所有服务
for service in "${services[@]}"; do
    build_service "$service"
done

# 启动所有服务
for service in "${services[@]}"; do
    start_service "$service"
done

echo '所有微服务启动完毕.'