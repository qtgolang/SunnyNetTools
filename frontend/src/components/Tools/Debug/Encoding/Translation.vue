<script>

export default { 
  computed: {},
  data() {
    return {
      DragFilePath: "",
      mode: "baidu",
    }
  },
  watch: {

  },
  methods: {
    empty() {
      this.$refs.Data.SetCode("")
      this.$refs.Res.SetCode("")
    },
  },
  mounted() {
    //this.$refs.Data.SetLanguage("text")
  },
};
</script>

<template>
  <div style=" width: 100%; height:100%;overflow:hidden">
    <div
        style=" width: 100%; height:60px;overflow:hidden; display: flex;justify-content: center;align-items: center;margin-right: 20px">
      <el-radio-group v-model="mode" class="ml-4">
        <el-radio value="baidu" size="large">百度翻译</el-radio>
        <el-radio value="sogou" size="large">搜狗翻译</el-radio>
        <el-radio value="youdao" size="large">有道翻译</el-radio>
      </el-radio-group>
    </div>
    <div style="width: calc(100% - 5px);height:calc(100% - 62px); margin-left: 2px;margin-right: 2px;margin-bottom: 5px">
      <iframe v-if="mode==='baidu'" src="https://fanyi.baidu.com/mtpe-individual/multimodal#/" width="100%"
              height="100%"
              style=" overflow-y: auto;"/>
      <iframe v-if="mode==='sogou'"
              src="https://fanyi.sogou.com/text?keyword=&transfrom=auto&transto=zh-CHS&model=general" width="100%"
              height="100%"
              style=" overflow-y: auto;"/>
      <iframe v-if="mode==='youdao'"
              src="https://fanyi.youdao.com/#/TextTranslate" width="100%"
              height="100%"
              style=" overflow-y: auto;"/>
    </div>

  </div>

</template>
<style>
.el-form-item {
  margin-bottom: 4px
}

.el-cascader-menu__wrap.el-scrollbar__wrap {
  height: 220px !important;
}
</style>