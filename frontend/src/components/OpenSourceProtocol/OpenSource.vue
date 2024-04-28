<template>
  <div style="position: relative;text-align: left;left: 20px">
    <el-tag class="ml-2" type="warning">MIT开源协议</el-tag>
    <br><br>
  </div>
  <div style="position: relative;text-align: left;left: 20px">
    <el-text class="mx-1" type="danger">　　特此向任何获得该软件副本或相关文档的人免费授予许可，可随意处理本软件，包括但不限于使用、复制、修改、合并、发布、分发、再许可和/或销售本软件的副本，并允许提供该软件的人可以按照下述条件对其进行操作：</el-text>
    <br><br>
  </div>
  <div style="position: relative;text-align: left;left: 20px">
    <el-text class="mx-1" type="danger">　　1. 本软件的所有副本或重要部分必须包含上述版权声明和本许可声明。</el-text>
    <br><br>
  </div>
  <div style="position: relative;text-align: left;left: 20px">
    <el-text class="mx-1" type="danger">　　2. 本软件按"原样"提供，不附带任何明示或暗示的保证，包括但不限于适销性、特定用途适应性和非侵权。在任何情况下，作者或版权持有人均不对任何索赔、损害或其他责任负责，无论是在合同诉讼、侵权行为或其他方面产生的、与本软件或使用或其他交易有关的或与之连接的行为。</el-text>
    <br><br><br>
  </div>

  <div style="position: relative;text-align:right">
    <el-tag class="ml-2" type="warning">版权所有 (C) 2024 秦天</el-tag>
    <br><br><br>
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