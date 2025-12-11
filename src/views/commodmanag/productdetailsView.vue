<template>
    <div class="product">
        <div class="meddata" v-if="option == '药品'">
            <div>
                <div class="meddatajs">商品类型：</div>
                <div class="meddatetxt">
                    <span
                        v-if="reviewtype.product_type == 'PRODUCT_TYPE_CPD' || reviewtype.product_type == 'PRODUCT_TYPE_HEALTH_CARE_PRODUCT' || reviewtype.product_type == 'PRODUCT_TYPE_NUTRITION'">药品</span>
                    <span v-if="reviewtype.product_type == 'PRODUCT_TYPE_SERVICE'">理疗服务</span>
                </div>
            </div>
            <div>
                <div class="meddatajs">商品分类：</div>
                <div class="meddatetxt">
                    <span v-if="reviewtype.product_type == 'PRODUCT_TYPE_CPD'">中成药</span>
                    <span v-if="reviewtype.product_type == 'PRODUCT_TYPE_HEALTH_CARE_PRODUCT'">保健品</span>
                    <span v-if="reviewtype.product_type == 'PRODUCT_TYPE_NUTRITION'">营养食品</span>
                    <span v-if="reviewtype.product_type == 'PRODUCT_TYPE_SERVICE'">理疗服务</span>

                </div>
            </div>

            <div>
                <div class="meddatajs">商品名称：</div>
                <div class="meddatetxt">{{ reviewtype.product_name }}</div>
            </div>
            <div v-if="reviewtype.product_type == 'PRODUCT_TYPE_CPD'">
                <div class="meddatajs">药品名称：</div>
                <div class="meddatetxt">{{ reviewtype.product_name }}</div>
            </div>
            <div v-if="reviewtype.product_type == 'PRODUCT_TYPE_CPD'">
                <div class="meddatajs" >药品类型：</div>
                <div class="meddatetxt">
                    <span v-if="reviewtype.is_otc">处方药</span>
                    <span v-if="!reviewtype.is_otc">非处方药</span>
                </div>
            </div>
            <div  v-if="reviewtype.product_type == 'PRODUCT_TYPE_CPD'">
                <div class="meddatajs">国药标子号：</div>
                <div class="meddatetxt">{{ reviewtype.approved_number }}</div>
            </div>
            <div v-if="reviewtype.product_type == 'PRODUCT_TYPE_CPD'">
                <div class="meddatajs">药品准效期：</div>
                <div class="meddatetxt">{{ fermitTime(reviewtype.drug_validity_period )}}</div>
            </div>
            <div>
                <div class="meddatajs">商品图片：</div>
                <div class="meddateimg">
                    <img :src="reviewtype.image_url" alt="" style="width: 150px;">
                </div>
            </div>
            <div>
                <div class="meddatajs">商品介绍：</div>
                <div class="meddatetxt">{{reviewtype.description}}</div>
            </div>
            <div>
                <div class="meddatajs">适应症候：</div>
                <div class="meddatetxt">
                    <span v-for="(v, i) in reviewtype.symptom_keys" :key="v">
                        {{ v }}
                        <span v-if="i < reviewtype.symptom_keys.length-1">,</span>
                    </span>
                </div>
            </div>
            <div>
                <div class="meddatajs">备注：</div>
                <div class="meddatetxt" v-if="reviewtype.remarks">{{ reviewtype.remarks }} </div>
                <div class="meddatetxt" v-else></div>
            </div>
        </div>
        <div class="meddata" v-if="option == '理疗'">
            <div>
                <div class="meddatajs">商品类型：</div>
                <div class="meddatetxt">药品</div>
            </div>
            <div>
                <div class="meddatajs">商品/服务名称：</div>
                <div class="meddatetxt">脾胃温热灸</div>
            </div>
            <div>
                <div class="meddatajs">商品图片：</div>
                <div class="meddateimg"></div>
            </div>
            <div>
                <div class="meddatajs">商品介绍：</div>
                <div class="meddatetxt">脾胃部及相关穴位艾灸，经络疏通</div>
            </div>
            <div>
                <div class="meddatajs">适应症候：</div>
                <div class="meddatetxt">血糖、血脂、心血亏虚、血瘀</div>
            </div>
            <div>
                <div class="meddatajs">备注：</div>
                <div class="meddatetxt"> 无 </div>
            </div>
        </div>
    </div>
</template>
<script setup>
import { ref } from 'vue';
import { useRoute } from 'vue-router';
import { getKeyMaplist } from '@/utils/storage'
import {settopname} from  '@/utils/sethometop'
import { fermitTime } from '@/utils/fermitTime'


//设置路由对象
const $router = useRoute();

const reviewtype = ref(JSON.parse(history.state.keyword))
console.log(reviewtype.value);
//模拟药品类型
const meddata = ref('药品')
//模拟方案类型
const option = ref('药品')

const getmaplist = () => {
    const list = []
    const a = getKeyMaplist()
    reviewtype.value.symptom_keys.forEach(v => {
        a.forEach(h => {
            if (Object.keys(h) == v) {
                console.log(Object.values(h));
                list.push(Object.values(h)[0])
            }
        })
    })
    reviewtype.value.symptom_keys = list
    console.log(reviewtype.value);

}
const topname = ()=>{
   settopname([{name:'商品管理',url:'/products'},{name:'商品详情',url:''}])
}
topname()
getmaplist()
</script>
<style lang="scss">
.meddata {
    padding: 40px 40px;
    background-color: #fff;
}

.meddata>div {
    font-size: 14px;
    margin-top: 30px;
    display: flex;

    .meddatajs {
        flex: 0.1;
        color: rgb(193, 199, 208);
    }

    .meddatetxt,
    .meddateimg {

        flex: 0.9;
        margin-left: 50px;
        border-bottom: 1px solid #ccc;
        min-height: 40px;
    }
}
</style>