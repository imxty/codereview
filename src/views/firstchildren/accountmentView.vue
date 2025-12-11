<template>
    <div class="account">
        <div class="found-general">
            <div class="found-cont steat">
                <div class="cont-text">
                    <div class="cont-flex">
                        <div class="shubi">
                            <span @click="status = '账号信息'" class="shubi">账号管理</span>
                            <span v-if="status == '验证原手机号' || status == '验证新手机号'">/修改手机号</span>
                            <span v-if="status == '修改密码'">/修改密码</span>
                            <span v-if="status == '组织手机号'">/修改组织客服手机号</span>
                        </div>
                        <div class="cont-fh">
                            <div @click="tohome('/general')" class="shubi" style="margin-right: 20px;">
                                返回后台
                            </div>
                            <div class="shubi" @click="tohome('login')">
                                退出
                            </div>
                        </div>
                    </div>
                </div>
                <div>
                    <!-- v-if="status == '抬头管理'" -->
                    <el-button type="primary elbutton newtt" @click="createdling"
                        v-if="status == '抬头管理'">新建抬头</el-button>
                </div>
                <div class="mainbody">
                    <div class="body-cont">
                        <div v-if="status == '账号信息'">
                            <div class="detacnt">
                                <div class="detacnttop">
                                    基本信息
                                </div>
                                <div class="cntdata">
                                    <div class="username">
                                        <div class="cntadataleft">
                                            账户名：
                                        </div>
                                        <div class="cntadataright">
                                            <div>
                                                {{ reviewtype?.name }}
                                            </div>
                                            <div>

                                            </div>
                                        </div>
                                    </div>
                                    <div class="username">
                                        <div class="cntadataleft">
                                            手机号：
                                        </div>
                                        <div class="cntadataright">
                                            <div>
                                                {{ geTel(reviewtype?.phone) }}
                                            </div>
                                            <div style="color:#1575EE;cursor: pointer;" @click="newphone">
                                                绑定新手机
                                            </div>
                                        </div>
                                    </div>
                                    <div class="username">
                                        <div class="cntadataleft">
                                            登录密码：
                                        </div>
                                        <div class="cntadataright">
                                            <div>
                                                ******
                                            </div>
                                            <div style="color:#1575EE;cursor: pointer;" @click="newyibu">
                                                修改密码
                                            </div>
                                        </div>
                                    </div>
                                    <div class="username">
                                        <div class="cntadataleft">
                                            组织客服手机号：
                                        </div>
                                        <div class="cntadataright">
                                            <div>
                                                {{ geTel(reviewtype?.organization_contact_phone) }}
                                            </div>
                                            <div style="color:#1575EE;cursor: pointer;" @click="neworgan">
                                                绑定新手机
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div v-if="status == '验证原手机号'" class="cont-data">
                            <div style="text-align: center;" class="toptext">
                                <p>为保障您的账号安全，更换绑定需验证原手机号</p>请输入已绑定手机号：<span style="color: #1575EE;">+86
                                    {{ geTel(reviewtype?.phone) }}</span>
                                收到的短信验证码
                            </div>
                            <el-form-item prop="account" class="el-form-items ">
                                <el-input v-model="phonedata.oldPhoneSmsCode" placeholder="请输入验证码" />
                                <el-button type="primary dingwei elbutton"
                                    @click="SendPhone('TEMPLATE_ACTION_MODIFY_ORGANIZATION_PHONE', reviewtype?.phone)"
                                    :disabled="duanxindj">{{ dxtxt }}</el-button>
                            </el-form-item>
                            <el-button type="primary elbutton elxyb"
                                @click="VerifyPhoneVerificationCode">下一步</el-button>
                        </div>
                        <div v-if="status == '验证新手机号'" class="cont-data">
                            <div>请验证新手机号</div>
                            <el-form-item prop="account" class="el-form-items ">
                                <el-input v-model="phonedata.newSafePhone" placeholder="请输入手机号" />
                            </el-form-item>
                            <el-form-item prop="account" class="el-form-items ">
                                <el-input v-model="phonedata.newPhoneSmsCode" placeholder="请输入验证码" />
                                <el-button type="primary elbutton dingwei"
                                    @click="newSendPhone('TEMPLATE_ACTION_MODIFY_ORGANIZATION_PHONE', phonedata.newSafePhone)"
                                    :disabled="duanxindj">{{ dxtxt }}</el-button>
                            </el-form-item>
                            <el-button type="primary elbutton elxyb" @click="UpdateOrganizationPhone">完成</el-button>
                        </div>
                        <div v-if="status == '修改密码'" class="cont-data">
                            <el-form-item prop="account" class="el-form-items ">
                                <el-input v-model="userInfo.onloadpassword" placeholder="请输入当前密码" />
                            </el-form-item>
                            <el-form-item prop="account" class="el-form-items">
                                <el-input v-model="userInfo.newpassword" placeholder="请输入新密码" />
                            </el-form-item>
                            <el-form-item prop="account" class="el-form-items ">
                                <el-input v-model="userInfo.password" placeholder="请再次输入新密码" />
                            </el-form-item>
                            <el-button type="primary elbutton elxyb" @click="UpdateOrganizationPassword">完成</el-button>
                        </div>
                        <div v-if="status == '抬头管理'" class="cont-data" style="padding-top: 10px;">
                            <el-table :data="tableData" style="width: 100%">
                                <el-table-column prop="title_name" label="抬头名称" width="240" />
                                <el-table-column prop="tax_id" label="税号" width="240" />
                                <el-table-column prop="title_type" label="类型">
                                    <template #default="scope">
                                        <span v-if="scope.row.title_type == 'TITLE_TYPE_SPECIAL'">
                                            专票
                                        </span>
                                        <span v-else>
                                            普票
                                        </span>
                                    </template>
                                </el-table-column>
                                <el-table-column label="操作">
                                    <template #default="scope">
                                        <el-button type="primary" plain @click="editdling(scope.row)">
                                            修改</el-button>
                                        <el-button type="danger" plain @click="deletefunction(scope.row)">删除</el-button>
                                    </template>
                                </el-table-column>
                            </el-table>
                        </div>
                        <div v-if="status == '组织手机号'" class="cont-data">
                            <el-form-item prop="account" class="el-form-items ">
                                <el-input v-model="organizationContactPhone" placeholder="请输入新手机号" />
                            </el-form-item>
                            <el-button type="primary elbutton elxyb" @click="UpdateOrganizationContactPhone">完成</el-button>
                        </div>
                    </div>
                    <div class="bottom-zhu">
                        如您需要续期或在使用中有其他问题，请联系大拇哥客服，电话：17807968188
                    </div>
                </div>
                <div class="fenye" v-if="status == '抬头管理'">
                    <el-pagination v-model:current-page="offset" v-model:page-size="pageSize4"
                        :page-sizes="[10, 30, 50, 100]" :small="true" :background="background"
                        layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange"
                        @current-change="handleCurrentChange" />
                </div>
            </div>
        </div>
        <el-dialog v-model="dialogVisible" width="50%" :destroy-on-close=true :show-close=false class="changjiantan"
            top="30px">
            <div class="newtop">
                <div class="toptxt">
                    请填写发票信息
                </div>
                <div class="newcnt">
                    <el-form ref="ruleFormRef" style="width: 90%" :model="ruleForm" :rules="rules" label-width="auto"
                        class="demo-ruleForm" :size="formSize" status-icon>
                        <el-form-item label="票据类型：" prop="titleType">
                            <el-radio-group v-model="ruleForm.titleType">
                                <el-radio label="TITLE_TYPE_SPECIAL">专票</el-radio>
                                <el-radio label="TITLE_TYPE_COMMON">普票</el-radio>
                            </el-radio-group>
                        </el-form-item>
                        <el-form-item label="名称：" prop="titleName">
                            <el-input v-model="ruleForm.titleName" />
                        </el-form-item>
                        <el-form-item label="税号：" prop="taxId">
                            <el-input v-model="ruleForm.taxId"
                                @input="(v) => (ruleForm.taxId = v.replace(/[^\w_]/g, ''))" />
                        </el-form-item>
                        <el-form-item label="单位地址：" prop="companyAddress"
                            v-if="ruleForm.titleType == 'TITLE_TYPE_SPECIAL'">
                            <el-input v-model="ruleForm.companyAddress" />
                        </el-form-item>
                        <el-form-item label="单位地址：" v-else>
                            <el-input v-model="ruleForm.companyAddress" />
                        </el-form-item>

                        <el-form-item label="电话：" prop="companyPhone" v-if="ruleForm.titleType == 'TITLE_TYPE_SPECIAL'">
                            <el-input v-model="ruleForm.companyPhone"
                                @input="(v) => (ruleForm.companyPhone = v.replace(/[^\d]/g, ''))" />
                        </el-form-item>
                        <el-form-item label="电话：" v-else>
                            <el-input v-model="ruleForm.companyPhone"
                                @input="(v) => (ruleForm.companyPhone = v.replace(/[^\d]/g, ''))" />
                        </el-form-item>
                        <el-form-item label="开户银行：" prop="bank" v-if="ruleForm.titleType == 'TITLE_TYPE_SPECIAL'">
                            <el-input v-model="ruleForm.bank" />
                        </el-form-item>
                        <el-form-item label="开户银行：" v-else>
                            <el-input v-model="ruleForm.bank" />
                        </el-form-item>
                        <el-form-item label="银行账户：" prop="bankAccount"
                            v-if="ruleForm.titleType == 'TITLE_TYPE_SPECIAL'">
                            <el-input v-model="ruleForm.bankAccount"
                                @input="(v) => (ruleForm.bankAccount = v.replace(/[^\d]/g, ''))" />
                        </el-form-item>
                        <el-form-item label="银行账户：" v-else>
                            <el-input v-model="ruleForm.bankAccount"
                                @input="(v) => (ruleForm.bankAccount = v.replace(/[^\d]/g, ''))" />
                        </el-form-item>
                        <div class="jsxx">
                            寄送信息
                        </div>
                        <el-form-item label="组织名称：" prop="organizationName">
                            <el-input v-model="ruleForm.organizationName" disabled />
                        </el-form-item>
                        <el-form-item label="联系电话：" prop="contactPhone">
                            <el-input v-model="ruleForm.contactPhone"
                                @input="(v) => (ruleForm.contactPhone = v.replace(/[^\d]/g, ''))" />
                        </el-form-item>
                        <el-form-item label="收票邮箱：" prop="contactEmail">
                            <el-input v-model="ruleForm.contactEmail" />
                        </el-form-item>
                        <el-form-item label="收票地址：" prop="acceptAddress">
                            <el-input v-model="ruleForm.acceptAddress" />
                        </el-form-item>
                    </el-form>
                </div>
            </div>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click.stop="dialogVisible = false" class="quxiaobut">取消</el-button>
                    <el-button type="primary" @click.stop="CreateInvoiceTitle(ruleFormRef)" class="tijiaobut"
                        v-if="xiugaixj == 'new'">
                        确认
                    </el-button>
                    <el-button type="primary" @click.stop="ModifyInvoiceTitle(ruleFormRef)" class="tijiaobut" v-else>
                        确认
                    </el-button>
                </span>
            </template>
        </el-dialog>
        <el-dialog v-model="daletsp" title="删除抬头" width="30%" :before-close="handleClose">
            <span>是否确定删除"{{ deleteinvoice.title_name }}"发票抬头?</span>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="daletsp = false" class="elbut">取消</el-button>
                    <el-button type="primary" @click="DeleteInvoiceTitle" class="elbut quern">
                        确认删除
                    </el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<script setup>
