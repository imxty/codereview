<template>
    <el-button type="primary" @click="tz('addProduct')" class="tjsp">添加商品</el-button>
    <div class="usecomcnt">
        <div class="usecnttop">
            <div class="cnttoplf">
                <el-button plain @click="getlist('PRODUCT_STATUS_UNSET', true)"
                    :class="{ butxz: butxz == 'PRODUCT_STATUS_UNSET' }" class="spglbut">全部</el-button>

                <el-button plain @click="getlist('PRODUCT_STATUS_UNUSED',false)"
                    :class="{ butxz: butxz == 'PRODUCT_STATUS_UNUSED' }" class="spglbut">未使用</el-button>

                <el-button plain @click="getlist('PRODUCT_STATUS_TO_BE_USED',false)"
                    :class="{ butxz: butxz == 'PRODUCT_STATUS_TO_BE_USED' }" class="spglbut">待使用</el-button>

                <el-button plain @click="getlist('PRODUCT_STATUS_USING',false)"
                    :class="{ butxz: butxz == 'PRODUCT_STATUS_USING' }" class="spglbut">使用中</el-button>

                <el-button plain @click="getlist('PRODUCT_STATUS_DELETED',false)"
                    :class="{ butxz: butxz == 'PRODUCT_STATUS_DELETED' }" class="spglbut">已删除</el-button>
            </div>
        </div>
        <div class="usecntdata">
            <el-table :data="listproduct" style="width: 100%">
                <el-table-column label="商品名">
                    <template #default="scope">
                        <div style="display: flex; align-items: center;cursor: pointer;" @click="routerps(scope.row)">
                            <div class="commdname">
                                <div class="nameimg">
                                    <img :src="scope.row.image_url" alt="" style="width: 100%; height: 100%;">
                                </div>
                                <div class="nametxt">
                                    <div>
                                        {{ scope.row.product_name }}
                                    </div>
                                    <div class="namedse">{{ scope.row.description }}</div>
                                </div>
                            </div>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="商品类型" width="240">
                    <template #default="scope">
                        <span
                            v-if="scope.row.product_type == 'PRODUCT_TYPE_CPD' || scope.row.product_type == 'PRODUCT_TYPE_HEALTH_CARE_PRODUCT' || scope.row.product_type == 'PRODUCT_TYPE_NUTRITION'">药品</span>
                        <span v-if="scope.row.product_type == 'PRODUCT_TYPE_SERVICE'">理疗服务</span>
                    </template>
                </el-table-column>
                <el-table-column label="状态" width="240">
                    <template #default="scope">
                        <span class="wei" :class="[{
                            wei: scope.row.status === '未使用',
                            shan: scope.row.status === '已删除'
                        }]" v-if="scope.row.product_status === 'PRODUCT_STATUS_UNUSED'">未使用</span>
                        <span class="shan" v-if="scope.row.product_status === 'PRODUCT_STATUS_DELETED'">已删除</span>
                        <span v-if="scope.row.product_status === 'PRODUCT_STATUS_USING'">使用中</span>
                        <span class="dai" v-if="scope.row.product_status === 'PRODUCT_STATUS_TO_BE_USED'">待使用</span>
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="260">
                    <template #default="scope">
                        <el-button type="primary" plain @click="routerfx(scope.row)">数据统计</el-button>
                        <el-button type="primary" plain v-if="scope.row.product_status === 'PRODUCT_STATUS_UNUSED'"
                            @click="routerpske(scope.row)">
                            编辑</el-button>
                        <el-button type="danger" plain v-if="scope.row.product_status === 'PRODUCT_STATUS_UNUSED'"
                            @click="deletshp(scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>

        </div>
        <div class="fenye">
            <el-pagination v-model:current-page="offset" v-model:page-size="pageSize4" :page-sizes="[10, 30, 50, 100]"
                :small="true" :disabled="disabled" :background="background"
                layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange"
                @current-change="handleCurrentChange" />
        </div>
        <el-dialog v-model="daletsp" title="删除商品" width="30%" :before-close="handleClose">
            <span>是否确定删除商品 &nbsp; {{ shanchusp.product_name }}?</span>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="daletsp = false" class="elbut">取消</el-button>
                    <el-button type="primary" @click="deletqr" class="elbut quern">
                        确认删除
                    </el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import {
    ListProductsRequest,
    DeleteProductRequest
} from '@/api/api'
import { getTenantId } from '@/utils/storage'
import { ElMessage } from 'element-plus';
import { useRouter } from 'vue-router';
import { settopname } from '../../utils/sethometop'

