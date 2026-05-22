<template>
  <div style="height: 100%; width: 100%; position: relative; overflow-y: hidden">
    <div style="height: 100%; width: 100%; position: relative">
      <div :style="styleMain">
        <div
            style="width: 100%; display: flex; align-items: center; position: relative; justify-content: center; text-align: center; flex-wrap: wrap"
        >
          <Icon
              v-for="(icon, index) in iconNames"
              :key="index"
              :name="icon.Name"
              @click="Tools(icon.Name, icon.thisOpen)"
          />
        </div>

        <el-divider v-if="systemIsWindows"/>

        <div
            v-if="systemIsWindows"
            style="width: 100%; display: flex; align-items: center; position: relative; justify-content: center; text-align: center; flex-wrap: wrap"
        >
          <CustomToolIcon
              v-for="(icon, index) in customToolsList"
              :key="index"
              :obj="icon"
              @contextmenu.prevent="menu(icon, $event)"
              @click="openTools(icon)"
          />
          <span v-if="customToolsList.length === 0">这里是添加的自定义程序</span>
        </div>
      </div>

      <div
          v-if="systemIsWindows"
          class="drop-area drop-zone"
          ref="drop"
          style="position: relative; cursor: pointer"
          @click="addDropFiles" data-file-drop-target
      >
        拖动文件到此处
      </div>

      <el-drawer v-model="drawer" :with-header="false" size="480px" title="I am the title">
        <div
            style="display: flex; align-items: center; justify-content: center; text-align: center; width: 100%; flex-wrap: wrap;margin-top: 50px">
          <CustomToolIcon :obj="ThisObj" noName="true"/>
          <h1 style="margin: 0; word-break: break-word; flex: 1 1 auto; min-width: 0">
            {{ ThisObj?.Name }}
          </h1>
        </div>

        <el-divider/>

        <el-form :model="form" label-width="auto" label-position="top">
          <el-form-item label="名称:">
            <el-input v-model="form.Name"/>
          </el-form-item>

          <el-form-item label="路径:">
            <el-input v-model="form.file"/>
          </el-form-item>

          <el-form-item label="启动参数:">
            <el-input v-model="form.args"/>
          </el-form-item>

          <div style="display: flex; justify-content: flex-end; margin-top: 10px">
            <el-form-item>
              <el-button @click="drawer = false">取消</el-button>
              <el-button @click="submit">确定</el-button>
            </el-form-item>
          </div>
        </el-form>
      </el-drawer>

      <el-drawer
          v-model="valueIsShowToThis"
          :with-header="false"
          size="480px"
          style="margin-top: 32px; height: calc(100% - 32px)"
      >
        <div
            style="display: flex; align-items: center; justify-content: center; text-align: center; width: 100%; flex-wrap: wrap">
          <Icon :name="ThisOpenToolsName" noName="true"/>
          <h1 style="word-break: break-word; flex: 1 1 auto; min-width: 0; text-align: left; margin: -10px 0 0">
            {{ ThisOpenToolsName }}
          </h1>
        </div>

        <el-divider style="margin-top: -0px"/>

        <div ref="divider" :style="getThisOpenToolsStyle">
          <div v-if="ThisOpenToolsName === '颜色标记'" style="height: 100%">
            <HomeTextMark/>
          </div>
        </div>
      </el-drawer>
    </div>
  </div>
</template>

<script>
import ContextMenu from "@imengyu/vue3-context-menu";
import Icon from "./icon.vue";
import CustomToolIcon from "./CustomTool.vue";
import {OpenTools} from "../CallbackEventsOn.js";
import {
  Config_GOOS_IsWindows,
  Config_IsDark,
  Config_IsRest,
  Config_SunnyNetIsStart,
  DropFilesEvent,
  getThisObject,
} from "../config/Config.js";
import {Plus} from "@element-plus/icons-vue";
import {Dialogs, Events} from "@wailsio/runtime";
import {ElNotification} from "element-plus";
import {
  CustomToolsAdd,
  CustomToolsDel,
  CustomToolsList,
  ExecCustomTools,
  ExportCert,
  GetPort,
  SaveCustomTools,
} from "../../../bindings/changeme/Service/appmain";
import {Config_Tools_CustomList, Config_Tools_SystemList, registerHotkeyFunction} from "../config/Keys";
import HomeTextMark from "./TextMark/HomeTextMark.vue";

