<template>
    <div class="addgood">
        <div class="add-top">
            <div class="toptxt " :class="[{ addtop_lf: comtyp === 1 }]" @click="comtyp = 1">
                选择商品类型
            </div>
            <div class="toptxt" :class="[{ addtop_lf: comtyp === 2 }]">
                填写商品信息
            </div>
        </div>
        <div class="addcnt" v-if="comtyp == 1">
            <div class="addcntdiv" @click="shangpintian('药品')">
                <div class="addcntimg">
                    <img src="../../assets/image/icon_adddrug.png" alt="">
                </div>
                <div>药品</div>
            </div>
            <div class="addcntdiv" @click="shangpintian('理疗')">
                <div class="addcntimg">
                    <img src="../../assets/image/icon_addservice.png" alt="">
                </div>
                <div>理疗服务方案</div>
            </div>
        </div>
        <div class="adddata" v-if="comtyp == 2">
            <div class="meddata" v-if="meddata == '药品'">
                <el-form ref="ruleFormRef" :model="CreateProduct.product" label-width="130px" :size="formSize"
                    :rules="rules">
                    <el-form-item label="商品分类：" required="true" style="text-align:left">
                        <el-select v-model="yaopin" placeholder="请选择" @change="aa">
                            <el-option v-for="item in options" :key="item.value" :label="item.label"
                                :value="item.value">
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="商品名称：" required="true" style="text-align:left" prop="productName">
                        <el-input placeholder="请输入商品名称" v-model="CreateProduct.product.productName" />
                    </el-form-item>
                    <el-form-item label="药品名称：" required="true" style="text-align:left"
                        v-if="yaopin == 'PRODUCT_TYPE_CPD'" prop="drugName">
                        <el-input placeholder="请输入药品名称" v-model="CreateProduct.product.drugName" />
                    </el-form-item>
                    <el-form-item label="药品类型：" required="true" style="text-align:left"
                        v-if="yaopin == 'PRODUCT_TYPE_CPD'">
                        <el-radio v-model="CreateProduct.product.isOtc" :label="true">处方药</el-radio>
                        <el-radio v-model="CreateProduct.product.isOtc" :label="false">非处方药</el-radio>
                    </el-form-item>
                    <el-form-item label="国药准字号：" required="true" style="text-align:left"
                        v-if="yaopin == 'PRODUCT_TYPE_CPD'" prop="approvedNumber">
                        <el-input placeholder="请输入国药准字号" v-model="CreateProduct.product.approvedNumber" />
                    </el-form-item>
                    <el-form-item label="健字号：" required="true" style="text-align:left"
                        v-if="yaopin == 'PRODUCT_TYPE_HEALTH_CARE_PRODUCT'" prop="approvedNumber">
                        <el-input placeholder="请输入健字号" v-model="CreateProduct.product.approvedNumber" />
                    </el-form-item>
                    <el-form-item label="食字号：" required="true" style="text-align:left"
                        v-if="yaopin == 'PRODUCT_TYPE_NUTRITION'" prop="approvedNumber">
                        <el-input placeholder="请输入食字号" v-model="CreateProduct.product.approvedNumber" />
                    </el-form-item>
                    <el-form-item label="药品准效期：" required="true" style="text-align:left"
                        v-if="yaopin == 'PRODUCT_TYPE_CPD'" prop="drugValidityPeriod">
                        <el-date-picker :editable="false" placeholder="请选择药品准效期" type="date"
                            v-model="CreateProduct.product.drugValidityPeriod"></el-date-picker>
                    </el-form-item>
                    <el-form-item label="商品图片：" required="true" style="text-align:left">
                        <div class="tableright">
                            <imageUploader @imageBase64="logoImageBase64" @onImageUploaded="logoImageUploaded">
                            </imageUploader>
                            <div class="imgtext">
                                <p>商家图标不允许涉及政治敏感与色情;</p>

                                <p>图片格式只支持：BMP、JPEG、JPG、GIF、PNG，大小不超过2M</p>
                            </div>
                        </div>
                    </el-form-item>
                    <el-form-item label="商品介绍：">
                        <el-input type="textarea" placeholder="请输入内容" v-model="CreateProduct.product.description"
                            maxlength="100" show-word-limit>
                        </el-input>
                    </el-form-item>
                    <el-form-item label="适应症候：" required="true">
                        <div class="checkcss">
                            <div class="checkquan checkquan1">
                                <el-checkbox :indeterminate="isIndeterminate" v-model="checkAll"
                                    @change="handleCheckAllChange">风险疾病</el-checkbox>
                            </div>
                            <div class="checkdan">
                                <el-checkbox-group v-model="checkedCities" @change="handleCheckedCitiesChange">
                                    <el-checkbox v-for="(city, value) in cities" :label="value" :key="city">{{ city
                                    }}</el-checkbox>
                                </el-checkbox-group>
                            </div>
                        </div>
                        <div class="checkcss">
                            <div class="checkquan">
                                <el-checkbox :indeterminate="dirisIndeterminate" v-model="dircheckAll"
                                    @change="dirhandleCheckAllChange">脏腑辨证</el-checkbox>
                            </div>
                            <div class="checkdan">
                                <el-checkbox-group v-model="dircheckedCities" @change="dirhandleCheckedCitiesChange">
                                    <el-checkbox v-for="(city, value) in dircities" :label="value" :key="city">{{ city
                                    }}</el-checkbox>
                                </el-checkbox-group>
                            </div>
                        </div>
                        <div class="checkcss tzcheck">
                            <div class="checkquan">
                                <el-checkbox :indeterminate="quedialectisIndeterminate" v-model="quedialectcheckAll"
                                    @change="dialecticsirhandleCheckAllChange">体质辨证</el-checkbox>
                            </div>
                            <div class="checkdan">
                                <el-checkbox-group v-model="quedialectcheckedCities"
                                    @change="dialecticsdirhandleCheckedCitiesChange">
                                    <el-checkbox v-for="(city, value) in quedialectcities" :label="value" :key="city">{{
                                        city
                                    }}</el-checkbox>
                                </el-checkbox-group>
                            </div>
                        </div>
                        <div class="checkcss">
                            <div class="checkquan">
                                <el-checkbox :indeterminate="phyisIndeterminate" v-model="phycheckAll"
                                    @change="phydirhandleCheckAllChange">理疗</el-checkbox>
                            </div>
                            <div class="checkdan">
                                <el-checkbox-group v-model="phycheckedCities" @change="phydirhandleCheckedCitiesChange">
                                    <el-checkbox v-for="(city, value) in phycities" :label="value" :key="city">{{ city
                                    }}</el-checkbox>
                                </el-checkbox-group>
                            </div>
                        </div>
                    </el-form-item>
                    <el-form-item label="备注：">
                        <el-input type="textarea" placeholder="请输入内容" v-model="CreateProduct.product.remarks"
                            maxlength="50" show-word-limit>
                        </el-input>
                    </el-form-item>
                </el-form>
                <el-button type="primary" class="meddatbut" @click="createprod(ruleFormRef)"
                    :disabled="butno">保存</el-button>
            </div>
            <div class="meddata" v-if="meddata == '理疗'">
                <el-form ref="ruleFormRef" :model="CreateProduct.product" label-width="130px" :size="formSize"
                    :rules="rules">
                    <el-form-item label="商品/服务名称：" required="true" style="text-align:left" prop="productName">
                        <el-input placeholder="请输入商品名称" v-model="CreateProduct.product.productName" />
                    </el-form-item>
                    <el-form-item label="商品图片：" required="true" style="text-align:left">
                        <div class="tableright">
                            <imageUploader @imageBase64="logoImageBase64" @onImageUploaded="logoImageUploaded">
                            </imageUploader>
                            <div class="imgtext">
                                <p>商家图标不允许涉及政治敏感与色情;</p>

                                <p>图片格式只支持：BMP、JPEG、JPG、GIF、PNG，大小不超过2M</p>
                            </div>
                        </div>
                    </el-form-item>
                    <el-form-item label="商品介绍：">
                        <el-input type="textarea" placeholder="请输入内容" v-model="CreateProduct.product.description"
                            maxlength="100" show-word-limit>
                        </el-input>
                    </el-form-item>
                    <el-form-item label="适应症候：" required="true">
                        <div class="checkcss">
                            <div class="checkquan checkquan1">
                                <el-checkbox :indeterminate="isIndeterminate" v-model="checkAll"
                                    @change="handleCheckAllChange">风险疾病</el-checkbox>
                            </div>
                            <div class="checkdan">
                                <el-checkbox-group v-model="checkedCities" @change="handleCheckedCitiesChange">
                                    <el-checkbox v-for="(city, value) in cities" :label="value" :key="city">{{ city
                                    }}</el-checkbox>
                                </el-checkbox-group>
                            </div>
                        </div>
                        <div class="checkcss">
                            <div class="checkquan">
                                <el-checkbox :indeterminate="dirisIndeterminate" v-model="dircheckAll"
                                    @change="dirhandleCheckAllChange">脏腑辨证</el-checkbox>
                            </div>
                            <div class="checkdan">
                                <el-checkbox-group v-model="dircheckedCities" @change="dirhandleCheckedCitiesChange">
                                    <el-checkbox v-for="(city, value) in dircities" :label="value" :key="city">{{ city
                                    }}</el-checkbox>
                                </el-checkbox-group>
                            </div>
                        </div>
                        <div class="checkcss tzcheck">
                            <div class="checkquan">
                                <el-checkbox :indeterminate="quedialectisIndeterminate" v-model="quedialectcheckAll"
                                    @change="dialecticsirhandleCheckAllChange">体质辨证</el-checkbox>
                            </div>
                            <div class="checkdan">
                                <el-checkbox-group v-model="quedialectcheckedCities"
                                    @change="dialecticsdirhandleCheckedCitiesChange">
                                    <el-checkbox v-for="(city, value) in quedialectcities" :label="value" :key="city">{{
                                        city
                                    }}</el-checkbox>
                                </el-checkbox-group>
                            </div>
                        </div>
                        <div class="checkcss">
                            <div class="checkquan">
                                <el-checkbox :indeterminate="phyisIndeterminate" v-model="phycheckAll"
                                    @change="phydirhandleCheckAllChange">理疗</el-checkbox>
                            </div>
                            <div class="checkdan">
                                <el-checkbox-group v-model="phycheckedCities" @change="phydirhandleCheckedCitiesChange">
                                    <el-checkbox v-for="(city, value) in phycities" :label="value" :key="city">{{ city
                                    }}</el-checkbox>
                                </el-checkbox-group>
                            </div>
                        </div>
                    </el-form-item>
                    <el-form-item label="备注：">
                        <el-input type="textarea" placeholder="请输入内容" v-model="CreateProduct.product.remarks"
                            maxlength="50" show-word-limit>
                        </el-input>
                    </el-form-item>
                </el-form>
                <el-button type="primary" class="meddatbut" @click="createprod(ruleFormRef)"
                    :disabled="butno">保存</el-button>
            </div>
        </div>
    </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import imageUploader from '@/components/common/imageUploader.vue';
