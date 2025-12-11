<template>
    <div class="membertop">
        <el-button type="primary" class="tjsp" @click.stop="tz('addProduct')">添加商品</el-button>
        <div class="sptj">
            方案上限：{{ (memoption?.length == undefined ? 0 : memoption.length) + (memoptions?.length == undefined ? 0 :
                memoptions.length) }}/20
        </div>
    </div>
    <div class="yfb" v-if="memoptions?.length > 0">
        <div class="fbdiv">
            已发布方案
        </div>
    </div>
    <div class="membercnt">
        <div class="memdata">
            <div v-for="(key, i) in memoptions" :key="key" class="memdiv" @click.stop.stop="tsxq(key)">
                <div class="onneame">
                    {{ fanganxb[i] }}
                </div>
                <div class="optdiv">
                    <div class="opttag" v-if="key.publish == true">使用中</div>
                    <div class="optdivtop">
                        <div class="optname">{{ key.treatment_name }}</div>
                        <div class="optdate">{{ key.created_time }}</div>
                        <div class="optstat"
                            :class="[{ txtred: key.treatment_status == 'TREATMENT_STATUS_UNAPPROVED' }, { txtbl: key.treatment_status == 'TREATMENT_STATUS_UNDER_REVIEW' }]">
                            <span v-if="key.treatment_status == 'TREATMENT_STATUS_DRAFT'">草稿</span>
                            <span v-if="key.treatment_status == 'TREATMENT_STATUS_UNDER_REVIEW'">审核中</span>
                            <span v-if="key.treatment_status == 'TREATMENT_STATUS_APPROVED'">已过审</span>
                            <span v-if="key.treatment_status == 'TREATMENT_STATUS_UNAPPROVED'">未过审</span>
                        </div>

                        <div @click.stop="">
                            <el-popover placement="bottom" :width="200" trigger="click" :content="key.review_comment">
                                <template #reference>
                                    <div v-if="key.treatment_status == 'TREATMENT_STATUS_UNAPPROVED'"
                                        style="margin-left: 15px;">审核意见 <img
                                            src="../../assets/icons/button_pulldown.svg" alt="" srcset=""
                                            style="width: 12px;"></div>
                                </template>
                            </el-popover>
                        </div>
                        <div v-if="key.is_published == undefined">未发布</div>
                        <div v-if="key.is_published">已发布</div>
                    </div>
                    <div class="optdivbom">
                        <div v-if="key.treatment_status == 'TREATMENT_STATUS_DRAFT'" class="caog">
                            <el-button type="success" plain class="opbut" @click.stop=" tstc(key)">提交审核</el-button>
                            <el-button type="primary" plain class="opbut" @click.stop="routerps(key)">编辑</el-button>
                            <el-button type="danger" plain class="opbut"
                                @click.stop="delettreatment(key)">删除</el-button>
                        </div>
                        <div v-if="key.treatment_status == 'TREATMENT_STATUS_UNDER_REVIEW'">
                            <el-button type="primary" plain class="opbut"
                                @click.stop="CancelTreatmentReview(key.treatment_id)">取消审核</el-button>
                        </div>
                        <div v-if="key.treatment_status == 'TREATMENT_STATUS_UNAPPROVED'">
                            <el-button type="primary" plain class="opbut" @click.stop="routerps(key)">编辑</el-button>
                        </div>
                        <div v-if="key.treatment_status == 'TREATMENT_STATUS_APPROVED' && key.is_published == undefined"
                            class="caog">
                            <el-button type="primary" plain @click.stop="PublishTreatment(key.treatment_id)"
                                class="opbut">使用</el-button>
                            <el-button type="primary" plain class="opbut" @click.stop="routerps(key)">编辑</el-button>
                            <el-button type="danger" plain @click.stop="delettreatment(key)"
                                class="opbut">删除</el-button>
                        </div>
                        <el-button type="warning" plain class="opbut" @click.stop="routerps(key, false)">复制</el-button>
                        <br />
                        <el-button type="primary" plain class="opbut" @click.stop="RemoveTreatment(key.treatment_id)"
                            v-if="key.is_published == true">下架</el-button>
                    </div>
                    <div v-if="key.is_published" class="optdivtips"><img src="../../assets/image/u76.png" alt=""
                            width="20">如需修改当前方案，请先解除商户与方案绑定关系，再进行下架</div>
                </div>
            </div>
        </div>
    </div>
    <div class="yfb">
        <div class="fbdiv">
            未发布方案
        </div>
    </div>
    <div class="membercnt">
        <div class="memdata">
            <div class="memdiv" style="padding-top: 20px;" @click.stop="tz('treatmentEdit')"
                v-if="(memoption?.length == undefined ? 0 : memoption.length) + (memoptions?.length == undefined ? 0 : memoptions.length) < 20">
                <div class="optdiv">
                    <img src="../../assets/image/icon_add.png" alt="" style="width: 36px;">
                    <div style="color:#C1C7D0;margin-top: 20px;">点击创建</div>
                </div>
            </div>
            <div v-for="(key, i) in memoption" :key="key" class="memdiv" @click.stop.stop="tsxq(key)">
                <div class="onneame">
                    {{ fanganxb[i] }}
                </div>
                <div class="optdiv">
                    <div class="opttag" v-if="key.publish == true">使用中</div>
                    <div class="optdivtop">
                        <div class="optname">{{ key.treatment_name }}</div>
                        <div class="optdate">{{ key.created_time }}</div>
                        <div class="optstat"
                            :class="[{ txtred: key.treatment_status == 'TREATMENT_STATUS_UNAPPROVED' }, { txtbl: key.treatment_status == 'TREATMENT_STATUS_UNDER_REVIEW' }]">
                            <span v-if="key.treatment_status == 'TREATMENT_STATUS_DRAFT'">草稿</span>
                            <span v-if="key.treatment_status == 'TREATMENT_STATUS_UNDER_REVIEW'">审核中</span>
                            <span v-if="key.treatment_status == 'TREATMENT_STATUS_APPROVED'">已过审</span>
                            <span v-if="key.treatment_status == 'TREATMENT_STATUS_UNAPPROVED'">未过审</span>
                        </div>

                        <div @click.stop="">
                            <el-popover placement="bottom" :width="200" trigger="click" :content="key.review_comment">
                                <template #reference>
                                    <div v-if="key.treatment_status == 'TREATMENT_STATUS_UNAPPROVED'"
                                        style="margin-left: 15px;">审核意见 <img
                                            src="../../assets/icons/button_pulldown.svg" alt="" srcset=""
                                            style="width: 12px;"></div>
                                </template>
                            </el-popover>
                        </div>
                        <div v-if="key.is_published == undefined">未发布</div>
                        <div v-if="key.is_published">已发布</div>
                    </div>
                    <div class="optdivbom">
                        <div v-if="key.treatment_status == 'TREATMENT_STATUS_DRAFT'" class="caog">
                            <el-button type="success" plain class="opbut" @click.stop=" tstc(key)">提交审核</el-button>
                            <el-button type="primary" plain class="opbut" @click.stop="routerps(key)">编辑</el-button>
                            <el-button type="danger" plain class="opbut"
                                @click.stop="delettreatment(key)">删除</el-button>
                        </div>
                        <div v-if="key.treatment_status == 'TREATMENT_STATUS_UNDER_REVIEW'">
                            <el-button type="primary" plain class="opbut"
                                @click.stop="CancelTreatmentReview(key.treatment_id)">取消审核</el-button>
                        </div>
                        <div v-if="key.treatment_status == 'TREATMENT_STATUS_UNAPPROVED'">
                            <el-button type="primary" plain class="opbut" @click.stop="routerps(key)">编辑</el-button>
                        </div>
                        <div v-if="key.treatment_status == 'TREATMENT_STATUS_APPROVED' && key.is_published == undefined"
                            class="caog">
                            <el-button type="primary" plain @click.stop="PublishTreatment(key.treatment_id)"
                                class="opbut">使用</el-button>
                            <el-button type="primary" plain class="opbut" @click.stop="routerps(key)">编辑</el-button>
                            <el-button type="danger" plain @click.stop="delettreatment(key)"
                                class="opbut">删除</el-button>
                        </div>
                        <el-button type="warning" plain class="opbut" @click.stop="routerps(key, false)">复制</el-button>
                    </div>
                    <div v-if="key.is_published" class="optdivtips"><img src="../../assets/image/u76.png" alt=""
                            width="20">如需修改当前方案，请添加其他方案并进行使用后再进行编辑</div>
                </div>
            </div>
        </div>
    </div>
    <el-dialog v-model="dialogVisible" title="确定提交审核" width="30%">
        <div class="submitcnt">
            <div class="subdata">
                <div class="subdatatop">
                    以下为此次提交方案，请确认
                </div>
                <div>
                    <el-table :data="shangpintishen" height="200" border style="width: 100%">
                        <el-table-column prop="keyname" label="症候" align="center">
                        </el-table-column>
                        <el-table-column prop="product" label="药品名称" align="center">
                    <template #default="scope">
                        {{ scope.row.productfood?.product_name || '--' }}
                    </template>
                </el-table-column>
                <el-table-column prop="product" label="理疗名称" align="center">
                    <template #default="scope">
                        {{ scope.row.productll?.product_name || '--'}}
                    </template>
                </el-table-column>
                    </el-table>
                </div>
            </div>
        </div>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click.stop="dialogVisible = false" class="elbtn">取消</el-button>
                <el-button type="primary" @click.stop="tishen" class="elbtn elqd">
                    确认提交审核
                </el-button>
            </span>
        </template>
    </el-dialog>
    <el-dialog v-model="limitation" title="提审受限" width="30%">
        <div class="submitcnt">
            <div class="limitation">
                对不起，每30天仅可发起一轮审核，当前时间段您已没有审核机会，请稍后再试
            </div>
        </div>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click.stop="limitation = false">确定</el-button>
            </span>
        </template>
    </el-dialog>
    <el-dialog v-model="detailfn" title="删除推荐方案" width="30%">
        <div class="submitcnt">
            <div class="limitation">
                确定删除方案 {{ detaildata.treatment_name }} 吗？
            </div>
        </div>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click.stop="detailfn = false" class="elbtn">取消</el-button>

                <el-button @click.stop="deletereatm" class="elbtn elqd">确定</el-button>
            </span>
        </template>
    </el-dialog>