import { ref } from "vue";
import { useStore } from 'vuex'
import {
    UpdateOrganizationPasswordRequest,
    UpdateOrganizationPhoneRequest,
    CheckOrganizationPhoneExistRequest,
    VerifyPhoneVerificationCodeRequest,
    SendPhoneVerificationCode,
    GetOrganizationDetailRequest,
    ListOrganizationInvoiceTitlesRequest,
    CreateInvoiceTitleRequest,
    ModifyInvoiceTitleRequest,
    DeleteInvoiceTitleRequest,
    UpdateOrganizationContactPhoneRequest
} from '@/api/api'
import { geTel } from "../../utils/phontel";
import { getTenantId } from '@/utils/storage'
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { elmessage } from '@/utils/popup'

const $store = useStore();
const $router = useRouter();
const status = ref('账号信息')
const time = ref(60)
const dxtxt = ref('获取验证码')
const duanxindj = ref(false)
const reviewtypes = ref(JSON.parse(history.state.keyword))
const timesphone = ref(null)

const reviewtype = ref()

const dialogVisible = ref(false)

const daletsp = ref(false)

const ruleFormRef = ref()

const xiugaixj = ref('new')

const organizationContactPhone = ref()
//默认每页几条数据
const pageSize4 = ref(10)
//默认总数据
const total = ref(0)
//分页偏移量
const offset = ref(1)

