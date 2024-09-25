import React, { useState, useMemo } from 'react';
import { styled, ThemeProvider, createTheme, alpha } from '@mui/material/styles';
import {
    Box,
    AppBar as MuiAppBar,
    Drawer as MuiDrawer,
    Toolbar,
    IconButton,
    Typography,
    Divider,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    Button,
    CssBaseline,
    useMediaQuery,
    Avatar,
    Tooltip,
} from '@mui/material';
import {
    ChevronLeft,
    AccountCircle,
    Brightness4,
    Brightness7,
    ExitToApp,
    Settings,
} from '@mui/icons-material';
import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { mainListItems, secondaryListItems } from './ListItems';

const drawerWidth = 280;

const AppBar = styled(MuiAppBar, {
    shouldForwardProp: (prop) => prop !== 'open',
})(({ theme, open }) => ({
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
    }),
    ...(open && {
        marginLeft: drawerWidth,
        width: `calc(100% - ${drawerWidth}px)`,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen,
        }),
    }),
}));

const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
    ({ theme, open }) => ({
        '& .MuiDrawer-paper': {
            position: 'relative',
            whiteSpace: 'nowrap',
            width: drawerWidth,
            transition: theme.transitions.create('width', {
                easing: theme.transitions.easing.sharp,
                duration: theme.transitions.duration.enteringScreen,
            }),
            boxSizing: 'border-box',
            ...(!open && {
                overflowX: 'hidden',
                transition: theme.transitions.create('width', {
                    easing: theme.transitions.easing.sharp,
                    duration: theme.transitions.duration.leavingScreen,
                }),
                width: theme.spacing(7),
                [theme.breakpoints.up('sm')]: {
                    width: theme.spacing(9),
                },
            }),
        },
    }),
);

const MotionBox = motion(Box);
const MotionIconButton = motion(IconButton);

const StyledListItem = styled(ListItem)(({ theme }) => ({
    borderRadius: theme.shape.borderRadius,
    marginBottom: theme.spacing(0.5),
    '&:hover': {
        backgroundColor: alpha(theme.palette.primary.main, 0.08),
    },
}));

const Logo = styled('img')(({ theme }) => ({
    height: 60,
    marginRight: theme.spacing(2),
    cursor: 'pointer',
    transition: theme.transitions.create(['transform'], {
        duration: theme.transitions.duration.shorter,
    }),
    '&:hover': {
        transform: 'scale(1.05)',
    },
}));

function Copyright(props) {
    return (
        <Typography variant="body2" color="text.secondary" align="center" {...props}>
            {'Copyright © '}
            <Link color="inherit" href="https://mui.com/">
                Your Website
            </Link>{' '}
            {new Date().getFullYear()}
            {'.'}
        </Typography>
    );
}

