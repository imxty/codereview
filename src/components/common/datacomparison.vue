<template>
    <div class="prodata">
        <div class="client prodatadiv">
            <div class="prodatatop">
                商户数据对比
            </div>
            <div class="prodatari" @click="tz('TenantCompare')">
                详情查询
            </div>
            <div class="qie">
                <div class="qieleft">
                    <div class="qieclient" :class="[{ select: select == 'customer_rank' }]" @click="Changelisttenant">
                        客户数据
                    </div>
                    <div class="measure" style="cursor: pointer;" :class="[{ select: select == 'measurement_rank' }]"
                        @click="Changelistnamber">测量次数
                    </div>
                </div>
                <div class="qieright">
                    <div class="qiedata">
                        <div class="datasho">
                            客户数据显示Top10
                        </div>
                        <div class="minutes">
                            <span class="minuico" style="margin-right: 5px;"></span>上月
                        </div>
                        <div class="minutex">
                            <span class="minuico" style="margin-right: 5px;"></span>本月
                        </div>
                    </div>
                </div>
            </div>
            <div class="clientcnt">
                <div class="zwsj" v-if="!datalength">
                    暂无数据
                </div>
                <div id="myEcharts" ref="box" style=" height: 300px" v-show="datalength"></div>
            </div>

        </div>
    </div>
</template>
<script setup>
import * as echarts from "echarts";
import { onMounted, ref } from "vue";
import { ListTopTenantRankRequest } from '@/api/api'
import { getTenantId } from '@/utils/storage'
import { useRouter } from 'vue-router';


const $router = useRouter();

//创建图表
let box = ref(null)
const select = ref('customer_rank')

const datalength = ref(true)
onMounted(async () => {
    const rs = {
        organizationId: getTenantId()
    }
    let dataname = []
    let datada = []
    let lastMonthAmounts = []
    const add = await ListTopTenantRankRequest(rs)
    console.log(add);
    if (add.data.rank_count) {
        add.data.rank_count.forEach(data => {
            const last_month_customer_amounts = data.last_month_customer_amounts == undefined ? 0 : data.last_month_customer_amounts
            const customer_amounts = data.customer_amounts == undefined ? 0 : data.customer_amounts
            console.log(data);
            dataname.push(data.store_name)
            datada.push(customer_amounts)
            lastMonthAmounts.push(last_month_customer_amounts)
        });
        datalength.value = true
    } else {
        datalength.value = false
    }
    let myecharts = echarts.init(box.value);
    let option = {
        color: ['#BAD6FD', '#3fa0f9', '#fac858', '#ee6666', '#73c0de', '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc'],
        xAxis: [
            {
                type: 'category',
                data: dataname,
                axisTick: {
                    show: false
                },
                axisLine: {
                    show: false,
                    lineStyle: { color: '#333' }
                },
            }
        ],
        yAxis: {
            type: 'value',
            axisLine: {
                show: false
            },


        },
        series: [
            {
                data: lastMonthAmounts,
                type: 'bar',
                barWidth: 30,
                itemStyle: {
                    borderRadius: [5, 5, 0, 0]
                },
            },
            {
                data: datada,
                type: 'bar',
                barWidth: 30,
                itemStyle: {
                    borderRadius: [5, 5, 0, 0]
                },
            },

        ]
    }
    console.log(lastMonthAmounts.length);
    myecharts.setOption(option)
});
const Changelistnamber = async () => {
    datalength.value = true
    select.value = 'measurement_rank'
    const rs = {
        organizationId: getTenantId()
    }
    let dataname = []
    let datada = []
    let lastMonthAmounts = []
    const add = await ListTopTenantRankRequest(rs)
    console.log(add);
    if (add.data.rank_count) {
        add.data.rank_count.forEach(data => {
            console.log(data);
            dataname.push(data.store_name)
            datada.push(data.measurement_amounts)
            lastMonthAmounts.push(data.last_month_measurement_amounts)
        });
        datalength.value = true
    } else {
        datalength.value = false
    }

    let myecharts = echarts.init(box.value);
    let option = {
        color: ['#BAD6FD', '#3fa0f9', '#fac858', '#ee6666', '#73c0de', '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc'],
        xAxis: [
            {
                type: 'category',
                data: dataname,
                axisTick: {
                    show: false,
                    alignWithLabel: true
                }
            }
        ],
        yAxis: {
            type: 'value',
            axisLine: {
                show: false
            },


        },
        series: [
            {
                data: lastMonthAmounts,
                type: 'bar',
                barWidth: 30,
                itemStyle: {
                    borderRadius: [5, 5, 0, 0]
                },
            },
            {
                data: datada,
                type: 'bar',
                barWidth: 30,
                itemStyle: {
                    borderRadius: [5, 5, 0, 0]
                },
            },

        ]
    }
    myecharts.setOption(option)
}
const Changelisttenant = async () => {
    datalength.value = true
    select.value = 'customer_rank'
    const rs = {
        organizationId: getTenantId()
    }
    let dataname = []
    let datada = []
    let lastMonthAmounts = []
    const add = await ListTopTenantRankRequest(rs)
    console.log(add);
    add.data.rank_count.forEach(data => {
        const last_month_customer_amounts = data.last_month_customer_amounts == undefined ? 0 : data.last_month_customer_amounts
        const customer_amounts = data.customer_amounts == undefined ? 0 : data.customer_amounts
        console.log(data);
        dataname.push(data.store_name)
        datada.push(customer_amounts)
        lastMonthAmounts.push(last_month_customer_amounts)
    });

    let myecharts = echarts.init(box.value);
    let option = {
        color: ['#BAD6FD', '#3fa0f9', '#fac858', '#ee6666', '#73c0de', '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc'],
        xAxis: [
            {
                type: 'category',
                data: dataname,
                axisTick: {
                    show: false,
                    alignWithLabel: true
                }
            }
        ],
        yAxis: {
            type: 'value',
            axisLine: {
                show: false
            },


        },
        series: [
            {
                data: lastMonthAmounts,
                type: 'bar',
                barWidth: 30,
                itemStyle: {
                    borderRadius: [5, 5, 0, 0]
                },
            },
            {
                data: datada,
                type: 'bar',
                barWidth: 30,
                itemStyle: {
                    borderRadius: [5, 5, 0, 0]
                },
            },

        ]
    }
    myecharts.setOption(option)
}