const tableData = ref([])

const deleteinvoice = ref()

const ruleForm = ref({
    organizationId: getTenantId(),
    titleType: 'TITLE_TYPE_SPECIAL',
    titleName: '',
    taxId: '',
    companyAddress: '',
    companyPhone: '',
    bank: '',
    bankAccount: '',
    organizationName: reviewtype.value?.name,
    contactPhone: '',
    contactEmail: '',
    acceptAddress: ''
})

const rules = ref({
    titleName: [
        { required: true, message: '名称不能为空', trigger: 'blur' },
        { max: 20, message: '字数超过上限', trigger: 'blur' },
    ],
    taxId: [
        {
            required: true,
            message: '税号不能为空',
            trigger: 'blur',
        },
    ],
    companyAddress: [
        {
            required: true,
            message: '单位地址不能为空',
            trigger: 'blur',
        },
        {
            required: true,
            max: 50,
            message: '字数不能超过50位',
            trigger: 'blur',
        },
    ],
    companyPhone: [
        {
            required: true,
            message: '电话不能为空',
            trigger: 'blur',
        },
        {
            required: true,
            max: 20,
            message: '请填写正确的电话号码',
            trigger: 'blur',
        },
    ],
    bank: [
        {
            required: true,
            message: '开户银行不能为空',
            trigger: 'blur',
        },
        {
            required: true,
            max: 50,
            message: '字数不能超过50位',
            trigger: 'blur',
        },
    ],
    titleType: [
        {
            required: true,
            message: '请选择一个类型',
            trigger: 'blur',
        },
    ],
    bankAccount: [
        {
            required: true,
            message: '银行账户不能为空',
            trigger: 'blur',
        },
        {
            required: true,
            max: 50,
            message: '字数不能超过50位',
            trigger: 'blur',
        },
    ],
    organizationName: [
        {
            required: true,
            message: '组织名称不能为空',
            trigger: 'blur',
        },
        {
            required: true,
            max: 50,
            message: '字数不能超过50位',
            trigger: 'blur',
        },
    ],
    contactPhone: [
        {
            required: true,
            message: '联系电话不能为空',
            trigger: 'blur',
        }, {
            required: true,
            min: 11,
            max: 11,
            message: '请填写正确的电话号码',
            trigger: 'blur',
        },
    ],
    contactEmail: [
        {
            type: 'email',
            required: true,
            message: '邮箱格式错误',
            trigger: 'blur',
        },
        {
            required: true,
            max: 50,
            message: '字数不能超过50位',
            trigger: 'blur',
        },
    ],
})

