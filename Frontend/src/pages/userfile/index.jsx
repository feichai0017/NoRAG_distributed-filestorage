import './index.css'
import React, { useState } from 'react';

const UserFiles = () => {
    const [username, setUsername] = useState('');
    const [limit, setLimit] = useState(10);
    const [files, setFiles] = useState([]);

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        if (name === 'username') {
            setUsername(value);
        } else if (name === 'limit') {
            setLimit(value);
        }
    };

    const handleFormSubmit = (e) => {
        e.preventDefault();

        const params = new URLSearchParams({ username, limit });

        fetch('/api/meta/query', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: params,
        })
            .then((response) => response.json())
            .then((data) => {
                setFiles(data);
            })
            .catch((error) => console.error('Error:', error));
    };

    return (
        <div className="file-display">
            <div className="user-files">
            <h1>User Uploaded Files</h1>
            <form id="queryForm" onSubmit={handleFormSubmit}>
                <label htmlFor="username">Username:</label>
                <input
                    type="text"
                    id="username"
                    name="username"
                    value={username}
                    onChange={handleInputChange}
                    required
                />
                <label htmlFor="limit">Number of Files:</label>
                <input
                    type="number"
                    id="limit"
                    name="limit"
                    min="1"
                    max="100"
                    value={limit}
                    onChange={handleInputChange}
                    required
                />
                <button type="submit">Search</button>
            </form>

            <h2>Files:</h2>
            <table>
                <thead>
                <tr>
                    <th>File Name</th>
                    <th>File Hash</th>
                    <th>File Size (Bytes)</th>
                    <th>Upload Time</th>
                    <th>Action</th>
                    <th>Delete</th>
                </tr>
                </thead>
                <tbody id="fileList">
                {files.map((fileMeta, index) => (
                    <tr key={index}>
                        <td>{fileMeta.FileName}</td>
                        <td>{fileMeta.FileHash}</td>
                        <td>{fileMeta.FileSize}</td>
                        <td>{fileMeta.UploadAt}</td>
                        <td>
                            <a
                                href={`/file/download?filehash=${fileMeta.FileHash}`}
                                className="download-btn"
                            >
                                Download
                            </a>
                        </td>
                        <td>
                            <a
                                href={`/file/delete?filehash=${fileMeta.FileHash}`}
                                className="download-btn"
                            >
                                Delete
                            </a>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
        </div>
    );
};

export default UserFiles;