<template>

  <div ref="targetElement" :style="isNoDisableClick">
    <div style="width: calc(100% - 20px);height: 100%;font-size: 15px;margin-top: 10px">
      <el-form
          ref="ruleFormRef"
          label-width="auto"
          class="demo-ruleForm"
          size="default"
          status-icon
      >
        <el-form-item label="查找内容：" prop="region">
          <Autocomplete ref="autoEdit" style="" @keyup.enter="handleKeyEnter"/>
        </el-form-item>
        <el-form-item label="查找类型：" prop="region">
          <FindType ref="FindType"/>
          <div
              style="margin-left: auto; text-align: right;display: flex; justify-content: flex-end;margin-top: 5px;margin-bottom: -15px">
            <el-checkbox-group v-model="FindConfig.check" style="" v-show="!isIntType && FindConfig.Type!=='pb' ">
              <el-checkbox value="不区分大小写" v-show="CaseSensitive">不区分大小写</el-checkbox>
              <el-checkbox value="删除空格后搜索">删除空格后搜索</el-checkbox>
            </el-checkbox-group>
            <div style="margin-left: 20px" v-show="FindConfig.Type==='pb' ">
              <span>忽略前: </span>
              <input type="number" style="width: 30px" value="0" min="0" ref="pbSkip"/>
              <span> 字节 </span>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="查找范围：" prop="region">
          <FindRange ref="FindRange"/>
        </el-form-item>

        <el-form-item label="搜索动作：" prop="region">
          <el-radio-group v-model="FindConfig.action" style="display: flex;width: 100%" ref="action">
            <el-radio-button size="small" value="color" style="flex: 1;width: 50%" border>颜色标记</el-radio-button>
            <el-radio-button size="small" value="hide" style="flex: 1;width: 50%" border>隐藏不相关项</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="颜色标记：" ref="Color" prop="region" v-show="FindConfig.action==='color'">
          <FindColor ref="FindColor" style="width: calc(100% - 170px)"/>
          <el-checkbox-group v-model="FindConfig.check" style="width: 150px;margin-left: 20px">
            <el-checkbox value="取消之前的颜色标记">取消之前的颜色标记</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label=" " prop="region">
          <el-button size="large" @click="findStart" style="width: 100%;height: 20px;margin-top: -10px">查找</el-button>
          <el-progress :percentage="per" color="#13EA00" style="width: 100%;height: 20px; " v-show="ShowSchedule"/>
        </el-form-item>

      </el-form>

    </div>
  </div>
</template>

<script>
import Autocomplete from "./Autocomplete.vue";
import FindType from "./FindType.vue";
import FindRange from "./FindRange.vue";
import FindColor from "./FindColor.vue";
import {
  Config_Find,
  Config_Find_Range_ALL,
  Config_Find_Range_Options_All,
  Config_Find_Range_Options_Websocket,
  Config_Find_Type_Options,
  Config_Find_Window,
  Config_SelectedRow,
  DisableClick,
  Find_Type_Del_Option
} from "../../config/Config.js";
import {FindSession} from "../../../../bindings/changeme/Service/appmain";
import {ElMessage} from "element-plus";
import {Events} from "@wailsio/runtime";

