<template>

  <div ref="bh" style="display: flex; align-items: center;" @click="ClickEve">
    <div
        style="cursor: pointer; display: flex; align-items: center;font-size: 15px;justify-content:center;width: 100px">
      <img :src="img" style="height: 25px;"><img>
      {{ Title }}
    </div>
  </div>
</template>

<script>
import {GetImage} from "../../Tools/image.js";
import {McpFuncRes, SetWorking} from "../../../../bindings/changeme/Service/appmain.js";
import {Tour_Add} from "../Tour";
import {Events} from "@wailsio/runtime";

export default {
  props: ["isHideHook"],
  data() {
    return {
      state: true,
      Title: "正在捕获",
      img: GetImage("开始捕获")
    }
  },
  methods: {
    ClickEve() {
      this.ClickEves(!this.state)
    },
    ClickEves(a) {
      this.state = a
      this.img = GetImage(this.state ? "开始捕获" : "停止捕获")
      this.Title = this.state ? "正在捕获" : "隐藏捕获"
      if (this.state) {
        try {
          document.getElementById("HookMessageText").innerText = "还没有捕获到数据";
        } catch (e) {
        }
      } else {
        try {
          document.getElementById("HookMessageText").innerText = "您隐藏了捕获数据";
        } catch (e) {

        }
      }
      SetWorking(!this.state).then(() => {
        this.isHideHook(!this.state)
      })
    }
  },
  mounted() {
    Tour_Add(this.$refs.bh, 2, "显示/隐藏捕获", "点击这里切换捕获状态,隐藏捕获后\n\n程序依旧在工作只是不显示到列表区域")

    Events.On("mcp", async (evt) => {
      const mcp = evt?.data ?? {};
      const reply = (text) => {
        mcp.res = text;
        typeof McpFuncRes === "function" && McpFuncRes(mcp);
      };
      try {
        const page = String(mcp.page ?? "").toLowerCase();
        const tag = String(mcp.tag ?? "").toLowerCase();
        if (page !== "main" || tag !== "home") return;
        const msg = String(mcp.msg ?? "");
        switch (msg) {
          case "start":
            if (this.state) {
              return reply("当前已经是开启状态");
            }
            this.ClickEves(true)
            return reply("设置成功");
          case "stop": {
            if (!this.state) {
              return reply("当前已经是关闭状态");
            }
            this.ClickEves(false)
            return reply("设置成功");
          }
          case "state":
            return reply(this.Title);
          default:
            return;
        }
      } catch (e) {
        // 异常兜底：避免对端一直等
        try {
          reply("处理失败");
        } catch (_) {
        }
      }
    });
  }
}
</script>
<style scoped>

</style>