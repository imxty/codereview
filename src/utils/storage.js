const s = JSON.stringify

const p = JSON.parse

// User
const accessTokenKey = 'hmaccessToken'
const refreshTokenKey = 'hmrefreshToken'
const tenantIdKey = 'hmtenantId'

//存储accessToken
export const setAccessTokenDetails = (tokenDetails) =>
    sessionStorage.setItem(accessTokenKey, s(tokenDetails))

//获取accessToken
export const getAccessToken = () => {
    if (sessionStorage.getItem(accessTokenKey) === 'undefined') {
        return ''
    } else {
        return p(sessionStorage.getItem(accessTokenKey))?.token
    }
}

export const getAccessTokenExpiredTimestamp = () =>
    Math.trunc(Date.parse(p(sessionStorage.getItem(accessTokenKey))?.expired_time) / 1000)

export const removeAccessTokenDetails = () =>
    sessionStorage.removeItem(accessTokenKey)


//存储refreshToken
export const setRefreshTokenDetails = (tokenDetails) =>
    localStorage.setItem(refreshTokenKey, s(tokenDetails))

//获取refreshToken
export const getRefreshToken = () => {
    if (localStorage.getItem(refreshTokenKey) === undefined) {
        return ''
    } else {
        return p(localStorage.getItem(refreshTokenKey))?.token
    }
}

export const removeRefreshTokenDetails = () =>
    localStorage.removeItem(refreshTokenKey)

export const getRefreshTokenExpiredTimestamp = () =>
    Math.trunc(Date.parse(p(localStorage.getItem(refreshTokenKey))?.expired_time) / 1000)
//存储tenantId
export const setTenantId = (tenantId) =>
    localStorage.setItem(tenantIdKey, tenantId)

//获取tenantId
export const getTenantId = () => localStorage.getItem(tenantIdKey)

//获取用户信息
export const getuserinfo = () => {
   return p(localStorage.getItem('vuex'))
}

//存储map与key
export const initSymptomKeyMap = (symptomKeyMapObj) => {
    const newlist = []
    console.log(symptomKeyMapObj);
    for (let v in symptomKeyMapObj.dirty_dialectics_key_map) {
        const a = {}
        a[v] = symptomKeyMapObj.dirty_dialectics_key_map[v]
        newlist.push(a)
    }
    for (let v in symptomKeyMapObj.physical_therapy_key_map) {
        const a = {}
        a[v] = symptomKeyMapObj.physical_therapy_key_map[v]
        newlist.push(a)
    }
    for (let v in symptomKeyMapObj.risk_disease_key_map) {
        const a = {}
        a[v] = symptomKeyMapObj.risk_disease_key_map[v]
        newlist.push(a)
    }
    for (let v in symptomKeyMapObj.physique_dialectics_key_map) {
        const a = {}
        a[v] = symptomKeyMapObj.physique_dialectics_key_map[v]
        newlist.push(a)
    }
    const symptomKeyMap = symptomKeyMapObj
    localStorage.setItem('smsymptomKeyMap', s(symptomKeyMap))
    localStorage.setItem('KeyMaplist', s(newlist))
}
export const getSymptomKeyMap = () => p(localStorage.getItem('smsymptomKeyMap'))
export const getKeyMaplist = () => p(localStorage.getItem('KeyMaplist'))
