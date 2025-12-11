<template>
    <div class="edgoods">
        <div class="edgoodsdata">
            <div class="meddata" v-if="meddata == '药品'">
                <el-form ref="form" :model="form" label-width="130px">
                    <el-form-item label="商品类型：" required="true" style="text-align:left">
                        <span
                            v-if="CreateProduct.product.product_type == 'PRODUCT_TYPE_CPD' || CreateProduct.product.product_type == 'PRODUCT_TYPE_HEALTH_CARE_PRODUCT' || CreateProduct.product.product_type == 'PRODUCT_TYPE_NUTRITION'">药品</span>
                        <span v-if="CreateProduct.product.product_type == 'PRODUCT_TYPE_SERVICE'">理疗服务</span>
                    </el-form-item>
                    <el-form-item label="商品分类：" required="true" style="text-align:left">
                        <span v-if="CreateProduct.product.product_type == 'PRODUCT_TYPE_CPD'">中成药</span>
                        <span v-if="CreateProduct.product.product_type == 'PRODUCT_TYPE_HEALTH_CARE_PRODUCT'">保健品</span>
                        <span v-if="CreateProduct.product.product_type == 'PRODUCT_TYPE_NUTRITION'">营养食品</span>
                        <span v-if="CreateProduct.product.product_type == 'PRODUCT_TYPE_SERVICE'">理疗服务</span>
                    </el-form-item>
                    <el-form-item label="商品名称：" required="true" style="text-align:left">
                        <el-input placeholder="请输入商品名称" v-model="CreateProduct.product.product_name" />
                    </el-form-item>
                    <el-form-item label="药品名称：" required="true" style="text-align:left"
                        v-if="yaopin == 'PRODUCT_TYPE_CPD'">
                        <el-input placeholder="请输入药品名称" v-model="CreateProduct.product.product_name" />
                    </el-form-item>
                    <el-form-item label="药品类型：" required="true" style="text-align:left"
                        v-if="yaopin == 'PRODUCT_TYPE_CPD'">
                        <el-radio v-model="CreateProduct.product.is_otc" :label="true">处方药</el-radio>
                        <el-radio v-model="CreateProduct.product.is_otc" :label="false">非处方药</el-radio>
                    </el-form-item>
                    <el-form-item label="国药准字号：" required="true" style="text-align:left"
                        v-if="yaopin == 'PRODUCT_TYPE_CPD'">
                        <el-input placeholder="请输入国药准字号" v-model="CreateProduct.product.approved_number" />
                    </el-form-item>
                    <el-form-item label="健字号：" required="true" style="text-align:left"
                        v-if="yaopin == 'PRODUCT_TYPE_HEALTH_CARE_PRODUCT'">
                        <el-input placeholder="请输入健字号" v-model="CreateProduct.product.approved_number" />
                    </el-form-item>
                    <el-form-item label="食字号：" required="true" style="text-align:left"
                        v-if="yaopin == 'PRODUCT_TYPE_NUTRITION'">
                        <el-input placeholder="请输入食字号" v-model="CreateProduct.product.approved_number" />
                    </el-form-item>
                    <el-form-item label="药品准效期：" required="true" style="text-align:left"
                        v-if="yaopin == 'PRODUCT_TYPE_CPD'">
                        <el-date-picker :editable="false" placeholder="请选择药品准效期" type="date"
                            v-model="CreateProduct.product.drug_validity_period"></el-date-picker>

                    </el-form-item>
                    <el-form-item label="商品图片：" required="true" style="text-align:left">
                        <div class="tableright">
                            <imageUploader @imageBase64="logoImageBase64" @onImageUploaded="logoImageUploaded"
                                :imgurl="CreateProduct.product.image_url">
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
                    <el-form-item label="适应症候：">
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
                <div class="meddatbut">
                    <el-button class="buttonel" @click="fanhui">取消</el-button>
                    <el-button class="buttonel butbc" type="primary" @click="updateore">保存</el-button>
                </div>

            </div>
            <div class="meddata" v-if="meddata == '理疗'">
                <el-form ref="form" :model="form" label-width="130px">
                    <el-form-item label="商品/服务名称：" required="true" style="text-align:left">
                        <el-input placeholder="请输入商品/服务名称" />
                    </el-form-item>
                    <el-form-item label="商品图片：" required="true" style="text-align:left">
                        <div class="tableright">
                            <imageUploader></imageUploader>
                            <div class="imgtext">
                                <p>商家图标不允许涉及政治敏感与色情;</p>
                                <p>图片格式只支持：BMP、JPEG、JPG、GIF、PNG，大小不超过2M</p>
                            </div>
                        </div>
                    </el-form-item>
                    <el-form-item label="商品介绍：">
                        <el-input type="textarea" placeholder="请输入内容" v-model="textarea" maxlength="100"
                            show-word-limit>
                        </el-input>
                    </el-form-item>
                    <el-form-item label="适应症候：">
                        <el-checkbox :indeterminate="isIndeterminate" v-model="checkAll"
                            @change="handleCheckAllChange">风险疾病</el-checkbox>
                        <div style="margin: 15px 0;"></div>
                        <el-checkbox-group v-model="checkedCities" @change="handleCheckedCitiesChange">
                            <el-checkbox v-for="city in cities" :label="city" :key="city">{{ city }}</el-checkbox>
                        </el-checkbox-group>
                    </el-form-item>
                    <el-form-item label="备注：">
                        <el-input type="textarea" placeholder="请输入内容" v-model="textarea" maxlength="100"
                            show-word-limit>
                        </el-input>
                    </el-form-item>
                </el-form>
            </div>
        </div>
    </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import imageUploader from '@/components/common/imageUploader.vue';
