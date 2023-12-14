<template>
  <el-radio-group v-model="radio">
    <el-tooltip
        class="box-item"
        effect="dark"
        content="（方法一：手动安装）建议尝试方法二"
        placement="bottom"
    >
      <el-radio-button label="方法一"/>
    </el-tooltip>

    <el-tooltip
        class="box-item"
        effect="dark"
        content="方法二：在终端中输入 命令行 安装"
        placement="bottom"
    >
      <el-radio-button label="方法二"/>
    </el-tooltip>

  </el-radio-group>
  <div v-if="radio!==''">
    <br> <br>
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
    <br> <br>
    <el-tag class="ml-2" type="danger">第三步</el-tag>
    <el-tag class="ml-2" type="info">将默认证书导出到桌面</el-tag>
    <br>
    <img :src="Image3"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第四步</el-tag>
    <el-tag class="ml-2" type="info">在桌面找到导出的证书文件</el-tag>
    <br>
    <img :src="Image4"/>
  </div>
  <div v-if="radio==='方法一'">
    <br> <br>
    <el-tag class="ml-2" type="danger">第五步</el-tag>
    <el-tag class="ml-2" type="info">双击打开证书文件</el-tag>
    <br>
    <el-tag class="ml-2" type="danger">第六步</el-tag>
    <el-tag class="ml-2" type="info">在弹出的信息框中输入密码,并且点击"修改密钥串"</el-tag>
    <br>
    <img :src="Image5"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第七步</el-tag>
    <el-tag class="ml-2" type="info">在"访达"中找到"应用程序"</el-tag>
    <br>
    <img :src="Image6"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第八步</el-tag>
    <el-tag class="ml-2" type="info">在"应用程序"中找到"实用工具"</el-tag>
    <br>
    <img :src="Image7"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第九步</el-tag>
    <el-tag class="ml-2" type="info">在"实用工具"中找到"钥匙串访问",并且双击打开</el-tag>
    <br>
    <img :src="Image8"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第十步</el-tag>
    <el-tag class="ml-2" type="info">在"钥匙串访问"中左侧找到"系统",右侧找到"证书"</el-tag>
    <br>
    <el-tag class="ml-2" type="danger">如果你在这一步中找不到"SunnyNet"证书,请尝试“方法二”</el-tag>
    <br>
    <img :src="Image9"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第十一步</el-tag>
    <el-tag class="ml-2" type="info">右侧下方找到"SunnyNet"并且右键点击选择"显示简介"</el-tag>
    <br>
    <img :src="Image10"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第十二步</el-tag>
    <el-tag class="ml-2" type="info">在弹出的窗口中选择"信任"</el-tag>
    <br>
    <img :src="Image11"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第十三步</el-tag>
    <el-tag class="ml-2" type="info">将"使用此证书时"改为"始终信任"后,关闭窗口</el-tag>
    <br>
    <img :src="Image12"/>
    <br>
    <img :src="Image13"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第十四步</el-tag>
    <el-tag class="ml-2" type="info">关闭窗口时要求输入密码并且点击"更新设置"</el-tag>
    <br>
    <img :src="Image14"/>

  </div>
  <div v-if="radio==='方法二'">
    <el-tag class="ml-2" type="danger">第五步</el-tag>
    <el-tag class="ml-2" type="info">在"访达"中找到"应用程序"</el-tag>
    <br>
    <img :src="Image6"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第六步</el-tag>
    <el-tag class="ml-2" type="info">在"应用程序"中找到"实用工具"</el-tag>
    <br>
    <img :src="Image7"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第七步</el-tag>
    <el-tag class="ml-2" type="info">在"实用工具"中找到"终端",并且双击打开</el-tag>
    <br>
    <img :src="Image28"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第八步</el-tag>
    <el-tag class="ml-2" type="info">在终端中输入:</el-tag>
    <br>
    <el-tooltip
        class="box-item"
        effect="dark"
        content="点击复制"
        placement="bottom"
    >
      <el-tag class="ml-2" type="info"
              @click="copyText('sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ')"
              style="cursor: pointer">sudo security add-trusted-cert -d -r trustRoot -k
        /Library/Keychains/System.keychain
      </el-tag>
    </el-tooltip>

    <br>
    <el-tag class="ml-2" type="danger">请注意在 粘贴/输入 完成后在最后敲一个空格</el-tag>
    <br>
    <img :src="Image29"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第九步</el-tag>
    <el-tag class="ml-2" type="info">将刚才导出到桌面的证书拖入“终端”，然后回车执行</el-tag>
    <br>
    <img :src="Image30"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第十步</el-tag>
    <el-tag class="ml-2" type="info">回车执行之后,在终端中输入密码,按下回车</el-tag>
    <br>
    <img :src="Image31"/>
    <br> <br>
    <el-tag class="ml-2" type="danger">第十一步</el-tag>
    <el-tag class="ml-2" type="info">在弹出的窗口中再次输入密码,并且点击"更新设置"</el-tag>
    <br>
    <img :src="Image32"/>
  </div>
  <br> <br>
  <div v-if="radio!==''">
    <el-tag class="ml-2" type="danger">到此证书安装完成</el-tag>
  </div>
</template>

<script>

import {GetDocImage} from "../image.js";
import {ClipboardSetText} from "../../../wailsjs/runtime/runtime.js";
import {ElMessage} from "element-plus";

export default {
  components: {Image},
  computed: {
    Image1() {
      return GetDocImage("Image1")
    },
    Image2() {
      return GetDocImage("Image2")
    },
    Image3() {
      return GetDocImage("Image3")
    },
    Image4() {
      return GetDocImage("Image4")
    },
    Image5() {
      return GetDocImage("Image5")
    },
    Image6() {
      return GetDocImage("Image6")
    },
    Image7() {
      return GetDocImage("Image7")
    },
    Image8() {
      return GetDocImage("Image8")
    },
    Image9() {
      return GetDocImage("Image9")
    },
    Image10() {
      return GetDocImage("Image10")
    },
    Image11() {
      return GetDocImage("Image11")
    },
    Image12() {
      return GetDocImage("Image12")
    },
    Image13() {
      return GetDocImage("Image13")
    },
    Image14() {
      return GetDocImage("Image14")
    },
    Image28() {
      return GetDocImage("Image28")
    },
    Image29() {
      return GetDocImage("Image29")
    },
    Image30() {
      return GetDocImage("Image30")
    },
    Image31() {
      return GetDocImage("Image31")
    },
    Image32() {
      return GetDocImage("Image32")
    }
  },
  data() {
    return {
      radio: ""
    }
  }, methods: {
    copyText(a) {
      ClipboardSetText(a).then(res => {
        if (res) {
          ElMessage({
            message: "复制成功",
            type: 'success',
          })
        } else {
          ElMessage({
            message: "复制失败",
            type: 'error',
          })
        }
      })
    }
  }
}
</script>