<template>
  <div  class="sidebar-logo-container">
      <div class="sidebar-logo-link flex-row-lc" key="expand" to="/">
        <img src="../../assets/image/smlogo.png" class="sidebar-logo" />
        <h2 class="sidebar-title">{{ title }}</h2>
      </div>
  </div>
</template>
<script>
export default {
  name: 'SidebarLogo',
  data() {
    return {
      title: '慧脉管理系统'
    }
  }
}
</script>
<style lang="scss" scoped>
.sidebarLogoFade-enter-active {
  transition: opacity 1.5s;
}

.sidebarLogoFade-enter,
.sidebarLogoFade-leave-to {
  opacity: 0;
}

.sidebar-logo-container {
  user-select: none;
  position: relative;
  width: 100%;
  height: 80px;
  line-height: 80px;
  text-align: center;
  overflow: hidden;

  & .sidebar-logo-link {
    height: 100%;
    width: 100%;
    // padding: 20px;

    & .sidebar-logo {
      width: 50px;
      vertical-align: middle;
      margin-right: 12px;
    }

    & .sidebar-title {
      display: inline-block;
      margin: 0;
      line-height: 50px;
      // font-weight: 600;
      // font-size: 20px;

      color: #333;
      vertical-align: middle;
    }
  }

  &.collapse {
    .sidebar-logo {
      margin-right: 0;
    }
  }
}
</style>