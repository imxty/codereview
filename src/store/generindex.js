import { createStore } from 'vuex';
import { GetOrganizationTenantSummaryRequest } from '@/api/api'
import { getTenantId } from '@/utils/storage'

//账单列表数据
export default createStore({
    state: {
        Organization: {
            total_tenant_amounts: 0,
            using_teant_amounts: 0,
            today_added_tenant_count: 0,
            tenants_expired_amounts: 0,
            customer_overdue_count: 0
        }
    },
    mutations: {
        async getlist(state) {
            const data = await GetOrganizationTenantSummaryRequest({ organizationId: getTenantId() })
               state.Organization.total_tenant_amounts = data.data.total_tenant_amounts === undefined ? '0' : data.data.total_tenant_amounts

               state.Organization.using_teant_amounts = data.data.using_teant_amounts === undefined ? '0' : data.data.using_teant_amounts

               state. Organization.today_added_tenant_count = data.data.today_added_tenant_count === undefined ? '0' : data.data.today_added_tenant_count

               state. Organization.tenants_expired_amounts = data.data.tenants_expired_amounts === undefined ? '0' : data.data.tenants_expired_amounts

               state. Organization.customer_overdue_count = data.data.customer_overdue_count === undefined ? '0' : data.data.customer_overdue_count
        },

    },
    actions: {

    }
})