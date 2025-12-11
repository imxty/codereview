<template>
    <div class="prepayme">
        <div class="kaitong">
            <div>
                <div class="kaitong-top">
                    <div class="bigtext">
                        开通年限（年）
                    </div>
                    <div class="divlr">
                        <el-input-number @keydown="handleKeydown" v-model="num" :min="1" :max="3" @change="handleChange" />
                    </div>
                    <div class="byfy">
                        本月已产生的费用，将在下个月账单中进行结算
                    </div>
                </div>
                <div>
                    <div class="dyts" v-if="reviewtype.length > 20">
                        <div>· 提交申请后，将由相关业务经理联系组织负责人进行费用沟通</div>
                        <div>
                            · 申请期间，未逾期的预付款商户可继续使用，新开通的用户将停用，申请通过后可正常进行使用
                        </div>
                    </div>
                    <div>
                        <el-button type="primary" class="cxbut" @click="SubmitTenantPrestoreAccount"
                            v-if="reviewtype.length > 20">立即申请</el-button>
                    </div>
                </div>
            </div>
            <div class="zhangdan" v-if="reviewtype.length <= 20">
                <div class="zdtext">
                    <div class="bigtext">
                        账单金额（元）
                    </div>
                    <div class="divlr zdjs">
                        <div>开通年份、商户数越多，价格越优惠</div>
                        <div>开通商户数超过20个，可与客服沟通获取优惠报价方案</div>
                        <div>原价：6.66元/天/商户</div>
                        <div>现价：<span style="color: #1575EE;">{{ ((money /
                            10000) / 365 / num / reviewtype.length).toFixed(2) }}</span>元/天/商户</div>
                    </div>
                </div>
                <div class="jiesuan">
                    <div class="jiner">
                        ¥ {{ money / 10000 }}
                    </div>
                    <div class="divlr">
                        <el-button type="primary" @click="PayTenantPrestoreAccount" class="ljjs">立即结算</el-button>
                    </div>
                </div>
            </div>
        </div>

        <div class="detacnt">
            <div>
                <div class="detacnttop">
                    开通预付款商户列表
                </div>
                <div class="cntdata">
                    <el-table :data="reviewtype" style="width: 100%">
                        <el-table-column prop="name" label="商户名称">
                        </el-table-column>
                        <el-table-column prop="tenant_id" label="商户ID">
                        </el-table-column>
                        <el-table-column prop="settle" label="开通日期">
                            <template #default="scope">
                                <span
                                    v-if="scope.row.tenant_status == 'TENANT_STATUS_PRESTORE' && scope.row.timeline.expired != true">{{
                                        fermitTime(scope.row.timeline.end_time, true) }}</span>
                                <span v-else>{{ fermitTime(new Date(), true) }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="day" label="到期日期">
                            <template #default="scope">
                                <span
                                    v-if="scope.row.tenant_status == 'TENANT_STATUS_PRESTORE' && scope.row.timeline.expired != true">{{
                                        datejia(fermitTime(scope.row.timeline.end_time, true), num) }}</span>
                                <span v-else>{{ datejia(fermitTime(new Date(), true), num) }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="amount" label="状态">
                            <template #default="scope">
                                <div v-if="scope.row.tenant_status == 'TENANT_STATUS_PRESTORE'" class="daishen">续费</div>
                                <div v-else class="shenhezhong">开通</div>
                            </template>
                        </el-table-column>
                        <el-table-column label="操作" width="130">
                            <el-button type="danger" plain @click="deletelist(scope)">删除</el-button>
                        </el-table-column>
                    </el-table>
                </div>
            </div>
        </div>
        <div class="commerss">
            <el-dialog v-model="pljh" title="打开微信扫一扫二维码" width="30%" @close="closetime">
                <div class="commerdg" v-loading="loading" style="text-align: center;">
                    <img :src="qrCodeImage" alt="">
                </div>
                <div style="text-align: center;">应付金额：<span style="color: #F56C6C;text-align: center;">¥ {{ money / 10000
                }}</span>元</div>
            </el-dialog>
        </div>
    </div>
</template>
 
<script setup>
import { ref, reactive, onUnmounted } from 'vue'
import { GetPrestorePriceRequest, SubmitTenantPrestoreAccountRequest, PayTenantPrestoreAccountRequest, CheckBillingPayStatusRequest } from '@/api/api'
import { getTenantId } from '@/utils/storage'
import { ElMessage } from 'element-plus';
import { useRouter } from 'vue-router';
import { fermitTime, getNextMonth } from '../../utils/fermitTime'
import QRCode from 'qrcode'
import { settopname } from '../../utils/sethometop'

const $router = useRouter();
const loading = ref(true)
const reviewtype = ref(JSON.parse(history.state.keyword))
console.log(reviewtype.value);
const money = ref('0.00')
const enddate = ref('')
const pljh = ref(false)
const qrCodeImage = ref('')
const billingId = ref(null)
const timestop = ref()
console.log(reviewtype);
const num = ref(1)
const handleChange = (value) => {
    console.log(value)
    getNext()
    GetPrestorePrice()
}

const GetPrestorePrice = async () => {
    const data = []
    reviewtype.value.forEach(e => {
        data.push(e.tenant_id)
    });
    const rs = {
        organizationId: getTenantId(),
        tenantIds: data,
        years: num.value
    }
    const mory = await GetPrestorePriceRequest(rs)
    money.value = mory.data.account
}
const SubmitTenantPrestoreAccount = async () => {
    const data = []
    reviewtype.value.forEach(e => {
        data.push(e.tenant_id)
    });
    const rs = {
        organizationId: getTenantId(),
        tenantIds: data,
        years: num.value
    }
    const ts = await SubmitTenantPrestoreAccountRequest(rs)
    if (ts.status === 200) {
        ElMessage({
            message: '提审成功',
            type: 'success',
        })
        window.history.go(-1)
    } else {
        console.log(ts);
        ElMessage({
            showClose: true,
            message: ts.data.detail,
            type: 'error',
        })
    }
}
const getNext = () => {
    enddate.value = getNextMonth(num.value)
    console.log(enddate.value);
}

const PayTenantPrestoreAccount = async () => {
    pljh.value = true
    const datas = []
    reviewtype.value.forEach(e => {
        datas.push(e.tenant_id)
    })
    const rs = {
        organizationId: getTenantId(),
        tenantIds: datas,
        years: num.value
    }
    const data = await PayTenantPrestoreAccountRequest(rs)
    qrCodeImage.value = await QRCode.toDataURL(data.data.pay_url)
    billingId.value = data.data.billing_id
    loading.value = false
    timestop.value = setInterval(async function () {
        console.log(billingId.value);
        const checkstatus = await CheckBillingPayStatusRequest({ billingId: billingId.value })
        if (checkstatus.data.success) {
            clearInterval(timestop.value);
            ElMessage({
                message: '支付成功',
                type: 'success',
            })
            pljh.value = false
            window.history.go(-1)
        }
    }, 5000)
    console.log(timestop.value);
    console.log(data);
}
const topname = () => {
    settopname([{ name: '商户管理', url: '/merchantManagement' }, { name: '开通预付款', url: '' }])
}

const deletelist = (i) => {
    console.log(i);
    reviewtype.value.splice(i, 1)
    if (reviewtype.value.length == 0) {
        window.history.go(-1)
    }
    if (reviewtype.value.length <= 20) {
        GetPrestorePrice()
    }
}
const closetime = () => {
    clearInterval(timestop.value)
}
onUnmounted(() => {
    clearInterval(timestop.value)
})
function datejia(time, num) {
    var timeFlagt = new Date(time)
    timeFlagt.setDate(timeFlagt.getDate() + 365 * num)
    var ss =
        (timeFlagt.getFullYear()) +
        "-" +
        (timeFlagt.getMonth() + 1 < 10
            ? "0" + (timeFlagt.getMonth() + 1)
            : timeFlagt.getMonth() + 1) +
        "-" +
        (timeFlagt.getDate() < 10 ? "0" + timeFlagt.getDate() : timeFlagt.getDate())
    return ss;
}
function handleKeydown(event) {
    if (event.key !== "ArrowUp" && event.key !== "ArrowDown") {
        event.preventDefault();
    }
}
topname()
getNext()
GetPrestorePrice()
</script> 
<style lang="scss">
.prepayme {
    .kaitong {
        width: 100%;
        background-color: #fff;
        border-radius: 10px;
        padding: 20px;

        .commerss {
            .el-dialog {
                border-radius: 10px;
                min-height: 350px;
            }

            .el-dialog__body {
                padding: 10px 20px;
            }

            .commerdg {
                padding-top: 30px;
                border-top: 1px solid #eeeeee;
                text-align: center;
                min-height: 230px;
                display: flex;
            }

            .pljh {
                color: #3CA1F8;
            }
        }

        .kaitong-top {
            display: flex;
            padding-bottom: 30px;
            border-bottom: 1px solid #eee;
            align-items: center;
        }

        .zhangdan {
            margin-top: 30px;

            .zdtext {
                display: flex;

                .zdjs {
                    font-size: 14px;
                    color: #666666;
                }
            }

            .jiesuan {
                display: flex;
                margin-top: 30px;
            }
        }
    }

    .detacnt {
        background-color: #fff;
        border: 1px solid #eee;
        border-radius: 10px;
        margin-top: 30px;
        height: 60vh;

        .detacnttop {
            font-size: 18px;
            font-weight: 700;
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

    .bigtext {
        font-size: 20px;
        color: #333333;
        font-weight: 650;
    }

    .jiner {
        font-size: 30px;
        color: #1575EE;
    }

    .divlr {
        margin-left: 30px;
    }

    .cxbut {
        width: 100px;
        height: 40px;
        background-color: #1575ee;
        border: none;
        margin-top: 30px;
    }

    .ljjs {
        width: 100px;
        height: 40px;
        background-color: #1575ee;
    }

    .dyts {
        font-size: 14px;
        color: #666666;
        line-height: 24px;
    }

    .byfy {
        font-size: 14px;
        color: #F1AA2E;
        background-color: rgba(255, 250, 232, 1);
        padding: 3px;
        margin-left: 10px;
    }



    .commerdg>div {
        margin-top: 10px;
        font-size: 14px;
    }
}
</style>