import { CreateProductRequest } from '@/api/api'
import { getTenantId, getSymptomKeyMap } from '@/utils/storage'
import { ElMessage } from 'element-plus';
import { settopname } from '@/utils/sethometop'

//当前处在那一步
const comtyp = ref(1)

//药品或理疗方案
const meddata = ref('药品')
//接收商品类型
const yaopin = ref('PRODUCT_TYPE_CPD')
//商品分类
const options = ref([{
    label: '中成药',
    value: 'PRODUCT_TYPE_CPD'
}, {
    label: '保健品',
    value: 'PRODUCT_TYPE_HEALTH_CARE_PRODUCT'
}, {
    label: '食品',
    value: 'PRODUCT_TYPE_NUTRITION'
}])

const butno = ref(false)
//表单验证
const ruleFormRef = ref()
const formSize = ref('default')
const rules = reactive({
    productName: [
        { required: true, message: '商品名称不能为空', trigger: 'blur' },
    ],
    drugName: [
        {
            required: true,
            message: '药品名称不能为空',
            trigger: 'blur',
        },
    ],
    approvedNumber: [
        {
            required: true,
            message: '商品字号不能为空',
            trigger: 'blur',
        },
    ],
    contactPhone: [
        {
            required: true,
            message: '联系人手机号不能为空',
            trigger: 'blur',
        }, { min: 11, max: 11, message: '请输入正确的手机号', trigger: 'blur' },
    ],
    drugValidityPeriod: [
        {
            required: true,
            message: '药品准效期不能为空',
            trigger: 'blur',
        }
    ],
})

