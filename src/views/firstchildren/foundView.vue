<template>
    <div class="founvie">
        <div class="found-general">
            <div class="found-cont steat">
                <div class="cont-text">
                    <span>新建商户</span>
                </div>
                <div class="mainbody">
                    <div class="body-cont">
                        <div class="body-heard" v-show="found == false">
                            <div class="messon">创建您的第一个商户</div>
                            <div class="mess">请填写您的商户信息</div>
                        </div>
                        <div class="body-heard" v-show="found">
                            <div class="mess">还需要邀请其他商户吗？</div>
                            <div class="messmin txtsm">您可以通过点击“复制链接”按钮，将邀请链接发送给他人</div>
                        </div>
                        <div v-show="conttext === '地址'">
                            <div class="body-from">
                                <div class="body-table">
                                    <div class="table-img">
                                        <div class="tebleleft">商户图标:</div>
                                        <div class="tableright">
                                            <imageUploader @imageBase64="logoImageBase64"
                                                @onImageUploaded="logoImageUploaded">
                                            </imageUploader>
                                            <div class="imgtext">
                                                <p>请上传商户门头照片/商户logo；</p>

                                                <p>图片不允许涉及政治敏感与色情;</p>

                                                <p>图片格式只支持：BMP、JPEG、JPG、GIF、PNG，大小不超过2M</p>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="table-input">
                                        <div class="tebleleft">商户名称:</div>
                                        <div class="tableright">
                                            <el-form-item prop="account" class="el-form-item ">
                                                <el-input placeholder="请输入商户名称" v-model="certifyForm.name" />
                                            </el-form-item>
                                        </div>
                                    </div>
                                    <div class="table-input">
                                        <div class="tebleleft">商户地址:</div>
                                        <div class="tableright">
                                            <el-select v-model="certifyForm.address.province" placeholder="请选择"
                                                @change="checksheng">
                                                <el-option v-for="item in provinceOptions()" :key="item.value"
                                                    :label="item.label" :value="item.value">
                                                </el-option>
                                            </el-select>
                                            <el-select v-model="certifyForm.address.city" placeholder="请选择"
                                                @change="checkshi" style="margin-left: 6px;">
                                                <el-option v-for="item in cityOptions()" :key="item.value"
                                                    :label="item.label" :value="item.value">
                                                </el-option>
                                            </el-select>
                                            <el-select v-model="certifyForm.address.district" placeholder="请选择"
                                                @change="checkxian" style="margin-left: 5px;">
                                                <el-option v-for="item in regionOptions()" :key="item.value"
                                                    :label="item.label" :value="item.value">
                                                </el-option>
                                            </el-select>
                                            <el-form-item prop="account" class="el-form-item "
                                                style="margin-top: 10px;">
                                                <el-input placeholder="具体地址" v-model="certifyForm.address.street" />
                                            </el-form-item>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div v-show="conttext === '营业执照'">
                            <div class="body-from">
                                <div class="body-table">
                                    <div class="table-img">
                                        <div class="tebleleft">营业执照图片:</div>
                                        <div class="tableright">
                                            <imageUploader @imageBase64="yingyeImageBase64"
                                                @onImageUploaded="yingyeImageUploaded"></imageUploader>
                                            <div class="imgtext">
                                                <p>请上传清晰的营业执照图片，保证以上商户信息与营业执照一致；</p>

                                                <p>图片不允许涉及政治敏感与色情;</p>

                                                <p>图片格式只支持：BMP、JPEG、JPG、GIF、PNG，大小不超过2M</p>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="table-input">
                                        <div class="tebleleft">社会信用代码:</div>
                                        <div class="tableright">
                                            <el-form-item prop="account" class="el-form-item ">
                                                <el-input v-model="certifyForm.socialCreditCode" />
                                            </el-form-item>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div v-show="conttext === '联系人'">
                            <div class="body-from">
                                <div class="body-table">
                                    <div class="table-input">
                                        <div class="tebleleft">联系人:</div>
                                        <div class="tableright">
                                            <el-form-item prop="account" class="el-form-item ">
                                                <el-input placeholder="请输入商户联系人姓名" v-model="certifyForm.contactName" />
                                            </el-form-item>
                                        </div>
                                    </div>
                                    <div class="table-input">
                                        <div class="tebleleft">联系人手机号:</div>
                                        <div class="tableright">
                                            <el-form-item prop="account" class="el-form-item ">
                                                <el-input placeholder="请输入联系人手机号，用于商户登录" v-model="userphone.phone" />
                                            </el-form-item>
                                        </div>
                                    </div>
                                    <div class="table-input">
                                        <div class="tebleleft">短信验证码:</div>
                                        <div class="tableright">
                                            <el-form-item prop="account" class="el-form-item ">
                                                <el-input placeholder="请输入验证码" v-model="userphone.smsCode"
                                                    style="width: 72%;" />
                                                <el-button type="primary"
                                                    @click="SendPhone(userphone, 'TEMPLATE_ACTION_BIND_TENANT_PHONE')"
                                                    :disabled="duanxindj"
                                                    style="width: 25%;margin-left: 3%;background-color: #3773e6;border: none;">{{
                                                        dxtxt }}</el-button>
                                            </el-form-item>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div v-show="conttext === '商户'">
                            <div class="body-from">
                                <div class="body-table">
                                    <el-form ref="form" :model="userInfo" :rules="rules" class="el-form">
                                        <div class="table-input">
                                            <div class="tebleleft">商户登录密码:</div>
                                            <div class="tableright">
                                                <el-form-item prop="newPlainPassword" class="el-form-item ">
                                                    <el-input show-password placeholder="请设置商户登录密码"
                                                        v-model="userInfo.newPlainPassword" />
                                                </el-form-item>
                                            </div>
                                        </div>
                                        <div class="table-input">
                                            <div class="tebleleft">确认密码:</div>
                                            <div class="tableright">
                                                <el-form-item prop="password" class="el-form-item ">
                                                    <el-input show-password placeholder="请再次输入商户登录密码"
                                                        v-model="userInfo.password" />
                                                </el-form-item>
                                            </div>
                                        </div>
                                    </el-form>
                                </div>
                            </div>
                        </div>
                        <div v-show="conttext === '添加其他'">
                            <div class="body-from">
                                <div class="body-table">
                                    <div class="body-link">
                                        <el-form-item prop="account" class="el-form-item linkinp">
                                            <el-input v-model="jmurl" />
                                        </el-form-item>
                                        <el-button type="primary" class="complete"
                                            @click="copyText(false)">复制链接</el-button>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="body-bom" v-show="found == false">
                            <span @click="push" style="cursor:pointer">跳过</span>
                            <el-button type="primary" @click="Toggle(conttext)"
                                style="background-color: #3773e6;border: none;width: 150px;height: 40px;">下一步</el-button>
                        </div>
                        <div class="body-bom" v-show="found">
                            <el-button type="primary" class="complete">进入后台</el-button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
