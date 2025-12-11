<template>
    <div class="tenantcom">
        <div class="custdatas">
            <div class="prodata">
                <div class="measureds prodatadiv">
                    <div class="prodatatop">
                        组织vip客户
                    </div>
                    <div class="prodatari">
                        <div class="summartxtrig"><span @click="monthout"> &lt;</span>
                            <span style="color: #409EFF;">
                                <el-date-picker v-model="Date3" type="month" placeholder="选择日期" class="xzdata"
                                    :picker-options="pickerOptions0" @change="GetOrganizationCustomerCompareSummary"
                                    :clearable="false">
                                </el-date-picker></span>
                            <span @click="monthadd"> > </span>
                        </div>
                    </div>
                    <div class="measurecnt">
                        <div class="measure">
                            <div class="iconimg">
                                <img src="../../assets/image/data_icon_dayvip.png" alt="" class="measureicon">
                            </div>
                            <div class="prodatatxt">当月新增vip客户数量</div>
                            <div class="prodatanbm" style="color: #3CA1F8;">
                                {{
                                    OrganizationCustomer.month_added_customer_count }}

                            </div>
                            <nav>所有商户</nav>
                        </div>
                        <div class="measure">
                            <div class="iconimg">
                                <img src="../../assets/image/data_icon_daynotvip.png" alt="" class="measureicon">
                            </div>
                            <div class="prodatatxt">同比</div>
                            <div class="prodatanbm" style="color: #9C88FF;"> <span
                                    v-if="OrganizationCustomer.year_on_year == 2147483647">
                                    --
                                </span>
                                <span v-else>
                                    {{ OrganizationCustomer.year_on_year }}%
                                </span>
                            </div>
                            <nav>所有商户</nav>
                        </div>
                        <div class="measure">
                            <div class="iconimg">
                                <img src="../../assets/image/data_icon_vip.png" alt="" class="measureicon">
                            </div>
                            <div class="prodatatxt">环比</div>
                            <div class="prodatanbm" style="color: #3CA1F8;">
                                <span v-if="OrganizationCustomer.month_on_month == 2147483647">
                                    --
                                </span>
                                <span v-else>
                                    {{ OrganizationCustomer.month_on_month }}%
                                </span>
                            </div>
                            <nav>所有商户</nav>
                        </div>

                    </div>
                </div>
                <div class="measureds prodatadiv">
                    <div class="prodatatop">
                        组织测量次数
                    </div>
                    <div class="prodatari">
                        <div class="summartxtrig"><span @click="monthouts"> &lt;</span>
                            <span style="color: #409EFF;"><el-date-picker v-model="Date4" type="month"
                                    placeholder="选择日期" class="xzdata" :picker-options="pickerOptions0"
                                    @change="GetOrganizationMonthlyReportCount" :clearable="false">
                                </el-date-picker></span>
                            <span @click="monthadds"> > </span>
                        </div>
                    </div>
                    <div class="measurecnt">
                        <div class="measure">
                            <div class="iconimg">
                                <img src="../../assets/image/data_icon_dayvip.png" alt="" class="measureicon">
                            </div>
                            <div class="prodatatxt">当月测量次数</div>
                            <div class="prodatanbm" style="color: #3CA1F8;">
                                {{ tenantmeaserement.month_added_customer_count }}
                            </div>
                            <nav>所有商户</nav>
                        </div>
                        <div class="measure">
                            <div class="iconimg">
                                <img src="../../assets/image/data_icon_daynotvip.png" alt="" class="measureicon">
                            </div>
                            <div class="prodatatxt">同比</div>
                            <div class="prodatanbm" style="color: #9C88FF;">
                                <span v-if="tenantmeaserement.year_on_year == 2147483647">
                                    --
                                </span>
                                <span v-else>
                                    {{ tenantmeaserement.year_on_year }}%
                                </span>
                            </div>
                            <nav>所有商户</nav>
                        </div>
                        <div class="measure">
                            <div class="iconimg">
                                <img src="../../assets/image/data_icon_vip.png" alt="" class="measureicon">
                            </div>
                            <div class="prodatatxt">环比</div>
                            <div class="prodatanbm" style="color: #3CA1F8;">
                                <span v-if="tenantmeaserement.month_on_month == 2147483647">
                                    --
                                </span>
                                <span v-else>
                                    {{ tenantmeaserement.month_on_month }}%
                                </span>
                            </div>
                            <nav>所有商户</nav>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="custtop">
            <div class="ownerright" style="margin-top: 15px;">
                <el-date-picker v-model="Date1" type="month" placeholder="选择日期" class="xzdata"
                    :picker-options="pickerOptions0" @change="changemonthh" :clearable="false">
                </el-date-picker>
                <el-date-picker v-model="Date2" type="month" placeholder="选择日期" class="xzdata"
                    :picker-options="pickerOptions1" @change="changeMonth" :clearable="false">
                </el-date-picker>
                <el-input v-model="tenantName" placeholder="商户名称" class="ownerinp"
                    @blur="tenantName = $event.target.value.trim()" :maxlength="20"></el-input>
                <div class="mb-2 ml-4">
                    <el-radio-group v-model="radio1">
                        <el-radio label="1">查询vip客户数量</el-radio>
                        <el-radio label="2">查询测量次数</el-radio>
                    </el-radio-group>

                </div>
                <el-button type="primary" class="owbtm" @click="dialogconse = true">导出为Excel</el-button>
                <el-button type="primary" class="owbtm" @click="SearchCustomerLastStatus">查询</el-button>
            </div>
        </div>
        <div class="custbody">
            <el-table ref="liecheck" :data="datalist" style="width: 100%">
                <el-table-column prop="date" label="月份" v-if="radio3 == '1'">
                    <template #default="scope">
                        {{ fermitTime(scope.row.date).substring(0, fermitTime(scope.row.date).length - 3) }}
                    </template>
                </el-table-column>
                <el-table-column prop="date" label="月份" v-else>
                    <template #default="scope">
                        {{ fermitTime(scope.row.time).substring(0, fermitTime(scope.row.time).length - 3) }}
                    </template>
                </el-table-column>
                <el-table-column prop="tenant_name" label="商户名称">
                </el-table-column>
                <el-table-column prop="month_count" label="当月vip客户数量/测量次数" v-if="radio3 == '1'" />
                <el-table-column prop="monthly_customer_measurement_count" label="当月vip客户数量/测量次数" v-else />
                <el-table-column prop="year_on_year" label="同比">
                    <template #default="scope">
                        <span v-if="scope.row.year_on_year == 2147483647">
                            --
                        </span>
                        <span v-else>
                            {{ scope.row.year_on_year }}%
                        </span>
                    </template>
                </el-table-column>
                <el-table-column prop="month_on_month" label="环比">
                    <template #default="scope">
                        <span v-if="scope.row.month_on_month == 2147483647">
                            --
                        </span>
                        <span v-else>
                            {{ scope.row.month_on_month }}%
                        </span>
                    </template>
                </el-table-column>
            </el-table>
            <div class="fenye">
                <el-pagination v-model:current-page="offset" v-model:page-size="pageSize4"
                    :page-sizes="[10, 30, 50, 100]" :small="true" :background="background"
                    layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange"
                    @current-change="handleCurrentChange" />
            </div>
        </div>
        <el-dialog v-model="dialogconse" title="导出" width="30%">
            <div class="diaurl">
                <div class="urltop">
                    选择月份： <el-date-picker v-model="Date5" type="month" placeholder="选择日期" class="xzdata"
                        :picker-options="pickerOptions0" :clearable="false">
                    </el-date-picker>
                </div>
            </div>
            <template #footer>
                <span class="dialog-footer-yq">
                    <el-button @click.stop="dialogconse = false" class="quxiaobut" style="height: 40px;">取消</el-button>
                    <el-button type="primary" @click.stop="excelapi" class="owbtm">
                        确定
                    </el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import { settopname } from '../../utils/sethometop'
