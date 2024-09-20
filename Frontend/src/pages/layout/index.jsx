import React, { useState, useMemo } from 'react';
import { styled, ThemeProvider, createTheme } from '@mui/material/styles';
import {
    Box,
    AppBar as MuiAppBar,
    Drawer as MuiDrawer,
    Toolbar,
    IconButton,
    Typography,
    Divider,
    List,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    Button,
    CssBaseline,
    useMediaQuery,
} from '@mui/material';
import {
    Menu,
    ChevronLeft,
    AccountCircle,
    Brightness4,
    Brightness7,
    ExitToApp,
} from '@mui/icons-material';
import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { mainListItems, secondaryListItems } from './ListItems';

const drawerWidth = 240;

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
                        main: '#f50057',
                    },
                    background: {
                        default: darkMode ? '#121212' : '#f5f5f5',
                        paper: darkMode ? '#1e1e1e' : '#ffffff',
                    },
                },
                typography: {
                    fontFamily: 'Roboto, Arial, sans-serif',
                },
                shape: {
                    borderRadius: 8,
                },
                components: {
                    MuiButton: {
                        styleOverrides: {
                            root: {
                                textTransform: 'none',
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

    return (
        <ThemeProvider theme={theme}>
            <Box sx={{ display: 'flex' }}>
                <CssBaseline />
                <AppBar position="absolute" open={open} elevation={0}>
                    <Toolbar sx={{ pr: '24px' }}>
                        <IconButton
                            edge="start"
                            color="inherit"
                            aria-label="open drawer"
                            onClick={toggleDrawer}
                            sx={{
                                marginRight: '36px',
                                ...(open && { display: 'none' }),
                            }}
                        >
                            <Menu />
                        </IconButton>
                        <Typography
                            component="h1"
                            variant="h6"
                            color="inherit"
                            noWrap
                            sx={{ flexGrow: 1 }}
                        >
                            {getPageTitle(location.pathname)}
                        </Typography>
                        <IconButton color="inherit" onClick={toggleDarkMode}>
                            {darkMode ? <Brightness7 /> : <Brightness4 />}
                        </IconButton>
                        <IconButton color="inherit" onClick={handleClick}>
                            <AccountCircle />
                        </IconButton>
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
                        <IconButton onClick={toggleDrawer}>
                            <ChevronLeft />
                        </IconButton>
                    </Toolbar>
                    <Divider />
                    <List component="nav" sx={{ flexGrow: 1 }}>
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
                        <Divider sx={{ my: 1 }} />
                        <AnimatePresence>
                            <MotionBox
                                initial={{ opacity: 0, y: -20 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, y: -20 }}
                                transition={{ duration: 0.3 }}
                            >
                                {secondaryListItems}
                            </MotionBox>
                        </AnimatePresence>
                    </List>
                    <Divider />
                    <List>
                        <MotionBox
                            initial={{ opacity: 0, y: -20 }}
                            animate={{ opacity: 1, y: 0 }}
                            exit={{ opacity: 0, y: -20 }}
                            transition={{ duration: 0.3 }}
                        >
                            <Button
                                fullWidth
                                startIcon={<ExitToApp />}
                                onClick={handleLogout}
                                sx={{ justifyContent: 'flex-start', pl: 2 }}
                            >
                                Logout
                            </Button>
                        </MotionBox>
                    </List>
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
                    >
                        <Outlet />
                    </MotionBox>
                    <Copyright sx={{ pt: 4 }} />
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
                    <Button onClick={handleLogoutCancel}>Cancel</Button>
                    <Button onClick={handleLogoutConfirm} autoFocus>
                        Logout
                    </Button>
                </DialogActions>
            </Dialog>
        </ThemeProvider>
    );
};

export default EnhancedLayout;