<template>
  <el-select v-model="value" placeholder="Select" size="small">
    <el-option
        v-for="item in options"
        :key="item.value"
        :label="item.label"
        :value="item.value"
        :disabled="item.disabled"
    />
  </el-select>
</template>

<script lang="ts" setup>
import {ref} from 'vue'
import {Config_Find_Range_Options_All} from "../../config/Config.js";

const value = ref('在全部范围寻找')
const options = ref(Config_Find_Range_Options_All)

function GetValue() {
  return value.value
}

function SetOptions(val) {
  return options.value = val
}

defineExpose({
  GetValue, SetOptions
});
</script>
