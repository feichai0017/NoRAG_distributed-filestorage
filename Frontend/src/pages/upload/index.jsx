// eslint-disable-next-line no-unused-vars
import React, { useState, useEffect, useRef } from 'react';
import { Upload, X, File, CheckCircle } from 'lucide-react';
import axios from 'axios';
import {
    cancelMultipartUploadAPI,
    completeMultipartUploadAPI,
    getMultipartUploadStatusAPI,
    initMultipartUploadAPI,
    uploadAPI,
    uploadPartAPI
} from "@/api/files.jsx";

const EnhancedFileUpload = () => {
    const [file, setFile] = useState(null);
    const [fileName, setFileName] = useState('');
    const [uploadMode, setUploadMode] = useState('normal');
    const [uploadProgress, setUploadProgress] = useState(0);
    const [isUploading, setIsUploading] = useState(false);
    const [uploadComplete, setUploadComplete] = useState(false);
    const cancelTokenSource = useRef(null);
    const [uploadId, setUploadId] = useState(null);
    const [dragActive, setDragActive] = useState(false);

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
        setUploadComplete(false);
    };

    const handleDrag = (e) => {
        e.preventDefault();
        e.stopPropagation();
        if (e.type === "dragenter" || e.type === "dragover") {
            setDragActive(true);
        } else if (e.type === "dragleave") {
            setDragActive(false);
        }
    };

    const handleDrop = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setDragActive(false);
        if (e.dataTransfer.files && e.dataTransfer.files[0]) {
            handleFileChange({ target: { files: e.dataTransfer.files } });
        }
    };

    const handleUpload = async (e) => {
        e.preventDefault();
        if (!file) {
            alert('Please select a file first!');
            return;
        }

        setIsUploading(true);
        setUploadProgress(0);
        setUploadComplete(false);
        cancelTokenSource.current = axios.CancelToken.source();

        const formData = new FormData();
        formData.append('file', file);
        formData.append('filename', fileName);

        try {
            let response;
            if (uploadMode === 'normal' || file.size <= 5 * 1024 * 1024) {
                response = await uploadAPI(formData, {
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
            setUploadComplete(true);
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
        const initResponse = await initMultipartUploadAPI(formData);
        const { uploadID, chunkSize, chunkCount } = initResponse.data;
        setUploadId(uploadID);

        for (let i = 0; i < chunkCount; i++) {
            const start = i * chunkSize;
            const end = Math.min(start + chunkSize, file.size);
            const chunk = file.slice(start, end);

            const chunkFormData = new FormData();
            chunkFormData.append('file', chunk, `${fileName}.part${i}`);
            chunkFormData.append('uploadid', uploadID);
            chunkFormData.append('index', i);

            await uploadPartAPI(chunkFormData, {
                cancelToken: cancelTokenSource.current.token,
                onUploadProgress: (progressEvent) => {
                    const percentCompleted = Math.round(((i * chunkSize + progressEvent.loaded) * 100) / file.size);
                    setUploadProgress(percentCompleted);
                }
            });
        }

        return completeMultipartUploadAPI({
            uploadid: uploadID,
            filehash: await calculateFileHash(file),
            filesize: file.size,
            filename: fileName
        });
    };

    const handleCancelUpload = async () => {
        if (cancelTokenSource.current) {
            cancelTokenSource.current.cancel('Upload cancelled by user');
        }
        if (uploadId) {
            await cancelMultipartUploadAPI({ uploadid: uploadId });
            setUploadId(null);
        }
        setIsUploading(false);
        setUploadProgress(0);
    };

    const calculateFileHash = async (file) => {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = async (e) => {
                try {
                    const arrayBuffer = e.target.result;
                    const hashBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer);
                    const hashArray = Array.from(new Uint8Array(hashBuffer));
                    const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
                    resolve(hashHex);
                } catch (error) {
                    reject(error);
                }
            };
            reader.onerror = error => reject(error);
            reader.readAsArrayBuffer(file);
        });
    };

    const checkUploadStatus = async () => {
        if (uploadId) {
            const status = await getMultipartUploadStatusAPI({ uploadid: uploadId });
            console.log('Upload status:', status.data);
        }
    };

    return (
        <div className="flex items-center justify-center min-h-screen bg-gradient-to-br from-gray-100 to-gray-200">
            <div className="w-full max-w-md bg-white rounded-2xl shadow-xl overflow-hidden transition-all duration-300 ease-in-out transform hover:scale-105">
                <div className="p-8">
                    <h2 className="text-3xl font-bold mb-6 text-center bg-clip-text text-transparent bg-gradient-to-r from-blue-500 to-purple-500">File Upload</h2>
                    {!file ? (
                        <div
                            className={`border-2 border-dashed rounded-xl p-8 text-center transition-all duration-300 ease-in-out ${dragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300'}`}
                            onDragEnter={handleDrag}
                            onDragLeave={handleDrag}
                            onDragOver={handleDrag}
                            onDrop={handleDrop}
                        >
                            <Upload className="mx-auto h-16 w-16 text-gray-400 mb-4" />
                            <p className="text-lg text-gray-600 mb-4">Drag and drop your file here, or</p>
                            <label htmlFor="fileInput" className="cursor-pointer inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-full text-white bg-gradient-to-r from-blue-500 to-purple-500 hover:from-blue-600 hover:to-purple-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-all duration-300 ease-in-out">
                                Choose File
                                <input
                                    id="fileInput"
                                    type="file"
                                    className="sr-only"
                                    onChange={handleFileChange}
                                    accept=".pdf,.txt,.doc,.docx"
                                />
                            </label>
                            <p className="mt-4 text-sm text-gray-500">Supported formats: PDF, TXT, DOC, DOCX</p>
                        </div>
                    ) : (
                        <form onSubmit={handleUpload} className="space-y-6">
                            <div className="flex items-center space-x-4 p-4 bg-gray-50 rounded-lg">
                                <File className="h-8 w-8 text-blue-500" />
                                <span className="text-lg font-medium text-gray-700 truncate flex-1">{fileName}</span>
                            </div>
                            <select
                                value={uploadMode}
                                onChange={(e) => setUploadMode(e.target.value)}
                                className="block w-full px-4 py-3 rounded-lg border-2 border-gray-200 focus:border-blue-500 focus:ring focus:ring-blue-200 focus:ring-opacity-50 transition-all duration-300 ease-in-out"
                            >
                                <option value="normal">Normal Upload</option>
                                <option value="multipart">Multipart Upload</option>
                            </select>
                            <div className="flex space-x-4">
                                <button
                                    type="submit"
                                    className={`flex-1 px-6 py-3 border border-transparent text-base font-medium rounded-full text-white bg-gradient-to-r from-blue-500 to-purple-500 hover:from-blue-600 hover:to-purple-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-all duration-300 ease-in-out ${isUploading ? 'opacity-50 cursor-not-allowed' : ''}`}
                                    disabled={isUploading}
                                >
                                    {isUploading ? 'Uploading...' : 'Upload'}
                                </button>
                                {isUploading && (
                                    <button
                                        type="button"
                                        onClick={handleCancelUpload}
                                        className="px-4 py-3 border-2 border-gray-300 text-base font-medium rounded-full text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-all duration-300 ease-in-out"
                                    >
                                        <X className="h-5 w-5" />
                                    </button>
                                )}
                            </div>
                        </form>
                    )}
                    {(isUploading || uploadComplete) && (
                        <div className="mt-6 space-y-4">
                            <div className="relative pt-1">
                                <div className="overflow-hidden h-2 text-xs flex rounded-full bg-blue-200">
                                    <div
                                        style={{ width: `${uploadProgress}%` }}
                                        className="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-gradient-to-r from-blue-500 to-purple-500 transition-all duration-300 ease-in-out"
                                    ></div>
                                </div>
                            </div>
                            <p className="text-sm text-gray-600 text-center">
                                {uploadComplete ? (
                                    <span className="flex items-center justify-center text-green-500">
                    <CheckCircle className="h-5 w-5 mr-2" />
                    Upload Complete
                  </span>
                                ) : (
                                    <span className="font-semibold">{`${uploadProgress}% Uploaded`}</span>
                                )}
                            </p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default EnhancedFileUpload;