const EnhancedLayout = () => {
    const [open, setOpen] = useState(false);
    const [darkMode, setDarkMode] = useState(false);
    const [logoutDialogOpen, setLogoutDialogOpen] = useState(false);
    const navigate = useNavigate();
    const location = useLocation();

    const theme = useMemo(
        () =>
            createTheme({
                palette: {
                    mode: darkMode ? 'dark' : 'light',
                    primary: {
                        main: '#3f51b5',
                    },
                    secondary: {
                        main: '#9c27b0',
                    },
                    background: {
                        default: darkMode ? '#121212' : '#f5f5f5',
                        paper: darkMode ? '#1e1e1e' : '#ffffff',
                    },
                    text: {
                        primary: darkMode ? '#ffffff' : '#333333',
                    },
                },
                typography: {
                    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif',
                },
                shape: {
                    borderRadius: 16,
                },
                components: {
                    MuiButton: {
                        styleOverrides: {
                            root: {
                                textTransform: 'none',
                                borderRadius: 12,
                                padding: '10px 20px',
                            },
                        },
                    },
                    MuiDrawer: {
                        styleOverrides: {
                            paper: {
                                backgroundColor: darkMode ? '#1e1e1e' : '#f0f0f0',
                                backgroundImage: 'none',
                            },
                        },
                    },
                    MuiListItemIcon: {
                        styleOverrides: {
                            root: {
                                minWidth: 40,
                                color: darkMode ? alpha('#ffffff', 0.7) : alpha('#000000', 0.7),
                            },
                        },
                    },
                    MuiListItemText: {
                        styleOverrides: {
                            primary: {
                                fontWeight: 500,
                            },
                        },
                    },
                },
            }),
        [darkMode]
    );

    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const toggleDrawer = () => {
        setOpen(!open);
    };

    const toggleDarkMode = () => {
        setDarkMode(!darkMode);
    };

    const getPageTitle = (path) => {
        if (path === "/") {
            return "Dashboard";
        }
        return path.substring(1).charAt(0).toUpperCase() + path.slice(2);
    };

    const handleClick = () => {
        navigate('userprofile');
    };

    const handleLogout = () => {
        setLogoutDialogOpen(true);
    };

    const handleLogoutConfirm = () => {
        setLogoutDialogOpen(false);
        console.log('User logged out');
        navigate('/login');
    };

    const handleLogoutCancel = () => {
        setLogoutDialogOpen(false);
    };

    const headerColor = darkMode ? 'text.primary' : 'rgba(0, 0, 0, 0.87)';

    return (
        <ThemeProvider theme={theme}>
            <Box sx={{ display: 'flex' }}>
                <CssBaseline />
                <AppBar
                    position="absolute"
                    open={open}
                    elevation={0}
                    sx={{
                        backdropFilter: 'blur(20px)',
                        backgroundColor: alpha(theme.palette.background.default, 0.8),
                        boxShadow: darkMode ? 'none' : '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06)',
                    }}
                >
                    <Toolbar sx={{ pr: '24px' }}>
                        <MotionBox
                            initial={{ opacity: 0, x: -20 }}
                            animate={{ opacity: 1, x: 0 }}
                            transition={{ duration: 0.5 }}
                            sx={{ display: 'flex', alignItems: 'center', flexGrow: 1 }}
                        >
                            <Logo src="/images/logo.png" alt="Logo" onClick={toggleDrawer} />
                            <Typography
                                component="h1"
                                variant="h5"
                                noWrap
                                sx={{
                                    fontWeight: 'bold',
                                    color: headerColor,
                                    textShadow: darkMode ? 'none' : '0 1px 2px rgba(0, 0, 0, 0.1)',
                                }}
                            >
                                {getPageTitle(location.pathname)}
                            </Typography>
                        </MotionBox>
                        <MotionIconButton
                            onClick={toggleDarkMode}
                            whileHover={{ scale: 1.1 }}
                            whileTap={{ scale: 0.9 }}
                            sx={{ color: headerColor, ml: 1 }}
                        >
                            {darkMode ? <Brightness7 /> : <Brightness4 />}
                        </MotionIconButton>
                        <Tooltip title="User Profile">
                            <MotionIconButton
                                onClick={handleClick}
                                whileHover={{ scale: 1.1 }}
                                whileTap={{ scale: 0.9 }}
                                sx={{ color: headerColor, ml: 1 }}
                            >
                                <Avatar sx={{ width: 32, height: 32, bgcolor: theme.palette.secondary.main }}>
                                    <AccountCircle />
                                </Avatar>
                            </MotionIconButton>
                        </Tooltip>
                    </Toolbar>
                </AppBar>
                <Drawer variant="permanent" open={open}>
                    <Toolbar
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'flex-end',
                            px: [1],
                        }}
                    >
                        <MotionIconButton onClick={toggleDrawer} whileHover={{ scale: 1.1 }} whileTap={{ scale: 0.9 }}>
                            <ChevronLeft />
                        </MotionIconButton>
                    </Toolbar>
                    <Divider />
                    <Box
                        sx={{
                            height: '100%',
                            display: 'flex',
                            flexDirection: 'column',
                            background: darkMode
                                ? 'linear-gradient(145deg, #1e1e1e 0%, #2d2d2d 100%)'
                                : 'linear-gradient(145deg, #f0f0f0 0%, #ffffff 100%)',
                            pt: 2,
                            overflowX: 'hidden',
                            overflowY: 'auto',
                        }}
                    >
                        <List component="nav" sx={{ px: 2, flexGrow: 1 }}>
                            <AnimatePresence>
                                <MotionBox
                                    initial={{ opacity: 0, y: -20 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    exit={{ opacity: 0, y: -20 }}
                                    transition={{ duration: 0.3 }}
                                >
                                    {mainListItems}
                                </MotionBox>
                            </AnimatePresence>
                            <Divider sx={{ my: 2 }} />
                            <AnimatePresence>
                                <MotionBox
                                    initial={{ opacity: 0, y: -20 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    exit={{ opacity: 0, y: -20 }}
                                    transition={{ duration: 0.3, delay: 0.1 }}
                                >
                                    {secondaryListItems}
                                </MotionBox>
                            </AnimatePresence>
                        </List>
                        <Box sx={{ p: 2 }}>
                            <StyledListItem
                                button
                                onClick={handleLogout}
                                sx={{
                                    background: `linear-gradient(45deg, ${theme.palette.primary.main} 30%, ${theme.palette.secondary.main} 90%)`,
                                    color: 'white',
                                    mb: 1,
                                    '&:hover': {
                                        background: `linear-gradient(45deg, ${theme.palette.primary.dark} 30%, ${theme.palette.secondary.dark} 90%)`,
                                    },
                                }}
                            >
                                <ListItemIcon>
                                    <ExitToApp sx={{ color: 'white' }} />
                                </ListItemIcon>
                                <ListItemText primary="Logout" />
                            </StyledListItem>
                            <StyledListItem button onClick={() => navigate('/settings')}>
                                <ListItemIcon>
                                    <Settings />
                                </ListItemIcon>
                                <ListItemText primary="Settings" />
                            </StyledListItem>
                        </Box>
                    </Box>
                </Drawer>
                <Box
                    component="main"
                    sx={{
                        backgroundColor: 'background.default',
                        flexGrow: 1,
                        height: '100vh',
                        overflow: 'auto',
                        transition: theme.transitions.create('background-color', {
                            duration: theme.transitions.duration.standard,
                        }),
                    }}
                >
                    <Toolbar />
                    <MotionBox
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: 20 }}
                        transition={{ duration: 0.3 }}
                        sx={{ p: 3 }}
                    >
                        <Outlet />
                    </MotionBox>
                    <Copyright sx={{ pt: 4, pb: 4 }} />
                </Box>
            </Box>
            <Dialog
                open={logoutDialogOpen}
                onClose={handleLogoutCancel}
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description"
                PaperProps={{
                    style: {
                        borderRadius: theme.shape.borderRadius,
                        boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
                    },
                }}
            >
                <DialogTitle id="alert-dialog-title">
                    {"Confirm Logout"}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText id="alert-dialog-description">
                        Are you sure you want to log out?
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button
                        onClick={handleLogoutCancel}
                        variant="outlined"
                        sx={{
                            borderColor: theme.palette.primary.main,
                            color: theme.palette.primary.main,
                            '&:hover': {
                                borderColor: theme.palette.primary.dark,
                                backgroundColor: alpha(theme.palette.primary.main, 0.04),
                            }
                        }}
                    >
                        Cancel
                    </Button>
                    <Button
                        onClick={handleLogoutConfirm}
                        variant="contained"
                        autoFocus
                        sx={{
                            background: `linear-gradient(45deg, ${theme.palette.primary.main} 30%, ${theme.palette.secondary.main} 90%)`,
                            color: 'white',
                            '&:hover': {
                                background: `linear-gradient(45deg, ${theme.palette.primary.dark} 30%, ${theme.palette.secondary.dark} 90%)`,
                            }
                        }}
                    >
                        Logout
                    </Button>
                </DialogActions>
            </Dialog>
        </ThemeProvider>
    );
};

export default EnhancedLayout;