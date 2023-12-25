<template>
  <div>
    <el-tag class="ml-2" type="warning">请尽量不要在WIFI 设置代理</el-tag>
    <br><br><br>
  </div>
  <div>
    <el-text class="mx-1" type="success">安卓推荐使用：</el-text>&nbsp;
    <el-tooltip
        class="box-item"
        effect="dark"
        content="点击去下载"
        placement="top"
    >
      <el-text class="mx-1" type="danger" @click="callUrl('https://wwc.lanzouy.com/irunh08i9plc')" style="cursor: pointer; text-decoration: underline;">Kitsuneb APP</el-text>
    </el-tooltip>
    &nbsp;&nbsp;&nbsp;
    <el-text class="mx-1" type="success">设置Socks代理</el-text>  <el-text class="mx-1" type="danger">注意：需设置全局代理</el-text>
    <br><br><br>
    <el-text class="mx-1" type="success">Ios推荐使用：</el-text>
    <el-tooltip
        class="box-item"
        effect="dark"
        content="自行想办法下载"
        placement="top"
    >
      <el-text class="mx-1" type="danger">Shadowrocket</el-text>&nbsp;
    </el-tooltip>&nbsp;&nbsp;、
    <el-tooltip
        class="box-item"
        effect="dark"
        content="自行想办法下载"
        placement="top"
    >
      <el-text class="mx-1" type="danger">Quantumult</el-text>&nbsp;
    </el-tooltip>&nbsp;&nbsp;
    <el-text class="mx-1" type="success">设置Socks5代理</el-text>  <el-text class="mx-1" type="danger">注意：需设置全局代理</el-text>
    <br><br>
    <el-text class="mx-1" type="warning">在IOS中,若是Socks5代理,工作不正常,可尝试选择HTTP/HTTPS代理类型</el-text>
  </div>

  <div>
    <br><br>
    <el-text class="mx-1" type="danger">除此之外还需要注意电脑/手机、模拟器的系统时间是否正确</el-text>
  </div>
</template>

<script>

import {BrowserOpenURL} from "../../../wailsjs/runtime/runtime.js";

export default {
  data() {
    return {
      Port: ""
    }
  }, mounted() {

  }, methods: {
    callUrl(a) {
      BrowserOpenURL(a)
    }
  }
}
</script>