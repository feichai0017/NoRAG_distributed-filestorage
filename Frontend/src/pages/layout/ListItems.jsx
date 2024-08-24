import * as React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ListSubheader from '@mui/material/ListSubheader';
import DashboardIcon from '@mui/icons-material/Dashboard';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import PaymentIcon from '@mui/icons-material/Payment';
import BarChartIcon from '@mui/icons-material/BarChart';
import AssignmentIcon from '@mui/icons-material/Assignment';
import {NavLink} from "react-router-dom";

const activeStyle = {
    textDecoration: "none",
    color: "inherit",
    '&.active': {
        backgroundColor: 'rgba(0, 0, 0, 0.08)',
        '& .MuiListItemIcon-root': {
            color: 'inherit',
        }
    }
};

export const mainListItems = (
    <React.Fragment>
        <ListItemButton component={NavLink} to="/" sx={activeStyle}>
            <ListItemIcon>
                <DashboardIcon />
            </ListItemIcon>
            <ListItemText primary="Dashboard" />
        </ListItemButton>
        <ListItemButton component={NavLink} to="/upload" sx={activeStyle}>
            <ListItemIcon>
                <AttachMoneyIcon/>
            </ListItemIcon>
            <ListItemText primary="Upload" />
        </ListItemButton>
        <ListItemButton component={NavLink} to="/queryfile" sx={activeStyle}>
            <ListItemIcon>
                <PaymentIcon />
            </ListItemIcon>
            <ListItemText primary="QueryFile" />
        </ListItemButton>
        <ListItemButton component={NavLink} to="/userfiles" sx={activeStyle}>
            <ListItemIcon>
                <BarChartIcon />
            </ListItemIcon>
            <ListItemText primary="UserFiles" />
        </ListItemButton>
    </React.Fragment>
);

export const secondaryListItems = (
    <React.Fragment>
        <ListSubheader component="div" inset>
            Saved reports
        </ListSubheader>
        <ListItemButton>
            <ListItemIcon>
                <AssignmentIcon />
            </ListItemIcon>
            <ListItemText primary="Current month" />
        </ListItemButton>
        <ListItemButton>
            <ListItemIcon>
                <AssignmentIcon />
            </ListItemIcon>
            <ListItemText primary="Last quarter" />
        </ListItemButton>
        <ListItemButton>
            <ListItemIcon>
                <AssignmentIcon />
            </ListItemIcon>
            <ListItemText primary="Year-end sale" />
        </ListItemButton>
    </React.Fragment>
);