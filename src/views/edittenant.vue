<template>
    <div class="editt">
        <div class="lnvite">
            <div class="lnvitetop">
                <div class="toptitle">慧脉后台管理系统</div>
                <div class="toptxt" v-if="shangdata?.tenant">请填写你的商户信息</div>
            </div>
            <div class="lnvitedata" v-if="shangdata?.tenant">
                <el-form ref="ruleFormRef" :model="certifyForm" :rules="rules" :size="formSize" label-width="130px">
                    <el-form-item label="商户图标：" width=90 required="true" style="text-align:left">
                        <div class="tableright">
                            <imageUploader @imageBase64="logoImageBase64" @onImageUploaded="logoImageUploaded"
                                :imgurl="shangdata?.tenant?.logo_url">
                            </imageUploader>
                            <div class="imgtext">
                                <p>请上传商户门头照片/商户logo；</p>

                                <p>图片不允许涉及政治敏感与色情;</p>

                                <p>图片格式只支持：BMP、JPEG、JPG、GIF、PNG，大小不超过2M</p>
                            </div>
                        </div>
                    </el-form-item>
                    <el-form-item label="商户名称：" required="true" style="text-align:left">
                        <el-input placeholder="请输入商户名称" v-model="certifyForm.name" v-if="stecncildata" />
                        <span v-else>{{ certifyForm.name }}</span>
                    </el-form-item>
                    <el-form-item label="商户地址：" required="true" style="text-align:left" :error="addressmessage">
                        <div class="dizhi" v-if="stecncildata">
                            <div class="dizhishi">
                                <el-select v-model="certifyForm.address.province" placeholder="请选择" @change="checksheng"
                                    @blur="addressblur">
                                    <el-option v-for="item in provinceOptions()" :key="item.value" :label="item.label"
                                        :value="item.value">
                                    </el-option>
                                </el-select>
                                <el-select v-model="certifyForm.address.city" placeholder="请选择" @change="checkshi"
                                    @blur="addressblur">
                                    <el-option v-for="item in cityOptions()" :key="item.value" :label="item.label"
                                        :value="item.value">
                                    </el-option>
                                </el-select>
                                <el-select v-model="certifyForm.address.district" placeholder="请选择" @change="checkxian"
                                    @blur="addressblur">
                                    <el-option v-for="item in regionOptions()" :key="item.value" :label="item.label"
                                        :value="item.value">
                                    </el-option>
                                </el-select>
                            </div>
                            <el-form-item class="el-form-item " style="margin-top: 20px;">
                                <el-input placeholder="具体地址" v-model="certifyForm.address.street" @blur="addressblur" />
                            </el-form-item>
                        </div>
                        <span v-else>{{ certifyForm.address.province }}{{ certifyForm.address.city }}{{
                            certifyForm.address.district }}{{ certifyForm.address.street }}</span>
                    </el-form-item>
                    <el-form-item label="营业执照：" required="true" style="text-align:left">
                        <div class="tableright" v-if="stecncildata">
                            <imageUploader @imageBase64="yingyeImageBase64" @onImageUploaded="yingyeImageUploaded"
                                :imgurl="shangdata?.tenant?.business_license_url">
                            </imageUploader>
                            <div class="imgtext">
                                <p>请上传清晰的营业执照图片，保证以上商户信息与营业执照一致；</p>

                                <p>图片不允许涉及政治敏感与色情;</p>

                                <p>图片格式只支持：BMP、JPEG、JPG、GIF、PNG，大小不超过2M</p>
                            </div>
                        </div>
                        <div v-else class="tableright">
                            <div>
                                <img :src="certifyForm?.businessLicenseUrl" alt="" style="width: 148px;">
                            </div>
                            <div class="imgtext">
                                <p>请上传清晰的营业执照图片，保证以上商户信息与营业执照一致；</p>

                                <p>图片不允许涉及政治敏感与色情;</p>

                                <p>图片格式只支持：BMP、JPEG、JPG、GIF、PNG，大小不超过2M</p>
                            </div>
                        </div>
                    </el-form-item>
                    <el-form-item label="社会信用代码：" required="true" style="text-align:left">
                        <el-input placeholder="请输入社会信用代码" v-model="certifyForm.socialCreditCode" v-if="stecncildata" />
                        <span v-else>{{ certifyForm.socialCreditCode }}</span>
                    </el-form-item>
                    <el-form-item label="联系人：" required="true" style="text-align:left" prop="contactName">
                        <el-input placeholder="请输入商户联系人姓名" v-model="certifyForm.contactName" />
                    </el-form-item>
                    <el-form-item label="联系人手机号：" required="true" style="text-align:left" prop="contactPhone">
                        <el-input placeholder="请输入联系人手机号，用于商户登录" v-model="certifyForm.contactPhone" />
                    </el-form-item>
                    <el-form-item label="短信验证码" required="true" style="text-align:left" :error="codemessage">
                        <div class="yzm">
                            <div class="sr">
                                <el-input placeholder="请输入验证码" v-model="userphone.smsCode" maxlength="6"
                                    @blur="checkcodeDate" />
                            </div>
                            <span>
                                |
                            </span>
                            <div class="hq">
                                <el-link type="primary" :underline="false" class="huoqu"
                                    @click="SendPhone(certifyForm.contactPhone, 'TEMPLATE_ACTION_BIND_TENANT_PHONE')"
                                    :disabled="duanxindj" style="color: #1575EE;">{{ dxtxt }}</el-link>
                            </div>
                        </div>
                    </el-form-item>
                    <el-form-item label="商户登录密码" required="true" style="text-align:left" :error="birthDateMessage">
                        <el-input show-password placeholder="请设置商户登录密码" v-model="userInfo.newPlainPassword"
                            @blur="checkBirthDate" />
                    </el-form-item>
                    <el-form-item label="确认密码" required="true" style="text-align:left" :error="qrpassword">
                        <el-input show-password placeholder="请再次输入商户登录密码" v-model="userInfo.password" @blur="usepassword" />
                    </el-form-item>
                </el-form>
            </div>
            <div class="lnvitebom" v-if="shangdata?.tenant">
                <el-button type="primary" class="lnvitebum" @click="CreateTenant(ruleFormRef)">提交审核</el-button>

                <div class="banquan">
                    版权所有：常州金姆健康科技有限公司
                </div>
            </div>
            <div v-if="!shangdata?.has_unauth_revision" class="meiyou">
                当前商户已提审成功或无需修改
            </div>
        </div>
    </div>