const $router = useRouter();
//商品数据
const listproduct = ref([])
const getid = ref(getTenantId())
//默认每页几条数据
const pageSize4 = ref(10)
//默认总数据
const total = ref(0)
//分页偏移量
const offset = ref(1)
const getAllll = ref(true)
//删除商品弹窗中的商品名称
const shanchusp = ref('')
const daletsp = ref(false)
const butxz = ref('PRODUCT_STATUS_UNSET')
const handleSizeChange = (val) => {
    getlist(butxz.value,getAllll.value)
}
const handleCurrentChange = (val) => {
    getlist(butxz.value,getAllll.value)
}
const getlist = async (date, getAll = false) => {
    getAllll.value = getAll
    butxz.value = date
    const rs = {
        organizationId: getid.value,
        pagination: {
            offset: (offset.value - 1) * pageSize4.value,
            size: pageSize4.value
        },
        getAll: getAll,
        productStatus: date
    }
    const data = await ListProductsRequest(rs)
    console.log(data);
    listproduct.value = data.data.products
    total.value = data.data.total_count
}
const tz = (to) => {
    $router.push(
        {
            name: to
        }
    )
}
const routerps = (rs) => {
    $router.push({
        name: 'prodetails',
        state: {
            keyword: JSON.stringify(rs)
        }
    })
}
const routerpske = (rs) => {
    $router.push({
        name: 'edlgoods',
        state: {
            keyword: JSON.stringify(rs)
        }
    })
}

const routerfx = (rs) => {
    $router.push({
        name: 'dataanalysis',
        state: {
            keyword: JSON.stringify(rs)
        }
    })
}

const deletshp = (rs) => {
    daletsp.value = true
    console.log(rs);
    shanchusp.value = rs
}
const deletqr = () => {
    const rs = {
        organizationId: getTenantId(),
        productId: shanchusp.value.product_id
    }
    DeleteProductRequest(rs).then(res => {
        console.log(res);
        ElMessage({
            message: '删除成功',
            type: 'success',
        })
        getlist('PRODUCT_STATUS_UNSET', true)
        daletsp.value = false
    })
        .catch(error => {
            console.log(error);
            ElMessage({
                showClose: true,
                message: error.detail,
                type: 'error',
            })
        })

}
const topname = () => {
    settopname([{ name: '商品管理', url: '' }])
}
topname()
getlist('PRODUCT_STATUS_UNSET', true)
</script>
<style lang="scss" scoped>
.tjsp {
    width: 118px;
    height: 40px;
    background-color: #1575ee;
    font-size: 14px;
    font-weight: 400;
    border: none;
}

.usecomcnt {
    background-color: #fff;
    padding: 20px 20px 20px 20px;
    margin-top: 30px;
    border: 1px solid #eeeeee;
    ;
    border-radius: 20px;

    .usecnttop {
        display: flex;
        justify-content: space-between;

        .cnttoprg {
            font-size: 14px;
            color: rgb(102, 102, 102);
        }
    }

    .el-dialog {
        border-radius: 10px;
    }

    .el-dialog__body {
        border-top: 1px solid #eee;
    }

    .elbut {
        width: 103px;
        height: 40px;
    }

    .quern {
        background-color: #1575ee;
    }

}

.commdname {
    display: flex;

    .nameimg {
        width: 106px;
        height: 100px;
    }

    .nametxt {
        font-size: 16px;
        color: #333333;
        padding-top: 10px;
        margin-left: 10px;

        .namedse {
            font-size: 14px;
            color: rgb(193, 199, 208);
            width: 180px;
            margin-top: 8px;
        }
    }
}

.wei {
    color: #3ca1f8
}

.shan {
    color: #c1c7d0;
}

.dai {
    color: #F1AA2E;
}

.usecomcnt {
    position: relative;
    min-height: 75vh;

    .fenye {
        margin-top: 30px;
        position: absolute;
        bottom: -55px;
        right: 30px;
    }
}

.spglbut {
    width: 68px;
    height: 40px;
    background-color: #fff;
}

.butxz {
    background-color: rgba(236, 246, 254, 1);
    color: #3CA1F8;
    border: 1px solid rgba(177, 217, 252, 1);
}

.usecntdata {
    margin-top: 20px;
}
</style>