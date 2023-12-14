<script>
import {GetImage} from "./image.js";

export default {
  setup(props) {
    let imageBase64 = "";
    const mode = props.params.data["断点模式"]
    if (mode === 1) {
      imageBase64 = GetImage("拦截上行")
    } else if (mode === 2) {
      imageBase64 = GetImage("拦截下行")
    } else {
      imageBase64 = GetImage(props.params.data["ico"])
    }
    return {
      imageBase64,
      value: props.params.value
    }
  }
};
</script>
<template>
  <div style="cursor: pointer; display: flex; align-items: center;">
    <img :src="imageBase64" style="height: 22px;width: 22px"/>&nbsp;{{ value }}
  </div>
</template>>