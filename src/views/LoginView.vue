<template>
   <div class="page-container loginsyle">
      <div class="login-page-containe">
         <div class="brand" style="width: 300px;">
            <img src="../assets/image/smlogo.png" alt="" class="logo" style="width: 56px;">
            <div class="text">
               <span class="title">
                  慧脉后台管理系统
               </span>
            </div>
         </div>
         <div class="loginform">
            <div class="login-from-wrap" v-show="signlogin === '登录'">
               <div class="login-from-title">
                  <span :class="`ac-number' ${logintext != '账号登录' && 'sm-text'}`" @click="Toggle('账号登录')"
                     style="cursor: pointer;margin-right: 20px;">账号登录</span> <span
                     :class="`ac-number' ${logintext != '手机号登录' && 'sm-text'}`" @click="Toggle('手机号登录')"
                     style="cursor: pointer;">手机号登录</span>
               </div>
               <div v-show="logintext === '账号登录'">
                  <el-form ref="form" :model="userInfo" :rules="rules" class="el-form" style="margin-top:30px">
                     <el-form-item prop="username" class="el-form-item ">
                        <el-input v-model="userInfo.username" placeholder="请输入账号"
                           @blur="userInfo.username = $event.target.value.trim()" />
                     </el-form-item>
                     <div class="input-cipher">
                        <el-form-item prop="password" class="el-form-item">
                           <el-input show-password v-model="userInfo.password" placeholder="请输入密码" style="width: 285px;" />
                           <div class="forgot-password" @click="routers"><span class="vertical">|</span> 忘记密码</div>
                        </el-form-item>

                     </div>
                  </el-form>
               </div>
               <div v-show="logintext === '手机号登录'">
                  <el-form ref="form" :model="userphone" :rules="rules" class="el-form" style="margin-top:30px">
                     <el-form-item prop="phone" class="el-form-item ">
                        <el-input v-model="userphone.phone" placeholder="+86 请输入手机号" />
                     </el-form-item>
                     <div class="input-cipher">
                        <el-form-item prop="smsCode" class="el-form-item">
                           <el-input v-model="userphone.smsCode" placeholder="请输入短信验证码" style="width: 285px;" />
                           <div class="forgot-password">
                              <el-link type="primary" :underline="false" :disabled="duanxindj"
                                 @click="SendPhone('TEMPLATE_ACTION_SIGNIN_ORGANIZATION', userphone, false)"
                                 style="color: #1575ee;font-weight: 400;">{{ dxtxt
                                 }}</el-link>
                           </div>
                        </el-form-item>
                     </div>
                  </el-form>
               </div>
               <div class="login-agree">
                  <el-checkbox v-model="checked"> </el-checkbox>
                  <span style="margin-left: 10px;">我已阅读并同意</span>
                  <router-link :to="{ path: '/userLicense' }"><a class="terms-link">《用户协议与隐私政策》</a></router-link>
               </div>
               <el-button @click="uslogin" class="el-form el-button" :disabled="disabled" v-loading="loading"
                  :element-loading-svg="svg" element-loading-svg-view-box="-10, -10, 50, 50">登录</el-button>
               <div class="signin">
                  没有账号？<span class="signin-text" @click="signlogins('注册')">免费注册</span>
               </div>
            </div>
            <div class="login-from-wrap" v-show="signlogin === '注册'">
               <div class="login-from-title">
                  <span class="ac-number">账号注册</span>
               </div>
               <el-form ref="form" :model="SignUp" :rules="rules" class="el-form" style="margin-top:10px">
                  <el-form-item prop="username" class="el-form-item " :error="usernameMsg">
                     <el-input v-model="SignUp.username" placeholder="请设置账号名（2-20个中英文字符）" @change="useryanz"
                        @blur="SignUp.username = $event.target.value.trim()" />
                  </el-form-item>
                  <el-form-item prop="Password" class="el-form-item">
                     <el-input show-password v-model="SignUp.Password" placeholder="请设置登录密码" />
                  </el-form-item>
                  <el-form-item prop="plainPassword" class="el-form-item">
                     <el-input show-password v-model="SignUp.plainPassword" placeholder="请再次输入密码" />
                  </el-form-item>
                  <el-form-item prop="phone" class="el-form-item" :error="phoneMsg">
                     <el-input v-model="SignUp.phone" placeholder="+86 请输入手机号" @change="phoneyanz" />
                  </el-form-item>
                  <div class="input-cipher">
                     <el-form-item prop="smsCode" class="el-form-item">
                        <el-input v-model="SignUp.smsCode" placeholder="请输入短信验证码" style="width: 285px;" />
                        <div class="forgot-password">
                           <el-link type="primary" :underline="false" :disabled="duanxindj"
                              @click="SendPhone('TEMPLATE_ACTION_SIGNUP_ORGANIZATION', SignUp)"
                              style="color: #1575ee;font-weight: 400;">{{ dxtxt }}</el-link>
                        </div>
                     </el-form-item>
                  </div>
               </el-form>
               <div class="login-agree">
                  <el-checkbox v-model="checked"> </el-checkbox>
                  <span style="margin-left: 10px;">我已阅读并同意</span>
                  <router-link :to="{ path: '/userLicense' }"><a class="terms-link">《用户协议与隐私政策》</a></router-link>
               </div>
               <el-button @click="signs" class="el-form el-button" :disabled="disabled">注册</el-button>
               <div class="signin">
                  已有账号？<span class="signin-text" @click="signlogins('登录')">登录</span>
               </div>
            </div>
         </div>
      </div>
   </div>
