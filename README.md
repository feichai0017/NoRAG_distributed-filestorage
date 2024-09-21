# Distributed File System on Cloud

![Status](https://img.shields.io/badge/status-active-success.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)

A cloud-based distributed file system designed for scalability, reliability, and performance.

---

## Table of Contents
- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [System Architecture](#system-architecture)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

---

## Overview

This project is a cloud-based distributed file system built using cutting-edge technologies to ensure robust data management and seamless file storage and retrieval processes across distributed environments.

---

## Tech Stack

### Operating System
- Linux (CentOS 9)

### Frontend
- **Framework:** React.js
- **Styling:** CSS, Bootstrap

### Backend
- **Framework:** Gin (Go web framework)
- **API Documentation:** Swagger

### Programming Languages
- Go
- JavaScript

### Microservices
- **File Management:** Handles file upload, download, and management operations
- **Framework:** go-micro
- **Communication:** gRPC

### Distributed Storage
- **Ceph:** Unified, distributed storage system
- **AWS S3:** Integrated for additional cloud storage capabilities

### Databases
- **MySQL 5.7:** Deployed using Docker with master-slave replication
- **Redis 6.2.7:** In-memory key-value store for caching

### Message Queue
- RabbitMQ

### Containerization and Orchestration
- Docker
- Kubernetes

---

## System Architecture

![System Architecture](/Architect.png)

The system uses a microservices architecture with loosely coupled components, enabling independent scaling and development. Containerization and orchestration manage resources efficiently and ensure seamless integration between services.

### MySQL Setup
- **Version:** MySQL 5.7
- **Deployment:** Docker container
- **Access:**
  - Master Instance: `localhost:3301`
  - Slave Instance: `localhost:3302`
  - Container Internal Port: `3306`

### Redis Setup
- **Version:** Redis 6.2.7
- **Deployment:** Docker container
- **Access:** `localhost:6379`

### Ceph Storage
- **Ceph Monitor:** `172.20.0.10/16`
- **Ceph OSD:** `172.20.0.11, 172.20.0.12, 172.20.0.13`
- **Ceph MGR:** `172.20.0.14:7000`
- **Ceph RGW:** `172.20.0.15:7480`

### AWS S3 Integration
- Ceph RGW configured to provide S3-compatible API
- Interact with Ceph using AWS CLI or SDKs

### RabbitMQ Setup
- **Version:** RabbitMQ 3.9.7
- **Deployment:** Docker container
- **Ports:** 5672, 15672, 25672

---

## Installation

1. **Clone the repository:**
 ```bash
 git clone https://github.com/feichai0017/distributed-file-system.git
 cd distributed-file-system
 ```

2. **Set up Docker containers:**

  ```shellscript
  docker-compose -f docker-compose-mysql.yml up -d
  ```


3. **Install and run Redis:**

  ```shellscript
  wget http://download.redis.io/releases/redis-6.2.7.tar.gz
  tar xzf redis-6.2.7.tar.gz
  cd redis-6.2.7
  make
  src/redis-server
  ```


4. **Configure Ceph storage:**
  Follow the official Ceph documentation to set up the distributed storage system.

5.**Start the application:**
To start all services, run the following command:

```shellscript
./service/start-all.sh
```

This script will initialize and start all necessary services for the distributed file system.




---

## Usage

The system can be accessed via a web interface or API. Users can upload, manage, and retrieve files with features such as file versioning, access controls, and real-time status updates.

---

## Contributing

Contributions are welcome. Please follow the guidelines in the `CONTRIBUTING.md` file to submit issues or pull requests.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

---

## Contact

For inquiries, please contact: [songguocheng348@gmail.com](mailto:songguocheng348@gmail.com)

  ```plaintext
  
  ```
