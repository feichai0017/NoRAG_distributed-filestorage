import React, { useState, useEffect } from 'react';
import { Container, Grid, Paper, Typography, List, ListItem, ListItemText, ListItemIcon, IconButton, Tooltip } from '@mui/material';
import { BarChart, CloudUpload, CloudDownload, InsertDriveFile, Refresh } from '@mui/icons-material';
import { ResponsiveContainer, BarChart as RechartsBarChart, Bar, XAxis, YAxis, Tooltip as RechartsTooltip, Legend, CartesianGrid } from 'recharts';
import { motion, AnimatePresence } from 'framer-motion';

const MotionPaper = motion(Paper);
const MotionListItem = motion(ListItem);

const EnhancedDashboard = () => {
    const [weeklyData, setWeeklyData] = useState([]);
    const [recentUploads, setRecentUploads] = useState([]);
    const [isLoading, setIsLoading] = useState(true);

    const fetchData = async () => {
        setIsLoading(true);
        // Simulating API calls to fetch data
        await new Promise(resolve => setTimeout(resolve, 1000)); // Simulated delay
        fetchWeeklyData();
        fetchRecentUploads();
        setIsLoading(false);
    };

    useEffect(() => {
        fetchData();
    }, []);

    const fetchWeeklyData = () => {
        // Simulated data - replace with actual API call
        const days = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
        const data = days.map(day => ({
            day,
            uploads: Math.floor(Math.random() * 10),
            downloads: Math.floor(Math.random() * 15),
        }));
        setWeeklyData(data);
    };

    const fetchRecentUploads = () => {
        // Simulated data - replace with actual API call
        const uploads = [
            { id: 1, name: 'document.pdf', size: '2.5 MB', date: '2023-06-01' },
            { id: 2, name: 'image.jpg', size: '1.8 MB', date: '2023-05-30' },
            { id: 3, name: 'spreadsheet.xlsx', size: '3.2 MB', date: '2023-05-28' },
        ];
        setRecentUploads(uploads);
    };

    const containerVariants = {
        hidden: { opacity: 0 },
        visible: {
            opacity: 1,
            transition: {
                when: "beforeChildren",
                staggerChildren: 0.1
            }
        }
    };

    const itemVariants = {
        hidden: { y: 20, opacity: 0 },
        visible: {
            y: 0,
            opacity: 1,
            transition: {
                type: "spring",
                stiffness: 100
            }
        }
    };

    return (
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
            <motion.div
                variants={containerVariants}
                initial="hidden"
                animate="visible"
            >
                <Grid container spacing={3}>
                    {/* Weekly Upload/Download Frequency Chart */}
                    <Grid item xs={12}>
                        <MotionPaper
                            elevation={3}
                            variants={itemVariants}
                            whileHover={{ scale: 1.02 }}
                            transition={{ type: "spring", stiffness: 400, damping: 10 }}
                            sx={{
                                p: 2,
                                display: 'flex',
                                flexDirection: 'column',
                                height: 400,
                            }}
                        >
                            <Typography component="h2" variant="h6" color="primary" gutterBottom>
                                Weekly Upload/Download Frequency
                            </Typography>
                            <ResponsiveContainer width="100%" height="100%">
                                <RechartsBarChart
                                    data={weeklyData}
                                    margin={{
                                        top: 20,
                                        right: 30,
                                        left: 20,
                                        bottom: 5,
                                    }}
                                >
                                    <CartesianGrid strokeDasharray="3 3" />
                                    <XAxis dataKey="day" />
                                    <YAxis />
                                    <RechartsTooltip />
                                    <Legend />
                                    <Bar dataKey="uploads" name="Uploads" fill="#8884d8" animationBegin={0} animationDuration={1500} />
                                    <Bar dataKey="downloads" name="Downloads" fill="#82ca9d" animationBegin={0} animationDuration={1500} />
                                </RechartsBarChart>
                            </ResponsiveContainer>
                        </MotionPaper>
                    </Grid>
                    {/* Recent Uploads */}
                    <Grid item xs={12}>
                        <MotionPaper
                            elevation={3}
                            variants={itemVariants}
                            sx={{
                                p: 2,
                                display: 'flex',
                                flexDirection: 'column',
                            }}
                        >
                            <Typography component="h2" variant="h6" color="primary" gutterBottom>
                                Recent Uploads
                                <Tooltip title="Refresh data">
                                    <IconButton onClick={fetchData} size="small" sx={{ ml: 1 }}>
                                        <Refresh />
                                    </IconButton>
                                </Tooltip>
                            </Typography>
                            <List>
                                <AnimatePresence>
                                    {recentUploads.map((upload, index) => (
                                        <MotionListItem
                                            key={upload.id}
                                            initial={{ opacity: 0, x: -20 }}
                                            animate={{ opacity: 1, x: 0 }}
                                            exit={{ opacity: 0, x: 20 }}
                                            transition={{ delay: index * 0.1 }}
                                            whileHover={{
                                                backgroundColor: 'rgba(0, 0, 0, 0.04)',
                                                scale: 1.02,
                                                transition: { duration: 0.2 }
                                            }}
                                        >
                                            <ListItemIcon>
                                                <InsertDriveFile />
                                            </ListItemIcon>
                                            <ListItemText
                                                primary={upload.name}
                                                secondary={`Size: ${upload.size} | Uploaded on: ${upload.date}`}
                                            />
                                        </MotionListItem>
                                    ))}
                                </AnimatePresence>
                            </List>
                        </MotionPaper>
                    </Grid>
                </Grid>
            </motion.div>
        </Container>
    );
};

export default EnhancedDashboard;