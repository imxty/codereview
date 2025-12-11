<template>
    <div class="detail">
        <div class="dataitop">
            <div class="detatoplf">
                <div class="fx" :class="[{ fxxz: symptomTypeLabel == '风险疾病' }]" @click="lablexz('风险疾病')">风险疾病</div>
                <div class="zf" :class="[{ zfxz: symptomTypeLabel == '脏腑辨证' }]" @click="lablexz('脏腑辨证')">脏腑辨证</div>
                <div class="tz" :class="[{ tzxz: symptomTypeLabel == '体质辨证' }]" @click="lablexz('体质辨证')">体质辨证</div>
                <div class="ll" :class="[{ llxz: symptomTypeLabel == '理疗' }]" @click="lablexz('理疗')">理疗</div>
            </div>
        </div>
        <div class="datailcnt">
            <table class="datatab" border="1">
                <tr class="tabtop">
                    <th class="tablelf">症候</th>
                    <th>药品名称</th>
                    <th>理疗名称</th>
                </tr>
                <tbody v-show="symptomTypeLabel == '风险疾病'">
                    <tr v-for="v in Treatment.treatmentItemsRiskyDisease" :key="v" class="tabcnt">
                        <td class="tablelf">{{ v.symptomName }}</td>
                        <td class="tablerg">
                            <div class="yp" v-if="v.productfood?.productName !== null">
                                <img :src="v.productfood?.image_url" alt="">
                                {{ v.productfood?.product_name }}
                            </div>
                            <div v-if="v.productfood?.productName === null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v)">
                                        选择商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item
                                                    v-if="v.productsfood == null || v.productsfood == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsfood" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v)">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div v-if="v.productfood?.productName !== null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v)">
                                        修改商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item v-if="v.productsfood == null || v.productsfood == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsfood" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v)">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div class="sc" v-if="v.productfood?.productName !== null">
                                <el-button type="danger" plain @click="deletesp(v)">删除</el-button>
                            </div>
                            <div v-if="v.lock == true" class="lock">系统锁定无法修改</div>
                        </td>
                        <td class="tablerg">
                            <div class="yp" v-if="v.productll?.productName !== null">
                                <img :src="v.productll?.image_url" alt="">
                                {{ v.productll?.product_name }}
                            </div>
                            <div v-if="v.productll?.productName === null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v, 'll')">
                                        选择商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item
                                                    v-if="v.productsll == null || v.productsll == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsll" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v,)">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div v-if="v.productll?.productName !== null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v, 'll')">
                                        修改商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item
                                                    v-if="v.productsll == null || v.productsll == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsll" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v, 'll')">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div class="sc" v-if="v.productll?.productName !== null">
                                <el-button type="danger" plain @click="deletesp(v, 'll')">删除</el-button>
                            </div>
                            <div v-if="v.lock == true" class="lock">系统锁定无法修改</div>
                        </td>
                    </tr>
                </tbody>
                <tbody v-show="symptomTypeLabel == '脏腑辨证'">
                    <tr v-for="v in Treatment.treatmentItemsDirtyDialectic" :key="v" class="tabcnt">
                        <td class="tablelf">{{ v.symptomName }}</td>
                        <td class="tablerg">
                            <div class="yp" v-if="v.productfood?.productName !== null">
                                <img :src="v.productfood?.image_url" alt="">
                                {{ v.productfood?.product_name }}
                            </div>
                            <div v-if="v.productfood?.productName === null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v)">
                                        选择商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item
                                                    v-if="v.productsfood == null || v.productsfood == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsfood" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v)">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div v-if="v.productfood?.productName !== null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v)">
                                        修改商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item v-if="v.productsfood == null || v.productsfood == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsfood" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v)">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div class="sc" v-if="v.productfood?.productName !== null">
                                <el-button type="danger" plain @click="deletesp(v)">删除</el-button>
                            </div>
                            <div v-if="v.lock == true" class="lock">系统锁定无法修改</div>
                        </td>
                        <td class="tablerg">
                            <div class="yp" v-if="v.productll?.productName !== null">
                                <img :src="v.productll?.image_url" alt="">
                                {{ v.productll?.product_name }}
                            </div>
                            <div v-if="v.productll?.productName === null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v, 'll')">
                                        选择商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item
                                                    v-if="v.productsll == null || v.productsll == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsll" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v,)">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div v-if="v.productll?.productName !== null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v, 'll')">
                                        修改商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item
                                                    v-if="v.productsll == null || v.productsll == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsll" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v, 'll')">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div class="sc" v-if="v.productll?.productName !== null">
                                <el-button type="danger" plain @click="deletesp(v, 'll')">删除</el-button>
                            </div>
                            <div v-if="v.lock == true" class="lock">系统锁定无法修改</div>
                        </td>
                    </tr>
                </tbody>
                <tbody v-show="symptomTypeLabel == '体质辨证'">
                    <tr v-for="v in Treatment.treatmentItemsPhysicalDialectics" :key="v" class="tabcnt">
                        <td class="tablelf">{{ v.symptomName }}</td>
                        <td class="tablerg">
                            <div class="yp" v-if="v.productfood?.productName !== null">
                                <img :src="v.productfood?.image_url" alt="">
                                {{ v.productfood?.product_name }}
                            </div>
                            <div v-if="v.productfood?.productName === null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v)">
                                        选择商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item
                                                    v-if="v.productsfood == null || v.productsfood == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsfood" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v)">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div v-if="v.productfood?.productName !== null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v)">
                                        修改商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item v-if="v.productsfood == null || v.productsfood == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsfood" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v)">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div class="sc" v-if="v.productfood?.productName !== null">
                                <el-button type="danger" plain @click="deletesp(v)">删除</el-button>
                            </div>
                            <div v-if="v.lock == true" class="lock">系统锁定无法修改</div>
                        </td>
                        <td class="tablerg">
                            <div class="yp" v-if="v.productll?.productName !== null">
                                <img :src="v.productll?.image_url" alt="">
                                {{ v.productll?.product_name }}
                            </div>
                            <div v-if="v.productll?.productName === null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v, 'll')">
                                        选择商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item
                                                    v-if="v.productsll == null || v.productsll == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsll" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v,)">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div v-if="v.productll?.productName !== null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v, 'll')">
                                        修改商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item
                                                    v-if="v.productsll == null || v.productsll == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsll" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v, 'll')">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div class="sc" v-if="v.productll?.productName !== null">
                                <el-button type="danger" plain @click="deletesp(v, 'll')">删除</el-button>
                            </div>
                            <div v-if="v.lock == true" class="lock">系统锁定无法修改</div>
                        </td>
                    </tr>
                </tbody>
                <tbody v-show="symptomTypeLabel == '理疗'">
                    <tr v-for="v in Treatment.treatmentItemsPhysicalTherapy" :key="v" class="tabcnt">
                        <td class="tablelf">{{ v.symptomName }}</td>
                        <td class="tablerg">
                            <div class="yp" v-if="v.productfood?.productName !== null">
                                <img :src="v.productfood?.image_url" alt="">
                                {{ v.productfood?.product_name }}
                            </div>
                            <div v-if="v.productfood?.productName === null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v)">
                                        选择商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item
                                                    v-if="v.productsfood == null || v.productsfood == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsfood" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v)">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div v-if="v.productfood?.productName !== null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v)">
                                        修改商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item v-if="v.productsfood == null || v.productsfood == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsfood" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v)">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div class="sc" v-if="v.productfood?.productName !== null">
                                <el-button type="danger" plain @click="deletesp(v)">删除</el-button>
                            </div>
                            <div v-if="v.lock == true" class="lock">系统锁定无法修改</div>
                        </td>
                        <td class="tablerg">
                            <div class="yp" v-if="v.productll?.productName !== null">
                                <img :src="v.productll?.image_url" alt="">
                                {{ v.productll?.product_name }}
                            </div>
                            <div v-if="v.productll?.productName === null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v, 'll')">
                                        选择商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item
                                                    v-if="v.productsll == null || v.productsll == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsll" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v,)">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div v-if="v.productll?.productName !== null" class="xz">
                                <el-dropdown type="primary" trigger="click">
                                    <el-button plain size="medium" type="primary" class="xzspbut"
                                        @click="getlistsymptom(v.mapkey, v, 'll')">
                                        修改商品
                                        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"
                                            data-v-ea893728="" width="20px">
                                            <path fill="currentColor" d="m192 384 320 384 320-384z"></path>
                                        </svg>
                                    </el-button>
                                    <template #dropdown>
                                        <div v-loading="loading">
                                            <el-scrollbar style="height: 370px; width: 250px;">
                                                <el-dropdown-item
                                                    v-if="v.productsll == null || v.productsll == undefined"
                                                    style="cursor:default">
                                                    您还未添加任何商品
                                                </el-dropdown-item>
                                                <el-dropdown-item v-for="h in v.productsll" :key="h">
                                                    <div class="xialasp" @click="getproduct(h, v, 'll')">
                                                        <div class="xialaimg">
                                                            <img :src="h.image_url" alt="">
                                                        </div>
                                                        <el-tooltip :content="h.product_name" placement="bottom"
                                                            effect="light">
                                                            <div class="xialatxt">
                                                                {{ h.product_name }}
                                                            </div>
                                                        </el-tooltip>
                                                    </div>
                                                </el-dropdown-item>
                                            </el-scrollbar>
                                        </div>
                                    </template>
                                </el-dropdown>
                            </div>
                            <div class="sc" v-if="v.productll?.productName !== null">
                                <el-button type="danger" plain @click="deletesp(v, 'll')">删除</el-button>
                            </div>
                            <div v-if="v.lock == true" class="lock">系统锁定无法修改</div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
        <div class="datailbtm">
            <div class="datailbtmlf">
                <el-form-item label="方案名称：" required="true" style="text-align:left">
                    <el-input placeholder="请输入设置方案名称" v-model="Treatment.treatmentName"
                        style="width: 469px;height: 40px;" />
                </el-form-item>
            </div>
            <div class="datailbtmrg">
                <el-button @click="Nextstep" v-if="symptomTypeLabel != '理疗'"
                    style="width: 103px;height: 40px;">下一步</el-button>
                <el-button type="primary" @click="CreateTreatment" v-else class="fabc" :disabled="butno">保存</el-button>
            </div>
        </div>
    </div>