</template>
<script setup>
import { reactive, ref } from 'vue';
import { useStore } from 'vuex'
import { useRouter } from 'vue-router';
import { elmessage } from '../utils/popup'
import { ElMessage } from 'element-plus';
import {
   SignUpRequest,
   SendPhoneVerificationCode,
   SignInByUsernameRequest,
   SignInByPhoneRequest,
   CheckOrganizationUsernameExistRequest,
   CheckOrganizationPhoneExistRequest,
   GetSymptomKeyMapRequest
} from '../api/api'
import { setAccessTokenDetails, setRefreshTokenDetails, setTenantId, initSymptomKeyMap } from '../utils/storage'
//设置vuex对象
const $store = useStore();
const usernameMsg = ref('')
const phoneMsg = ref('')
const time = ref(60)
const dxtxt = ref('获取验证码')
const duanxindj = ref(false)
const timesphone = ref(null)
const exist = ref(false)
//设置路由对象
const $router = useRouter();
const loading = ref(false)
const svg = `
        <path class="path" d="
          M 30 15
          L 28 17
          M 25.61 25.61
          A 15 15, 0, 0, 1, 15 30
          A 15 15, 0, 1, 1, 27.99 7.5
          L 15 15
        " style="stroke-width: 4px; fill: rgba(0, 0, 0, 0)"/>
      `
//定义接受用户账号密码
const userInfo = ref({
   username: '',
   password: ''
})
//手机号登录
const userphone = ref({
   phone: '',
   smsCode: '',
   txId: ''
})
//账号注册
const SignUp = ref({
   username: '',
   Password: '',
   plainPassword: '',
   phone: '',
   smsCode: '',
   txId: ''
})
//定义登录方式来判断当前事什么登录方式
const logintext = ref('账号登录')
//控制登陆按钮可用
const checked = ref(false)
const signlogin = ref('登录')

//进入登陆页面先进行退出登陆
const logout = () => {
   $store.dispatch('logout')
}
//定义登录与注册按钮是否可点击防止多次点击调用
const disabled = ref(false)

