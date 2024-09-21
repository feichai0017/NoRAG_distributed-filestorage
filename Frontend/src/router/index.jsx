import Login from "../pages/login";
import SignUp from "../pages/signup";
import {createBrowserRouter} from "react-router-dom";
// import AuthRoute from "../components/AuthRoute";
import Dashboard from "../pages/dashboard";
import Layout from "@/pages/layout/index.jsx";
import KnowledgeBaseManager from "@/pages/knowledgeBase/index.jsx";
import QueryFile from "@/pages/queryfile/index.jsx";
import UserFiles from "@/pages/userfile/index.jsx";
import UserProfile from "@/pages/userProfile/index.jsx";
import Settings from "@/pages/setting/index.jsx";


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
                path: 'knowledge-base',
                element: <KnowledgeBaseManager/>
            },
            {
                path:'queryfile',
                element: <QueryFile/>
            },
            {
                path:'userfiles',
                element: <UserFiles/>
            }
        ]
    },
    {
        path: 'login',
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
    },
    {
        path: 'settings',
        element: <Settings/>
    }
]);


export default router;