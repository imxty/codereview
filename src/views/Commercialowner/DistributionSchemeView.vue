<template>
    <div class="distrbution">
        <div class="measureds prodatadiv">
            <div class="prodatatop">
                商户列表
            </div>
            <div class="prodatari">
                如果商户已有方案，将被新方案替代
            </div>
            <div class="measurecnt">
                <div v-if="sftru == false" class="meadata">
                    <div v-for="(item, index) in tenantlist" :key="item" class="distrdivs">
                        <div v-if="index < 20">
                            {{ item.name }}
                        </div>
                        <div style="margin-left: 25px;margin-top: 2px;" @click="deletedata(index)">
                            <img src="../../assets/icons/x.svg" alt="">
                        </div>
                    </div>
                </div>

                <div v-if="sftru == false" @click="sftru = true" class="meatext">
                    显示全部已选项
                </div>
                <div v-if="sftru == true" @click="sftru = false" class="meatext">
                    收起全部已选项
                </div>
            </div>
            <div class="measuredata" v-if="sftru == true">
                <el-table ref="liecheck" :data="tenantlist" style="width: 100%">
                    <el-table-column prop="name" label="商户名称/商户ID">
                        <template #default="scope">
                            <div>
                                <div class="lienametxt">
                                    {{ scope.row.name }}
                                </div>
                                <div>
                                    {{ scope.row.tenant_id }}
                                </div>
                            </div>
                        </template>
                    </el-table-column>

                    <el-table-column prop="phone" label="方案">
                        <template #default="scope">
                            <div>
                                <div>
                                    {{ scope.row.contact_name }}
                                </div>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" :filter-method="filterTag" filter-placement="bottom-end" width="300">
                        <template #default="scope">
                            <div>
                                <el-button type="primary" plain @click.stop="deletedata(scope.$index)">
                                    删除</el-button>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </div>
        <div class="measureds prodatadiv" v-if="sftru == false">
            <div class="prodatatop">
                当前可使用方案
            </div>
            <div class="prodatadata">
                <div class="datatop">
                    <div class="topdiv">
                        <div>方案名称</div>
                        <div style="flex: 0.4;">方案上架时间</div>
                    </div>
                    <div class="topdiv" style="margin-left: -20px;">
                        <div>方案名称</div>
                        <div style="flex: 0.4;">方案上架时间</div>
                    </div>
                </div>
                <el-radio-group v-model="treatmentId">
                    <div class="prodatadatagroup">
                        <div v-for="(item, index) in selectdata" :key="index" class="groupdiv">
                            <el-radio :label="item.treatment_id" size="large">{{ item.treatment_name }} <span>{{
                                fermitTime(item.created_time) }}</span> </el-radio>
                        </div>
                    </div>
                </el-radio-group>
            </div>
        </div>
        <div class="prodatabom">
            <el-button type="primary" class="owbtm" @click.stop=" SubmitTreatmentToTenant">分配方案</el-button>
        </div>
    </div>
</template>
<script setup>
import { ref, reactive, onUnmounted } from 'vue'
import { getTenantId } from '@/utils/storage'
import { ElMessage } from 'element-plus';
import { SubmitTreatmentToTenantRequest, ListTreatmentsRequest } from '@/api/api'
import { useRouter } from 'vue-router';
import { settopname } from '../../utils/sethometop'
import { fermitTime } from '@/utils/fermitTime'


const tenantlist = ref(JSON.parse(history.state.keyword))
console.log(tenantlist.value);

const tenantIds = ref([])

const sftru = ref(false)

const radio1 = ref('1')

const treatmentId = ref('')
console.log(treatmentId.value);

const selectdata = ref([])
const SubmitTreatmentToTenant = async () => {
    let tid = []
    tenantlist.value.forEach(element => {
        tid.push(element.tenant_id)
    });
    const res = {
        organizationId: getTenantId(),
        tenantIds: tid,
        treatmentId: treatmentId.value
    }
    const data = await SubmitTreatmentToTenantRequest(res)
    if (data.status === 200) {
        ElMessage({
            message: '分配方案成功',
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
const deletedata = (index) => {
    tenantlist.value.splice(index, 1)
    if (tenantlist.value.length == 0) {
        window.history.go(-1)
    }
}

const ListTreatments = async () => {
    const data = await ListTreatmentsRequest({
        organizationId: getTenantId(),
        isPublish: true
    })
    console.log(data);
    selectdata.value = data.data?.treatments
}
ListTreatments()
const topname = () => {
    settopname([{ name: '商户管理', url: '/merchantManagement' }, { name: '分配方案', url: '' }])
}
topname()
</script>

<style lang="scss">
.distrbution {

    .prodatabom {
        display: flex;
        justify-content: flex-end;
    }

    .owbtm {
        background-color: rgba(21, 117, 238, 1);
        height: 40px;
        padding: 0 25px;
    }

    .meatext {
        font-size: 14px;
        color: #1575EE;
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
        margin-bottom: 30px;

        .prodatatop {
            height: 70px;
            line-height: 70px;
            font-weight: 700;
            font-size: 18px;
            padding-left: 40px;
            border-bottom: 1px solid rgba(238, 238, 238, 1);
        }

        .el-radio-group {
            width: 100%;
        }

        .prodatadata {

            .datatop {
                width: 90%;
                display: flex;
                justify-content: space-between;
                margin: auto;

                .topdiv {
                    flex: 0.5;
                    display: flex;
                    justify-content: space-evenly;
                }
            }

            .prodatadatagroup {
                width: 100%;
                display: flex;
                flex-wrap: wrap;

                .groupdiv {
                    width: 45%;
                    margin-left:-9px;
                    .el-radio {
                        width: 100%;
                        justify-content: center;

                        .el-radio__label {
                            flex: 0.7;
                            display: flex;
                            justify-content: space-around;
                        }
                    }
                }
            }
        }

        .prodatari {
            font-weight: 650;
            font-style: normal;
            color: #1575EE;
            position: absolute;
            right: 10px;
            top: 25px;
            font-size: 14px;
            cursor: pointer;
        }

        .measurecnt {
            display: flex;
            flex-wrap: wrap;
            align-items: baseline;
            flex: 1;
            padding: 20px 40px;

            .meadata {
                display: flex;
                flex-wrap: wrap;
                align-items: baseline;
            }

            .distrdivs {
                height: 37px;
                background-color: #ddebfb;
                color: #2F84EF;
                padding: 0px 6px;
                display: flex;
                align-items: center;
                margin-right: 10px;
                margin-bottom: 25px;
            }
        }
    }
}
</style>