</template>
<script setup>
import { reactive, ref } from 'vue';
import { getTenantId, getSymptomKeyMap } from '@/utils/storage'
import { elmessage } from '@/utils/popup'
import {
    CreateTreatmentRequest,
    ListSymptomProductsRequest,
    UpdateTreatmentRequest
} from '@/api/api'
import { settopname } from '@/utils/sethometop'
import { ElMessage } from 'element-plus';

const symptomTypeLabel = ref('风险疾病')
const loading = ref(false)
const data = ref(getSymptomKeyMap())

const butno = ref(false)

//点击编辑进入页面传递到数据
let reviewtype = ref(JSON.parse(history.state?.keyword))
const reviewtycz = ref(JSON.parse(history.state?.keycz))
if (reviewtycz.value == false) {
    reviewtype.value.treatment.treatment_name = null
}
console.log(reviewtype.value);
const currentDiffProductSelected = ref([])
const Treatment = ref({
    treatmentId: null,
    treatmentName: reviewtype.value.treatment.treatment_name,
    reviewPass: null,
    reviewComment: null,
    isPublished: null,
    treatmentStatus: null,
    treatmentItemsRiskyDisease: [],
    treatmentItemsDirtyDialectic: [],
    treatmentItemsPhysicalTherapy: [],
    treatmentItemsPhysicalDialectics: [],
    createdTime: [],
    organizationId: null,
})
const lablexz = (rs) => {
    symptomTypeLabel.value = rs
}

