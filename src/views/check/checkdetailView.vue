<template>
    <div class="checkdeta">
        <div class="detatop">
            <div class="topdate">
                <div>
                    <el-date-picker v-model="value1" type="month" placeholder="选择日期" class="xzdata"
                        :picker-options="pickerOptions0" @change="changemonthh" :clearable="false">
                    </el-date-picker>
                    <el-date-picker v-model="value2" type="month" placeholder="选择日期" class="xzdata"
                        :picker-options="pickerOptions1" @change="changeMonth" :clearable="false">
                    </el-date-picker>
                    <el-button type="primary" class="cxbut" :disabled="disabledbut" @click="getlist">查询</el-button>
                </div>
            </div>
        </div>
        <div class="detacnt">
            <div class="detacnttopbor">
                <div class="detacnttop">
                    <div :class="[{ xuanzhong: zdtable == '月结' }]" class="detatopnm" style="cursor: pointer;"
                       >续期明细 <span class="detacnttopbtm"></span>
                    </div>
                </div>
            </div>
            <div class="cntdata">
                <el-table ref="liecheck" :data="checkList" style="width: 100%"
                    @selection-change="handleSelectionChange" :row-key="getRowKeys" @select="selectTable"
                    class="no-multiple" v-if="zdtable != '发票'">
                    <el-table-column prop="month" label="月份">
                        <template #default="scope">
                            {{ new Date(scope.row.date).getFullYear() }}-<span
                                v-if="(new Date(scope.row.date).getMonth() + 1) < 10">0{{ new
                                    Date(scope.row.date).getMonth() + 1 }}
                            </span>
                            <span v-else>
                                {{ new Date(scope.row.date).getMonth() + 1 }}
                            </span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="content" label="续期商户数">
                        <template #default="scope">
                            商户激活 {{ scope.row.tenant_count }} 个
                        </template>
                    </el-table-column>
                    <el-table-column prop="days" label="合计年限">

                        <template #default="scope">
                            <span>
                                {{ scope.row.year_count  }}年
                            </span>
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="130">
                        <template #default="scope">
                            <div class="cz">
                                <div @click="fanhui(scope.row)" style="cursor: pointer;">详细</div>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
            <div class="fenye">
                <el-pagination v-model:current-page="offset" v-model:page-size="pageSize4"
                    :page-sizes="[10, 30, 50, 100]" :small="true" :background="background"
                    layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange"
                    @current-change="handleCurrentChange" />
            </div>
        </div>
    </div>
</template>

<script setup>
import { reactive, ref, onUnmounted } from 'vue';
import { getTenantId } from '@/utils/storage'
import {
    SearchBillingRequest,
    PayBillingRequest,
    CheckBillingPayStatusRequest,
    GetBillingRequest,
    ListOrganizationInvoicesRequest,
    SearchTenantSubscriptionsPaginationRequest,
    GetOrganizationSubscriptionsMonthlyRequest
} from '@/api/api'
import { settopname } from '@/utils/sethometop'
import { useRouter } from 'vue-router';
import QRCode from 'qrcode'
import { ElMessage } from 'element-plus';
import { getDateByDays } from '@/utils/fermitTime'
import { getoneyear } from '../../utils/fermitTime'
import { useStore } from 'vuex'
import { elmessage } from '@/utils/popup'
import data from 'china-area-data';

const $store = useStore();

const $router = useRouter();
//默认总数据
const total = ref(0)
//分页偏移量

const disabledbut = ref(false)
const zdtable = ref('月结')
const offset = ref(1)
const newdate = ref(new Date())
const oldDate = ref(new Date(newdate.value.getFullYear(), newdate.value.getMonth() + 1, 1))
const value1 = ref(new Date(newdate.value.getFullYear(), newdate.value.getMonth() - 1, 1))
const value2 = ref(new Date())
const getIndex = ref(null)
const billsworld = ref(null)
const pageSize4 = ref(10)

const multipleSelection = ref([])
const liecheck = ref([])
//二维码弹窗
const pljh = ref(false)
//支付账单的数据
const PayBillingdata = ref([])


const billlistdata = ref(0)
const prestore = ref(false)
//支付按钮是否可用
const disabled = ref(true)
const pickerOptions0 = ref((date) => {
    return false
})
const pickerOptions1 = ref()
const handleSizeChange = (val) => {
    getlist()
}
const handleCurrentChange = (val) => {
    offset.value = val
    console.log(`第: ${val}页`)
    getlist()
}
const checkList = ref([])

const selectTable = (selection, row) => {
    console.log(selection);
    multipleSelection.value = [row]
    billsworld.value = row
    console.log(row);
    getIndex.value = row.index
    if (selection.length > 0) {
        disabled.value = false
    } else {
        disabled.value = true
    }
}

