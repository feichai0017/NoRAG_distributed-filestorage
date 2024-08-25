import {request} from "@/utils"


export function queryAPI(formData){
    return  request({
        url:'/api/meta',
        method:'POST',
        data: formData
    })
}

export function queryAllAPI(formData){
    return  request({
        url:'/api/meta/query',
        method:'POST',
        data: formData
    })
}


export function uploadAPI(formData){
    return  request({
        url:'/api/upload',
        method:'POST',
        data: formData
    })
}


export function downloadAPI(formData){
    return  request({
        url:'/api/download',
        method:'POST',
        data: formData
    })
}


export function deleteAPI(formData){
    return  request({
        url:'/api/delete',
        method:'POST',
        data: formData
    })
}