//验证两次密码是否一致
let validatePass2 = (rule, value, callback) => {
   if (value === "") {
      callback(new Error("请再次输入密码"));
   } else if (value !== SignUp.value.Password) {
      callback(new Error("两次输入密码不一致!"));
   } else {
      callback();
   }
};

//表单验证
const rules = reactive({
   username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
   Password: [{ required: true, message: "请输入密码", trigger: 'blur' }],
   plainPassword: [{ required: true, trigger: 'blur', validator: validatePass2 }]
})

const useryanz = async () => {
   if (SignUp.value.username != '') {
      const data = await CheckOrganizationUsernameExistRequest({ username: SignUp.value.username })
      if (data?.data.exist) {
         exist.value = true
      } else {
         exist.value = false
      }
      if (exist.value == true) {
         usernameMsg.value = '账号已存在'
      } else {
         usernameMsg.value = ''
      }
      console.log(data);
   }
}
const phoneyanz = async () => {
   if (SignUp.value.phone.length != 11) {
      phoneMsg.value = '请输入正确格式的手机号'
   } else {
      if (SignUp.value.phone != '') {
         const data = await CheckOrganizationPhoneExistRequest({  phone: SignUp.value.phone })
         if (data?.data.exist) {
            phoneMsg.value = '手机号已存在'
         } else {
            phoneMsg.value = ''
         }
         console.log(data);
      }
   }
}
//注册点击事件
const signs = async () => {
   if (SignUp.value.username == '' || SignUp.value.Password == '' || SignUp.value.plainPassword == '' || SignUp.value.phone == '' || SignUp.value.smsCode == '') {
      elmessage('请完善账号信息')
   } else if (exist.value == true) {
      elmessage('账号已存在')
   } else if (SignUp.value.Password != SignUp.value.plainPassword) {
      return
   } else {
      const rs = {
         username: SignUp.value.username,
         plainPassword: SignUp.value.plainPassword,
         phone: SignUp.value.phone,
         smsCode: SignUp.value.smsCode,
         txId: SignUp.value.txId
      }
      console.log(rs);
      if (checked.value == false) {
         elmessage('请勾选服务协议条款')
      } else {
         const add = await SignUpRequest(rs)
         console.log(add);
         if (add.status == 200) {
            ElMessage({
               message: '注册成功',
               type: 'success',
            })
            signlogin.value = '登录'
            SignUp.value = {
               username: '',
               Password: '',
               plainPassword: '',
               phone: '',
               smsCode: '',
               txId: ''
            }
         } else {
            elmessage(add.data.detail)
         }
      }
   }
}

