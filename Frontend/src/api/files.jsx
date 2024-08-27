import {request} from "@/utils"


export function queryAPI(formData){
    return  request({
        url:'/api/file/meta',
        method:'POST',
        data: formData
    })
}

export function queryAllAPI(formData){
    return  request({
        url:'/api/file/meta/query',
        method:'POST',
        data: formData
    })
}


export function uploadAPI(formData){
    return  request({
        url:'/api/file/upload',
        method:'POST',
        data: formData
    })
}


export function downloadAPI(formData){
    return  request({
        url:'/api/file/download',
        method:'POST',
        data: formData
    })
}


export function deleteAPI(formData){
    return  request({
        url:'/api/file/delete',
        method:'POST',
        data: formData
    })
}

export function initMultipartUploadAPI(formData) {
    return request({
        url: '/api/mpupload/init',
        method: 'POST',
        data: formData
    })
}

export function uploadPartAPI(formData) {
    return request({
        url: '/api/mpupload/uploadpart',
        method: 'POST',
        data: formData
    })
}

export function completeMultipartUploadAPI(formData) {
    return request({
        url: '/api/mpupload/complete',
        method: 'POST',
        data: formData
    })
}

export function cancelMultipartUploadAPI(formData) {
    return request({
        url: '/api/mpupload/cancel',
        method: 'POST',
        data: formData
    })
}

export function getMultipartUploadStatusAPI(formData) {
    return request({
        url: '/api/mpupload/status',
        method: 'POST',
        data: formData
    })
}