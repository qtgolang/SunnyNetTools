<script>
import {Config_Find_Window_Hide, Config_SelectedRow} from "../config/Config.js";
import Request_Message from "./Request/Request_Message.vue";

export default {
  components: {
    Request_Message,
  },

  data() {
    return {
      get SelectedRow() {
        return Config_SelectedRow.value
      },
      set SelectedRow(value) {
        Config_SelectedRow.value = value
      },
    };
  },
  computed: {
    isSelectedRow() {
      Config_Find_Window_Hide.value()
      try {
        if (this.SelectedRow != null) {
          if (parseInt(this.SelectedRow["Theology"]) !== 0 && !isNaN(parseInt(this.SelectedRow["Theology"]))) {
            return true
          }
        }
      } catch (e) {
      }
      return false
    }
  }
  ,
  methods: {},
  watch: {
    "SelectedRow"(newValue, oldValue) {
      // Config_agGrid_API.value.openToolPanel("reqs");
    }
  },
  mounted() {
  },
};
</script>

<template>
  <div class="ag-column-panel" style="position: relative" id="myToolsId">
    <div style="display: flex; justify-content: center; align-items: center; height: 100vh;" v-show="!isSelectedRow">
      <span>请选择一个会话</span>
    </div>
    <div style="display: flex; justify-content: center; align-items: center; height: 100vh;" v-show="isSelectedRow">
      <Request_Message style="height: 100%;width: 100%"/>
    </div>
  </div>

</template>