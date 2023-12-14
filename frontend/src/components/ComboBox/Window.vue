<template>
  <div class="top-component" :style="WindowStyleZIndex">
    <div class="ag-panel ag-default-panel ag-dialog ag-ltr ag-popup-child" tabindex="-1" role="dialog"
         aria-label="Range Chart"
         :style="WindowStyle" ref="Window">
      <div ref="eContentWrapper" class="ag-panel-content-wrapper ag-default-panel-content-wrapper"
           style="height: 100px;">
        <div class="ag-chart ag-ltr" tabindex="-1">
          <div ref="eChartContainer" tabindex="-1" class="ag-chart-components-wrapper ag-chart-menu-visible">
            <div ref="eEmpty" class="ag-chart-empty-text ag-unselectable" style=" overflow-y: auto;">
              <div style="width: 90%">
                <div style="display: block;cursor: pointer; margin-bottom: 10px;" @click="Call(1)">设置系统IE代理</div>

                <div style="display: block;cursor: pointer; " @click="Call(2)">取消系统IE代理</div>
              </div>

            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: ['size', "func"],
  watch: {
    size(Value) {
      this.NSize = Value
      this.setWindowSize(Value.top, Value.left, Value.width, Value.height)
    },
    func(Value) {
      this.CallFunction = Value
    }
  },
  components: {},
  computed: {
    WindowStyleZIndex() {
      return "z-index: 1100000;"
    }
    ,
  },
  data() {
    return {
      WindowStyle: "", NSize: {}, CallFunction: null
    }
  },
  mounted() {
    this.setWindowSize(30, 132, 200, 200)
  },
  beforeUnmount() {

  },
  methods: {
    setWindowSize(top, left, width, height) {
      let t = "top: " + top + "px; left: " + left + "px; width: " + width + "px; max-width: " + width + "px; min-width: " + width + "px; height: " + height + "px; max-height: " + height + "px; min-height: " + height + "px;"
      t = t.replaceAll("pxpx;", "px;")
      this.WindowStyle = t
    }, Call(index) {
      if (this.CallFunction !== null) {
        this.CallFunction(index)
      }
    }
  }
}
</script>
<style scoped>
.top-component {
  position: absolute;
  top: 0;
  left: 0;
}
</style>