if (reviewtype.value.treatment.treatment_id) {
    const getRiskyDisease = () => {
        const a = data.value.risk_disease_key_map
        const risk = reviewtype.value.treatment.treatment_items_risky_disease
        if (reviewtype.value.treatment.treatment_items_risky_disease != undefined) {
            for (let v in a) {
                for (let i = 0; i < risk.length; i++) {
                    if (risk[i].symptom == v) {
                        console.log(risk[i]);
                        let lllist = null
                        let foodlist = null
                        risk[i].products.forEach(eva => {
                            if (eva.product_type == 'PRODUCT_TYPE_SERVICE') {
                                lllist = eva
                                console.log(lllist);
                            } else {
                                foodlist = eva
                                console.log(foodlist);

                            }
                        })
                        if (lllist == null) {
                            lllist = {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            }
                        }
                        if(foodlist == null) {
                            foodlist = {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            }
                        }
                        Treatment.value.treatmentItemsRiskyDisease.push({
                            symptomName: a[v],
                            productfood: foodlist,
                            productll: lllist,
                            mapkey: v,
                            keyproduct: null,
                            productsll: null,
                            productsfood: null
                        })
                        console.log(Treatment.value.treatmentItemsPhysicalTherapy);
                        break
                    } else if (i == risk.length - 1) {
                        Treatment.value.treatmentItemsRiskyDisease.push({
                            symptomName: a[v],
                            productfood: {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            },
                            productll: {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            },
                            mapkey: v,
                            keyproduct: null,
                            productsll: null,
                            productsfood: null
                        })
                    }
                }
            }
        } else {
            for (let v in a) {
                Treatment.value.treatmentItemsRiskyDisease.push({
                    symptomName: a[v],
                    productfood: {
                        productType: null,
                        productName: null,
                        drugName: null,
                        isOtc: null,
                        approvedNumber: null,
                        drugValidityPeriod: null,
                        description: null,
                        symptomKeys: null,
                        remarks: null
                    },
                    productll: {
                        productType: null,
                        productName: null,
                        drugName: null,
                        isOtc: null,
                        approvedNumber: null,
                        drugValidityPeriod: null,
                        description: null,
                        symptomKeys: null,
                        remarks: null
                    },
                    mapkey: v,
                    keyproduct: null,
                    productsll: null,
                    productsfood: null
                })
            }
        }
        console.log(Treatment.value);

    }
    const getDirtyDialectic = () => {
        const a = data.value.dirty_dialectics_key_map
        delete a.Z0022
        const risk = reviewtype.value.treatment.treatment_items_dirty_dialectic
        console.log(reviewtype);
        console.log(risk);
        if (reviewtype.value.treatment.treatment_items_dirty_dialectic != undefined) {
            for (let v in a) {
                for (let i = 0; i < risk.length; i++) {
                    if (risk[i].symptom == v) {
                        let lllist = null
                        let foodlist = null
                        risk[i].products.forEach(eva => {
                            if (eva) {
                                if (eva.product_type == 'PRODUCT_TYPE_SERVICE') {
                                    lllist = eva
                                } else {
                                    foodlist = eva
                                }
                            }
                        })
                        if (lllist == null) {
                            lllist = {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            }
                        }
                        if(foodlist == null) {
                            foodlist = {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            }
                        }
                        Treatment.value.treatmentItemsDirtyDialectic.push({
                            symptomName: a[v],
                            productfood: foodlist,
                            productll: lllist,
                            mapkey: v,
                            keyproduct: null,
                            productsll: null,
                            productsfood: null
                        })
                        break
                    } else if (i == risk.length - 1) {
                        Treatment.value.treatmentItemsDirtyDialectic.push({
                            symptomName: a[v],
                            productfood: {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            },
                            productll: {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            },
                            mapkey: v,
                            keyproduct: null,
                            productsll: null,
                            productsfood: null
                        })
                    }
                }
            }
        } else {
            for (let v in a) {
                Treatment.value.treatmentItemsDirtyDialectic.push({
                    symptomName: a[v],
                    productfood: {
                        productType: null,
                        productName: null,
                        drugName: null,
                        isOtc: null,
                        approvedNumber: null,
                        drugValidityPeriod: null,
                        description: null,
                        symptomKeys: null,
                        remarks: null
                    },
                    productll: {
                        productType: null,
                        productName: null,
                        drugName: null,
                        isOtc: null,
                        approvedNumber: null,
                        drugValidityPeriod: null,
                        description: null,
                        symptomKeys: null,
                        remarks: null
                    },
                    mapkey: v,
                    keyproduct: null,
                    productsll: null,
                    productsfood: null
                })
            }
        }
    }
    const gethysicalTherapy = () => {
        const a = data.value.physical_therapy_key_map
        const risk = reviewtype.value.treatment.treatment_items_physical_therapy
        if (reviewtype.value.treatment.treatment_items_physical_therapy != undefined) {
            for (let v in a) {
                for (let i = 0; i < risk.length; i++) {
                    if (risk[i].symptom == v) {
                        let lllist = null
                        let foodlist = null
                        risk[i].products.forEach(eva => {
                            if (eva.product_type == 'PRODUCT_TYPE_SERVICE') {
                                lllist = eva
                                console.log(lllist);
                            } else {
                                foodlist = eva
                                console.log(foodlist);

                            }

                        })
                        if (lllist == null) {
                            lllist = {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            }
                        }
                        if(foodlist == null) {
                            foodlist = {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            }
                        }
                        Treatment.value.treatmentItemsPhysicalTherapy.push({
                            symptomName: a[v],
                            productfood: foodlist,
                            productll: lllist,
                            mapkey: v,
                            keyproduct: null,
                            productsll: null,
                            productsfood: null
                        })

                        break
                    } else if (i == risk.length - 1) {
                        Treatment.value.treatmentItemsPhysicalTherapy.push({
                            symptomName: a[v],
                            productfood: {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            },
                            productll: {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            },
                            mapkey: v,
                            keyproduct: null,
                            productsll: null,
                            productsfood: null
                        })
                    }
                }
            }
        } else {
            for (let v in a) {
                Treatment.value.treatmentItemsPhysicalTherapy.push({
                    symptomName: a[v],
                    productfood: {
                        productType: null,
                        productName: null,
                        drugName: null,
                        isOtc: null,
                        approvedNumber: null,
                        drugValidityPeriod: null,
                        description: null,
                        symptomKeys: null,
                        remarks: null
                    },
                    productll: {
                        productType: null,
                        productName: null,
                        drugName: null,
                        isOtc: null,
                        approvedNumber: null,
                        drugValidityPeriod: null,
                        description: null,
                        symptomKeys: null,
                        remarks: null
                    },
                    mapkey: v,
                    keyproduct: null,
                    productsll: null,
                    productsfood: null
                })
            }
        }

    }
    const gethysicaldialect = () => {
        const a = data.value.physique_dialectics_key_map
        const risk = reviewtype.value.treatment.treatment_items_physical_dialectics
        if (reviewtype.value.treatment.treatment_items_physical_dialectics != undefined) {
            for (let v in a) {
                for (let i = 0; i < risk.length; i++) {
                    if (risk[i].symptom == v) {
                        let lllist = null
                        let foodlist = null
                        risk[i].products.forEach(eva => {
                            if (eva) {
                                if (eva.product_type == 'PRODUCT_TYPE_SERVICE') {
                                    lllist = eva
                                } else {
                                    foodlist = eva
                                }
                            }
                        })
                        if (lllist == null) {
                            lllist = {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            }
                        }
                        if(foodlist == null) {
                            foodlist = {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            }
                        }
                        Treatment.value.treatmentItemsPhysicalDialectics.push({
                            symptomName: a[v],
                            productfood: foodlist,
                            productll: lllist,
                            mapkey: v,
                            keyproduct: null,
                            productsll: null,
                            productsfood: null
                        })
                        break
                    } else if (i == risk.length - 1) {
                        Treatment.value.treatmentItemsPhysicalDialectics.push({
                            symptomName: a[v],
                            productfood: {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            },
                            productll: {
                                productType: null,
                                productName: null,
                                drugName: null,
                                isOtc: null,
                                approvedNumber: null,
                                drugValidityPeriod: null,
                                description: null,
                                symptomKeys: null,
                                remarks: null
                            },
                            mapkey: v,
                            keyproduct: null,
                            productsll: null,
                            productsfood: null
                        })
                    }
                }
            }
        } else {
            for (let v in a) {
                Treatment.value.treatmentItemsPhysicalDialectics.push({
                    symptomName: a[v],
                    productfood: {
                        productType: null,
                        productName: null,
                        drugName: null,
                        isOtc: null,
                        approvedNumber: null,
                        drugValidityPeriod: null,
                        description: null,
                        symptomKeys: null,
                        remarks: null
                    },
                    productll: {
                        productType: null,
                        productName: null,
                        drugName: null,
                        isOtc: null,
                        approvedNumber: null,
                        drugValidityPeriod: null,
                        description: null,
                        symptomKeys: null,
                        remarks: null
                    },
                    mapkey: v,
                    keyproduct: null,
                    productsll: null,
                    productsfood: null
                })
            }
        }

    }
    getRiskyDisease()
    getDirtyDialectic()
    gethysicalTherapy()
    gethysicaldialect()
}
else {
    const getRiskyDisease = () => {
        const a = data.value.risk_disease_key_map
        for (let v in a) {
            Treatment.value.treatmentItemsRiskyDisease.push({
                symptomName: a[v],
                productfood: {
                    productType: null,
                    productName: null,
                    drugName: null,
                    isOtc: null,
                    approvedNumber: null,
                    drugValidityPeriod: null,
                    description: null,
                    symptomKeys: null,
                    remarks: null
                },
                productll: {
                    productType: null,
                    productName: null,
                    drugName: null,
                    isOtc: null,
                    approvedNumber: null,
                    drugValidityPeriod: null,
                    description: null,
                    symptomKeys: null,
                    remarks: null
                },
                mapkey: v,
                keyproduct: null,
                productsll: null,
                productsfood: null
            })
        }
    }
    const getDirtyDialectic = () => {
        const a = data.value.dirty_dialectics_key_map
        delete a.Z0022
        for (let v in a) {
            Treatment.value.treatmentItemsDirtyDialectic.push({
                symptomName: a[v],
                productfood: {
                    productType: null,
                    productName: null,
                    drugName: null,
                    isOtc: null,
                    approvedNumber: null,
                    drugValidityPeriod: null,
                    description: null,
                    symptomKeys: null,
                    remarks: null
                },
                productll: {
                    productType: null,
                    productName: null,
                    drugName: null,
                    isOtc: null,
                    approvedNumber: null,
                    drugValidityPeriod: null,
                    description: null,
                    symptomKeys: null,
                    remarks: null
                },
                mapkey: v,
                keyproduct: null,
                productsll: null,
                productsfood: null
            })
        }
    }
    const gethysicalTherapy = () => {
        const a = data.value.physical_therapy_key_map
        for (let v in a) {
            Treatment.value.treatmentItemsPhysicalTherapy.push({
                symptomName: a[v],
                productfood: {
                    productType: null,
                    productName: null,
                    drugName: null,
                    isOtc: null,
                    approvedNumber: null,
                    drugValidityPeriod: null,
                    description: null,
                    symptomKeys: null,
                    remarks: null
                },
                productll: {
                    productType: null,
                    productName: null,
                    drugName: null,
                    isOtc: null,
                    approvedNumber: null,
                    drugValidityPeriod: null,
                    description: null,
                    symptomKeys: null,
                    remarks: null
                },
                mapkey: v,
                keyproduct: null,
                productsll: null,
                productsfood: null
            })
        }
        console.log(Treatment.value);
    }
    const gethysicaldialect = () => {
        const a = data.value.physique_dialectics_key_map
        console.log(a);
        for (let v in a) {
            Treatment.value.treatmentItemsPhysicalDialectics.push({
                symptomName: a[v],
                productfood: {
                    productType: null,
                    productName: null,
                    drugName: null,
                    isOtc: null,
                    approvedNumber: null,
                    drugValidityPeriod: null,
                    description: null,
                    symptomKeys: null,
                    remarks: null
                },
                productll: {
                    productType: null,
                    productName: null,
                    drugName: null,
                    isOtc: null,
                    approvedNumber: null,
                    drugValidityPeriod: null,
                    description: null,
                    symptomKeys: null,
                    remarks: null
                },
                mapkey: v,
                keyproduct: null,
                productsll: null,
                productsfood: null
            })
        }
        console.log(Treatment.value);
    }
    getRiskyDisease()
    getDirtyDialectic()
    gethysicalTherapy()
    gethysicaldialect()
}

