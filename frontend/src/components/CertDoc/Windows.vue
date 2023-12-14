<template>
  <el-tag class="ml-2" type="danger">第一步</el-tag>
  <el-tag class="ml-2" type="info">在找到软件界面"设置"</el-tag>
  <br>
  <img :src="Image1"/>
  <br>
  <br>
  <el-tag class="ml-2" type="danger">第二步</el-tag>
  <el-tag class="ml-2" type="info">在"设置"中找到"SSL证书"</el-tag>
  <br>
  <img :src="Image2"/>
  <br>  <br>
  <el-tag class="ml-2" type="danger">第三步</el-tag>
  <el-tag class="ml-2" type="info">点击"安装默认证书"</el-tag>
  <br>
  <img :src="Image15"/>
  <br>  <br>
  <el-tag class="ml-2" type="danger">第四步</el-tag>
  <el-tag class="ml-2" type="info">弹出信息"安装结果"（成功安装）</el-tag>
  <br>
  <img :src="Image16"/>
  <br>  <br>
  <el-tag class="ml-2" type="danger">或者</el-tag>
  <el-tag class="ml-2" type="info">弹出信息"已在储存中"（成功安装）</el-tag>
  <br>
  <img :src="Image17"/>
  <br>
  <el-tag class="ml-2" type="danger">到此证书安装完成</el-tag>
</template>

<script>
import {GetDocImage} from "../image.js";

export default {
  data() {
    return {}
  },
  computed: {
    Image1() {
      return GetDocImage("Image1")
    },
    Image2(){
      return GetDocImage("Image2")
    },
    Image15(){
      return GetDocImage("Image15")
    },
    Image16(){
      return GetDocImage("Image16")
    },
    Image17(){
      return GetDocImage("Image17")
    },
  }
}
</script>