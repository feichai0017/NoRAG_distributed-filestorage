"use client"

import React, { useState } from 'react';
import { Search, FileText, AlertCircle } from 'lucide-react';
import { queryAPI } from "@/api/files";
import {
    Box,
    Typography,
    TextField,
    Button,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    Paper,
    Alert,
    InputAdornment,
    Autocomplete,
    useTheme
} from '@mui/material';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { motion, AnimatePresence } from 'framer-motion';

const EnhancedQueryFile = () => {
    const theme = useTheme();
    const [filehash, setFilehash] = useState('');
    const [results, setResults] = useState([]);
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [searchHistory] = useState(['example1', 'example2', 'example3']);
    const [suggestions] = useState(['suggestion1', 'suggestion2', 'suggestion3']);

    // 固定的搜索统计数据
    const stats = [
        { name: 'Mon', searches: 4 },
        { name: 'Tue', searches: 7 },
        { name: 'Wed', searches: 2 },
        { name: 'Thu', searches: 5 },
        { name: 'Fri', searches: 9 },
        { name: 'Sat', searches: 3 },
        { name: 'Sun', searches: 6 },
    ];

    const handleInputChange = (event, newValue) => {
        setFilehash(newValue);
    };

    const handleFormSubmit = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        setError('');
        setResults([]);

        try {
            const data = await queryAPI({ filehash: filehash });
            if (!data || data.error) {
                setError(data ? data.error : 'No data received');
            } else {
                setResults(data);
            }
        } catch (error) {
            setError('Error fetching data');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Box
            sx={{
                minHeight: '100vh',
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'flex-start',
                p: 4,
                bgcolor: 'background.default',
                color: 'text.primary',
            }}
        >
            <Paper
                elevation={24}
                sx={{
                    width: '100%',
                    maxWidth: 'md',
                    borderRadius: 4,
                    overflow: 'hidden',
                    transition: 'all 0.3s ease-in-out',
                    '&:hover': {
                        transform: 'scale(1.02)',
                    },
                    bgcolor: 'background.paper',
                }}
            >
                <Box sx={{ p: 4 }}>
                    <Typography
                        variant="h4"
                        component="h1"
                        align="center"
                        sx={{
                            mb: 3,
                            fontWeight: 'bold',
                            background: theme.palette.mode === 'dark'
                                ? 'linear-gradient(45deg, #90caf9 30%, #64b5f6 90%)'
                                : 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
                            WebkitBackgroundClip: 'text',
                            WebkitTextFillColor: 'transparent'
                        }}
                    >
                        File Search
                    </Typography>
                    <form onSubmit={handleFormSubmit}>
                        <Autocomplete
                            freeSolo
                            options={[...searchHistory, ...suggestions]}
                            renderInput={(params) => (
                                <TextField
                                    {...params}
                                    fullWidth
                                    variant="outlined"
                                    placeholder="Enter filehash"
                                    required
                                    InputProps={{
                                        ...params.InputProps,
                                        startAdornment: (
                                            <InputAdornment position="start">
                                                <Search color={theme.palette.text.secondary} />
                                            </InputAdornment>
                                        ),
                                    }}
                                />
                            )}
                            value={filehash}
                            onChange={handleInputChange}
                            sx={{ mb: 3 }}
                        />
                        <Button
                            type="submit"
                            fullWidth
                            variant="contained"
                            disabled={isLoading}
                            sx={{
                                py: 1.5,
                                background: theme.palette.mode === 'dark'
                                    ? 'linear-gradient(45deg, #90caf9 30%, #64b5f6 90%)'
                                    : 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
                                color: theme.palette.mode === 'dark' ? 'text.primary' : 'white',
                                '&:hover': {
                                    background: theme.palette.mode === 'dark'
                                        ? 'linear-gradient(45deg, #82b1ff 30%, #448aff 90%)'
                                        : 'linear-gradient(45deg, #1e88e5 30%, #1cb5e0 90%)',
                                },
                            }}
                        >
                            {isLoading ? 'Searching...' : 'Search'}
                        </Button>
                    </form>

                    <AnimatePresence>
                        {error && (
                            <motion.div
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, y: -20 }}
                                transition={{ duration: 0.3 }}
                            >
                                <Alert
                                    severity="error"
                                    icon={<AlertCircle />}
                                    sx={{ mt: 3 }}
                                >
                                    {error}
                                </Alert>
                            </motion.div>
                        )}
                    </AnimatePresence>

                    <AnimatePresence>
                        {results.length > 0 && (
                            <motion.div
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, y: -20 }}
                                transition={{ duration: 0.3 }}
                            >
                                <Box sx={{ mt: 4 }}>
                                    <Typography variant="h5" component="h2" sx={{ mb: 2, fontWeight: 'medium' }}>
                                        Search Results
                                    </Typography>
                                    <List>
                                        {results.map((result, index) => (
                                            <motion.div
                                                key={index}
                                                initial={{ opacity: 0, x: -20 }}
                                                animate={{ opacity: 1, x: 0 }}
                                                transition={{ duration: 0.3, delay: index * 0.1 }}
                                            >
                                                <ListItem
                                                    sx={{
                                                        bgcolor: 'background.paper',
                                                        borderRadius: 2,
                                                        mb: 2,
                                                        transition: 'all 0.3s ease-in-out',
                                                        '&:hover': {
                                                            boxShadow: 3,
                                                        },
                                                    }}
                                                >
                                                    <ListItemIcon>
                                                        <FileText color={theme.palette.primary.main} />
                                                    </ListItemIcon>
                                                    <ListItemText
                                                        primary={result.file_name}
                                                        secondary={
                                                            <>
                                                                <Typography component="span" variant="body2" color="text.primary">
                                                                    Size: {result.file_size}
                                                                </Typography>
                                                                <br />
                                                                <Typography component="span" variant="body2" color="text.primary">
                                                                    Location: {result.location}
                                                                </Typography>
                                                            </>
                                                        }
                                                    />
                                                </ListItem>
                                            </motion.div>
                                        ))}
                                    </List>
                                </Box>
                            </motion.div>
                        )}
                    </AnimatePresence>
                </Box>
            </Paper>

            <Paper
                elevation={24}
                sx={{
                    width: '100%',
                    maxWidth: 'md',
                    borderRadius: 4,
                    overflow: 'hidden',
                    mt: 4,
                    p: 4,
                    transition: 'all 0.3s ease-in-out',
                    '&:hover': {
                        transform: 'scale(1.02)',
                    },
                    bgcolor: 'background.paper',
                }}
            >
                <Typography variant="h5" component="h2" sx={{ mb: 2, fontWeight: 'medium', color: 'text.primary' }}>
                    Search Statistics
                </Typography>
                <Box sx={{ height: 300 }}>
                    <ResponsiveContainer width="100%" height="100%">
                        <LineChart data={stats}>
                            <CartesianGrid strokeDasharray="3 3" stroke={theme.palette.divider} />
                            <XAxis dataKey="name" stroke={theme.palette.text.primary} />
                            <YAxis stroke={theme.palette.text.primary} />
                            <Tooltip
                                contentStyle={{
                                    backgroundColor: theme.palette.background.paper,
                                    color: theme.palette.text.primary,
                                    border: `1px solid ${theme.palette.divider}`,
                                }}
                            />
                            <Line type="monotone" dataKey="searches" stroke={theme.palette.primary.main} strokeWidth={2} dot={{ r: 4 }} activeDot={{ r: 8 }} />
                        </LineChart>
                    </ResponsiveContainer>
                </Box>
            </Paper>
        </Box>
    );
};

export default EnhancedQueryFile;