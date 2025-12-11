<template>
   <div>
      <summarize></summarize>
      <datacomparison></datacomparison>
   </div>
</template>

<script setup>
import summarize from '../../components/common/summarize.vue';
import datacomparison from '../../components/common/datacomparison.vue';
import {settopname} from  '../../utils/sethometop'
import { useStore } from 'vuex'

//设置vuex对象
const $store = useStore();
console.log(1);

const topname = ()=>{
   settopname([{name:'概况'}])
   $store.commit('setactiveMenuKey', 0)
}
topname()
</script>