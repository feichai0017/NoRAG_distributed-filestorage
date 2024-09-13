import React, { useState, useEffect } from 'react';
import { Download, Trash2, FileText, AlertCircle } from 'lucide-react';
import { deleteAPI, downloadAPI, queryAllAPI } from "@/api/files.jsx";
import { getToken } from "@/utils/index.jsx";
import { ThemeProvider, useTheme } from '@mui/material/styles';
import {
    Box,
    Typography,
    TextField,
    Button,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    IconButton,
    CircularProgress,
    Alert,
} from '@mui/material';

const EnhancedUserFiles = () => {
    const [limit, setLimit] = useState(10);
    const [files, setFiles] = useState([]);
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const theme = useTheme();

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
        <ThemeProvider theme={theme}>
            <Box sx={{
                minHeight: '100vh',
                bgcolor: 'background.default',
                py: 6,
                px: { xs: 2, sm: 3, md: 4 },
            }}>
                <Paper elevation={24} sx={{
                    maxWidth: 'xl',
                    mx: 'auto',
                    borderRadius: 4,
                    overflow: 'hidden',
                }}>
                    <Box sx={{ p: 4 }}>
                        <Typography variant="h4" component="h1" align="center" sx={{
                            mb: 3,
                            fontWeight: 'bold',
                            background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
                            WebkitBackgroundClip: 'text',
                            WebkitTextFillColor: 'transparent'
                        }}>
                            User Uploaded Files
                        </Typography>

                        <Box component="form" sx={{ mb: 4 }}>
                            <Typography variant="subtitle1" sx={{ mb: 1 }}>
                                Number of Files:
                            </Typography>
                            <Box sx={{ display: 'flex' }}>
                                <TextField
                                    type="number"
                                    id="limit"
                                    name="limit"
                                    inputProps={{ min: 1, max: 100 }}
                                    value={limit}
                                    onChange={handleInputChange}
                                    required
                                    sx={{ flexGrow: 1, mr: 1 }}
                                />
                                <Button
                                    onClick={fetchFiles}
                                    variant="contained"
                                    sx={{
                                        background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
                                        color: 'white',
                                    }}
                                >
                                    Fetch Files
                                </Button>
                            </Box>
                        </Box>

                        {error && (
                            <Alert severity="error" icon={<AlertCircle />} sx={{ mb: 2 }}>
                                {error}
                            </Alert>
                        )}

                        {isLoading ? (
                            <Box sx={{ display: 'flex', justifyContent: 'center', my: 8 }}>
                                <CircularProgress />
                            </Box>
                        ) : (
                            <TableContainer component={Paper} sx={{ maxHeight: 440 }}>
                                <Table stickyHeader aria-label="user files table">
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>File Name</TableCell>
                                            <TableCell>File Hash</TableCell>
                                            <TableCell>File Size (Bytes)</TableCell>
                                            <TableCell>Upload Time</TableCell>
                                            <TableCell>Actions</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {files.map((fileMeta, index) => (
                                            <TableRow key={index} hover>
                                                <TableCell>
                                                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                                                        <FileText sx={{ mr: 1, color: 'text.secondary' }} />
                                                        <Typography variant="body2">{fileMeta.FileName}</Typography>
                                                    </Box>
                                                </TableCell>
                                                <TableCell>
                                                    <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                                                        {fileMeta.FileHash}
                                                    </Typography>
                                                </TableCell>
                                                <TableCell>
                                                    <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                                                        {fileMeta.FileSize}
                                                    </Typography>
                                                </TableCell>
                                                <TableCell>
                                                    <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                                                        {fileMeta.UploadAt}
                                                    </Typography>
                                                </TableCell>
                                                <TableCell>
                                                    <IconButton
                                                        onClick={() => handleDownload(fileMeta.FileHash)}
                                                        color="primary"
                                                    >
                                                        <Download />
                                                    </IconButton>
                                                    <IconButton
                                                        onClick={() => handleDelete(fileMeta.FileHash)}
                                                        color="error"
                                                    >
                                                        <Trash2 />
                                                    </IconButton>
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </TableContainer>
                        )}
                    </Box>
                </Paper>
            </Box>
        </ThemeProvider>
    );
};

export default EnhancedUserFiles;