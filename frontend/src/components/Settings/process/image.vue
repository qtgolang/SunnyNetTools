<script>
//

import {CallGoDo} from "../../CallbackEventsOn.js";

export default {
  props: ['params'],
  data() {
    return {
      checkList: false,
      PID: "null",
    }
  },
  created() {
    this.PID = "" + this.params.data["PID"];
    this.checkList = this.params.data["box"] === 1;
  },
  mounted() {

  }, methods: {
    onCheckChange(value) {
      this.params.data["box"] = value ? 1 : 0;

      console.log("onCheckChange", value, this.PID);
      CallGoDo("进程驱动添加PID", {PID: this.PID, isSelected: value})

    }
  }
};
</script>
<template>
  <el-checkbox v-model="checkList" @change="onCheckChange">{{ PID }}</el-checkbox>
</template>