import React, { useState, useEffect, useRef } from 'react';
import { Upload, X, File, CheckCircle } from 'lucide-react';
import { ThemeProvider, useTheme } from '@mui/material/styles';
import { Box, Typography, Button, Select, MenuItem, LinearProgress } from '@mui/material';
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
    const abortController = useRef(new AbortController());
    const [uploadId, setUploadId] = useState(null);
    const [dragActive, setDragActive] = useState(false);
    const theme = useTheme();

    useEffect(() => {
        return () => {
            if (abortController.current) {
                abortController.current.abort();
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
        abortController.current = new AbortController();

        const formData = new FormData();
        formData.append('file', file);
        formData.append('filename', fileName);

        try {
            let response;
            if (uploadMode === 'normal' || file.size <= 5 * 1024 * 1024) {
                response = await uploadAPI(formData, {
                    signal: abortController.current.signal,
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
            if (error.name === 'AbortError') {
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
                signal: abortController.current.signal,
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
        if (abortController.current) {
            abortController.current.abort();
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
        <ThemeProvider theme={theme}>
            <Box
                sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    minHeight: '100vh',
                    bgcolor: 'background.default',
                    color: 'text.primary',
                }}
            >
                <Box
                    sx={{
                        width: '100%',
                        maxWidth: 'md',
                        bgcolor: 'background.paper',
                        borderRadius: 4,
                        boxShadow: 24,
                        overflow: 'hidden',
                        transition: 'all 0.3s ease-in-out',
                        '&:hover': {
                            transform: 'scale(1.05)',
                        },
                    }}
                >
                    <Box sx={{ p: 4 }}>
                        <Typography variant="h4" component="h2" align="center" sx={{ mb: 3, background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                            File Upload
                        </Typography>
                        {!file ? (
                            <Box
                                sx={{
                                    border: 2,
                                    borderStyle: 'dashed',
                                    borderRadius: 2,
                                    p: 4,
                                    textAlign: 'center',
                                    transition: 'all 0.3s ease-in-out',
                                    borderColor: dragActive ? 'primary.main' : 'grey.300',
                                    bgcolor: dragActive ? 'action.hover' : 'background.paper',
                                }}
                                onDragEnter={handleDrag}
                                onDragLeave={handleDrag}
                                onDragOver={handleDrag}
                                onDrop={handleDrop}
                            >
                                <Upload sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
                                <Typography variant="h6" sx={{ mb: 2 }}>Drag and drop your file here, or</Typography>
                                <Button
                                    component="label"
                                    variant="contained"
                                    sx={{
                                        background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
                                        color: 'white',
                                        borderRadius: 50,
                                    }}
                                >
                                    Choose File
                                    <input
                                        type="file"
                                        hidden
                                        onChange={handleFileChange}
                                        accept=".pdf,.txt,.doc,.docx"
                                    />
                                </Button>
                                <Typography variant="body2" sx={{ mt: 2, color: 'text.secondary' }}>
                                    Supported formats: PDF, TXT, DOC, DOCX
                                </Typography>
                            </Box>
                        ) : (
                            <form onSubmit={handleUpload}>
                                <Box sx={{ display: 'flex', alignItems: 'center', p: 2, bgcolor: 'action.hover', borderRadius: 2, mb: 2 }}>
                                    <File sx={{ fontSize: 32, color: 'primary.main', mr: 2 }} />
                                    <Typography variant="body1" sx={{ flexGrow: 1, overflow: 'hidden', textOverflow: 'ellipsis' }}>
                                        {fileName}
                                    </Typography>
                                </Box>
                                <Select
                                    value={uploadMode}
                                    onChange={(e) => setUploadMode(e.target.value)}
                                    fullWidth
                                    sx={{ mb: 2 }}
                                >
                                    <MenuItem value="normal">Normal Upload</MenuItem>
                                    <MenuItem value="multipart">Multipart Upload</MenuItem>
                                </Select>
                                <Box sx={{ display: 'flex', gap: 2 }}>
                                    <Button
                                        type="submit"
                                        variant="contained"
                                        fullWidth
                                        disabled={isUploading}
                                        sx={{
                                            background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
                                            color: 'white',
                                            borderRadius: 50,
                                        }}
                                    >
                                        {isUploading ? 'Uploading...' : 'Upload'}
                                    </Button>
                                    {isUploading && (
                                        <Button
                                            onClick={handleCancelUpload}
                                            variant="outlined"
                                            sx={{ borderRadius: 50 }}
                                        >
                                            <X />
                                        </Button>
                                    )}
                                </Box>
                            </form>
                        )}
                        {(isUploading || uploadComplete) && (
                            <Box sx={{ mt: 3 }}>
                                <LinearProgress
                                    variant="determinate"
                                    value={uploadProgress}
                                    sx={{
                                        height: 8,
                                        borderRadius: 5,
                                        bgcolor: 'action.hover',
                                        '& .MuiLinearProgress-bar': {
                                            borderRadius: 5,
                                            background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
                                        },
                                    }}
                                />
                                <Typography variant="body2" align="center" sx={{ mt: 1 }}>
                                    {uploadComplete ? (
                                        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'success.main' }}>
                                            <CheckCircle sx={{ mr: 1 }} />
                                            Upload Complete
                                        </Box>
                                    ) : (
                                        `${uploadProgress}% Uploaded`
                                    )}
                                </Typography>
                            </Box>
                        )}
                    </Box>
                </Box>
            </Box>
        </ThemeProvider>
    );
};

export default EnhancedFileUpload;