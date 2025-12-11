<template>
    <div>
        <logo></logo>
        <div>
            <el-col>
                <img src="" alt="">
                <el-menu :default-active="activeMenuKey" class="el-menu-vertical-demo" :select="selectmenu">
                    <el-submenu>
                        <el-menu-item-group>
                            <el-menu-item v-for="(menu, vm) in menu" :index="vm" :key="menu.name"
                                @click="gorouter(menu, vm)">
                                <div class="itemclass">
                                    <img :src="`${getAssetsFile(icons[vm])}`" alt=""
                                        style="width: 22px;height: 22px;margin-right: 10px;" v-show="activeMenuKey != vm">
                                    <img :src="`${getAssetsFile(iconst[vm])}`" alt=""
                                        style="width: 20px;height: 20px;margin-right: 10px;" v-show="activeMenuKey == vm">
                                    <span>{{ menu.meta.title }}</span>
                                </div>
                            </el-menu-item>
                        </el-menu-item-group>
                    </el-submenu>
                </el-menu>
            </el-col>
        </div>
    </div>
</template>
<script setup>
import { createApp, reactive, ref, toRefs, onBeforeMount } from 'vue'
import logo from '@/components/normalSidebar/Logo.vue'
import { useStore } from 'vuex'
import router from '@/router';

//设置vuex对象
const $store = useStore();
const { menu } = toRefs($store.state)

const { activeMenuKey } = toRefs($store.state)
const activeMenuKeys = ref(JSON.parse(localStorage.getItem('vuex')).activeMenuKey)
console.log(import.meta.url);

//设置未被电击导航图标
const icons = ref(['sidebar_icon_overview.png',
    'sidebar_icon_shop.png',
    'sidebar_icon_commodity.png',
    'sidebar_icon_programme.png',
    'sidebar_icon_bill.png',
])

//设置被点击的导航图标
const iconst = ref(['sidebar_icon_overview_selected.png',
    'sidebar_icon_shop_selected.png',
    'sidebar_icon_commodity_selected.png',
    'sidebar_icon_programme_selected.png',
    'sidebar_icon_bill_selected.png',
])
function getAssetsFile(url) {
    return new URL(`../../assets/image/${url}`, import.meta.url).href;
};
//点击切换页面
const gorouter = (value, i) => {
    $store.commit('setactiveMenuKey', i)
    activeMenuKeys.value = JSON.parse(localStorage.getItem('vuex')).activeMenuKey
    router.push(value.path);
}

const selectmenu = (index, indexPath) => {
    index.click(index)
    gorouter(index.child)
}
onBeforeMount(() => {
    activeMenuKey.value = $store.state.activeMenuKey
})
</script>
<style lang="scss" scoped>
.el-menu-item {
    margin-right: 20px;
    border-radius: 0px 10px 10px 0px;
}

.el-menu-item.is-active {
    background-color: #1575ee !important;
    color: #fff !important;

}

.el-menu--vertical {
    border-right: none !important;
}

.itemclass{
    display: flex;
    align-items: center;
    margin-left: 40px;
}
</style>