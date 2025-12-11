<template>
    <div class="summar-cnt">
        <div class="activate sumdiv">
            <div class="summar-txt">
                使用中商户数/总商户数
            </div>
            <div class="summar-nab">
                {{ Organization.using_teant_amounts }}/{{ Organization.total_tenant_amounts }}
            </div>
        </div>
        <div class="overdue sumdiv">
            <div class="summar-txt">
                当日新增商户数量
            </div>
            <div class="summar-nab">
                {{ Organization.today_added_tenant_count }}
            </div>
        </div>
        <div class="expire sumdiv">
            <div class="summar-txt">
                本月即将到期商户数量
            </div>
            <div class="summar-nab">
                {{ Organization.tenants_expired_amounts }}
            </div>
        </div>
    </div>
    <div class="prodata">
        <div class="client prodatadiv">
            <div class="prodatatop">
                客户数据
            </div>
            <div class="prodatari" @click="tz('CustomerData')">
                详情查询
            </div>
            <div class="clientcnt">
                <div class="clientdiv">
                    <div class="iconimg">
                        <img src="../../assets/image/data_dayvip.png" alt="" class="clienticon">
                    </div>
                    <div class="prodatatxt">累计超30天未检测vip客户数量</div>
                    <div class="prodatanbm" style="color: #3CA1F8;">{{ Organization.customer_overdue_count }}</div>
                    <nav>所有商户</nav>
                </div>
            </div>
        </div>
        <div class="measureds prodatadiv">
            <div class="prodatatop">
                测量数据
            </div>
            <div class="prodatari" @click="tz('MeasurementData')">
                详情查询
            </div>
            <div class="measurecnt">
                <div class="measure">
                    <div class="iconimg">
                        <img src="../../assets/image/data_icon_dayvip.png" alt="" class="measureicon">
                    </div>
                    <div class="prodatatxt">当日体验测量次数</div>
                    <div class="prodatanbm" style="color: #3CA1F8;">{{
                    OrganizationMeasurement.today_temp_customer_measurement_count }}
                    </div>
                    <nav>所有商户</nav>
                </div>
                <div class="measure">
                    <div class="iconimg">
                        <img src="../../assets/image/data_icon_daynotvip.png" alt="" class="measureicon">
                    </div>
                    <div class="prodatatxt">当日vip测量次数</div>
                    <div class="prodatanbm" style="color: #9C88FF;">{{
                    OrganizationMeasurement.today_customer_measurement_count }}
                    </div>
                    <nav>所有商户</nav>
                </div>
                <div class="measure">
                    <div class="iconimg">
                        <img src="../../assets/image/data_icon_vip.png" alt="" class="measureicon">
                    </div>
                    <div class="prodatatxt">累计vip测量次数</div>
                    <div class="prodatanbm" style="color: #3CA1F8;">{{
                    OrganizationMeasurement.year_customer_measurement_count }}
                    </div>
                    <nav>所有商户</nav>
                </div>
                <div class="measure">
                    <div class="iconimg">
                        <img src="../../assets/image/data_icon_notvip.png" alt="" class="measureicon">
                    </div>
                    <div class="prodatatxt">累计散客测量次数</div>
                    <div class="prodatanbm" style="color: #9C88FF;">{{
                        OrganizationMeasurement.year_temp_customer_measurement_count }}
                    </div>
                    <nav>所有商户</nav>
                </div>
            </div>
        </div>
    </div>
</template>
<script setup>
import { ref } from 'vue'
import stort from '@/store/generindex'
import { useRouter } from 'vue-router';
import { getTenantId } from '@/utils/storage'
import { GetOrganizationMeasurementSummaryRequest } from '@/api/api'


const $router = useRouter();

const Organization = ref(stort.state.Organization)
const OrganizationMeasurement = ref({
    today_temp_customer_measurement_count: 0,
    today_customer_measurement_count: 0,
    year_customer_measurement_count: 0,
    year_temp_customer_measurement_count: 0,
})

const GetOrganizationMeasurementSummary = async () => {
    const data = await GetOrganizationMeasurementSummaryRequest({ organizationId: getTenantId() })

    OrganizationMeasurement.value.today_temp_customer_measurement_count = data.data.today_temp_customer_measurement_count === undefined ? '0' : data.data.today_temp_customer_measurement_count

    OrganizationMeasurement.value.today_customer_measurement_count = data.data.today_customer_measurement_count === undefined ? '0' : data.data.today_customer_measurement_count

    OrganizationMeasurement.value.year_customer_measurement_count = data.data.year_customer_measurement_count === undefined ? '0' : data.data.year_customer_measurement_count

    OrganizationMeasurement.value.year_temp_customer_measurement_count = data.data.year_temp_customer_measurement_count === undefined ? '0' : data.data.year_temp_customer_measurement_count

}

const giedata = () => {
    stort.commit('getlist')
    GetOrganizationMeasurementSummary()
}
const tz = (to) => {
    $router.push(
        {
            name: to
        }
    )
}
giedata()
</script>
<style lang="scss" scoped>
.summar-cnt {
    display: flex;
    justify-content: space-between;

    .sumdiv {
        width: 35%;
        height: 160px;
        border-radius: 10px;
        margin-left: 20px;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        font-weight: 400;
        font-style: normal;
        color: #FFFFFF;
        text-align: center;
        background-size: cover;

        .summar-txt {
            font-size: 15px;
        }

        .summar-nab {
            font-size: 36px;
        }
    }

    .activate {
        background-image: url(../../assets/icons/overview_pic_1.svg);
        margin-left: 0;
    }

    .overdue {
        background-image: url(../../assets/icons/overview_pic_2.svg);

    }

    .expire {
        background-image: url(../../assets/icons/overview_pic_3.svg);

    }
}

.prodata {
    margin-top: 30px;
    display: flex;
    justify-content: space-between;

    .client {
        margin-right: 40px;
        width: 35%;
        height: 228px;
    }

    .measureds {
        width: 60%;
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

        .clientcnt,
        .measurecnt {
            display: flex;
            flex: 1;
            justify-content: space-around;

            .measure,
            .clientdiv {
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: center;

                .iconimg {
                    height: 30px;
                    line-height: 30px;
                }
            }

            .clienticon {
                width: 30px;
                height: 27px;
            }

            .measureicon {
                width: 12px;
                height: 12px;
            }

            .prodatatxt {
                margin-top: 10px;
                font-size: 13px;
            }

            .prodatanbm {
                font-weight: 700;
                font-style: normal;
                font-size: 30px;

                text-align: center;
            }
        }
    }
}

nav {
    font-weight: 400;
    font-style: normal;
    font-size: 12px;
    color: #C1C7D0;
    text-align: center;
}
</style>