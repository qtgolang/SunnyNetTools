<script>
import {Config_IsRest} from "../../config/Config";
import {ElNotification} from "element-plus";
import {SetKeys} from "../../../../bindings/changeme/Service/appmain";
import {Config_Keys, Config_Keys_ID, GetKeysArray, keydownEventToString, RestKeys} from "../../config/Keys";

export default {
  data() {
    return {
      set isRest(val) {
        Config_IsRest.value = val;
      },
      get isRest() {
        return Config_IsRest.value;
      },
      set keysID(val) {
        Config_Keys_ID.value = val
      },
      get keysID() {
        this.thisId = {}
        for (let i = 0; i < Config_Keys_ID.value.length; i++) {
          const obj = Config_Keys.get(Config_Keys_ID.value[i])
          if (obj) {
            this.thisId[Config_Keys_ID.value[i]] = obj.value
          } else {
            this.thisId[Config_Keys_ID.value[i]] = ""
          }
        }
        return Config_Keys_ID.value
      },
      thisId: {}
    }
  },
  watch: {
    isRest() {
      RestKeys()
    }
  },
  methods: {
    init() {

    },
    check(row, key) {
      const array = GetKeysArray()
      for (let i = 0; i < array.length; i++) {
        if (array[i].Name === row.Name) {
          continue
        }
        if (array[i].value === key) {
          return array[i].Name
        }
      }
      return "";
    },
    handleKeyDown(id, event) {
      window.isEditKeyDown.value = true;
      try {
        const row = Config_Keys.get(id)
        if (!row) {
          return
        }
        row.start = true
        event.stopPropagation()
        event.returnValue = false;
        event.preventDefault();
        event.stopImmediatePropagation();
        const obj = keydownEventToString(event);
        row.altKey = obj.altKey
        row.ctrlKey = obj.ctrlKey
        row.shiftKey = obj.shiftKey
        row.key = obj.key
        row.keyCode = obj.keyCode
        row.value = obj.value
        this.thisId[id] = row.value
      } catch (err) {
      }
    },
    handleKeyUp(id, event) {
      const row = Config_Keys.get(id)
      if (!row) {
        window.isEditKeyDown.value = false;
        return
      }
      event.stopPropagation()
      event.returnValue = false;
      event.preventDefault();
      event.stopImmediatePropagation();
      if (row.start === true) {
        row.start = false
      } else {
        row.start = false
        window.isEditKeyDown.value = false;
        return;
      }
      const k = row.value === '' ? "" : this.check(row, row.value);
      if (k === "") {
        this.Save(row.Name, row.value.length < 1)
        window.isEditKeyDown.value = false;
        return;
      }
      row.altKey = false
      row.ctrlKey = false
      row.shiftKey = false
      row.key = ""
      row.value = ""
      row.keyCode = 0
      this.thisId[id] = ""
      this.Save(row.Name, true)
      window.isEditKeyDown.value = false;
    },
    Save(Name, Empty) {
      for (let i = 0; i < Config_Keys_ID.value.length; i++) {
        const obj = Config_Keys.get(Config_Keys_ID.value[i])
        if (obj) {
          obj.value = this.thisId[Config_Keys_ID.value[i]]
        }
      }
      const message = "快捷键 [ " + Name + " ] " + (Empty ? "已清空\n请重新输入" : "已设置")
      SetKeys(JSON.stringify(GetKeysArray())).then(() => {
        ElNotification({
          showClose: true,
          message: message,
          type: Empty ? 'error' : "success",
          position: 'bottom-right',
          customClass: 'multiline-message'
        })
      })
    },
    GetLimitName(id) {
      const obj = Config_Keys.get(id)
      if (obj) {
        return obj.Name
      }
      return "未知"
    },
  },
  computed: {},
  mounted() {
    this.init()
    const arr = this.$el.getElementsByClassName("el-table__header-wrapper")
    if (arr.length > 0) {
      arr[0].style.display = "none"
    }
  }
}
</script>

<template>
  <div>
    <el-table :data="keysID" style="background-color: #f0f0f0;">
      <el-table-column label="" width="180">
        <template #default="scope">
          <span style="margin-left: 10px;">{{ GetLimitName(scope.row) }}</span>
        </template>
      </el-table-column>
      <el-table-column>
        <template #default="scope">
          <el-input v-model="thisId[scope.row]" placeholder="请按下快捷键" readonly
                    @keydown="handleKeyDown(scope.row,$event)" @keyup="handleKeyUp(scope.row, $event)"
                    @blur="handleKeyUp(scope.row, $event)" @focus="handleKeyDown(scope.row,$event)"/>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
