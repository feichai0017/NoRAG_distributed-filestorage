

# Distributed File System on Cloud

## Overview

This project is a cloud-based distributed file system designed for scalability, reliability, and performance. The system is built using cutting-edge technologies to ensure robust data management and seamless file storage and retrieval processes across distributed environments.

## Tech Stack

- **Operating System:** Linux (CentOS 9)
- **Frontend:**
    - **Framework:** React.js
    - **Styling:** CSS, Bootstrap
- **Backend:**
    - **Framework:** Gin - A lightweight web framework for building high-performance APIs in Go.
    - **API Documentation:** Swagger - A tool for documenting and testing APIs.
- **Programming Language:** Go, JavaScript
- **Micro Services:**
    - **File Management:** Handles file upload, download, and management operations.
    - **Framework:** go-micro - A pluggable microservices framework for Go.
    - **Communication:** gRPC - A high-performance, open-source RPC framework.
- **Distributed Storage:**
    - **Ceph:** A unified, distributed storage system designed for excellent performance, reliability, and scalability.
    - **AWS S3:** Integrated for additional cloud storage capabilities and global accessibility.
- **Database:**
    - **MySQL** - A widely-used relational database management system. Deployed using Docker with the official MySQL 5.7 image. The configuration includes master-slave replication for data redundancy and failover capabilities.
    - **Redis** - An in-memory key-value store, used for caching and real-time analytics. Deployed using Docker with the official Redis 6.2.7 image.
- **Message Queue:** RabbitMQ - A robust message broker that facilitates asynchronous communication between distributed components.
- **Containerization and Orchestration:**
    - **Docker:** Containers are used to package the application components and their dependencies, ensuring consistency across different environments.
    - **Kubernetes:** Manages containerized applications in a clustered environment, ensuring high availability, scalability, and fault tolerance.

## System Architecture
![](/usr/local/Distributed_system/cloud_distributed_storage/microservice_interact_archi.png)


The system is designed with a microservices architecture, where each component is loosely coupled, enabling independent scaling and development. The architecture leverages containerization and orchestration to manage resources efficiently and ensure seamless integration between services.

### MySQL Setup

- **Version:** MySQL 5.7
- **Deployment:** Docker container using the official MySQL image.
- **Configuration:**
    - Default configuration with master-slave replication enabled.
    - Configuration files are mounted to the container using Docker volumes, allowing for easy updates and persistence.
- **Access:**
    - **Master Instance:** `localhost:3301`
    - **Slave Instance:** `localhost:3302`
    - **Container Internal Port:** `3306`

### Redis Setup

- **Version:** Redis 6.2.7
- **Deployment:** Docker container using the official Redis image.
- **Configuration:** Default settings.
- **Access:** `localhost:6379`

### Ceph Storage

Ceph is used to handle distributed file storage, offering high scalability and reliability. It provides seamless integration with other components in the system, ensuring efficient data storage and retrieval.
- **Access:** 
- **Ceph Monitor:** `172.20.0.10/16`
- **Ceph OSD:** `172.20.0.11, 172.20.0.12, 172.20.0.13`
- **Ceph MGR:** `172.20.0.14:7000`
- **Ceph RGW:** `172.20.0.15:7480`

### AWS S3 Integration
**Ceph Configuration:**
    <li>Ceph can be configured to use its RADOS Gateway (RGW) to provide an S3-compatible API.
    <li>Configure Ceph RGW to interact with AWS S3 using the S3 API, enabling data redundancy across your Ceph cluster and AWS.
**Usage:**
    <li>You can interact with Ceph just as you would with AWS S3, using tools like the AWS CLI or SDKs, by pointing them to your Ceph RGW endpoint.

### RabbitMQ Setup
**Version:** RabbitMQ 3.9.7
**Deployment:** Docker container using the official RabbitMQ image.
  <li>docker run -d --hostname rabbit-server --name rabbit -p 5672:5672 -p 15672:15672 -p 25672:25672 -v /data/rabbitmq:/var/lib/rabbitmq rabbitmq:management<li>
**Configuration:** Default settings.

## Installation

To set up the system, follow these steps:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-repo/distributed-file-system.git
   cd distributed-file-system
   ```

2. **Set up Docker containers:**
    - For MySQL:
      ```bash
      docker-compose -f docker-compose-mysql.yml up -d
      ```

3. **Redis 6.2.7 Installation and Running Instructions:**

   1. Download and Install Redis 6.2.7

       ```bash
         wget http://download.redis.io/releases/redis-6.2.7.tar.gz
         tar xzf redis-6.2.7.tar.gz
         cd redis-6.2.7
         make
       ```

   2. Start the Redis server:
       ```bash
       src/redis-server
       redis-server /etc/redis/redis.conf
       redis-cli
       ```

4**Configure Ceph storage:**
   Follow the official Ceph documentation to set up the distributed storage system.

5**Start the application:**
   ```bash
   go run main.go
   ```

## Usage

The system can be accessed via a web interface or API, where users can upload, manage, and retrieve files. The file management interface provides features such as file versioning, access controls, and real-time status updates.

## Contributing

Contributions to this project are welcome. Please follow the guidelines in the `CONTRIBUTING.md` file to submit issues or pull requests.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.

## Contact
````
songguocheng348@gmail.com
