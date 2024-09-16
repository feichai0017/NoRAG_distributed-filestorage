#!/bin/bash

ROOT_DIR=/usr/local/Distributed_system/cloud_distributed_storage/Backend/
#ROOT_DIR=/data/imooc/src/filestore-server
services="
dbproxy
upload
download
transfer
account
apigw
"

echo -e "\033[32m停止运行微服务容器... \033[0m"
for service in $services
do
    app=`sudo docker ps -a | grep "hub.fileserver.com/filestore/${service}" | awk '{print $1}'`
    if [[ $app != "" ]];then
        echo $app | xargs sudo docker stop
    fi
done