import React, { useState, useEffect } from 'react';
import {
    Container, Grid, Paper, Typography, Box, CircularProgress,
    useTheme, alpha, Button, useMediaQuery
} from '@mui/material';
import {
    CheckCircle as CheckCircleIcon,
    Error as ErrorIcon,
    Storage as StorageIcon,
    Search as SearchIcon,
    CloudQueue as CloudIcon,
    Dns as DnsIcon,
    ArrowBack as ArrowBackIcon
} from '@mui/icons-material';
import {
    LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer
} from 'recharts';
import { useNavigate } from 'react-router-dom';

const services = [
    { name: 'MySQL', icon: StorageIcon, color: '#8884d8' },
    { name: 'Elasticsearch', icon: SearchIcon, color: '#82ca9d' },
    { name: 'MinIO', icon: CloudIcon, color: '#ffc658' },
    { name: 'Ceph', icon: DnsIcon, color: '#ff8042' }
];

const mockFetchServiceData = () => {
    return new Promise((resolve) => {
        setTimeout(() => {
            resolve(services.map(service => ({
                ...service,
                status: Math.random() > 0.8 ? 'down' : 'up',
                data: Array.from({ length: 20 }, (_, i) => ({
                    time: i * 5,
                    throughput: Math.floor(Math.random() * 1000),
                    latency: Math.floor(Math.random() * 100)
                }))
            })));
        }, 1000);
    });
};

export default function Settings() {
    const [serviceData, setServiceData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const theme = useTheme();
    const navigate = useNavigate();
    const isSmallScreen = useMediaQuery(theme.breakpoints.down('md'));

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                setError(null);
                const data = await mockFetchServiceData();
                setServiceData(data);
            } catch (err) {
                console.error("Error fetching service data:", err);
                setError("Failed to fetch service data. Please try again later.");
            } finally {
                setLoading(false);
            }
        };

        fetchData();
        const interval = setInterval(fetchData, 30000);

        return () => clearInterval(interval);
    }, []);

    const getStatusColor = (status) => {
        return status === 'up' ? theme.palette.success.main : theme.palette.error.main;
    };

    const StatusIcon = ({ status }) => {
        return status === 'up' ? (
            <CheckCircleIcon sx={{ color: theme.palette.success.main }} />
        ) : (
            <ErrorIcon sx={{ color: theme.palette.error.main }} />
        );
    };

    if (error) {
        return (
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <Typography color="error">{error}</Typography>
            </Container>
        );
    }

    return (
        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 4 }}>
                <Button
                    variant="contained"
                    startIcon={<ArrowBackIcon />}
                    onClick={() => navigate(-1)}
                    sx={{
                        mr: 2,
                        background: `linear-gradient(45deg, ${theme.palette.primary.main} 30%, ${theme.palette.secondary.main} 90%)`,
                        color: 'white',
                        '&:hover': {
                            background: `linear-gradient(45deg, ${theme.palette.primary.dark} 30%, ${theme.palette.secondary.dark} 90%)`,
                        }
                    }}
                >
                    Back
                </Button>
                <Typography
                    variant="h3"
                    component="h1"
                    sx={{
                        color: theme.palette.primary.main,
                        fontWeight: 'bold',
                        textShadow: `2px 2px 4px ${alpha(theme.palette.primary.main, 0.2)}`
                    }}
                >
                    Service Health Dashboard
                </Typography>
            </Box>
            {loading ? (
                <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
                    <CircularProgress size={60} thickness={4} />
                </Box>
            ) : (
                <Grid container spacing={4}>
                    {serviceData.map((service) => (
                        <Grid item xs={12} lg={6} key={service.name}>
                            <Paper
                                elevation={6}
                                sx={{
                                    p: 3,
                                    display: 'flex',
                                    flexDirection: 'column',
                                    height: isSmallScreen ? 400 : 500,
                                    background: `linear-gradient(45deg, ${alpha(theme.palette.primary.main, 0.05)} 0%, ${alpha(theme.palette.secondary.main, 0.05)} 100%)`,
                                    borderRadius: 4,
                                    overflow: 'hidden',
                                    boxShadow: `0 10px 30px -5px ${alpha(theme.palette.primary.main, 0.2)}`,
                                }}
                            >
                                <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
                                    <Box sx={{
                                        p: 1,
                                        borderRadius: '50%',
                                        backgroundColor: alpha(service.color, 0.1),
                                        mr: 2
                                    }}>
                                        {React.createElement(service.icon, { sx: { fontSize: 40, color: service.color } })}
                                    </Box>
                                    <Box sx={{ flexGrow: 1 }}>
                                        <Typography variant="h5" sx={{ fontWeight: 'bold', color: service.color }}>
                                            {service.name}
                                        </Typography>
                                        <Box sx={{ display: 'flex', alignItems: 'center', mt: 0.5 }}>
                                            <StatusIcon status={service.status} />
                                            <Typography sx={{ ml: 1, color: getStatusColor(service.status), fontWeight: 'medium' }}>
                                                {service.status.toUpperCase()}
                                            </Typography>
                                        </Box>
                                    </Box>
                                </Box>
                                <Box sx={{ flexGrow: 1, minHeight: 0 }}>
                                    <ResponsiveContainer width="100%" height="100%">
                                        <LineChart
                                            data={service.data}
                                            margin={{
                                                top: 5,
                                                right: 30,
                                                left: 20,
                                                bottom: 5,
                                            }}
                                        >
                                            <CartesianGrid strokeDasharray="3 3" stroke={alpha(theme.palette.text.primary, 0.1)} />
                                            <XAxis
                                                dataKey="time"
                                                stroke={theme.palette.text.secondary}
                                                label={{ value: 'Time (minutes)', position: 'insideBottomRight', offset: -10, fill: theme.palette.text.secondary }}
                                            />
                                            <YAxis
                                                yAxisId="left"
                                                stroke={theme.palette.text.secondary}
                                                label={{ value: 'Throughput', angle: -90, position: 'insideLeft', fill: theme.palette.text.secondary }}
                                            />
                                            <YAxis
                                                yAxisId="right"
                                                orientation="right"
                                                stroke={theme.palette.text.secondary}
                                                label={{ value: 'Latency (ms)', angle: 90, position: 'insideRight', fill: theme.palette.text.secondary }}
                                            />
                                            <Tooltip
                                                contentStyle={{
                                                    backgroundColor: theme.palette.background.paper,
                                                    border: `1px solid ${theme.palette.divider}`,
                                                    borderRadius: 4,
                                                }}
                                            />
                                            <Line
                                                yAxisId="left"
                                                type="monotone"
                                                dataKey="throughput"
                                                stroke={service.color}
                                                strokeWidth={2}
                                                dot={false}
                                                activeDot={{ r: 8 }}
                                            />
                                            <Line
                                                yAxisId="right"
                                                type="monotone"
                                                dataKey="latency"
                                                stroke={alpha(service.color, 0.5)}
                                                strokeWidth={2}
                                                dot={false}
                                            />
                                        </LineChart>
                                    </ResponsiveContainer>
                                </Box>
                            </Paper>
                        </Grid>
                    ))}
                </Grid>
            )}
        </Container>
    );
}