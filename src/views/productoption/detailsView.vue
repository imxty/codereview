<template>
    <div class="detail">
        <div class="dataitop">
            <div class="detatoplf">
                方案名称：{{ reviewtype.name }}
            </div>
            <div class="detatoprg">
                创建时间：{{ reviewtype.time }}
            </div>
        </div>
        <div class="datailcnt">
            <el-table :data="reviewtype.data" border style="width: 100%">
                <el-table-column prop="keyname" label="症候" align="center" width="300">
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
        <div class="datailbtm">
            <el-button plain @click="fanhui">返回列表</el-button>
        </div>
    </div>
</template>
<script setup>
import { reactive, ref } from 'vue';
import { settopname } from '@/utils/sethometop'
import { useRouter } from 'vue-router';

const $router = useRouter();
const reviewtype = ref(JSON.parse(history.state?.keyword))
console.log(reviewtype.value);
const topname = () => {
    settopname([{ name: '商品推荐方案', url: '/recommendedProposal' }, { name: '方案详情', url: '' }])
}

const fanhui = () => {
    $router.push({
        name: 'Recomscheme',
    })
}
topname()
</script>
<style lang="scss">
.detail {
    background-color: #fff;
    padding: 20px;

    .dataitop {
        display: flex;
        justify-content: space-between;
        padding: 20px 40px;

        .detatoplf {
            font-weight: 700;
            font-size: 18px;
        }

        .detatoprg {
            font-size: 14px;
            color: #666666;
        }
    }

    .datatab {
        width: 100%;
        text-align: center;

        .tablelf {
            width: 240px;
        }

        .tablerg {
            position: relative;
        }
    }

    table {
        border-collapse: collapse;
    }

    table tr {
        border: 1px solid #c1c7d0;
        ;
    }

    th {
        color: #c1c7d0;
        padding: 10px 0;
    }

    td {
        font-size: 14px;
        padding: 10px 0;
        color: #333333;
        overflow: hidden;
    }
}

.lock {
    font-size: 14px;
    color: #f56c6c;
    text-align: center;
    width: 117px;
    height: 30px;
    background-color: #fde2e2;
    line-height: 30px;
    position: absolute;
    right: -2px;
    top: 0;
    border-radius: 5px;
}

.replace {
    width: 316px;
    height: 30px;
    font-size: 14px;
    color: #3ca1f8;
    background-color: #b1d9fc;
    line-height: 30px;
    position: absolute;
    right: -2px;
    top: 0;
    border-radius: 5px;
}

.datailbtm {
    text-align: right;
    margin-top: 50px;
}
</style>