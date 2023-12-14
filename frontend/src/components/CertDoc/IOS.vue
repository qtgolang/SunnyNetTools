<template>
  <el-tag class="ml-2" type="danger">第一步</el-tag>
  <el-tag class="ml-2" type="info">确认手机与电脑在同一局域网</el-tag>
  <br>
  <el-tag class="ml-2" type="danger">第二步</el-tag>
  <el-tag class="ml-2" type="info">确认电脑已关闭防火墙</el-tag>
  <br>
  <img :src="Image18"/>
  <br>  <br>
  <el-tag class="ml-2" type="danger">第三步</el-tag>
  <el-tag class="ml-2" type="info">手机上打开"Safari浏览器"</el-tag>
  <br>
  <img :src="Image19"/>
  <br>  <br>
  <el-tag class="ml-2" type="danger">第四步</el-tag>
  <el-tag class="ml-2" type="info">在浏览器输入"电脑的IP"+本软件运行的端口号</el-tag>
  <br>
  <img :src="Image20"/>
  <br>  <br>
  <el-tag class="ml-2" type="danger">第五步</el-tag>
  <el-tag class="ml-2" type="info">点击下载SunnyRoot证书</el-tag>
  <br>
  <img :src="Image21"/>
  <br> <br>
  <el-tag class="ml-2" type="danger">第六步</el-tag>
  <el-tag class="ml-2" type="info">弹出"正尝试下载描述文件",点击"允许"</el-tag>
  <br>
  <img :src="Image22"/>
  <br> <br>
  <el-tag class="ml-2" type="danger">第七步</el-tag>
  <el-tag class="ml-2" type="info">下载完成，点击“关闭”</el-tag>
  <br>
  <img :src="Image23"/>
  <br> <br>
  <el-tag class="ml-2" type="danger">第八步</el-tag>
  <el-tag class="ml-2" type="info">打开"设置"->"通用"->"VPN与设备管理" 点击"SunnyNet"</el-tag>
  <br>
  <img :src="Image24"/>
  <br> <br>
  <el-tag class="ml-2" type="danger">第九步</el-tag>
  <el-tag class="ml-2" type="info">点击"安装",直到安装完成</el-tag>
  <br>
  <img :src="Image25"/>
  <br> <br>
  <el-tag class="ml-2" type="danger">第十步</el-tag>
  <el-tag class="ml-2" type="info">打开"设置"->"通用"->"关于本机" 最下方点击"证书信任设置"</el-tag>
  <br> <br>
  <el-tag class="ml-2" type="danger">第十一步</el-tag>
  <el-tag class="ml-2" type="info">找到"SunnyNet"</el-tag>
  <br>
  <img :src="Image26"/>
  <br> <br>
  <el-tag class="ml-2" type="danger">第十二步</el-tag>
  <el-tag class="ml-2" type="info">打开信任</el-tag>
  <br>
  <img :src="Image27"/>
  <br><el-tag class="ml-2" type="danger">到此证书安装完成</el-tag>
</template>

<script>
import {GetDocImage} from "../image.js";

export default {
  data() {
    return {}
  },
  computed: {
    Image18() {
      return GetDocImage("Image18")
    },
    Image19() {
      return GetDocImage("Image19")
    },
    Image20() {
      return GetDocImage("Image20")
    },
    Image21() {
      return GetDocImage("Image21")
    },
    Image22() {
      return GetDocImage("Image22")
    },
    Image23() {
      return GetDocImage("Image23")
    },
    Image24() {
      return GetDocImage("Image24")
    },
    Image25() {
      return GetDocImage("Image25")
    },
    Image26() {
      return GetDocImage("Image26")
    },
    Image27() {
      return GetDocImage("Image27")
    },
  }
}
</script>