</template>
<script setup>
import imageUploader from '@/components/common/imageUploader.vue';
import { getAreaOptions, TextToCode } from '@/utils/areadDataOptions'
import { GetUnAuthTenantRevisionRequest, ReCreateTenantRequest, SendPhoneVerificationCode, GetTenantStencilRequest, UploadImageRequest } from '@/api/api';
import { ref, reactive } from 'vue'
import { elmessage } from '../utils/popup'
import { useRouter } from 'vue-router';
import { rulesz } from '../utils/rules'

const $router = useRouter();

const shangdata = ref(null)
const qiqu = ref(null)
const tenantid = ref($router.currentRoute.value.query.t)
console.log(tenantid.value);
const certifyForm = ref({
    organizationId: null,
    tenantId: null,
    logo: {
        mime: null,
        image: null,
        filename: null
    },
    name: null,
    contactName: null,
    address: {
        province: '',
        city: '',
        district: '',
        street: ''
    },
    
    contactPhone: null,
    businessLicense: {
        mime: null,
        image: null,
        filename: null
    },
    socialCreditCode: null,
    logoUrl: null,
    businessLicenseUrl: null,
    safePhone: null
})
const stecncildata = ref(true)
//短信60秒间隔
const time = ref(60)
const dxtxt = ref('获取验证码')
const duanxindj = ref(false)
const timesphone = ref(null)
const stencidata = ref(null)
//表单验证
const formSize = ref('default')
const ruleFormRef = ref()
//密码校验
const birthDateMessage = ref('')
//短信验证
const codemessage = ref('')

