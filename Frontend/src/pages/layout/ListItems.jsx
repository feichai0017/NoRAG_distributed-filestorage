import * as React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ListSubheader from '@mui/material/ListSubheader';
import { Dashboard, CloudUpload, Search, Folder, Assignment } from '@mui/icons-material';
import { NavLink } from "react-router-dom";

const activeStyle = {
    textDecoration: "none",
    color: "inherit",
    '&.active': {
        backgroundColor: 'rgba(0, 0, 0, 0.08)',
        '& .MuiListItemIcon-root': {
            color: 'inherit',
        }
    },
    transition: 'all 0.3s',
    '&:hover': {
        backgroundColor: 'rgba(0, 0, 0, 0.04)',
        transform: 'translateX(5px)',
    }
};

export const mainListItems = (
    <React.Fragment>
        <ListItemButton component={NavLink} to="/" sx={activeStyle}>
            <ListItemIcon>
                <Dashboard />
            </ListItemIcon>
            <ListItemText primary="Dashboard" />
        </ListItemButton>
        <ListItemButton component={NavLink} to="/knowledge-base" sx={activeStyle}>
            <ListItemIcon>
                <CloudUpload />
            </ListItemIcon>
            <ListItemText primary="Knowledge-Base" />
        </ListItemButton>
        <ListItemButton component={NavLink} to="/queryfile" sx={activeStyle}>
            <ListItemIcon>
                <Search />
            </ListItemIcon>
            <ListItemText primary="Query File" />
        </ListItemButton>
        <ListItemButton component={NavLink} to="/userfiles" sx={activeStyle}>
            <ListItemIcon>
                <Folder />
            </ListItemIcon>
            <ListItemText primary="User Files" />
        </ListItemButton>
    </React.Fragment>
);

export const secondaryListItems = (
    <React.Fragment>
        <ListSubheader component="div" inset>
            Saved reports
        </ListSubheader>
        <ListItemButton sx={activeStyle}>
            <ListItemIcon>
                <Assignment />
            </ListItemIcon>
            <ListItemText primary="Current month" />
        </ListItemButton>
        <ListItemButton sx={activeStyle}>
            <ListItemIcon>
                <Assignment />
            </ListItemIcon>
            <ListItemText primary="Last quarter" />
        </ListItemButton>
        <ListItemButton sx={activeStyle}>
            <ListItemIcon>
                <Assignment />
            </ListItemIcon>
            <ListItemText primary="Year-end sale" />
        </ListItemButton>
    </React.Fragment>
);