import { createStore } from 'vuex';
import { ListTenantsRequest, GetTenantStencilRequest } from '@/api/api'
import { getTenantId } from '@/utils/storage'

//账单列表数据
export default createStore({
    state: {
        commoddata: [],
        TenantSt: false,
        cities:[]
    },
    mutations: {
        CATEGORYLIST(state,rs){
            state.commoddata = rs
        }
    },
    actions: {
        async getlist({commit},rs) {
            const tenste = await GetTenantStencilRequest({ organizationId: getTenantId() })
            const data = await ListTenantsRequest(rs)
            commit('CATEGORYLIST',data.data.tenants)
        },

    }
})