//接收密码
const userInfo = ref({
    onloadpassword: '',
    newpassword: '',
    password: ''
})

const phonedata = ref({
    organizationId: getTenantId(),
    oldPhoneSmsCode: null,
    oldPhoneSmsTxId: null,
    newSafePhone: null,
    newPhoneSmsCode: null,
    newPhoneSmsTxId: null
}
)

//点击下一步进度条变换
const newphone = () => {
    status.value = '验证原手机号'
    phonedata.value = {
        organizationId: getTenantId(),
        oldPhoneSmsCode: null,
        oldPhoneSmsTxId: null,
        newSafePhone: null,
        newPhoneSmsCode: null,
        newPhoneSmsTxId: null
    }
}

const newyibu = () => {
    status.value = '修改密码'
    userInfo.value = {
        onloadpassword: '',
        newpassword: '',
        password: ''
    }
}

const neworgan = () => {
    status.value = '组织手机号'
}
const deletefunction = (data) => {
    deleteinvoice.value = data
    daletsp.value = true
}

const handleSizeChange = (val) => {
    ListOrganizationInvoiceTitles()
}
const handleCurrentChange = (val) => {
    ListOrganizationInvoiceTitles()
}
const newfapiao = () => {
    status.value = '抬头管理'
    ListOrganizationInvoiceTitles()
}