<script setup>
import { ref, reactive } from 'vue';
import imageUploader from '@/components/common/imageUploader.vue';
import { getAreaOptions, TextToCode } from '@/utils/areadDataOptions'
import { SendPhoneVerificationCode, VerifyPhoneVerificationCodeRequest, CreateTenantRequest, UploadImageRequest } from '@/api/api'
import { elmessage } from '@/utils/popup'
import { getTenantId, getuserinfo } from '@/utils/storage'
import { ElMessage } from 'element-plus';
import useClipboard from 'vue-clipboard3'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router';

const $store = useStore();

const { toClipboard } = useClipboard()

const userInfo = ref({
    newPlainPassword: '',
    password: ''
})
const $router = useRouter();

const time = ref(60)
const dxtxt = ref('获取验证码')
const duanxindj = ref(false)
const timesphone = ref(null)
const jmurl = ref('https://res.jinmuhealth.com' + import.meta.env.VITE_APP_PUBLIC_PATH + 'index.html#/lnvite' + '?o=' + getTenantId() + '&on=' + encodeURI(getuserinfo()?.userinfo.name))
console.log(getuserinfo());
const userphone = ref({

    phone: '',
    smsCode: '',
    txId: '',
    templateAction: 'TEMPLATE_ACTION_BIND_TENANT_PHONE'
})
const certifyForm = ref({
    tenantId: null,
    logo: {
        mime: '',
        image: '',
        filename: ''
    },
    name: '',
    contactName: '',
    address: {
        province: '',
        city: '',
        district: '',
        street: ''
    },

    contactPhone: '',
    businessLicense: {
        mime: '',
        image: '',
        filename: ''
    },
    socialCreditCode: '',
    logoUrl: '',
    businessLicenseUrl: '',
    safePhone: ''
})

