<template>
    <div class="fixed-header">
        <div class="navbar">
            <div class="el-breadcrumb app-breadcrumb breadcrumb-container">
                <el-breadcrumb separator="/">
                    <el-breadcrumb-item v-for="(v, i) in hometop" :to="{ path: v.url }" :key="i">{{
                        v.name }}</el-breadcrumb-item>
                </el-breadcrumb>
            </div>
            <div class="right-menu">
                <div class="el-dropdown">
                    <div class="tenant-field el-dropdown-selfdefine">
                        <el-popover placement="bottom" :width="287" trigger="click">
                            <template #reference>
                                <div class="tenant-info">
                                    <div class="tenant-name">
                                        {{ usinfo.name }}
                                    </div>
                                    <div class="account-expired-time">
                                        组织ID：{{ usinfo.organization_id }}
                                    </div>
                                    <div>
                                        <img src="../assets/icons/button_pulldown.svg" alt="" class="topjt">
                                    </div>
                                </div>
                            </template>
                            <div class="topnamexl">
                                <div class="xldata dianxiugai" @click="touser">
                                    <div class="namexlleft ">账号管理</div>
                                    <div class="namexlright ">
                                        {{ usinfo.name }}<span><img src="../assets/icons/u82.svg" alt=""
                                                style="width: 6px; height: 7px;margin-left: 20px;"></span>
                                    </div>
                                </div>
                                <div class="xldata dianxiugai" @click="loingout">
                                    <div class="namexlleft ">退出登录</div>
                                    <div class="namexlright "><img src="../assets/image/u78.png" alt=""
                                            style="width: 14px; height: 14px;"></div>
                                </div>
                            </div>
                        </el-popover>
                    </div>
                </div>
                <div class="el-dropdown">
                    <div class="notification-field el-dropdown-selfdefine">
                        <div class="el-badge">
                            <span>
                                <span class="caution" v-if="listReview.review_notifications"></span>
                                <el-popover :width="600"
                                    popper-style="box-shadow: rgb(14 18 22 / 35%) 0px 10px 38px -10px, rgb(14 18 22 / 20%) 0px 10px 20px -15px; padding: 20px;">
                                    <template #reference>
                                        <img src="@/assets/image/icon_topbar_news.png" alt=""
                                            style="width: 30px;height: 30px;">
                                    </template>
                                    <template #default>
                                        <div class="demo-rich-conent"
                                            style="display: flex; gap: 16px; flex-direction: column">
                                            <div>
                                                <div class="el-scrollbar" style="height: 320px;">
                                                    <div class="el-scrollbar__wrap"
                                                        style="margin-bottom: -15px; margin-right: -15px;">
                                                        <div class="el-scrollbar__view">
                                                            <div v-for="v in listReview.review_notifications" :key="v"
                                                                class="el-dropdown-menu__item">
                                                                <div class="notification-item-container">
                                                                    <div class="notification-item-content">
                                                                        <div class="notification-title">
                                                                            <span class="notification-dot not-pass"
                                                                                :style="{ background: v.pass ? '#3ca1f8' : '#f56c6c' }"></span>
                                                                            <span>{{ v.notification_title }}  <span
                                                                                v-if="v.review_type == 'REVIEW_TYPE_PRODUCT'">
                                                                                <span v-if="!v.pass">审核失败</span>
                                                                                <span v-if="v.pass">审核通过</span>
                                                                            </span></span>
                                                                           
                                                                        </div>
                                                                        <div class="shyj" v-if="!v.pass">
                                                                            审核意见：{{ v.review_comment }}
                                                                        </div>
                                                                        <div class="notification-date">{{
                                                                            fermitTime(v.created_time) }}</div>
                                                                    </div>
                                                                </div>
                                                                <el-button type="primary" plain
                                                                    @click="ConfirmReviewNotification(v.notification_id)">知道了</el-button>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </template>
                                </el-popover>
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { createApp, reactive, ref, toRefs, defineProps } from 'vue'
import { useStore } from 'vuex'
import { getTenantId } from '@/utils/storage'
import { fermitTime } from '@/utils/fermitTime'
import { useRouter } from 'vue-router';
import { ListReviewNotificationsRequest, ConfirmReviewNotificationRequest } from '@/api/api'

const $router = useRouter();

//设置vuex对象
const $store = useStore();

//获取头部信息
const { hometop } = toRefs($store.state)
//获取用户信息
const usinfo = ref($store.state.userinfo)
//审核信息
const listReview = ref([])

const ListReviewNotifications = async () => {
    const data = await ListReviewNotificationsRequest({ organizationId: getTenantId() })
    listReview.value = data.data
    console.log(listReview.value);
}