const tz = (to) => {
    $router.push(
        {
            name: to
        }
    )
}

</script>
<style lang="scss">
.prodata {
    margin-top: 30px;
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

    .prodatatop {
        height: 70px;
        line-height: 70px;
        font-weight: 700;
        font-size: 18px;
        padding-left: 40px;
        border-bottom: 1px solid rgba(238, 238, 238, 1);
        ;
    }

    .prodatari {
        font-weight: 650;
        font-style: normal;
        color: #1575EE;
        position: absolute;
        right: 10px;
        top: 25px;
        cursor: pointer;
    }

    .qie {
        display: flex;
        justify-content: space-between;
        padding-left: 40px;
        margin-top: 20px;
        ;
        height: 30px;

        .qieleft {
            display: flex;
            font-weight: 400;
            font-style: normal;
            font-size: 16px;

            .qieclient {
                margin-right: 30px;
                cursor: pointer;
            }
        }

        .qieright {
            width: 300px;

            .qiedata {
                display: flex;
                justify-content: space-evenly;

                .datasho {
                    font-size: 11px;
                    color: #C1C7D0;
                }

                .minuico {
                    display: block;
                    width: 11px;
                    height: 11px;
                }

                .minutes {
                    font-size: 11px;
                    display: flex;
                    color: #C1C7D0;
                    align-items: center;

                    .minuico {
                        background-color: #BAD6FD;
                    }
                }

                .minutex {
                    font-size: 11px;
                    display: flex;
                    align-items: center;

                    .minuico {
                        background-color: #3FA0F9;
                    }
                }


            }
        }

        .select {
            color: #1575EE;
            border-bottom: 3px solid #1575EE;
        }
    }
}

.zwsj {
    height: 300px;
    text-align: center;
    line-height: 300px;
    color: #C1C7D0;
}
</style>