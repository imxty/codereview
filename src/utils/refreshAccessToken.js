import {RefreshAccessTokenRequest} from '../api/api'
import { setAccessTokenDetails, getTenantId, getRefreshToken } from '../utils/storage'
import stort from '../store/index'
import { elmessage } from '../utils/popup'

//刷新token并保存
export async function refreshAccessToken() {
    const refreshToken = getRefreshToken()
    if (refreshToken) {
        const accessToken  = await RefreshAccessTokenRequest({
          tenantId: getTenantId(),
          refreshToken: refreshToken
        })
        console.log(stort);
        setAccessTokenDetails(accessToken.data.access_token)
        return accessToken
    } else {
      return null
    }
  }