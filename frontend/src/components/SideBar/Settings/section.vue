<template>
  <div class="ag-chart-format-section">
    <div class="ag-group ag-charts-format-top-level-group ag-group-item-alignment-center" data-ref="axisGroup">
      <div
          class="ag-group-title-bar ag-charts-format-top-level-group-title-bar ag-unselectable"
          style="cursor: pointer"
          @click="onHeaderClick"
      >
        <span class="ag-group-title-bar-icon ag-charts-format-top-level-group-title-bar-icon">
          <span :class="iconClass"></span>
        </span>

        <span
            v-if="!useButtonTitle"
            class="ag-group-title ag-charts-format-top-level-group-title"
            style="cursor: pointer"
        >
          {{ title }}
        </span>

        <slot v-else name="title" />
      </div>

      <div v-show="modelValue" :style="bodyStyle">
        <slot />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "FormatSection",
  props: {
    title: { type: String, default: "" },
    modelValue: { type: Boolean, default: false },

    // 是否可折叠：true=点击切换展开/收起；false=点击只触发 toggle 事件
    collapsible: { type: Boolean, default: true },

    // true 时用 slot(title) 渲染标题（例如 el-button）
    useButtonTitle: { type: Boolean, default: false },

    // 内容区域样式
    bodyStyle: {
      type: Object,
      default: () => ({ height: "auto", width: "100%" }),
    },
  },
  emits: ["update:modelValue", "toggle"],
  computed: {
    iconClass() {
      // 非折叠入口固定显示 closed 图标
      if (!this.collapsible) return "ag-icon ag-icon ag-icon-tree-closed";
      return this.modelValue ? "ag-icon ag-icon ag-icon-tree-open" : "ag-icon ag-icon ag-icon-tree-closed";
    },
  },
  methods: {
    onHeaderClick() {
      // 非折叠入口：不改值，只抛事件
      if (!this.collapsible) {
        this.$emit("toggle");
        return;
      }

      const next = !this.modelValue;
      this.$emit("update:modelValue", next);
      this.$emit("toggle", next);
    },
  },
};
</script>
