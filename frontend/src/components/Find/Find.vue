<template>

  <div style="width: 100%;height: 100%;position: absolute;">
    <div style="justify-items: center;top:30px;position: absolute;left: 20px;width: 500px">
      查找：
      <Autocomplete ref="autoEdit" style="width: 200px"/>
      &nbsp;
      查找类型：
      <FindType ref="FindType" style="width: 130px"/>
    </div>
    <div style="justify-items: center;top:70px;position: absolute;left: 90px;width: 500px">
      <el-checkbox-group v-model="checkList" style="text-align: left;">
        <el-checkbox label="不区分大小写" v-show="CaseSensitive"/>
        <el-checkbox label="删除空格后搜索"/>
      </el-checkbox-group>
    </div>
    <div style="justify-items: center;top:110px;position: absolute;left: 0px;width: 500px">
      <span style="position: absolute;left: 25px;top:4px">查找范围：</span>
      <FindRange ref="FindRange" style="position: absolute;left: 90px;width: 400px"/>
    </div>
    <div style="justify-items: center;top:150px;position: absolute;left: 0px;width: 500px">
      <span style="position: absolute;left: 25px;top:4px;background-color: #ff0000;color: #0a0a0a"
            ref="color">颜色标记：</span>
      <FindColor ref="FindColor" style="position: absolute;left: 90px;width: 199px"/>
      <el-checkbox-group v-model="checkList" style="position: absolute;left: 272px;width: 199px;top:-4px">
        <el-checkbox label="取消之前的颜色标记"/>
      </el-checkbox-group>
    </div>
    <div style="justify-items: center;top:75px;position: absolute;left: 165px;width: 500px" v-show="isPbSkip">
      <span>忽略前: </span>
      <input type="number" style="width: 30px" value="0" min="0" ref="pbSkip"/>
      <span> 字节 </span>
    </div>
    <el-button size="large"
               style="top:190px;position: absolute; margin: 0 auto;left: 90px;width: 200px;height: 20px;"
               @click="findStart" v-show="start===false">查找
    </el-button>
    <el-button size="large"
               style="top:190px;position: absolute; margin: 0 auto;left: 90px;width: 200px;height: 20px;"
               v-show="start" disabled>查找中...
    </el-button>
    <el-progress v-show="isShowSchedule" :percentage="per" color="#13EA00"
                 style="top:230px;left: 90px;width: 450px;"/>
  </div>
</template>

<script>
import Autocomplete from "./Autocomplete.vue";
import FindType from "./FindType.vue";
import FindRange from "./FindRange.vue";
import FindColor from "./FindColor.vue";
import {CallGoDo} from "../CallbackEventsOn.js";
import {ElMessage} from "element-plus";

export default {
  computed: {
    isShowSchedule() {
      return this.ShowSchedule
    }
  },
  components: {FindColor, FindRange, FindType, Autocomplete},
  data() {
    return {
      checkList: ['取消之前的颜色标记'],
      per: 100,
      start: false,
      //显示搜索进度
      ShowSchedule: false,
      CaseSensitive: true,
      isPbSkip: false
    }
  },
  methods: {
    SetFindType(a) {
      this.isPbSkip = false
      if (a === "UTF8" || a === "GBK") {
        this.CaseSensitive = true
      } else {
        if (a === "pb") {
          this.isPbSkip = true
        }
        this.CaseSensitive = false
        this.checkList = this.checkList.filter(item => item.toLowerCase() !== '不区分大小写');
      }
      this.$refs.autoEdit.SetAutoType(a)
    },
    SetFindColor(a) {
      this.$refs.color.style.backgroundColor = a
    },
    findStart() {
      const obj = {
        Options: this.checkList.join("|"),
        Value: this.$refs.autoEdit.GetValue(),
        Type: this.$refs.FindType.GetType(),
        Range: this.$refs.FindRange.GetValue(),
        Color: this.$refs.FindColor.GetValue(),
        ProtoSkip: this.$refs.pbSkip.value
      }
      this.start = true
      this.per = 0
      this.ShowSchedule = true
      this.$refs.autoEdit.Set(obj.Type, obj.Value)
      window.vm.List.SetEmptySearchMenuVisible(false)
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
    },
    SetPosition(position) {
      this.per = position
    },
    SetFocus() {
      this.$refs.autoEdit.focus()
    }
  },
  mounted() {
    window.vm.Find = this
    this.$nextTick(() => {
      this.SetFocus()
    })
  }
}

</script>