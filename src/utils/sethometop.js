import { useStore } from 'vuex'
export function settopname (value){
    const $store = useStore();
    $store.commit('setHometop', value)
    
}