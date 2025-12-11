<template>
    <div class="found-general">
        <div class="found-cont steat">
            <div class="cont-text">
                <span>忘记密码</span>
            </div>
            <div class="mainbody">
                <div class="body-cont">
                    <div class="body-heard">
                        <el-steps :active="active" align-center>
                            <el-step title="验证手机号"></el-step>
                            <el-step title="重置登录密码"></el-step>
                            <el-step title="重置成功"></el-step>
                        </el-steps>
                    </div>
                    <div v-show="active == 1">
                        <div class="body-from">
                            <div class="number">
                                <el-select v-model="countryCode" slot="prepend" placeholder="+86">
                                    <el-option label="+86" value="1"></el-option>
                                </el-select>
                                <el-input placeholder="请输入手机号" v-model="phoneInfo.phoneNumber" class="input-with-select">
                                </el-input>
                            </div>
                            <el-form-item prop="password" class="el-form-item yanzheng">
                                <el-input v-model="phoneInfo.verificationCode" placeholder="请输入验证码" />
                                <el-button type="primary" class="yanzh"
                                    @click="sendVerificationCode(phoneInfo, 'TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD')"
                                    :disabled="isVerificationBtnDisabled">{{ verificationBtnText }}</el-button>
                            </el-form-item>
                            <div class="body-bom">
                                <el-button type="primary" class="complete" @click="verifyAndNextStep">下一步</el-button>
                            </div>
                        </div>
                    </div>
                    <div v-show="active == 2" class="passww">
                        <div class="body-from">
                            <el-form ref="form" :model="passwordInfo" :rules="rules" class="el-form" style="margin-top:10px">
                                <el-form-item prop="newPlainPassword" class="el-form-item">
                                    <el-input show-password v-model="passwordInfo.newPlainPassword" placeholder="请设置登录密码" />
                                </el-form-item>
                                <el-form-item prop="password" class="el-form-item">
                                    <el-input show-password v-model="passwordInfo.password" placeholder="请再次输入密码" />
                                </el-form-item>
                            </el-form>
                            <div class="body-bom">
                                <el-button type="primary" class="complete" @click="resetPassword">完成</el-button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
<script setup>
import { ref, reactive } from "vue";
import { elmessage } from '@/utils/popup'
import { SendPhoneVerificationCode, VerifyPhoneVerificationCodeRequest, ResetOrganizationPasswordRequest } from '@/api/api'
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';

// 定义进度条进度
const active = ref(1)
const $router = useRouter();

// 验证码相关状态
const countdown = ref(60)
const verificationBtnText = ref('获取验证码')
const isVerificationBtnDisabled = ref(false)
const countdownTimer = ref(null)
const countryCode = ref('')

// 密码相关信息
const passwordInfo = ref({
    newPlainPassword: '',
    password: ''
})

// 手机号相关信息
const phoneInfo = ref({
    phoneNumber: '',
    verificationCode: '',
    txId: ''
})

// 验证两次密码是否一致
const validatePass2 = (rule, value, callback) => {
    if (value === "") {
        callback(new Error("请再次输入密码"));
    } else if (value !== passwordInfo.value.newPlainPassword) {
        callback(new Error("两次输入密码不一致!"));
    } else {
        callback();
    }
};

// 表单验证规则
const rules = reactive({
    newPlainPassword: [{ required: true, message: "请输入密码", trigger: 'blur' }],
    password: [{ required: true, trigger: 'blur', validator: validatePass2 }]
})

/**
 * 验证手机号格式是否合规（中国大陆手机号）
 * @param {string} phoneNumber - 手机号
 * @returns {boolean} 是否合规
 */
const validatePhoneNumber = (phoneNumber) => {
    // 中国大陆手机号正则：11位数字，以1开头
    const phoneReg = /^1[3-9]\d{9}$/;
    // 去除可能的空格
    const cleanPhone = phoneNumber.trim();
    
    if (!cleanPhone) {
        elmessage('请输入手机号');
        return false;
    }
    
    if (!phoneReg.test(cleanPhone)) {
        elmessage('请输入有效的中国大陆手机号');
        return false;
    }
    
    return true;
};

// 发送验证码点击事件
const sendVerificationCode = async (data, template) => {
    // 先禁用按钮防止重复点击
    isVerificationBtnDisabled.value = true;
    
    // 1. 先验证手机号格式是否合规
    if (!validatePhoneNumber(data.phoneNumber)) {
        isVerificationBtnDisabled.value = false;
        return;
    }

    try {
        const params = {
            phone: data.phoneNumber.trim(), // 去除空格后传入接口
            templateAction: template,
            language: 'LANGUAGE_SIMPLIFIED_CHINESE'
        };
        
        const res = await SendPhoneVerificationCode(params);
        
        // 只有接口成功返回后才开始倒计时
        if (res.status === 200) {
            data.txId = res.data.tx_id;
            
            // 启动倒计时
            countdownTimer.value = setInterval(() => {
                countdown.value--;
                if (countdown.value > 0) {
                    verificationBtnText.value = `重新发送${countdown.value}s`;
                } else {
                    // 重置倒计时状态
                    countdown.value = 60;
                    verificationBtnText.value = '获取验证码';
                    clearInterval(countdownTimer.value);
                    isVerificationBtnDisabled.value = false;
                }
            }, 1000);
        } else {
            // 接口返回非成功状态时启用按钮
            elmessage(res.data.detail || '验证码发送失败');
            isVerificationBtnDisabled.value = false;
        }
    } catch (error) {
        // 异常情况处理
        elmessage('请求失败，请稍后重试');
        isVerificationBtnDisabled.value = false;
    }
};

