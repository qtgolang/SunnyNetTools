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

const value = ref('UTF8')
const options = [
  {
    value: 'UTF8',
    label: '字符串(UTF8)',
  },
  {
    value: 'GBK',
    label: '字符串(GBK)',
  },
  {
    value: 'Hex',
    label: '十六进制',
  },
  {
    value: 'Base64',
    label: 'Base64',
  },
  {
    value: 'pb',
    label: 'ProtoBuf',
  },
  {
    value: '整数4',
    label: '4字节整数',
  },
  {
    value: '整数8',
    label: '8字节整数',
  },
  {
    value: '浮点数4',
    label: '4字节浮点数',
  },
  {
    value: '浮点数8',
    label: '8字节浮点数',
  },
]

function GetType() {
  return value.value
}

watch(value, (newValue, oldValue) => {
  window.vm.Find.SetFindType(newValue)
});

defineExpose({
  GetType
});
</script>
