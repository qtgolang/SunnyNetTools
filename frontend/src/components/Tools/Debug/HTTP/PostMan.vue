<script xmlns="">
import VueText from "../../Theme/Text.vue";
import VueText2 from "./TextR.vue";
import VueText3 from "./RawHeader.vue";
import VueText4 from "../../../Home/Request/tool/Raw.vue";
import {ElLoading, ElNotification} from "element-plus";
import {Pointer, Promotion} from "@element-plus/icons-vue";
import {
  AppStartSelectRequest,
  ClipboardReadAll,
  ClipboardWriteAll,
  GetLocalServerPATH
} from "../../../../../bindings/changeme/Service/appmain";
import {bytesToRequest, parseHTTPHeaders} from "../../../config/SunnyNetInfoApi";
import {base64ToBytes, bytesToBase64, bytesToString, StringToBytes} from "../../../config/encoding";
import {NewBytesIO} from "../../../config/BytesIO";
import {getImageType} from "../../../config/CheckImageType";
import {Events} from "@wailsio/runtime";
import {GetGenerateCodeListDebugToolsMenu} from "../../../config/GenerateCode";

export default {
  components: {VueText, VueText2, VueText3, VueText4},
  data() {
    return {
      form: {
        url: '',
        mode: 'GET',
        dataType: 'text',
        proxyType: "HTTP",
        proxy: "",
        outTime: "10",
        fakeIP: false,
        redirect: false,
      },
      xzReq: false,
      isImg: false,
      ImgBody: "",
      options: [],
      optionsValue: "",
    }
  },
  computed: {
    Promotion() {
      return Promotion
    },
    Pointer() {
      return Pointer
    },
  },
  methods: {
    paste() {
      ClipboardReadAll().then(res => {
        const obj = bytesToRequest(StringToBytes(res))
        if (obj != null) {
          let header = '';
          let cookie = '';
          let host = '';
          const keys = Object.keys(obj.Header);
          for (let i = 0; i < keys.length; i++) {
            if (keys[i].toLowerCase() === "host") {
              host = obj.Header[keys[i]][0]
            }
            if (keys[i].toLowerCase() === "cookie") {
              cookie += obj.Header[keys[i]].join(";") + ";"
              continue
            }
            for (let j = 0; j < obj.Header[keys[i]].length; j++) {
              header += keys[i] + ": " + obj.Header[keys[i]][j] + "\n"
            }
          }

          if (!obj.URL.toLowerCase().startsWith("http")) {
            this.form.url = "https://" + host + obj.URL
          } else {
            this.form.url = obj.URL
          }
          this.form.mode = obj.Method
          this.form.dataType = 'text'
          const body = bytesToString(Base64ToBytes(obj.data))
          try {
            JSON.parse(body)
            this.$refs.RequestText.SetLanguage("json")
          } catch (e) {
            this.$refs.RequestText.SetLanguage("plaintext")
          }
          this.$refs.RequestText.SetCode(body)
          this.$refs.RequestHeader.SetCode(header)
          this.$refs.RequestCookies.SetCode(cookie)
        }
      })
    },
    async Send() {
      const loading = ElLoading.service({
        lock: true,
        text: '请求正在发送中....',
        target: document.getElementById('postMan'),
      });

      try {
        const proxyAddress = `${this.form.proxyType}://${this.form.proxy}`;

        const headerBytes = StringToBytes(this.$refs.RequestHeader.GetCode());
        const headerStream = NewBytesIO(headerBytes);
        const headers = parseHTTPHeaders(headerStream);

        const requestBodyBase64 = bytesToBase64(StringToBytes(this.$refs.RequestText.GetCode()));

        const requestPayload = {
          url: this.form.url,
          method: this.form.mode,
          header: headers,
          outTime: parseInt(this.form.outTime),
          redirect: this.form.redirect,
          disguise: this.form.fakeIP,
          proxyIP: proxyAddress,
          body: requestBodyBase64,
          bodyType: this.form.dataType,
        };

        const serverPath = await GetLocalServerPATH();
        const response = await fetch(`${serverPath}/DoHTTPRequest`, {
          method: 'POST',
          body: JSON.stringify(requestPayload),
        });

        const responseText = await response.text();

        try {
          const responseJson = JSON.parse(responseText);

          if (responseJson.errorLevel !== 0) {
            this.handleErrorResponse(responseJson);
          } else {
            this.handleSuccessResponse(responseJson);
          }
        } catch (parseError) {
          this.displayErrorResult(responseText);
        }
      } catch (err) {
        this.displayErrorResult(err.message || '未知错误');
      } finally {
        loading.close();
      }
    },

    handleErrorResponse(responseJson) {
      let errorMessage = "请求失败:\n\n" + (responseJson.Error || '未知错误');

      if (responseJson.errorLevel === 1) {
        errorMessage = "请求失败:\n\n您填写的请求 HEX 数据不正确,请检查\n\n" + responseJson.Error;
      } else if (responseJson.errorLevel === 2) {
        errorMessage = "请求失败:\n\n您填写的请求 Base64 数据不正确,请检查\n\n" + responseJson.Error;
      }

      this.$refs.ResponseText.SetCode(errorMessage);
      this.$refs.ResponseRaw.SetCode(errorMessage);
      this.$refs.ResponseCookies.SetCode(errorMessage);
    },

    handleSuccessResponse(responseJson) {
      const responseBody = bytesToString(Base64ToBytes(responseJson.body));
      const headers = responseJson.header;
      const headerKeys = Object.keys(headers);

      let headerText = '';
      let cookieText = '';
      let language = 'plaintext';

      for (const key of headerKeys) {
        if (key.toLowerCase() === 'set-cookie') {
          cookieText += headers[key].join("\n\n") + "\n\n";
          continue;
        }

        if (key.toLowerCase() === 'content-type') {
          const contentTypeParts = headers[key].join(";").replaceAll(";", "/").split("/");
          if (contentTypeParts.length >= 2) {
            language = contentTypeParts[0].toLowerCase() === 'image'
                ? (contentTypeParts[0] + '/' + contentTypeParts[1]).toLowerCase()
                : contentTypeParts[1].toLowerCase();
          }
        }

        for (const value of headers[key]) {
          headerText += `${key}: ${value}\n`;
        }
      }
      this.isImg = language.includes("image");
      if (this.isImg) {
        const ImageType = getImageType(base64ToBytes(responseJson.body));
        this.ImgBody = "data:image/" + ImageType + ";base64," + responseJson.body
      }
      const rawResponse = `${responseJson.proto} ${responseJson.status}\n${headerText}\n\n${responseBody}`;

      this.$refs.ResponseText.SetCode(responseBody);
      this.$refs.ResponseRaw.SetCode(rawResponse);
      this.$refs.ResponseCookies.SetCode(cookieText);

      this.$refs.ResponseRaw.SetLanguage(language);
      this.$refs.ResponseText.SetLanguage(language);
    },

    displayErrorResult(message) {
      const errorText = "请求失败:\n\n" + message;
      this.$refs.ResponseText.SetCode(errorText);
      this.$refs.ResponseRaw.SetCode(errorText);
      this.$refs.ResponseCookies.SetCode(errorText);
    },

    async GenerateCode(Language, Type) {
      const headerBytes = StringToBytes(this.$refs.RequestHeader.GetCode());
      const headerStream = NewBytesIO(headerBytes);
      const headers = parseHTTPHeaders(headerStream);
      const requestBodyBase64 = bytesToBase64(StringToBytes(this.$refs.RequestText.GetCode()));
      const requestPayload = {
        url: this.form.url,
        method: this.form.mode,
        header: headers,
        outTime: parseInt(this.form.outTime),
        redirect: this.form.redirect,
        disguise: this.form.fakeIP,
        proxyIP: "",
        body: requestBodyBase64,
        bodyType: this.form.dataType,
        Language: Language,
        Type: Type,
      };

      const serverPath = await GetLocalServerPATH();
      const response = await fetch(`${serverPath}/GenerateCode`, {
        method: 'POST',
        body: JSON.stringify(requestPayload),
      });
      if (response === null || response === undefined) {
        ElNotification({
          position: 'bottom-right',
          message: '生成代码失败\n请检查接口函数是否编写正确？',
          type: 'warning',
          customClass: 'multiline-message'
        })
        return
      }
      const responseText = await response.text();
      if (responseText === null || responseText === undefined || responseText === "") {
        ElNotification({
          position: 'bottom-right',
          message: '生成代码失败\n请检查接口函数是否编写正确？',
          type: 'warning',
          customClass: 'multiline-message'
        })
        return
      }
      ElNotification({
        position: 'bottom-right',
        showClose: true,
        message: '代码生成成功\n\n已复制到剪辑版',
        type: 'success',
        customClass: 'multiline-message'
      })
    },
    async codeDebugToolsFunc(func) {
      const headerBytes = StringToBytes(this.$refs.RequestHeader.GetCode());
      const headerStream = NewBytesIO(headerBytes);
      const headers = parseHTTPHeaders(headerStream);
      let res = "";
      {
        res = func({
          isTCP: false,
          isWebsocket: false,
          isHTTP: true,
          Request: {
            Method: this.form.mode,
            URL: this.form.url,
            Header: headers,
            Body: StringToBytes(this.$refs.RequestText.GetCode()),
            ServerIP: "",
            TmpFile: ""
          }
        })
      }
      if (res === null || res === undefined || res === "") {
        ElNotification({
          position: 'bottom-right',
          message: '生成代码失败\n请检查接口函数是否编写正确？',
          type: 'warning',
          customClass: 'multiline-message'
        })
        return
      }
      ClipboardWriteAll(res)
      ElNotification({
        position: 'bottom-right',
        showClose: true,
        message: '代码生成成功\n\n已复制到剪辑版',
        type: 'success',
        customClass: 'multiline-message'
      })
    },
    GenerateCodeList() {
      this.options = GetGenerateCodeListDebugToolsMenu(this.GenerateCode, this.codeDebugToolsFunc)
    },
    findByPath(data, path) {
      let current = data;
      let result = null;
      for (let value of path) {
        result = current.find(item => item.value === value);
        if (!result) return null;
        current = result.children || [];
      }
      return result;
    }
  },
  watch: {
    xzReq(n) {
      AppStartSelectRequest(n)
    },
    optionsValue(n) {
      if (n === "") {
        return;
      }
      if (this.form.url === "") {
        ElNotification({
          position: 'bottom-right',
          message: '生成代码失败\n请先填写URL',
          type: 'warning',
          customClass: 'multiline-message'
        })
      } else {
        const obj = this.findByPath(this.options, n)
        if (obj) {
          if (obj.action) {
            obj.action()
          }
        }
      }
      requestAnimationFrame(() => {
        this.optionsValue = ""
      })
    }
  },
  mounted() {
    this.$refs.RequestText.SetReadOnly(false)
    this.$refs.RequestCookies.SetReadOnly(false)
    this.$refs.RequestHeader.SetReadOnly(false)
    this.$refs.ResponseCookies.SetReadOnly(false)
    this.$refs.ResponseRaw.SetReadOnly(false)
    this.$refs.ResponseText.SetReadOnly(false)
    Events.On("SetSelectRequest", (obj) => {
      const url = obj.data[0];
      const Method = obj.data[1];
      const Header = obj.data[2];
      const Body = obj.data[3];
      const BodyType = obj.data[4];


      let header = '';
      let cookie = '';
      const keys = Object.keys(Header);
      for (let i = 0; i < keys.length; i++) {
        if (keys[i].toLowerCase() === "cookie") {
          cookie += Header[keys[i]].join(";") + ";"
          continue
        }
        for (let j = 0; j < Header[keys[i]].length; j++) {
          header += keys[i] + ": " + Header[keys[i]][j] + "\n"
        }
      }

      this.form.url = url
      this.form.mode = Method
      this.form.dataType = BodyType
      try {
        JSON.parse(Body)
        this.$refs.RequestText.SetLanguage("json")
      } catch (e) {
        this.$refs.RequestText.SetLanguage("plaintext")
      }
      if (BodyType === "hex") {
        this.$refs.RequestText.SetCode(Body)
      } else {
        this.$refs.RequestText.SetCode(bytesToString(base64ToBytes(Body)))
      }
      this.$refs.RequestHeader.SetCode(header)
      this.$refs.RequestCookies.SetCode(cookie)
    })
    window.GenerateCodeList = this.GenerateCodeList
  },
};
</script>

