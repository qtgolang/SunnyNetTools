<script>

import {Tour_Add} from "../Tour";

export default {
  data() {
    return {
      WayContent: [],
    }
  },
  methods: {},
  mounted() {
    Tour_Add(this.$refs.xy, 1000, "开源协议", "点击这里查看开源协议,尊重作者！\n\n禁止一切违法用途，否则一切后果与作者无关！")
  }
}
</script>

<template>
  <div ref="xy">
    <el-popover
        placement="top-start"
        :width="750"
        trigger="click"
        popper-style="box-shadow: rgb(14 18 22 / 35%) 0px 10px 38px -10px, rgb(14 18 22 / 20%) 0px 10px 20px -15px; padding: 20px;"
    >
      <el-text type="success" size="large"
               style="display: flex;justify-content: center;align-content: center;text-align: center; font-weight: bold;font-size: 16px;">
        　　SunnyNet 是完全免费开源的软件,任何收费行为均为骗子,谨防上当
      </el-text>
      <el-divider/>
      <div style="margin-right: 30px">
        <div style="position: relative;text-align: left;left: 20px">
          <el-tag class="ml-2" type="warning">MIT开源协议</el-tag>
          <br><br>
        </div>
        <div style="position: relative;text-align: left;left: 20px">
          <el-text class="mx-1" type="danger">
            　　特此向任何获得该软件副本或相关文档的人免费授予许可，可随意处理本软件，包括但不限于使用、复制、修改、合并、发布、分发、再许可和/或销售本软件的副本，并允许提供该软件的人可以按照下述条件对其进行操作：
          </el-text>
          <br><br>
        </div>
        <div style="position: relative;text-align: left;left: 20px">
          <el-text class="mx-1" type="danger">　　1. 本软件的所有副本或重要部分必须包含上述版权声明和本许可声明。</el-text>
          <br><br>
        </div>
        <div style="position: relative;text-align: left;left: 20px">
          <el-text class="mx-1" type="danger">　　2.
            本软件按"原样"提供，不附带任何明示或暗示的保证，包括但不限于适销性、特定用途适应性和非侵权。在任何情况下，作者或版权持有人均不对任何索赔、损害或其他责任负责，无论是在合同诉讼、侵权行为或其他方面产生的、与本软件或使用或其他交易有关的或与之连接的行为。
          </el-text>
          <br><br><br>
        </div>
        <div style="position: relative;text-align:right">
          <el-tag class="ml-2" type="warning">版权所有 (C) 2024 秦天</el-tag>
          <br>
        </div>

      </div>
      <template #reference>
        <button class="ag-button ag-advanced-filter-builder-button" data-ref="eBuilderFilterButton" tabindex="0">
          <span data-ref="eBuilderFilterButtonIcon">
            <span class="ag-icon ag-icon-linked" style="margin-right: 5px"/>
          </span>
          <span class="ag-advanced-filter-builder-button-label">开源协议</span>
        </button>
      </template>
    </el-popover>
  </div>
</template>
