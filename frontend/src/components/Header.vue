<template>
  <div class="ag-header-group-cell-label " @dblclick="clickWindowButton(2)" :style="backStyle " >

    <div style="display: flex; align-items: center;" @dblclick.stop>
      <div style="cursor: pointer; display: flex; align-items: center;">

        <el-menu
            :default-active="activeIndex"
            class="el-menu-demo"
            mode="horizontal"
            :ellipsis="false"
            @select="handleSelect"
        >
          <el-sub-menu index="文件">
            <template #title>
              <img :src="logo">
            </template>
            <el-menu-item index="打开文件" @click="OpenFile()" :disabled="Stop">打开文件</el-menu-item>
            <el-sub-menu index="保存文件" :disabled="Stop">
              <template #title>保存文件</template>
              <el-menu-item index="保存选中的文件" @click="SaveToFile(false)">保存选中的文件</el-menu-item>
              <el-menu-item index="保存全部" @click="SaveToFile(true)">保存全部</el-menu-item>
            </el-sub-menu>
          </el-sub-menu>

          <el-menu-item index="设置">
            <div style="display: flex; align-items: center;position: relative;top:0px">
              <div style="cursor: pointer; display: flex; align-items: center;">
                <el-tooltip class="item" effect="dark"
                            content="程序设置"
                            placement="top">
                  <el-icon>
                    <Setting/>
                  </el-icon>
                </el-tooltip>

              </div>
            </div>
          </el-menu-item>

          <el-menu-item index="清除全部数据" :disabled="Stop">
            <div style="display: flex; align-items: center;">
              <div style="cursor: pointer; display: flex; align-items: center;position: relative;top:0px">

                <el-tooltip class="item" effect="dark"
                            content="清空所有记录"
                            placement="top">
                  <el-icon>
                    <Delete/>
                  </el-icon>
                </el-tooltip>

              </div>
            </div>
          </el-menu-item>
          <el-menu-item index="清除全部过滤条件" :disabled="Stop">
            <div style="display: flex; align-items: center;">
              <div style="cursor: pointer; display: flex; align-items: center;position: relative;top:0px">

                <el-tooltip class="item" effect="dark"
                            content="清空全部过滤条件"
                            placement="top">
                  <el-icon>
                    <Filter/>
                  </el-icon>
                </el-tooltip>

              </div>
            </div>
          </el-menu-item>
          <el-menu-item index="全部放行">
            <div style="display: flex; align-items: center;">
              <div style="cursor: pointer; display: flex; align-items: center;position: relative;top:1px">
                <el-tooltip class="item" effect="dark"
                            content="全部放行"
                            placement="top">
                  <el-icon>
                    <CaretRight/>
                  </el-icon>
                </el-tooltip>

              </div>
            </div>
          </el-menu-item>
          <el-menu-item v-if="IsWindows" index="进程驱动">
            <div style="display: flex; align-items: center;">
              <div style="cursor: pointer; display: flex; align-items: center;position: relative;top:1px">
                <el-tooltip class="item" effect="dark"
                            content="进程驱动"
                            placement="top">
                  <el-icon>
                    <Monitor/>
                  </el-icon>
                </el-tooltip>

              </div>
            </div>
          </el-menu-item>
          <el-menu-item index="脚本编辑">
            <div style="display: flex; align-items: center;">
              <div style="cursor: pointer; display: flex; align-items: center;position: relative;top:1px">
                <el-tooltip class="item" effect="dark"
                            content="脚本编辑"
                            placement="top">
                  <el-icon>
                    <EditPen/>
                  </el-icon>
                </el-tooltip>

              </div>
            </div>
          </el-menu-item>
          <el-menu-item index="自动滚动">
            <div style="display: flex; align-items: center;">
              <div style="cursor: pointer; display: flex; align-items: center;position: relative;top:1px">
                <el-icon v-show="AutoRollShow===false">
                  <CircleCloseFilled/>
                </el-icon>
                <el-icon v-show="AutoRollShow">
                  <SuccessFilled/>
                </el-icon>
                <el-tooltip class="item" content="列表自动跟随显示" placement="top">
                  <span>自动滚动显示</span>
                </el-tooltip>
              </div>
            </div>
          </el-menu-item>
        </el-menu>
      </div>

    </div>
    <PartitionOperator @dblclick.stop/>
    <!-- 文本对比 -->
    <div style="display: flex; align-items: center;" @dblclick.stop>
      <div @click="clickTextCompare" style="cursor: pointer; display: flex; align-items: center;">
        <svg :class="svgStyle" xmlns="http://www.w3.org/2000/svg" width="14" height="16" viewBox="0 0 24 24">
          <polyline points="4 7 4 4 20 4 20 7" fill="none" stroke-linecap="round" stroke-linejoin="round"
                    stroke-width="2"/>
          <line x1="9" y1="20" x2="15" y2="20" fill="none" stroke-linecap="round" stroke-linejoin="round"
                stroke-width="2"/>
          <line x1="12" y1="4" x2="12" y2="20" fill="none" stroke-linecap="round" stroke-linejoin="round"
                stroke-width="2"/>
        </svg>
        <div style="width: 2px"></div>
        文本对比
      </div>
    </div>
    <PartitionOperator @dblclick.stop/>
    <!-- 证书安装教程 -->
    <div style="display: flex; align-items: center;" @dblclick.stop>
      <div @click="clickDocCompare" style="cursor: pointer; display: flex; align-items: center;">
        <el-icon>
          <Flag/>
        </el-icon>
        <div style="width: 2px"></div>
        证书安装教程
      </div>
    </div>
    <PartitionOperator @dblclick.stop/>
    <!-- 开源协议 -->
    <div style="display: flex; align-items: center;" @dblclick.stop>
      <div @click="clickOpenSourceProtocol" style="cursor: pointer; display: flex; align-items: center;">
        <el-icon>
          <Reading/>
        </el-icon>
        <div style="width: 2px"></div>
        开源协议
      </div>
    </div>


    <div style="position: absolute; right: 95px; cursor: pointer;z-index: 1000000;width: 16px;height: 16px"
         @dblclick.stop>
      <el-popover
          placement="top-start"
          :width="200"
          trigger="hover"
          popper-style="box-shadow: rgb(14 18 22 / 35%) 0px 10px 38px -10px, rgb(14 18 22 / 20%) 0px 10px 20px -15px; padding: 20px;"
      >
        <el-table :data="WayContent">
          <el-table-column width="150" property="ip" label="当前内网IP"/>
        </el-table>
        <template #reference>
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48">
            <path
                d="M24,0A24,24,0,1,0,48,24,24,24,0,0,0,24,0ZM35.32,35.8A40.15,40.15,0,0,0,37,25h9a21.91,21.91,0,0,1-5.85,13.94A24.94,24.94,0,0,0,35.32,35.8ZM8,39.1A21.92,21.92,0,0,1,2,25h9a40.05,40.05,0,0,0,1.71,10.94A25,25,0,0,0,8,39.1Zm4.64-26.81A40.21,40.21,0,0,0,11,23H2A21.91,21.91,0,0,1,7.85,9.08,25,25,0,0,0,12.66,12.28ZM25,15a24.92,24.92,0,0,0,8.51-1.85A38.76,38.76,0,0,1,35,23H25Zm0-2V2.1c3.2.61,6.05,4.1,7.88,9.12A22.9,22.9,0,0,1,25,13ZM23,2.1V13a22.91,22.91,0,0,1-7.88-1.74C17,6.19,19.8,2.71,23,2.1ZM23,15v8H13a38.75,38.75,0,0,1,1.48-9.87A24.93,24.93,0,0,0,23,15ZM13,25H23v8.2a24.9,24.9,0,0,0-8.44,1.89A38.63,38.63,0,0,1,13,25ZM23,35.23V45.9c-3.15-.6-6-4-7.8-8.9A22.89,22.89,0,0,1,23,35.23ZM25,45.9V35.22a22.93,22.93,0,0,1,7.85,1.66C31,41.85,28.18,45.3,25,45.9Zm0-12.7V25H35a38.7,38.7,0,0,1-1.51,10A24.94,24.94,0,0,0,25,33.2ZM37,23a40.21,40.21,0,0,0-1.64-10.72,24.94,24.94,0,0,0,4.8-3.21A21.91,21.91,0,0,1,46,23ZM38.71,7.66a23,23,0,0,1-4,2.71,21,21,0,0,0-4.5-7.48A22,22,0,0,1,38.71,7.66ZM13.3,10.36a23,23,0,0,1-4-2.71,22,22,0,0,1,8.52-4.76A21,21,0,0,0,13.3,10.36ZM9.47,40.5a23,23,0,0,1,3.92-2.65,20.82,20.82,0,0,0,4.42,7.25A22,22,0,0,1,9.47,40.5Zm25.2-2.79a23,23,0,0,1,4,2.65,22,22,0,0,1-8.5,4.75A21,21,0,0,0,34.67,37.71Z"
                fill="#0797E1"/>
          </svg>
        </template>
      </el-popover>


    </div>


    <div style="position: absolute; right: 60px; cursor: pointer;z-index: 1000000;font-size: 17px;top: 8px"
         @click="clickWindowButton(1)" @dblclick.stop>
      <el-icon>
        <SemiSelect/>
      </el-icon>
    </div>
    <div style="position: absolute; right: 33px; cursor: pointer;z-index: 1000000;font-size: 17px;top: 5px"
         @click="clickWindowButton(2)" @dblclick.stop>
      <el-icon v-if="Maximise===false">
        <TopRight/>
      </el-icon>
      <el-icon v-if="Maximise===true">
        <BottomLeft/>
      </el-icon>
    </div>
    <div style="position: absolute; right: 6px; cursor: pointer;z-index: 1000000;font-size: 17px;top: 5px"
         @click="clickWindowButton(3)" @dblclick.stop>
      <el-icon>
        <SwitchButton/>
      </el-icon>
      <!--       <span class="ag-icon ag-icon-cross ag-panel-title-bar-button-icon"></span>   -->
    </div>
  </div>
  <Strings ref="Strings" v-show="ShowSetting" :show="ShowSetting"   @keydown.stop="handleKeyDown"/>
  <Doc ref="Doc" v-show="ShowDocCompare" :show="ShowDocCompare" />
  <TextCompare ref="TextCompare" v-show="ShowTextCompare" :show="ShowTextCompare"   @keydown.stop="handleKeyDown"/>
  <OpenSourceProtocol ref="OpenSourceProtocol" v-show="ShowOpenSourceProtocol" :show="ShowOpenSourceProtocol" />
