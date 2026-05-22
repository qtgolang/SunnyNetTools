<script>
import Home from "./Home.vue";
import Tour from "./Home/Tour.vue";

export default {
  props: ["Title", "Show", "CancelShow"],
  components: {Home, Tour},
}
</script>

<template>
  <div>
    <Tour/>
    <Home/>
  </div>
</template>