<template>
  <div style="width: 100%; height: 100%" id="postMan">
    <div style="width: 100%; ">
      <el-input
          v-model="form.url"
          style="max-width: calc(100% - 10px);"
          placeholder="请输入请求地址"
          class="input-with-select"
      >
        <template #prepend>
          <el-select v-model="form.mode" placeholder="Select" style="width: 100px;">
            <el-option label="GET" value="GET"/>
            <el-option label="POST" value="POST"/>
            <el-option label="PUT" value="PUT"/>
          </el-select>
        </template>
        <template #append>
          <el-tooltip class="item" content="从剪辑版粘贴请求" placement="top">
            <el-button :icon="Pointer" @click="paste"/>
          </el-tooltip>
          <div style="margin-left: 20px;margin-right: 20px">|</div>
          <el-tooltip class="item" content="立即发送请求" placement="top">
            <el-button :icon="Promotion" @click="Send"/>
          </el-tooltip>
        </template>
      </el-input>
      <div style="width: 100%;display: flex;margin-top: 5px;margin-bottom: 5px;margin-left: 5px">
        <el-input
            v-model="form.proxy"
            style="max-width: 100%;margin-right: 5px"
            class="input-with-select"
            placeholder="(账号:密码@IP:端口) 例如(admin:pass@127.0.0.1:8888)或(IP:端口) 例如(127.0.0.1:8888)"
        >
          <template #prepend>
            <el-select v-model="form.proxyType" placeholder="Select" style="width: 100px">
              <el-option label="HTTP" value="HTTP"/>
              <el-option label="Socks" value="Socks"/>
            </el-select>
          </template>
        </el-input>

        <div style="display: flex;margin-right: 5px">
          <el-cascader
              placeholder="生成代码"
              v-model="optionsValue"
              :options="options"
          />
        </div>

        <el-input
            v-model="form.outTime"
            style="max-width: 200px;justify-content: center;"
            placeholder="10"
        >
          <template #prepend>超时时间</template>
          <template #append>秒</template>
        </el-input>
        <el-tooltip class="item" content="使用伪装IP" placement="bottom">
          <el-checkbox v-model="form.fakeIP" style="margin-left: 10px;width: 50px">伪装IP</el-checkbox>
        </el-tooltip>
        <el-tooltip class="item" content="是否禁止重定向" placement="bottom">
          <el-checkbox v-model="form.redirect" style="left: 0;width: 75px">禁止重定向</el-checkbox>
        </el-tooltip>
        <el-tooltip class="item" content="开启后,抓包工具选择的请求将同步到这里" placement="bottom">
          <el-checkbox v-model="xzReq" style="margin-left: 0;margin-right: 15px">选择请求</el-checkbox>
        </el-tooltip>

      </div>
    </div>
    <div style="width: 100%; height: 100%;display: flex; margin-top: 0;">
      <div ref="tabs" style="width: 50%; height: calc(100% - 78px);">
        <el-tabs type="border-card" style="height: 100%">
          <el-tab-pane label="提交数据" style="height: 100%">
            <div style="width: 100%; height: 100%; display: flex; flex-direction: column;">
              <el-select v-model="form.dataType" style="width: 100%; margin-bottom: 10px;" size="small">
                <el-option label="文本发送" value="text"/>
                <el-option label="十六进制解析后发送" value="hex"/>
                <el-option label="Base64解析后发送" value="base64"/>
              </el-select>
              <div style="flex: 1;">
                <VueText ref="RequestText" Language="'text'" Name="data"
                         :Overflow="true"
                         style="height: 100%;"/>
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane label="提交Cookies" style="height: 100%;">
            <div style="width: 100%;height: 100%;display: flex; flex-direction: column;">
              <VueText ref="RequestCookies" Language="'text'" :zIndex="'1000'"
                       :Overflow="true"
                       Name="data" style="margin-top: -1px;margin-left: -1px;height: 100%;"/>
            </div>
          </el-tab-pane>
          <el-tab-pane label="提交协议头" style="height: 100%;">
            <div style="width: 100%;height:100%;display: flex; flex-direction: column;">
              <VueText3 ref="RequestHeader" Language="'text'"
                        :Overflow="true"
                        :NoShowEncoding="true"
                        Name="data" style="margin-top: -1px;margin-left: -1px;height: 100%"/>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
      <div ref="tabs2" style="width: 50%; height: calc(100% - 78px);left: 50%;">
        <el-tabs type="border-card" style="height: 100%;">
          <el-tab-pane label="原始响应" style="height: 100%">
            <div style="width: 100%;height:100%">
              <VueText4 ref="ResponseRaw" :glyphMargin="false" :readOnly="false" Language="'text'"
                        :Overflow="true"
                        Name="data" style="margin-top: -1px;margin-left: -1px;height: 100%"/>
            </div>
          </el-tab-pane>

          <el-tab-pane label="响应Cookies" style="height: 100%;">
            <div style="width: 100%;height:100%">
              <VueText2 ref="ResponseCookies" :glyphMargin="false" :readOnly="false" Language="'text'"
                        :Overflow="true"
                        Name="data" style="margin-top: -1px;margin-left: -1px;"/>
            </div>
          </el-tab-pane>
          <el-tab-pane label="响应文本" style="height: 100%;">
            <div style="width: 100%;height:100%">
              <VueText3 ref="ResponseText" :glyphMargin="false" :readOnly="false" Language="'text'"
                        :Overflow="true"
                        Name="data" style="margin-top: -1px;margin-left: -1px;"/>
            </div>
          </el-tab-pane>
          <el-tab-pane label="响应图片" v-if="isImg" style="height: 100%;">
            <div style="width: 100%;height:100%">
              <img
                  ref="imgObj"
                  :src="ImgBody"
                  style="height: 100%; width: 100%; object-fit: contain; object-position: center; display: block;"
                  alt=""/>
            </div>
          </el-tab-pane>


        </el-tabs>
      </div>
    </div>
  </div>

</template>