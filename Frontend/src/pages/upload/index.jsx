import './index.css'
import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';

const Upload = () => {
    const [file, setFile] = useState(null);
    const [fileName, setFileName] = useState('');
    const [uploadMode, setUploadMode] = useState('normal');
    const [uploadProgress, setUploadProgress] = useState(0);
    const [isUploading, setIsUploading] = useState(false);
    const cancelTokenSource = useRef(null);

    useEffect(() => {
        return () => {
            if (cancelTokenSource.current) {
                cancelTokenSource.current.cancel('Component unmounted');
            }
        };
    }, []);

    const handleFileChange = (e) => {
        const selectedFile = e.target.files[0];
        setFile(selectedFile);
        setFileName(selectedFile.name);
    };

    const handleButtonClick = () => {
        document.getElementById('fileID').click();
    };

    const handleUpload = async (e) => {
        e.preventDefault();
        if (!file) {
            alert('Please select a file first!');
            return;
        }

        setIsUploading(true);
        setUploadProgress(0);
        cancelTokenSource.current = axios.CancelToken.source();

        const formData = new FormData();
        formData.append('file', file);
        formData.append('filename', fileName);

        try {
            let response;
            if (uploadMode === 'normal' || file.size <= 5 * 1024 * 1024) {
                response = await axios.post('/api/upload', formData, {
                    cancelToken: cancelTokenSource.current.token,
                    onUploadProgress: (progressEvent) => {
                        const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
                        setUploadProgress(percentCompleted);
                    }
                });
            } else {
                response = await initiateMultipartUpload(formData);
            }
            console.log(response.data);
            alert('File uploaded successfully!');
        } catch (error) {
            if (axios.isCancel(error)) {
                console.log('Upload cancelled');
            } else {
                console.error('Upload failed:', error);
                alert('Upload failed. Please try again.');
            }
        } finally {
            setIsUploading(false);
        }
    };

    const initiateMultipartUpload = async (formData) => {
        const initResponse = await axios.post('/api/mpupload/init', formData, {
            cancelToken: cancelTokenSource.current.token
        });
        const { uploadID, chunkSize, chunkCount } = initResponse.data;

        for (let i = 0; i < chunkCount; i++) {
            const start = i * chunkSize;
            const end = Math.min(start + chunkSize, file.size);
            const chunk = file.slice(start, end);

            const chunkFormData = new FormData();
            chunkFormData.append('file', chunk, `${fileName}.part${i}`);
            chunkFormData.append('uploadid', uploadID);
            chunkFormData.append('index', i);

            await axios.post('/api/mpupload/uploadpart', chunkFormData, {
                cancelToken: cancelTokenSource.current.token,
                onUploadProgress: (progressEvent) => {
                    const percentCompleted = Math.round(((i * chunkSize + progressEvent.loaded) * 100) / file.size);
                    setUploadProgress(percentCompleted);
                }
            });
        }

        return axios.post('/api/mpupload/complete', {
            uploadid: uploadID,
            filehash: await calculateFileHash(file),
            filesize: file.size,
            filename: fileName
        }, {
            cancelToken: cancelTokenSource.current.token
        });
    };

    const calculateFileHash = async (file) => {
        return Math.random().toString(36).substring(7);
    };

    const handleCancelUpload = () => {
        if (cancelTokenSource.current) {
            cancelTokenSource.current.cancel('Upload cancelled by user');
        }
    };

    return (
        <div className="upload-file">
            <div className="card">
                <h3>Upload Files</h3>
                {!file ? (
                    <div className="drop_box">
                        <header>
                            <h4>Select File here</h4>
                        </header>
                        <p>Files Supported: PDF, TEXT, DOC, DOCX</p>
                        <input type="file" id="fileID" style={{display: 'none'}} onChange={handleFileChange}/>
                        <button className="btn" onClick={handleButtonClick}>Choose File</button>
                    </div>
                ) : (
                    <form onSubmit={handleUpload}>
                        <div className="form">
                            <h4>{fileName}</h4>
                            <select
                                value={uploadMode}
                                onChange={(e) => setUploadMode(e.target.value)}
                            >
                                <option value="normal">Normal Upload</option>
                                <option value="multipart">Multipart Upload</option>
                            </select>
                            <button type="submit" className="btn" disabled={isUploading}>
                                {isUploading ? 'Uploading...' : 'Upload'}
                            </button>
                            {isUploading && (
                                <button type="button" className="btn cancel" onClick={handleCancelUpload}>
                                    Cancel Upload
                                </button>
                            )}
                        </div>
                    </form>
                )}
                {isUploading && (
                    <div className="progress-bar">
                        <div className="progress" style={{width: `${uploadProgress}%`}}></div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default Upload;