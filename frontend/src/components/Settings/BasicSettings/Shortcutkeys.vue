<template>
  <el-table :data="tableData" style="width: 100%" :class="IsDark">
    <el-table-column label="快捷键设置" width="280">
      <template #default="scope">
        <span style="margin-left: 10px">{{ scope.row.Name }}</span>
      </template>
    </el-table-column>
    <el-table-column>
      <template #default="scope">
        <el-input v-model="scope.row.value" placeholder="请按下快捷键" readonly
                  @keydown="handleKeyDown(scope.row,$event)"/>
      </template>
    </el-table-column>
  </el-table>
</template>

<script lang="ts" setup>

import {computed, ref} from "vue";

const IsDark = computed(() => {
  if (window.Theme.IsDark) {
    return "table_drak"
  }
  return "table_ffff";
});

interface User {
  Name: string
  ctrlKey: boolean
  altKey: boolean
  shiftKey: boolean
  metaKey: boolean
  key: string
  value: string
}

const tableData: User[] = [
  {
    Name: '搜索/查找',
    ctrlKey: true,
    altKey: false,
    shiftKey: false,
    metaKey: false,
    key: "F",
    value: "CTRL + F",
  },
  {
    Name: '全部放行',
    ctrlKey: false,
    altKey: false,
    shiftKey: true,
    metaKey: false,
    key: "Z",
    value: "Shift + Z",
  },
  {
    Name: '放行当前请求',
    ctrlKey: true,
    altKey: false,
    shiftKey: false,
    metaKey: false,
    key: "Z",
    value: "CTRL + Z",
  },
  {
    Name: '设置/取消IE代理',
    ctrlKey: false,
    altKey: false,
    shiftKey: false,
    metaKey: false,
    key: "F12",
    value: "F12",
  },
  {
    Name: '清空全部记录',
    ctrlKey: true,
    altKey: false,
    shiftKey: false,
    metaKey: false,
    key: "X",
    value: "CTRL + X",
  },
  {
    Name: '复制',
    ctrlKey: true,
    altKey: false,
    shiftKey: false,
    metaKey: false,
    key: "C",
    value: "CTRL + C",
  },
]
const tableDataMap = {}
tableDataMap["搜索/查找"] = tableData[0]
tableDataMap["全部放行"] = tableData[1]
tableDataMap["放行当前请求"] = tableData[2]
tableDataMap["设置/取消IE代理"] = tableData[3]
tableDataMap["清空全部记录"] = tableData[4]
tableDataMap["复制"] = tableData[5]
window.KeysStrings = tableDataMap
const handleKeyDown = (row: User, event: KeyboardEvent) => {
  event.stopPropagation()
  const mKey = event.key.toUpperCase()
  if ("CONTROL" === mKey) {
    return
  }
  if ("BACKSPACE" == mKey) {
    row.altKey = false
    row.ctrlKey = false
    row.shiftKey = false
    row.key = ""
    row.value = ""
    return;
  }
  row.altKey = event.altKey
  row.ctrlKey = event.ctrlKey
  row.shiftKey = event.shiftKey
  row.key = event.key.toUpperCase()
  let obj = ""
  {
    if (row.ctrlKey) {
      obj = "CTRL"
    }
    if (row.altKey) {
      if (obj === "") {
        obj = "ALT"
      } else {
        obj += " + ALT"
      }
    }
    if (row.shiftKey) {
      if (obj === "") {
        obj = "Shift"
      } else {
        obj += " + Shift"
      }
    }

    if (row.key) {
      if (obj === "") {
        obj = row.key
      } else {
        obj += " + " + row.key
      }
    }
  }
  row.value = obj
  // 在这里可以访问到 scope.row 的值
}

defineExpose({
  tableDataMap
});
</script>

<style>

.table_drak {
  --el-table-tr-bg-color: #2D3436 !important;
}

.table_ffff {
  --el-table-tr-bg-color: #ffffff !important;
}
</style>