</template>
<script setup>
import { reactive, ref } from "vue";
import { useRouter } from 'vue-router';
import { getTenantId, getKeyMaplist } from '@/utils/storage'
import {
    ListTreatmentsRequest,
    GetTreatmentReviewStatusRequest,
    SubmitTreatmentToReviewRequest,
    GetTreatmentRequest,
    DeleteTreatmentRequest,
    CancelTreatmentReviewRequest,
    PublishTreatmentRequest,
    GetOrganizationProductRecommendationStatusRequest,
    ModifyProductRecommendationRequest,
    RemoveTreatmentRequest
} from '@/api/api'
import { fermitTime } from '@/utils/fermitTime'
import { ElMessage } from 'element-plus';
import { settopname } from '@/utils/sethometop'
import { elmessage } from '@/utils/popup'

const $router = useRouter();
const mapkeyy = ref(getKeyMaplist())
//提交审核弹窗开关
const dialogVisible = ref(false)
const shoulicishu = ref(15)
const fanganxb = ref(['方案一', '方案二', '方案三', '方案四', '方案五', '方案六', '方案七', '方案八', '方案九', '方案十', '方案十一', '方案十二', '方案十三', '方案十四', '方案十五', '方案十六', '方案十七', '方案十八', '方案十九', '方案二十'])
//提审受限开关
const limitation = ref(false)
//商品推荐方案列表
const memoption = ref([])
//商品推荐方案列表
const memoptions = ref([])
//提审方案数据
const tijiaodata = ref(null)
//提审商品数据
const shangpintishen = ref([])
//删除方案弹窗
const detailfn = ref(false)
//删除方案数据
const detaildata = ref(null)
const tz = (to) => {
    $router.push({
        name: to,
        state: {
            keyword: JSON.stringify({ treatment: {} }),
            keycz: JSON.stringify(false)
        }
    })
}
const getlist = () => {
    const rs = {
        organizationId: getTenantId(),
        isPublish: false
    }
    ListTreatmentsRequest(rs)
        .then(res => {
            console.log(res);
            memoption.value = res.data.treatments == undefined ? [] : res.data.treatments
            console.log(memoption.value);

            memoption.value.forEach(v => {
                v.created_time = fermitTime(v.created_time)
            })
        }).catch(error => {
            console.log(error);
        })

    ListTreatmentsRequest({
        organizationId: getTenantId(),
        isPublish: true
    })
        .then(res => {
            console.log(res);
            memoptions.value = res.data.treatments == undefined ? [] : res.data.treatments
            memoptions.value.forEach(v => {
                v.created_time = fermitTime(v.created_time)
            })
        }).catch(error => {
            console.log(error);
        })

}
const tstc = (re) => {
    const rs = {
        organizationId: getTenantId(),
        treatmentId: re.treatment_id
    }
    GetTreatmentRequest(rs)
        .then(
            res => {
                shangpintishen.value = []
                tijiaodata.value = res.data.treatment
                dialogVisible.value = true
                tijiaodata.value.treatment_items_risky_disease?.forEach(v => {
                    mapkeyy.value.forEach(h => {
                        if (Object.keys(h) == v.symptom) {
                            let lllist = null
                            let foodlist = null
                            v.products.forEach(eva => {
                                if (eva.product_type == 'PRODUCT_TYPE_SERVICE') {
                                    lllist = eva
                                    console.log(lllist);
                                } else {
                                    foodlist = eva
                                    console.log(foodlist);
                                }
                            })
                            shangpintishen.value.push({
                                keyname: Object.values(h)[0],
                                productfood: foodlist,
                                productll: lllist,
                            })
                        }
                    })
                });
                tijiaodata.value.treatment_items_dirty_dialectic?.forEach(v => {
                    mapkeyy.value.forEach(h => {
                        if (Object.keys(h) == v.symptom) {
                            let lllist = null
                            let foodlist = null
                            v.products.forEach(eva => {
                                if (eva.product_type == 'PRODUCT_TYPE_SERVICE') {
                                    lllist = eva
                                    console.log(lllist);
                                } else {
                                    foodlist = eva
                                    console.log(foodlist);
                                }
                            })
                            shangpintishen.value.push({
                                keyname: Object.values(h)[0],
                                productfood: foodlist,
                                productll: lllist,
                            })
                        }
                    })
                });
                tijiaodata.value.treatment_items_physical_therapy?.forEach(v => {
                    mapkeyy.value.forEach(h => {
                        if (Object.keys(h) == v.symptom) {
                            let lllist = null
                            let foodlist = null
                            v.products.forEach(eva => {
                                if (eva.product_type == 'PRODUCT_TYPE_SERVICE') {
                                    lllist = eva
                                    console.log(lllist);
                                } else {
                                    foodlist = eva
                                    console.log(foodlist);
                                }
                            })
                            shangpintishen.value.push({
                                keyname: Object.values(h)[0],
                                productfood: foodlist,
                                productll: lllist,
                            })
                        }
                    })
                });
                tijiaodata.value.treatment_items_physical_dialectics?.forEach(v => {
                    mapkeyy.value.forEach(h => {
                        if (Object.keys(h) == v.symptom) {
                            let lllist = null
                            let foodlist = null
                            v.products.forEach(eva => {
                                if (eva.product_type == 'PRODUCT_TYPE_SERVICE') {
                                    lllist = eva
                                    console.log(lllist);
                                } else {
                                    foodlist = eva
                                    console.log(foodlist);
                                }
                            })
                            shangpintishen.value.push({
                                keyname: Object.values(h)[0],
                                productfood: foodlist,
                                productll: lllist,
                            })
                        }
                    })
                });
                console.log(shangpintishen.value);
            })
        .catch(error => {
            console.log(error);
        })
}
const tsxq = (re) => {
    const rs = {
        organizationId: getTenantId(),
        treatmentId: re.treatment_id
    }
    GetTreatmentRequest(rs)
        .then(
            res => {
                shangpintishen.value = []
                tijiaodata.value = res.data.treatment
                console.log(res.data.treatment);
                tijiaodata.value.treatment_items_risky_disease?.forEach(v => {
                    mapkeyy.value.forEach(h => {
                        if (Object.keys(h) == v.symptom) {
                            let lllist = null
                            let foodlist = null
                            v.products.forEach(eva => {
                                if (eva.product_type == 'PRODUCT_TYPE_SERVICE') {
                                    lllist = eva
                                    console.log(lllist);
                                } else {
                                    foodlist = eva
                                    console.log(foodlist);
                                }
                            })
                            shangpintishen.value.push({
                                keyname: Object.values(h)[0],
                                productfood: foodlist,
                                productll: lllist,
                            })
                        }
                    })
                });
                tijiaodata.value.treatment_items_dirty_dialectic?.forEach(v => {
                    mapkeyy.value.forEach(h => {
                        if (Object.keys(h) == v.symptom) {
                            let lllist = null
                            let foodlist = null
                            v.products.forEach(eva => {
                                if (eva.product_type == 'PRODUCT_TYPE_SERVICE') {
                                    lllist = eva
                                    console.log(lllist);
                                } else {
                                    foodlist = eva
                                    console.log(foodlist);
                                }
                            })
                            shangpintishen.value.push({
                                keyname: Object.values(h)[0],
                                productfood: foodlist,
                                productll: lllist,
                            })
                        }
                    })
                });
                tijiaodata.value.treatment_items_physical_therapy?.forEach(v => {
                    mapkeyy.value.forEach(h => {
                        if (Object.keys(h) == v.symptom) {
                            let lllist = null
                            let foodlist = null
                            v.products.forEach(eva => {
                                if (eva.product_type == 'PRODUCT_TYPE_SERVICE') {
                                    lllist = eva
                                    console.log(lllist);
                                } else {
                                    foodlist = eva
                                    console.log(foodlist);
                                }
                            })
                            shangpintishen.value.push({
                                keyname: Object.values(h)[0],
                                productfood: foodlist,
                                productll: lllist,
                            })
                        }
                    })
                });
                tijiaodata.value.treatment_items_physical_dialectics?.forEach(v => {
                    mapkeyy.value.forEach(h => {
                        if (Object.keys(h) == v.symptom) {
                            let lllist = null
                            let foodlist = null
                            v.products.forEach(eva => {
                                if (eva.product_type == 'PRODUCT_TYPE_SERVICE') {
                                    lllist = eva
                                    console.log(lllist);
                                } else {
                                    foodlist = eva
                                    console.log(foodlist);
                                }
                            })
                            shangpintishen.value.push({
                                keyname: Object.values(h)[0],
                                productfood: foodlist,
                                productll: lllist,
                            })
                        }
                    })
                });
                $router.push({
                    name: 'detailsView',
                    state: {
                        keyword: JSON.stringify({ name: tijiaodata.value.treatment_name, time: fermitTime(tijiaodata.value.created_time), data: shangpintishen.value })
                    }
                })
                console.log(shangpintishen.value);
            })
        .catch(error => {
            console.log(error);
        })
}
const tishen = () => {
    const rs = {
        organizationId: getTenantId(),
        treatmentId: tijiaodata.value.treatment_id
    }
    SubmitTreatmentToReviewRequest(rs)
        .then(res => {
            if (res.status == 200) {
                ElMessage({
                    message: '提审成功',
                    type: 'success',
                })
                getlist()
                dialogVisible.value = false
            } else {
                elmessage(res.data.detail)
            }
        })
}