export default {
  components: {HomeTextMark, Plus, Icon, CustomToolIcon},

  data() {
    return {
      openToolsHeight: 0, // 内嵌工具区域距离顶部的偏移量，用于计算高度
      valueIsShowToThis: false, // 右侧工具抽屉是否显示
      drawer: false, // 自定义工具编辑抽屉是否显示
      ThisObj: null, // 当前正在编辑的自定义工具对象
      ThisOpenToolsName: "", // 当前打开的内嵌工具名称
      form: {Name: "", args: "", file: ""}, // 编辑表单数据
    };
  },

  computed: {
    // 是否处于“重置/刷新”状态（透传 Config_IsRest）
    isRest: {
      get() {
        return Config_IsRest.value;
      },
      set(val) {
        Config_IsRest.value = val;
      },
    },

    // 当前是否 Windows（透传 Config_GOOS_IsWindows）
    IsWindows: {
      get() {
        return Config_GOOS_IsWindows.value;
      },
      set(val) {
        Config_GOOS_IsWindows.value = val;
      },
    },

    // 自定义工具列表（透传 Config_Tools_CustomList）
    customList: {
      get() {
        return Config_Tools_CustomList.value;
      },
      set(val) {
        Config_Tools_CustomList.value = val;
      },
    },

    // 供模板使用：系统是否 Windows
    systemIsWindows() {
      return this.IsWindows;
    },

    // 主区域样式：Windows 下预留底部拖拽区域高度
    styleMain() {
      const h = this.IsWindows ? " - 56px" : "";
      return `height: calc(100%${h}); width: 100%; overflow-y: auto; position: relative; margin-top: 10px; margin-bottom: 10px`;
    },

    // 供模板使用：自定义工具列表
    customToolsList() {
      return this.customList;
    },

    // 内嵌工具容器高度：用 divider 的 top 偏移来算剩余空间
    getThisOpenToolsStyle() {
      return `height: calc(100% - ${this.openToolsHeight}px)`;
    },

    // 系统工具图标列表，并在首次渲染时注册热键（保持原逻辑）
    iconNames() {
      const array = [];

      Config_Tools_SystemList.forEach((node) => {
        if (!this.IsWindows && node.Windows) return;

        array.push(node);

        if (!node.thisOpen && !node.register) {
          // 这里保持原行为：只给非内嵌工具注册热键，并避免重复注册
          registerHotkeyFunction(node.ID, () => this.Tools(node.Name));
          node.register = true;
        }
      });

      return array;
    },
  },

  watch: {
    // 触发刷新时重新拉取自定义工具并注册热键
    isRest() {
      this.customList = []; // 清空旧数据，避免重复

      CustomToolsList().then((list) => {
        list.forEach((item) => {
          registerHotkeyFunction(item.ID, () => this.openTools(item)); // 为自定义工具注册热键
          this.customList.push(item); // 写回到全局配置列表
        });
      });
    },
  },

  mounted() {
    // 后端回调：添加工具结果
    Events.On("addTools", (e) => {
      const name = e.data[0];

      if (name === "-1") {
        // 添加失败：后端返回错误信息
        this.notify("warning", e.data[1]);
        return;
      }

      // 添加成功：插入到列表
      const obj = e.data[1];
      this.customToolsList.push(obj);
      this.notify("success", `添加成功：${name}`);
    });

    // 后端回调：拖拽文件事件
    Events.On("DropFilesEvent", (e) => {
      // 只处理 Main 窗口的拖拽
      if (e.data[0] !== "Main") return;

      // 一次只允许拖一个
      if (e.data[1].length !== 1) {
        this.notify("warning", "拖入失败\n\n大佬！一次只能拖入一个文件哦！");
        return;
      }

      const file = e.data[1][0];

      // 统一用 DropFilesEvent 去拿真正的落点元素
      DropFilesEvent((element) => {
        try {
          if (!this.$refs.drop.contains(element)) {
            const cb = getThisObject("registerDropFiles");
            if (cb) cb(element, file);
            return;
          }
        } catch (ex) {
          const cb = getThisObject("registerDropFiles");
          if (cb) cb(element, file);
          return;
        }

        // 落点在 drop 区域：仅允许 exe
        if (!file.toLowerCase().endsWith(".exe")) {
          this.notify("warning", "拖入失败\n\n大佬！目前只能拖入EXE文件哦！");
          return;
        }

        CustomToolsAdd(file);
      });
    });

    // 初始化：拉取自定义工具并注册热键
    CustomToolsList().then((list) => {
      list.forEach((item) => {
        registerHotkeyFunction(item.ID, () => this.openTools(item)); // 注册热键
        this.customList.push(item); // 写回配置列表
      });
    });
  },

  methods: {
    // 统一通知入口，避免到处复制 ElNotification
    notify(type, message) {
      ElNotification({
        position: "bottom-right",
        message,
        type,
        customClass: "multiline-message",
      });
    },

    // 判断是否 exe
    isExe(path) {
      return typeof path === "string" && path.toLowerCase().endsWith(".exe");
    },

    // 选择文件并添加自定义程序
    addDropFiles() {
      const options = {
        CanChooseFiles: true,
        AllowsMultipleSelection: false,
        TreatsFilePackagesAsDirectories: true,
        AllowsOtherFiletypes: false,
        Filters: [{DisplayName: "可执行文件", Pattern: "*.exe"}],
        Title: "请选择要添加的程序",
      };

      Dialogs.OpenFile(options)
          .then((selectedFile) => {
            // 只允许 exe
            if (!this.isExe(selectedFile)) {
              this.notify("warning", "添加失败\n\n大佬！目前只能添加EXE文件哦！");
              return;
            }
            CustomToolsAdd(selectedFile);
          })
          .catch(() => {
            // 用户取消选择
            this.notify("warning", "未选择文件！");
          });
    },

    // 右键菜单：删除/编辑自定义工具
    menu(obj, event) {
      event.preventDefault();

      ContextMenu.showContextMenu({
        items: [
          {
            label: "删除",
            onClick: () => this.deleteCustomTool(obj),
          },
          {
            label: "编辑",
            onClick: () => this.editCustomTool(obj),
          },
        ],
        customClass: Config_IsDark.value ? "my-dark-menu" : "",
        x: event.x,
        y: event.y,
      });
    },

    // 删除自定义工具
    deleteCustomTool(obj) {
      CustomToolsDel(obj.ID).then((ok) => {
        if (!ok) {
          this.notify("warning", "删除失败");
          return;
        }

        this.notify("success", "删除成功");

        // 过滤掉被删除项并写回配置
        this.customList = this.customToolsList.filter((item) => item.ID !== obj.ID);
      });
    },

    // 打开编辑抽屉并填充表单
    editCustomTool(obj) {
      this.form.Name = obj.Name;
      this.form.file = obj.File;
      this.form.args = obj.Args;
      this.ThisObj = obj;
      this.drawer = true;
    },

    // 执行自定义工具
    openTools(obj) {
      ExecCustomTools(obj.ID).then((res) => {
        if (res === "") {
          this.notify("success", "执行成功  [ " + obj.Name + " ] ");
          return;
        }
        this.notify("error", `执行失败  [ " + obj.Name + " ] \n\n${res}`);
      });
    },

    // 系统工具点击入口
    Tools(Name, thisOpen) {
      // 需要在抽屉内打开的工具
      if (thisOpen) {
        this.ThisOpenToolsName = Name;
        this.valueIsShowToThis = true;
        this.getOpenToolsStyle(); // 计算内嵌区域高度
        return;
      }

      // 导出证书
      if (Name === "导出证书") {
        this.exportCert();
        return;
      }

      // 代码生成
      if (Name === "代码生成") {
        OpenTools("代码生成", true, "code");
        return;
      }

      // 文本对比
      if (Name === "文本对比") {
        OpenTools("文本对比", true, "diffText");
        return;
      }

      // 证书安装
      if (Name === "证书安装") {
        this.openInstallCert();
        return;
      }

      // 其他工具按原逻辑打开
      OpenTools(Name, true, "");
    },

    // 导出 SunnyNet 证书（保持原逻辑：路径拼接 + 自动补 .cer）
    exportCert() {
      const options = {
        Filename: "SunnyNet.cer",
        TreatPackagesAsDirectories: true,
        CanCreateDirectories: true,
        Filters: [{DisplayName: "证书", Pattern: "*.cer"}],
        Title: "导出SunnyNet证书",
      };

      Dialogs.SaveFile(options).then((selectedFiles) => {
        const filePath = selectedFiles.substring(0, selectedFiles.lastIndexOf("/"));
        const fileName = selectedFiles.split("/").pop();
        const finalFileName = fileName.includes(".") ? fileName : `${fileName}.cer`;
        const finalFilePath = `${filePath}/${finalFileName}`;

        ExportCert(finalFilePath).then((ok) => {
          if (ok) {
            this.notify("success", "导出SunnyNet证书成功");
            return;
          }
          this.notify("warning", "导出SunnyNet证书失败");
        });
      });
    },

    // 打开证书安装页面（保持原逻辑：必须端口已启动，拿端口后打开本地 URL）
    openInstallCert() {
      if (Config_SunnyNetIsStart.value === false) {
        this.notify("warning", "打开证书安装失败\n您当前程序的工作端口未启动成功\n请修改端口后再试!!");
        return;
      }

      GetPort().then((port) => {
        const url = `http://localhost:${port}/install.html`;
        OpenTools("证书安装", true, url);
      });
    },

    // 保存编辑后的自定义工具
    submit() {
      if (!this.ThisObj) {
        this.drawer = false;
        return;
      }

      // 回写表单到对象（保持原逻辑）
      this.ThisObj.Name = this.form.Name;
      this.ThisObj.File = this.form.file;
      this.ThisObj.Args = this.form.args;

      const mm = this.ThisObj;
      const iconBak = mm.Icon;

      // 保存时按原逻辑临时清空 Icon
      mm.Icon = "";

      SaveCustomTools(mm).then((ok) => {
        if (!ok) {
          this.notify("warning", "保存失败"); // 提示失败
          mm.Icon = iconBak; // 失败时还原 Icon，避免对象被污染
          return;
        }
        this.notify("success", "保存成功"); // 提示成功
        const next = this.customToolsList.filter((item) => item.ID !== mm.ID);
        mm.Icon = iconBak;
        next.push(mm);
        this.customList = []; // 清空一次，确保触发视图更新
        requestAnimationFrame(() => {
          this.customList = next; // 写回新列表
        });
      });

      this.drawer = false; // 关闭抽屉
    },

    getOpenToolsStyle() {
      try {
        if (this.openToolsHeight !== 0) return;
        this.openToolsHeight = this.$refs.divider.getBoundingClientRect().top + window.scrollY;
      } catch (e) {
        requestAnimationFrame(() => this.getOpenToolsStyle());
      }
    },
  },
};
</script>

