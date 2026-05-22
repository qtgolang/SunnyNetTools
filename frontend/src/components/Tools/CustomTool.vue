<template>
  <div class="icon-wrapper">
    <img :src="Icon" alt="" class="icon"/>
    <el-text class="mx-1" v-show="noName!=='true'">{{ Name }}</el-text>
  </div>
</template>


<script>
import {Plus} from "@element-plus/icons-vue";

export default {
  components: {Plus},
  computed: {
    Icon() {
      return "data:image/png;base64," + this.obj.Icon
    }
  },
  props: ['obj', "noName"],
  data() {
    return {
      Name: ""
    }
  },
  mounted() {
    this.Name = this.obj.Name
  }
}
</script>

<style scoped>
.icon-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: auto;
  min-width: 80px; /* 设置最小宽度为80px */
  cursor: pointer; /* 添加此行 */
  transition: transform 0.2s ease-in-out;
  margin-left: 5px;
  margin-top: 5px;
  margin-bottom: 10px;
  margin-right: 0px;
}

.icon-wrapper:hover {
  transform: scale(1.05);
}

.icon {
  width: 40px;
  height: 40px;
}

.name {
  margin-top: 4px;
}
</style>