//风险疾病check
const checkAll = ref(false)
const checkedCities = ref([])
const cities = ref([])
const isIndeterminate = ref(false)
//脏腑辨证check
const dircheckAll = ref(false)
const dircheckedCities = ref([])
const dircities = ref([])
const dirisIndeterminate = ref(false)
//理疗check
const phycheckAll = ref(false)
const phycheckedCities = ref([])
const phycities = ref([])
const phyisIndeterminate = ref(false)
//体质check
const quedialectcheckAll = ref(false)
const quedialectcheckedCities = ref([])
const quedialectcities = ref([])
const quedialectisIndeterminate = ref(false)
const radio = ref(true)
const aa = () => {
    console.log(yaopin.value);
}
const CreateProduct = ref({
    organizationId: getTenantId(),
    product: {
        productType: null,
        productName: null,
        drugName: null,
        isOtc: true,
        approvedNumber: null,
        drugValidityPeriod: null,
        description: null,
        symptomKeys: null,
        remarks: null
    },
    uploadingImage: {
        mime: '',
        image: '',
        filename: ''
    }
})
//处理上传图片
const logoImageBase64 = (val) => {
    console.log(val);
    const imageInfo = val.split(/data:(.*?);base64,(.*?)/g).filter(Boolean)
    const [mime, image] = imageInfo
    CreateProduct.value.uploadingImage.mime = mime
    CreateProduct.value.uploadingImage.image = image
}
//获取上传图片名称
const logoImageUploaded = (val) => {
    CreateProduct.value.uploadingImage.filename = val.name.split('.')[0]
}
//获取map值
const getmep = () => {
    const data = getSymptomKeyMap()
    const key = data.risk_disease_key_map
    const dirkey = data.dirty_dialectics_key_map
    const phykey = data.physical_therapy_key_map
    const dialecticstkey = data.physique_dialectics_key_map
    delete dirkey.Z0022
    cities.value = key
    dircities.value = dirkey
    quedialectcities.value = dialecticstkey
    const Z0004 = data.dirty_dialectics_key_map.Z0004
    const Z0008 = data.dirty_dialectics_key_map.Z0008
    const Z0011 = data.dirty_dialectics_key_map.Z0011
    delete dirkey.Z0004
    delete dirkey.Z0008
    delete dirkey.Z0011
    dircities.value.Z0004 = Z0004
    dircities.value.Z0008 = Z0008
    dircities.value.Z0011 = Z0011
    console.log(data.dirty_dialectics_key_map);
    phycities.value = phykey
}
//风险疾病点击事件
const handleCheckAllChange = (val) => {
    const aa = Object.keys(cities.value)
    checkedCities.value = val ? aa : [];
    isIndeterminate.value = false;
}
const handleCheckedCitiesChange = (val) => {
    const aa = Object.keys(cities.value)
    let checkedCount = val.length;
    checkAll.value = checkedCount === aa.length;
    isIndeterminate.value = checkedCount > 0 && checkedCount < aa.length;
}