<style scoped>
.drop-area {
  width: calc(100% - 5px);
  height: 30px;
  border: 2px dashed #ccc;
  text-align: center;
  line-height: 30px;
  color: #666;
  transition: border-color 0.2s;
}

.drop-area.dragging {
  border-color: #42b983;
  background-color: #f0faff;
}
</style>

<style>
/* 整个菜单背景和边框 */
.my-dark-menu {
  background-color: #2b2b2b !important;
  border: 1px solid #444 !important;
}

/* 菜单项（正常状态） */
.my-dark-menu .mx-context-menu-item {
  border-color: #444 !important;
  color: #a9c5ec !important;
}

/* 菜单项的文字颜色 */
.my-dark-menu .mx-context-menu-item-label {
  color: #fff !important;
}

/* 鼠标 hover 高亮效果 */
.my-dark-menu .mx-context-menu-item:hover {
  background-color: #3a3a3a !important;
}

/* 悬停时文字颜色 */
.my-dark-menu .mx-context-menu-item:hover .mx-context-menu-item-label {
  color: #fff !important;
}

/* 点击/聚焦时选中效果 */
.my-dark-menu .mx-context-menu-item:active,
.my-dark-menu .mx-context-menu-item:focus {
  background-color: #444 !important;
  color: #fff !important;
}

.drop-zone {
  border: 2px dashed #ccc;
  text-align: center;
  transition: all 0.2s ease;
}

.drop-zone.file-drop-target-active {
  border-color: #007bff;
  background-color: rgba(0, 123, 255, 0.1);
}
</style>
