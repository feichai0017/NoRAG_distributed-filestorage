import React, { useState } from 'react';
import {
    Paper,
    Typography,
    Slider,
    TextField,
    Button,
    Box,
    List,
    ListItem,
    ListItemText,
} from '@mui/material';
import { motion } from 'framer-motion';

const MotionPaper = motion(Paper);

function KnowledgeBaseRetrievalTesting({ knowledgeBase }) {
    const [similarityThreshold, setSimilarityThreshold] = useState(0.5);
    const [topK, setTopK] = useState(3);
    const [testQuery, setTestQuery] = useState('');
    const [results, setResults] = useState([]);

    const handleRunTest = () => {
        // Simulate retrieval test
        const simulatedResults = [
            { id: 1, content: 'This is a sample result 1', similarity: 0.85 },
            { id: 2, content: 'This is a sample result 2', similarity: 0.75 },
            { id: 3, content: 'This is a sample result 3', similarity: 0.65 },
        ];
        setResults(simulatedResults);
    };

    return (
        <MotionPaper
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.3 }}
            elevation={0}
            sx={{ p: 3, borderRadius: 2, bgcolor: 'background.paper', transition: 'background-color 0.3s' }}
        >
            <Typography variant="h5" gutterBottom>Retrieval Testing</Typography>
            <Typography variant="subtitle1" gutterBottom>Similarity Threshold</Typography>
            <Slider
                value={similarityThreshold}
                onChange={(e, newValue) => setSimilarityThreshold(newValue)}
                step={0.01}
                min={0}
                max={1}
                valueLabelDisplay="auto"
                sx={{ mb: 2 }}
            />
            <Typography variant="subtitle1" gutterBottom>Top K</Typography>
            <Slider
                value={topK}
                onChange={(e, newValue) => setTopK(newValue)}
                step={1}
                min={1}
                max={10}
                valueLabelDisplay="auto"
                sx={{ mb: 2 }}
            />
            <TextField
                label="Test Query"
                variant="outlined"
                fullWidth
                multiline
                rows={4}
                value={testQuery}
                onChange={(e) => setTestQuery(e.target.value)}
                sx={{ mb: 2 }}
            />
            <Button variant="contained" color="primary" onClick={handleRunTest}>
                Run Test
            </Button>
            {results.length > 0 && (
                <Box sx={{ mt: 3 }}>
                    <Typography variant="h6" gutterBottom>Results:</Typography>
                    <List>
                        {results.map((result) => (
                            <ListItem key={result.id}>
                                <ListItemText
                                    primary={result.content}
                                    secondary={`Similarity: ${result.similarity.toFixed(2)}`}
                                />
                            </ListItem>
                        ))}
                    </List>
                </Box>
            )}
        </MotionPaper>
    );
}

export default KnowledgeBaseRetrievalTesting;