export default {
  computed: {
    isShowSchedule() {
      return this.ShowSchedule
    },

    isNoDisableClick() {
      return (this.isDisableClick ? "pointer-events: none;" : "") + "display: block"
    },
    isIntType() {
      const n = this.FindConfig.Type;
      return n === "整数4" || n === "整数8" || n === "浮点数4" || n === "浮点数8"
    },
  },
  watch: {
    "FindConfig.action"(n) {
      if (n === "color") {
        Config_Find_Window.value.setNewHeight(Config_Find_Window.value.Height + 30)
      } else {
        Config_Find_Window.value.setNewHeight(Config_Find_Window.value.Height - 30)
      }
    },
    "FindConfig.value"(n) {
      let types = JSON.parse(JSON.stringify(Config_Find_Type_Options))
      const num = parseInt(n, 10);
      if (isNaN(num)) {
        types = Find_Type_Del_Option(types, "整数4")
        types = Find_Type_Del_Option(types, "整数8")
        types = Find_Type_Del_Option(types, "浮点数4")
        types = Find_Type_Del_Option(types, "浮点数8")
        const m = this.$refs.FindType.GetType();
        if (m === "整数4" || m === "整数8" || m === "浮点数4" || m === "浮点数8") {
          this.$refs.FindType.SetType("UTF8")
        }
      }
      if (!/^([0-9A-Fa-f]{2})+$/.test(n)) {
        types = Find_Type_Del_Option(types, "Hex")
        if (this.$refs.FindType.GetType() === "Hex") {
          this.$refs.FindType.SetType("UTF8")
        }
      }
      this.$refs.FindType.SetOptions(types)
    },
    "FindAllRange"(all) {
      this.$refs.FindRange.SetOptions(all ? Config_Find_Range_Options_All : Config_Find_Range_Options_Websocket)
    },
    "FindConfig.Color"(color) {
      const firstChild = this.$refs.Color.$el.firstElementChild;
      if (firstChild && firstChild.firstElementChild) {
        const target = firstChild.firstElementChild;
        target.style.color = this.$refs.FindColor.GetColorValue(color)
        target.style.fontWeight = "bold";
      }
    }
  },
  components: {FindColor, FindRange, FindType, Autocomplete},
  data() {
    return {
      get FindConfig() {
        return Config_Find.value
      },
      set FindConfig(value) {
        Config_Find.value = value
      },
      get FindAllRange() {
        return Config_Find_Range_ALL.value
      },
      set FindAllRange(value) {
        Config_Find_Range_ALL.value = value
      },
      get isDisableClick() {
        return DisableClick.value
      },
      set isDisableClick(v) {
        DisableClick.value = v
      },
      checkList: ['取消之前的颜色标记', "不区分大小写", "删除空格后搜索"],
      per: 0,
      start: false,
      //显示搜索进度
      ShowSchedule: false,
      CaseSensitive: true,
      isPbSkip: false,
      mainStyle: "",
      resizeObserver: null,
    }
  },
  methods: {
    handleKeyEnter(e) {
      this.findStart()
    },
    SetFindColor(a) {
      this.$refs.color.style.backgroundColor = a
    },
    findStart() {
      const obj = {
        Options: this.FindConfig.check.join("|"),
        Value: this.$refs.autoEdit.GetValue(),
        Type: this.$refs.FindType.GetType(),
        Range: this.$refs.FindRange.GetValue(),
        Color: this.$refs.FindColor.GetColorObj(),
        ProtoSkip: this.$refs.pbSkip.value,
        action: this.FindConfig.action,
        Theology: -1,
      }
      this.$refs.autoEdit.Set(obj.Type, obj.Value)
      if (obj.action === "hide") {
        Config_Find_Window.value.CompleteCallback()
      }
      if (!Config_Find_Range_ALL.value) {
        try {
          obj.Theology = parseInt(Config_SelectedRow.value["Theology"])
          if (obj.Theology < 1) {
            ElMessage.error("未知错误")
            return
          }
        } catch (e) {
          ElMessage.error("未知错误")
          return
        }
      }
      this.start = true
      this.per = 0
      this.ShowSchedule = true
      DisableClick.value = true

      const isCancelColor = obj.Options.indexOf("取消之前的") !== -1
      FindSession(obj).then((res) => {
        setTimeout(() => {
          this.ShowSchedule = false
        }, 1000)
        Config_Find_Window.value.SearchDone(res, this.FindConfig.action, this.$refs.FindColor.getKeyValue(), isCancelColor)
        DisableClick.value = false
      })
      // window.vm.List.SetEmptySearchMenuVisible(false)
      /*
      CallGoDo("查找", obj).then(res => {
        if (res) {
          if (res.LastSearchResult) {
            for (let i = 0; i < res.LastSearchResult.length; i++) {
              const LastSearchResult = res.LastSearchResult[i]
              const obj = window.vm.List.RowDataHashMap[LastSearchResult]
              if (obj) {
                obj.data.color.search = null
                obj.setData(obj.data)
              }
            }
          }
          let rowIndex = -1;
          if (res.SearchResult) {
            for (let i = 0; i < res.SearchResult.length; i++) {
              const LastSearchResult = res.SearchResult[i]
              const obj = window.vm.List.RowDataHashMap[LastSearchResult]
              if (obj) {
                obj.data.color.search = res.Color
                if (window.Theology === obj.data.Theology) {
                  if (obj.data.color) {
                    obj.data.color.selected = {dark: res.Color, right: res.Color}
                  } else {
                    obj.data.color = {selected: {dark: res.Color, right: res.Color}}
                  }
                  window.vm.List.RefreshRenderedNodes()
                }
                obj.setData(obj.data)
                if (rowIndex === -1) {
                  rowIndex = obj.rowIndex
                }
              }
            }
            if (rowIndex !== -1) {
              console.log(rowIndex)
              window.UI.FindWindow = false
              window.vm.List.ListFollowShow = false
              this.$nextTick(() => {
                window.vm.List.SetEmptySearchMenuVisible(true)
                window.vm.List.agGridApi.ensureIndexVisible(rowIndex, 'middle');
              })
            }
          }
          if (rowIndex === -1) {
            ElMessage({
              message: "没有搜索结果",
              type: 'warning',
            })
          }
        }
        this.$nextTick(() => {
          this.per = 100
          this.start = false
          this.ShowSchedule = false
          window.vm.List.agSelectedLine = null
          window.vm.List.RefreshRenderedNodes()
          this.$nextTick(() => {
            window.vm.List.agSelectedLine = window.vm.List.RowDataHashMap[window.Theology]
          })
        })
      })
      */
    },
    SetPosition(position) {
      this.per = position
    },
    SetFocus() {
      this.$refs.autoEdit.focus()
    },
  },
  mounted() {
    //window.vm.Find = this
    this.$nextTick(() => {
      this.SetFocus()
    })
    //init
    {
      const firstChild = this.$refs.Color.$el.firstElementChild;
      if (firstChild && firstChild.firstElementChild) {
        const target = firstChild.firstElementChild;
        target.style.color = this.$refs.FindColor.GetValue();
        target.style.fontWeight = "bold";
      }
    }
    {
      Events.On("FindSearchProgress", (obj) => {
        let o=obj.data[0];
        if (!o){
          o=obj.data;
        }
        this.per = parseInt(o)
      })
    }
  }
}

</script>
<style>
/* ✅ 让 `span` 填充整个按钮 */
.el-radio-button__inner {
  width: 100%;
  display: flex;
  justify-content: center; /* ✅ 让文字居中 */
}
</style>