//脏腑辨证点击事件
const dirhandleCheckAllChange = (val) => {
    const aa = Object.keys(dircities.value)
    dircheckedCities.value = val ? aa : [];
    dirisIndeterminate.value = false;
}
const dirhandleCheckedCitiesChange = (val) => {
    const aa = Object.keys(dircities.value)
    let checkedCount = val.length;
    dircheckAll.value = checkedCount === aa.length;
    dirisIndeterminate.value = checkedCount > 0 && checkedCount < aa.length;
}

//理疗点击事件
const phydirhandleCheckAllChange = (val) => {
    const aa = Object.keys(phycities.value)
    phycheckedCities.value = val ? aa : [];
    phyisIndeterminate.value = false;
}
const phydirhandleCheckedCitiesChange = (val) => {
    const aa = Object.keys(phycities.value)
    let checkedCount = val.length;
    phycheckAll.value = checkedCount === aa.length;
    phyisIndeterminate.value = checkedCount > 0 && checkedCount < aa.length;
}

//体质点击事件
const dialecticsirhandleCheckAllChange = (val) => {
    const aa = Object.keys(quedialectcities.value)
    quedialectcheckedCities.value = val ? aa : [];
    quedialectisIndeterminate.value = false;
}
const dialecticsdirhandleCheckedCitiesChange = (val) => {
    const aa = Object.keys(quedialectcities.value)
    let checkedCount = val.length;
    quedialectcheckAll.value = checkedCount === aa.length;
    quedialectisIndeterminate.value = checkedCount > 0 && checkedCount < aa.length;
}