import { elmessage } from '@/utils/popup'
import { getTenantId, getSymptomKeyMap } from '@/utils/storage'
import { GetOrganizationCustomerCompareSummaryRequest, SearchOrganizationTenantCompareRequest, ListOrganizationTenantMonthlyReportCountRequest, GetOrganizationMonthlyReportCountRequest, SearchOrganizationTenantDownloadRequest } from '@/api/api'
import { fermitTime } from '@/utils/fermitTime'
import { exportExcel } from '@/utils/exportExcel.js'

//商户名称
const tenantName = ref('')
//搜索商户输入框
const sminpt = ref('')
//默认每页几条数据
const pageSize4 = ref(10)
const dialogconse = ref(false)
const datasco = ref(true)
//默认总数据
const total = ref(0)
//分页偏移量
const offset = ref(1)

const radio1 = ref('1')
const radio3 = ref('1')
const OrganizationCustomer = ref(
    {
        month_added_customer_count: 0,
        year_on_year: 0,
        month_on_month: 0,
    }
)
const tenantmeaserement = ref(
    {
        month_added_customer_count: 0,
        year_on_year: 0,
        month_on_month: 0,
    }
)

const disabledbut = ref(false)
const newdate = ref(new Date())
const Date1 = ref(new Date(newdate.value.getFullYear(), newdate.value.getMonth() - 1, 1))
const Date2 = ref(new Date())
const Date3 = ref(new Date())
const Date4 = ref(new Date())
const Date5 = ref(new Date())
const datalist = ref()