const getlistsymptom = async (key, productss, type) => {
    console.log('获取商品列表');
    
    const rs = {
        organizationId: getTenantId(),
        symptomKey: key
    }
    if (type == 'll') {
        if (productss.productsll == null || productss.productsll == undefined) {
            loading.value = true
            const data = await ListSymptomProductsRequest(rs)
            loading.value = false
            console.log(data);
            const datatype = []
            data.data.products.forEach(v => {
                if (v.product_type == 'PRODUCT_TYPE_SERVICE') {
                    datatype.push(v)
                }
            })
            productss.productsll = datatype
            console.log(productss);
        }
    } else {
        console.log(123);
        
        if (productss.productsfood == null || productss.productsfood == undefined) {
            loading.value = true
            const data = await ListSymptomProductsRequest(rs)
            loading.value = false
            console.log(data);
            const datatype = []
            data.data.products.forEach(v => {
                if (v.product_type != 'PRODUCT_TYPE_SERVICE') {
                    datatype.push(v)
                }
            })
            productss.productsfood = datatype
            console.log(productss);
            console.log(Treatment.value);
        }
    }

}

const getproduct = (rs, product) => {
    console.log(rs);
    console.log(product);

    if (currentDiffProductSelected.value.length < 10) {
        if (rs.product_type == 'PRODUCT_TYPE_SERVICE') {
            if (product.productll.productName == '') {
                product.productll = rs
                let ab = false
                currentDiffProductSelected.value.forEach(v => {
                    if (v == rs.product_id) {
                        ab = true
                    } else {
                        ab = false
                    }
                })
                if (ab == false) {
                    currentDiffProductSelected.value.push(rs.product_id)
                }
            } else {
                product.productll = rs
            }
        } else {
            if (product.productfood.productName == '') {
                product.productfood = rs
                let ab = false
                currentDiffProductSelected.value.forEach(v => {
                    if (v == rs.product_id) {
                        ab = true
                    } else {
                        ab = false
                    }
                })
                if (ab == false) {
                    currentDiffProductSelected.value.push(rs.product_id)
                }
            } else {
                product.productfood = rs
            }
        }

    } else {
        elmessage('选择商品失败，每个方案限制使用10个商品')
    }
    console.log(Treatment.value.treatmentItemsRiskyDisease);

    console.log(Treatment.value);
}
const deletesp = (rs, type) => {
    currentDiffProductSelected.value = currentDiffProductSelected.value.filter(item => item !== rs.product.product_id)
    console.log(currentDiffProductSelected.value, rs);
    if (type == 'll') {
        rs.productll = {
            productType: null,
            productName: null,
            drugName: null,
            isOtc: null,
            approvedNumber: null,
            drugValidityPeriod: null,
            description: null,
            symptomKeys: null,
            remarks: null
        }
    } else {
        rs.productfood = {
            productType: null,
            productName: null,
            drugName: null,
            isOtc: null,
            approvedNumber: null,
            drugValidityPeriod: null,
            description: null,
            symptomKeys: null,
            remarks: null
        }
    }

}

