<script>
import {ElNotification} from "element-plus";
import {ReplaceBodyUpdate} from "../../../../bindings/changeme/Service/appmain";
import {Config_agGrid_API, registerThisObject} from "../../config/Config";

export default {
  props: ['params'],
  watch: {
    "value"(m, l) {
      if (m === this.value2) {
        return
      }
      if (this.Params.data['替换类型'] === "拦截请求") {
        this.Params.data['新数据'] = "拦截请求[此项不用填写]"
      } else if (this.Params.data['新数据'] === "拦截请求[此项不用填写]") {
        this.Params.data['新数据'] = "[双击修改-新数据]"
      }
      this.Params.data['查找范围'] = m
      const node = this.Params.data;
      ReplaceBodyUpdate(parseInt(node.id), node["替换类型"], node["查找范围"], node["旧数据"], node["新数据"], node["注释"]).then(res => {
        if (res) {
          node["状态"] = "已生效"
          if (this.isIntervalGreaterThan100ms(new Date().getTime(), this.t)) {
            ElNotification({
              showClose: true,
              message: '拦截/替换规则已更新',
              type: 'success',
              position: 'bottom-right',
            })
          }
          this.t = new Date().getTime()
        } else {
          node["状态"] = "未生效"
          if (this.isIntervalGreaterThan100ms(new Date().getTime(), this.t)) {
            ElNotification({
              position: 'bottom-right',
              showClose: true,
              message: '拦截/替换规则更新失败,请检查数据格式是否正确?',
              type: 'warning',
            })
          }
          this.t = new Date().getTime()
        }
        Config_agGrid_API.value.applyTransaction({update: [node]});
      })
    },
  },
  data() {
    return {
      Params: null,
      value: "任意",
      value2: "任意",
      items: [],
      t: null,
      options: [
        {
          label: 'HTTP请求/响应',
          value: 'HTTP请求/响应',
          tooltip: '如果 在HTTP请求/响应 任意数据中,匹配上 旧数据则拦截此请求',
          isBreak: true
        },
        {
          label: 'HTTP请求-全部',
          value: 'HTTP请求-全部',
          tooltip: '如果 在HTTP请求-任意数据中,匹配上 旧数据则拦截此请求',
          isBreak: true
        },
        {
          label: 'HTTP响应-全部',
          value: 'HTTP响应-全部',
          tooltip: '如果 在HTTP响应-任意数据中,匹配上 旧数据则拦截此请求',
          isBreak: true
        },
        {
          label: 'HTTP请求-URL',
          value: 'HTTP请求-URL',
          tooltip: '如果 在HTTP请求-的URL中,匹配上 旧数据则拦截此请求',
          isBreak: true
        },
        {
          label: 'HTTP请求-协议头',
          value: 'HTTP请求-协议头',
          tooltip: '如果 在HTTP请求-的协议头中,匹配上 旧数据则拦截此请求',
          isBreak: true
        },
        {
          label: 'HTTP请求-提交数据',
          value: 'HTTP请求-提交数据',
          tooltip: '如果 在HTTP请求-的POST提交数据中,匹配上 旧数据则拦截此请求',
          isBreak: true
        },
        {
          label: 'HTTP响应-协议头',
          value: 'HTTP响应-协议头',
          tooltip: '如果 在HTTP响应-的协议头中,匹配上 旧数据则拦截此请求',
          isBreak: true
        },
        {
          label: 'HTTP响应-响应数据',
          value: 'HTTP响应-响应数据',
          tooltip: '如果 在HTTP响应-的数据中,匹配上 旧数据则拦截此请求',
          isBreak: true
        },
        {label: '任意', value: '任意', tooltip: '在任何类型的数据中替换', isBreak: false},
        {label: 'Socket-任意', value: 'Socket-任意', tooltip: 'TCP/UDP/Websocket-发送/接收 数据中替换', isBreak: false},
        {label: 'TCP-全部', value: 'TCP-全部', tooltip: 'TCP-发送/接收 数据中替换', isBreak: false},
        {label: 'UDP-全部', value: 'UDP-全部', tooltip: 'UDP-发送/接收 数据中替换', isBreak: false},
        {label: 'Websocket-全部', value: 'Websocket-全部', tooltip: 'Websocket-发送/接收 数据中替换', isBreak: false},
        {label: 'HTTP请求-全部', value: 'HTTP请求-全部', tooltip: 'HTTP/S-请求/响应 任意数据中替换', isBreak: false},
        {label: 'HTTP请求-URL', value: 'HTTP请求-URL', tooltip: 'HTTP/S请求-URL中替换', isBreak: false},
        {label: 'HTTP请求-协议头', value: 'HTTP请求-协议头', tooltip: 'HTTP/S请求-协议头中替换', isBreak: false},
        {label: 'HTTP请求-提交数据', value: 'HTTP请求-提交数据', tooltip: 'HTTP/S请求-提交数据中替换', isBreak: false},
        {label: 'HTTP响应-全部', value: 'HTTP响应-全部', tooltip: 'HTTP/S响应-全部数据中替换', isBreak: false},
        {label: 'HTTP响应-协议头', value: 'HTTP响应-协议头', tooltip: 'HTTP/S响应-协议头中替换', isBreak: false},
        {label: 'HTTP响应-响应数据', value: 'HTTP响应-响应数据', tooltip: 'HTTP/S响应-响应数据中替换', isBreak: false},
        {label: 'TCP-发送', value: 'TCP-发送', tooltip: 'TCP-发送数据中替换', isBreak: false},
        {label: 'TCP-接收', value: 'TCP-接收', tooltip: 'TCP-接收数据中替换', isBreak: false},
        {label: 'UDP-发送', value: 'UDP-发送', tooltip: 'UDP-发送数据中替换', isBreak: false},
        {label: 'UDP-接收', value: 'UDP-接收', tooltip: 'UDP-接收数据中替换', isBreak: false},
        {label: 'Websocket-发送', value: 'Websocket-发送', tooltip: 'Websocket-发送数据中替换', isBreak: false},
        {label: 'Websocket-接收', value: 'Websocket-接收', tooltip: 'Websocket-接收数据中替换', isBreak: false},
      ]
    }
  },
  created() {
    this.Params = this.params;
    this.Params = this.params;
    this.value2 = this.Params.data['查找范围']
    this.value = this.Params.data['查找范围']
  },
  computed: {
    isBreak() {
      return this.Params.data['替换类型'] === "拦截请求"
    },
    BreakItems() {
      const array = [];
      for (const item of this.options) {
        if (item.isBreak === true) {
          array.push(item)
        }
      }
      return array
    },
    NoBreakItems() {
      const array = [];
      for (const item of this.options) {
        if (item.isBreak === false) {
          array.push(item)
        }
      }
      return array
    }
  },
  mounted() {
    registerThisObject(this.params.data.id + "|查找类型", () => {
      this.t = new Date().getTime()
      if (this.Params.data['替换类型'] === "拦截请求") {
        this.Params.data['查找范围'] = "HTTP请求/响应"
        this.Params.data['新数据'] = "拦截请求[此项不用填写]"
        this.value = "HTTP请求/响应"
        return;
      }
      if (this.Params.data['新数据'] === "拦截请求[此项不用填写]") {
        this.Params.data['新数据'] = "[双击修改-新数据]"
      }

      for (const item of this.options) {
        if (item.isBreak === true) {
          if (item.value === this.value) {
            const array = [];
            for (const item2 of this.options) {
              if (item2.isBreak === false) {
                array.push(item2)
              }
            }
            this.value = array[0].value
            return;
          }
        }
      }
    })
  },
  methods: {
    isIntervalGreaterThan100ms(time1, time2) {
      if (time1 === undefined || time2 === undefined) {
        return false;
      }
      if (time1 === null || time2 === null) {
        return false;
      }
      const difference = Math.abs(time1 - time2); // 计算时间差的绝对值
      return difference > 100; // 判断差值是否大于 100 毫秒
    },
  }
};
</script>
<template>
  <div>
    <div v-show="isBreak">
      <el-select v-model="value" class="m-2" placeholder="Select" size="small"
                 style="position: relative;top: -2px;left: -8px; width: 140px">
        <el-option
            v-for="item in BreakItems"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          <!-- 使用 el-tooltip 包裹 option 内容 -->
          <el-tooltip
              class="item"
              effect="dark"
              :content="item.tooltip"
              placement="right">
            <span>{{ item.label }}</span>
          </el-tooltip>
        </el-option>

      </el-select>
    </div>
    <div v-show="!isBreak">
      <el-select v-model="value" class="m-2" placeholder="Select" size="small"
                 style="position: relative;top: -2px;left: -8px; width: 140px">
        <el-option
            v-for="item in NoBreakItems"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          <!-- 使用 el-tooltip 包裹 option 内容 -->
          <el-tooltip
              class="item"
              effect="dark"
              :content="item.tooltip"
              placement="right">
            <span>{{ item.label }}</span>
          </el-tooltip>
        </el-option>

      </el-select>
    </div>
  </div>
</template>
