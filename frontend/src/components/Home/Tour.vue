<script xmlns="">
import {IsOpenTour, TourList} from "./Tour";

export default {
  data() {
    return {
      get open() {
        return IsOpenTour.value
      },
      set open(v) {
        IsOpenTour.value = v
      },
    }
  },
  computed: {
    getTourList() {
      return TourList.value
    },
    IsOpen() {
      return IsOpenTour.value
    }
  },
  mounted() {
  }
}
</script>

<template>
  <div v-if="IsOpen" >
    <el-tour v-model="open" :close-on-press-escape="false" :show-close="true" style="--wails-draggable: drag">
      <el-tour-step
          v-for="(step, index) in getTourList"
          :key="index"
          :target="step.target"
          :title="step.title"
          type="primary" :mask="true"
      >
        <div v-html="step.description.replace(/\n/g, '<br>')"></div>
      </el-tour-step>
    </el-tour>
  </div>
</template>
