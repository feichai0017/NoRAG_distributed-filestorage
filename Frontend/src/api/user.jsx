import {request} from "@/utils"


export function loginAPI(formData){
    return  request({
        url:'/authorizations',
        method:'POST',
        data: formData
    })
}

export function signUpAPI(formData){
    return  request({
        url:'/register',
        method:'POST',
        data: formData
    })
}


export function getProfileAPI(){
    return  request({
        url:'/user/profile',
        method:'GET',
    })
}