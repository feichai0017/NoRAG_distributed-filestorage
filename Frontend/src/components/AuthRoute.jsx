import {getToken} from "../utils";
import {Navigate} from "react-router-dom";


// eslint-disable-next-line react/prop-types
function AuthRoute({children}) {

    const token = getToken();
    if (token) {
        return <>{children}</>
    }else {
        return <Navigate to={"/login"} replace/>;
    }
}


export default AuthRoute;