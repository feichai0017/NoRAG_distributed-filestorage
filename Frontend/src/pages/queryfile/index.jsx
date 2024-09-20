import React, { useState } from 'react';
import {
    Box,
    Typography,
    TextField,
    Button,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    Checkbox,
    Paper,
    Divider,
    useTheme,
    alpha,
    IconButton,
} from '@mui/material';
import {
    Search as SearchIcon,
    Description as FileIcon,
    Folder as FolderIcon,
    ExpandMore as ExpandMoreIcon,
    ChevronRight as ChevronRightIcon,
    AccountTree as MindMapIcon,
    Menu as MenuIcon,
} from '@mui/icons-material';
import { motion, AnimatePresence } from 'framer-motion';

const MotionPaper = motion(Paper);
const MotionBox = motion(Box);

const TreeView = ({ data, onSelect, selectedItems }) => {
    const [expanded, setExpanded] = useState({});

    const toggleExpand = (id) => {
        setExpanded(prev => ({ ...prev, [id]: !prev[id] }));
    };

    const handleSelect = (node) => {
        if (node.children) {
            const childIds = node.children.map(child => child.id);
            if (childIds.every(id => selectedItems.includes(id))) {
                onSelect(childIds, false);
            } else {
                onSelect(childIds, true);
            }
        } else {
            onSelect([node.id], !selectedItems.includes(node.id));
        }
    };

    const renderTree = (nodes) => (
        <List>
            {nodes.map((node) => (
                <motion.div
                    key={node.id}
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.3 }}
                >
                    <ListItem
                        button
                        onClick={() => handleSelect(node)}
                        sx={{
                            pl: node.level * 2,
                            borderRadius: 1,
                            mb: 0.5,
                            '&:hover': {
                                bgcolor: 'action.hover',
                            },
                        }}
                    >
                        <ListItemIcon onClick={(e) => {
                            e.stopPropagation();
                            if (node.children) toggleExpand(node.id);
                        }}>
                            {node.children ? (
                                expanded[node.id] ? <ExpandMoreIcon /> : <ChevronRightIcon />
                            ) : (
                                <FileIcon />
                            )}
                        </ListItemIcon>
                        <Checkbox
                            edge="start"
                            checked={node.children
                                ? node.children.every(child => selectedItems.includes(child.id))
                                : selectedItems.includes(node.id)
                            }
                            tabIndex={-1}
                            disableRipple
                            sx={{ '& .MuiSvgIcon-root': { fontSize: 20 } }}
                        />
                        <ListItemText primary={node.name} primaryTypographyProps={{ variant: 'body2' }} />
                    </ListItem>
                    {node.children && expanded[node.id] && renderTree(node.children)}
                </motion.div>
            ))}
        </List>
    );

    return renderTree(data);
};

