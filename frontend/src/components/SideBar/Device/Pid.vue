<script>
import {AgGridVue} from "ag-grid-vue3";
import {Config_Theme_agGrid} from "../../config/Config.js";
import {AG_GRID_LOCALE_CN} from "../../config/AG_ZH_CN.js";
import {Events} from "@wailsio/runtime";
import {ProcessAddPid, ProcessDelPid} from "../../../../bindings/changeme/Service/appmain";
import {attachMcpDevicePid} from "../../config/mcpDeviceSync.js";
import MustTcp from "../Settings/MustTcp.vue";
import Host from "../Settings/Host.vue";
import BaseSettings from "../Settings/BaseSettings.vue";
import Keys from "../Settings/Keys.vue";
import HTTPSProto from "../Settings/HTTPSProto.vue";
import Way from "../Settings/Way.vue";

export default {
  components: {Way, HTTPSProto, Keys, BaseSettings, Host, MustTcp, 'ag-grid-vue': AgGridVue},
  data() {
    return {
      Stopped: null,
      agGridApi: null,
      rowData: [],
      addValue: null,
      RowId: 0,
      overlayNoRowsTemplate: `<span style="padding: 20px;" >还没有添加任何进程</span>`,
      defaultColDef: {
        flex: 1,
        sortable: false,
        suppressHeaderMenuButton: true,
        suppressHeaderContextMenu: true,
      },
      ReadOnly: false,
      RowNodes: [],
      CellForClipboard: null,
      previousSelectedIds: new Set(),
      FindValue: "",
      gridOptions: {
        cellSelection: true,
        suppressMovableColumns: true,
        stopEditingWhenCellsLoseFocus: true, // 失去焦点时自动结束编辑
        getRowId: (params) => params.data.id,
        getContextMenuItems: this.MenuEvent,
        localeText: AG_GRID_LOCALE_CN,
        isExternalFilterPresent: this.isExternalFilterPresent,
        doesExternalFilterPass: this.doesExternalFilterPass,
        onSelectionChanged: (event) => {
          const currentSelectedRows = event.api.getSelectedRows();
          const currentSelectedIds = new Set(currentSelectedRows.map(row => row.id));
          // 新增的 ID（在当前有，但之前没有）
          const added = [...currentSelectedIds].filter(id => !this.previousSelectedIds.has(id));
          // 移除的 ID（之前有，但当前没有）
          const removed = [...this.previousSelectedIds].filter(id => !currentSelectedIds.has(id));
          removed.forEach((id) => {
            ProcessDelPid(parseInt(id))
          });
          added.forEach((id) => {
            ProcessAddPid(parseInt(id))
          });
          // 更新为当前选择
          this.previousSelectedIds = currentSelectedIds;
        },
        columnDefs: [
          {
            field: "PID", tooltipField: 'PID',
            minWidth: 70,
            width: 70,
            maxWidth: 70,
          },
          {
            field: "进程名称", tooltipField: '进程名称',
            minWidth: 210,
            width: 310,
          },
        ],
        rowSelection: {
          mode: 'multiRow',
          checkboxes: true,
          headerCheckbox: false,
          copySelectedRows: false,
          enableClickSelection: true,
        },
        onCellSelectionChanged: (params) => {
          if (params.started) {
            this.RowNodes.forEach((node) => {
              node.setSelected(false);
            });
          }
          if (params.finished) {
            const rangeSelections = params.api.getCellRanges();
            const selectedRowNodes = [];
            let ___c = true;
            rangeSelections.forEach(function (range) {
              const _a = range.startRow.rowIndex;
              const _b = range.endRow.rowIndex;

              let startRow = _a;
              let endRow = _b;
              if (_a > _b) {
                startRow = _b
                endRow = _a
                ___c = false;
              } else {
                ___c = true;
              }
              params.api.forEachNodeAfterFilter(function (node) {
                if (node.rowIndex >= startRow && node.rowIndex <= endRow) {
                  selectedRowNodes.push(node);
                  node.setSelected(true);
                }
              });
            });
            this.RowNodes = selectedRowNodes;
          }
        },
        onCellEditingStarted: (params) => {
        },
        onCellEditingStopped: (params) => {

        },
      }
    }
  },
  watch: {
    FindValue(n) {
      this.agGridApi.onFilterChanged()
    }
  },
  computed: {
    agTheme() {
      return Config_Theme_agGrid.value
    },
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.api;
    attachMcpDevicePid(this);
    Events.On("DeviceUpdateProcessesList", (obj) => {
      const pid = parseInt(obj.data[0])
      const Name = obj.data[1]
      const isDel = obj.data[2]
      if (!isDel) {
        this.addLine(pid, Name)
        return
      }
      const node = this.agGridApi.getRowNode(pid + "")
      if (node) {
        this.agGridApi.applyTransaction({remove: [node.data]});
      }
    })
  },
  methods: {
    isExternalFilterPresent(params) {
      return this.FindValue !== "";
    },
    doesExternalFilterPass(node) {
      if (node.data['PID'].indexOf(this.FindValue) < 0) {
        if (node.data['进程名称'].toLowerCase().indexOf(this.FindValue.toLowerCase()) < 0) {
          return false;
        }
      }
      return true;
    },
    addLine(pid, Name) {
      const data = [{
        "进程名称": Name,
        "PID": pid + "",
        id: pid + "",
      }];
      this.agGridApi.applyTransaction({add: data});
    },
    SetReadOnly(readOnly) {
      this.ReadOnly = readOnly;
      this.gridOptions.columnDefs.forEach((obj) => {
        obj.editable = !readOnly
      })
      this.agGridApi.setGridOption('columnDefs', this.gridOptions.columnDefs);
    },
    MenuEvent(params) {
      return [];
    },
    deselectAll() {
      this.agGridApi.deselectAll();
    },
    Empty() {
      this.RowNodes = [];
      this.agGridApi.setGridOption("rowData", []);
    },
  }
}
</script>

<template>
  <div style="height: 100%;width: 100%">
    <el-input v-model="FindValue" style="width: 100%;height: 30px" placeholder="请输入 PID/进程名称,进行查找(不区分大小写)"/>
    <ag-grid-vue ref="agGrid"
                 :theme="agTheme"
                 :rowData="rowData"
                 style="height: calc(100% - 30px);width: 100%"
                 :grid-options="gridOptions"
                 :loading="false"
                 :allowContextMenuWithControlKey="true"
                 :defaultColDef="defaultColDef"
                 :overlayNoRowsTemplate="overlayNoRowsTemplate"
                 :suppressCutToClipboard="true"

    />
  </div>
</template>
<style>
.white-svg path {
  stroke: white;
}
</style>
