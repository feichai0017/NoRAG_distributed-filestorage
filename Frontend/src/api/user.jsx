import { request } from "@/utils";

export function loginAPI(formData) {
    return request({
        url: '/user/login',
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        data: JSON.stringify(formData)
    });
}

export function signUpAPI(formData) {
    return request({
        url: '/user/signup',
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        data: JSON.stringify(formData)
    });
}

export function getProfileAPI() {
    return request({
        url: '/user/info',
        method: 'GET',
        headers: {
            'Accept': 'application/json'
        }
    });
}