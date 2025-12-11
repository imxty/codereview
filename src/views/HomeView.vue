<template>
  <div class="container" v-if="tokenValid">
    <aside class="sidebar-container">
      <sidebar />
    </aside>
    <main class="main-container">
      <indiv :info="topBarLists" />
      <div class="app-main">
        <div class="page-wrap">
          <RouterView />
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onUnmounted, watchEffect } from 'vue'
import { useStore } from 'vuex'
import sidebar from '@/components/normalSidebar/normalsidebar.vue'
import indiv from '@/components/individual.vue'
import { GetSymptomKeyMapRequest } from '@/api/api'
import { refreshAccessToken } from '../utils/refreshAccessToken'
import {
  getRefreshTokenExpiredTimestamp, 
  getAccessTokenExpiredTimestamp, 
  getRefreshToken, 
  initSymptomKeyMap 
} from '../utils/storage'

// 状态管理
const store = useStore()
const tokenValid = ref(false)
const topBarLists = ref(store.state.hometop)
const refreshTimer = ref(null)

// 初始化应用
const initApp = async () => {
  try {
    // 验证用户状态
    const vuexData = JSON.parse(localStorage.getItem('vuex') || '{}')
    if (!vuexData.userinfo?.name || !getRefreshToken()) {
      store.dispatch('logout')
      return
    }

    // 刷新令牌并初始化数据
    await refreshAccessToken()
    store.dispatch('getrouter')
    
    // 获取症状映射表
    const symptomRes = await GetSymptomKeyMapRequest()
    if (symptomRes.status === 200) {
      initSymptomKeyMap(symptomRes?.data, true)
    } else {
      console.warn('获取症状映射表失败', symptomRes)
    }

    tokenValid.value = true
  } catch (error) {
    console.error('应用初始化失败:', error)
    store.dispatch('logout')
  }
}

// 令牌自动刷新逻辑
const setupTokenRefresh = () => {
  // 清理旧定时器
  if (refreshTimer.value) {
    clearInterval(refreshTimer.value)
  }

  refreshTimer.value = setInterval(async () => {
    try {
      const currentTime = Math.trunc(Date.now() / 1000)
      const accessExpiry = getAccessTokenExpiredTimestamp()
      const refreshExpiry = getRefreshTokenExpiredTimestamp()

      // 检查刷新令牌有效性
      if (currentTime >= refreshExpiry) {
        console.log('刷新令牌已过期，退出登录')
        store.dispatch('logout')
        clearInterval(refreshTimer.value)
        return
      }

      // 提前4分钟刷新访问令牌
      if (currentTime > accessExpiry - 240) {
        console.log('自动刷新访问令牌')
        await refreshAccessToken()
      }
    } catch (error) {
      console.error('令牌刷新失败:', error)
    }
  }, 5000)
}

// 监听路由变化时更新顶部导航
watchEffect(() => {
  topBarLists.value = store.state.hometop
})

// 初始化
initApp()
setupTokenRefresh()

// 清理资源
onUnmounted(() => {
  if (refreshTimer.value) {
    clearInterval(refreshTimer.value)
  }
})
</script>

<style lang="scss">
.container {
  display: flex;
}

.sidebar-container {
  width: 270px;
}

.main-container {
  flex: 1 1 0%;
  overflow: hidden;
  min-height: 100%;
  transition: margin-left 0.28s;
  position: relative;
}

.app-main {
  min-height: calc(100vh - 50px);
  width: 100%;
  position: relative;
  padding-top: 80px;
  height: 100vh;
  overflow: auto;
  background-color: #f8f8fa;
}

.page-wrap {
  padding: 30px 40px;
  padding-bottom: 80px;
}
</style>