import { ElMessage } from 'element-plus';
import { getTenantId, getSymptomKeyMap } from '@/utils/storage'
import {
    UpdateProductRequest
} from '@/api/api'
import { settopname } from '../../utils/sethometop'


//药品或理疗方案
const meddata = ref('药品')
//接收商品类型
const yaopin = ref('PRODUCT_TYPE_SERVICE')

const reviewtype = ref(JSON.parse(history.state.keyword))
console.log(reviewtype.value);

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
const uploadingImage = ref({
    mime: '',
    image: '',
    filename: ''
})
const CreateProduct = ref({
    organizationId: getTenantId(),
    product: reviewtype.value
})
//处理上传图片
const logoImageBase64 = (val) => {
    console.log(val);
    const imageInfo = val.split(/data:(.*?);base64,(.*?)/g).filter(Boolean)
    const [mime, image] = imageInfo
    uploadingImage.value.mime = mime
    uploadingImage.value.image = image
}
//获取上传图片名称
const logoImageUploaded = (val) => {
    uploadingImage.value.filename = val.name.split('.')[0]
}
//获取map值
const getmep = () => {
    CreateProduct.value.product.is_otc = CreateProduct.value.product.is_otc == undefined ? false : CreateProduct.value.product.is_otc
    const data = getSymptomKeyMap()
    let key = data.risk_disease_key_map
    let dirkey = data.dirty_dialectics_key_map
    let phykey = data.physical_therapy_key_map
    let dialecticstkey = data.physique_dialectics_key_map
    delete dirkey.Z0022
    cities.value = key
    dircities.value = dirkey
    phycities.value = phykey
    quedialectcities.value = dialecticstkey
    reviewtype.value.symptom_keys.forEach(v => {
        for (let j in key) {
            if (j == v) {
                checkedCities.value.push(j)
                handleCheckedCitiesChange(j)
            }
        }
        for (let j in dirkey) {
            if (j == v) {
                dircheckedCities.value.push(j)
                dirhandleCheckedCitiesChange(j)
            }
        }
        for (let j in phykey) {
            if (j == v) {
                phycheckedCities.value.push(j)
                phydirhandleCheckedCitiesChange(j)
            }
        }
        for (let j in dialecticstkey) {
            if (j == v) {
                quedialectcheckedCities.value.push(j)
                dialecticsdirhandleCheckedCitiesChange(j)
            }
        }
    });
    yaopin.value = reviewtype.value.product_type
    console.log(data);
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
const updateore = () => {
    CreateProduct.value.product.symptom_keys = [...checkedCities.value, ...dircheckedCities.value, ...phycheckedCities.value, ...quedialectcheckedCities.value]
    if (uploadingImage.value.mime == '') {
        UpdateProductRequest(CreateProduct.value)
            .then(res => {
                console.log(res);
                ElMessage({
                    message: '商品更新成功',
                    type: 'success',
                })
                window.history.go(-1)
            })
            .catch(error => {
                console.log(error);
                ElMessage({
                    showClose: true,
                    message: error.detail,
                    type: 'error',
                })

            })
    } else {
        CreateProduct.value.uploadingImage = uploadingImage.value
        UpdateProductRequest(CreateProduct.value).then(res => {
            ElMessage({
                message: '商品更新成功',
                type: 'success',
            })
            window.history.go(-1)
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
    console.log(CreateProduct.value);

}
const topname = () => {
    settopname([{ name: '商品管理', url: '/products' }, { name: '编辑商品', url: '' }])
}
const fanhui = () => {
    window.history.go(-1)
}
topname()
getmep()
</script>
<style lang="scss">
.edgoods {
    background-color: #fff;
    padding: 20px 40px;

    .edgoodsdata {
        width: 60%;
        margin: auto;
    }
}

.meddata {
    padding-bottom: 80px;
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
        bottom: 0px;
        display: flex;

        .buttonel {
            width: 103px;
            height: 40px;
        }

        .butbc {
            background-color: #1575ee;
        }
    }

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
</style>