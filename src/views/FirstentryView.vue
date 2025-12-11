<template>
    <div class="found-top">
        <div class="found-header steat">
            <img src="https://files.axshare.com/gsc/MH9NX5/b5/c0/40/b5c0405b649748ec89f111341a7b5afe/images/%E6%96%B0%E5%BB%BA%E5%95%86%E6%88%B7%EF%BC%88%E5%90%8D%E7%A7%B0%E3%80%81%E5%9C%B0%E5%9D%80%EF%BC%89/u22.png?pageId=fff809c8-5e20-4968-b536-12afe001a042"
                alt="" class="found-log">
            <span class="header-text">
              慧脉管理系统
            </span>
        </div>
    </div>
   <RouterView/>
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
</style>