//处理上传图片
const logoImageBase64 = async (val) => {
    console.log(val);
    const imageInfo = val.split(/data:(.*?);base64,(.*?)/g).filter(Boolean)
    const [mime, image] = imageInfo
    certifyForm.value.logo.mime = mime
    certifyForm.value.logo.image = image
    const data = await UploadImageRequest({ image: certifyForm.value.logo })
    console.log(data);
    certifyForm.value.logoUrl = data.data.image_url
}
const yingyeImageBase64 = async (val) => {
    console.log(val);
    const imageInfo = val.split(/data:(.*?);base64,(.*?)/g).filter(Boolean)
    const [mime, image] = imageInfo
    certifyForm.value.businessLicense.mime = mime
    certifyForm.value.businessLicense.image = image
    const data = await UploadImageRequest({ image: certifyForm.value.businessLicense })
    certifyForm.value.businessLicenseUrl = data.data.image_url
}

//获取上传图片名称
const logoImageUploaded = (val) => {
    certifyForm.value.logo.filename = val.name.split('.')[0]
    console.log(certifyForm.value);
}
const yingyeImageUploaded = (val) => {
    certifyForm.value.businessLicense.filename = val.name.split('.')[0]
    console.log(certifyForm.value);
}

//点击省重置市和县
const checksheng = () => {
    certifyForm.value.address.city = certifyForm.value.address.district = ''
}

//点击市重置县
const checkshi = () => {
    certifyForm.value.address.district = ''
}

//获取省级地区
const provinceOptions = () => {
    return getAreaOptions()
}

//获取市级地区
const cityOptions = () => {
    if (certifyForm.value.address.province) {
        return getAreaOptions(TextToCode[certifyForm.value.address.province]?.code)
    }
    return []

}

//获取县级地区
const regionOptions = () => {
    if (certifyForm.value.address.city) {
        return getAreaOptions(
            TextToCode[certifyForm.value.address.province][
                certifyForm.value.address.city
            ]?.code
        )
    }
    return []
}

//发送注册短信点击事件
const SendPhone = async (data, tem) => {
    duanxindj.value = true
    if (data.phone != '') {
        if (data.phone.length != 11) {
            elmessage('手机号格式错误')
            duanxindj.value = false
        } else {
            const rs = {

                phone: data.phone,
                templateAction: tem,
                language: 'LANGUAGE_SIMPLIFIED_CHINESE'
            }
            await SendPhoneVerificationCode(rs).then((res) => {
                console.log(res);
                if (res?.status === 200) {
                    userphone.value.txId = res.data.tx_id
                }
                else {
                    elmessage(res.data.detail)
                }
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
            })
        }
    } else {
        elmessage('请输入手机号')
        duanxindj.value = false
    }
}
//初始化进入页面
const conttext = ref('地址')
//是否创建完成
const found = ref(false)

//验证两次密码是否一致
let validatePass2 = (rule, value, callback) => {
    if (value === "") {
        callback(new Error("请再次输入密码"));
    } else if (value !== userInfo.value.newPlainPassword) {
        callback(new Error("两次输入密码不一致!"));
    } else {
        callback();
    }
};

//表单验证
const rules = reactive({
    newPlainPassword: [{ required: true, message: "请输入密码", trigger: 'blur' }],
    password: [{ required: true, trigger: 'blur', validator: validatePass2 }]
})

