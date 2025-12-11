<template>
    <div class="custdata">
        <div class="custtop">
            <div class="ownerright">
                <el-input v-model="tenantName" placeholder="商户名称" class="ownerinp"
                    @blur="tenantName = $event.target.value.trim()" :maxlength="20"></el-input>
                <el-select v-model="yaopin" placeholder="请选择" @change="aa" class="ownerinp">
                    <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value">
                    </el-option>
                </el-select>
                <el-input v-model="customerName" placeholder="客户姓名" class="ownerinp"
                    @blur="customerName = $event.target.value.trim()" :maxlength="20"></el-input>
                <el-input v-model="customerPhone" placeholder="客户手机号" class="ownerinp"
                    @blur="customerPhone = $event.target.value.trim()" :maxlength="20"></el-input>

            </div>
            <div class="ownerright" style="margin-top: 15px;">
                <el-date-picker v-model="Date1" type="day" placeholder="选择日期" class="xzdata"
                    :picker-options="pickerOptions0" @change="changemonthh" :clearable="false">
                </el-date-picker>
                <el-date-picker v-model="Date2" type="day" placeholder="选择日期" class="xzdata"
                    :picker-options="pickerOptions1" @change="changeMonth" :clearable="false">
                </el-date-picker>
                <el-select :placeholder="lists" class="ownerinp selectinp" multiple>
                    <el-option class="optioninp">
                        <div class="checkquan">
                            <el-checkbox :indeterminate="dirisIndeterminate" v-model="dircheckAll"
                                @change="dirhandleCheckAllChange">全部</el-checkbox>
                        </div>
                    </el-option>
                    <el-option class="optioninp" v-for="(city, value) in dircities" :key="city">
                        <div class="checkdan">
                            <el-checkbox-group v-model="dircheckedCities" @change="dirhandleCheckedCitiesChange">
                                <el-checkbox :label="value">{{ city
                                    }}</el-checkbox>
                            </el-checkbox-group>
                        </div>
                    </el-option>
                </el-select>
                <el-button type="primary" class="owbtm" @click="SearchReports">搜索</el-button>
            </div>
        </div>
        <div class="custbody">
            <el-table ref="liecheck" :data="ruquest?.reports" style="width: 100%">
                <el-table-column prop="tenant_name" label="商户名称" />
                <el-table-column prop="customerype" label="客户类型">
                    <template #default="scope">
                        <span v-if="scope.row.is_customer">vip客户</span>
                        <span v-else>体验客户</span>
                    </template>
                </el-table-column>
                <el-table-column prop="customer_name" label="客户名称" >
                <template #default="scope">
                    <span v-if="scope.row.customer_name">{{ scope.row.customer_name }}</span>
                    <span v-else>--</span>
                </template>
                </el-table-column>
                <el-table-column prop="customer_phone" label="客户手机号" >
                <template #default="scope">
                    <span v-if="scope.row.customer_phone">{{ scope.row.customer_phone }}</span>
                    <span v-else>--</span>
                </template>
                </el-table-column>
                <el-table-column prop="gender" label="性别">
                    <template #default="scope">
                        <span v-if="scope.row.gender=='GENDER_MALE'">男</span>
                        <span v-else>女</span>
                    </template>
                </el-table-column>
                <el-table-column prop="phone" label="脏腑辨证" >
                <template #default="scope">
                    <!-- {{ scope.row.dirty_dialectics }} -->
                    {{ getmaplist(scope.row.dirty_dialectics) }}
                </template>
                </el-table-column>
                <el-table-column prop="create_time" label="测量时间">
                    <template #default="scope">
                        {{ bjTime(new Date(scope.row.create_time)) }}
                    </template>
                </el-table-column>
                <el-table-column label="操作" :filter-method="filterTag" filter-placement="bottom-end" width="100">
                    <template #default="scope">
                        <el-button type="primary" plain @click="clickreport(scope.row.report_id,scope.row.tenant_id)">查看</el-button>
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
        <el-dialog v-model="dialogFormVisible" width="430">
         <span>
            <iframe :src=reporturl width="100%" height="600px" frameborder="0" v-if="dialogFormVisible"></iframe>
         </span>
      </el-dialog>
    </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import { settopname } from '../../utils/sethometop'
import { elmessage } from '@/utils/popup'
import { getTenantId, getSymptomKeyMap, getKeyMaplist } from '@/utils/storage'
import { SearchReportsRequestRequest } from '@/api/api'
import { bjTime } from '@/utils/fermitTime'

