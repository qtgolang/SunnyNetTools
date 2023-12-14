<template>
  <div class="common-layout">
    <el-container>
      <el-container>
        <el-main>
          <div ref="window" style="overflow-x: hidden;overflow-y: hidden">
            <div ref="b">
              <el-radio-group v-model="sslMode" class="ml-4">
                <el-radio label="默认证书" size="large">使用默认证书</el-radio>
                <el-radio label="自定义证书" size="large">使用自定义证书</el-radio>
              </el-radio-group>
            </div>
            <div v-show="sslMode==='自定义证书'" style="top: 10px;position: relative;display: inline-block">
              <div style="width: 518px;">
                <el-input style="width: 100%;text-align: left;top:5px;position: relative" v-model="CaFilePath"
                          placeholder="请选择CA文件" class="input-with-select" disabled>
                  <template #prepend> 导入CA证书：&nbsp;</template>
                  <template #append>
                    <span style="cursor:pointer;width: 20px;" @click="selectCA">...</span>
                  </template>
                </el-input>
                <el-input style="width: 100%;text-align: left;top:10px;position: relative" v-model="KeyFilePath"
                          placeholder="请选择KEY文件" class="input-with-select" disabled>
                  <template #prepend>导入证书KEY：</template>
                  <template #append>
                    <span style="cursor:pointer;width: 20px;" @click="selectKEY">...</span>
                  </template>
                </el-input>
                <div style="width: 100%;text-align: center;top:15px;position: relative">
                  <el-button v-show="CaFilePath!==''" @click="InstallCert(1)">安装选择的CA证书到系统</el-button>
                  <el-button v-show="CaFilePath!=='' && KeyFilePath!==''" @click="InstallCert(2)"
                             style="width: 134px;position: relative">
                    应用导入的证书
                  </el-button>
                </div>
                <el-divider style="width: 100%;text-align: left;top:10px;position: relative;">创建自定义证书
                </el-divider>
                <div style="width: 100%;text-align: center;top:20px;position: relative;display: flex">
                  <el-form-item label="　颁发给：">
                    <el-input v-model="createCert.domain" placeholder="证书域名 或 英文名" style="width: 175px;"/>
                  </el-form-item>
                  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
                  <el-form-item label="　　国家：">
                    <el-input v-model="createCert.country" placeholder="证书所属的国家(CN)" style="width: 175px;"/>
                  </el-form-item>
                </div>
                <div style="width: 100%;text-align: center;top:20px;position: relative;display: flex">
                  <el-form-item label="公司名称：">
                    <el-input v-model="createCert.company" placeholder="公司名称(Sunny)" style="width: 175px;"/>
                  </el-form-item>
                  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
                  <el-form-item label="部门名称：">
                    <el-input v-model="createCert.department" placeholder="所属的部门名称(Sunny)"
                              style="width: 175px;"/>
                  </el-form-item>
                </div>
                <div style="width: 100%;text-align: center;top:20px;position: relative;display: flex">
                  <el-form-item label="　所在省：">
                    <el-input v-model="createCert.province" placeholder="签发机构所在省(BeiJing)"
                              style="width: 175px;"/>
                  </el-form-item>
                  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
                  <el-form-item label="　所在市：">
                    <el-input v-model="createCert.city" placeholder="签发机构所在市(BeiJing)" style="width: 175px;"/>
                  </el-form-item>
                </div>
                <div style="width: 100%;text-align: center;top:20px;position: relative;display: flex">
                  <el-form-item label="到期时间：">
                    <el-input v-model="createCert.outTime" placeholder="到期时间/天(例:3650)" style="width: 175px;"/>
                  </el-form-item>
                  <el-form-item label="　　　　　">
                    &nbsp;&nbsp;&nbsp;&nbsp;
                    <el-button style="width: 170px;" @click="CreateCA">创建证书</el-button>
                  </el-form-item>

                </div>
              </div>
            </div>
            <div v-show="sslMode==='默认证书'" style="top: 10px;position: relative;display: inline-block">
              <div style="width: 518px;">
                <el-button v-if="IsWindows" @click="InstallCert(3)">安装默认证书</el-button>
                <el-button @click="InstallCert(4)">保存默认证书到桌面</el-button>
              </div>
            </div>
            <div>
              &nbsp;
            </div>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>
<script>

import {CallGoDo} from "../../CallbackEventsOn.js";

export default {
  components: {},
  computed: {
    IsWindows() {
      if (window.Theme) {
        return window.Theme.GOOS === "windows"
      }
      return false
    },
  },
  data() {
    return {
      sslMode: "默认证书",
      CaFilePath: "",
      KeyFilePath: "",
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
  methods: {
    CreateCA() {
      CallGoDo("创建证书", this.createCert)
    },
    selectCA() {
      const obj = {
        Title: "请选择证书CA文件",
        Filters: [
          {Name: "证书CA文件", Pattern: "*.ca;*.cer;*.pem"},
        ]
      }
      CallGoDo("选择文件", obj).then(res => {
        if (res !== '') {
          this.CaFilePath = res
        }
      })
    },
    selectKEY() {
      const obj = {
        Title: "请选择证书密钥文件",
        Filters: [
          {Name: "证书Key文件", Pattern: "*.key"},
        ]
      }
      CallGoDo("选择文件", obj).then(res => {
        if (res !== '') {
          this.KeyFilePath = res
        }
      })
    },
    InstallCert(mode) {
      if (mode === 3) {
        CallGoDo("安装默认证书", null)
        return
      }
      if (mode === 4) {
        CallGoDo("保存默认证书到桌面", null)
        return
      }
      if (mode === 1) {
        CallGoDo("安装CA证书", {CaFilePath: this.CaFilePath})
        return
      }
      CallGoDo("导入证书", {CaFilePath: this.CaFilePath, KeyFilePath: this.KeyFilePath})
    }
  },
  mounted() {

  },
  watch: {
    'sslMode': (newVal, oldVal) => {
      if (newVal === "默认证书") {
        CallGoDo("应用默认证书", null).then(res => {
          if (!res) {
            this.sslMode = "自定义证书"
          }
        })
      }
    },
  }
}
</script>