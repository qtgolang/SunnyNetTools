<template>

  <div style="width: 50%;">
    <div style="height: 270px;">
      <ProcessList ref="ProcessName"/>
    </div>
    <div :style="getAnyProcessInputStyle1">
      <el-input v-model="AddName" placeholder="要添加的进程名称,例如: qq.exe">
        <template #append>
          <el-tooltip class="item" effect="dark" content="确定添加" placement="top">
            <el-icon style="cursor: pointer" @click="addProcessName">
              <Plus/>
            </el-icon>
          </el-tooltip>
        </template>
      </el-input>
    </div>
    <div :style="getAnyProcessInputStyle2">
      <el-input placeholder="要添加的进程名称,例如: qq.exe">
        <template #append>
          <el-tooltip class="item" effect="dark" content="已经开启捕获任意进程,无法添加" placement="top">
            <el-icon>
              <Plus/>
            </el-icon>
          </el-tooltip>
        </template>
      </el-input>
    </div>
    <div :style="getAnyProcessInputStyle3">
      <el-input placeholder="正在添加,请稍等...">
        <template #append>
          <el-tooltip class="item" effect="dark" content="正在添加请稍等" placement="top">
            <el-icon>
              <Plus/>
            </el-icon>
          </el-tooltip>
        </template>
      </el-input>
    </div>

  </div>
  <div style="width: 50%;">
    <div style="height: 270px;">
      <ProcessList2 ref="ProcessPid"/>
    </div>
    <el-checkbox v-model="anyProcess">捕获任意进程</el-checkbox>
  </div>
</template>


<script>

import ProcessList2 from "./processList2.vue";
import ProcessList from "./processList.vue";
import {CallGoDo} from "../../CallbackEventsOn.js";
import {ElMessage} from "element-plus";

export default {
  components: {ProcessList, ProcessList2},
  watch: {
    "anyProcess"(newValue, _) {
      if (newValue) {
        this.$refs.ProcessPid.RowDataHashMap = {}
        this.$refs.ProcessPid.RowData = []
        this.$refs.ProcessPid.agGridApi.setRowData(this.$refs.ProcessPid.RowData);
        this.$refs.ProcessPid.overlayNoRowsTemplate = `<span style="padding: 20px;" id="HookMessageText">已经开启捕获任意进程</span>`
        this.$refs.ProcessName.overlayNoRowsTemplate = `<span style="padding: 20px;" id="HookMessageText">已经开启捕获任意进程</span>`
        this.$refs.ProcessName.Empty()
      } else {
        this.$refs.ProcessPid.overlayNoRowsTemplate = `<span style="padding: 20px;" id="HookMessageText">无内容</span>`
        this.$refs.ProcessName.overlayNoRowsTemplate = `<span style="padding: 20px;" id="HookMessageText">未添加任何进程名</span>`
        this.$refs.ProcessName.Empty()
        this.$refs.ProcessPid.RowDataHashMap = {}
        this.$refs.ProcessPid.RowData = []
        this.$refs.ProcessPid.agGridApi.setRowData(this.$refs.ProcessPid.RowData);
      }
      this.$refs.ProcessPid.agGridApi.setFilterModel(null); // 清空过滤器条件
      this.$refs.ProcessName.MenuItems[2].visible = !newValue
      this.$refs.ProcessName.MenuItems[3].visible = !newValue
      this.$nextTick(() => {
        CallGoDo("进程驱动添加进程名", {Name: "{OpenALL}", isSet: newValue})
      })
    },
  },
  computed: {
    getAnyProcessInputStyle1() {
      return this.anyProcess ? "display:none" : this.Adding ? "display:none" : ""
    },
    getAnyProcessInputStyle2() {
      return !this.anyProcess ? "display:none" : ""
    },
    getAnyProcessInputStyle3() {
      return this.Adding === false ? "display:none" : this.anyProcess ? "display:none" : ""
    }
  },
  data() {
    return {
      //要添加的进程名
      AddName: "",
      //是否正在枚举进程
      isEnumerateProcesses: false,

      //捕获任意进程
      anyProcess: false,
      Adding: false
    }

  },
  mounted() {
    setInterval(() => {
      if (!this.isEnumerateProcesses) {
        if (window.UI.Settings && window.vm.Settings.Title.indexOf("进程拦截") !== -1 && !this.anyProcess && window.vm.Settings.LoadDrive) {
          this.isEnumerateProcesses = true
          this.EnumerateProcesses()
          this.isEnumerateProcesses = false
        }
      }
    }, 1000)
  },
  methods: {
    addProcessName() {
      const pName = this.AddName.trim().toLowerCase()
      if (pName === '') {
        ElMessage({
          message: "请输入要添加的进程名",
          type: 'error',
        })
        return
      }
      this.Adding = true
      CallGoDo("进程驱动添加进程名", {Name: pName, isSet: true}).then(res => {
        this.$refs.ProcessName.AddLine(pName)
        this.Adding = false
        this.AddName = ""
      })
    },
    EnumerateProcesses() {
      CallGoDo("枚举进程", null).then(res => {
            if (res) {
              let update = []
              let addList = []
              let Delete = []
              for (let pid in this.$refs.ProcessPid.RowDataHashMap) {
                if (!res.hasOwnProperty(pid)) {
                  //如果列表中的在新的数组中不存在则删除
                  Delete.push(this.$refs.ProcessPid.RowDataHashMap[pid].data)
                  delete this.$refs.ProcessPid.RowDataHashMap[pid]
                } else {
                  //如果列表中的PID在新的数组中存在,则对比进程名是否变化,如果变化,则更新
                  const ps = this.$refs.ProcessPid.RowDataHashMap[pid].data
                  if (ps['进程名'] !== res[pid]) {
                    update.push({PID: pid, '进程名': res[pid]})
                  }
                }
              }
              this.$refs.ProcessPid.Delete(Delete)
              for (let pid in res) {
                if (!this.$refs.ProcessPid.RowDataHashMap.hasOwnProperty(pid)) {
                  //如果新数组中的PID,在列表中不存在，则新增
                  if (pid !== 0 && pid !== '0') {
                    addList.push({PID: pid, '进程名': res[pid]})
                  }
                }
              }
              this.$refs.ProcessPid.AddLines(addList)
              this.$refs.ProcessPid.Update(update)
            }
          }
      )
    }
  }
}

</script>