#!/bin/bash

stop_process() {
    sleep 1
    pid=`ps aux | grep -v grep | grep "service/bin" | grep $1 | awk '{print $2}'`
    if [[ $pid != '' ]]; then
	ps aux | grep -v grep | grep "service/bin" | grep $1 | awk '{print $2}' | xargs kill
        echo -e "\033[32m已关闭: \033[0m" "$1"
        return 1
    else
        echo -e "\033[31m并未启动: \033[0m" "$1"
        return 0
    fi
}


services="
apigw
account
transfer
download
upload
dbproxy
"

# 关闭service
for sname in $services
do
    stop_process $sname
done

echo "执行完毕."