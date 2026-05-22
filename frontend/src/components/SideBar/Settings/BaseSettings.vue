<script>
import Table from "../../Tools/table.vue"
import {
  Config_HomeTextMark,
  Config_IsRest,
  Config_Status_Info,
  Config_SunnyNetIsStart,
  DefaultRowData,
  getThisObject,
  ObjString
} from "../../config/Config";
import {ElMessage, ElNotification} from "element-plus";
import {attachMcpSettingsReload} from "../../config/mcpSettingsSync.js";
import {
  AppIsSetPort,
  AuthenticationList,
  AuthenticationRemove,
  AuthenticationUpdate,
  CreateAuthentication,
  GetBaseSettingsValue,
  GetHomeTextMark,
  GetPort,
  IsDark,
  ResetALLConfig,
  SetAuthMode,
  SetDisableCache,
  SetDisableTCP,
  SetDisableUDP,
  SetIsDark,
  SetLimitRequestSize,
  SetPort
} from "../../../../bindings/changeme/Service/appmain";

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
      Port: 0,
      Size: 1024000,
      DisableTCP: false,
      DisableUDP: false,
      DisableBrowserCache: false,
      authentication: false,
    }
  },
  watch: {
    isRest() {
      GetPort().then(port => {
        this.Port = port
        SetPort(parseInt(this.Port), false).then(error => {
          AppIsSetPort().then(State => {
            const o = getThisObject("SetIEProxyState")
            if (o) {
              o(State)
            }
          })
          if (error === "") {
            Config_SunnyNetIsStart.value = true
            Config_Status_Info.value = `启动成功：端口号[${this.Port}]`;
          } else {
            Config_SunnyNetIsStart.value = false
            Config_Status_Info.value = error;
          }
        })
      })
      IsDark().then(isDark => {
        SetIsDark(isDark)
      })
      {
        GetBaseSettingsValue().then((obj) => {
          this.DisableTCP = obj[0]
          this.DisableUDP = obj[1]
          this.DisableBrowserCache = obj[2]
          this.authentication = obj[3]
          this.Size = obj[4]
        })
      }
      this.$refs.authentication.Empty()
    },
    authentication(n) {
      SetAuthMode(n)
    },
    DisableTCP(n) {
      SetDisableTCP(n)
    },
    DisableUDP(n) {
      SetDisableUDP(n)
    },
    DisableBrowserCache(n) {
      SetDisableCache(n)
    },
    reloadBaseSettingsFromBackend() {
      return GetBaseSettingsValue().then((obj) => {
        this.DisableTCP = obj[0];
        this.DisableUDP = obj[1];
        this.DisableBrowserCache = obj[2];
        this.authentication = obj[3];
        this.Size = obj[4];
      });
    },
    Size(n) {
      if (n < 1024) {
        ElNotification({
          position: 'bottom-right',
          showClose: true,
          message: '限制请求提交大小\n范围错误 不能小于1024',
          type: 'warning',
          customClass: 'multiline-message'
        })
        return
      }
      SetLimitRequestSize(n)
    },
  },
  methods: {
    init() {
      GetPort().then(port => {
        this.Port = port;
      })
      //身份验证模式 初始化表格信息
      {
        this.$refs.authentication.agGridApi.hideOverlay();
        this.$refs.authentication.agGridApi.setGridOption('overlayNoRowsTemplate', `<span style="padding: 20px;" id="HookMessageText">您还没有添加身份验证的账号</span>`);
        this.$refs.authentication.agGridApi.showNoRowsOverlay();
        this.$refs.authentication.agGridApi.setGridOption('columnDefs', [
          {
            field: "账号", tooltipField: '账号',
            minWidth: 170,
            width: 170,
            maxWidth: 170,
            editable: true,
          },
          {
            field: "密码", tooltipField: '密码',
            minWidth: 170,
            maxWidth: 170,
            width: 170,
            editable: true,
          },
        ]);
        this.$refs.authentication.Stopped = this.AuthenticationEdit
        this.$refs.authentication.addValue = this.addAuthentication
        this.$refs.authentication.DeleteID = this.delAuthentication
      }
      {
        GetBaseSettingsValue().then((obj) => {
          this.DisableTCP = obj[0]
          this.DisableUDP = obj[1]
          this.DisableBrowserCache = obj[2]
          this.authentication = obj[3]
          this.Size = obj[4]
        })
      }
    },
    //应用修改端口
    submitPort() {
      if (parseInt(this.Port) < 1 || parseInt(this.Port) > 65534) {
        ElNotification({
          position: 'bottom-right',
          showClose: true,
          message: '修改工作端口\n你输入的端口号范围错误 1-65534',
          type: 'warning',
          customClass: 'multiline-message'
        })
        return
      }
      GetPort().then(port => {
        if (parseInt(this.Port) === parseInt(port)) {
          ElNotification({
            position: 'bottom-right',
            showClose: true,
            message: '端口没有发生变化！！',
            type: 'warning',
            customClass: 'multiline-message'
          })
          return
        }
        SetPort(parseInt(this.Port), false).then(error => {
          AppIsSetPort().then(State => {
            const o = getThisObject("SetIEProxyState")
            if (o) {
              o(State)
            }
          })
          if (error === "") {
            ElNotification({
              position: 'bottom-right',
              showClose: true,
              message: '修改工作端口\n' + `启动成功：端口号[${this.Port}]`,
              type: 'success',
              customClass: 'multiline-message'
            })
            Config_SunnyNetIsStart.value = true
            Config_Status_Info.value = `启动成功：端口号[${this.Port}]`;
          } else {
            Config_SunnyNetIsStart.value = false
            ElNotification({
              position: 'bottom-right',
              showClose: true,
              message: '修改工作端口\n重新启动失败:' + error,
              type: 'warning',
              customClass: 'multiline-message'
            })
            Config_Status_Info.value = error;
          }
        })
      })
    },
    //重置所有配置
    ResetAll() {
      ElNotification({
        position: 'bottom-right',
        showClose: true,
        message: '正在重置...',
        type: 'success',
        customClass: 'multiline-message'
      })
      setTimeout(() => {
        ResetALLConfig().then(() => {

          GetHomeTextMark().then((res) => {
            Config_HomeTextMark.clear()
            try {
              JSON.parse(res).forEach((item) => {
                Config_HomeTextMark.set(item.id, item)
              })
            } catch (e) {
            }
            if (Config_HomeTextMark.size < 1) {
              for (let i = 0; i < DefaultRowData.length; i++) {
                Config_HomeTextMark.set(DefaultRowData[i].id, DefaultRowData[i])
              }
            }
            const func = getThisObject("ColorMarking")
            if (func) {
              func()
            }
            Config_IsRest.value++
            ElNotification({
              position: 'bottom-right',
              showClose: true,
              message: '所有设置已重置',
              type: 'success',
              customClass: 'multiline-message'
            })
          })

        })
      }, 1500)
    },
    //删除身份验证账号
    delAuthentication(id) {
      AuthenticationRemove(parseInt(id))
    },
    //添加身份验证账号密码
    addAuthentication() {
      CreateAuthentication().then(id => {
        this.$refs.authentication.agGridApi.applyTransaction({
          add: [{
            "账号": "账号:双击修改",
            "密码": '密码:双击修改',
            id: (id + "")
          }]
        });
      })
    },
    //身份验证码模式完成编辑
    AuthenticationEdit(params) {
      const array = [];
      const isUserExist = (user) => {
        return array.some(item => item.user === user);
      }
      let isExist = false;
      this.$refs.authentication.agGridApi.forEachNode(node => {
        const user = ObjString(node.data["账号"]).trim();
        const pass = ObjString(node.data["密码"]).trim();
        const id = ObjString(node.data.id).trim();
        if (!isUserExist(user)) {
          array.push({user: user, pass: pass, id: id})
          return
        }
        isExist = true
      });
      if (isExist) {
        this.$refs.authentication.Empty()
        const lodArray = [];
        array.forEach(item => {
          lodArray.push({"账号": item.user, "密码": item.pass, id: item.id});
        })
        this.$refs.authentication.agGridApi.applyTransaction({add: lodArray});
        ElMessage.error('您添加的账号已经存在,请重新添加！')
        return
      }
      const node = params.node;
      const user = ObjString(node.data["账号"]).trim();
      const pass = ObjString(node.data["密码"]).trim();
      const id = ObjString(node.data.id).trim();
      AuthenticationUpdate(parseInt(id), user, pass).then(ok => {
        if (ok) {
          ElNotification({
            position: 'bottom-right',
            showClose: true,
            message: '身份验证账号更新完成',
            type: 'success',
            customClass: 'multiline-message'
          })
        } else {
          ElNotification({
            position: 'bottom-right',
            showClose: true,
            message: '身份验证账号更新失败',
            type: 'warning',
            customClass: 'multiline-message'
          })
          AuthenticationRemove(id).then(() => {
            this.$refs.authentication.agGridApi.applyTransaction({remove: [node.data]});
          })
        }
      })
    },
  },
  computed: {},
  mounted() {
    this.init()
    attachMcpSettingsReload("base", () => this.reloadBaseSettingsFromBackend());
    AuthenticationList().then(list => {
      list.forEach(item => {
        this.$refs.authentication.agGridApi.applyTransaction({
          add: [{
            "账号": item.User,
            "密码": item.Pass,
            id: (item.ID + "")
          }]
        });
      })
    })
  }
}
</script>

