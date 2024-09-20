import React, { useState } from 'react';
import {
    Paper,
    Typography,
    Button,
    Box,
    TextField,
    List,
    ListItem,
    ListItemText,
    IconButton,
    LinearProgress,
    InputAdornment,
} from '@mui/material';
import { Upload, Delete, CheckCircle, Search } from '@mui/icons-material';
import { motion } from 'framer-motion';

const MotionPaper = motion(Paper);

function KnowledgeBaseDataset({ knowledgeBase, onUpdate }) {
    const [files, setFiles] = useState([]);
    const [uploadProgress, setUploadProgress] = useState({});
    const [searchQuery, setSearchQuery] = useState('');

    const handleFileUpload = (event) => {
        const newFiles = Array.from(event.target.files);
        setFiles([...files, ...newFiles]);
        newFiles.forEach((file) => {
            simulateFileUpload(file);
        });
        onUpdate({
            ...knowledgeBase,
            documentCount: knowledgeBase.documentCount + newFiles.length,
        });
    };

    const simulateFileUpload = (file) => {
        setUploadProgress((prev) => ({ ...prev, [file.name]: 0 }));
        const interval = setInterval(() => {
            setUploadProgress((prev) => {
                const newProgress = Math.min((prev[file.name] || 0) + 10, 100);
                if (newProgress === 100) {
                    clearInterval(interval);
                }
                return { ...prev, [file.name]: newProgress };
            });
        }, 500);
    };

    const handleDeleteFile = (fileName) => {
        setFiles(files.filter((file) => file.name !== fileName));
        setUploadProgress((prev) => {
            const newProgress = { ...prev };
            delete newProgress[fileName];
            return newProgress;
        });
        onUpdate({
            ...knowledgeBase,
            documentCount: knowledgeBase.documentCount - 1,
        });
    };

    const filteredFiles = files.filter((file) =>
        file.name.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <MotionPaper
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.3 }}
            elevation={0}
            sx={{ p: 3, borderRadius: 2, bgcolor: 'background.paper', transition: 'background-color 0.3s' }}
        >
            <Typography variant="h5" gutterBottom>Dataset</Typography>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Button
                    variant="contained"
                    component="label"
                    startIcon={<Upload />}
                >
                    Upload File
                    <input
                        type="file"
                        hidden
                        onChange={handleFileUpload}
                        multiple
                    />
                </Button>
                <TextField
                    placeholder="Search documents"
                    variant="outlined"
                    size="small"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <Search />
                            </InputAdornment>
                        ),
                    }}
                />
            </Box>
            <List>
                {filteredFiles.map((file) => (
                    <ListItem key={file.name}>
                        <ListItemText primary={file.name} />
                        {uploadProgress[file.name] === 100 ? (
                            <CheckCircle color="success" />
                        ) : (
                            <Box sx={{ width: '100%', mr: 1 }}>
                                <LinearProgress variant="determinate" value={uploadProgress[file.name] || 0} />
                            </Box>
                        )}
                        <IconButton edge="end" aria-label="delete" onClick={() => handleDeleteFile(file.name)}>
                            <Delete />
                        </IconButton>
                    </ListItem>
                ))}
            </List>
        </MotionPaper>
    );
}

export default KnowledgeBaseDataset;