const handleSizeChange = (val) => {
    if (radio1.value == '1') {
        SearchOrganizationTenantCompare()
    } else {
        ListOrganizationTenantMonthlyReportCount()
    }
}
const handleCurrentChange = (val) => {
    offset.value = val
    console.log(`第: ${val}页`)
    if (radio1.value == '1') {
        SearchOrganizationTenantCompare()
    } else {
        ListOrganizationTenantMonthlyReportCount()
    }
}
const changemonthh = () => {
    let date1 = new Date(Date1.value).getTime()
    let date2 = new Date(Date2.value).getTime()
    let currentDate = new Date(Date2.value)
    currentDate.setFullYear(currentDate.getFullYear() - 1);
    if (date2 != '') {
        if (date2 < date1) {
            elmessage('开始时间不能大于结束时间')
            Date1.value = null
        }
        if (date1 < (currentDate.getTime())) {
            elmessage('时间间隔不能超过一年')
            Date1.value = null
        }
        if (Date1.value == null || Date2.value == null) {
            disabledbut.value = true
        } else {
            disabledbut.value = false
        }
    }

}
const changeMonth = (value) => {
    let date = new Date(value);
    let month = (date.getMonth() + 1).toString().padStart(2, '0');
    let year = date.getFullYear();
    let day = new Date(year, month, 0);
    let endDate = new Date(new Date(day.toLocaleDateString()).getTime() + 24 * 60 * 60 * 1000 - 1);
    Date2.value = endDate
    let date1 = new Date(Date1.value).getTime()
    let date2 = new Date(Date2.value).getTime()
    console.log(new Date(Date1.value), new Date(Date2.value));
    let currentDate = new Date(Date1.value)
    currentDate.setFullYear(currentDate.getFullYear() + 1);
    if (date1 != '') {
        if (date2 < date1) {
            elmessage('结束时间不能小于开始时间')
            Date2.value = null
        }
        if (date2 > (currentDate.getTime())) {
            elmessage('时间间隔不能超过一年')
            Date2.value = null
        }
        if (Date1.value == null || Date2.value == null) {
            disabledbut.value = true
        }
        else {
            disabledbut.value = false
        }
    }
}
const GetOrganizationCustomerCompareSummary = async () => {
    const rs = {
        organizationId: getTenantId(),
        date: new Date(Date3.value)
    }
    const data = await GetOrganizationCustomerCompareSummaryRequest(rs)
    OrganizationCustomer.value.month_added_customer_count = data.data.month_added_customer_count === undefined ? '0' : data.data.month_added_customer_count

    OrganizationCustomer.value.year_on_year = data.data.year_on_year === undefined ? '0' : data.data.year_on_year

    OrganizationCustomer.value.month_on_month = data.data.month_on_month === undefined ? '0' : data.data.month_on_month
}
const GetOrganizationMonthlyReportCount = async () => {
    Date4.value = new Date(Date4.value)
    const rs = {
        organizationId: getTenantId(),
        measurementTime: {
            year: Date4.value.getFullYear(),
            month: Date4.value.getMonth() + 1,
            day: Date4.value.getDate(),
        }
    }
    const data = await GetOrganizationMonthlyReportCountRequest(rs)
    tenantmeaserement.value.month_added_customer_count = data.data.monthly_customer_measurement_count === undefined ? '0' : data.data.monthly_customer_measurement_count

    tenantmeaserement.value.year_on_year = data.data.year_on_year === undefined ? '0' : data.data.year_on_year

    tenantmeaserement.value.month_on_month = data.data.month_on_month === undefined ? '0' : data.data.month_on_month
}
const SearchOrganizationTenantCompare = async () => {
    const res = {
        organizationId: getTenantId(),
        tenantName: tenantName.value,
        startTime: Date1.value,
        endTime: Date2.value,
        pagination: {
            offset: (offset.value - 1) * pageSize4.value,
            size: pageSize4.value
        }
    }
    const data = await SearchOrganizationTenantCompareRequest(res)
    datalist.value = data?.data.tenants
    total.value = data?.data.total_count == undefined ? 0 : data?.data.total_count
}
const ListOrganizationTenantMonthlyReportCount = async () => {
    const res = {
        organizationId: getTenantId(),
        tenantName: tenantName.value,
        startTime: Date1.value,
        endTime: Date2.value,
        pagination: {
            offset: (offset.value - 1) * pageSize4.value,
            size: pageSize4.value
        }
    }
    const data = await ListOrganizationTenantMonthlyReportCountRequest(res)
    datalist.value = data?.data.counts
    total.value = data?.data.total_count == undefined ? 0 : data?.data.total_count
}

