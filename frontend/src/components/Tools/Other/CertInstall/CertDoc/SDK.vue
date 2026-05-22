<template>
  <div>
    当前软件源码
    <br><br>
    <el-tag class="ml-2" type="success">https://github.com/qtgolang/SunnyNet</el-tag>
    <br><br>
    获取SDK及源码、文档,请访问 以下任意地址
    <br><br>
    <el-tag class="ml-2" type="danger" style="cursor: pointer" @click="copy('https://esunny.vip')">https://esunny.vip
    </el-tag>
    <br>
    <el-tag class="ml-2" type="danger" style="cursor: pointer" @click="copy('https://www.esunny.vip')">
      https://www.esunny.vip
    </el-tag>
    <br>
    <el-tag class="ml-2" type="danger" style="cursor: pointer" @click="copy('https://github.esunny.vip')">
      https://github.esunny.vip
    </el-tag>
    <br>
    <el-tag class="ml-2" type="danger" style="cursor: pointer" @click="copy('https://gitee.com/qtr/SunnyNet')">
      https://gitee.com/qtr/SunnyNet
    </el-tag>
    <br>
    <el-tag class="ml-2" type="danger" style="cursor: pointer" @click="copy('https://github.com/qtgolang/SunnyNet')">
      https://github.com/qtgolang/SunnyNet
    </el-tag>
    <br> <br>
    <el-text class="mx-1" type="success">QQ交流群 :</el-text>
    <el-tag class="ml-2" type="danger" style="cursor: pointer" @click="copy('751406884')">751406884</el-tag>
    <el-tag class="ml-2" type="danger" style="cursor: pointer" @click="copy('545120699')">545120699</el-tag>
    <el-tag class="ml-2" type="danger" style="cursor: pointer" @click="copy('170902713')">170902713</el-tag>
    <br>
    <br>
    <el-text class="mx-1" type="success">QQ频道:</el-text>
    <el-tag class="ml-2" type="danger" style="cursor: pointer" @click="copy('https://pd.qq.com/g/SunnyNetV5')">
      https://pd.qq.com/g/SunnyNetV5
    </el-tag>
  </div>
</template>

<script>

import {ClipboardWriteAll} from "../../../../../../bindings/changeme/Service/appmain";
import {ElNotification} from "element-plus";

export default {
  data() {
    return {
      Port: ""
    }
  }, mounted() {

  }, methods: {
    copy(a) {
      ClipboardWriteAll(a).then(res => {
        if (res !== '') {
          ElNotification({
            position: 'bottom-left',
            message: '复制失败\n\n' + res,
            type: 'warning',
            customClass: 'multiline-message'
          })
          return;
        }
        ElNotification({
          position: 'bottom-left',
          message: '复制成功\n\n已复制到剪辑版',
          type: 'success',
          customClass: 'multiline-message'
        })
      })
    }
  }
}
</script>