<template xmlns="http://www.w3.org/1999/html">
  <el-tabs v-model="activeName" class="demo-tabs"
           style="position: relative;left: 10px;display: inline-grid;justify-content: center;align-items: center;">
    <el-tab-pane label="Windows" name="Windows"/>
    <el-tab-pane label="MacOs" name="MacOs"/>
    <el-tab-pane label="Ios" name="Ios"/>
    <el-tab-pane label="Android" name="Android"/>
    <el-tab-pane label="注意事项" name="注意事项"/>
    <el-tab-pane label="导出证书" name="导出证书"/>
    <el-tab-pane label="源码" name="源码"/>
    <el-tab-pane label="开源协议" name="开源协议"/>
    <el-tab-pane label="捐助开发者" name="捐助开发者"/>
  </el-tabs>
  <div style="position: relative;height: calc(100% - 60px);width: 100%;overflow-y: auto; overflow-x: hidden">
    <el-container style="height: 100%;width: 100%">
      <el-container style="height: 100%;width: 100%">
        <el-main style="height: 100%">
          <div v-if="activeName==='Windows'">
            <Windows/>
          </div>
          <div v-if="activeName==='MacOs'">
            <MacOs/>
          </div>
          <div v-if="activeName==='Ios'">
            <IOS/>
          </div>
          <div v-if="activeName==='Android'" style="height: 100%;width: 100%">
            <Android1/>
          </div>
          <div v-if="activeName==='注意事项'">
            <h3 id="sunnynet-是完全开源的软件任何收费行为均为骗子谨防上当" style="color:red">SunnyNet
              是完全开源的软件,任何收费行为均为骗子,谨防上当</h3>
            <attention/>
          </div>
          <div v-if="activeName==='导出证书'"
               style="display: flex; justify-content: center; align-items: center; height: 80%;">
            <el-button @click="exportCert">导出证书</el-button>
            <el-button @click="checkCert">检查证书</el-button>
          </div>
          <div v-if="activeName==='源码'">
            <h3 id="sunnynet-是完全开源的软件任何收费行为均为骗子谨防上当" style="color:red">SunnyNet
              是完全开源的软件,任何收费行为均为骗子,谨防上当</h3>
            <SDK/>
          </div>
          <div v-if="activeName==='开源协议'">
            <OpenSource/>
          </div>
          <div v-if="activeName==='捐助开发者'">
            <el-text class="mx-1" type="success">在兴趣的驱动下,写一个</el-text>
            <el-text class="mx-1" type="danger"> 免费</el-text>
            <el-text class="mx-1" type="success"> 的东西，有欣喜，也还有汗水，希望你喜欢我的作品，同时也能支持一下。
            </el-text>
            <br>
            <br>
            <img :src="ImagePay" style="width: 100%;max-width: 1200px;display: block;margin: auto;">
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>

import Windows from "./CertDoc/Windows.vue";
import Attention from "./CertDoc/attention.vue";
import MacOs from "./CertDoc/MacOs.vue";
import IOS from "./CertDoc/IOS.vue";
import SDK from "./CertDoc/SDK.vue";
import OpenSource from "./CertDoc/OpenSource.vue";
import Android1 from "./CertDoc/Android.vue";
import {GetDocImage} from "../../image";
import {Dialogs} from "@wailsio/runtime";
import {AppCheckSunnyNet, ExportCert} from "../../../../../bindings/changeme/Service/appmain";
import {ElNotification} from "element-plus";

export default {
  components: {
    Android1, IOS, MacOs, Attention, Windows, SDK, OpenSource,
  },
  data() {
    return {
      activeName: "导出证书",
    }
  },
  mounted() {
  },
  beforeUnmount() {
  },
  computed: {
    ImagePay() {
      return GetDocImage("Imagepay")
    },
  },
  methods: {
    checkCert() {
      AppCheckSunnyNet().then((isInstall) => {
        if (isInstall) {
          ElNotification({
            position: 'bottom-right',
            message: 'SunnyNet 证书已安装！',
            type: 'success',
            customClass: 'multiline-message'
          })
        } else {
          ElNotification({
            position: 'bottom-right',
            message: 'SunnyNet 证书未安装！\n\n部分功能可能异常！！',
            type: 'warning',
            customClass: 'multiline-message'
          })
        }
      })
    },
    exportCert() {
      const options = {
        TreatPackagesAsDirectories: true,
        CanCreateDirectories: true,
        Filters: [
          {DisplayName: "证书", Pattern: "*.cer"},
        ],
        Filename: 'SunnyNet.cer',  // 默认文件名
        Title: "导出SunnyNet证书",
      };
      Dialogs.SaveFile(options).then((selectedFiles) => {
        const filePath = selectedFiles.substring(0, selectedFiles.lastIndexOf('/'));
        const fileName = selectedFiles.split('/').pop();
        const finalFileName = fileName.includes('.') ? fileName : `${fileName}.cer`;
        const finalFilePath = `${filePath}/${finalFileName}`;
        ExportCert(finalFilePath).then((ok) => {
          if (ok) {
            ElNotification({
              position: 'bottom-right',
              message: "导出SunnyNet证书成功",
              type: 'success',
              customClass: 'multiline-message'
            })
          } else {
            ElNotification({
              position: 'bottom-right',
              message: "导出SunnyNet证书失败",
              type: 'warning',
              customClass: 'multiline-message'
            })
          }
        })
      })
    },
  },

}
</script>