const ListOrganizationInvoiceTitles = async () => {
    const rs = {
        organizationId: getTenantId(),
        pagination: {
            offset: (offset.value - 1) * pageSize4.value,
            size: pageSize4.value
        },
    }
    const data = await ListOrganizationInvoiceTitlesRequest(rs)
    tableData.value = data.data.titles
    total.value = data.data.total_count == undefined ? '0' : data.data.total_count
}


const DeleteInvoiceTitle = async () => {
    console.log(deleteinvoice);
    const rs = {
        organizationId: getTenantId(),
        invoiceTitleId: deleteinvoice.value.title_id,
    }
    const data = await DeleteInvoiceTitleRequest(rs)
    if (data.status == 200) {
        ElMessage({
            message: '删除成功',
            type: 'success',
        })
        ListOrganizationInvoiceTitles()
        daletsp.value = false
    } else {
        elmessage(res.data.detail)
    }
}

const CreateInvoiceTitle = async (formEl) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
        if (valid) {
            const rs = {
                invoiceTitle: ruleForm.value,
                organizationId: getTenantId(),
            }
            const data = await CreateInvoiceTitleRequest(rs)
            if (data.status == 200) {
                ElMessage({
                    message: '新建成功',
                    type: 'success',
                })
                ListOrganizationInvoiceTitles()
                dialogVisible.value = false
            } else {
                elmessage(data.data.detail)
            }
            console.log(data);
        } else {
            console.log('error submit!', fields)
        }
    })
}

const ModifyInvoiceTitle = async (formEl) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
        if (valid) {
            const rs = {
                invoiceTitle: ruleForm.value,
                organizationId: getTenantId(),
            }
            const data = await ModifyInvoiceTitleRequest(rs)
            if (data.status == 200) {
                ElMessage({
                    message: '修改成功',
                    type: 'success',
                })
                ListOrganizationInvoiceTitles()
                dialogVisible.value = false
            } else {
                elmessage(data.data.detail)
            }
            console.log(data);
        } else {
            console.log('error submit!', fields)
        }
    })
}

const createdling = () => {
    xiugaixj.value = 'new'
    ruleForm.value = {
        organizationId: getTenantId(),
        titleType: 'TITLE_TYPE_SPECIAL',
        titleName: '',
        taxId: '',
        companyAddress: '',
        companyPhone: '',
        bank: '',
        bankAccount: '',
        organizationName: reviewtype.value?.name,
        contactPhone: '',
        contactEmail: '',
        acceptAddress: ''
    }
    dialogVisible.value = true
}

