#!/bin/bash

# Consul API 地址
CONSUL_API="http://localhost:8500/v1"

# 函数：从 Consul 获取服务信息
get_services() {
    curl -s "${CONSUL_API}/agent/services" | jq -r 'to_entries[] | .value | "\(.Service)|\(.ID)|\(.Port)"'
}

# 函数：停止指定端口的进程
stop_process() {
    local port=$1
    local pid=$(lsof -ti:$port)
    if [ -n "$pid" ]; then
        echo "正在停止端口 $port 的进程 (PID: $pid)..."
        kill -15 $pid
        sleep 2
        if kill -0 $pid 2>/dev/null; then
            echo "进程未响应，强制终止..."
            kill -9 $pid
        fi
        echo "进程已停止"
    else
        echo "端口 $port 没有运行的进程"
    fi
}

# 函数：从 Consul 注销服务
deregister_service() {
    local service_id=$1
    echo "从 Consul 注销服务: $service_id"
    curl -X PUT "${CONSUL_API}/agent/service/deregister/${service_id}"
}

# 主程序
echo "正在从 Consul 获取服务信息..."
services=$(get_services)

if [ -z "$services" ]; then
    echo "没有找到注册的服务"
    exit 0
fi

echo "开始停止服务进程并注销服务..."
echo "$services" | while IFS='|' read -r service_name service_id port; do
    echo "处理服务: $service_name (ID: $service_id, Port: $port)"
    stop_process $port
    deregister_service $service_id
    echo "------------------------"
done

echo "所有操作完成"

# 最后检查是否还有相关进程在运行
echo "检查是否还有相关进程在运行..."
remaining_services=$(get_services)
if [ -n "$remaining_services" ]; then
    echo "警告: 以下服务仍在 Consul 中注册:"
    echo "$remaining_services"
else
    echo "所有服务已从 Consul 中注销"
fi

echo "检查是否还有相关端口被占用..."
echo "$services" | while IFS='|' read -r service_name service_id port; do
    pid=$(lsof -ti:$port)
    if [ -n "$pid" ]; then
        echo "警告: 端口 $port 仍被进程 (PID: $pid) 占用"
    fi
done

echo "脚本执行完毕"