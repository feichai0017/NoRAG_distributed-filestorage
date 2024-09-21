# Distributed File System on Cloud

![Status](https://img.shields.io/badge/status-active-success.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)

A cloud-based distributed file system designed for scalability, reliability, and performance, with advanced storage, retrieval capabilities, and service discovery.

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

This project is a cloud-based distributed file system built using cutting-edge technologies to ensure robust data management, seamless file storage and retrieval processes across distributed environments, advanced document processing, search capabilities, and efficient service discovery and management.

---

## Tech Stack

### Operating System
- x86 Linux 

### Frontend
- **Framework:** React.js
- **Styling:** CSS, Bootstrap

### Backend
- **Framework:** Gin (Go web framework)
- **API Documentation:** Swagger

### Programming Languages
- Go
- JavaScript
- Python

### Microservices
- **File Management:** Handles file upload, download, and management operations
- **Framework:** go-micro
- **Communication:** gRPC

### Service Discovery and Configuration
- **Consul:** Service registry, discovery, and distributed key-value store

### Distributed Storage
- **Block Storage:**
  - Ceph: For frequently modified documents
  - AWS EBS: Backup for block storage
- **Object Storage:**
  - MinIO: Primary object storage for immutable files (e.g., PDFs, images)
  - AWS S3: Backup for object storage

### Databases
- **MySQL 5.7:** Deployed using Docker with master-slave replication
- **Redis 6.2.7:** In-memory key-value store for caching

### Search and Analytics
- **Elasticsearch:** Full-text search and vector database for document indexing

### Document Processing
- **Python:** RAG (Retrieval-Augmented Generation) service and document scanning functionality

### Message Queue
- RabbitMQ

### Containerization and Orchestration
- Docker
- Kubernetes

---

## System Architecture

![System Architecture](/Architect.png)

The system uses a microservices architecture with loosely coupled components, enabling independent scaling and development. Containerization and orchestration manage resources efficiently and ensure seamless integration between services. Consul provides service discovery and configuration management across the distributed system.

### Consul Setup
- **Version:** Consul 1.11.0
- **Deployment:** Docker container
- **UI Access:** `localhost:8500`
- **DNS Interface:** `localhost:8600`

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

### Ceph Block Storage
- **Ceph Monitor:** `172.20.0.10/16`
- **Ceph OSD:** `172.20.0.11, 172.20.0.12, 172.20.0.13`
- **Ceph MGR:** `172.20.0.14:7000`
- **Ceph RBD:** For block storage of frequently modified documents

### MinIO Object Storage
- **Deployment:** Docker container
- **Access:** `localhost:9000`
- **Console:** `localhost:9001`

### AWS Integration
- **AWS EBS:** Backup for Ceph block storage
- **AWS S3:** Backup for MinIO object storage

### Elasticsearch Setup
- **Version:** Elasticsearch 7.14.0
- **Deployment:** Docker container
- **Access:** `localhost:9200`

### Python Services
- **RAG Service:** Retrieval-Augmented Generation for advanced document processing
- **Document Scanning:** OCR and document analysis capabilities

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
   ```bash
   docker-compose up -d
   ```

3. **Start Consul:**
   ```bash
   docker run -d --name=consul -p 8500:8500 -p 8600:8600/udp consul:1.11.0 agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
   ```

4. **Configure Ceph block storage:**
   Follow the official Ceph documentation to set up the block storage system.

5. **Set up MinIO object storage:**
   ```bash
   docker run -p 9000:9000 -p 9001:9001 minio/minio server /data --console-address ":9001"
   ```

6. **Configure AWS services:**
   Set up AWS EBS and S3 for backup purposes using AWS CLI or SDKs.

7. **Start Elasticsearch:**
   ```bash
   docker run -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.14.0
   ```

8. **Set up Python services:**
   ```bash
   pip install -r requirements.txt
   python rag_service.py
   python document_scanner.py
   ```

9. **Register services with Consul:**
   Update your service configurations to register with Consul upon startup.

10. **Start the application:**
    To start all services, run the following command:
    ```bash
    ./service/start-all.sh
    ```

---

## Usage

The system can be accessed via a web interface or API. Users can upload, manage, and retrieve files with features such as file versioning, access controls, and real-time status updates. The system provides advanced search capabilities, document processing, and intelligent retrieval using RAG technology. Services are dynamically discovered and managed through Consul, ensuring high availability and scalability.

---

## Contributing

Contributions are welcome. Please follow the guidelines in the `CONTRIBUTING.md` file to submit issues or pull requests.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

---

## Contact

For inquiries, please contact: [songguocheng348@gmail.com](mailto:songguocheng348@gmail.com)