//添加商品
const createprod = async (formEl) => {
    // 如果按钮已经处于禁用状态，直接返回，防止重复执行
    if (butno.value) return;

    // 禁用按钮
    butno.value = true;

    try {
        if (!formEl) {
            console.log(formEl);
            ElMessage({
                showClose: true,
                message: '请填写完整信息',
                type: 'error',
            });
            // 恢复按钮状态
            butno.value = false;
            return;
        }

        // 使用Promise包装validate，确保异步流程正确处理
        await new Promise((resolve) => {
            formEl.validate(async (valid, fields) => {
                if (valid) {
                    console.log('submit!');
                    CreateProduct.value.product.symptomKeys = [
                        ...checkedCities.value,
                        ...dircheckedCities.value,
                        ...phycheckedCities.value,
                        ...quedialectcheckedCities.value
                    ];

                    if (CreateProduct.value.uploadingImage.image === '') {
                        ElMessage({
                            showClose: true,
                            message: '请上传图片',
                            type: 'error',
                        });
                        resolve();
                        return;
                    } else if (CreateProduct.value.product.symptomKeys.length === 0) {
                        ElMessage({
                            showClose: true,
                            message: '请选择适应症候',
                            type: 'error',
                        });
                        resolve();
                        return;
                    }

                    CreateProduct.value.product.productType = yaopin.value;
                    try {
                        const data = await CreateProductRequest(CreateProduct.value);

                        if (data.status === 200) {
                            ElMessage({
                                message: '添加商品成功',
                                type: 'success',
                            });
                            window.history.go(-1);
                        } else {
                            ElMessage({
                                showClose: true,
                                message: data.data?.detail || '添加商品失败',
                                type: 'error',
                            });
                        }
                    } catch (requestError) {
                        console.error('接口请求错误:', requestError);
                        ElMessage({
                            showClose: true,
                            message: '接口请求失败，请稍后重试',
                            type: 'error',
                        });
                    } finally {
                        resolve();
                    }
                } else {
                    console.log('error submit!', fields);
                    resolve();
                }
            });
        });
    } catch (error) {
        console.error('提交表单时发生错误:', error);
        ElMessage({
            showClose: true,
            message: '提交过程中发生错误',
            type: 'error',
        });
    } finally {
        // 所有操作完成后才恢复按钮状态
        butno.value = false;
    }
};

//填写商品信息
const shangpintian = (data) => {
    if (data == '理疗') {
        yaopin.value = 'PRODUCT_TYPE_SERVICE'
    } else {
        yaopin.value = 'PRODUCT_TYPE_CPD'
    }
    meddata.value = data
    comtyp.value = 2

}
const topname = () => {
    settopname([{ name: '商品管理', url: '/products' }, { name: '添加商品', url: '' }])
}
topname()
getmep()
</script>
<style lang="scss">
.addgood {
    background-color: #fff;
    min-height: 84vh;
    padding: 60px 100px;
}

.add-top {
    display: flex;
    border: 1px solid #c1c7d0;
    border-radius: 5px;
    overflow: hidden;
    width: 80%;
    margin: 0 auto;

    .toptxt {
        width: 50%;
        height: 40px;
        text-align: center;
        line-height: 40px;
    }

    .addtop_lf {
        background-color: #006be5;
        color: #fff;
    }
}

.addcnt {
    margin-top: 140px;
    display: flex;
    justify-content: space-around;
    padding: 0 100px;

    .addcntdiv {
        width: 300px;
        height: 300px;
        border: 1px solid #c1c7d0;
        border-radius: 10px;
        font-size: 18px;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        cursor: pointer;

        img {
            width: 94px;
            height: 94px;
        }
    }
}

.adddata {
    margin-top: 30px;
    padding: 0px 140px 40px;

    .tableright {

        flex: 0.9;
        display: flex;
        align-items: center;

        .imgtext {
            margin-left: 20px;
            font-weight: 400;
            font-style: normal;
            font-size: 14px;
            color: rgb(193, 199, 208);
            line-height: 28px;
        }
    }
}

.addgood .meddata {
    position: relative;

    .checkcss {
        display: flex;
        flex-direction: column;

        .el-checkbox {
            min-width: 80px;
        }

        .checkdan {
            .el-checkbox__label {
                color: #999999;
                font-weight: 400;

            }
        }

        .checkquan {
            margin-top: 10px;
        }

        .checkquan1 {
            margin-top: 0px;
        }
    }

    .tzcheck {
        .el-checkbox {
            min-width: 210px;
        }
    }

    .meddatbut {
        position: absolute;
        right: 0;
        bottom: -80px;
        width: 103px;
        height: 40px;
        background-color: #1575ee;
    }
}
</style>