//发送注册短信点击事件
const SendPhone = async (tem, phone, bola = true) => {
   duanxindj.value = true
   if (phone.phone != '') {
      if (phone.phone.length != 11) {
         elmessage('请输入正确格式的手机号')
         duanxindj.value = false
      } else {
         if (bola) {
            const data = await CheckOrganizationPhoneExistRequest({  phone: phone.phone })
            if (data?.data.exist) {
               phoneMsg.value = '手机号已存在'
               duanxindj.value = false
            } else {
               console.log(data);
               const rs = {
                  
                  phone: phone.phone,
                  templateAction: tem,
                  language: 'LANGUAGE_SIMPLIFIED_CHINESE'
               }
               console.log(rs);
               await SendPhoneVerificationCode(rs).then((res) => {
                  console.log(res);
                  if (res.status === 200) {
                     phone.txId = res.data.tx_id
                     duanxindj.value = true
                     timesphone.value = setInterval(function () {
                        time.value--;
                        if (time.value > 0) {
                           dxtxt.value = '重新发送' + time.value + 's';
                        } else {
                           time.value = 60;//当减到0时赋值为60
                           dxtxt.value = '获取验证码';
                           clearInterval(timesphone.value);//清除定时器
                           duanxindj.value = false
                        }
                     }, 1000)
                  }
                  else {
                     elmessage(res.data.detail)
                     duanxindj.value = false
                  }
               })
            }
         } else {
            const rs = {
               phone: phone.phone,
               templateAction: tem,
               language: 'LANGUAGE_SIMPLIFIED_CHINESE'
            }
            console.log(rs);
            await SendPhoneVerificationCode(rs).then((res) => {
               console.log(res);
               if (res.status === 200) {
                  phone.txId = res.data.tx_id
                  duanxindj.value = true
                  timesphone.value = setInterval(function () {
                     time.value--;
                     if (time.value > 0) {
                        dxtxt.value = '重新发送' + time.value + 's';
                     } else {
                        time.value = 60;//当减到0时赋值为60
                        dxtxt.value = '获取验证码';
                        clearInterval(timesphone.value);//清除定时器
                        duanxindj.value = false
                     }
                  }, 1000)
               }
               else {
                  elmessage(res.data.detail)
                  duanxindj.value = false
               }
            })
         }
      }
   } else {
      elmessage('请输入手机号')
      duanxindj.value = false
   }

}
async function signinadmin(as) {
   $store.commit('setUserInfo', as?.organization)
   setAccessTokenDetails(as?.access_token)
   setRefreshTokenDetails(as?.refresh_token)
   setTenantId(as?.organization.organization_id)
   const data = await GetSymptomKeyMapRequest()
   console.log(data.data);
   const kkk = true
   if (data.status === 200) {
      initSymptomKeyMap(data?.data, kkk)
   }
   //设置路由信息
   $store.dispatch('getrouter')
   if (as?.organization.has_tenant) {
      $router.push('/general')
   } else {
      $router.push({ name: 'createTenant' })
   }
}
//登陆功能
const uslogin = async () => {
   if (checked.value == false) {
      elmessage('请勾选服务协议条款')
   }
   else if (logintext.value == '账号登录') {
      if (userInfo.value.username == '' || userInfo.value.password == '') {
         elmessage('请输入账号或密码')
      } else {
         loading.value = true
         disabled.value = true
         await SignInByUsernameRequest(userInfo.value)
            .then((add) => {
               console.log(add);
               if (add?.status == 500) {
                  loading.value = false
                  disabled.value = false
                  elmessage(add.data.detail)
               } else {
                  loading.value = false
                  disabled.value = false
                  if (add.data) {
                     signinadmin(add.data)
                  }
               }
            })
      }
   } else if (logintext.value == '手机号登录') {
      disabled.value = true
      console.log(userphone.value);
      await SignInByPhoneRequest(userphone.value)
         .then((add) => {
            console.log(add);
            if (add.status === 500) {
               disabled.value = false
               elmessage(add.data.detail)
            } else {
               disabled.value = false
               if (add.data) {
                  signinadmin(add.data)
               }
            }
         })
   }
}
const routers = () => {
   $router.push({ name: 'findPassword' })
}
const Toggle = (key) => {
   logintext.value = key
}
const signlogins = (key) => {
   signlogin.value = key
}
logout()
</script>
<style lang="scss">
.loginsyle {
   .el-input__inner {
      width: 100%;
      display: block;
      height: 60px;
   }

   // 更改element的默认样式
   .el-input {
      --el-input-focus-border-color: none;
   }

   //更改输入框样式
   .el-input__wrapper {
      box-shadow: none !important;

      &:hover {
         box-shadow: none;
      }


   }

   .is-focus {
      box-shadow: none !important;
   }

   .el-input__wrapper {
      border: none;
      border-radius: 0;
      border-bottom: 1px solid #eee;
      width: 100%;
      height: 60px;
      outline-style: none;
      outline-width: 0;
      padding: 0;
      margin-top: 3px;

      &:hover {
         border-bottom: 1px solid #dcdfe6;
      }
   }

   .el-form {
      .el-form-item {
         margin-bottom: 15px;
      }
   }

   //更改按钮样式
   .el-button {
      width: 100%;
      height: 50px;
      border-radius: 10px;
      margin-top: 10px;
      background-color: #1575ee;
      color: #fff;
      -webkit-transition: .1s;
      transition: .1s;
      text-align: center;
      cursor: pointer;
      outline: none;
      margin-bottom: 10px;
      font-family: Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB, Microsoft YaHei, "\5FAE\8F6F\96C5\9ED1", Arial, sans-serif;
      font-weight: 400;
      font-style: normal;
      font-size: 18px;
   }



   .el-form-item.is-error .el-input__wrapper {
      // 0 0 0 1px var(--el-color-danger) inset
      box-shadow: none !important;
      border-bottom: 1px solid var(--el-color-danger);
   }
}

