import React, { useState, useEffect } from 'react';
import {
    Container, Grid, Paper, Typography, List, ListItem, ListItemText,
    ListItemIcon, IconButton, Tooltip, Skeleton, useTheme, useMediaQuery
} from '@mui/material';
import {
    BarChart as BarChartIcon, CloudUpload, CloudDownload, InsertDriveFile,
    Refresh, TrendingUp, TrendingDown
} from '@mui/icons-material';
import {
    ResponsiveContainer, BarChart as RechartsBarChart, Bar, XAxis, YAxis,
    Tooltip as RechartsTooltip, Legend, CartesianGrid
} from 'recharts';
import { motion, AnimatePresence } from 'framer-motion';

const MotionPaper = motion(Paper);
const MotionListItem = motion(ListItem);

const EnhancedDashboard = () => {
    const [weeklyData, setWeeklyData] = useState([]);
    const [recentUploads, setRecentUploads] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const fetchData = async () => {
        setIsLoading(true);
        await new Promise(resolve => setTimeout(resolve, 1500)); // Simulated delay
        fetchWeeklyData();
        fetchRecentUploads();
        setIsLoading(false);
    };

    useEffect(() => {
        fetchData();
    }, []);

    const fetchWeeklyData = () => {
        const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
        const data = days.map(day => ({
            day,
            uploads: Math.floor(Math.random() * 10),
            downloads: Math.floor(Math.random() * 15),
        }));
        setWeeklyData(data);
    };

    const fetchRecentUploads = () => {
        const uploads = [
            { id: 1, name: 'Q2 Financial Report.pdf', size: '2.5 MB', date: '2023-06-01', trend: 'up' },
            { id: 2, name: 'Product Launch Presentation.pptx', size: '5.8 MB', date: '2023-05-30', trend: 'down' },
            { id: 3, name: 'Customer Feedback Analysis.xlsx', size: '1.2 MB', date: '2023-05-28', trend: 'up' },
            { id: 4, name: 'Team Building Event Photos.zip', size: '15.7 MB', date: '2023-05-26', trend: 'down' },
            { id: 5, name: 'Project Roadmap.docx', size: '0.8 MB', date: '2023-05-24', trend: 'up' },
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

    const chartColors = {
        uploads: theme.palette.primary.main,
        downloads: theme.palette.secondary.main,
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
                            whileHover={{ scale: 1.01 }}
                            transition={{ type: "spring", stiffness: 400, damping: 10 }}
                            sx={{
                                p: 3,
                                display: 'flex',
                                flexDirection: 'column',
                                height: isMobile ? 300 : 400,
                                borderRadius: 2,
                                overflow: 'hidden',
                            }}
                        >
                            <Typography component="h2" variant="h6" color="primary" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
                                <BarChartIcon sx={{ mr: 1 }} />
                                Weekly Activity Overview
                            </Typography>
                            {isLoading ? (
                                <Skeleton variant="rectangular" height="100%" animation="wave" />
                            ) : (
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
                                        <RechartsTooltip
                                            contentStyle={{
                                                backgroundColor: theme.palette.background.paper,
                                                border: `1px solid ${theme.palette.divider}`,
                                                borderRadius: 4,
                                            }}
                                        />
                                        <Legend />
                                        <Bar
                                            dataKey="uploads"
                                            name="Uploads"
                                            fill={chartColors.uploads}
                                            animationBegin={0}
                                            animationDuration={1500}
                                            radius={[4, 4, 0, 0]}
                                        />
                                        <Bar
                                            dataKey="downloads"
                                            name="Downloads"
                                            fill={chartColors.downloads}
                                            animationBegin={0}
                                            animationDuration={1500}
                                            radius={[4, 4, 0, 0]}
                                        />
                                    </RechartsBarChart>
                                </ResponsiveContainer>
                            )}
                        </MotionPaper>
                    </Grid>
                    {/* Recent Uploads */}
                    <Grid item xs={12}>
                        <MotionPaper
                            elevation={3}
                            variants={itemVariants}
                            sx={{
                                p: 3,
                                display: 'flex',
                                flexDirection: 'column',
                                borderRadius: 2,
                                overflow: 'hidden',
                            }}
                        >
                            <Typography component="h2" variant="h6" color="primary" gutterBottom sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <span style={{ display: 'flex', alignItems: 'center' }}>
                  <CloudUpload sx={{ mr: 1 }} />
                  Recent Uploads
                </span>
                                <Tooltip title="Refresh data">
                                    <IconButton onClick={fetchData} size="small" sx={{ ml: 1 }}>
                                        <Refresh />
                                    </IconButton>
                                </Tooltip>
                            </Typography>
                            <List>
                                <AnimatePresence>
                                    {isLoading ? (
                                        [...Array(5)].map((_, index) => (
                                            <ListItem key={index} divider>
                                                <ListItemIcon>
                                                    <Skeleton variant="circular" width={24} height={24} />
                                                </ListItemIcon>
                                                <ListItemText
                                                    primary={<Skeleton width="60%" />}
                                                    secondary={<Skeleton width="40%" />}
                                                />
                                            </ListItem>
                                        ))
                                    ) : (
                                        recentUploads.map((upload, index) => (
                                            <MotionListItem
                                                key={upload.id}
                                                initial={{ opacity: 0, x: -20 }}
                                                animate={{ opacity: 1, x: 0 }}
                                                exit={{ opacity: 0, x: 20 }}
                                                transition={{ delay: index * 0.1 }}
                                                whileHover={{
                                                    backgroundColor: theme.palette.action.hover,
                                                    scale: 1.02,
                                                    transition: { duration: 0.2 }
                                                }}
                                                divider
                                            >
                                                <ListItemIcon>
                                                    <InsertDriveFile color="primary" />
                                                </ListItemIcon>
                                                <ListItemText
                                                    primary={upload.name}
                                                    secondary={`Size: ${upload.size} | Uploaded on: ${upload.date}`}
                                                />
                                                <Tooltip title={upload.trend === 'up' ? 'Trending Up' : 'Trending Down'}>
                                                    <IconButton size="small" color={upload.trend === 'up' ? 'success' : 'error'}>
                                                        {upload.trend === 'up' ? <TrendingUp /> : <TrendingDown />}
                                                    </IconButton>
                                                </Tooltip>
                                            </MotionListItem>
                                        ))
                                    )}
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