const CreateTreatment = () => {
    butno.value = true
    if (Treatment.value.treatmentName == '' || Treatment.value.treatmentName == null || Treatment.value.treatmentName == undefined) {
        ElMessage({
            showClose: true,
            message: '请输入方案名称',
            type: 'error',
        })
        butno.value = false
    } else {
        const treatmentItemsRiskyDisease = []
        const treatmentItemsDirtyDialectic = []
        const treatmentItemsPhysicalTherapy = []
        const treatmentItemsPhysicalDialectics = []
        Treatment.value.treatmentItemsRiskyDisease.forEach(v => {
            const datalist = []
            let databale = false
            if (v.productfood.product_id != null) {
                datalist.push(v.productfood)
                databale = true
            }
            if (v.productll.product_id != null) {
                datalist.push(v.productll)
                databale = true
            }
            if (databale) {
                delete v.productfood.symptom_keys;
                delete v.productll.symptom_keys;
                treatmentItemsRiskyDisease.push(
                    {
                        symptom: v.mapkey,
                        products: datalist
                    }
                )
            }

        })
        Treatment.value.treatmentItemsDirtyDialectic.forEach(v => {
            const datalist = []
            let databale = false
            if (v.productfood.product_id != null) {
                datalist.push(v.productfood)
                databale = true
            }
            if (v.productll.product_id != null) {
                datalist.push(v.productll)
                databale = true
            }
            if (databale) {
                delete v.productfood.symptom_keys;
                delete v.productll.symptom_keys;
                treatmentItemsDirtyDialectic.push(
                    {
                        symptom: v.mapkey,
                        products: datalist
                    }
                )
            }
        })
        Treatment.value.treatmentItemsPhysicalTherapy.forEach(v => {
            const datalist = []
            let databale = false
            if (v.productfood.product_id != null) {
                datalist.push(v.productfood)
                databale = true
            }
            if (v.productll.product_id != null) {
                datalist.push(v.productll)
                databale = true
            }
            if (databale) {
                console.log(v);
                delete v.productfood.symptom_keys;
                delete v.productll.symptom_keys;
                treatmentItemsPhysicalTherapy.push(
                    {
                        symptom: v.mapkey,
                        products: datalist
                    }
                )
            }
        })
        Treatment.value.treatmentItemsPhysicalDialectics.forEach(v => {
            const datalist = []
            let databale = false
            if (v.productfood.product_id != null) {
                datalist.push(v.productfood)
                databale = true
            }
            if (v.productll.product_id != null) {
                datalist.push(v.productll)
                databale = true
            }
            if (databale) {
                delete v.productfood.symptom_keys;
                delete v.productll.symptom_keys;
                treatmentItemsPhysicalDialectics.push(
                    {
                        symptom: v.mapkey,
                        products: datalist
                    }
                )
            }
        })
        const rs = {
            organizationId: getTenantId(),
            treatment: {
                treatmentName: Treatment.value.treatmentName,
                treatmentItemsRiskyDisease: treatmentItemsRiskyDisease,
                treatmentItemsDirtyDialectic: treatmentItemsDirtyDialectic,
                treatmentItemsPhysicalTherapy: treatmentItemsPhysicalTherapy,
                treatmentItemsPhysicalDialectics: treatmentItemsPhysicalDialectics
            }
        }

        if (reviewtycz.value) {
            const add = {
                treatment: {
                    organizationId: getTenantId(),
                    treatmentId: reviewtype.value.treatment.treatment_id,
                    treatmentName: Treatment.value.treatmentName,
                    treatmentItemsRiskyDisease: treatmentItemsRiskyDisease,
                    treatmentItemsDirtyDialectic: treatmentItemsDirtyDialectic,
                    treatmentItemsPhysicalTherapy: treatmentItemsPhysicalTherapy,
                    treatmentItemsPhysicalDialectics: treatmentItemsPhysicalDialectics

                }
            }
            if (add.treatment.treatmentItemsRiskyDisease.length == 0 && add.treatment.treatmentItemsDirtyDialectic.length == 0 && add.treatment.treatmentItemsPhysicalTherapy.length == 0 && add.treatment.treatmentItemsPhysicalDialectics.length == 0) {
                ElMessage({
                    showClose: true,
                    message: '请上传商品',
                    type: 'error',
                })
                butno.value = false

            } else {
                UpdateTreatmentRequest(add)
                    .then(res => {
                        console.log(res);
                        if (res.status === 200) {
                            ElMessage({
                                message: '方案更新成功',
                                type: 'success',
                            })
                            window.history.go(-1)
                            butno.value = false

                        } else {
                            ElMessage({
                                showClose: true,
                                message: res.data.detail,
                                type: 'error',
                            })
                            butno.value = false

                        }

                    })
            }

        } else {
            console.log(rs.treatment);
            if (rs.treatment.treatmentItemsDirtyDialectic.length == 0 && rs.treatment.treatmentItemsPhysicalTherapy.length == 0 && rs.treatment.treatmentItemsRiskyDisease.length == 0 && rs.treatment.treatmentItemsPhysicalDialectics.length == 0) {
                ElMessage({
                    showClose: true,
                    message: '请上传商品',
                    type: 'error',
                })
                butno.value = false

            } else {
                CreateTreatmentRequest(rs)
                    .then(res => {
                        console.log(res);
                        if (res.status === 200) {
                            ElMessage({
                                message: '方案创建成功成功',
                                type: 'success',
                            })
                            window.history.go(-1)
                            butno.value = false

                        } else {
                            ElMessage({
                                showClose: true,
                                message: res.data.detail,
                                type: 'error',
                            })
                            butno.value = false

                        }
                    })
            }

        }
    }

}
const topname = () => {
    settopname([{ name: '商品推荐方案', url: '/recommendedProposal' }, { name: '创建方案', url: '' }])
}
const Nextstep = () => {
    console.log(1);
    console.log(symptomTypeLabel.value);
    if (symptomTypeLabel.value == '风险疾病') {
        symptomTypeLabel.value = '脏腑辨证'
    } else if (symptomTypeLabel.value == '脏腑辨证') {
        console.log(1);
        symptomTypeLabel.value = '体质辨证'
    } else if (symptomTypeLabel.value == '体质辨证') {
        console.log(1);
        symptomTypeLabel.value = '理疗'
    }
}
topname()
</script>
<style lang="scss">
.block-col-2 .demonstration {
    display: block;
    color: var(--el-text-color-secondary);
    font-size: 14px;
    margin-bottom: 20px;
}