const EnhancedAISearch = () => {
    const theme = useTheme();
    const [selectedItems, setSelectedItems] = useState([]);
    const [searchQuery, setSearchQuery] = useState('');
    const [searchResults, setSearchResults] = useState(null);
    const [sidebarOpen, setSidebarOpen] = useState(true);

    const fileStructure = [
        {
            id: 'kb1',
            name: 'Computer Science',
            level: 0,
            children: [
                { id: 'file1', name: 'Practical Task.pdf', level: 1 },
                { id: 'file2', name: 'Lecture Notes.docx', level: 1 },
            ],
        },
        {
            id: 'kb2',
            name: 'Mathematics',
            level: 0,
            children: [
                { id: 'file3', name: 'Calculus Basics.pdf', level: 1 },
                { id: 'file4', name: 'Linear Algebra.pdf', level: 1 },
            ],
        },
    ];

    const handleSelect = (ids, isSelected) => {
        if (isSelected) {
            setSelectedItems(prev => [...new Set([...prev, ...ids])]);
        } else {
            setSelectedItems(prev => prev.filter(id => !ids.includes(id)));
        }
    };

    const handleSearch = () => {
        if (selectedItems.length === 0) {
            alert('Please select at least one file or knowledge base');
            return;
        }
        setSearchResults({
            title: `AI-generated answer for "${searchQuery}"`,
            content: `This is a simulated AI-generated answer based on the content of the selected files:
      ${selectedItems.join(', ')}.
      
      The search query "${searchQuery}" has been processed across all selected files.
      
      In a real implementation, the AI would analyze the content of all selected files,
      understand the context of the query, and provide a comprehensive response
      synthesizing information from all relevant sources.`
        });
    };

    const handleGenerateMindMap = () => {
        console.log('Generating mind map for:', selectedItems);
        // Implement mind map generation logic here
    };

    return (
        <Box sx={{ display: 'flex', height: '100vh', bgcolor: 'background.default' }}>
            {/* Sidebar */}
            <AnimatePresence>
                {sidebarOpen && (
                    <MotionPaper
                        initial={{ x: -300 }}
                        animate={{ x: 0 }}
                        exit={{ x: -300 }}
                        transition={{ type: 'spring', stiffness: 300, damping: 30 }}
                        elevation={0}
                        sx={{
                            width: 300,
                            height: '100%',
                            overflowY: 'auto',
                            bgcolor: 'background.paper',
                            borderRight: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
                            position: 'relative',
                            zIndex: 1,
                        }}
                    >
                        <Typography variant="h6" sx={{ p: 2, fontWeight: 'bold', color: 'primary.main' }}>
                            Knowledge Base
                        </Typography>
                        <Divider sx={{ opacity: 0.6 }} />
                        <TreeView
                            data={fileStructure}
                            onSelect={handleSelect}
                            selectedItems={selectedItems}
                        />
                    </MotionPaper>
                )}
            </AnimatePresence>

            {/* Main Content */}
            <Box sx={{ flexGrow: 1, p: 3, overflowY: 'auto' }}>
                <MotionBox
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 0.5 }}
                    sx={{ display: 'flex', alignItems: 'center', mb: 3 }}
                >
                    <IconButton onClick={() => setSidebarOpen(!sidebarOpen)} sx={{ mr: 2 }}>
                        <MenuIcon />
                    </IconButton>
                    <TextField
                        fullWidth
                        variant="outlined"
                        placeholder="Enter your search query"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        InputProps={{
                            startAdornment: <SearchIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                            sx: { borderRadius: 30, bgcolor: 'background.paper' },
                        }}
                    />
                    <Button
                        variant="contained"
                        onClick={handleSearch}
                        disabled={selectedItems.length === 0}
                        sx={{
                            ml: 2,
                            borderRadius: 30,
                            boxShadow: 'none',
                            '&:hover': { boxShadow: 'none' },
                        }}
                    >
                        Search
                    </Button>
                </MotionBox>

                {selectedItems.length > 0 && (
                    <MotionBox
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.5 }}
                        sx={{ mb: 3, p: 2, bgcolor: alpha(theme.palette.primary.main, 0.05), borderRadius: 2 }}
                    >
                        <Typography variant="body2">
                            Selected: <strong>{selectedItems.length} item(s)</strong>
                        </Typography>
                    </MotionBox>
                )}

                <AnimatePresence>
                    {searchResults && (
                        <MotionPaper
                            elevation={0}
                            sx={{ p: 3, borderRadius: 2, bgcolor: alpha(theme.palette.background.paper, 0.8) }}
                            initial={{ y: 50, opacity: 0 }}
                            animate={{ y: 0, opacity: 1 }}
                            exit={{ y: 50, opacity: 0 }}
                            transition={{ type: 'spring', stiffness: 300, damping: 30 }}
                        >
                            <Typography variant="h5" gutterBottom sx={{ fontWeight: 'bold', color: 'primary.main' }}>
                                {searchResults.title}
                            </Typography>
                            <Typography variant="body1" sx={{ whiteSpace: 'pre-line', mb: 3 }}>
                                {searchResults.content}
                            </Typography>
                            <Button
                                variant="outlined"
                                startIcon={<MindMapIcon />}
                                onClick={handleGenerateMindMap}
                                sx={{ borderRadius: 30 }}
                            >
                                Generate Mind Map
                            </Button>
                        </MotionPaper>
                    )}
                </AnimatePresence>
            </Box>
        </Box>
    );
};

export default EnhancedAISearch;