const Toggle = async (key) => {
    if (key === '地址') {
        if (certifyForm.value.logo.image == '' || certifyForm.value.name == '' || certifyForm.value.address.province == '' || certifyForm.value.address.city == '' || certifyForm.value.address.street == '') {
            elmessage('请填写完整的商户信息')
        } else {
            conttext.value = "营业执照"
        }
    }
    else if (key === '营业执照') {
        if (certifyForm.value.businessLicense.image == '' || certifyForm.value.socialCreditCode == '') {
            elmessage('请填写营业执照或社会信用代码')
        } else {
            conttext.value = "联系人"
        }
    }
    else if (key === '联系人') {
        if (certifyForm.value.contactName == '') {
            elmessage('联系人不能为空')
        } else
            if (userphone.value.phone === '' || userphone.value.smsCode === '') {
                elmessage('手机号或验证码为空')
            } else if (userphone.value.phone.length != 11) {
                elmessage('手机号格式错误')
            } else if (userphone.value.txId == '') {
                elmessage('验证码错误请重新获取')
            } else {
                await VerifyPhoneVerificationCodeRequest(userphone.value).then((res) => {
                    console.log(res);
                    if (res.status === 200) {
                        console.log('验证通过');
                        conttext.value = "商户"
                    }
                    else {
                        elmessage(res.data.detail)
                    }
                })
            }

    }
    else if (key === '商户') {
        if (userInfo.value.newPlainPassword == '') {
            elmessage('请输入商户密码')
        } else if (userInfo.value.newPlainPassword.length < 8 || userInfo.value.newPlainPassword.length > 16) {
            elmessage('密码格式错误请输入8到16位密码')
        } else if (userInfo.value.password == '') {
            elmessage('请再次输入商户密码')
        } else if (userInfo.value.newPlainPassword != userInfo.value.password) {
            return
        } else {
            certifyForm.value.contactPhone = certifyForm.value.safePhone = userphone.value.phone
            delete certifyForm.value.logo
            delete certifyForm.value.businessLicense
            const rs = {
                organizationId: localStorage.getItem('hmtenantId'),
                tenant: certifyForm.value,
                plainPassword: userInfo.value.newPlainPassword,
                smsCode: userphone.value.smsCode,
                txId: userphone.value.txId
            }
            console.log(rs);
            await CreateTenantRequest(rs).then((res) => {
                console.log(res);
                if (res.status === 200) {
                    console.log('验证通过');
                    conttext.value = '添加其他'
                }
                else {
                    elmessage(res.data.detail)
                    if (!certifyForm.value.logo) {
                        certifyForm.value.logo = {
                            mime: null,
                            image: null,
                            filename: null
                        }
                    }
                    if (!certifyForm.value.businessLicense) {
                        certifyForm.value.businessLicense = {
                            mime: null,
                            image: null,
                            filename: null
                        }
                    }
                }
            })
        }

    }
    if (key === '添加其他') {
        $router.push('/general')
    }
}
const copyText = async () => {
    let link = jmurl.value
    try {
        await toClipboard(link)
        ElMessage({
            message: '成功复制到粘贴板',
            type: 'success',
        })
    } catch {
        ElMessage({
            message: '复制失败',
            type: 'success',
        })
    }
}
const push = () => {
    $store.dispatch('getrouter')
    $router.push('/general')
}
</script>
<style lang="scss" scoped>
.found-general {
    background-color: rgba(246, 247, 248, 1);
    height: 90vh;

    .el-select {
        width: 189px;
    }

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
                .body-heard {
                    text-align: center;
                    font-size: 28px;
                    font-weight: 700;
                    padding-top: 70px;

                    .mess {
                        font-size: 18px;
                        font-weight: 400;
                    }

                    .messmin {
                        color: #C1C7D0;
                    }
                }
            }

            .body-from {
                width: 700px;
                margin: 50px auto;

                .body-table {
                    .table-img {
                        display: flex;
                        align-items: center;

                        .tebleleft {
                            text-align: right;
                        }

                        .tableright {
                            padding-left: 30px;
                            flex: 0.9;
                            display: flex;
                            align-items: center;

                            .imgtext {
                                margin-left: 20px;
                                font-weight: 400;
                                font-style: normal;
                                font-size: 14px;
                                color: rgb(193, 199, 208);
                                line-height: 28px;
                            }
                        }
                    }

                    .table-input {
                        margin-top: 20px;
                        display: flex;
                        align-items: center;

                        .tableright {
                            padding-left: 30px;
                            width: 100%;
                            flex: 0.9;
                        }
                    }

                    .body-link {
                        display: flex;

                        .linkinp {
                            flex: 1;
                            height: 40px;
                            margin-right: 20px;

                            .el-input__wrapper {
                                height: 40px;
                            }
                        }

                        .complete {
                            width: 126px;
                        }
                    }
                }
            }

            .body-bom {
                margin: 0 50px;
                width: 300px;
                margin: 50px auto;
                text-align: center;
                padding-bottom: 50px;
                display: flex;
                justify-content: space-around;
                align-items: center;
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
    width: 150px;
    height: 40px;
    background-color: rgba(21, 117, 238, 1);
}

.el-form-item {
    margin-bottom: 0px;
}
</style>