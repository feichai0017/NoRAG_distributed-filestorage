import React, { useState } from 'react';
import {
    Box,
    Typography,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    IconButton,
    TextField,
    InputAdornment,
    Select,
    MenuItem,
    Breadcrumbs,
    Pagination,
} from '@mui/material';
import {
    Search as SearchIcon,
    GetApp as DownloadIcon,
    Visibility as PreviewIcon,
    Description as FileIcon,
    Folder as FolderIcon,
} from '@mui/icons-material';

const rootFiles = [
    { name: 'knowledgebase', type: 'folder' },
    { name: 'document1.txt', type: 'file', size: '10 KB', uploadDate: '2023-01-01' },
];

const knowledgebaseFiles = [
    { name: 'computer', type: 'folder' },
    { name: 'biology', type: 'folder' },
];

const computerFiles = [
    { name: 'Practical Task.pdf', type: 'file', size: '178.17 KB', uploadDate: '20/09/2024 21:47:15', knowledgeBase: 'computer' },
    { name: 'Computer Science 101.docx', type: 'file', size: '2.5 MB', uploadDate: '19/09/2024 15:30:00', knowledgeBase: 'computer' },
];

export default function FileManager() {
    const [currentPath, setCurrentPath] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [page, setPage] = useState(1);

    const getCurrentFiles = () => {
        if (currentPath.length === 0) return rootFiles;
        if (currentPath.length === 1 && currentPath[0] === 'knowledgebase') return knowledgebaseFiles;
        if (currentPath.length === 2 && currentPath[1] === 'computer') return computerFiles;
        return [];
    };

    const files = getCurrentFiles();

    const handleFileClick = (fileName) => {
        setCurrentPath([...currentPath, fileName]);
    };

    const handleBreadcrumbClick = (index) => {
        setCurrentPath(currentPath.slice(0, index + 1));
    };

    const filteredFiles = files.filter(file =>
        file.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const pageCount = Math.ceil(filteredFiles.length / itemsPerPage);

    return (
        <Box sx={{ width: '100%', p: 3 }}>
            <Breadcrumbs aria-label="breadcrumb" sx={{ mb: 2 }}>
                <Typography
                    color="inherit"
                    style={{ cursor: 'pointer' }}
                    onClick={() => setCurrentPath([])}
                >
                    root
                </Typography>
                {currentPath.map((path, index) => (
                    <Typography
                        key={path}
                        color={index === currentPath.length - 1 ? 'text.primary' : 'inherit'}
                        style={{ cursor: 'pointer' }}
                        onClick={() => handleBreadcrumbClick(index)}
                    >
                        {path}
                    </Typography>
                ))}
            </Breadcrumbs>

            <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <TextField
                    placeholder="Search files"
                    variant="outlined"
                    size="small"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <SearchIcon />
                            </InputAdornment>
                        ),
                    }}
                />
                <Select
                    value={itemsPerPage}
                    onChange={(e) => setItemsPerPage(e.target.value)}
                    size="small"
                >
                    <MenuItem value={10}>10 items/page</MenuItem>
                    <MenuItem value={20}>20 items/page</MenuItem>
                    <MenuItem value={50}>50 items/page</MenuItem>
                </Select>
            </Box>

            <TableContainer component={Paper} elevation={0}>
                <Table sx={{ minWidth: 650 }} aria-label="file table">
                    <TableHead>
                        <TableRow>
                            <TableCell>Name</TableCell>
                            <TableCell>Upload Date</TableCell>
                            <TableCell>Size</TableCell>
                            <TableCell>Knowledge Base</TableCell>
                            <TableCell>Actions</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {filteredFiles
                            .slice((page - 1) * itemsPerPage, page * itemsPerPage)
                            .map((file) => (
                                <TableRow key={file.name}>
                                    <TableCell component="th" scope="row">
                                        <Box sx={{ display: 'flex', alignItems: 'center' }}>
                                            {file.type === 'folder' ? (
                                                <FolderIcon sx={{ mr: 1, color: 'primary.main' }} />
                                            ) : (
                                                <FileIcon sx={{ mr: 1, color: 'error.main' }} />
                                            )}
                                            <Typography
                                                style={{ cursor: file.type === 'folder' ? 'pointer' : 'default' }}
                                                onClick={() => file.type === 'folder' && handleFileClick(file.name)}
                                            >
                                                {file.name}
                                            </Typography>
                                        </Box>
                                    </TableCell>
                                    <TableCell>{file.uploadDate || '-'}</TableCell>
                                    <TableCell>{file.size || '-'}</TableCell>
                                    <TableCell>
                                        {file.knowledgeBase && (
                                            <Box sx={{ display: 'inline-block', bgcolor: 'primary.main', color: 'white', px: 1, py: 0.5, borderRadius: 1 }}>
                                                {file.knowledgeBase}
                                            </Box>
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        {file.type === 'file' && (
                                            <>
                                                <IconButton size="small" aria-label="download">
                                                    <DownloadIcon />
                                                </IconButton>
                                                <IconButton size="small" aria-label="preview">
                                                    <PreviewIcon />
                                                </IconButton>
                                            </>
                                        )}
                                    </TableCell>
                                </TableRow>
                            ))}
                    </TableBody>
                </Table>
            </TableContainer>

            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mt: 2 }}>
                <Typography variant="body2">Total: {filteredFiles.length} item(s)</Typography>
                <Pagination
                    count={pageCount}
                    page={page}
                    onChange={(event, value) => setPage(value)}
                    color="primary"
                />
            </Box>
        </Box>
    );
}