const delettreatment = (data) => {
    detaildata.value = data
    detailfn.value = true
}

const deletereatm = () => {
    console.log(detaildata.value);
    const rs = {
        organizationId: getTenantId(),
        treatmentId: detaildata.value.treatment_id
    }
    DeleteTreatmentRequest(rs)
        .then(res => {
            if (res.status == 200) {
                ElMessage({
                    message: '删除成功',
                    type: 'success',
                })
                getlist()
                detailfn.value = false
            } else {
                elmessage(res.data.detail)
            }
        })
}
//编辑跳转
const routerps = (re, cz = true) => {
    console.log(memoption.value);

    let memoptiona = memoption.value.length == undefined ? 0 : memoption.value.length
    let memoptionsa = memoptions.value.length == undefined ? 0 : memoptions.value.length
    if (cz == false) {
        if (memoptiona + memoptionsa >= 20) {
            ElMessage({
                showClose: true,
                message: '您的方案已达上限，无法进行复制，请先删除部分暂未使用的方案',
                type: 'error',
            })
            return
        }
    }
    const rs = {
        organizationId: getTenantId(),
        treatmentId: re.treatment_id
    }
    GetTreatmentRequest(rs)
        .then(res => {
            $router.push({
                name: 'treatmentEdit',
                state: {
                    keyword: JSON.stringify(res.data),
                    keycz: JSON.stringify(cz),
                }
            })
        })
        .catch(error => {
            console.log(error);
        })
}
const CancelTreatmentReview = (id) => {
    const rs = {
        organizationId: getTenantId(),
        treatmentId: id
    }
    CancelTreatmentReviewRequest(rs)
        .then(res => {
            if (res.status == 200) {
                ElMessage({
                    message: '取消审核成功',
                    type: 'success',
                })
                getlist()
            } else {
                ElMessage({
                    showClose: true,
                    message: res.data.detail,
                    type: 'error',
                })
            }
        })
}
const PublishTreatment = (id) => {
    const rs = {
        organizationId: getTenantId(),
        treatmentId: id
    }
    PublishTreatmentRequest(rs)
        .then(res => {
            if (res.status == 200) {
                ElMessage({
                    message: '发布推荐方案成功成功',
                    type: 'success',
                })
                getlist()
            } else {
                ElMessage({
                    showClose: true,
                    message: res.data.detail,
                    type: 'error',
                })
            }
        })
}
const RemoveTreatment = (id) => {
    const rs = {
        organizationId: getTenantId(),
        treatmentId: id
    }
    RemoveTreatmentRequest(rs)
        .then(res => {
            if (res.status == 200) {
                ElMessage({
                    message: '下架推荐方案成功',
                    type: 'success',
                })
                getlist()
            } else {
                ElMessage({
                    showClose: true,
                    message: res.data.detail,
                    type: 'error',
                })
            }
        })
}

