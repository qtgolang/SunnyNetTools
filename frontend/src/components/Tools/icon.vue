<template>
  <div class="icon-wrapper">
    <img :src="icon" alt="" class="icon" v-if="this.name !== '导出证书'"/>
    <el-icon v-if="this.name === '导出证书'" style="height: 40px;width: 40px">
      <BottomRight style="width: 50%;height: 50%;"/>
    </el-icon>
    <el-text class="mx-1" :style="getColorStyle" v-show="noName!=='true'">{{ Name }}</el-text>
  </div>
</template>


<script>
import {GetImage2} from "./image.js";
import {BottomRight, Plus} from "@element-plus/icons-vue";

export default {
  components: {BottomRight, Plus},
  computed: {
    getColorStyle() {
      //内置的图标文字颜色设置
      return "";
      //return Config_IsDark.value ? 'color: #15e15a' : "color: #f421a8"
    }
  },
  props: ['name', "noName"],
  watch: {
    name(value) {
      this.icon = GetImage2(value)
      this.Name = value
    },
  },
  data() {
    return {
      icon: '',
      Name: ""
    }
  },
  mounted() {
    if (this.name !== "导出证书") {
      this.icon = GetImage2(this.name)
    }
    this.Name = this.name
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