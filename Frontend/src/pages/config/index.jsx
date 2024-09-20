import React, { useState, useEffect } from 'react';
import {
    Paper,
    Typography,
    TextField,
    Select,
    MenuItem,
    Slider,
    Box,
    Switch,
    Button,
    IconButton,
    CircularProgress,
} from '@mui/material';
import { motion } from 'framer-motion';
import { Upload } from '@mui/icons-material';

const MotionPaper = motion(Paper);

function KnowledgeBaseConfiguration({ knowledgeBase, onUpdate }) {
    const [localKnowledgeBase, setLocalKnowledgeBase] = useState({
        coverImage: '',
        name: '',
        description: '',
        language: 'english',
        embeddingModel: 'bge-large',
        splitTokenCount: 128,
        apiEnabled: false,
    });
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (knowledgeBase) {
            setLocalKnowledgeBase(prevState => ({
                ...prevState,
                ...knowledgeBase,
            }));
            setLoading(false);
        }
    }, [knowledgeBase]);

    const handleImageUpload = (event) => {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                const newCoverImage = e.target.result;
                setLocalKnowledgeBase(prevState => ({
                    ...prevState,
                    coverImage: newCoverImage,
                }));
                onUpdate({ ...localKnowledgeBase, coverImage: newCoverImage });
            };
            reader.readAsDataURL(file);
        }
    };

    const handleInputChange = (field, value) => {
        setLocalKnowledgeBase(prevState => ({
            ...prevState,
            [field]: value,
        }));
        onUpdate({ ...localKnowledgeBase, [field]: value });
    };

    if (loading) {
        return (
            <Box display="flex" justifyContent="center" alignItems="center" height="100%">
                <CircularProgress />
            </Box>
        );
    }

    return (
        <MotionPaper
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.3 }}
            elevation={0}
            sx={{ p: 3, borderRadius: 2, bgcolor: 'background.paper', transition: 'background-color 0.3s' }}
        >
            <Typography variant="h5" gutterBottom>Configuration</Typography>

            <Box sx={{ mb: 3, position: 'relative', width: '100%', height: 200, borderRadius: 2, overflow: 'hidden' }}>
                <img
                    src={localKnowledgeBase.coverImage || '/placeholder.svg'}
                    alt="Knowledge Base Cover"
                    style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                />
                <input
                    accept="image/*"
                    id="cover-image-upload"
                    type="file"
                    hidden
                    onChange={handleImageUpload}
                />
                <label htmlFor="cover-image-upload">
                    <IconButton
                        component="span"
                        sx={{
                            position: 'absolute',
                            bottom: 8,
                            right: 8,
                            bgcolor: 'rgba(255, 255, 255, 0.8)',
                            '&:hover': { bgcolor: 'rgba(255, 255, 255, 0.9)' },
                        }}
                    >
                        <Upload />
                    </IconButton>
                </label>
            </Box>

            <TextField
                label="Knowledge Base Name"
                variant="outlined"
                fullWidth
                value={localKnowledgeBase.name}
                onChange={(e) => handleInputChange('name', e.target.value)}
                sx={{ mb: 2 }}
            />
            <TextField
                label="Description"
                variant="outlined"
                fullWidth
                multiline
                rows={4}
                value={localKnowledgeBase.description}
                onChange={(e) => handleInputChange('description', e.target.value)}
                sx={{ mb: 2 }}
            />
            <Select
                label="Language"
                variant="outlined"
                fullWidth
                value={localKnowledgeBase.language}
                onChange={(e) => handleInputChange('language', e.target.value)}
                sx={{ mb: 2 }}
            >
                <MenuItem value="english">English</MenuItem>
                <MenuItem value="chinese">Chinese</MenuItem>
            </Select>
            <Typography variant="subtitle1" gutterBottom>Embedding Model</Typography>
            <Select
                label="Embedding Model"
                variant="outlined"
                fullWidth
                value={localKnowledgeBase.embeddingModel}
                onChange={(e) => handleInputChange('embeddingModel', e.target.value)}
                sx={{ mb: 2 }}
            >
                <MenuItem value="bge-large">BGE-large</MenuItem>
                <MenuItem value="openai">OpenAI</MenuItem>
            </Select>
            <Typography variant="subtitle1" gutterBottom>Split Token Count</Typography>
            <Slider
                value={localKnowledgeBase.splitTokenCount}
                onChange={(e, newValue) => handleInputChange('splitTokenCount', newValue)}
                step={1}
                min={1}
                max={512}
                valueLabelDisplay="auto"
                sx={{ mb: 2 }}
            />
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="subtitle1">Enable API</Typography>
                <Switch
                    checked={localKnowledgeBase.apiEnabled}
                    onChange={(e) => handleInputChange('apiEnabled', e.target.checked)}
                />
            </Box>
            <Button variant="contained" color="primary" onClick={() => onUpdate(localKnowledgeBase)}>
                Save Changes
            </Button>
        </MotionPaper>
    );
}

export default KnowledgeBaseConfiguration;