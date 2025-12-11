<template>
    <div class="custdata">
        <div class="custtop">
            <div class="ownerright">
                <el-input v-model="tenantName" placeholder="商户名称" class="ownerinp"
                    @blur="tenantName = $event.target.value.trim()" :maxlength="20"></el-input>
                <el-input v-model="customerName" placeholder="客户姓名" class="ownerinp"
                    @blur="customerName = $event.target.value.trim()" :maxlength="20"></el-input>
                <el-input v-model="customerPhone" placeholder="客户手机号" class="ownerinp"
                    @blur="customerPhone = $event.target.value.trim()" :maxlength="20"></el-input>
                <el-button type="primary" class="owbtm" @click="SearchCustomerLastStatus">搜索</el-button>
            </div>
        </div>
        <div class="custbody">
            <el-table ref="liecheck" :data="customerlast?.customers" style="width: 100%">
                <el-table-column prop="tenant_name" label="商户名称" />
                <el-table-column prop="customer_name" label="客户名称" />
                <el-table-column prop="customer_phone" label="客户手机号" />
                <el-table-column prop="overdue_count" label="距离上次检查天数">
                    <template #default="scope">
                        <span v-if="scope.row.overdue_count">
                            {{ scope.row.overdue_count }}
                        </span>
                        <span v-else>
                            0
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
    </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import { settopname } from '../../utils/sethometop'
import { SearchCustomerLastStatusRequest } from '@/api/api'
import { getTenantId } from '@/utils/storage'

//商户名称
const tenantName = ref('')
//商户名称
const customerName = ref('')
//商户名称
const customerPhone = ref('')
//默认每页几条数据
const pageSize4 = ref(10)
//默认总数据
const total = ref(0)
//使用中商品
const shiyongz = ref(0)
//分页偏移量
const offset = ref(1)
//数据列表
const customerlast = ref()

const handleSizeChange = (val) => {
    SearchCustomerLastStatus()
}
const handleCurrentChange = (val) => {
    SearchCustomerLastStatus()
}
const SearchCustomerLastStatus = async () => {
    const rs = {
        organizationId: getTenantId(),
        tenantName: tenantName.value,
        customerName: customerName.value,
        customerPhone: customerPhone.value,
        overdueMask:'30',
        pagination: {
            offset: (offset.value - 1) * pageSize4.value,
            size: pageSize4.value
        },
    }
    const data = await SearchCustomerLastStatusRequest(rs)
    customerlast.value = data.data
    total.value = data?.data.total_count==undefined?0:data?.data.total_count

}

const topname = () => {
    settopname([{ name: '概况', url: '/' }, { name: '客户数据详情', url: '' }])
}
SearchCustomerLastStatus()
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