//确认密码验证
const qrpassword = ref('')
const rules = reactive(rulesz)
const phoneMsg = ref('')

const addressmessage = ref('')
const userInfo = ref({
    newPlainPassword: '',
    password: ''
})
const userphone = ref({
    smsCode: '',
    txId: '',
    templateAction: 'TEMPLATE_ACTION_BIND_TENANT_PHONE'
})
const checkBirthDate = () => {
    if (userInfo.value.newPlainPassword.length < 6 || userInfo.value.newPlainPassword.length > 18) {
        birthDateMessage.value = '请输入6到18位密码'
    } else if (userInfo.value.newPlainPassword != '') {
        birthDateMessage.value = ''
    } else {
        birthDateMessage.value = '请输入登录密码'

    }
}
const checkcodeDate = () => {
    if (userphone.value.smsCode != '') {
        codemessage.value = ''
    } else {
        codemessage.value = '请输入短信验证码'
    }
}
const usepassword = () => {
    if (userInfo.value.password == "") {
        qrpassword.value = "请再次输入密码";
    } else if (userInfo.value.password !== userInfo.value.newPlainPassword) {
        qrpassword.value = "两次输入密码不一致!";
    } else {
        qrpassword.value = '';
    }
};
function getUrlParams(key) {
    let reg = new RegExp("(^|&)" + key + "=([^&]*)(&|$)");
    let r = window.location.search.substr(1)
        .match(reg);
    if (r != null)
        return unescape(r[2]);
    return null;
}
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
const yingyeImageBase64 = (val) => {
    console.log(val);
    const imageInfo = val.split(/data:(.*?);base64,(.*?)/g).filter(Boolean)
    const [mime, image] = imageInfo
    certifyForm.value.businessLicense.mime = mime
    certifyForm.value.businessLicense.image = image
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

const GetTenantEntity = () => {
    GetUnAuthTenantRevisionRequest({ tenantId: tenantid.value }).then(res => {
        console.log(res);
        if (res.status == 200) {
            if (res.data.has_unauth_revision) {
                shangdata.value = res.data
                console.log(shangdata);
                certifyForm.value = {
                    organizationId: shangdata.value?.tenant.organization_id,
                    tenantId: shangdata.value?.tenant.tenant_id,
                    logo: {
                        mime: null,
                        image: null,
                        filename: null
                    },
                    
                    name: shangdata.value?.tenant.name,
                    contactName: shangdata.value?.tenant.contact_name,
                    address: {
                        province: shangdata.value?.tenant.address.province,
                        city: shangdata.value?.tenant.address.city,
                        district: shangdata.value?.tenant.address.district == undefined ? '' : shangdata.value?.tenant.address.district,
                        street: shangdata.value?.tenant.address.street
                    },
                    contactPhone: shangdata.value?.tenant.contact_phone,
                    businessLicense: {
                        mime: null,
                        image: null,
                        filename: null
                    },
                    logoUrl: shangdata.value?.tenant.logo_url,
                    businessLicenseUrl: shangdata.value?.tenant.business_license_url,
                    socialCreditCode: shangdata.value?.tenant.social_credit_code,
                }
                if (res.data.tenant.use_stencil && res.data.tenant.is_stencil_modifiable) {
                    stecncildata.value = true
                }
            } else {
                shangdata.value = false
            }
        }
        console.log(certifyForm.value);
    })
}
//发送创建短信点击事件
const SendPhone = async (data, tem) => {
    duanxindj.value = true
    if (data != '') {
        const rs = {
            
            phone: data,
            templateAction: tem,
            language: 'LANGUAGE_SIMPLIFIED_CHINESE'
        }
        await SendPhoneVerificationCode(rs).then((res) => {
            console.log(res);
            if (res?.status === 200) {
                userphone.value.txId = res.data.tx_id
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
    } else {
        elmessage('请输入手机号')
        duanxindj.value = false
    }
}
const CreateTenant = async (formEl) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
        if (valid) {
            if (certifyForm.value.address.province != null && certifyForm.value.address.city != null && certifyForm.value.address.street != null) {
                addressmessage.value = ''
            }
            checkBirthDate()
            checkcodeDate()
            usepassword()
            if (certifyForm.value.logo.image == null && certifyForm.value.logoUrl == null) {
                elmessage('请上传商户图标')
            } else if (certifyForm.value.businessLicense.image == null && certifyForm.value.businessLicenseUrl == null) {
                elmessage('请上传营业执照')
            } else if (certifyForm.value.address.province == '' || certifyForm.value.address.city == '' || certifyForm.value.address.street == '') {
                elmessage('请输入完整地址')
                addressmessage.value = '请输入完整地址'
            } else if (userInfo.value.newPlainPassword == '') {
                return
            } else if (userInfo.value.newPlainPassword.length < 6 || userInfo.value.newPlainPassword.length > 18) {
                return
            } else if (userInfo.value.password !== userInfo.value.newPlainPassword) {
                return
            } else if (userphone.value.txId == '') {
                elmessage('验证码错误请重新获取')
            }
            else {
                certifyForm.value.safePhone = certifyForm.value.contactPhone
                delete certifyForm.value.logo

                if (certifyForm.value.businessLicense.mime == null) {
                    delete certifyForm.value.businessLicense
                }
                const rs = {
                    tenant: certifyForm.value,
                    plainPassword: userInfo.value.newPlainPassword,
                    smsCode: userphone.value.smsCode,
                    txId: userphone.value.txId
                }
                console.log(rs);
                await ReCreateTenantRequest(rs).then((res) => {
                    console.log(res);
                    if (res.status === 200) {
                        console.log('提审成功');
                        GetTenantEntity()
                        location.reload()
                    }
                    else {
                        elmessage(res.data.detail)
                        if (!certifyForm.value.businessLicense) {
                            certifyForm.value.businessLicense = {
                                mime: null,
                                image: null,
                                filename: null
                            }
                        }
                        if (!certifyForm.value.logo) {
                            certifyForm.value.logo = {
                                mime: null,
                                image: null,
                                filename: null
                            }
                        }
                    }
                })
            }
        } else {
            checkBirthDate()
            checkcodeDate()
            usepassword()
            if (certifyForm.value.address.province == '' || certifyForm.value.address.city == '' || certifyForm.value.address.street == '') {
                addressmessage.value = '请输入完整地址'
                console.log(certifyForm.value);
            } else {
                addressmessage.value = ''

            }

        }
    })
}
GetTenantEntity()
</script>
<style lang="scss">
.lnvite {
    padding: 20px;
    max-width: 800px;
    margin: 0 auto;

    .lnvitetop {
        text-align: center;

        .toptitle {
            font-size: 20px;
            margin-top: 40px;
        }

        .toptxt {
            font-size: 15px;
            margin-top: 10px;
            margin-bottom: 30px;
        }
    }
}

.lnvitedata {
    .el-form-item {
        display: block;
    }

    .el-form-item__label {
        display: block;
    }
}

.editt .tableright {
    display: flex;
    align-items: center;

    .imgtext {
        margin-left: 10px;
        font-size: 12px;
        color: #C1C7D0;
        line-height: 22px;
    }
}

.rightdizhi {
    flex-direction: column;
}

.lnvitebom {

    .lnvitebum {
        width: 100%;
        height: 40px;
        margin-bottom: 40px;
    }

    .banquan {
        margin: 20px;
        text-align: center;
        font-size: 12px;
        color: #C1C7D0
    }
}

.yzm {
    display: flex;
    flex: 1;
    box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color)) inset;

    .sr {
        flex: 1;

        .el-input__wrapper {
            box-shadow: none;
        }
    }

    .hq {
        color: #1575EE;
        font-size: 14px;
        width: 100px;
        text-align: center;
    }

    span {
        font-size: 14px;
        color: #1575EE;
    }
}

.el-input__wrapper {
    background-color: transparent;
}

.meiyou {
    height: 60vh;
    font-size: 18px;
    color: rgb(193, 199, 208);
    display: flex;

    align-items: center;
    justify-content: center;
}
</style>