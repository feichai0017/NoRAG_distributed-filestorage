import {createSlice} from "@reduxjs/toolkit";
import {getToken, removeToken, setToken as _setToken} from "../../utils";
import {getProfileAPI, loginAPI} from "@/api/user.jsx";

const userStore = createSlice({
    name: "user",
    initialState: {
        token: getToken() || '',
        userInfo: {}
    },
    reducers:{
        setToken: (state, action) => {
            state.token = action.payload;
            _setToken(action.payload);
        },
        setUserInfo: (state, action) => {
            state.userInfo = action.payload;
        },
        clearUserInfo: (state) => {
            state.userInfo = {};
            state.token = '';
            removeToken();
        }
    }
});



const {setToken,setUserInfo,clearUserInfo} = userStore.actions



const fetchLogin = (loginForm) => {
    return async (dispatch) => {
        const res = await loginAPI(loginForm);
        dispatch(setToken(res.data.token));
    }
}

const fetchUserInfo = () => {
    return async (dispatch) => {
        const res = await getProfileAPI();
        dispatch(setUserInfo(res.data));
    }
}

export {fetchLogin, fetchUserInfo,clearUserInfo}


const userReducer = userStore.reducer;

export default userReducer;