const SearchCustomerLastStatus = () => {
    if (radio1.value == 1) {
        SearchOrganizationTenantCompare()
        datasco.value = true
        radio3.value = radio1.value
    } else {
        ListOrganizationTenantMonthlyReportCount()
        datasco.value = false
        radio3.value = radio1.value
    }
}

const monthout = () => {
    const data = new Date(Date3.value)
    Date3.value = data.setMonth(data.getMonth() - 1)
    GetOrganizationCustomerCompareSummary()
    console.log(new Date(Date3.value));

}
const monthadd = () => {
    const data = new Date(Date3.value)
    Date3.value = data.setMonth(data.getMonth() + 1)
    GetOrganizationCustomerCompareSummary()
    console.log(new Date(Date3.value));

}
const monthouts = () => {
    const data = new Date(Date4.value)
    Date4.value = data.setMonth(data.getMonth() - 1)
    GetOrganizationMonthlyReportCount()
    console.log(new Date(Date4.value));

}
const monthadds = () => {
    const data = new Date(Date4.value)
    Date4.value = data.setMonth(data.getMonth() + 1)
    GetOrganizationMonthlyReportCount()
    console.log(new Date(Date4.value));

}
const topname = () => {
    settopname([{ name: '概况', url: '/' }, { name: '测量数据详情', url: '' }])
}

