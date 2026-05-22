<template>
  <el-select v-model="value" placeholder="Select">
    <el-option
        v-for="item in options"
        :key="item.value"
        :label="item.label"
        :value="item.value"
    />
  </el-select>
</template>

<script lang="ts" setup>
import {ref, watch} from 'vue'

import {Config_Find, Config_Find_Type_Options} from "../../config/Config.js";


const value = ref('auto')
const options = ref(Config_Find_Type_Options)

function GetType() {
  return value.value
}

function SetType(val) {
  return value.value = val
}

function SetOptions(val) {
  return options.value = val
}

watch(value, (newVal, oldVal) => {
  Config_Find.value.Type = newVal
},);
defineExpose({
  GetType, SetOptions, SetType
});
</script>
