import React, { useState } from 'react';
import { useTheme } from '@mui/material/styles';
import {
    Box,
    Typography,
    Button,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    TextField,
    AppBar,
    Toolbar,
    Tabs,
    Tab,
    Grid,
    Card,
    CardContent,
    CardMedia,
    IconButton,
} from '@mui/material';
import { Add, ArrowBack, Close } from '@mui/icons-material';
import { motion, AnimatePresence } from 'framer-motion';
import KnowledgeBaseConfiguration from '../config';
import KnowledgeBaseDataset from '../dataset';
import KnowledgeBaseRetrievalTesting from '../retrievalTesting';

const StyledTab = ({ label, ...props }) => {
    const theme = useTheme();
    return (
        <Tab
            label={label}
            sx={{
                textTransform: 'none',
                minWidth: 0,
                fontWeight: theme.typography.fontWeightRegular,
                marginRight: theme.spacing(4),
                transition: theme.transitions.create('color', {
                    duration: theme.transitions.duration.shortest,
                }),
                '&:hover': {
                    color: theme.palette.primary.main,
                    opacity: 1,
                },
                '&.Mui-selected': {
                    color: theme.palette.primary.main,
                    fontWeight: theme.typography.fontWeightMedium,
                },
            }}
            {...props}
        />
    );
};

const MotionCard = motion(Card);

