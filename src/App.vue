<template>
  <el-config-provider :locale="locale">
    <RouterView></Routerview>
  </el-config-provider>
</template>
<script>
import { ElConfigProvider } from 'element-plus'
import zhCn from "element-plus/dist/locale/zh-cn.mjs";

export default {
  name: 'App',
  components: {
    ElConfigProvider
  },
  setup() {
    return {
      locale: zhCn
    }
  },
  mounted() {
    // 检测浏览器路由改变页面不刷新问题,hash模式的工作原理是hashchange事件
    window.addEventListener('hashchange', () => {
      console.log(window.location.hash, 'window.location.hash')
      let currentPath = window.location.hash.slice(1)
      console.log(currentPath, 'currentPath')
      if (this.$route.path !== currentPath) {
        this.$router.push(currentPath)
      }
    }, false)
  }
}
</script>
