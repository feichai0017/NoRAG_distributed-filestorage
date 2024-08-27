import {request} from "@/utils"


export function loginAPI(formData){
    return  request({
        url:'/api/user/login',
        method:'POST',
        data: formData
    })
}

export function signUpAPI(formData){
    return  request({
        url:'/api/user/signup',
        method:'POST',
        data: formData
    })
}


export function getProfileAPI(){
    return  request({
        url:'/api/user/info',
        method:'GET',
    })
}