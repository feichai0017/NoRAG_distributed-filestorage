#!/bin/bash

# 设置工作目录
WORK_DIR="/usr/local/Distributed_system/cloud_distributed_storage/Backend/"
cd "$WORK_DIR" || { echo "无法切换到工作目录"; exit 1; }

# 创建日志目录
LOG_DIR="/tmp/log/filestore-server"
mkdir -p "$LOG_DIR" || { echo "无法创建日志目录"; exit 1; }

# 停止进程
stop_process() {
    local service=$1
    echo "正在停止 $service 服务..."
    local pid=$(pgrep -f "service/$service")
    if [[ -n $pid ]]; then
        kill -15 $pid  # 首先尝试优雅关闭
        for i in {1..10}; do  # 等待最多10秒
            if ! pgrep -f "service/$service" > /dev/null; then
                echo -e "\033[32m已关闭\033[0m $service"
                return 0
            fi
            sleep 1
        done
        # 如果进程仍在运行，强制终止
        kill -9 $pid
        if ! pgrep -f "service/$service" > /dev/null; then
            echo -e "\033[33m已强制关闭\033[0m $service"
            return 0
        else
            echo -e "\033[31m关闭失败\033[0m $service"
            return 1
        fi
    else
        echo -e "\033[33m服务未运行\033[0m $service"
        return 0
    fi
}

# 清理 Consul 注册
clean_consul() {
    if command -v consul &> /dev/null; then
        echo "正在清理 Consul 注册..."
        for service in "${services[@]}"; do
            consul services deregister -id "go.micro.service.$service" &> /dev/null
            echo "已从 Consul 注销 $service"
        done
    else
        echo "Consul CLI 不可用，跳过 Consul 清理"
    fi
}

# 注销 Consul 服务
deregister_service() {
    local service=$1
    echo "正在从 Consul 注销服务: $service"

    # 获取匹配的服务 ID
    local service_ids=$(curl -s http://localhost:8500/v1/agent/services | jq -r 'to_entries[] | select(.value.Service | contains("'$service'")) | .key')

    if [ -z "$service_ids" ]; then
        echo "未找到匹配的服务 ID"
    else
        for id in $service_ids; do
            echo "注销服务 ID: $id"
            local response=$(curl -s -w "%{http_code}" -X PUT "http://localhost:8500/v1/agent/service/deregister/$id")
            if [[ $response == "200" ]]; then
                echo "成功注销服务 ID: $id"
            else
                echo "注销服务 ID 失败: $id, HTTP 状态码: $response"
            fi
        done
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
failed_services=()
for service in "${services[@]}"; do
    if ! stop_process "$service"; then
        failed_services+=("$service")
    fi
    deregister_service "$service"
done

# 清理 Consul 注册
clean_consul

# 报告结果
if [ ${#failed_services[@]} -eq 0 ]; then
    echo "所有微服务已成功停止."
else
    echo "以下服务停止失败: ${failed_services[*]}"
    exit 1
fi