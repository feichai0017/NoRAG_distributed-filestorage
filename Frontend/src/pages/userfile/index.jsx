// eslint-disable-next-line no-unused-vars
import React, { useState, useEffect } from 'react';
import { Download, Trash2, FileText, AlertCircle } from 'lucide-react';
import { deleteAPI, downloadAPI, queryAllAPI } from "@/api/files.jsx";
import { getToken } from "@/utils/index.jsx";

const EnhancedUserFiles = () => {
    const [limit, setLimit] = useState(10);
    const [files, setFiles] = useState([]);
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const handleInputChange = (e) => {
        setLimit(e.target.value);
    };

    const fetchFiles = async () => {
        setIsLoading(true);
        setError('');
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
                setFiles([]);
                setError('Invalid data received from server');
            }
        } catch (error) {
            console.error('Error fetching files:', error);
            setFiles([]);
            setError('Failed to fetch files. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchFiles();
    }, [limit]);

    const handleDownload = async (fileHash) => {
        try {
            const response = await downloadAPI({ filehash: fileHash });
            console.log('Download response:', response);
            // Handle the download response here
        } catch (error) {
            console.error('Error downloading file:', error);
            setError('Failed to download file. Please try again.');
        }
    };

    const handleDelete = async (fileHash) => {
        try {
            const response = await deleteAPI({ filehash: fileHash });
            console.log('Delete response:', response);
            setFiles(files.filter(file => file.FileHash !== fileHash));
        } catch (error) {
            console.error('Error deleting file:', error);
            setError('Failed to delete file. Please try again.');
        }
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-gray-100 to-gray-200 py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-7xl mx-auto bg-white rounded-2xl shadow-xl overflow-hidden">
                <div className="p-8">
                    <h1 className="text-3xl font-bold mb-6 text-center bg-clip-text text-transparent bg-gradient-to-r from-blue-500 to-purple-500">User Uploaded Files</h1>

                    <form id="queryForm" className="mb-8">
                        <label htmlFor="limit" className="block text-sm font-medium text-gray-700 mb-2">Number of Files:</label>
                        <div className="flex items-center">
                            <input
                                type="number"
                                id="limit"
                                name="limit"
                                min="1"
                                max="100"
                                value={limit}
                                onChange={handleInputChange}
                                required
                                className="flex-1 px-4 py-2 border border-gray-300 rounded-l-md focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                            />
                            <button
                                onClick={fetchFiles}
                                className="px-4 py-2 border border-transparent text-sm font-medium rounded-r-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                            >
                                Fetch Files
                            </button>
                        </div>
                    </form>

                    {error && (
                        <div className="mb-4 p-4 bg-red-100 border-l-4 border-red-500 rounded-r-md">
                            <div className="flex items-center">
                                <AlertCircle className="h-5 w-5 text-red-500 mr-2" />
                                <p className="text-red-700">{error}</p>
                            </div>
                        </div>
                    )}

                    {isLoading ? (
                        <div className="flex justify-center items-center h-64">
                            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
                        </div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="min-w-full divide-y divide-gray-200">
                                <thead className="bg-gray-50">
                                <tr>
                                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">File Name</th>
                                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">File Hash</th>
                                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">File Size (Bytes)</th>
                                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Upload Time</th>
                                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                                </tr>
                                </thead>
                                <tbody className="bg-white divide-y divide-gray-200">
                                {files.map((fileMeta, index) => (
                                    <tr key={index} className="hover:bg-gray-50 transition-colors duration-200">
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <div className="flex items-center">
                                                <FileText className="h-5 w-5 text-gray-400 mr-2" />
                                                <span className="text-sm font-medium text-gray-900">{fileMeta.FileName}</span>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span className="text-sm text-gray-500">{fileMeta.FileHash}</span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span className="text-sm text-gray-500">{fileMeta.FileSize}</span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span className="text-sm text-gray-500">{fileMeta.UploadAt}</span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                                            <button
                                                onClick={() => handleDownload(fileMeta.FileHash)}
                                                className="text-blue-600 hover:text-blue-900 mr-4 transition-colors duration-200"
                                            >
                                                <Download className="h-5 w-5" />
                                            </button>
                                            <button
                                                onClick={() => handleDelete(fileMeta.FileHash)}
                                                className="text-red-600 hover:text-red-900 transition-colors duration-200"
                                            >
                                                <Trash2 className="h-5 w-5" />
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                                </tbody>
                            </table>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default EnhancedUserFiles;