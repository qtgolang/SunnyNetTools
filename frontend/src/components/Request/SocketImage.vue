<script>
import {GetImage} from ".././image.js";

export default {
  setup(props) {
    let imageBase64 =GetImage(props.params.data["ico"]);
    return {
      imageBase64,
      value: props.params.value
    }
  }
};
</script>
<template>
  <div style="cursor: pointer; display: flex; align-items: center;"><img :src="imageBase64"
                                                                         style="height: 22px;width: 22px"/>&nbsp;{{
      value
    }}
  </div>
</template>>