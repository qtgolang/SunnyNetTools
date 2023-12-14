<template>
  <div>
    当前软件源码
    <br><br>
    <el-tag class="ml-2" type="warning">https://github.com/qtgolang/SunnyNetTools</el-tag>
    <br><br>
    获取SDK及源码、文档,请访问 以下任意地址
    <br><br>
    <el-tag class="ml-2" type="warning">https://esunny.vip</el-tag>
    <br>
    <el-tag class="ml-2" type="warning">https://www.esunny.vip</el-tag>
    <br>
    <el-tag class="ml-2" type="warning">https://github.esunny.vip</el-tag>
    <br>
    <el-tag class="ml-2" type="warning">https://gitee.com/qtr/SunnyNet</el-tag>
    <br>
    <el-tag class="ml-2" type="warning">https://github.com/qtgolang/SunnyNet</el-tag>
    <br> <br>
    <el-text class="mx-1" type="success">QQ交流群(一):</el-text> <el-tag class="ml-2" type="danger">751406884</el-tag><br>
    <el-text class="mx-1" type="success">QQ交流群(二):</el-text> <el-tag class="ml-2" type="danger">545120699</el-tag>
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