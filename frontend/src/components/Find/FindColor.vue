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
import {ref, watch} from 'vue'

const value = ref('红色')
const options = [
  {
    value: '黄色',
    label: '黄色',
  },
  {
    value: '蓝色',
    label: '蓝色',
  },
  {
    value: '红色',
    label: '红色',
  },
  {
    value: '绿色',
    label: '绿色',
  },
  {
    value: '紫色',
    label: '紫色',
  },
]

function GetValue() {
  return GetColorValue(value.value)
}

function GetColorValue(w) {
  switch (w) {
    case "黄色":
      return "#ffe100"
    case "蓝色":
      return "#005ffd"
    case "红色":
      return "#fa0000"
    case "绿色":
      return "#029802"
    case "紫色":
      return "#6905fa"
    default:
      return "#ffe100"
  }
}

watch(value, (newValue, oldValue) => {
  window.vm.Find.SetFindColor(GetColorValue(newValue))
});

defineExpose({
  GetValue
});
</script>
