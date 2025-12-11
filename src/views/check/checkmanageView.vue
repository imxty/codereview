<template>
    <div class="check">
        <div class="checkmatop">
            <div class="checkmalf">
                <div class="checkdate">
                    <div>结算周期：{{ reviewtype.date.substring(0, 7).replaceAll('-',
                        '') }}01-{{ reviewtype.date.substring(0, 7).replaceAll('-', '')
                        }}{{ getLastDayOfNaturalMonth(reviewtype.date).getDate() }}</div>
                </div>
                <div>
                    <el-button type="primary" class="owbtm" @click="Monthlyendexcil">下载明细</el-button>
                </div>
            </div>

        </div>
        <div class="checkmacnt">
            <div class="checcnttop">
                <span>激活商户账单详情</span>
            </div>
            <div class="checcntdata">
                <el-table :data="tableData.subcriptions" style="width: 100%" :row-click="openDetails">
                    <el-table-column prop="tenant_name" label="商户名">
                    </el-table-column>
                    <el-table-column prop="tenant_id" label="商户ID">
                    </el-table-column>
                    <el-table-column prop="years" label="开通年数">
                    </el-table-column>
                    <el-table-column prop="created_at" label="开通日期">
                        <template #default="scope">
                            <span>
                                {{ getoneyear(scope.row.created_at, +0) }}
                            </span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="dayprice" label="到期日期">
                        <template #default="scope">
                            <span>
                                {{ getoneyear(scope.row.expired_at, +0) }}
                            </span>
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
import { reactive, ref } from 'vue';
import { settopname } from '@/utils/sethometop'
import { exportExcel } from '@/utils/exportExcel.js'
import { getTenantId } from '@/utils/storage'
import {
    GetBillingRequest, SearchTenantSubscriptionsPaginationRequest, SearchTenantSubscriptionsRequestRequest
} from '@/api/api'
import { getoneyear, getLastDayOfNaturalMonth } from '../../utils/fermitTime'

//申请发票开关
const limitation = ref(false)
//默认总数据
const total = ref(0)
const pageSize4 = ref(10)
const offset = ref(1)
const reviewtype = ref(JSON.parse(history.state?.keyword))
console.log(reviewtype.value);
const handleSizeChange = (val) => {
    GetBilling()
}
const handleCurrentChange = (val) => {
    offset.value = val
    console.log(`第: ${val}页`)
    GetBilling()
}
//账单数据
const tableData = ref([])


const topname = () => {
    settopname([{ name: '账单管理', url: '/billList' }, { name: '续费详情', url: '' }])
}

const GetBilling = async () => {
    const rs = {
        organizationId: getTenantId(),
        startTime: new Date(reviewtype.value.date),
        endTime: new Date(getLastDayOfNaturalMonth(reviewtype.value.date)),
        pagination: {
            offset: (offset.value - 1) * pageSize4.value,
            size: pageSize4.value
        }
    }
    const data = await SearchTenantSubscriptionsPaginationRequest(rs)
    tableData.value = data.data
    total.value = data.data.total_count == undefined ? 0 : data.data.total_count

}

//导出月结数据为excil表格
const Monthlyendexcil = async (data) => {
    let exc_data = [];
    const billinglist = await SearchTenantSubscriptionsRequestRequest({
        organizationId: getTenantId(),
        startTime: new Date(reviewtype.value.date),
        endTime: new Date(getLastDayOfNaturalMonth(reviewtype.value.date)),
    })

    exc_data = [['商户名', '商户ID', '开通年数', '开通日期', '到期日期']]
    let lenght = 0
    billinglist.data.subcriptions.forEach(v => {
        console.log(v);
        lenght++
        let list = []
        list.push(v.tenant_name)
        list.push(v.tenant_id)
        list.push(v.years)
        list.push( getoneyear(v.start_time, +0))
        list.push( getoneyear(v.end_time, +0))
        exc_data.push(list)
    });
    console.log(billinglist);
    const billingdate = new Date(reviewtype.value.date)
    exportExcel(billingdate.getFullYear() + '年' + (billingdate.getMonth() + 1) + '月' + lenght + '个商户', exc_data)
}
topname()
GetBilling()
</script>
<style lang="scss">
.check {
    .owbtm {
        background-color: rgba(21, 117, 238, 1);
        height: 40px;
        padding: 0 25px;
    }

    .checkmatop {
        display: flex;
        justify-content: space-between;
        background: #fff;
        align-items: center;
        padding: 40px;
        border: 1px solid #eeeeee;
        border-radius: 10px;

        .checkmalf {
            display: flex;
            flex: 0.9;
            justify-content: space-between;
            align-items: center;
            height: 81px;

            .checkam,
            .checkdate {
                height: 81px;
                display: flex;
                flex-direction: column;
                justify-content: space-around;
            }

            .checkam {
                color: #999999;
                font-size: 13px;

                .amount {
                    font-size: 22px;
                    color: #333333;
                }
            }

            .checkdate {
                font-size: 14px;
                flex: 0.8;

            }
        }
    }

    .checkmacnt {
        margin-top: 40px;
        background-color: #fff;
        border: 1px solid #eeeeee;
        border-radius: 10px;
        position: relative;

        .fenye {
            display: flex;
            justify-content: flex-end;
            position: absolute;
            bottom: -50px;
            right: 0;
        }

        .checcnttop {
            font-size: 18px;
            color: #333333;
            font-weight: 650;
            padding: 20px 30px;
            border-bottom: 1px solid #eeeeee;
        }

        .checcntdata {
            padding: 20px;
        }
    }

    .limitation {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 20px;
    }
}

/*选中某行时的背景色*/
.el-table body tr.current-row>td {
    color: #28A458;
    background: rgb(197, 213, 255) !important;
}
</style>