const touser = () => {
    $router.push({
        name: 'accountInformation',
        state: {
            keyword: JSON.stringify(usinfo.value)
        }
    })
}
const ConfirmReviewNotification = async (id) => {
    const data = await ConfirmReviewNotificationRequest({ notificationId: id })
    console.log(data);
    await ListReviewNotifications()
}
const loingout = () => {
    $store.dispatch('logout')
}
ListReviewNotifications()
</script>
<style lang="scss">
.fixed-header {
    position: absolute;
    top: 0;
    right: 0;
    z-index: 1000;
    width: 100%;
    height: 80px;
    -webkit-transition: width .28s;
    transition: width .28s;

    .navbar {
        height: 100%;
        overflow: hidden;
        position: relative;
        background: #fff;
        border-bottom: 1px solid #eee;
        display: flex;
        -webkit-box-pack: justify;
        -ms-flex-pack: justify;
        justify-content: space-between;
        align-items: center;
        font-size: 18px;

        .app-breadcrumb.el-breadcrumb {
            font-size: 14px;
            line-height: 80px;
            margin-left: 40px;
        }

        .right-menu {
            height: 100%;
            width: 370px;
            display: flex;
            cursor: pointer;
            user-select: none;
            justify-content: flex-end;

            .el-dropdown {
                display: inline-block;
                position: relative;
                color: #666;
                font-size: 14px;

                .tenant-field {
                    height: 100%;
                    width: 210px;
                    border-right: 1px solid #eee;
                    border-left: 1px solid #eee;
                    display: -webkit-box;
                    display: -ms-flexbox;
                    display: flex;
                    -webkit-box-pack: space-evenly;
                    -ms-flex-pack: space-evenly;
                    justify-content: space-evenly;
                    -webkit-box-align: center;
                    -ms-flex-align: center;
                    align-items: center;

                    .el-avatar {
                        display: inline-block;
                        -webkit-box-sizing: border-box;
                        box-sizing: border-box;
                        text-align: center;
                        overflow: hidden;
                        color: #fff;
                        background: #c0c4cc;
                        width: 40px;
                        height: 40px;
                        line-height: 40px;
                        font-size: 14px;

                        img {
                            display: block;
                            height: 100%;
                            min-width: 100%;
                        }
                    }

                    .el-avatar--circle {
                        border-radius: 50%;
                    }

                    .tenant-info {
                        font-family: Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB, Microsoft YaHei, "\5FAE\8F6F\96C5\9ED1", Arial, sans-serif;
                        font-weight: 400;
                        font-style: normal;

                        .tenant-name {
                            width: 150px;
                            overflow: hidden;
                            white-space: nowrap;
                            text-overflow: ellipsis;
                            font-size: 18px;
                            color: #333;
                        }

                        .account-expired-time {
                            font-size: 12px;
                            color: #3ca1f8;
                            margin-top: 5px;
                        }

                    }
                }

                .notification-field {
                    height: 100%;
                    width: 80px;
                    display: -webkit-box;
                    display: -ms-flexbox;
                    display: flex;
                    -webkit-box-pack: center;
                    -ms-flex-pack: center;
                    justify-content: center;
                    -webkit-box-align: center;
                    -ms-flex-align: center;
                    align-items: center;

                    .el-badge {
                        position: relative;
                        vertical-align: middle;
                        display: inline-block;
                    }
                }
            }
        }
    }
}

.el-dropdown-menu {
    position: absolute;
    top: 0;
    left: 0;
    z-index: 10;
    padding: 10px 0;
    margin: 5px 0;
    background-color: #fff;
    border: 1px solid #ebeef5;
    border-radius: 4px;
    -webkit-box-shadow: 0 2px 12px 0 rgba(0, 0, 0, .1);
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, .1);
}

.account-menu.el-dropdown-menu {
    border-radius: 10px;
}

.el-dropdown-menu__item {
    list-style: none;
    line-height: 18px;
    padding: 10px 20px;
    margin: 0;
    font-size: 14px;
    color: #666;
    cursor: pointer;
    outline: none;
    justify-content: space-between;
    width: 100%;
}

.account-menu.el-dropdown-menu .el-dropdown-menu__item {
    width: 290px;
    display: flex;
    justify-content: space-between;
}

.el-menu--vertical {
    border-right: none !important;
}

.el-menu-item-group__title {
    padding: 0 !important;
}

.el-scrollbar {
    overflow: hidden;
    position: relative;
}

.el-scrollbar__wrap {
    overflow: scroll;
    height: 100%;
    overflow-x: hidden;
}

.el-dropdown-menu__item {
    list-style: none;
    line-height: 18px;
    padding: 10px 20px;
    margin: 0;
    font-size: 14px;
    color: #666;
    cursor: pointer;
    outline: none;
}

.notifications-menu.el-dropdown-menu .el-dropdown-menu__item {
    line-height: normal;
}

.notifications-menu.el-dropdown-menu .notification-item-container {
    margin-bottom: 10px;
    border-bottom: 1px solid #ebeef5;
    display: flex;
    -webkit-box-pack: justify;
    -ms-flex-pack: justify;
    justify-content: space-between;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    height: 100px;
}

.notifications-menu.el-dropdown-menu .notification-item-container .notification-item-content {
    color: #666;
}

.notification-title {
    font-family: "Helvetica Neue", "Helvetica", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "\5FAE\8F6F\96C5\9ED1", "Arial", sans-serif;
    font-weight: 400;
    font-style: normal;
    font-size: 16px;
    margin-bottom: 10px;
    color: #333;
}

.notification-dot {
    height: 10px;
    width: 10px;
    border-radius: 50%;
    display: inline-block;
    margin-right: 10px;
}

.notification-date {
    font-family: "Helvetica Neue", "Helvetica", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "\5FAE\8F6F\96C5\9ED1", "Arial", sans-serif;
    font-weight: 400;
    font-style: normal;
    font-size: 12px;
    color: #c1c7d0;
}

.caution {
    display: block;
    width: 11px;
    height: 11px;
    background-color: #F56C6C;
    position: absolute;
    border-radius: 10px;
    left: 25px;
    top: -5px;
}

.xldata {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    padding: 0 20px;
    cursor: pointer;
}

.dianxiugai:hover {
    background-color: rgba(246, 247, 248, 1);
}

.el-breadcrumb__item {
    font-size: 18px;
    color: #666666;
    font-weight: 700;

    .el-breadcrumb__inner,
    i.s-link {
        font-weight: 700;
    }
}

.el-breadcrumb__inner {
    font-weight: 700;
}

.lienametxt {
    display: flex;
    align-items: center;
}

.topjt {
    width: 12px;
    position: absolute;
    right: 15px;
    top: 35px;
}
</style>