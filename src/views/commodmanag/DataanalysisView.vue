<template>
    <div class="dataanal">
        <div class="detacnt">
            <div class="detacnttop">
                基本信息
            </div>
            <div class="cntdata">
                <div class="cntbody">
                    <div class="bodyimg">
                        <img :src="reviewtype.image_url" alt="" style="width:180px;height: 180px;">
                    </div>
                    <div class="bodybg">
                        <div>
                            <img src="../../assets/image/icon_exposure.png" alt="" style="width: 30px;height: 30px;">
                        </div>
                        <div>
                            商品累计曝光次数
                        </div>
                        <div>
                            {{ datalist?.product_instant_total_count }}
                        </div>
                    </div>
                    <div class="bodydata">
                        <el-table :data="datalist?.symptom_exposure" style="width: 100%" :border="true"
                            :header-cell-style="{ 'text-align': 'center' }" :cell-style="{ 'text-align': 'center' }"
                            max-height="250">
                            <el-table-column prop="symptom_name" label="相关症候">
                                <template #default="scope">
                                    {{ keynames(scope.row.symptom_name) }}
                                </template>
                            </el-table-column>
                            <el-table-column label="推荐商品" prop="exposure_count"></el-table-column>
                        </el-table>
                    </div>
                </div>
            </div>
        </div>
        <div class="detacnt">
            <div class="cntdata">
                <div id="main" ref="line" class="line">
                </div>
            </div>
        </div>
    </div>
</template>
<script setup>
import { ref, onMounted } from 'vue';
import * as echarts from "echarts";
import { settopname } from '@/utils/sethometop'
import {
    GetProductStatisticsRequest
} from '@/api/api'
import { getTenantId, getKeyMaplist } from '@/utils/storage'

const line = ref()
const reviewtype = ref(JSON.parse(history.state.keyword))
const datalist = ref()
console.log(reviewtype.value);
const topname = () => {
    settopname([{ name: '商品管理', url: '/products' }, { name: '数据统计', url: '' }])
}

const GetProductStatistics = async () => {
    const rs = {
        organizationId: getTenantId(),
        productId: reviewtype.value.product_id
    }
    const data = await GetProductStatisticsRequest(rs)
    data.data.productExposure = data.data.productExposure == undefined ? 0 : data.data.productExposure
    console.log(data);
    datalist.value = data.data
    const exposudata = []
    const exposuname = []
    data.data.product_exposure.forEach(v => {
        console.log(v);
        exposudata.push(v.exposure_count)
        exposuname.push(v.exposed_date.month + '-' + v.exposed_date.day)
    })
    var lineChart = echarts.init(line.value);
    var lineOption = (lineOption = {
        tooltip: {
            trigger: "axis",
        },

        legend: {
            textStyle: {
                color: "#4c9bfd", // 图例文字颜色
            },
            right: "10%", // 距离右边10%
            // 如果series 里面设置了name，那么此时图里组件的data可以省略！！！
            // data: ['Email', 'Union Ads']
        },
        grid: {
            top: "20%",
            left: "3%",
            right: "4%",
            bottom: "3%", //显示边框
            borderColor: "#012f4a", // 边框颜色
            containLabel: true,
        },

        xAxis: {
            type: "category",
            boundaryGap: false, // 去除轴内间距
            axisTick: {
                show: false, // 去除刻度线
            },
            axisLabel: {
                color: "#C1C7D0", // 文本颜色
            },
            axisLine: {
                show: false, // 去除轴线
            },
            data: exposuname,
        },
        yAxis: {
            type: "value",
            axisTick: {
                show: false, // 去除刻度
            },
            axisLabel: {
                color: "#C1C7D0", // 文字颜色
            },
            splitLine: {
                lineStyle: {
                    color: "#C1C7D0", // 分割线颜色
                },
            },
        },

        color: ["#50adf5"], //两条曲线改变颜色
        series: [
            {
                symbol: 'circle',
                type: "line",
                stack: "总量",
                data: exposudata,
                itemStyle: {
                    normal: {
                        color: '#50adf5',//拐点颜色
                        borderColor: '#50adf5',//拐点边框颜色
                        borderWidth: 4//拐点边框大小
                    },
                    emphasis: {          //突出效果配置(鼠标置于拐点上时)
                        borderWidth: 10,         //  阴影渐变范围控制
                    },
                },
                areaStyle: {
                    color: {
                        type: 'linear',
                        x: 0,
                        y: 0,
                        x2: 0,
                        y2: 1,
                        colorStops: [{
                            offset: 0, color: '#dff' // 起始颜色和透明度
                        }, {
                            offset: 0.7, color: 'rgba(77,201,251,0)' // 结束颜色和透明度
                        }]
                    }
                },
            },
        ],
    });
    lineChart.setOption(lineOption);
}
const keynames = (data) => {
    let name = ''
    const v = getKeyMaplist()
    v.forEach(v => {
        if (Object.keys(v) == data) {
            name = v[data]
        }
    });
    console.log(name);
    return name
}
onMounted(() => {
    GetProductStatistics()
});
keynames()
topname()
</script>
<style lang="scss">
.dataanal {
    .detacnt {
        background-color: #fff;
        border: 1px solid #eee;
        border-radius: 10px;
        margin-top: 30px;

        .detacnttop {
            font-size: 18px;
            font-weight: 700;
            padding: 20px;
            border-bottom: 1px solid #eee;
        }

        .cntdata {
            padding: 20px;

            .cz {
                color: #1575ee;
                display: flex;
                justify-content: space-between;
            }
        }
    }

    .cntbody {
        display: flex;
        justify-content: flex-start;
        align-items: center;

        .bodyimg {
            margin-right: 90px;
        }

        .bodybg {
            display: flex;
            justify-content: space-around;
            flex-direction: column;
            align-items: center;
            height: 100px;
            margin-right: 90px;
        }

        .bodydata {
            flex: 1;
        }
    }

    .line {
        width: 100%;
        height: 470px;
    }
}
</style>