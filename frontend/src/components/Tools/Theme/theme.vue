<template>
  <div class="demo-collapse"
       style="margin: 0px;left: 0px;position: absolute;width: 100%;height: 100%;top: 0px;">
    <TitleBar Title="主题调色"></TitleBar>
    <div style="  display: flex;  align-items: center;  justify-content: center;">
      <el-radio-group v-model="from.IsDark">
        <el-radio value="true" size="large">暗黑配色</el-radio>
        <el-radio value="false" size="large">明亮配色</el-radio>
      </el-radio-group>
    </div>
    <el-collapse v-model="activeNames" accordion>
      <el-collapse-item title="　　主窗口-暗黑配色" name="1" v-show="Config_IsDark">
        <div style="  display: flex;  flex-wrap: wrap; align-items: center;  justify-content: center;">
          <VueText ref="drakModeText" style="width: 100%;height: 300px"></VueText>
        </div>
        <div style="  display: flex;  flex-wrap: wrap; align-items: center;  justify-content: center;">
          <el-button style="width: 90%" @click="applyColor('应用暗黑配色')">应用配色</el-button>
        </div>
        <div style="  display: flex;  flex-wrap: wrap; align-items: center;  justify-content: center;">
          <el-button style="width: 90%" @click="applyColor('默认暗黑配色1')">恢复默认暗黑配色1</el-button>
        </div>
        <div style="  display: flex;  flex-wrap: wrap; align-items: center;  justify-content: center;">
          <el-button style="width: 90%" @click="applyColor('默认暗黑配色2')">恢复默认暗黑配色2</el-button>
        </div>
        <div style="  display: flex;  flex-wrap: wrap; align-items: center;  justify-content: center;">
          <el-button style="width: 90%" @click="applyColor('默认暗黑配色3')">恢复默认暗黑配色3</el-button>
        </div>
        <div style="  display: flex;  flex-wrap: wrap; align-items: center;  justify-content: center;">
          <el-button style="width: 90%" @click="Design('暗黑')">点击配置自定义主题</el-button>
        </div>
      </el-collapse-item>
      <el-collapse-item title="　　主窗口-明亮配色" name="2" v-show="!Config_IsDark">
        <div style="  display: flex;  flex-wrap: wrap; align-items: center;  justify-content: center;">
          <VueText ref="lightModeText" style="width: 100%;height: 300px;"></VueText>
        </div>
        <div style="  display: flex;  flex-wrap: wrap; align-items: center;  justify-content: center;">
          <el-button style="width: 90%" @click="applyColor('应用明亮配色')">应用配色</el-button>
        </div>
        <div style="  display: flex;  flex-wrap: wrap; align-items: center;  justify-content: center;">
          <el-button style="width: 90%" @click="applyColor('默认明亮配色1')">恢复默认明亮配色1</el-button>
        </div>
        <div style="  display: flex;  flex-wrap: wrap; align-items: center;  justify-content: center;">
          <el-button style="width: 90%" @click="applyColor('默认明亮配色2')">恢复默认明亮配色2</el-button>
        </div>
        <div style="  display: flex;  flex-wrap: wrap; align-items: center;  justify-content: center;">
          <el-button style="width: 90%" @click="applyColor('默认明亮配色3')">恢复默认明亮配色3</el-button>
        </div>
        <div style="  display: flex;  flex-wrap: wrap; align-items: center;  justify-content: center;">
          <el-button style="width: 90%" @click="Design('明亮')">点击配置自定义主题</el-button>
        </div>
      </el-collapse-item>
      <el-collapse-item title="　　主窗口-列表项配色" name="3">
        <ListColorConfig ref="ColorConfig" style="width: 100%;height: 500px" :theme="agTheme"/>
      </el-collapse-item>
    </el-collapse>

  </div>
</template>

<script>
import VueText from "./Text.vue";
import {
  GetAgGridDarkTheme,
  GetAgGridLightTheme,
  McpFuncRes,
  Theme
} from "../../../../bindings/changeme/Service/appmain.js";
import ListColorConfig from "./ListColorConfig/ListColorConfig.vue";
import {Config_IsDark, Config_Theme_agGrid, setTheme} from "../../config/Config";
import {Events} from "@wailsio/runtime";
import TitleBar from "../../TitleBar/TitleBar.vue";