</template>

<script>
import AgOption from './option.vue';
import SwitchTheme from './SwitchTheme.vue';
import Strings from './Settings.vue';
import TextCompare from './TextCompare/TextCompare.vue';
import ComboBox from './ComboBox.vue';
import PartitionOperator from "./PartitionOperator.vue";
import '../../wailsjs/runtime/runtime.js';
import {EventsOn, WindowMinimise, WindowToggleMaximise} from "../../wailsjs/runtime/runtime.js";
import {Do} from "../../wailsjs/go/main/App.js";
import {CallGoDo, EventsDo, StrBase64Encode} from "./CallbackEventsOn.js";
import {CircleCloseFilled, SuccessFilled} from '@element-plus/icons-vue'
import {ElMessage} from "element-plus";
import Doc from "./CertDoc/Doc.vue";
import OpenSourceProtocol from "./OpenSourceProtocol/OpenSourceProtocol.vue";


export default {
  components: {
    OpenSourceProtocol,
    Doc,
    PartitionOperator,
    AgOption, SwitchTheme, Strings, TextCompare, ComboBox
  },
  computed: {
    IsWindows() {
      if (window.Theme) {
        return window.Theme.GOOS === "windows"
      }
      return false
    },
    logo() {
      let c = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAG0AAAAeCAYAAAAmTpA5AAAACXBIWXMAAAsTAAALEwEAmpwYAAAF+mlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPD94cGFja2V0IGJlZ2luPSLvu78iIGlkPSJXNU0wTXBDZWhpSHpyZVN6TlRjemtjOWQiPz4gPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iQWRvYmUgWE1QIENvcmUgNS42LWMxNDUgNzkuMTYzNDk5LCAyMDE4LzA4LzEzLTE2OjQwOjIyICAgICAgICAiPiA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPiA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIiB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iIHhtbG5zOmRjPSJodHRwOi8vcHVybC5vcmcvZGMvZWxlbWVudHMvMS4xLyIgeG1sbnM6cGhvdG9zaG9wPSJodHRwOi8vbnMuYWRvYmUuY29tL3Bob3Rvc2hvcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RFdnQ9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZUV2ZW50IyIgeG1wOkNyZWF0b3JUb29sPSJBZG9iZSBQaG90b3Nob3AgQ0MgMjAxOSAoV2luZG93cykiIHhtcDpDcmVhdGVEYXRlPSIyMDIzLTExLTAxVDAyOjAzOjM4KzA4OjAwIiB4bXA6TW9kaWZ5RGF0ZT0iMjAyMy0xMS0wMVQwMjoxMzo0MCswODowMCIgeG1wOk1ldGFkYXRhRGF0ZT0iMjAyMy0xMS0wMVQwMjoxMzo0MCswODowMCIgZGM6Zm9ybWF0PSJpbWFnZS9wbmciIHBob3Rvc2hvcDpDb2xvck1vZGU9IjMiIHBob3Rvc2hvcDpJQ0NQcm9maWxlPSJzUkdCIElFQzYxOTY2LTIuMSIgeG1wTU06SW5zdGFuY2VJRD0ieG1wLmlpZDoyZTQ1MmE2NC0xOWM0LTU2NGEtODk5My0xZTVkNTEzYjU5NzYiIHhtcE1NOkRvY3VtZW50SUQ9ImFkb2JlOmRvY2lkOnBob3Rvc2hvcDo5MTg4NTg1Yi03M2U3LTJlNGYtODcyNS1lODhjYTRjNmI0ODkiIHhtcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD0ieG1wLmRpZDpmYWMxYjYxZS03MTMyLTFlNDQtYTA0NS1kN2EyNzNkN2UyNjYiPiA8eG1wTU06SGlzdG9yeT4gPHJkZjpTZXE+IDxyZGY6bGkgc3RFdnQ6YWN0aW9uPSJjcmVhdGVkIiBzdEV2dDppbnN0YW5jZUlEPSJ4bXAuaWlkOmZhYzFiNjFlLTcxMzItMWU0NC1hMDQ1LWQ3YTI3M2Q3ZTI2NiIgc3RFdnQ6d2hlbj0iMjAyMy0xMS0wMVQwMjowMzozOCswODowMCIgc3RFdnQ6c29mdHdhcmVBZ2VudD0iQWRvYmUgUGhvdG9zaG9wIENDIDIwMTkgKFdpbmRvd3MpIi8+IDxyZGY6bGkgc3RFdnQ6YWN0aW9uPSJzYXZlZCIgc3RFdnQ6aW5zdGFuY2VJRD0ieG1wLmlpZDoyZTQ1MmE2NC0xOWM0LTU2NGEtODk5My0xZTVkNTEzYjU5NzYiIHN0RXZ0OndoZW49IjIwMjMtMTEtMDFUMDI6MTM6NDArMDg6MDAiIHN0RXZ0OnNvZnR3YXJlQWdlbnQ9IkFkb2JlIFBob3Rvc2hvcCBDQyAyMDE5IChXaW5kb3dzKSIgc3RFdnQ6Y2hhbmdlZD0iLyIvPiA8L3JkZjpTZXE+IDwveG1wTU06SGlzdG9yeT4gPC9yZGY6RGVzY3JpcHRpb24+IDwvcmRmOlJERj4gPC94OnhtcG1ldGE+IDw/eHBhY2tldCBlbmQ9InIiPz5bGtmzAAARi0lEQVRoge1ae3RV1Z3+9j6Pe+4z94YkQELIawhRYkBtCZjCALUiI7ZYtUwDQ1BAWp0VW6owljpddFlmtONQ7Qy2TiyDQKSCYxsGEATkYQJJO8IoBBIJIOT9Jvfm3nte+zd/kJsmQE06i7ZTh2+tu+45v/Pbr/Ptx7d/+zAiwk38eYH/qStwE78/5KEcnvzJzkH3DAIWFFjMJas6/avgfAwBpqVIDpOz7/h6o9Ue04AzbKHb78TLxX8JQ5WgGjYWlVTB3R2FrsmQ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"
      if (!this.theme) {
        c = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAG0AAAAeCAYAAAAmTpA5AAAACXBIWXMAAAsTAAALEwEAmpwYAAAF+mlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPD94cGFja2V0IGJlZ2luPSLvu78iIGlkPSJXNU0wTXBDZWhpSHpyZVN6TlRjemtjOWQiPz4gPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iQWRvYmUgWE1QIENvcmUgNS42LWMxNDUgNzkuMTYzNDk5LCAyMDE4LzA4LzEzLTE2OjQwOjIyICAgICAgICAiPiA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPiA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIiB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iIHhtbG5zOmRjPSJodHRwOi8vcHVybC5vcmcvZGMvZWxlbWVudHMvMS4xLyIgeG1sbnM6cGhvdG9zaG9wPSJodHRwOi8vbnMuYWRvYmUuY29tL3Bob3Rvc2hvcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RFdnQ9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZUV2ZW50IyIgeG1wOkNyZWF0b3JUb29sPSJBZG9iZSBQaG90b3Nob3AgQ0MgMjAxOSAoV2luZG93cykiIHhtcDpDcmVhdGVEYXRlPSIyMDIzLTExLTAxVDAyOjAzOjM4KzA4OjAwIiB4bXA6TW9kaWZ5RGF0ZT0iMjAyMy0xMS0wMVQwMjoxNDoyMCswODowMCIgeG1wOk1ldGFkYXRhRGF0ZT0iMjAyMy0xMS0wMVQwMjoxNDoyMCswODowMCIgZGM6Zm9ybWF0PSJpbWFnZS9wbmciIHBob3Rvc2hvcDpDb2xvck1vZGU9IjMiIHBob3Rvc2hvcDpJQ0NQcm9maWxlPSJzUkdCIElFQzYxOTY2LTIuMSIgeG1wTU06SW5zdGFuY2VJRD0ieG1wLmlpZDpmNzdhNGI1NS05ZjQzLTdlNGMtOTQ0Yy1jZjEyZDllYTY0MzgiIHhtcE1NOkRvY3VtZW50SUQ9ImFkb2JlOmRvY2lkOnBob3Rvc2hvcDozODAwN2UxOC04N2U4LWNmNDAtYjZlYi0zMzY4MDdiNmVkZWIiIHhtcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD0ieG1wLmRpZDozMWE2MzA4NC1kNjA3LWFjNDYtOTA2OS1kZWQ2OGFkY2Q2ZDMiPiA8eG1wTU06SGlzdG9yeT4gPHJkZjpTZXE+IDxyZGY6bGkgc3RFdnQ6YWN0aW9uPSJjcmVhdGVkIiBzdEV2dDppbnN0YW5jZUlEPSJ4bXAuaWlkOjMxYTYzMDg0LWQ2MDctYWM0Ni05MDY5LWRlZDY4YWRjZDZkMyIgc3RFdnQ6d2hlbj0iMjAyMy0xMS0wMVQwMjowMzozOCswODowMCIgc3RFdnQ6c29mdHdhcmVBZ2VudD0iQWRvYmUgUGhvdG9zaG9wIENDIDIwMTkgKFdpbmRvd3MpIi8+IDxyZGY6bGkgc3RFdnQ6YWN0aW9uPSJzYXZlZCIgc3RFdnQ6aW5zdGFuY2VJRD0ieG1wLmlpZDpmNzdhNGI1NS05ZjQzLTdlNGMtOTQ0Yy1jZjEyZDllYTY0MzgiIHN0RXZ0OndoZW49IjIwMjMtMTEtMDFUMDI6MTQ6MjArMDg6MDAiIHN0RXZ0OnNvZnR3YXJlQWdlbnQ9IkFkb2JlIFBob3Rvc2hvcCBDQyAyMDE5IChXaW5kb3dzKSIgc3RFdnQ6Y2hhbmdlZD0iLyIvPiA8L3JkZjpTZXE+IDwveG1wTU06SGlzdG9yeT4gPC9yZGY6RGVzY3JpcHRpb24+IDwvcmRmOlJERj4gPC94OnhtcG1ldGE+IDw/eHBhY2tldCBlbmQ9InIiPz7fJdD4AAAOfklEQVRoge1aa5BcxXX+um/f59x5P/fB7grZIENZMg/xtEnAhe0EQiUVP2K7Ijt2ChwcwAQKquzEVbhwYocYGzsGXDEQTOGkMK5gHGKwQ2RjhKQ1IYBiSQhJSCuh3Z3Hzu68597bffJjZ4bRatFILhFhSl/V1PQ9fbr7dH99+p4+M4yIcAK/XeDH24ATOHqIQQrXfeuxg54ZFALoCJgjjDZ9W3E+SoAf6Jrpc3ZDpN7a6voe7EaA+ZiNb177O/AMDYYnse67kwjN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"
      }
      return c
    },
    backStyle() {
      let c = "height: 30px;background-color: #202020;"
      if (!this.theme) {
        c += "background-color: #f0f0f0;"
      } else {
        c += ""
      }
      return c + "--wails-draggable:drag"
    },
    svgStyle() {
      if (!this.theme) {
        return "dr2"
      } else {
        return "dr1"
      }
    },
    cvgStyle() {
      if (!this.theme) {
        return ""
      } else {
        return "fill: white;stroke: white"
      }
    },
    ShowSetting() {
      return this.UISettings
    },
    ShowTextCompare() {
      return this.TextCompare
    },
    ShowDocCompare() {
      return this.DocCompare
    },
    ShowOpenSourceProtocol() {
      return this.OpenSourceProtocol
    },
  },
  data() {
    return {
      Stop: false,
      AutoRollShow: false,
      WayContent: [
        {
          ip: '暂未获取到',
        },
      ],
      activeIndex: "打开文件",
      Maximise: false,
      get theme() {
        return window.Theme.IsDark
      },
      set theme(newValue) {
        window.Theme.IsDark = newValue
      },
      get UISettings() {
        return window.UI.Settings
      },
      set UISettings(newValue) {
        window.UI.Settings = newValue
      },
      get TextCompare() {
        return window.UI.TextCompare
      },
      set TextCompare(newValue) {
        window.UI.TextCompare = newValue
      },
      get DocCompare() {
        return window.UI.DocCompare
      },
      set DocCompare(newValue) {
        window.UI.DocCompare = newValue
      },
      get OpenSourceProtocol() {
        return window.UI.OpenSourceProtocol
      },
      set OpenSourceProtocol(newValue) {
        window.UI.OpenSourceProtocol = newValue
      },
      WindowSize: {min: 0, max: 0},
    }
  },
  methods: {
    ShowScriptEditing() {
      this.UISettings = false
      CallGoDo("加载配置", null).then(res => {
        this.$nextTick(() => {
          this.UISettings = true
          window.SetUILevel("Settings")
          window.vm.Settings.activeName = "脚本编辑"
        })
      })
    },
    ShowDrive() {
      this.UISettings = false
      CallGoDo("加载配置", null).then(res => {
        this.$nextTick(() => {
          this.UISettings = true
          window.SetUILevel("Settings")
          window.vm.Settings.activeName = "进程拦截"
        })
      })
    },
    rollShow() {
      this.AutoRollShow = !this.AutoRollShow
      window.vm.List.ListFollowShow = this.AutoRollShow
    },
    SetAutoRollShow(a) {
      this.AutoRollShow = a
    },
    handleSelect(key, path) {
      this.activeIndex = ""
      if (key === "设置") {
        this.clickSettings()
        return
      }
      // 清除全部数据
      if (key === "清除全部数据") {
        this.clickRemoveAll(1)
        return
      }
      // 清除全部过滤条件
      if (key === "清除全部过滤条件") {
        this.clickRemoveAll(2)
        return
      }
      //全部放行
      if (key === "全部放行") {
        this.ReleaseAll()
        return
      }
      //进程驱动
      if (key === "进程驱动") {
        this.ShowDrive()
        return
      }

      //脚本编辑
      if (key === "脚本编辑") {
        this.ShowScriptEditing()
        return
      }
      //自动滚动
      if (key === "自动滚动") {
        this.rollShow()
        return
      }


      console.log(key, path, this.activeIndex, "ok")

    },
    clickWindowButton(index) {
      if (index === 1) {
        WindowMinimise()
      } else if (index === 2) {
        let tl = window.vm.List.getTools();
        if (tl) {
          const w = parseInt(tl.style.width.replace('px', ''))

          if (this.Maximise) {
            this.WindowSize.max = w;
          } else {
            this.WindowSize.min = w;
          }
        }
        this.Maximise = !this.Maximise
        WindowToggleMaximise()
        this.$nextTick(() => {
          if (this.Maximise) {
            if (this.WindowSize.max < 30) {
              tl.style.width = '300px'
            } else {
              tl.style.width = this.WindowSize.max + 'px'
            }
          } else {
            tl.style.width = this.WindowSize.min + 'px'
          }

        })
      } else if (index === 3) {
        const objs = {}
        const filterModel = window.vm.List.agGridApi.getFilterModel();
        // 遍历过滤器模型，获取每个列的过滤器信息
        for (const colId in filterModel) {
          objs[colId] = filterModel[colId]
        }
        CallGoDo("CloseWindow", {
          Filter: StrBase64Encode(JSON.stringify(objs)),
          KeysStrings: StrBase64Encode(JSON.stringify(window.KeysStrings)),
          StorageColumns: StrBase64Encode(JSON.stringify(window.vm.List.$refs.agGrid.gridOptions.columnApi.getColumnState()))
        })
      }

    },
    clickSettings() {
      this.UISettings = false
      CallGoDo("加载配置", null).then(res => {
        this.$nextTick(() => {
          this.UISettings = true
          window.SetUILevel("Settings")
          window.vm.Settings.activeName = ""
        })
      })
    },
    clickTextCompare() {
      this.TextCompare = false
      this.$nextTick(() => {
        this.TextCompare = true
        window.SetUILevel("TextComparison")
      })
    },
    clickDocCompare() {
      this.DocCompare = false
      this.$nextTick(() => {
        this.DocCompare = true
        window.SetUILevel("DocCompare")
      })
    },
    clickOpenSourceProtocol() {
      this.OpenSourceProtocol = false
      this.$nextTick(() => {
        this.OpenSourceProtocol = true
        window.SetUILevel("OpenSourceProtocol")
      })
    },
    SaveToFile(ALL) {
      const obj = {
        Title: "请选择文件保存位置",
        Filters: [
          {Name: "SunnyNet抓包文件", Pattern: "*.syn"}
        ]
      }
      CallGoDo("保存文件对话框", obj).then(res => {
        if (res !== '') {
          const array = []
          if (!ALL) {
            for (let i = 0; i < window.vm.List.agSelectedArray.length; i++) {
              array.push(window.vm.List.agSelectedArray[i].data['Theology'])
            }
          }
          this.Stop = true
          CallGoDo("保存文件", {Path: res, ALL: ALL, Data: array}).then(res => {
            this.Stop = false
            if (res) {
              ElMessage({
                message: "文件已储存",
                type: 'success',
              })
            } else {
              ElMessage({
                message: "保存文件失败",
                dangerouslyUseHTMLString: true,
                type: 'error',
              })
            }
          })
        }
      })
    },
    OpenFile() {
      const obj = {
        Title: "请选择抓包记录文件",
        Filters: [
          {Name: "SunnyNet抓包文件", Pattern: "*.syn"},
        ]
      }
      CallGoDo("选择文件", obj).then(res => {
        if (res !== '') {
          this.Stop = true
          CallGoDo("打开记录文件", {Path: res}).then(res => {
            this.Stop = false
            if (res) {
              ElMessage({
                message: "文件已载入",
                type: 'success',
              })
            } else {
              ElMessage({
                message: "载入文件失败",
                dangerouslyUseHTMLString: true,
                type: 'error',
              })
            }
          })
        }
      })
    },
    ReleaseAll() {
      CallGoDo("全部放行", null)
    },
    clickRemoveAll(mode) {
      if (mode === 2) {
        this.$nextTick(() => {
          window.vm.List.agGridApi.setFilterModel([]); // 清空过滤器条件
          const ok = Object.keys(window.vm.List.RowDataHashMap).length < 1
          if (ok) {
            window.vm.List.RowData = [{"序号": 1001, "ico": "error"}]
            window.vm.List.RowDataHashMap = {}
            window.vm.List.agSelectedLine = null
            window.vm.List.agGridApi.setRowData(window.vm.List.RowData);
            window.vm.List.index = 0
          }
          this.$nextTick(() => {
            window.vm.List.RefreshRenderedNodes();
            if (ok) {
              window.vm.List.RowData = []
              window.vm.List.RowDataHashMap = {}
              window.vm.List.agSelectedLine = null
              window.vm.List.agGridApi.setRowData(window.vm.List.RowData);
              window.vm.List.index = 0
            }
            const columnFilter = window.vm.List.agGridApi.getFilterInstance('响应长度');
            columnFilter.setModel({
              type: 'notContains',
              filter: '0/0'
            });
            window.vm.List.agGridApi.onFilterChanged();
          })
        })
        return
      }
      CallGoDo("清空", null).then(res => {
        window.vm.List.RowData = []
        window.vm.List.RowDataHashMap = {}
        window.vm.List.agSelectedLine = null
        window.vm.List.agGridApi.setRowData(window.vm.List.RowData);
        window.vm.List.index = 0

      })
    },
    handleKeyDown(event) {
      event.stopPropagation();
    }
  },
  mounted() {
    try {
      EventsOn("Do", (Request) => {
            EventsDo(Request)
          },
      )
    } catch (e) {
      console.log("UpdateList error", e)
    }
    window.vm.Header = this

    this.$nextTick(() => {
      CallGoDo("获取内网IP", null).then(res => {
        const objs = []
        for (let i = 0; i < res.length; i++) {
          if (res[i].startsWith("169.")) {
            continue
          }
          objs.push({
            ip: res[i],
          })
        }
        if (res.length < 1) {
          this.WayContent = [{ip: "未获取到"}]
          return
        }
        this.WayContent = objs
      })

    })
  }
}
</script>
<style scoped>
/*深色模式显示*/
.dr1 {
  fill: #cccccc;
  stroke: #cccccc
}

/*浅色模式显示*/
.dr2 {
  fill: #1e1d1d;
  stroke: #1e1d1d
}

.el-menu-item {
  padding: 5px
}
</style>