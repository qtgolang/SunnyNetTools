<script>
import {ElNotification} from "element-plus";
import {ReplaceBodyUpdate} from "../../../../bindings/changeme/Service/appmain";
import {Config_agGrid_API, getThisObject} from "../../config/Config";

export default {
  props: ['params'],
  watch: {
    "value"(m, l) {
      getThisObject(this.Params.data.id + "|查找类型")()
      if (m === this.value2) {
        return
      }
      if (m === "拦截请求") {
        this.Params.data['查找范围'] = "HTTP请求/响应"
        this.Params.data['新数据'] = "拦截请求[此项不用填写]"
      } else if (this.Params.data['新数据'] === "拦截请求[此项不用填写]") {
        this.Params.data['新数据'] = "[双击修改-新数据]"
      }
      this.Params.data['替换类型'] = m
      getThisObject(this.Params.data.id + "|查找类型")()
      const node = this.Params.data;
      ReplaceBodyUpdate(parseInt(node.id), node["替换类型"], node["查找范围"], node["旧数据"], node["新数据"], node["注释"]).then(res => {
        if (res) {
          node["状态"] = "已生效"
          ElNotification({
            showClose: true,
            message: '请求拦截/替换规则已更新',
            type: 'success',
            position: 'bottom-right',
          })
        } else {
          node["状态"] = "未生效"
          ElNotification({
            position: 'bottom-right',
            showClose: true,
            message: '请求拦截/替换规则更新失败,请检查数据格式是否正确?',
            type: 'warning',
          })
        }
        Config_agGrid_API.value.applyTransaction({update: [node]});
      })
    },
  },
  data() {
    return {
      Params: null,
      value: "字符串(UTF8)",
      value2: "字符串(UTF8)",
    }
  },
  created() {
    this.Params = this.params;
    this.value2 = this.Params.data['替换类型']
    this.value = this.Params.data['替换类型']
  },
  mounted() {

  }, methods: {}
};
</script>
<template>
  <el-select v-model="value" class="m-2" placeholder="Select" size="small"
             style="position: relative;top: -2px;left: -8px; width: 105px">

    <el-tooltip class="item" effect="dark" content="旧数据和新数据都将使用 字符串(UTF8) 解码后替换"
                placement="right">
      <el-option key="字符串(UTF8)" label="字符串(UTF8)" value="字符串(UTF8)"/>
    </el-tooltip>
    <el-tooltip class="item" effect="dark" content="旧数据和新数据都将使用 字符串(GBK) 解码后替换"
                placement="right">
      <el-option key="字符串(GBK)" label="字符串(GBK)" value="字符串(GBK)"/>
    </el-tooltip>
    <el-tooltip class="item" effect="dark" content="旧数据和新数据都将使用 Base64 解码后替换" placement="right">
      <el-option key="Base64" label="Base64" value="Base64"/>
    </el-tooltip>

    <el-tooltip class="item" effect="dark" content="旧数据和新数据都将使用十六进制解码后替换" placement="right">
      <el-option key="十六进制" label="十六进制" value="十六进制"/>
    </el-tooltip>

    <el-tooltip class="item" effect="dark" content="如果旧数据匹配上,则拦截,需要手动放行" placement="right">
      <el-option key="拦截请求" label="拦截请求" value="拦截请求"/>
    </el-tooltip>
  </el-select>
</template>
