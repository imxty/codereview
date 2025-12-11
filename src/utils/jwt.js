import { Base64 } from "js-base64";
import { KEYUTIL, KJUR, hextob64 } from 'jsrsasign';
import { plaintextkeytest, plaintextkeydev, plaintextkeystag,plaintextkeyproduction } from "../../Plaintextkey"
import { getAccessToken, getAccessTokenExpiredTimestamp } from "./storage"
// import jose from "jose";

let apppem = {}
if (import.meta.env.MODE == 'testing') {
    apppem = plaintextkeytest
}
if (import.meta.env.MODE == 'development') {
    apppem = plaintextkeydev
}
if (import.meta.env.MODE == 'staging') {
    apppem = plaintextkeystag
}

if (import.meta.env.MODE == 'prod') {
    apppem = plaintextkeyproduction
}

//替换特殊符号
function shimBase64URL(s) {
    let r = s.replaceAll('+', '-')
    r = r.replaceAll('/', '_')
    r = r.replaceAll('=', '')
    return r
}

function encodeBase64(v) {
    return Base64.encodeURL(JSON.stringify(v))
}
/** 加密 */
function encrypt(txt) {
    // 方式1: 先建立 key 对象, 构建 signature 实例, 传入 key 初始化 -> 签名
    const key = KEYUTIL.getKey(apppem.key);
    // 创建 Signature 对象，设置签名编码算法
    const signature = new KJUR.crypto.Signature({
        alg: 'SHA256withRSA'
    });
    // 传入key实例, 初始化signature实例
    signature.init(key);
    // 传入待加密字符串
    signature.updateString(txt);
    // 签名, 得到16进制字符结果
    let a = signature.sign();
    let sign = hextob64(a);
    return shimBase64URL(sign);
}


export default {
    jwt() {
        //获取当前时间戳
        const getdate = Math.round(new Date().getTime() / 1000);

        let header = {
            "alg": "RS256",
            "kid": apppem.kid,
            "typ": "JWT"
        }

        let payload = {
            "exp": getdate + 270,
            "iat": getdate - 30,
            "iss": "JinmuHealth",
            "sub": apppem.app_id,
            "access_token": getAccessToken()
        }
        let base64Token = encodeBase64(header) + '.' + encodeBase64(payload)
        let signatureBase64 = encrypt(base64Token)
        let jwt = base64Token + '.' + signatureBase64
        return jwt

    }, isJSON(str) {
        if (typeof str === 'string') {
            try {
                JSON.parse(str)
                return true
            } catch (e) {
                return false
            }
        } else if (typeof str === 'object') {
            try {
                return true
            } catch (e) {
                return false
            }
        }
    }
}

