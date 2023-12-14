<template>
  <div class="common-layout">
    <el-container>
      <el-container>
        <el-main>
          <div ref="window">
            <el-radio-group v-model="MustTcp" class="ml-4">
              <el-tooltip class="item" effect="dark"
                          content="开启后所有请求都将不会解密SSL,全部转为TCP请求处理"
                          placement="top">
                <el-radio label="ALL" size="large">全部强制走TCP</el-radio>
              </el-tooltip>
              <el-tooltip class="item" effect="dark"
                          content="自定义配置,强制走TCP的IP/域名,配置后符合条件的请求将不会解密SSL"
                          placement="top">
                <el-radio label="custom" size="large">自定义配置</el-radio>
              </el-tooltip>

            </el-radio-group>
            <div v-show="MustTcp==='custom'" >
                <Text  style="height: 200px;width: 100%;position: relative" />
              <el-button style="position: relative;top: 5px">保存</el-button>
            </div>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>
<script>

import Text from "../../Request/Text.vue";

export default {
  components: {Text},
  computed: {},
  data() {
    return {
      sslMode: "默认证书",
      CaFilePath: "",
      KeyFilePath: "",
      MustTcp: "ALL",
      createCert: {
        //域名
        domain: "",
        //国家
        country: "",
        //公司
        company: "",
        //部门
        department: "",
        //省
        province: "",
        //市
        city: "",
        //到期时间/天
        outTime: "",
      }
    }
  },
  methods: {},
  mounted() {

  },
  watch: {
    'MustTcp': (newVal, oldVal) => {
      // 在这里执行 sslMode 变化后的操作
      console.log('sslMode 变化了:', newVal);
    },
  }
}
</script>