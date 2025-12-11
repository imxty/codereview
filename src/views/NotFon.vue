<template>
  <div class="found-top">
    <div class="found-header steat">
      <img
        src="../assets/image/smlogo.png"
        alt="" class="found-log">
      <span class="header-text">
        慧脉管理系统
      </span>
    </div>
  </div>
  <div class="found-data">
    <img src="../assets/image/notfon.jpg" alt="">
    <div style="font-size: 13px;">出错了，您访问的页面不存在…</div>
    <router-link
         :to="{ path: '/' }"> 
         <el-button type="primary" class="fanhuibtn">返回</el-button></router-link>
   
  </div>
</template>
<script setup>

</script>
<style lang="scss">
.steat {
  width: 1000px;
  margin: 0 auto;
}

.found-top {
  height: 10vh;
  line-height: 10vh;

  .found-header {
    display: flex;
    align-items: center;

    .found-log {
      border-width: 0px;
      width: 52px;
      height: 52px;
    }

    .header-text {
      font-weight: 700;
      font-style: normal;
      font-size: 28px;
    }
  }
}

.found-data {
  height: 90vh;
  background-color: rgba(246, 247, 248, 1);
  display: flex;
  flex-direction: column;
  align-items: center;

  .fanhuibtn {
    width: 100px;
    height: 40px;
    background-color: #006be5;
    margin-top: 20px;
  }
}
</style>