.detail {
    background-color: #fff;

    .dataitop {
        display: flex;
        justify-content: space-between;
        padding: 20px 20px;

        .detatoplf {
            display: flex;
            font-size: 14px;
            color: #333333;

            .fx,
            .zf,
            .ll,
            .tz {
                width: 144px;
                height: 40px;
                background-size: auto 100%;
                background-repeat: no-repeat;
                line-height: 40px;
                text-align: center;
                cursor: pointer;
            }

            .fx {
                background-image: url(../../assets/image/step_1.png);
            }

            .fxxz {
                background-image: url(../../assets/image/step_1_selected.png);
                color: #fff;
            }

            .zf {
                background-image: url(../../assets/image/step_2.png);
                position: relative;
                left: -15px;
            }

            .zfxz {
                background-image: url(../../assets/image/step_2_selected.png);
                color: #fff;
            }

            .tz {
                background-image: url(../../assets/image/step_2.png);
                position: relative;
                left: -30px;
            }

            .tzxz {
                background-image: url(../../assets/image/step_2_selected.png);
                color: #fff;
            }

            .ll {
                background-image: url(../../assets/image/step_3.png);
                position: relative;
                left: -50px;
            }

            .llxz {
                background-image: url(../../assets/image/step_3_selected.png);
                color: #fff;
            }
        }

        .detatoprg {
            font-size: 14px;
            color: #666666;
        }
    }

    .datailcnt {
        padding: 20px 20px;
    }

    .datailbtm {
        padding: 20px;
        display: flex;
        justify-content: space-between;
        align-items: center;

        .el-form-item__label {
            width: 106px;
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
            height: 100px;
            min-width: 550px;

            .yp {
                position: absolute;
                left: 40px;
                top: 6px;
                font-weight: 700;
                font-style: normal;
                color: #333333;
                font-size: 18px;
                display: flex;
                align-items: center;

                img {
                    width: 88px;
                    height: 88px;
                    margin-right: 20px;
                }
            }

            .xz {
                width: 92px;
                position: absolute;
                right: 120px;
                top: 35px;
            }

            .sc {
                width: 92px;
                position: absolute;
                right: 15px;
                top: 35px;
            }
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

.el-dropdown-menu {
    z-index: 3000 !important;
}

.xzspbut:hover {
    background-color: #ecf5ff;
    color: #3ca1f8;
}

.xialasp {
    width: 100%;
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: space-around;
    cursor: pointer;

    .xialaimg {
        padding-top: 13px;

        img {
            width: 48px;
            height: 48px;
        }
    }

    .xialatxt {
        font-size: 18px;
        font-weight: 650;
        width: 130px;
        text-overflow: ellipsis;
        overflow: hidden;
    }
}

.datailbtmlf {
    width: 469px;
}

.fabc {
    width: 103px;
    height: 40px;
    background-color: #006be5;
}
</style>