<template>
  <div style="display: flex; flex-direction: column; gap: 5px; margin: 5px;">
    <!-- 程序工作端口 -->
    <div style="display: flex; position: relative; align-items: center; margin-left: 15px; width: 100%;">
      <span>程序工作端口：</span>
      <el-input
          v-model="Port"
          type="number"
          min="1"
          max="65534"
          size="small"
          placeholder="请输入端口号"
          style="width: 100px; margin-right: 10px; text-align: center;"
          input-style="text-align: center;"
      ></el-input>
      <el-tooltip content="确定修改程序运行的端口号、支持Socket5、HTTP、HTTPS代理协议" placement="top">
        <el-button icon="Check" size="small" circle @click="submitPort"/>
      </el-tooltip>
      <el-tooltip
          class="item"
          effect="dark"
          content="重置所有配置"
          placement="top"
      >
        <el-button
            icon="QuestionFilled"
            size="small"
            circle
            style="margin-left: auto; margin-right: 20px;"
            @click="ResetAll"
        />
      </el-tooltip>
    </div>

    <!-- 禁止TCP/UDP/浏览器缓存 -->
    <div style="display: flex; position: relative; gap:30px; margin-left: 15px">
      <el-tooltip placement="right">
        <template #content>
          <div style="white-space: normal; line-height: 1.4;">
            开启后将禁止非HTTP/S的TCP连接
            <br>
            <br>
            某些APP,将先尝试TCP请求,如果TCP请求失败才会发送HTTP请求
            <br>
            这种场景下有用
          </div>
        </template>
        <div
            class="ag-labeled ag-label-align-left ag-toggle-button ag-input-field ag-group-title-bar ag-charts-format-sub-level-group-title-bar ag-unselectable ag-selected">
          <div class="ag-input-field-label ag-label ag-toggle-button-label" style="margin-right: 10px">
            禁止TCP
          </div>
          <el-switch v-model="DisableTCP" size="small"/>
        </div>
      </el-tooltip>
    </div>
    <div style="display: flex; position: relative; gap:30px; margin-left: 15px">
      <el-tooltip placement="right">
        <template #content>
          <div style="white-space: normal; line-height: 1.4;">
            [手机端设置Socket5代理有效]
            <br>
            [PC加载驱动有效]
            <br><br>
            将禁止发送、接收UDP数据，例如某手APP，
            <br>
            若不禁用UDP，某些关键数据将捕获不到
          </div>
        </template>
        <div
            class="ag-labeled ag-label-align-left ag-toggle-button ag-input-field ag-group-title-bar ag-charts-format-sub-level-group-title-bar ag-unselectable ag-selected">
          <div class="ag-input-field-label ag-label ag-toggle-button-label" style="margin-right: 10px">
            禁止UDP
          </div>
          <el-switch v-model="DisableUDP" size="small"/>
        </div>
      </el-tooltip>
    </div>
    <div style="display: flex; position: relative; gap:30px; margin-left: 15px">
      <el-tooltip placement="right">
        <template #content>
          <div style="white-space: normal; line-height: 1.4;">
            禁止浏览器缓存(实验阶段)
            <br><br>
            让浏览器不要缓存文件,每次请求都重新加载所有文件
            <br>
            否则浏览器有缓存后可能不会从服务器获取新的资源
          </div>
        </template>
        <div
            class="ag-labeled ag-label-align-left ag-toggle-button ag-input-field ag-group-title-bar ag-charts-format-sub-level-group-title-bar ag-unselectable ag-selected">
          <div class="ag-input-field-label ag-label ag-toggle-button-label" style="margin-right: 10px">
            禁止浏览器缓存
          </div>
          <el-switch v-model="DisableBrowserCache" size="small"/>
        </div>
      </el-tooltip>
    </div>
    <!-- 身份验证模式/限制提交请求大小 -->
    <div style="position: relative; gap:30px; justify-content:inherit;">
      <div style="display: flex; position: relative; gap:30px; margin-left: 15px">
        <el-tooltip placement="right">
          <template #content>
            <div style="white-space: normal; line-height: 1.4;">
              若开启身份验证模式后
              <br>
              <br>
              客户端只能设置Socket5代理
              <br>
              如果客户端设置的是HTTP、HTTPS代理,将会被拒绝请求
            </div>
          </template>
          <div
              class="ag-labeled ag-label-align-left ag-toggle-button ag-input-field ag-group-title-bar ag-charts-format-sub-level-group-title-bar ag-unselectable ag-selected">
            <div class="ag-input-field-label ag-label ag-toggle-button-label" style="margin-right: 10px">
              开启身份验证模式
            </div>
            <el-switch v-model="authentication" size="small"/>
          </div>
        </el-tooltip>
      </div>
      <div style="width: 100%;height: 300px" v-show="authentication">
        <Table ref="authentication"/>
      </div>
    </div>
    <div style="display: flex; position: relative; gap:30px; margin-left: 25px">
      <el-tooltip placement="right">
        <template #content>
          <div style="white-space: normal; line-height: 1.4;">
            如果请求是 HTTP / HTTPS1.1 / HTTPS2.0
            <br>
            <br>
            设置后,POST 请求提交数据，超过此限制后,不解析原始数据,以便优化内存
            <br>
            例如上传一个1G的文件,若是不限制,抓包软件将会加载这个1G的文件,可能会导致一些问题
            <br>
          </div>
        </template>
        <div
            class="ag-labeled ag-label-align-left ag-toggle-button ag-input-field ag-group-title-bar ag-charts-format-sub-level-group-title-bar ag-unselectable ag-selected">
          <div class="ag-input-field-label ag-label ag-toggle-button-label"
               style="margin-right: 5px;margin-left: -10px">
            限制请求提交大小:
          </div>
          <input class="ag-input-field-input ag-number-field-input" type="number" min="10240"
                 v-model="Size" style="width: 70px; text-align: center; padding: 0 5px;"/>
          <div class="ag-input-field-label ag-label ag-toggle-button-label" style="margin-left: 5px">字节</div>
        </div>
      </el-tooltip>
    </div>
  </div>
</template>
