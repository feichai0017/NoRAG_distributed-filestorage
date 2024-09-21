import * as React from 'react';
import { styled } from '@mui/material/styles';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ListSubheader from '@mui/material/ListSubheader';
import { Dashboard, CloudUpload, Search, Folder, Assignment } from '@mui/icons-material';
import { NavLink } from "react-router-dom";
import { motion } from "framer-motion";

const StyledListItemButton = styled(motion(ListItemButton))(({ theme }) => ({
    borderRadius: theme.shape.borderRadius,
    marginBottom: theme.spacing(0.5),
    transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
    '&.active': {
        backgroundColor: theme.palette.mode === 'dark'
            ? 'rgba(255, 255, 255, 0.08)'
            : 'rgba(0, 0, 0, 0.08)',
        '& .MuiListItemIcon-root': {
            color: theme.palette.primary.main,
        },
        '& .MuiListItemText-primary': {
            fontWeight: 600,
            color: theme.palette.mode === 'dark'
                ? theme.palette.primary.light
                : theme.palette.primary.main,
        },
    },
    '&:hover': {
        backgroundColor: theme.palette.mode === 'dark'
            ? 'rgba(255, 255, 255, 0.05)'
            : 'rgba(0, 0, 0, 0.04)',
    },
}));

const StyledListItemIcon = styled(ListItemIcon)(({ theme }) => ({
    minWidth: 40,
    color: theme.palette.mode === 'dark'
        ? theme.palette.text.secondary
        : theme.palette.text.primary,
}));

const StyledListSubheader = styled(ListSubheader)(({ theme }) => ({
    backgroundColor: 'transparent',
    color: theme.palette.mode === 'dark'
        ? theme.palette.text.secondary
        : theme.palette.text.primary,
    fontWeight: 700,
    fontSize: '0.75rem',
    letterSpacing: '0.08em',
    textTransform: 'uppercase',
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(1),
}));

const MotionListItemButton = ({ to, icon, primary }) => (
    <StyledListItemButton
        component={NavLink}
        to={to}
        whileHover={{ x: 5 }}
        whileTap={{ scale: 0.95 }}
    >
        <StyledListItemIcon>
            {icon}
        </StyledListItemIcon>
        <ListItemText
            primary={primary}
            primaryTypographyProps={{
                style: { fontWeight: 500 }
            }}
        />
    </StyledListItemButton>
);

export const mainListItems = (
    <React.Fragment>
        <MotionListItemButton to="/" icon={<Dashboard />} primary="Dashboard" />
        <MotionListItemButton to="/knowledge-base" icon={<CloudUpload />} primary="Knowledge Base" />
        <MotionListItemButton to="/queryfile" icon={<Search />} primary="Query File" />
        <MotionListItemButton to="/userfiles" icon={<Folder />} primary="User Files" />
    </React.Fragment>
);

export const secondaryListItems = (
    <React.Fragment>
        <StyledListSubheader component="div" inset>
            Saved reports
        </StyledListSubheader>
        <MotionListItemButton to="/reports/current" icon={<Assignment />} primary="Current month" />
        <MotionListItemButton to="/reports/last-quarter" icon={<Assignment />} primary="Last quarter" />
        <MotionListItemButton to="/reports/year-end" icon={<Assignment />} primary="Year-end sale" />
    </React.Fragment>
);