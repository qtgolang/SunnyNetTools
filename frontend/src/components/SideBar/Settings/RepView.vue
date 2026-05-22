<script>
export default {
  props: ['params'],
  data() {
    return {
      value: this.params.value,
      options: [
        {
          value: '全部数据',
          label: '全部数据',
        },
        {
          value: 'HTTP请求-全部',
          label: 'HTTP请求-全部',
        },
        {
          value: 'HTTP请求-URL',
          label: 'HTTP请求-URL',
        },
        {
          value: 'HTTP请求-协议头',
          label: 'HTTP请求-协议头',
        },
        {
          value: 'HTTP请求-提交数据',
          label: 'HTTP请求-提交数据',
        },
      ]
    }
  },
  methods: {},
  mounted() {

  }
};
</script>
<template>
  <div style="cursor: pointer; display: flex; align-items: center;margin-left:-7px;margin-right:-7px;height: 100%"
       ref="el">
    <el-select v-model="value" placeholder="Select" size="small" style="width: 100%">
      <el-option
          v-for="item in options"
          :key="item.value"
          :label="item.label"
          :value="item.value"
          :disabled="item.disabled"

      />
    </el-select>
  </div>
</template>>