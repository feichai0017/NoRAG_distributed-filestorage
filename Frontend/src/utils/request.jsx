import axios from "axios";
import {getToken, removeToken} from "./token";
import router from "@/router";



const request = axios.create({
    baseURL: "http://192.168.0.200:8081/",
    timeout: 10000,
    withCredentials: true,
});


request.interceptors.request.use((config)=> {
    const token = getToken();
    if (token){
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config
}, (error)=> {
    return Promise.reject(error)
})

// 添加响应拦截器
request.interceptors.response.use((response)=> {
    // 2xx 范围内的状态码都会触发该函数。

    return response.data
}, (error)=> {
    // 超出 2xx 范围的状态码都会触发该函数。

    console.log(error);
    if (error.response && error.response.status === 401) {
        removeToken();
        router.navigate('/login');
        window.location.reload();
    }
    return Promise.reject(error)
})


export {request}