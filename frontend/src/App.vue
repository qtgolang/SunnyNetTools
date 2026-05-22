<script>
import Home from "./components/Main.vue";

export default {
  components: {Home},
  data() {
    return {
      IsDark: true,
    }
  },
  mounted() {

  }
}
</script>

<template>
  <div class="fullscreen-div" ref="app">
    <Home style="width: 100%;height: 100%"/>
  </div>
</template>

<style>
.fullscreen-div {
  position: fixed; /* 让 div 相对于视口固定 */
  top: 0;
  left: 0;
  width: 100vw; /* 100% 视口宽度 */
  height: 100vh; /* 100% 视口高度 */
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.ag-side-bar {
  max-width: 100% !important;
}

.el-table {
  --el-table-tr-bg-color: var(--ag-background-color) !important;
  --el-table-text-color: var(--ag-cell-text-color) !important;
}

.home-root {
  display: flex;
  align-items: center; /* 垂直居中 */
  justify-content: flex-start; /* 水平靠左 */
  background-color: var(--ag-background-color);
  color: var(--ag-cell-text-color);
}
html.dark {
  --el-fill-color-blank: transparent !important;
}

</style>