export default {
  computed: {
    Config_IsDark() {
      return Config_IsDark.value
    },
    agTheme() {
      return Config_Theme_agGrid.value
    },
  },
  components: {TitleBar, ListColorConfig, VueText},
  data() {
    return {
      url: "https://www.ag-grid.com/theme-builder/",
      activeNames: '1',
      from: {
        IsDark: "true",
      },
      get IsDark() {
        document.documentElement.className = Config_IsDark.value ? "dark" : "light";
        if (this.from.IsDark) {
          document.documentElement.style.backgroundColor = 'rgb(57,57,57)';
          document.documentElement.style.color = 'rgb(255,255,255)';
        } else {
          document.documentElement.style.backgroundColor = 'rgb(230,226,226)';
          document.documentElement.style.color = 'rgb(66,66,66)';
        }

        return Config_IsDark.value
      },
      set IsDark(value) {
        document.documentElement.className = value ? "dark" : "light";
        if (this.from.IsDark) {
          document.documentElement.style.backgroundColor = 'rgb(57,57,57)';
          document.documentElement.style.color = 'rgb(255,255,255)';
        } else {
          document.documentElement.style.backgroundColor = 'rgb(230,226,226)';
          document.documentElement.style.color = 'rgb(66,66,66)';
        }
        Config_IsDark.value = value
      },
    }
  },
  watch: {
    IsDark(n, l) {
      this.from.IsDark = n ? "true" : "false";
      document.documentElement.className = n ? "dark" : "light";
      if (this.from.IsDark) {
        document.documentElement.style.backgroundColor = 'rgb(57,57,57)';
        document.documentElement.style.color = 'rgb(255,255,255)';
      } else {
        document.documentElement.style.backgroundColor = 'rgb(230,226,226)';
        document.documentElement.style.color = 'rgb(66,66,66)';
      }
    },
    "from.IsDark"(n, l) {
      Config_IsDark.value = n === 'true'
      if (this.activeNames === '1') {
        this.activeNames = "2";
      } else if (this.activeNames === '2') {
        this.activeNames = "1";
      }
    },
  },
  methods: {
    Design(name) {
      if (name !== "明亮") {
        Theme(this.IsDark, name, this.$refs.drakModeText.GetCode())
        return
      }
      Theme(this.IsDark, name, this.$refs.lightModeText.GetCode())
    },
    applyColor(name) {
      if (name === "应用暗黑配色") {
        setTheme(this.$refs.drakModeText.GetCode(), this.IsDark)
        Theme(this.IsDark, name, this.$refs.drakModeText.GetCode())
        return
      }
      if (name === "应用明亮配色") {
        setTheme(this.$refs.lightModeText.GetCode(), this.IsDark)
        Theme(this.IsDark, name, this.$refs.lightModeText.GetCode())
        return
      }
      Theme(this.IsDark, name, "").then((code) => {
        if ((name + "").indexOf("暗黑") !== -1) {
          this.$refs.drakModeText.SetCode(code)
        } else {
          this.$refs.lightModeText.SetCode(code)
        }
        this.applyColor(this.IsDark ? "应用暗黑配色" : "应用明亮配色")
      })
    },
  },
  mounted() {
    document.documentElement.className = this.IsDark ? "dark" : "light";
    this.$refs.drakModeText.SetReadOnly(false)
    this.$refs.drakModeText.SetLanguage("json")
    this.$refs.lightModeText.SetReadOnly(false)
    this.$refs.lightModeText.SetLanguage("json")
    GetAgGridDarkTheme().then((res) => {
      this.$refs.drakModeText.SetCode(res)
    })
    GetAgGridLightTheme().then((res) => {
      this.$refs.lightModeText.SetCode(res)
    })
    Events.On("SetIsDark", (obj) => {
      const dark = obj.data[0] === true || obj.data === true;
      if (dark !== Config_IsDark.value) Config_IsDark.value = dark;
    })
    Events.On("mcp", async (evt) => {
      const mcp = evt?.data ?? {};
      // 统一回包：确保只要进到 handler，就能在需要时返回
      const reply = (text) => {
        mcp.res = text;
        typeof McpFuncRes === "function" && McpFuncRes(mcp);
      };

      try {
        const page = String(mcp.page ?? "").toLowerCase();
        const tag = String(mcp.tag ?? "").toLowerCase();
        if (page !== "theme" || tag !== "theme") return; 
        const msg = String(mcp.msg ?? "").toLowerCase();
        const toggleActiveNames = () => {
          // 只在 '1' 和 '2' 之间切换，其他值不动
          if (this.activeNames === "1") this.activeNames = "2";
          else if (this.activeNames === "2") this.activeNames = "1";
        };

        const setTheme = (isDark) => {
          this.from.IsDark = isDark;
          Config_IsDark.value = isDark;
          toggleActiveNames();
          reply(isDark ? "设置暗色主题成功" : "设置明亮主题成功");
        };

        switch (msg) {
          case "setdark":
            return setTheme(true);

          case "setlight":
            return setTheme(false);

          case "state":
            return reply(Config_IsDark.value ? "Dark" : "Light");

          default:
            // 不认识的命令：这里保持沉默（和原逻辑一致：不回包）
            return;
        }
      } catch (e) {
        // 异常兜底：避免对端一直等
        try {
          reply("处理失败");
        } catch (_) {}
      }
    });
  }
}

</script>
<style>
.light {
  background-color: white !important;
}
</style>