//商户名称
const tenantName = ref('')
//商户名称
const customerName = ref('')
//商户名称
const customerPhone = ref('')
//搜索商户输入框
const sminpt = ref('')
//下拉多选
const lists = ref('')
//默认每页几条数据
const pageSize4 = ref(10)
//默认总数据
const total = ref(0)
//使用中商品
const shiyongz = ref(0)
//脏腑辨证check
const dircheckAll = ref(false)
const dircheckedCities = ref([])
const dircities = ref([])
const dirisIndeterminate = ref(false)
//分页偏移量
const offset = ref(1)
const disabledbut = ref(false)
const newdate = ref(new Date())
const Date1 = ref(new Date(newdate.value.getFullYear(), newdate.value.getMonth() - 1, 1))
const Date2 = ref(new Date())
const ruquest = ref()
const dialogFormVisible = ref(false)
const reporturl = ref('')

const yaopin = ref('CUSTOMER_TYPE_BOTH')
//客户类型
const options = ref([{
    label: '全部',
    value: 'CUSTOMER_TYPE_BOTH'
}, {
    label: 'vip客户',
    value: 'CUSTOMER_TYPE_CUSTOMER'
}, {
    label: '体验客户',
    value: 'CUSTOMER_TYPE_TEMP'
}])

const getmaplist = (keys) => {
    console.log(keys);
    if(keys){
        const list = []
    const a = getKeyMaplist()
    keys.forEach(v => {
        a.forEach(h => {
            if (Object.keys(h) == v) {
                console.log(Object.values(h));
                list.push(Object.values(h)[0])
            }
        })
    })
    return list.join(',')
    }else{
        return '--'
    }
}
const handleSizeChange = (val) => {
    SearchReports()
}
const handleCurrentChange = (val) => {
    SearchReports()
}
//获取map值
const getmep = () => {
    const data = getSymptomKeyMap()
    const key = data.risk_disease_key_map
    const dirkey = data.dirty_dialectics_key_map
    const phykey = data.physical_therapy_key_map
    delete dirkey.Z0022
    dircities.value = dirkey
    const Z0004 = data.dirty_dialectics_key_map.Z0004
    const Z0008 = data.dirty_dialectics_key_map.Z0008
    const Z0011 = data.dirty_dialectics_key_map.Z0011
    delete dirkey.Z0004
    delete dirkey.Z0008
    delete dirkey.Z0011
    dircities.value.Z0004 = Z0004
    dircities.value.Z0008 = Z0008
    dircities.value.Z0011 = Z0011
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
    // let date = new Date(value);
    // let month = (date.getMonth() + 1).toString().padStart(2, '0');
    // let year = date.getFullYear();
    // let day = new Date(year, month, 0);
    // let endDate = new Date(new Date(day.toLocaleDateString()).getTime() + 24 * 60 * 60 * 1000 - 1);
    Date2.value = new Date(new Date(new Date(Date2.value).toLocaleDateString()).getTime() + 24 * 60 * 60 * 1000 - 1);
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

const dirhandleCheckAllChange = (val) => {
    const aa = Object.keys(dircities.value)
    dircheckedCities.value = val ? aa : [];
    dirisIndeterminate.value = false;
    console.log(dircheckedCities.value);
    let text = []

    for (let key in dircheckedCities.value) {
        text.push(dircities.value[dircheckedCities.value[key]])
    }
    lists.value = text.join('、')
    if (lists.value == '') {
        lists.value = '请选择'
    }
    console.log(lists.value);
}
const dirhandleCheckedCitiesChange = (val) => {
    const aa = Object.keys(dircities.value)
    let checkedCount = val.length;
    dircheckAll.value = checkedCount === aa.length;
    dirisIndeterminate.value = checkedCount > 0 && checkedCount < aa.length;
    console.log(val);
    let text = []
    val.forEach(element => {
        console.log(element);
        text.push(dircities.value[element])
    })
    console.log(text);
    lists.value = text.join('、')
    if (lists.value == '') {
        lists.value = '请选择'
    }
    console.log(lists.value);
}

const SearchReports = async () => {
    const res = {
        tenantName: tenantName.value,
        customerName: customerName.value,
        customerPhone: customerPhone.value,
        startTime: Date1.value,
        endTime: Date2.value,
        dirtyDialectics: dircheckedCities.value,
        organizationId: getTenantId(),
        customerType:yaopin.value,
        pagination: {
            offset: (offset.value - 1) * pageSize4.value,
            size: pageSize4.value
        },
    }
    const data = await SearchReportsRequestRequest(res)
    ruquest.value = data?.data
    total.value = data?.data.total_count
}

const clickreport = (report_id,tenant_id) => {
    dialogFormVisible.value = true
    reporturl.value = ''
    reporturl.value = 'https://res.jinmuhealth.com' + import.meta.env.VITE_APP_REPORT_URL + 'index.html#/private/report/' + report_id + '?' + 'tenant_id=' + tenant_id + "&v=" + Math.random();    
    console.log(reporturl.value);
    
}
const topname = () => {
    settopname([{ name: '概况', url: '/' }, { name: '测量数据详情', url: '' }])
}
SearchReports()
getmep()
topname()
</script>
<style lang="scss">
.custdata {
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

        .selectinp {
            width: 490px;
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

}
</style>