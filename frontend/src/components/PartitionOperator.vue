<template>
  &nbsp;
  <svg width="20" height="24" xmlns="http://www.w3.org/2000/svg" stroke-linejoin="round" stroke-linecap="round"
       stroke-width="2" stroke="#000" fill="none" :class="cvgStyle">
    <g>
      <line transform="rotate(90 10.1132 12.9434)" id="svg_1" y2="12.9434" x2="17.11321" y1="12.9434" x1="3.11321"/>
    </g>
  </svg>
  &nbsp;
</template>
<script>

export default {
  components: {},
  computed: {
    cvgStyle() {
      if (!this.theme) {
        return "dr2"
      } else {
        return "dr1"
      }
    },
  },
  data() {
    return {
      get theme() {
        return window.Theme.IsDark
      },
      set theme(newValue) {
        window.Theme.IsDark = newValue
      }
    }
  }, methods: {}
}
</script>
<style scoped>
/*深色模式显示*/
.dr1 {
  fill: #5d5c5c;
  stroke: #5d5c5c
}

/*浅色模式显示*/
.dr2 {
  fill: #d0cfcf;
  stroke: #d0cfcf
}
</style>