<script>
import Table from "../../Tools/table.vue";
import {Config_IsRest, ObjString} from "../../config/Config";
import {ElMessage, ElNotification} from "element-plus";
import {ReplaceHostList, ReplaceHostRemove, ReplaceHostUpdate} from "../../../../bindings/changeme/Service/appmain";
import {attachMcpConfigReload} from "../../config/mcpRulesSync.js";

export default {
  components: {Table},
  data() {
    return {
      set isRest(val) {
        Config_IsRest.value = val;
      },
      get isRest() {
        return Config_IsRest.value;
      },
    }
  },
  watch: {
    isRest() {
      this.Empty()
      ReplaceHostList().then(list => {
        const array = []
        for (let i = 0; i < list.length; i++) {
          const item = list[i];
          array.push({
            "旧的": item["Lod"],
            "新的": item["New"],
            "注释": item["Note"],
            id: item["ID"] + "",
          })
        }
        this.agGridApi.applyTransaction({add: array});
      })
    }
  },
  methods: {
    reloadHostFromBackend() {
      if (!this.agGridApi) {
        return Promise.resolve();
      }
      this.Empty();
      return ReplaceHostList().then((list) => {
        const array = [];
        for (let i = 0; i < list.length; i++) {
          const item = list[i];
          array.push({
            "旧的": item["Lod"],
            "新的": item["New"],
            "注释": item["Note"],
            id: item["ID"] + "",
          });
        }
        this.agGridApi.applyTransaction({add: array});
      });
    },
    Empty() {
      this.agGridApi.setGridOption("rowData", []);
    },
    init() {
      //初始化表格信息
      {
        this.$refs.Host.agGridApi.hideOverlay();
        this.$refs.Host.agGridApi.setGridOption('overlayNoRowsTemplate', `<span style="padding: 20px;" id="HookMessageText">您还没有添加任何HOST替换</span>`);
        this.$refs.Host.agGridApi.showNoRowsOverlay();
        this.$refs.Host.agGridApi.setGridOption('columnDefs', [
          {
            field: "旧的", tooltipField: '旧的',
            minWidth: 150,
            width: 150,
            maxWidth: 570,
            editable: true,
          },
          {
            field: "新的", tooltipField: '新的',
            minWidth: 150,
            maxWidth: 550,
            width: 500,
            editable: true,
          },
          {
            field: "注释", tooltipField: '注释',
            minWidth: 450,
            maxWidth: 450,
            width: 450,
            editable: true,
          },
        ]);
        this.$refs.Host.Stopped = this.HostEdit
        this.$refs.Host.addValue = this.addHost
        this.$refs.Host.DeleteID = this.delHost
      }
    },
    addHost() {
      this.$refs.Host.agGridApi.applyTransaction({
        add: [{
          "旧的": "www.test.com",
          "新的": 'www.new.com',
          "注释": '',
          id: (this.$refs.Host.rowId++) + ""
        }]
      });
      ElMessage.success('已添加,请双击修改')
    },
    delHost(id) {
      ReplaceHostRemove(parseInt(id))
    },
    HostEdit(params) {
      const array = [];
      const isUserExist = (user) => {
        return array.some(item => item.Lod === user);
      }
      let isExist = false;
      this.$refs.Host.agGridApi.forEachNode(node => {
        const Lod = ObjString(node.data["旧的"]).trim();
        const New = ObjString(node.data["新的"]).trim();
        const Note = ObjString(node.data["注释"]).trim();
        if (!isUserExist(Lod)) {
          array.push({Lod: Lod, New: New, Note: Note})
          return
        }
        isExist = true
      });
      if (isExist) {
        this.$refs.Host.Empty()
        const lodArray = [];
        array.forEach(item => {
          lodArray.push({
            "旧的": item.Lod,
            "新的": item.New,
            "注释": item.Note,
            id: (this.$refs.Host.rowId++) + ""
          });
        })
        this.$refs.Host.agGridApi.applyTransaction({add: lodArray});
        ElMessage.error('您添加的旧 Host 已存在,请重新添加')
        return
      }
      const node = params.node.data;
      ReplaceHostUpdate(parseInt(node.id), node["旧的"], node["新的"], node["注释"]).then(() => {
        ElNotification({
          position: 'bottom-right',
          message: 'Host替换已经更新',
          type: 'success',
          customClass: 'multiline-message'
        })
      })
    },
  },
  mounted() {
    this.agGridApi = this.$refs.Host.agGridApi;
    this.init()
    attachMcpConfigReload("host", () => this.reloadHostFromBackend());
    this.reloadHostFromBackend();
  }
}
</script>

<template>
  <div>
    <Table ref="Host" style="height: 300px"/>
  </div>
</template>
