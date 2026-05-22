<script>

import {ClipboardWriteAll} from "../../../../bindings/changeme/Service/appmain";
import {ElNotification} from "element-plus";

export default {
  data() {
    return {
      WayContent: [],
    }
  },
  methods: {
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
  },
  mounted() {
  }
}
</script>

<template>
  <div>
    <el-popover
        placement="top-start"
        :width="590"
        trigger="click"
        popper-style="box-shadow: rgb(14 18 22 / 35%) 0px 10px 38px -10px, rgb(14 18 22 / 20%) 0px 10px 20px -15px; padding: 20px;"
    >
      <el-text type="success" size="large"
               style="display: flex;justify-content: center;align-content: center;text-align: center; font-weight: bold;font-size: 16px;margin-right: 20px">
        　　SunnyNet 是完全免费开源的软件,任何收费行为均为骗子,谨防上当
      </el-text>
      <el-divider/>
      <div style="margin-right: 30px;">
        当前软件源码
        <el-tag class="ml-2" type="success">https://github.com/qtgolang/SunnyNet</el-tag>
        <br><br>
        获取SDK及源码、文档,请访问 以下任意地址
        <br><br>
        <div style="display: flex; flex-wrap: wrap;">
          <el-tag style="margin-top: 5px; margin-left: 10px;cursor: pointer" type="danger"
                  @click="copy('https://esunny.vip')">https://esunny.vip
          </el-tag>
          <el-tag style="margin-top: 5px; margin-left: 10px;cursor: pointer" type="danger"
                  @click="copy('https://www.esunny.vip')">https://www.esunny.vip
          </el-tag>
          <el-tag style="margin-top: 5px; margin-left: 10px;cursor: pointer" type="danger"
                  @click="copy('https://github.esunny.vip')">https://github.esunny.vip
          </el-tag>
          <el-tag style="margin-top: 5px; margin-left: 10px;cursor: pointer" type="danger"
                  @click="copy('https://gitee.com/qtr/SunnyNe')">https://gitee.com/qtr/SunnyNet
          </el-tag>
          <el-tag style="margin-top: 5px; margin-left: 10px;cursor: pointer" type="danger"
                  @click="copy('https://github.com/qtgolang/SunnyNet')">https://github.com/qtgolang/SunnyNet
          </el-tag>
        </div>
        <br>
        <div>
          <el-text class="mx-1" style="margin-top: 3px; margin-left: 10px" type="success">QQ交流群 :</el-text>
          <el-tag class="ml-2" style="margin-top: 0; margin-left: 10px;cursor: pointer" type="danger"
                  @click="copy('751406884')">751406884
          </el-tag>
          <el-tag class="ml-2" style="margin-top: 0; margin-left: 10px;cursor: pointer" type="danger"
                  @click="copy('545120699')">545120699
          </el-tag>
          <el-tag class="ml-2" style="margin-top: 0; margin-left: 10px;cursor: pointer" type="danger"
                  @click="copy('170902713')">170902713
          </el-tag>
        </div>
        <br>
        <el-text class="mx-1" style="margin-top: 0; margin-left: 23px;" type="success">QQ频道 :</el-text>
        <el-tag class="ml-2" style="margin-top: -5px; margin-left: 10px;cursor: pointer" type="danger"
                @click="copy('https://pd.qq.com/g/SunnyNetV5')">https://pd.qq.com/g/SunnyNetV5
        </el-tag>
      </div>
      <template #reference>
        <button class="ag-button ag-advanced-filter-builder-button" data-ref="eBuilderFilterButton" tabindex="0">
          <span data-ref="eBuilderFilterButtonIcon">
            <span class="ag-icon ag-icon-eye" style="margin-right: 5px"/>
          </span>
          <span class="ag-advanced-filter-builder-button-label" style="margin-right: 10px">源代码</span>
        </button>
      </template>
    </el-popover>
  </div>
</template>
