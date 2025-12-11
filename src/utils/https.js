import axios from 'axios'
import jwt from './jwt'
import stort from '../store/index'
import { elmessage } from '../utils/popup'
// import qs from 'qs'

console.log(import.meta.env);
axios.interceptors.request.use(config => {
  // loading
  return config
}, error => {
  return Promise.reject(error)
})

axios.interceptors.response.use(response => {
  return response
}, error => {
  return Promise.resolve(error.response)
})

export default {
  post(url, data) {
    return axios({
      method: 'post',
      baseURL: import.meta.env.VITE_APP_API,//baseURL里面填服务器地址
      url,
      data: data,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json;charset=utf-8',
        Authorization: `Bearer ${jwt.jwt()}`
      }
    }).then(
      (response) => {
        var regexp = /errcode: 1003/g
        var regegq = /missing accesstoken/g
        if (regexp.test(response.data.detail)) {
          stort.dispatch('logout')
          elmessage('登录过期请重新登录')
        }
        if (regegq.test(response.data.detail)) {
          stort.dispatch('logout')
          elmessage('登录过期请重新登录')
        }
        return response
      }
    ).catch(
      (res) => {
        console.log('失败');
        return res
      }
    )
  }
}