const callback = (row) => {
    return row.billing_status != 'BILLING_STATUS_SETTLEMENT'
}
const handleSelectionChange = (val) => {
    if (val.length > 1) {
        liecheck.value.clearSelection()
        liecheck.value.toggleRowSelection(val.pop())
    }
}
const getRowKeys = (row) => {
    return row.tenant_id;
}
const getlist = async () => {
    const rs = {
        organizationId: getTenantId(),
        startTime: value1.value,
        endTime: value2.value,
        pagination: {
            offset: (offset.value - 1) * pageSize4.value,
            size: pageSize4.value
        }
    }
    const data = await GetOrganizationSubscriptionsMonthlyRequest(rs)
    checkList.value = data.data.subscriptions
    total.value = checkList.value.total_count == undefined ? 0 : checkList.value.total_count


}
const topname = () => {
    settopname([{ name: '账单管理', url: '' }])
    $store.commit('setactiveMenuKey', 4)
}
const fanhui = (data) => {
    $router.push({
        name: 'check',
        state: {
            keyword: JSON.stringify(data)
        }
    })

}
const changeMonth = (value) => {
    let date = new Date(value);
    let month = (date.getMonth() + 1).toString().padStart(2, '0');
    let year = date.getFullYear();
    let day = new Date(year, month, 0);
    let endDate = new Date(new Date(day.toLocaleDateString()).getTime() + 24 * 60 * 60 * 1000 - 1);
    value2.value = endDate
    let date1 = new Date(value1.value).getTime()
    let date2 = new Date(value2.value).getTime()
    console.log(new Date(value1.value), new Date(value2.value));
    if (date1 != '') {
        if (date2 < date1) {
            elmessage('结束时间不能小于开始时间')
            value2.value = null
        }
        if (value1.value == null || value2.value == null) {
            disabledbut.value = true
        }
        else {
            disabledbut.value = false
        }
    }
}
const changemonthh = () => {
    let date1 = new Date(value1.value).getTime()
    let date2 = new Date(value2.value).getTime()
    if (date2 != '') {
        if (date2 < date1) {
            elmessage('开始时间不能大于结束时间')
            value1.value = null
        }
        if (value1.value == null || value2.value == null) {
            disabledbut.value = true
        } else {
            disabledbut.value = false
        }
    }

}
const tableRowClassName = ({ row, rowIndex }) => {
    //把每一行的索引放进row
    row.index = rowIndex;
}
const selectedHighlight = ({ row, rowIndex }) => {
    if ((getIndex.value) === rowIndex) {
        return {
            "background-color": "#CAE1FF"
        };
    }
}


topname()
getlist()
changeMonth(value2.value)

</script>

<style lang="scss">
.checkdeta {
    .commers {
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
        }

        .pljh {
            color: #3CA1F8;
        }
    }

    .xmtop {
        width: 330px;
        display: flex;
        position: fixed;
        top: 24px;
        left: 400px;
        z-index: 999;
        font-size: 13px;
        color: #999999;
        justify-content: space-between;
        align-items: flex-end;

        .yfk {
            font-size: 13px;
            line-height: 10px;
            color: #1575EE;
            cursor: pointer;
        }

        .tian {
            font-size: 22px;
            color: #1575EE;
        }
    }

    .fenye {
        display: flex;
        justify-content: flex-end;
        position: absolute;
        bottom: -50px;
        right: 0;
    }

    .detatop {
        .toptips {
            font-size: 14px;
            color: rgb(102, 102, 102);
            margin-bottom: 20px;

            .zdts {
                width: 10px;
                height: 10px;
                display: inline-block;
                background-color: #957ef9;
                border-radius: 10px;
                margin-right: 5px;
            }
        }

        .xzdata {
            width: 240px;
            height: 40px;
            margin-right: 10px;
        }

        .cxbut {
            width: 80px;
            height: 40px;
            background-color: #1575ee;
        }
    }

    .detacnt {
        background-color: #fff;
        border: 1px solid #eee;
        border-radius: 10px;
        margin-top: 30px;
        min-height: 60vh;
        position: relative;
        padding-bottom: 40px;

        .detacnttop {
            width: 440px;
            display: flex;
            font-size: 18px;
            font-weight: 700;
            padding-left: 10px;
            color: #C1C7D0;
            justify-content: space-between;

            div {
                cursor: default;
                padding: 10px 25px;
            }

            .detatopnms {
                position: relative;

                .detacnttopbtm {
                    display: block;
                    width: 50px;
                    height: 2px;
                    background-color: #1575EE;
                    position: absolute;
                    bottom: 0px;
                    left: 60px;
                }
            }

            .detatopnmss {
                position: relative;

                .detacnttopbtm {
                    display: block;
                    width: 50px;
                    height: 2px;
                    background-color: #1575EE;
                    position: absolute;
                    bottom: 0px;
                    left: 20px;
                }
            }

            .detatopnm {
                position: relative;

                .detacnttopbtm {
                    display: block;
                    width: 50px;
                    height: 2px;
                    background-color: #1575EE;
                    position: absolute;
                    bottom: 0px;
                    left: 36px;
                }
            }

            .xuanzhong {
                font-size: 18px;
                color: #1575EE;
            }
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

    .detabtm {
        display: flex;
        justify-content: flex-end;
        align-items: center;
        height: 80px;

        .btmlf {
            display: flex;
            color: #333333;
            font-size: 13px;
            align-items: center;
            justify-content: space-around;
            flex: 0.4;

            .zk {
                width: 210px;
                height: 40px;
            }
        }

        .btmrg {
            display: flex;
            flex: 0.5;
            justify-content: flex-end;
            align-items: center;

            .amount {
                height: 80px;
                font-size: 13px;
                display: flex;
                flex-direction: column;
                justify-content: space-around;

                .amount-yuan {
                    color: #333333;
                }

                .mount {
                    font-size: 22px;
                    color: #1575EE;
                }

                .mountday {
                    font-size: 12px;
                    color: #999999;
                }
            }

            .ljjs {
                width: 100px;
                height: 40px;
                margin-top: 100px;
                background-color: #1575ee;
            }
        }
    }


    .commerdg>div {
        margin-top: 10px;
        font-size: 14px;

    }

    .detacnttopbor {
        border: 1px solid #eee;
    }

    .topdate {
        display: flex;
        justify-content: space-between;
    }
}

//关闭表头全选按钮
.el-table .el-table__header-wrapper .el-checkbox {
    display: none;
}

.xianzong {
    background-color: #91c9fc !important;
    border: none;
}

.weixuan {
    width: 13px;
    height: 13px;
    border: #1575EE 1px solid;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #1575EE;
    border-radius: 3px;
    cursor: not-allowed;
}
</style>