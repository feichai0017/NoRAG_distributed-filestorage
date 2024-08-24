import {request} from "@/utils"


export function loginAPI(formData){
    return  request({
        url:'/api/login',
        method:'POST',
        data: formData
    })
}

export function signUpAPI(formData){
    return  request({
        url:'/api/signup',
        method:'POST',
        data: formData
    })
}


export function getProfileAPI(){
    return  request({
        url:'/api/info',
        method:'GET',
    })
}