<template>
  <div id="IEProxyDiv" ref="IE"
       style="width: 180px;cursor: pointer; font-size: 15px;justify-content:center; display: flex;">
    <div v-show="isSettings===false"
         style="width: 180px;cursor: pointer; font-size: 15px;justify-content:center; display: flex;" @click="Click">
      {{ title }}
    </div>
    <div v-show="isSettings" style="width: 100px;cursor: pointer; ">正在设置中...</div>
  </div>

</template>

<script>


import {CancelIEProxy, GOOS, McpFuncRes, SetIEProxy} from "../../../../bindings/changeme/Service/appmain.js";
import {Keys_System_id_Cancel_IE_Agent, registerHotkeyFunction} from "../../config/Keys";
import {ElMessage, ElNotification} from "element-plus";
import {Config_GOOS_IsWindows, registerThisObject} from "../../config/Config";
import {Tour_Add} from "../Tour";
import {Events} from "@wailsio/runtime";

export default {
  data() {
    return {
      state: false,
      title: "未设置系统代理",
      isSettings: false,
      MTitle: "",
    }
  },
  watch: {
    "state": function (newVal, oldVal) {
      this.isSettings = false
      if (this.state) {
        this.title = "已设置系统" + this.MTitle + "代理"
      } else {
        this.title = "未设置系统" + this.MTitle + "代理"
      }
    }
  },
  methods: {
    Click() {
      if (this.isSettings) {
        return
      }
      if (!this.state) {
        SetIEProxy().then(ok => {
          this.isSettings = false
          if (ok) {
            this.state = true
          } else {
            if (!Config_GOOS_IsWindows.value) {
              ElNotification({
                showClose: true,
                message: "设置系统代理失败\n检查是否输入密码\r\n密码是否输入正确?",
                type: 'error',
                position: 'bottom-right',
                customClass: 'multiline-message'
              })
            } else {
              ElMessage.warning("设置系统" + this.MTitle + "代理失败")
            }
          }
        })
      } else {
        CancelIEProxy().then(ok => {
          this.isSettings = false
          if (ok) {
            this.state = false
          } else {
            if (!Config_GOOS_IsWindows.value) {
              ElNotification({
                showClose: true,
                message: "设置系统代理失败\n检查是否输入密码\r\n密码是否输入正确?",
                type: 'error',
                position: 'bottom-right',
                customClass: 'multiline-message'
              })
            } else {
              ElMessage.warning("设置系统" + this.MTitle + "代理失败")
            }
          }
        })
      }
    }
  },
  mounted() {
    registerHotkeyFunction(Keys_System_id_Cancel_IE_Agent, this.Click)
    registerThisObject("SetIEProxyState", (o) => {
      this.state = o
    })
    if (Config_GOOS_IsWindows.value) {
      this.MTitle = "IE"
    }
    GOOS().then(IsWindows => {
      if (IsWindows) {
        this.MTitle = "IE"
      }
      if (this.state) {
        this.title = "已设置系统" + this.MTitle + "代理"
      } else {
        this.title = "未设置系统" + this.MTitle + "代理"
      }

      Tour_Add(this.$refs.IE, 1, "设置IE代理", "点击这里 设置/取消IE代理")
    })
    Events.On("mcp", async (evt) => {
      try {
        const mcp = evt.data || {};
        const page = String(mcp.page || "").toLowerCase();
        const tag = String(mcp.tag || "").toLowerCase();
        const msg = String(mcp.msg || "").toLowerCase();

        if (page !== "main") return;
        if (tag !== "systemproxy") return;

        const reply = (text) => {
          mcp.res = text;
          McpFuncRes(mcp);
        };

        const isWindows = !!Config_GOOS_IsWindows.value;

        const failMsg = (actionText) => {
          if (!isWindows) return `${actionText}系统代理失败\n检查是否输入密码\r\n密码是否输入正确?`;
          return `${actionText}系统${this.MTitle}代理失败`;
        };

        if (msg === "set") {
          if (this.state) return reply("当前系统代理已经是开启状态");

          const ok = await SetIEProxy();
          this.isSettings = false;

          if (ok) {
            this.state = true;
            return reply("系统代理已开启");
          }
          return reply(failMsg("设置"));
        }

        if (msg === "cancel") {
          if (!this.state) return reply("当前系统代理已经是取消状态");

          const ok = await CancelIEProxy();
          this.isSettings = false;

          if (ok) {
            this.state = false;
            return reply("系统代理已取消");
          }
          return reply(failMsg("取消"));
        }

        if (msg === "state") {
          return reply(this.state ? "当前系统代理是开启状态" : "当前系统代理是取消状态");
        }

      } catch (e) {
        // 可选：异常也回包，避免 Go 端一直等
        try {
          const mcp = evt.data || {};
          mcp.res = "处理失败";
          McpFuncRes(mcp);
        } catch (_) {}
      }
    });
  }
}
</script>
<style scoped>

</style>