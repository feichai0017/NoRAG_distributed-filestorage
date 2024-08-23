import Login from "../pages/login";
import SignUp from "../pages/signup";
import {createBrowserRouter} from "react-router-dom";
import AuthRoute from "../components/AuthRoute";
import Dashboard from "../pages/dashboard";
import Layout from "@/pages/layout/index.jsx";
import Income from "@/pages/income/index.jsx";
import {Payments} from "@mui/icons-material";
import Reports from "@/pages/report/index.jsx";
import UserProfile from "@/pages/userProfile/index.jsx";


const router = createBrowserRouter([
    {
        path: '/',
        element: (
            // <AuthRoute>
                <Layout/>
            // </AuthRoute>
        ),
        children : [
            {
                path: '/',
                index: true,
                element: <Dashboard/>
            },
            {
                path: 'income',
                element: <Income/>
            },
            {
                path:'payment',
                element: <Payments/>
            },
            {
                path:'report',
                element: <Reports/>
            }
        ]
    },
    {
        path: '/login',
        element: <Login/>,
        index: true
    },
    {
        path:'signup',
        element: <SignUp/>
    },
    {
        path: 'userprofile',
        element: <UserProfile/>
    }
]);


export default router;