<template>

  <div class="demo-collapse" style="">
    <el-collapse accordion>
      <el-collapse-item title="安卓7.0以下" name="1">
        <div>
          打开自带浏览器 输入网址 http://电脑IP:{{ this.Port }}
        </div>
        <div>
          下载后安装即可
        </div>
      </el-collapse-item>
      <el-collapse-item title="安卓7.0及以上" name="2">
        <div
            style="position: relative; left: 10px;display: inline-grid; justify-content: center; align-items: center;">
          <div style="text-align: left">
            1.在本软件设置中找到SSL证书,导出默认证书
          </div>
          <div style="text-align: left">
            2.并且改名为
            <el-tag class="ml-2" type="danger">298807fb.0</el-tag>
          </div>
          <div style="text-align: left">
            3.将 298807fb.0 证书 push到手机： /sdcard/298807fb.0
            <div style="position: relative; left: 20px;">
              <div style="text-align: left">
                3.1.
                <el-tag class="ml-2" type="info">adb push 298807fb.0 /sdcard/298807fb.0</el-tag>
              </div>
              <div style="text-align: left">
                3.2.然后依次执行以下命令
              </div>
              <div style="text-align: left">
                3.3.
                <el-tag class="ml-2" type="info">adb shell</el-tag>
              </div>
              <div style="text-align: left">
                3.4.
                <el-tag class="ml-2" type="info">su</el-tag>
              </div>
              <div style="text-align: left">
                3.5.
                <el-tag class="ml-2" type="info">mount -o rw,remount /system</el-tag>
              </div>
              <div style="text-align: left;width: calc(100% - 30px)">
                3.6.
                <el-text class="mx-1" type="danger">如果上面这个命令出错(以下步骤不用看了) 表示你手机系统无法挂载system分区为可读写
                  ，将不能安装证书,[请改用面具 安装模块Move Certificates
                  <el-tooltip
                      class="box-item"
                      effect="dark"
                      content="点击去蓝奏下载"
                      placement="top"
                  >
                    <el-text class="mx-1" type="success" @click="callUrl('https://wwc.lanzouy.com/iynzR0buml1g')"
                             style="cursor: pointer; text-decoration: underline;">蓝奏下载
                    </el-text>
                  </el-tooltip>
                  ,
                  <el-tooltip
                      class="box-item"
                      effect="dark"
                      content="点击去GitHub下载"
                      placement="top"
                  >
                    <el-text class="mx-1" type="success"
                             @click="callUrl('https://github.com/Magisk-Modules-Repo/movecert')"
                             style="cursor: pointer; text-decoration: underline;">GitHub下载
                    </el-text>
                  </el-tooltip>

                  ]
                  安装后使用7.0以下的安装方式即可
                </el-text>
              </div>
              <div style="text-align: left">
                3.7.
                <el-tag class="ml-2" type="info">mv /sdcard/298807fb.0 /system/etc/security/cacerts/298807fb.0</el-tag>
              </div>
              <div style="text-align: left">
                3.8.
                <el-tag class="ml-2" type="info">cd /system/etc/security/cacerts</el-tag>
              </div>
              <div style="text-align: left">
                3.9.
                <el-tag class="ml-2" type="info">chmod 644 298807fb.0</el-tag>
              </div>
              <div style="text-align: left">
                3.10.
                <el-tag class="ml-2" type="info">reboot</el-tag>
              </div>
            </div>
          </div>
        </div>
      </el-collapse-item>
      <el-collapse-item title="雷电9模拟器" name="3">
        <div
            style="position: relative; left: 10px;display: inline-grid; justify-content: center; align-items: center;">
          <div style="text-align: left">
            1.下载所需文件
            <el-tooltip
                class="box-item"
                effect="dark"
                content="点击去下载"
                placement="top"
            >
              <el-text class="mx-1" type="success" @click="callUrl('https://wwf.lanzouy.com/ivcGA0g1wxle')"
                       style="cursor: pointer; text-decoration: underline;">[Magisk_v25.2.apk] [Magisk
                Terminal Emulator_17.7.apk] [MT2.12.0.apk] [app-debug.apk] [Move_Certificates-v1.9(10).zip]
              </el-text>
            </el-tooltip>

          </div>
          <div style="text-align: left">
            2.新建模拟器（先不要运行）选择设置
          </div>
          <div style="text-align: left">
            3.性能设置 -> 磁盘共享 选择
            <el-tag class="ml-2" type="danger">( System.vmdk 可写入 )</el-tag>
          </div>
          <div style="text-align: left">
            4.性能设置 -> 分辨率 -> 选择 手机版 540x960
            <el-tag class="ml-2" type="danger">( 非必须 )</el-tag>
          </div>
          <div style="text-align: left">
            5.其他设置 -> ROOT权限 ->
            <el-tag class="ml-2" type="danger">开启</el-tag>
          </div>
          <div style="text-align: left">
            6.保存设置 -> 启动模拟器
          </div>
          <div style="text-align: left">
            7.安装刚刚下载的所需文件中的
            <el-tag class="ml-2" type="danger">[ Magisk_v25.2.apk ] [ Magisk Terminal Emulator_17.7.apk ] [ MT2.12.0.apk
              ]
            </el-tag>
          </div>
          <div style="text-align: left">
            8.将 所需文件中的
            <el-tag class="ml-2" type="danger"> [app-debug.apk] [Move_Certificates-v1.9(10).zip]</el-tag>
            放到共享目录
          </div>
          <div style="text-align: left">
            9.打开
            <el-tag class="ml-2" type="danger">Magisk Terminal Emulator</el-tag>
            ->若提示需要什么权限或ROOT权限 全部同意
          </div>
          <div style="text-align: left">
            10.打开后 输入 "
            <el-tag class="ml-2" type="danger">m</el-tag>
            ->然后按回车键
          </div>
          <div style="text-align: left">
            11.然后 输入 "
            <el-tag class="ml-2" type="danger">y</el-tag>
            "然后按回车键 ->可能申请ROOT权限 选择 永久记住选择 然后 允许
          </div>
          <div style="text-align: left">
            12.然后 输入 "
            <el-tag class="ml-2" type="danger">1</el-tag>
            "然后按回车键
          </div>
          <div style="text-align: left">
            13.然后 输入 "
            <el-tag class="ml-2" type="danger">x</el-tag>
            "然后按回车键
          </div>
          <div style="text-align: left">
            14.然后 输入 "
            <el-tag class="ml-2" type="danger">/sdcard/Pictures/app-debug.apk</el-tag>
            "然后按回车键 (不能粘贴，手动输入一下)
          </div>
          <div style="text-align: left">
            15.然后 输入 "
            <el-tag class="ml-2" type="danger">1</el-tag>
            "然后按回车键 （没有出现红色的相关文字就是成功了）
          </div>
          <div style="text-align: left">
            16.然后 回到桌面 打开
            <el-tag class="ml-2" type="danger">MT管理器</el-tag>
          </div>
          <div style="text-align: left">
            17.找到这个文件并且 删除
            <el-tag class="ml-2" type="danger">/system/xbin/su</el-tag>
          </div>
          <div style="text-align: left">
            18.
            <el-tag class="ml-2" type="danger">重启模拟器</el-tag>
          </div>
          <div style="text-align: left">
            19.打开
            <el-tag class="ml-2" type="danger">Magisk</el-tag>
            APP ->
            <el-tag class="ml-2" type="danger"> 右下角 模块</el-tag>
            ->
            <el-tag class="ml-2" type="danger">从本地安装</el-tag>
            ->
            <el-tag class="ml-2" type="danger">选 择左上角三个横杠</el-tag>
            ->
            <el-tag class="ml-2" type="danger">文件管理器</el-tag>
            ->
            <el-tag class="ml-2" type="danger">Pictures</el-tag>
            ->
            <el-tag class="ml-2" type="danger">Move_Certificates-v1.9(10).zip</el-tag>
          </div>
          <div style="text-align: left">
            20.
            <el-tag class="ml-2" type="danger">重启模拟器</el-tag>
          </div>
          <div style="text-align: left">
            21.模拟器中 打开浏览器 输入http://IP:端口(你电脑的内网IP，和 软件的端口号)例如 http://192.168.31.111:{{
              Port
            }}
          </div>
          <div style="text-align: left">
            22.下载证书 安装
          </div>
          <div style="text-align: left">
            22.安装成功后重启模拟器
            <el-tag class="ml-2" type="danger">( 然后如果有需要的话，可以删除面具 和在模拟器关闭ROOT 均不影响抓包 )
            </el-tag>
          </div>
        </div>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script>
import {CallGoDo} from "../CallbackEventsOn.js";
import {BrowserOpenURL} from "../../../wailsjs/runtime/runtime.js";

export default {
  data() {
    return {
      Port: ""
    }
  }, mounted() {
    CallGoDo("获取运行端口", null).then(res => {
      this.Port = res
    })
  }, methods: {
    callUrl(a) {
      BrowserOpenURL(a)
    }
  }
}
</script>