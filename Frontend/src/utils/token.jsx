
const TOKENKEY = 'token'

function setToken(token){
    localStorage.setItem(TOKENKEY, token)
}

function getToken(){
    return localStorage.getItem(TOKENKEY)
}


function removeToken(){
    localStorage.removeItem(TOKENKEY);
}

export {setToken,getToken,removeToken}