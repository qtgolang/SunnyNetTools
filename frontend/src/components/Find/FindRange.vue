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

const value = ref('全部')
const options = [
  {
    value: '全部',
    label: '全部',
  },
  {
    value: 'HTTP请求',
    label: 'HTTP/HTTPS 请求',
  },
  {
    value: 'HTTP响应',
    label: 'HTTP/HTTPS 响应',
  },
  {
    value: 'socketSend',
    label: 'TCP/UDP/WebSocket 发送',
  },
  {
    value: 'socketRec',
    label: 'TCP/UDP/WebSocket 接收',
  },

  {
    value: 'socketAll',
    label: 'TCP/UDP/WebSocket 发送/接收',
  },
]

function GetValue() {
  return value.value
}

defineExpose({
  GetValue
});
</script>