const ModifyProductRecommendation = () => {
    ModifyProductRecommendationRequest({ organizationId: getTenantId() })
}
const topname = () => {
    settopname([{ name: '商品推荐方案', url: '' }])
}
topname()
getlist()
</script>
<style lang="scss">
.membertop {
    display: flex;
    justify-content: space-between;

    .sptj {
        display: flex;
        align-items: center;
        font-size: 14px;
        color: #666666;
    }

    .tjsp {
        width: 118px;
        height: 40px;
        background-color: #1575ee;
        font-size: 14px;
        font-weight: 400;
        border: none;
    }
}

.yfb {
    .fbdiv {
        font-size: 15px;
        color: #1575EE;
        border-bottom: #1575ee 2px solid;
        width: 15%;
        margin: 30px 0;
        padding-bottom: 10px;
    }
}

.membercnt {
    .memdata {
        display: flex;
        flex-wrap: wrap;

        .memdiv {
            width: 18%;
            font-size: 14px;
            font-weight: 400;
            margin-right: 1.3vw;
            cursor: pointer;

            .onneame {
                color: rgb(193, 199, 208);
            }

            .optdiv {
                height: 400px;
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: center;
                background-color: #fff;
                border: 1px solid #ccc;
                border-radius: 10px;
                position: relative;
                overflow: hidden;
                margin-top: 10px;

                .opttag {
                    position: absolute;
                    top: -1px;
                    right: -1px;
                    width: 60px;
                    height: 24px;
                    text-align: center;
                    line-height: 24px;
                    font-size: 12px;
                    border-radius: 5px;
                    color: rgb(21, 117, 238);
                    background-color: #eaf4fb;
                }

                .optdivtop {
                    flex: 0.5;
                    display: flex;
                    flex-direction: column;
                    align-items: center;
                    justify-content: center;
                    box-sizing: content-box;
                    color: rgb(102, 102, 102);

                    .optname {
                        color: #333333;
                        font-weight: 700;
                        margin-top: 10px;
                        margin-bottom: 20px;
                    }
                }

                .optdivbom {
                    flex: 0.5;

                    .caog {
                        display: flex;
                        flex-direction: column;
                        align-items: center;
                    }

                    .opbut {
                        width: 85px;
                        height: 32px;
                        margin-top: 10px;
                        margin-left: 0;
                    }

                }
            }
        }

        .memdiv:first-child {
            margin-left: 0;
        }
    }
}


.txtred {
    color: #F56C6C;
    ;
}

.txtbl {
    color: #3CA1F8;
}

.subtips {
    font-size: 14px;
    color: rgb(245, 108, 108);
}

.subdatatop {
    margin: 30px 0;
}

.el-table--border::after,
.el-table__border-left-patch {
    height: 80%;
}

.el-dialog {
    border-radius: 10px;
}

.dialog-footer {
    .elbtn {
        width: 103px;
        height: 40px;
        font-weight: 400;
        font-style: normal;
        font-size: 14px;
    }

    .elqd {
        margin-left: 10px;
        background-color: #1575ee;
        color: #fff;
    }
}

.optdivtips {
    font-size: 12px;
    color: #C1C7D0;
    font-weight: 400;
    line-height: 16px;
    padding: 10px;
    text-align: center;

    img {
        width: 16px;
        position: relative;
        bottom: -3px;
    }
}
</style>