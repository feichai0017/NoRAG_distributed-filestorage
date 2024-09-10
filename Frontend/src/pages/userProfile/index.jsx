// eslint-disable-next-line no-unused-vars
import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
    Box, Typography, TextField, Button, Avatar, IconButton,
    Card, CardContent, Grid, Select, MenuItem,
    Container, Paper, Chip, useMediaQuery
} from '@mui/material';
import { ThemeProvider, createTheme, styled } from '@mui/material/styles';
import {
    Edit as EditIcon,
    Email as EmailIcon,
    AccessTime as AccessTimeIcon,
    ArrowBack as ArrowBackIcon,
    Save as SaveIcon,
    CloudUpload as CloudUploadIcon
} from '@mui/icons-material';
import {useNavigate} from "react-router-dom";

const MotionContainer = motion(Container);
const MotionPaper = motion(Paper);
const MotionCard = motion(Card);
const MotionTypography = motion(Typography);

const theme = createTheme(); // Create a default theme

const StyledAvatar = styled(Avatar)(({ theme }) => ({
    width: 120,
    height: 120,
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

const EnhancedUserProfile = () => {
    const navigate = useNavigate();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const [name, setName] = useState('John Doe');
    const [role, setRole] = useState('UI Developer');
    const [email, setEmail] = useState('johndoe@example.com');
    const [country, setCountry] = useState('United States');
    const [timezone, setTimezone] = useState('GMT-05:00');
    const [bio, setBio] = useState('I\'m a software developer based in New York, USA. My goal is to solve UI problems with neat CSS without using too much JavaScript.');
    const [skills, setSkills] = useState(['React', 'JavaScript', 'CSS', 'HTML']);

    const handleSave = () => {
        // Implement save functionality
        console.log('Profile saved');
    };

    const handleAddSkill = (skill) => {
        setSkills([...skills, skill]);
    };

    const handleRemoveSkill = (skillToRemove) => {
        setSkills(skills.filter(skill => skill !== skillToRemove));
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
        <ThemeProvider theme={theme}>
            <MotionContainer maxWidth="md" sx={{ py: 4 }}>
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

                        <MotionPaper elevation={3} sx={{ p: 3, mb: 4 }} variants={itemVariants}>
                            <Grid container spacing={3} alignItems="center">
                                <Grid item xs={12} sm={4} sx={{ position: 'relative' }}>
                                    <StyledAvatar alt={name} src="/placeholder.svg" />
                                    <StyledIconButton size="small">
                                        <EditIcon fontSize="small" />
                                    </StyledIconButton>
                                </Grid>
                                <Grid item xs={12} sm={8}>
                                    <TextField
                                        fullWidth
                                        label="Name"
                                        value={name}
                                        onChange={(e) => setName(e.target.value)}
                                        variant="outlined"
                                        margin="normal"
                                    />
                                    <TextField
                                        fullWidth
                                        label="Role"
                                        value={role}
                                        onChange={(e) => setRole(e.target.value)}
                                        variant="outlined"
                                        margin="normal"
                                    />
                                </Grid>
                            </Grid>
                        </MotionPaper>

                        <Grid container spacing={4}>
                            <Grid item xs={12} md={6}>
                                <MotionCard variants={itemVariants}>
                                    <CardContent>
                                        <Typography variant="h6" gutterBottom>Contact Information</Typography>
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
                                        <TextField
                                            fullWidth
                                            label="Country"
                                            value={country}
                                            onChange={(e) => setCountry(e.target.value)}
                                            variant="outlined"
                                            margin="normal"
                                        />
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
                                    </CardContent>
                                </MotionCard>
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <MotionCard variants={itemVariants}>
                                    <CardContent>
                                        <Typography variant="h6" gutterBottom>Bio</Typography>
                                        <TextField
                                            fullWidth
                                            multiline
                                            rows={4}
                                            value={bio}
                                            onChange={(e) => setBio(e.target.value)}
                                            variant="outlined"
                                            margin="normal"
                                        />
                                    </CardContent>
                                </MotionCard>
                            </Grid>
                        </Grid>

                        <MotionCard variants={itemVariants} sx={{ mt: 4 }}>
                            <CardContent>
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
                            </CardContent>
                        </MotionCard>

                        <MotionCard variants={itemVariants} sx={{ mt: 4 }}>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>Portfolio Projects</Typography>
                                <Button
                                    variant="outlined"
                                    component="label"
                                    startIcon={<CloudUploadIcon />}
                                    fullWidth={isMobile}
                                >
                                    Upload Project
                                    <input type="file" hidden />
                                </Button>
                            </CardContent>
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
                            >
                                Save Changes
                            </Button>
                        </Box>
                    </motion.div>
                </AnimatePresence>
            </MotionContainer>
        </ThemeProvider>
    );
};

export default EnhancedUserProfile;