<template>
  <div class="ag-tabs ag-chart-tabbed-menu ag-focus-managed">
    <div  ref="BodyRect" class="ag-tabs-body ag-chart-tabbed-menu-body" style="">

    </div>
  </div>
</template>
<script>

export default {
  props: ['Theology', 'RequestWay'],
  watch: {
    Theology(value) {
      this.RequestTheology = value

    },
    RequestWay(value) {
      this.DisplayHTTPHeaderWebSocket = value.DisplayHTTPHeaderWebSocket
      this.DisplayHTTPHeader = value.DisplayHTTPHeader
    },
  },
  computed: {
  },
  data() {
    return {
    };
  },
  methods: {

  },
  mounted() {

  }
}
</script>