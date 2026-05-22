<script>
import {GetProxyDns, SetProxyDns} from "../../../../bindings/changeme/Service/appmain";
import {Config_IsRest} from "../../config/Config";
import {attachMcpSettingsReload} from "../../config/mcpSettingsSync.js";

export default {
  data() {
    return {
      set isRest(val) {
        Config_IsRest.value = val;
      },
      get isRest() {
        return Config_IsRest.value;
      },
      radio: '',//不重要
      drawer2: false,
      remoteDns: "223.5.5.5:853",
      isInit: false,
    }
  },
  methods: {
    reloadDnsFromBackend() {
      this.isInit = false;
      return GetProxyDns().then((proxy) => {
        if (proxy === "localhost" || proxy === "local") {
          this.radio = "local";
        } else if (proxy === "remote") {
          this.radio = "remote";
        } else {
          this.radio = "remotes";
          this.remoteDns = proxy;
        }
        setTimeout(() => {
          this.isInit = true;
        }, 100);
      });
    },
    applyRemoteDns() {
      let dns = ""
      if (this.radio === "local") {
        dns = "localhost"
      } else if (this.radio === "remote") {
        dns = "remote"
      } else {
        dns = this.remoteDns
      }
      if (!this.isInit) {
        return;
      }
      SetProxyDns(dns).then(() => {
        this.$message({
          message: "DNS更改成功",
          type: "success",
        })
      })
    },
  },
  watch: {
    radio(val) {
      this.applyRemoteDns()
    },
    isRest() {
      this.isInit = false
      GetProxyDns().then(proxy => {
        if (proxy === "localhost" || proxy === "local") {
          this.radio = "local"
        } else if (proxy === "remotes") {
          this.radio = "remotes"
          this.remoteDns = proxy
        } else {
          this.radio = "remote"
        }
        setTimeout(() => {
          this.isInit = true
        }, 100)
      })
    }
  },
  mounted() {
    attachMcpSettingsReload("proxy_dns", () => this.reloadDnsFromBackend());
    GetProxyDns().then(proxy => {
      if (proxy === "localhost"|| proxy === "local") {
        this.radio = "local"
      } else if (proxy === "remotes") {
        this.radio = "remotes"
        this.remoteDns = proxy
      } else {
        this.radio = "remote"
      }
      setTimeout(() => {
        this.isInit = true
      }, 100)
    })
  }
}
</script>

<template>
  <div style="display: flex; position: relative; gap:0px; justify-content: center;">
    <el-radio v-model="radio" value="local">本地解析</el-radio>
    <el-radio v-model="radio" value="remote">远程解析</el-radio>
    <el-radio v-model="radio" value="remotes">远程服务器解析</el-radio>
  </div>
  <div
      style="display: flex; position: relative; gap:50px; margin-top: 15px;margin-bottom: 10px;justify-content: center;"
      v-if="radio==='remotes'">
    <div style="display: flex; gap: 10px">
      <el-input v-model="remoteDns" placeholder="请输入远程DNS" input-style="text-align: center;"/>
      <el-button @click="applyRemoteDns">更改设置</el-button>
    </div>
  </div>
  <div
      style="display: flex; position: relative; gap:50px; margin-top: 15px;margin-bottom: 10px;justify-content: center;">
    <el-button @click="drawer2 = true">查看DNS设置说明</el-button>
  </div>
  <el-drawer v-model="drawer2" direction="rtl" size="50%">
    <template #header>
      <h4>DNS解析设置说明</h4>
    </template>
    <template #default>
      <div>
        <div>
          <el-divider content-position="left">情况1</el-divider>
          <span>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;你没有使用全局上游代理，也没有对请求单独设置代理,这种情况下，无论你设置的是那种模式，都只会使用<span style="color: #f80013;font-size: 20px">本地解析</span>!
          </span>
          <br>
          <br>
          <el-divider content-position="left">情况2</el-divider>
          <span>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;2.你使用全局上游代理，或对请求单独设置代理,这种情况下使用 <span style="color: #f80013;font-size: 20px">本地解析</span> 模式,将通过本地网络进行DNS解析
          </span>
          <br>
          <el-divider content-position="left">情况3</el-divider>
          <span>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span style="color: #10eab2">你使用了全局上游代理或请求单独设置代理,这种情况下，你设置以下3种模式会有区别</span>
          </span>
          <br>
          <br>
          <span>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;1.使用 <span style="color: #f80013;font-size: 20px">本地解析</span> 模式，你要访问的目标地址，通过你本地DNS解析出的IP，可能会被服务器拒绝连接。这时候你需要尝试 <span style="color: #02d4f8;font-size: 16px">远程解析 或 远程服务器解析</span>
          </span>
          <br>
          <br>
          <span>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;2.使用 <span style="color: #f80013;font-size: 20px">远程解析</span> 模式，你所使用的代理服务器可能存在无法解析的情况。这时你应该尝试 <span style="color: #02d4f8;font-size: 16px">远程服务器解析</span>
          </span>
          <br>
          <br>
          <span>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;3.使用 <span style="color: #f80013;font-size: 20px">远程服务器解析</span> 模式， 会使用你设置的代理，连接到远程DNS服务器进行查询并且解析，可能会导致首次访问变慢
          </span>
          <br>
          <el-divider content-position="left">总结</el-divider>

          <span>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  <span style="color: #9ca4e1;font-size: 20px" >选择使用什么模式，视情况而定！！</span>
          </span>
          <br>
          <br>
          <span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>
          <el-tag type="warning">
            * 仅支持,TLS 的远程DNS服务器,也就是853端口,53端口的DNS暂不支持
          </el-tag>
        </div>
      </div>
    </template>
    <template #footer>
      <div style="flex: auto">
        <el-button @click="drawer2 = false">朕知道了</el-button>
      </div>
    </template>
  </el-drawer>
</template>
