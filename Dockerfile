FROM centos:9
LABEL authors="root"

# Stage 1: MySQL Master
FROM mysql:5.7 AS master

# 设置环境变量
ENV MYSQL_ROOT_PASSWORD=119742
ENV MYSQL_DATABASE=filestore
ENV MYSQL_USER=master
ENV MYSQL_PASSWORD=119742

# 启动 MySQL 服务
CMD ["mysqld"]

# Stage 2: MySQL Slave
FROM mysql:5.7 AS slave

# 将 master 的 binlog 目录挂载到 slave 上
VOLUME /var/lib/mysql/master-binlog

# 设置环境变量
ENV MYSQL_ROOT_PASSWORD=119742
ENV MYSQL_DATABASE=filestore
ENV MYSQL_USER=slave
ENV MYSQL_PASSWORD=119742

# 启动 MySQL 服务
CMD ["mysqld"]



ENTRYPOINT ["top", "-b"]