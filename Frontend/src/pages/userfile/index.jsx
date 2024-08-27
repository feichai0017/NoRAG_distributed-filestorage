import './index.css'
import React, { useState, useEffect } from 'react';
import {deleteAPI, downloadAPI, queryAllAPI} from "@/api/files.jsx";
import {getToken} from "@/utils/index.jsx";

const UserFiles = () => {
    const [limit, setLimit] = useState(10);
    const [files, setFiles] = useState([]);

    const handleInputChange = (e) => {
        setLimit(e.target.value);
    };

    const fetchFiles = async () => {
        try {
            const limitInt = parseInt(limit, 10);
            console.log("Token before request:", getToken());
            const data = await queryAllAPI({
                limit: limitInt,
            });
            if (Array.isArray(data)) {
                setFiles(data);
            } else {
                console.error('Received data is not an array:', data);
                setFiles([]); // 设置为空数组
            }
        } catch (error) {
            console.error('Error fetching files:', error);
            setFiles([]); // 出错时设置为空数组
        }

    };

    useEffect(() => {
        fetchFiles();
    }, [limit]); // 组件加载时获取文件列表


    return (
        <div className="file-display">
            <div className="user-files">
                <h1>User Uploaded Files</h1>
                <form id="queryForm" >
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
                                <button
                                    onClick={async () => {
                                        try {
                                            const response = await downloadAPI({filehash: fileMeta.FileHash});
                                            // 处理下载响应
                                            console.log('Download response:', response);
                                        } catch (error) {
                                            console.error('Error downloading file:', error);
                                        }
                                    }}
                                    className="download-btn"
                                >
                                    Download
                                </button>
                            </td>
                            <td>
                                <button
                                    onClick={async () => {
                                        try {
                                            const response = await deleteAPI({filehash: fileMeta.FileHash});
                                            // 处理删除响应
                                            console.log('Delete response:', response);
                                            // 从文件列表中移除已删除的文件
                                            setFiles(files.filter(file => file.FileHash !== fileMeta.FileHash));
                                        } catch (error) {
                                            console.error('Error deleting file:', error);
                                        }
                                    }}
                                    className="delete-btn"
                                >
                                    Delete
                                </button>
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