.page-container {
   width: 100vw;
   height: 100vh;
   background: url('@/assets/image/login_bg_zuzhi.png') no-repeat;
   background-size: cover;

   .login-page-containe {
      width: 100%;
      height: 100%;
      display: flex;

      align-items: center;
      justify-content: space-around;

      .brand .logo {
         float: left;
         height: 30%;
         width: 30%;
         margin-right: 20px;
      }

      .brand .text {
         display: inline;
         color: #fff;

         .title {
            display: block;
            font-size: 26px;
            line-height: 50px;
            font-weight: 400;
         }

         .desc {
            font-size: 14px;
         }
      }

      .loginform {
         position: relative;
         width: 500px;
         height: 680px;
         background-color: #fff;
         border-radius: 10px;

         .login-from-wrap {
            padding: 35px 65px;

            .login-from-title {
               margin-top: 30px;
               font-family: Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB, Microsoft YaHei, "\5FAE\8F6F\96C5\9ED1", Arial, sans-serif;
               font-weight: 700;
               font-style: normal;
               font-size: 30px;
               color: #333;

               .ac-number {
                  margin-right: 30px;
               }
            }

            .input-cipher {
               display: flex;
               align-items: center;
               position: relative;
               border-bottom: 1px solid #eee;

               .el-form-item {
                  margin-bottom: 0;
                  flex: 1;
               }

               .el-input__wrapper {
                  border: none;
               }

               .forgot-password {
                  color: #1575EE;
                  font-size: 13px;
                  font-weight: 400;
                  font-style: normal;
                  text-align: right;
                  cursor: pointer;
                  position: absolute;
                  right: 10px;

                  .vertical {
                     color: #C1C7D0;
                     margin-right: 5px;
                  }
               }
            }

            .forger-tool {
               float: right;
               font-family: Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB, Microsoft YaHei, "\5FAE\8F6F\96C5\9ED1", Arial, sans-serif;
               font-weight: 400;
               font-style: normal;
               font-size: 13px;
               -webkit-transition: .1s;
               transition: .1s;
               color: #1575ee;
               cursor: pointer;
            }

            .login-agree {
               font-family: Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB, Microsoft YaHei, "\5FAE\8F6F\96C5\9ED1", Arial, sans-serif;
               font-weight: 400;
               font-style: normal;
               font-size: 14px;
               margin: 0 auto;
               margin-top: 10%;
               color: #c1c7d0;
               display: -webkit-box;
               display: -ms-flexbox;
               display: flex;
               -webkit-box-pack: center;
               -ms-flex-pack: center;
               justify-content: center;

               .terms-link {
                  color: #1575ee;
                  cursor: pointer;
               }

               .el-checkbox {
                  margin-top: -7px;
               }
            }

            .signin {
               position: absolute;
               bottom: 5%;
               margin-left: 25%;
               color: #c1c7d0;
               font-size: 13px;

               .signin-text {
                  color: #1575ee;
                  cursor: pointer;
               }
            }
         }
      }
   }


}

.sm-text {
   font-size: 20px;
   font-weight: 400;
   color: rgb(153, 153, 153);
}

::v-deep .el-input__wrapper:focus {
   box-shadow: none !important;
}
</style>