const editdling = (data) => {
    xiugaixj.value = 'edit'
    console.log(data);
    ruleForm.value = {
        titleId: data.title_id,
        organizationId: data.organization_id,
        titleType: data.title_type,
        titleName: data.title_name,
        taxId: data.tax_id,
        companyAddress: data.company_address,
        companyPhone: data.company_phone,
        bank: data.bank,
        bankAccount: data.bank_account,
        organizationName: data.organization_name,
        contactPhone: data.contact_phone,
        contactEmail: data.contact_email,
        acceptAddress: data.accept_address === undefined ? '' : data.accept_address
    }
    dialogVisible.value = true
}
const UpdateOrganizationPassword = async () => {
    const rs = {
        organizationId: getTenantId(),
        oldPlainPassword: userInfo.value.onloadpassword,
        newPlainPassword: userInfo.value.newpassword
    }
    const data = await UpdateOrganizationPasswordRequest(rs)
    if (data.status === 200) {
        ElMessage({
            message: '修改密码成功',
            type: 'success',
        })
        $router.push({ name: 'login' })
    } else {
        ElMessage({
            showClose: true,
            message: data.data.detail,
            type: 'error',
        })
    }
}
//发送注册短信点击事件
const SendPhone = async (tem, phone) => {
    duanxindj.value = true
    if (!/^1[3456789]\d{9}$/.test(phone)) {
        elmessage('请输入正确的手机号')
        return;
    }
    if (phone != '') {
        const rs = {

            phone: phone,
            templateAction: tem,
            language: 'LANGUAGE_SIMPLIFIED_CHINESE'
        }
        console.log(rs);
        await SendPhoneVerificationCode(rs).then((res) => {
            console.log(res);
            if (res.status === 200) {
                phonedata.value.oldPhoneSmsTxId = res.data.tx_id
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
    } else {
        elmessage('请输入手机号')
    }
}
const newSendPhone = async (tem, phone) => {
    if (!/^1[3456789]\d{9}$/.test(phone)) {
        elmessage('请输入正确的手机号')
        return;
    }
    if (phone != '') {
        const data = await CheckOrganizationPhoneExistRequest({ phone: phone })
        if (data?.data.exist) {
            elmessage('手机号已存在')
        } else {
            console.log(data);
            const rs = {

                phone: phone,
                templateAction: tem,
                language: 'LANGUAGE_SIMPLIFIED_CHINESE'
            }
            console.log(rs);
            await SendPhoneVerificationCode(rs).then((res) => {
                console.log(res);
                if (res.status === 200) {
                    phonedata.value.newPhoneSmsTxId = res.data.tx_id
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
                }
            })
        }
    } else {
        elmessage('请输入手机号')
    }
}
const submitForm = async (formEl) => {
    if (!formEl) return
    await formEl.validate((valid, fields) => {
        if (valid) {
            console.log('submit!')
        } else {
            console.log('error submit!', fields)
        }
    })
}

const resetForm = (formEl) => {
    if (!formEl) return
    formEl.resetFields()
}
const UpdateOrganizationPhone = async () => {
    if (phonedata.value.newSafePhone.length != 11) {
        elmessage('手机号码格式错误')
    } else if (!phonedata.value.newPhoneSmsTxId) {
        elmessage('验证码错误请重新获取')
    } else if (phonedata.value.newPhoneSmsCode === '') {
        elmessage('验证码为空')
    } else {
        const data = await UpdateOrganizationPhoneRequest(phonedata.value)
        if (data.status === 200) {
            ElMessage({
                message: '修改手机号成功',
                type: 'success',
            })
            GetOrganizationDetail()
            status.value = '账号信息'
        } else {
            ElMessage({
                showClose: true,
                message: data.data.detail,
                type: 'error',
            })
        }
    }

}

//验证手机号
const VerifyPhoneVerificationCode = async () => {
    if (!phonedata.value.oldPhoneSmsTxId) {
        elmessage('验证码错误请重新获取')
    }
    else if (phonedata.value.oldPhoneSmsCode === '') {
        elmessage('验证码为空')
    } else {
        const rs = {
            areaCode: phonedata.value.newAreaCode,
            phone: reviewtype.value.phone,
            smsCode: phonedata.value.oldPhoneSmsCode,
            txId: phonedata.value.oldPhoneSmsTxId,
            templateAction: 'TEMPLATE_ACTION_MODIFY_ORGANIZATION_PHONE'
        }
        await VerifyPhoneVerificationCodeRequest(rs).then((res) => {
            console.log(res);
            if (res.status === 200) {
                console.log('验证通过')
                time.value = 60;//当减到0时赋值为60
                dxtxt.value = '获取验证码';
                clearInterval(timesphone.value);//清除定时器
                duanxindj.value = false
                status.value = '验证新手机号'
            }
            else {
                elmessage(res.data.detail)
            }
        })
    }
}
//更改组织联系人手机号
const UpdateOrganizationContactPhone = async () => {
    if (organizationContactPhone.value.length != 11) {
        elmessage('手机号码格式错误')
    } else {
        const res = {
            organizationId: getTenantId(),
            organizationContactPhone: organizationContactPhone.value
        }
        const data = await UpdateOrganizationContactPhoneRequest(res)
        if (data.status === 200) {
            ElMessage({
                message: '修改手机号成功',
                type: 'success',
            })
            GetOrganizationDetail()
            status.value = '账号信息'
        } else {
            ElMessage({
                showClose: true,
                message: data.data.detail,
                type: 'error',
            })
        }
    }
}

const GetOrganizationDetail = async () => {
    const data = await GetOrganizationDetailRequest({ organizationId: getTenantId() })
    console.log(data);
    reviewtype.value = data.data.organization
}
const tohome = (data) => {
    $router.push(data)
    console.log($router);

}
GetOrganizationDetail()
if (reviewtypes.value == 'invoice') {
    ListOrganizationInvoiceTitles()
    status.value = '抬头管理'
}
</script>

<style lang="scss">
// 更改element的默认样式
.account {
    .elbut {
        width: 103px;
        height: 40px;
    }

    .quern {
        background-color: #1575ee;
    }

    .newcnt {
        margin-top: 30px;
        width: 110%;
        display: flex;

        .el-input__wrapper {
            height: 40px;
        }
    }

    .found-general {
        background-color: rgba(246, 247, 248, 1);
        min-height: 89vh;
        padding-bottom: 20px;

        .found-cont {
            position: relative;

            .el-input {
                --el-input-focus-border-color: none;
                height: 40px;
            }

            .el-input__wrapper {
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

            .el-input__wrapper {
                box-shadow: none !important;

                &:hover {
                    box-shadow: none;
                }
            }

            .fenye {
                margin-top: 30px;
                position: absolute;
                bottom: -55px;
                right: 0px;

                .el-input__wrapper {
                    height: 30px;
                    box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color)) inset !important;
                    border-radius: 2px;
                }

                .el-input__suffix {
                    margin-right: 10px;
                }
            }

            .cont-text {
                width: 100%;
                height: 8vh;
                line-height: 8vh;
                font-size: 24px;
                font-weight: 400;
                font-style: normal;
                text-align: left;
                border-bottom: 1px solid #333333;



                .cont-flex {
                    display: flex;
                    justify-content: space-between;

                    .cont-fh {
                        display: flex;
                        justify-content: space-between;
                        color: #1575EE;
                        font-size: 13px;
                    }
                }
            }

            .mainbody {
                background-color: #fff;
                margin-top: 30px;
                position: relative;
                border-radius: 20px;

                .body-cont {
                    box-sizing: content-box;
                    min-height: 550px;

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

                .body-from {
                    width: 380px;
                    margin: 50px auto;

                    .el-form-item {
                        margin-top: 40px;
                        min-height: 40px;

                        .el-input__inner {
                            font-size: 16px;
                            height: 40px;
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

    .el-form-items {
        margin-bottom: 0px;
        position: relative;
        width: 380px;
        margin-top: 20px;
        min-height: 40px;
        box-shadow: none;

        .dingwei {
            position: absolute;
            right: 0;
            height: 40px;
            width: 120px;
        }

    }

    .cont-data {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 30px;
    }

    .elbutton {
        background-color: #1575ee;
        border-radius: 10px;
        border: none;
    }

    .elxyb {
        width: 380px;
        height: 50px;
        margin-top: 50px;
    }

    .newtt {
        width: 118px;
        height: 40px;
        margin-top: 20px;
        border-radius: 8px !important;
    }

    .toptext {
        font-size: 18px;
        line-height: 36px;
        color: #333333;
    }

    .bottom-zhu {
        font-size: 13px;
        color: #C1C7D0;
        text-align: center;
        position: relative;
        top: -30px;

    }

    .detacnt {
        background-color: #fff;
        border: 1px solid #eee;
        border-radius: 10px;
        margin-top: 30px;
        height: 590px;

        .detacnttop {
            font-size: 13px;
            padding: 20px;
            border-bottom: 1px solid #eee;
        }

        .cntdata {
            padding: 10px;

            .cz {
                color: #1575ee;
                display: flex;
                justify-content: space-between;
            }
        }
    }

    .username {
        display: flex;
        margin: 20px;
        font-size: 14px;
    }

    .cntadataright {
        display: flex;
        flex: 1;
        justify-content: space-between;
        border-bottom: 1px solid #C1C7D0;
        color: #333333;
        height: 40px;
    }

    .cntadataleft {
        flex: 0.2;
        color: #C1C7D0;
    }

    .shubi {
        cursor: pointer;
    }

    .changjiantan {
        .el-dialog__footer {
            text-align: center;
        }

        .el-button {
            width: 103px;
            height: 40px;
        }

        .tijiaobut {
            background-color: #1575ee;
        }
    }

    .newtop {
        padding: 10px;

        .el-dialog__header {
            display: none;
        }

        .shenhe {
            color: #C1C7D0;
            font-size: 14px;
        }

    }

    .toptxt {
        text-align: center;
        font-size: 18px;
        font-weight: 400;
        color: #333333;
    }

    .el-form-item {
        align-items: center;
    }

    .jsxx {
        height: 40px;
        display: flex;
        margin-top: 30px;
        margin-bottom: 30px;
        border-bottom: 1px solid #ebebeb;
        font-size: 18px;
    }
}
</style>