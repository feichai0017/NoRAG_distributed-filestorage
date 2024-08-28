import {request} from "@/utils"


export function loginAPI(formData){
    return  request({
        url:'/user/login',
        method:'POST',
        data: formData
    })
}

export function signUpAPI(formData){
    return  request({
        url:'/user/signup',
        method:'POST',
        data: formData
    })
}


export function getProfileAPI(){
    return  request({
        url:'/user/info',
        method:'GET',
    })
}