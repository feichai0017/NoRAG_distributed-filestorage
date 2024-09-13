import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
    Box, Typography, TextField, Button, Avatar, IconButton,
    Card, CardContent, Grid, Select, MenuItem,
    Container, Paper, Chip, useMediaQuery, Snackbar, Alert, CircularProgress,
    Tabs, Tab, Fade
} from '@mui/material';
import { ThemeProvider, useTheme, styled } from '@mui/material/styles';
import {
    Edit as EditIcon,
    Email as EmailIcon,
    AccessTime as AccessTimeIcon,
    ArrowBack as ArrowBackIcon,
    Save as SaveIcon,
    CloudUpload as CloudUploadIcon,
    Person as PersonIcon,
    Work as WorkIcon,
    Info as InfoIcon
} from '@mui/icons-material';
import { useNavigate } from "react-router-dom";

const MotionContainer = motion(Container);
const MotionPaper = motion(Paper);
const MotionCard = motion(Card);
const MotionTypography = motion(Typography);

const StyledAvatar = styled(Avatar)(({ theme }) => ({
    width: 150,
    height: 150,
    border: `4px solid ${theme.palette.background.paper}`,
    boxShadow: theme.shadows[3],
}));

const StyledIconButton = styled(IconButton)(({ theme }) => ({
    position: 'absolute',
    bottom: 0,
    right: 0,
    backgroundColor: theme.palette.background.paper,
    '&:hover': {
        backgroundColor: theme.palette.action.hover,
    },
}));

const TabPanel = ({ children, value, index, ...other }) => (
    <div
        role="tabpanel"
        hidden={value !== index}
        id={`profile-tabpanel-${index}`}
        aria-labelledby={`profile-tab-${index}`}
        {...other}
    >
        {value === index && (
            <Fade in={value === index}>
                <Box sx={{ p: 3 }}>{children}</Box>
            </Fade>
        )}
    </div>
);