// 验证手机短信并进入下一步
const verifyAndNextStep = async () => {
    // 先验证手机号格式
    if (!validatePhoneNumber(phoneInfo.value.phoneNumber)) {
        return;
    }
    
    if (!phoneInfo.value.verificationCode) {
        elmessage('请输入验证码');
        return;
    }

    try {
        const params = {
            phone: phoneInfo.value.phoneNumber.trim(),
            smsCode: phoneInfo.value.verificationCode,
            txId: phoneInfo.value.txId,
            templateAction: 'TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD'
        };
        
        const res = await VerifyPhoneVerificationCodeRequest(params);
        
        if (res.status === 200) {
            console.log('验证通过');
            active.value = 2;
        } else {
            elmessage(res.data.detail);
        }
    } catch (error) {
        elmessage('验证失败，请稍后重试');
    }
};

// 完成密码重置
const resetPassword = async () => {
    // 先验证手机号格式（双重保障）
    if (!validatePhoneNumber(phoneInfo.value.phoneNumber)) {
        return;
    }
    
    if (!passwordInfo.value.newPlainPassword || !passwordInfo.value.password) {
        elmessage('密码不能为空');
        return;
    }

    try {
        const params = {
            phone: phoneInfo.value.phoneNumber.trim(),
            newPlainPassword: passwordInfo.value.newPlainPassword,
            txId: phoneInfo.value.txId
        };
        
        const res = await ResetOrganizationPasswordRequest(params);
        
        if (res.status === 200) {
            ElMessage({
                message: '重置密码成功',
                type: 'success',
            });
            $router.push({ name: 'login' });
        } else {
            elmessage(res.data.detail);
        }
    } catch (error) {
        elmessage('密码重置失败，请稍后重试');
    }
};

// 组件卸载时清除定时器（防止内存泄漏）
import { onUnmounted } from 'vue';
onUnmounted(() => {
    if (countdownTimer.value) {
        clearInterval(countdownTimer.value);
    }
});
</script>

<style lang="scss">

    .body-from {
        width: 380px;
        margin: 50px auto;

        .el-form-item {
            margin-top: 40px;

            .el-input__inner {
                font-size: 16px;
            }
        }

        .body-bom {
            margin: 0 50px;
            width: 380px;
            margin: 70px auto;
            text-align: center;
            padding-bottom: 50px;
            display: flex;
            justify-content: space-around;
            align-items: center;
        }

        .number {
            display: flex;

            .el-tooltip__trigger {
                width: 70px;
            }
        }

        .yanzheng {
            position: relative;

            .yanzh {
                position: absolute;
                right: 10px;
                border-radius: 10px;
                width: 120px;
                height: 40px;
                font-size: 14px;
                font-weight: 400;
                font-style: normal;
                background-color: rgba(21, 117, 238, 1);
            }
        }
    }

    .el-input {
        --el-input-focus-border-color: none;
    }

    .el-input__wrapper {
        box-shadow: none;

        &:hover {
            box-shadow: none;
        }
    }

    .is-focus {
        box-shadow: none;
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

    .found-general {
        background-color: rgba(246, 247, 248, 1);
        height: 100%;
        padding-bottom: 20px;

        .found-cont {
            .cont-text {
                width: 100%;
                height: 8vh;
                line-height: 8vh;
                font-size: 24px;
                font-weight: 400;
                font-style: normal;
                text-align: left;
                border-bottom: 1px solid #333333;
            }

            .mainbody {
                background-color: #fff;
                margin-top: 30px;

                .body-cont {
                    box-sizing: content-box;
                    height: 659px;

                    .body-heard {
                        text-align: center;
                        font-size: 28px;
                        font-weight: 700;
                        padding-top: 70px;

                        .el-step__line {
                            top: 22px;
                        }

                        .mess {
                            font-size: 18px;
                            font-weight: 400;
                        }

                        .messmin {
                            color: #C1C7D0;
                        }
                    }
                }
            }
        }
    }

    .tebleleft {
        font-weight: 400;
        font-style: normal;
        font-size: 14px;
        text-align: right;
        flex: 0.14;
    }

    .complete {
        border-radius: 10px;
        width: 380px;
        height: 50px;
        font-size: 18px;
        background-color: rgba(21, 117, 238, 1);
    }

    .el-form-item {
        margin-bottom: 0px;
    }

    .el-form-item.is-error .el-input__wrapper {
        // 0 0 0 1px var(--el-color-danger) inset
        box-shadow: none;
        border-bottom: 1px solid var(--el-color-danger);
    }
</style>