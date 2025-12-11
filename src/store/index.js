import { createStore } from 'vuex';
import VuexPer from 'vuex-persist'
import { admin } from '../router/routers'
import router from '../router';

//使用插件将vuex数据保存到本地
const vuexLocal = new VuexPer({
  storage: window.localStorage
})

//创建vuex对象并导出
export default createStore({
  state: {
    //用户信息
    userinfo: {},
    //路由菜单
    menu: [],
    //首页头部当前页面名称
    hometop: [],
    //当前选菜单栏选中
    activeMenuKey: 0
  },
  mutations: {
    //设置用户信息
    setUserInfo(state, payload) {
      state.userinfo = payload;
    },

    //设置菜单信息
    setMenu(state, payload) {
      state.menu = payload;
    },

    //清除用户相关状态（仅负责修改state）
    clearUserState(state) {
      state.userinfo = {};
      state.menu = [];
      state.hometop = '';
    },

    //设置头部页面信息
    setHometop(state, payload) {
      state.hometop = payload
    },
    setactiveMenuKey(state, payload) {
      state.activeMenuKey = payload
    },
    //设置动态路由相关状态
    setRouterState(state) {
      state.menu = admin
      state.hometop = admin[state.activeMenuKey].meta.title
    }
  },
  actions: {
    //退出登陆操作（处理副作用）
    logout({ commit }) {
      try {
        if (window.localStorage) {
          localStorage.removeItem('refreshToken');
        }
        //调用mutation修改状态
        commit('clearUserState');
        //路由跳转放在action中
        router.push({
          path: '/login',
          replace: true
        });
      } catch (error) {
        console.error('退出登录清理缓存失败：', error);
      }
    },
    //处理动态路由相关逻辑（处理副作用）
    getrouter({ commit, dispatch, state }) {
      //判断是否存储过用户信息,没有则退出到登陆页面
      if (JSON.parse(localStorage.getItem('vuex'))?.userinfo?.name == undefined) {
        //调用action处理路由跳转
        dispatch('logout')
      } else if (JSON.parse(localStorage.getItem('vuex')).userinfo.name != "") {
        //调用mutation修改状态
        commit('setRouterState');
      }
    }
  },
  plugins: [vuexLocal.plugin]
})