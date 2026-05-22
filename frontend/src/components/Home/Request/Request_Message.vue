<template>

  <div style="width: 100%;height: 100%">
    <div style="width: 100%;height: 100%" v-show="isHTTPMessage">
      <HTTP_Message style="width: 100%;height: 100%"/>
    </div>
    <div style="width: 100%;height: 100%" v-show="isTCPMessage||isUDPMessage">
      <Socket_Message style="width: 100%;height: 100%"/>
    </div>
  </div>
</template>
<script>
import {Config_SelectedRow} from "../../config/Config.js";
import HTTP_Message from "./Message/HTTP_Message.vue";
import Socket_Message from "./Message/Socket_Message.vue";

export default {
  components: {Socket_Message, HTTP_Message},
  computed: {
    isHTTPMessage() {
      if (!this.isSelectedRow()) return false;
      return this.SelectedRow["请求地址"].toLowerCase().startsWith("http");
    },
    isTCPMessage() {
      if (!this.isSelectedRow()) return false;
      return this.SelectedRow["方式"].toLowerCase().indexOf("tcp") !== -1;
    },
    isUDPMessage() {
      if (!this.isSelectedRow()) return false;
      return this.SelectedRow["方式"].toLowerCase().indexOf("udp") !== -1;
    },
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
  methods: {
    isSelectedRow() {
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
  },
  watch: {
    "SelectedRow"(newValue) {

    }
  },
  mounted() {
    /*

          GetHTTPSession(parseInt(newValue["Theology"])).then((res) => {
            console.log(res)
          })
     */
  }
}
</script>