function KnowledgeBaseManager() {
    const theme = useTheme();
    const [knowledgeBases, setKnowledgeBases] = useState([]);
    const [selectedKnowledgeBase, setSelectedKnowledgeBase] = useState(null);
    const [activeTab, setActiveTab] = useState(0);
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const [newKnowledgeBaseName, setNewKnowledgeBaseName] = useState('');
    const [kbToDelete, setKbToDelete] = useState(null);

    const handleCreateKnowledgeBase = () => {
        const newKnowledgeBase = {
            id: Date.now(),
            name: newKnowledgeBaseName,
            createdAt: new Date().toISOString(),
            documentCount: 0,
            coverImage: `https://source.unsplash.com/random/400x200?${newKnowledgeBaseName}`,
        };
        setKnowledgeBases([...knowledgeBases, newKnowledgeBase]);
        setIsDialogOpen(false);
        setNewKnowledgeBaseName('');
    };

    const handleSelectKnowledgeBase = (kb) => {
        setSelectedKnowledgeBase(kb);
        setActiveTab(0);
    };

    const handleDeleteKnowledgeBase = (kb) => {
        setKbToDelete(kb);
        setIsDeleteDialogOpen(true);
    };

    const confirmDeleteKnowledgeBase = () => {
        setKnowledgeBases(knowledgeBases.filter(kb => kb.id !== kbToDelete.id));
        setIsDeleteDialogOpen(false);
        setKbToDelete(null);
    };

    const handleTabChange = (event, newValue) => {
        setActiveTab(newValue);
    };

    const handleKnowledgeBaseUpdate = (updatedKnowledgeBase) => {
        setKnowledgeBases(prevKnowledgeBases =>
            prevKnowledgeBases.map(kb =>
                kb.id === updatedKnowledgeBase.id ? updatedKnowledgeBase : kb
            )
        );
        setSelectedKnowledgeBase(updatedKnowledgeBase);
    };

    return (
        <Box sx={{ display: 'flex', flexDirection: 'column', height: '100vh', bgcolor: 'background.default', transition: 'background-color 0.3s' }}>
            <AppBar position="static" color="default" elevation={0} sx={{ borderBottom: `1px solid ${theme.palette.divider}` }}>
                <Toolbar>
                    {selectedKnowledgeBase && (
                        <IconButton edge="start" color="inherit" onClick={() => setSelectedKnowledgeBase(null)} sx={{ mr: 2 }}>
                            <ArrowBack />
                        </IconButton>
                    )}
                    <Typography variant="h6" component="div" sx={{ flexGrow: 1, fontWeight: 'bold' }}>
                        {selectedKnowledgeBase ? selectedKnowledgeBase.name : 'Knowledge Base Manager'}
                    </Typography>
                    {!selectedKnowledgeBase && (
                        <Button
                            variant="contained"
                            startIcon={<Add />}
                            onClick={() => setIsDialogOpen(true)}
                            sx={{ borderRadius: 20, ml: 2 }}
                        >
                            Create Knowledge Base
                        </Button>
                    )}
                </Toolbar>
                {selectedKnowledgeBase && (
                    <Tabs value={activeTab} onChange={handleTabChange} aria-label="knowledge base tabs">
                        <StyledTab label="Configuration" />
                        <StyledTab label="Dataset" />
                        <StyledTab label="Retrieval Testing" />
                    </Tabs>
                )}
            </AppBar>

            <Box sx={{ flexGrow: 1, overflow: 'auto', p: 3 }}>
                <AnimatePresence mode="wait">
                    {!selectedKnowledgeBase ? (
                        <motion.div
                            key="knowledge-bases"
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            exit={{ opacity: 0, y: -20 }}
                            transition={{ duration: 0.3 }}
                        >
                            <Grid container spacing={3}>
                                {knowledgeBases.map((kb) => (
                                    <Grid item xs={12} sm={6} md={4} key={kb.id}>
                                        <MotionCard
                                            whileHover={{ scale: 1.05, boxShadow: theme.shadows[8] }}
                                            whileTap={{ scale: 0.95 }}
                                            sx={{ cursor: 'pointer', position: 'relative', overflow: 'hidden' }}
                                        >
                                            <CardMedia
                                                component="img"
                                                height="140"
                                                image={kb.coverImage}
                                                alt={kb.name}
                                                onClick={() => handleSelectKnowledgeBase(kb)}
                                            />
                                            <CardContent onClick={() => handleSelectKnowledgeBase(kb)}>
                                                <Typography variant="h6" component="div" noWrap>
                                                    {kb.name}
                                                </Typography>
                                                <Typography variant="body2" color="text.secondary">
                                                    Created: {new Date(kb.createdAt).toLocaleDateString()}
                                                </Typography>
                                                <Typography variant="body2" color="text.secondary">
                                                    Documents: {kb.documentCount}
                                                </Typography>
                                            </CardContent>
                                            <IconButton
                                                sx={{
                                                    position: 'absolute',
                                                    top: 8,
                                                    right: 8,
                                                    bgcolor: 'rgba(255, 255, 255, 0.7)',
                                                    '&:hover': { bgcolor: 'rgba(255, 255, 255, 0.9)' },
                                                }}
                                                onClick={(e) => {
                                                    e.stopPropagation();
                                                    handleDeleteKnowledgeBase(kb);
                                                }}
                                            >
                                                <Close color="error" />
                                            </IconButton>
                                        </MotionCard>
                                    </Grid>
                                ))}
                            </Grid>
                        </motion.div>
                    ) : (
                        <>
                            {activeTab === 0 && (
                                <KnowledgeBaseConfiguration
                                    knowledgeBase={selectedKnowledgeBase}
                                    onUpdate={handleKnowledgeBaseUpdate}
                                />
                            )}
                            {activeTab === 1 && (
                                <KnowledgeBaseDataset
                                    knowledgeBase={selectedKnowledgeBase}
                                    onUpdate={handleKnowledgeBaseUpdate}
                                />
                            )}
                            {activeTab === 2 && (
                                <KnowledgeBaseRetrievalTesting
                                    knowledgeBase={selectedKnowledgeBase}
                                />
                            )}
                        </>
                    )}
                </AnimatePresence>
            </Box>

            <Dialog open={isDialogOpen} onClose={() => setIsDialogOpen(false)}>
                <DialogTitle>Create New Knowledge Base</DialogTitle>
                <DialogContent>
                    <TextField
                        autoFocus
                        margin="dense"
                        id="name"
                        label="Knowledge Base Name"
                        type="text"
                        fullWidth
                        variant="outlined"
                        value={newKnowledgeBaseName}
                        onChange={(e) => setNewKnowledgeBaseName(e.target.value)}
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setIsDialogOpen(false)}>Cancel</Button>
                    <Button onClick={handleCreateKnowledgeBase} variant="contained">Create</Button>
                </DialogActions>
            </Dialog>

            <Dialog open={isDeleteDialogOpen} onClose={() => setIsDeleteDialogOpen(false)}>
                <DialogTitle>Delete Knowledge Base</DialogTitle>
                <DialogContent>
                    <Typography>Are you sure you want to delete this knowledge base? This action cannot be undone.</Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setIsDeleteDialogOpen(false)}>Cancel</Button>
                    <Button onClick={confirmDeleteKnowledgeBase} variant="contained" color="error">Delete</Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
}

export default KnowledgeBaseManager;