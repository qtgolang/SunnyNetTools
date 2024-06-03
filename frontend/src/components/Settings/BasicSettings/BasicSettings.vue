<template>
  <div ref="window">
    <div style="text-align: center;position: relative;left: 0px;top: 0px;height: 30px" ref="a">
      <span>&nbsp;&nbsp;程序工作端口：</span>
      <el-input v-model="Port" type="number" min="0" size="small" placeholder="请输入端口号"
                style="width: 120px;position: relative;top: 0px"/>
      &nbsp;
      <el-tooltip class="item" effect="dark" content="确定修改程序运行的端口号、支持Socket5、HTTP、HTTPS代理协议"
                  placement="top">
        <el-button icon="Check" size="small" circle style="position: relative;top: 0px" @click="submitPort"/>
      </el-tooltip>
    </div>

    <div style="text-align: center;position: relative;left: 0px;top: 0px;height: 30px">
    </div>
    <div style="">
      <el-tooltip class="item" effect="dark"
                  content="[请使用Pro版本] 开启后每个请求的TLS指纹将随机变化"
                  placement="top">
        <el-checkbox label="随机TLS指纹" disabled/>
      </el-tooltip>
      <el-tooltip class="item" effect="dark"
                  content="[设置Socket5代理有效] 将禁止非HTTP/S的TCP连接,某些APP,将先尝试TCP请求,如果TCP请求失败才会发送HTTP请求,这种场景下有用"
                  placement="top">
        <el-checkbox v-model="Option.DisableTCP" label="禁用TCP"/>
      </el-tooltip>
      <el-tooltip class="item" effect="dark"
                  content="[手机端设置Socket5代理有效、PC加载驱动有效] 将禁止发送、接收UDP数据,例如某手APP,若不禁用UDP,某些关键数据将捕获不到"
                  placement="top">
        <el-checkbox v-model="Option.DisableUDP" label="禁用UDP"/>
      </el-tooltip>
      <el-tooltip class="item" effect="dark" content="让浏览器不要缓存文件,每次请求都重新加载所有文件" placement="top">
        <el-checkbox v-model="Option.DisableBrowserCache" label="禁止浏览器缓存"/>
      </el-tooltip>
      <el-tooltip class="item" effect="dark"
                  content="开启后客户端只能设置Socket5代理、如果客户端设置的是HTTP、HTTPS代理,将会被拒绝请求"
                  placement="top">
        <el-checkbox v-model="Option.authentication" label="开启身份验证模式"/>
      </el-tooltip>
      <el-tooltip class="item" effect="dark" content="重置所有配置"
                  placement="top">
        <el-button icon="QuestionFilled" size="small" circle style="position: relative;top: -3px;left: 30px"
                   @click="ResetAll"/>
      </el-tooltip>
    </div>

    <div style="text-align: left;position: relative;left: 95px;">
      <List ref="UserInfo" v-show="Option.authentication" :style="getParentElementWidth"/>
    </div>
    <div style="text-align: left;position: relative;left: 95px;">
      <Shortcutkeys ref="Keys" :style="getParentElementWidth2"></Shortcutkeys>
    </div>

  </div>
</template>
<script>
import {Check} from '@element-plus/icons-vue'
import List from "./List.vue";
import {CallGoDo, SunnyErrorReplaceAll} from "../../CallbackEventsOn.js";
import {ElMessage, ElMessageBox} from "element-plus";
import Shortcutkeys from "./Shortcutkeys.vue";

export default {
  components: {Shortcutkeys, List},
  computed: {
    Check() {
      return Check
    },
    getParentElementWidth() {
      return "width: " + (window.Size.Settings.Width - (97 * 2)) + "px;height: " + (window.Size.Settings.Height - 220) + "px"
    },
    getParentElementWidth2() {
      return "width: " + (window.Size.Settings.Width - (97 * 2)) + "px"
    }
  },
  data() {
    return {
      Port: 8089,
      Option: {
        DisableUDP: false,
        DisableTCP: false,
        authentication: false,
        DisableBrowserCache: false
      }
    }
  },
  watch: {
    'Option.DisableTCP': (newVal, oldVal) => {
      CallGoDo("禁止TCP", {DisableTCP: newVal})
    },
    'Option.DisableUDP': (newVal, oldVal) => {
      CallGoDo("禁止UDP", {DisableUDP: newVal})
    },
    'Option.authentication': (newVal, oldVal) => {
      CallGoDo("身份验证模式", {authentication: newVal})
    },
    'Option.DisableBrowserCache': (newVal, oldVal) => {
      CallGoDo("禁止缓存", {DisableBrowserCache: newVal})
    }
  },
  methods: {
    ResetAll() {
      CallGoDo("重置所有配置", {Port: this.Port}).then(res => {
        this.Option.DisableTCP = false
        this.Option.DisableUDP = false
      })
    },
    submitPort() {
      CallGoDo("修改端口号", {Port: this.Port}).then(res => {
        if (res) {
          if (res.ok) {
            ElMessage({
              message: "修改端口号成功",
              type: 'success',
            })
            return
          }
          if (res.err === "Port error") {
            ElMessage({
              message: "输入了错误的端口号！1-65535",
              dangerouslyUseHTMLString: true,
              type: 'error',
            })
            return
          }
          ElMessageBox.alert(
              "端口修改成功,但启动失败:<br><br>" + SunnyErrorReplaceAll(res.err),
              "请更换端口",
              {
                dangerouslyUseHTMLString: true,
                confirmButtonText: '好的',
                closeOnClickModal: true, // 设置点击遮罩层关闭消息框
                closeOnPressEscape: true, // 设置按下 ESC 键关闭消息框
              }
          )
        }
      })
    },
  },
  mounted() {

  }
}
</script>