const excelapi = async () => {
    let exc_data = [];
    const dateacs = new Date(Date5.value)
    // 获取当前月份的1号0点00分00秒
    let kaisdate = new Date(dateacs.getFullYear(), dateacs.getMonth(), 1, 0, 0, 0);
    // 获取当前月份的最后一天23点59分59秒
    // 下个月的0号即当前月的最后一天
    let enddata = new Date(dateacs.getFullYear(), dateacs.getMonth() + 1, 0, 23, 59, 59);
    const res = {
        organizationId: getTenantId(),
        startTime: kaisdate,
        endTime: enddata,
    }
    const data = await SearchOrganizationTenantDownloadRequest(res)
    exc_data = [['商户名字', '联系人', '电话', '检验次数', 'vip数量']]
    data.data.tenants.forEach(v => {
        v.measurement_count = v.measurement_count === undefined ? '0' : v.measurement_count
        v.customer_count = v.customer_count === undefined ? '0' : v.customer_count
        let arr = [];
        arr.push(v.tenant_name, v.contact_name, v.contact_phone, v.measurement_count, v.customer_count);
        exc_data.push(arr)
        console.log(v);
        
    })
    let tenantdatamonth = dateacs.getMonth() + 1
    if (tenantdatamonth > 12) tenantdatamonth = 1
    console.log(exc_data);
    
    exportExcel(dateacs.getFullYear() + '年' + tenantdatamonth + '月商户数据.', exc_data)
    dialogconse.value = false
}
GetOrganizationCustomerCompareSummary()
GetOrganizationMonthlyReportCount()
SearchOrganizationTenantCompare()
topname()
</script>
<style lang="scss">
.tenantcom {
    .custdatas {
        .summar-cnt {
            display: flex;
            justify-content: space-between;

            .sumdiv {
                width: 50%;
                height: 228px;
                border-radius: 10px;
                margin-left: 20px;
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: space-around;
                font-weight: 400;
                font-style: normal;
                color: #666;
                text-align: center;
                background-size: cover;
                margin-bottom: 30px;
                background-color: #fff;

                .summar-txt {
                    width: 100%;
                    display: flex;
                    justify-content: space-between;
                    font-size: 18px;
                    font-weight: 650;
                    border-bottom: 1px solid #EEEEEE;
                    padding-bottom: 20px;

                    .summartxtlf {
                        margin-left: 30px;
                    }
                }
            }

            .summar-nab {
                font-size: 14px;
                display: flex;
                justify-content: space-around;
                width: 100%;
            }
        }
    }

    .prodata {
        margin-top: 30px;
        display: flex;
        justify-content: space-between;

        .measureds {
            width: 48%;
        }

        .prodatadiv {
            display: flex;
            flex-direction: column;
            border: 1px solid rgba(238, 238, 238, 1);
            background-color: #fff;
            font-style: normal;
            font-weight: 400;
            border-radius: 10px;
            position: relative;
            height: 228px;

            .summartxtrig {
                margin-right: 30px;
                font-size: 13px;
                display: flex;
                justify-content: space-around;

                span {
                    display: inline-block;
                }

                .el-input__wrapper {
                    box-shadow: none
                }

                .el-input__prefix {
                    display: none;
                }

                .el-input {
                    width: 80px;
                    position: relative;
                    top: -5px
                }
            }

            .prodatatop {
                height: 70px;
                line-height: 70px;
                font-weight: 700;
                font-size: 18px;
                padding-left: 40px;
                border-bottom: 1px solid rgba(238, 238, 238, 1);
            }

            .prodatari {
                font-weight: 650;
                font-style: normal;
                color: #1575EE;
                position: absolute;
                right: 10px;
                top: 25px;
                width: 30%;
                cursor: pointer;
            }

            .clientcnt,
            .measurecnt {
                display: flex;
                flex: 1;
                justify-content: space-around;

                .measure,
                .clientdiv {
                    display: flex;
                    flex-direction: column;
                    align-items: center;
                    justify-content: center;

                    .iconimg {
                        height: 30px;
                        line-height: 30px;
                    }
                }

                .clienticon {
                    width: 30px;
                    height: 27px;
                }

                .measureicon {
                    width: 12px;
                    height: 12px;
                }

                .prodatatxt {
                    margin-top: 10px;
                    font-size: 13px;
                }

                .prodatanbm {
                    font-weight: 700;
                    font-style: normal;
                    font-size: 30px;

                    text-align: center;
                }
            }
        }
    }

    nav {
        font-weight: 400;
        font-style: normal;
        font-size: 12px;
        color: #C1C7D0;
        text-align: center;
    }

    .ownerright {
        display: flex;
        flex: 0.5;

        .ownerinp {
            width: 240px;
            height: 40px;
            margin-right: 10px;

            .el-input--suffix {
                height: 40px;
            }
        }

        .xzdata {
            width: 240px;
            height: 40px;
            margin-right: 10px;
        }
    }

    .owbtm {
        background-color: rgba(21, 117, 238, 1);
        height: 40px;
        padding: 0 25px;
    }

    .custbody {
        margin-top: 30px;
        background-color: #fff;
        min-height: 60vh;
        padding: 20px;
        padding-top: 40px;
        position: relative;
    }

    .fenye {
        display: flex;
        justify-content: flex-end;
        position: absolute;
        bottom: -50px;
        right: 0;
    }

    .mb-2 {
        display: flex;

        .el-radio-group {
            display: inline-flex;
            flex-wrap: wrap;
            font-size: 0;
            flex-direction: column;
            align-items: flex-start;

            .el-radio {
                height: 20px;
            }

            .el-radio__label {
                font-size: 13px;
                color: #333;
            }
        }
    }
}
</style>