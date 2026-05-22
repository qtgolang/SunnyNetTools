<template>
  <el-tag class="ml-2" type="danger">第一步</el-tag>
  <el-tag class="ml-2" type="info">找到导出的证书文件</el-tag>
  <br>
  <img :src="Image1"/>
  <br>
  <br>
  <el-tag class="ml-2" type="danger">第二步</el-tag>
  <el-tag class="ml-2" type="info">双击打开并选择"安装证书"</el-tag>
  <br>
  <img :src="Image2"/>
  <br>  <br>
  <el-tag class="ml-2" type="danger">第三步</el-tag>
  <el-tag class="ml-2" type="info">选择"本地计算机",然后点击"下一步"</el-tag>
  <br>
  <img :src="Image3"/>
  <br>  <br><br>
  <el-tag class="ml-2" type="danger">第四步</el-tag>
  <el-tag class="ml-2" type="danger">按照图片中标记的顺序选择!!</el-tag>
  <el-tag class="ml-2" type="danger">按照图片中标记的顺序选择!!</el-tag>
  <el-tag class="ml-2" type="danger">按照图片中标记的顺序选择!!</el-tag>
  <br><br><br>
  <img :src="Image16"/>
  <br>  <br>
  <el-tag class="ml-2" type="danger">第五步</el-tag>
  <el-tag class="ml-2" type="info">点击"完成"</el-tag>
  <br>
  <img :src="Image17"/>
  <br>  <br>
  <el-tag class="ml-2" type="danger">第六步</el-tag>
  <el-tag class="ml-2" type="info">弹出信息框(安装成功)点击"确定"</el-tag>
  <br>
  <img :src="Image172"/>
  <br>
  <el-tag class="ml-2" type="danger">到此证书安装完成</el-tag>
</template>

<script>
import {GetDocImage} from "../../../image";

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
    Image3(){
      return GetDocImage("Image3")
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
    Image172(){
      return GetDocImage("Image172")
    },
  }
}
</script>