export default function EnhancedUserProfile() {
    const navigate = useNavigate();
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const [name, setName] = useState('');
    const [role, setRole] = useState('');
    const [email, setEmail] = useState('');
    const [country, setCountry] = useState('');
    const [timezone, setTimezone] = useState('');
    const [bio, setBio] = useState('');
    const [skills, setSkills] = useState([]);
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isSaving, setIsSaving] = useState(false);
    const [successMessage, setSuccessMessage] = useState('');
    const [tabValue, setTabValue] = useState(0);

    useEffect(() => {
        const fetchUserData = async () => {
            setIsLoading(true);
            try {
                // Simulating an API call
                const response = await fetch('https://api.example.com/user-profile');
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const userData = await response.json();
                setName(userData.name);
                setRole(userData.role);
                setEmail(userData.email);
                setCountry(userData.country);
                setTimezone(userData.timezone);
                setBio(userData.bio);
                setSkills(userData.skills);
                console.log('User data fetched successfully:', userData);
            } catch (err) {
                console.error('Error fetching user data:', err);
                setError('Failed to load user data. Please refresh the page and try again.');
            } finally {
                setIsLoading(false);
            }
        };

        fetchUserData();
    }, []);

    const handleSave = async () => {
        setIsSaving(true);
        setError(null);
        try {
            // Simulating an API call
            const response = await fetch('https://api.example.com/update-profile', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ name, role, email, country, timezone, bio, skills }),
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const result = await response.json();
            console.log('Profile saved successfully:', result);
            setSuccessMessage('Profile updated successfully!');
        } catch (err) {
            console.error('Error saving profile:', err);
            setError('Failed to save profile. Please try again.');
        } finally {
            setIsSaving(false);
        }
    };

    const handleAddSkill = (skill) => {
        if (skill && !skills.includes(skill)) {
            setSkills([...skills, skill]);
        }
    };

    const handleRemoveSkill = (skillToRemove) => {
        setSkills(skills.filter(skill => skill !== skillToRemove));
    };

    const handleTabChange = (event, newValue) => {
        setTabValue(newValue);
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

    if (isLoading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
                <CircularProgress />
            </Box>
        );
    }

    return (
        <ThemeProvider theme={theme}>
            <MotionContainer maxWidth="lg" sx={{ py: 4 }}>
                <AnimatePresence>
                    <motion.div
                        variants={containerVariants}
                        initial="hidden"
                        animate="visible"
                        exit="hidden"
                    >
                        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
                            <Button
                                startIcon={<ArrowBackIcon />}
                                variant="outlined"
                                component={motion.button}
                                whileHover={{ scale: 1.05 }}
                                whileTap={{ scale: 0.95 }}
                                onClick={() => navigate('/')}
                            >
                                Back to Dashboard
                            </Button>
                            <MotionTypography variant="h4" component="h1" variants={itemVariants}>
                                My Profile
                            </MotionTypography>
                        </Box>

                        <MotionPaper elevation={3} sx={{ p: 3, mb: 4, bgcolor: 'background.paper', borderRadius: 2 }} variants={itemVariants}>
                            <Grid container spacing={3} alignItems="center">
                                <Grid item xs={12} sm={4} sx={{ display: 'flex', justifyContent: 'center' }}>
                                    <Box sx={{ position: 'relative' }}>
                                        <StyledAvatar alt={name} src="/placeholder.svg" />
                                        <StyledIconButton size="small">
                                            <EditIcon fontSize="small" />
                                        </StyledIconButton>
                                    </Box>
                                </Grid>
                                <Grid item xs={12} sm={8}>
                                    <Typography variant="h4" gutterBottom>{name}</Typography>
                                    <Typography variant="h6" color="text.secondary" gutterBottom>{role}</Typography>
                                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                                        <EmailIcon sx={{ mr: 1, color: 'text.secondary' }} />
                                        <Typography variant="body1">{email}</Typography>
                                    </Box>
                                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                                        <AccessTimeIcon sx={{ mr: 1, color: 'text.secondary' }} />
                                        <Typography variant="body1">{`${country}, ${timezone}`}</Typography>
                                    </Box>
                                </Grid>
                            </Grid>
                        </MotionPaper>

                        <MotionCard variants={itemVariants} sx={{ mb: 4, borderRadius: 2 }}>
                            <Tabs value={tabValue} onChange={handleTabChange} aria-label="profile tabs" centered>
                                <Tab icon={<PersonIcon />} label="Personal Info" />
                                <Tab icon={<WorkIcon />} label="Skills & Projects" />
                                <Tab icon={<InfoIcon />} label="Bio" />
                            </Tabs>

                            <TabPanel value={tabValue} index={0}>
                                <Grid container spacing={3}>
                                    <Grid item xs={12} sm={6}>
                                        <TextField
                                            fullWidth
                                            label="Name"
                                            value={name}
                                            onChange={(e) => setName(e.target.value)}
                                            variant="outlined"
                                            margin="normal"
                                        />
                                    </Grid>
                                    <Grid item xs={12} sm={6}>
                                        <TextField
                                            fullWidth
                                            label="Role"
                                            value={role}
                                            onChange={(e) => setRole(e.target.value)}
                                            variant="outlined"
                                            margin="normal"
                                        />
                                    </Grid>
                                    <Grid item xs={12} sm={6}>
                                        <TextField
                                            fullWidth
                                            label="Email"
                                            value={email}
                                            onChange={(e) => setEmail(e.target.value)}
                                            variant="outlined"
                                            margin="normal"
                                            InputProps={{
                                                startAdornment: <EmailIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                                            }}
                                        />
                                    </Grid>
                                    <Grid item xs={12} sm={6}>
                                        <TextField
                                            fullWidth
                                            label="Country"
                                            value={country}
                                            onChange={(e) => setCountry(e.target.value)}
                                            variant="outlined"
                                            margin="normal"
                                        />
                                    </Grid>
                                    <Grid item xs={12}>
                                        <TextField
                                            fullWidth
                                            select
                                            label="Timezone"
                                            value={timezone}
                                            onChange={(e) => setTimezone(e.target.value)}
                                            variant="outlined"
                                            margin="normal"
                                            InputProps={{
                                                startAdornment: <AccessTimeIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                                            }}
                                        >
                                            <MenuItem value="GMT-05:00">Eastern Time (ET) — GMT-05:00</MenuItem>
                                            <MenuItem value="GMT-08:00">Pacific Time (PT) — GMT-08:00</MenuItem>
                                        </TextField>
                                    </Grid>
                                </Grid>
                            </TabPanel>

                            <TabPanel value={tabValue} index={1}>
                                <Typography variant="h6" gutterBottom>Skills</Typography>
                                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 2 }}>
                                    {skills.map((skill) => (
                                        <Chip
                                            key={skill}
                                            label={skill}
                                            onDelete={() => handleRemoveSkill(skill)}
                                            color="primary"
                                            variant="outlined"
                                        />
                                    ))}
                                </Box>
                                <TextField
                                    fullWidth
                                    label="Add a skill"
                                    variant="outlined"
                                    onKeyPress={(e) => {
                                        if (e.key === 'Enter') {
                                            handleAddSkill(e.target.value);
                                            e.target.value = '';
                                        }
                                    }}
                                />
                                <Typography variant="h6" gutterBottom sx={{ mt: 4 }}>Portfolio Projects</Typography>
                                <Button
                                    variant="outlined"
                                    component="label"
                                    startIcon={<CloudUploadIcon />}
                                    fullWidth={isMobile}
                                >
                                    Upload Project
                                    <input type="file" hidden />
                                </Button>
                            </TabPanel>

                            <TabPanel value={tabValue} index={2}>
                                <TextField
                                    fullWidth
                                    multiline
                                    rows={6}
                                    label="Bio"
                                    value={bio}
                                    onChange={(e) => setBio(e.target.value)}
                                    variant="outlined"
                                    margin="normal"
                                />
                            </TabPanel>
                        </MotionCard>

                        <Box sx={{ mt: 4, display: 'flex', justifyContent: 'flex-end' }}>
                            <Button
                                variant="contained"
                                color="primary"
                                startIcon={<SaveIcon />}
                                onClick={handleSave}
                                component={motion.button}
                                whileHover={{ scale: 1.05 }}
                                whileTap={{ scale: 0.95 }}
                                disabled={isSaving}
                            >
                                {isSaving ? 'Saving...' : 'Save Changes'}
                            </Button>
                        </Box>
                    </motion.div>
                </AnimatePresence>
            </MotionContainer>
            <Snackbar open={!!error} autoHideDuration={6000} onClose={() => setError(null)}>
                <Alert onClose={() => setError(null)} severity="error" sx={{ width: '100%' }}>
                    {error}
                </Alert>
            </Snackbar>
            <Snackbar open={!!successMessage} autoHideDuration={6000} onClose={() => setSuccessMessage('')}>
                <Alert onClose={() => setSuccessMessage('')} severity="success" sx={{ width: '100%' }}>
                    {successMessage}
                </Alert>
            </Snackbar>
        </ThemeProvider>
    );
}