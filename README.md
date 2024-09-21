<div className="max-w-4xl mx-auto p-6 bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-200">
      <h1 className="text-4xl font-bold mb-4 text-center">🌐 Distributed File System on Cloud</h1>
      
      <p className="text-center mb-6">
        <img src="https://img.shields.io/badge/status-active-success.svg" alt="Status" className="inline-block mr-2" />
        <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License" className="inline-block mr-2" />
        <img src="https://img.shields.io/badge/version-1.0.0-blue.svg" alt="Version" className="inline-block" />
      </p>

      <h2 className="text-2xl font-semibold mb-4">📋 Table of Contents</h2>
      <ul className="list-disc pl-6 mb-6">
        <li><a href="#overview" className="text-blue-500 hover:underline">Overview</a></li>
        <li><a href="#tech-stack" className="text-blue-500 hover:underline">Tech Stack</a></li>
        <li><a href="#system-architecture" className="text-blue-500 hover:underline">System Architecture</a></li>
        <li><a href="#installation" className="text-blue-500 hover:underline">Installation</a></li>
        <li><a href="#usage" className="text-blue-500 hover:underline">Usage</a></li>
        <li><a href="#contributing" className="text-blue-500 hover:underline">Contributing</a></li>
        <li><a href="#license" className="text-blue-500 hover:underline">License</a></li>
        <li><a href="#contact" className="text-blue-500 hover:underline">Contact</a></li>
      </ul>

      <h2 id="overview" className="text-2xl font-semibold mb-4">🔍 Overview</h2>
      <p className="mb-6">
        This project is a cloud-based distributed file system designed for scalability, reliability, and performance. 
        The system is built using cutting-edge technologies to ensure robust data management and seamless file storage 
        and retrieval processes across distributed environments.
      </p>

      <h2 id="tech-stack" className="text-2xl font-semibold mb-4">🛠️ Tech Stack</h2>
      <div className="grid grid-cols-2 gap-4 mb-6">
        <div>
          <h3 className="text-xl font-semibold mb-2">🖥️ Frontend</h3>
          <ul className="list-disc pl-6">
            <li>React.js</li>
            <li>CSS, Bootstrap</li>
          </ul>
        </div>
        <div>
          <h3 className="text-xl font-semibold mb-2">⚙️ Backend</h3>
          <ul className="list-disc pl-6">
            <li>Gin (Go web framework)</li>
            <li>Swagger (API documentation)</li>
          </ul>
        </div>
        <div>
          <h3 className="text-xl font-semibold mb-2">🔧 Microservices</h3>
          <ul className="list-disc pl-6">
            <li>go-micro</li>
            <li>gRPC</li>
          </ul>
        </div>
        <div>
          <h3 className="text-xl font-semibold mb-2">💾 Storage</h3>
          <ul className="list-disc pl-6">
            <li>Ceph</li>
            <li>AWS S3</li>
          </ul>
        </div>
        <div>
          <h3 className="text-xl font-semibold mb-2">🗄️ Databases</h3>
          <ul className="list-disc pl-6">
            <li>MySQL 5.7</li>
            <li>Redis 6.2.7</li>
          </ul>
        </div>
        <div>
          <h3 className="text-xl font-semibold mb-2">🐳 Containerization</h3>
          <ul className="list-disc pl-6">
            <li>Docker</li>
            <li>Kubernetes</li>
          </ul>
        </div>
      </div>

      <h2 id="system-architecture" className="text-2xl font-semibold mb-4">🏗️ System Architecture</h2>
      <img src="/usr/local/Distributed_system/cloud_distributed_storage/microservice_interact_archi.png" alt="System Architecture" className="w-full mb-4" />
      <p className="mb-6">
        The system is designed with a microservices architecture, where each component is loosely coupled, 
        enabling independent scaling and development. The architecture leverages containerization and 
        orchestration to manage resources efficiently and ensure seamless integration between services.
      </p>

      <h2 id="installation" className="text-2xl font-semibold mb-4">🚀 Installation</h2>
      <div className="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg mb-6">
        <p className="font-semibold mb-2">1. Clone the repository:</p>
        <pre className="bg-gray-200 dark:bg-gray-600 p-2 rounded">
          <code>
            git clone https://github.com/your-repo/distributed-file-system.git
            cd distributed-file-system
          </code>
        </pre>
        
        <p className="font-semibold mt-4 mb-2">2. Set up Docker containers:</p>
        <pre className="bg-gray-200 dark:bg-gray-600 p-2 rounded">
          <code>
            docker-compose -f docker-compose-mysql.yml up -d
          </code>
        </pre>
        
        <p className="font-semibold mt-4 mb-2">3. Install and run Redis:</p>
        <pre className="bg-gray-200 dark:bg-gray-600 p-2 rounded">
          <code>
            wget http://download.redis.io/releases/redis-6.2.7.tar.gz
            tar xzf redis-6.2.7.tar.gz
            cd redis-6.2.7
            make
            src/redis-server
          </code>
        </pre>
        
        <p className="font-semibold mt-4 mb-2">4. Configure Ceph storage:</p>
        <p>Follow the official Ceph documentation to set up the distributed storage system.</p>
        
        <p className="font-semibold mt-4 mb-2">5. Start the application:</p>
        <pre className="bg-gray-200 dark:bg-gray-600 p-2 rounded">
          <code>
            go run main.go
          </code>
        </pre>
      </div>

      <h2 id="usage" className="text-2xl font-semibold mb-4">📘 Usage</h2>
      <p className="mb-6">
        The system can be accessed via a web interface or API, where users can upload, manage, and retrieve files. 
        The file management interface provides features such as file versioning, access controls, and real-time status updates.
      </p>

      <h2 id="contributing" className="text-2xl font-semibold mb-4">🤝 Contributing</h2>
      <p className="mb-6">
        Contributions to this project are welcome. Please follow the guidelines in the <code>CONTRIBUTING.md</code> file 
        to submit issues or pull requests.
      </p>

      <h2 id="license" className="text-2xl font-semibold mb-4">📄 License</h2>
      <p className="mb-6">
        This project is licensed under the MIT License. See the <code>LICENSE</code> file for more details.
      </p>

      <h2 id="contact" className="text-2xl font-semibold mb-4">📧 Contact</h2>
      <p className="mb-6">
        For any inquiries, please contact: <a href="mailto:songguocheng348@gmail.com" className="text-blue-500 hover:underline">songguocheng348@gmail.com</a>
      </p>
    </div>
