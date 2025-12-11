import { createRouter, createWebHistory, createWebHashHistory } from 'vue-router';

const routes = [
    {
        path: '/',
        name: 'home',
        component: () => import('@/views/HomeView.vue'),
        children: [
            {
                path: '/',
                redirect: '/general'
            },
            {
                path: '/general',
                name: 'general',
                component: () => import('@/views/admin/GeneralizeView.vue'),
                meta: {
                    title: '概况'
                }, 
                children: [
                    {
                        path: '',
                        name: 'Generalize',
                        component: () => import('@/views/general/GeneralizeView.vue'),
                    },
                    {
                        path: 'MeasurementData',
                        name: 'MeasurementData',
                        component: () => import('@/views/general/MeasurementDataView.vue'),
                    },
                    {
                        path: 'CustomerData',
                        name: 'CustomerData',
                        component: () => import('@/views/general/CustomerDataView.vue'),
                    },
                    {
                        path: 'TenantCompare',
                        name: 'TenantCompare',
                        component: () => import('@/views/general/TenantCompareView.vue'),
                    },
                    
                ]
            },
            {
                path: '/merchantManagement',
                name: 'merchantManagement',
                component: () => import('@/views/admin/MerchantView.vue'),
                meta: {
                    title: '商户管理'
                },
                children: [
                    {
                        path: '',
                        name: 'commmange',
                        component: () => import('@/views/Commercialowner/CommermanageView.vue'),
                    },
                    {
                        path: 'prepayment',
                        name: 'prepayment',
                        component: () => import('@/views/Commercialowner/PrepaymentView.vue'),
                    },
                    {
                        path: 'DistributionScheme',
                        name: 'DistributionScheme',
                        component: () => import('@/views/Commercialowner/DistributionSchemeView.vue'),
                    }
                ]
            },
            {
                path: '/products',
                name: 'products',
                component: () => import('@/views/admin/MemberView.vue'),
                meta: {
                    title: '商品管理'
                },
                children: [
                    {
                        path: '',
                        name: 'commduse',
                        component: () => import('@/views/commodmanag/commoduseView.vue'),
                        meta: {
                            title: '商品管理'
                        },
                    },
                    {
                        path: 'addProduct',
                        name: 'addProduct',
                        component: () => import('@/views/commodmanag/addgoodsView.vue'),
                        meta: {
                            title: '添加商品'
                        },
                    },
                    {
                        path: 'edlgoods',
                        name: 'edlgoods',
                        component: () => import('@/views/commodmanag/editorialgoodsView.vue'),
                        meta: {
                            title: '编辑商品'
                        },
                    },
                    {
                        path: 'prodetails',
                        name: 'prodetails',
                        component: () => import('@/views/commodmanag/productdetailsView.vue'),
                        meta: {
                            title: '编辑商品'
                        },
                    },
                    {
                        path: 'dataanalysis',
                        name: 'dataanalysis',
                        component: () => import('@/views/commodmanag/DataanalysisView.vue'),
                        meta: {
                            title: '编辑商品'
                        },
                    }

                ]
            },
            {
                path: '/recommendedProposal',
                name: 'recommendedProposal',
                component: () => import('@/views/admin/MembershipView.vue'),
                meta: {
                    title: '商品推荐方案'
                },

                children: [
                    {
                        path: '',
                        name: 'Recomscheme',
                        component: () => import('@/views/productoption/RecomschemeView.vue'),

                    },

                    {
                        path: 'detailsView',
                        name: 'detailsView',
                        component: () => import('@/views/productoption/detailsView.vue'),
                        meta: {
                            title: '详情'
                        },
                    },
                    {
                        path: 'treatmentEdit',
                        name: 'treatmentEdit',
                        component: () => import('@/views/productoption/createschemeView.vue'),
                        meta: {
                            title: '添加方案'
                        },
                    },
                ]
            },
            {
                path: '/billList',
                name: 'billList',
                component: () => import('@/views/admin/CheckView.vue'),
                meta: {
                    title: '账单管理'
                },
                children: [
                    {
                        path: '',
                        name: 'checkdetai',
                        component: () => import('@/views/check/checkdetailView.vue'),
                        meta: {
                            title: '账单详情'
                        },
                    },
                    {
                        path: 'check',
                        name: 'check',
                        component: () => import('@/views/check/checkmanageView.vue'),
                        meta: {
                            title: '账单详情'
                        },
                    }
                ]
            },
        ]

    },
    {
        path: '/login',
        name: 'login',
        component: () => import('@/views/LoginView.vue')
    }, {
        path: '/account',
        name: 'account',
        component: () => import('@/views/FirstentryView.vue'),
        children: [
            {
                path: 'createTenant',
                name: 'createTenant',
                component: () => import('@/views/firstchildren/foundView.vue'),
            },
            {
                path: 'findPassword',
                name: 'findPassword',
                component: () => import('@/views/firstchildren/forgotpasswordView.vue'),

            },
            {
                path: 'accountInformation',
                name: 'accountInformation',
                component: () => import('@/views/firstchildren/accountmentView.vue'),

            }
        ]

    },
    {
        path: '/lnvite',
        name: 'lnvite',
        component: () => import('@/views/LnviteView.vue')
    },
    {
        path: '/edittenant',
        name: 'edittenant',
        component: () => import('@/views/edittenant.vue')
    },
    {
        path: '/userLicense',
        name: 'userLicense',
        component: () => import('@/views/userLicense.vue')
    },
    {
        path: '/:pathMatch(.*)',
        name: '404',
        component: () => import('@/views/NotFon.vue')
    }
]

//将默认路由表添加到路由中
const router = createRouter({
    history: createWebHashHistory(),
    routes
})
export default router