<template>
  <div :style="WindowStyle">
    <VueJsonEditor
        ref="JsonEdit"
        v-model:json="jsonData"
        :json="jsonData"
        :readOnly="ReadOnly"
        :darkTheme="getTheme"
        :mainMenuBar="false"
        style="height: 100%; width: 100%;"
        :mode="Mode"
    ></VueJsonEditor>
  </div>
</template>
<script>
import VueJsonEditor from '../vue3-ts-jsoneditor';
import {PbJsonConvert, protobufToJson, UInt8ToStr} from "../CallbackEventsOn.js";


export default {
  props: ['readOnly', 'height', 'width'],
  watch: {
    readOnly(newValue) {
      this.ReadOnly = newValue
    },
    height(newValue) {
      this.Rect.height = newValue
      this.WindowStyle = "height: " + this.Rect.height + ";width: " + this.Rect.width + ";position: fixed;"
    },
    width(newValue) {
      this.Rect.width = newValue
      this.WindowStyle = "height: " + this.Rect.height + ";width: " + this.Rect.width + ";position: fixed;"
    },
  },
  components: {
    VueJsonEditor
  },
  computed: {
    getTheme() {
      return this.theme
    }
  },
  data() {
    return {
      Rect: {
        height: "0px",
        width: "0px",
      },
      mountedOK: false,
      WindowStyle: "height: 100%; width: 100%;",
      Raw: "",
      jsonData: null,
      ReadOnly: true,
      ReadOnly2: true,
      Mode: "tree",
      get theme() {
        return window.Theme.IsDark
      },
      set theme(newValue) {
        window.Theme.IsDark = newValue
      },
      AddElement: {
        Box: null,
        Input: null,
        Span1: null,
        ALL: null
      }
    }
  },
  methods: {
    SetCode(newValue) {
      this.jsonData = null
      this.Raw = newValue
      if (this.AddElement.Box) {
        if (this.AddElement.Box.checked) {
          this.pbConvert(this.Raw, this.AddElement.Input.value)
          return
        }
      }
      try {
        let _Body = UInt8ToStr(newValue, "utf-8")
        if (_Body.indexOf("�") !== -1) {
          _Body = UInt8ToStr(newValue, "gbk")
        }
        const c = _Body
        if (!c.includes('{') && !c.includes('[')) {
          return;
        }
        if (c.length > 10240) {
          this.jsonData = {"Data": "此JSON文件过大,请使用其他工具分析"}
          this.SetReadOnly(true)
          this.$nextTick(() => {
            this.eventClick()
          })
          return;
        }
        this.jsonData = JSON.parse(c)
        this.$nextTick(() => {
          this.eventClick()
        })
      } catch (e) {

      }
    },
    GetCode() {
      if (this.AddElement.Box.checked) {
        let _Body = UInt8ToStr(this.Raw, "utf-8")
        if (_Body.indexOf("�") !== -1) {
          _Body = UInt8ToStr(this.Raw, "gbk")
        }
        return _Body
      }
      return JSON.stringify(this.jsonData)
    },
    IsHasModify() {
      try {
        let _Body = UInt8ToStr(this.Raw, "utf-8")
        if (_Body.indexOf("�") !== -1) {
          _Body = UInt8ToStr(this.Raw, "gbk")
        }
        const obj1 = JSON.parse(_Body)
        const obj2 = JSON.parse(JSON.stringify(this.jsonData))
        return !this.deepEqual(obj1, obj2)
      } catch (e) {

      }
      return false
    },
    deepEqual(obj1, obj2) {
      // 检查两个对象的引用是否相等
      if (obj1 === obj2) {
        return true;
      }

      // 检查两个对象的类型是否相同
      if (typeof obj1 !== "object" || typeof obj2 !== "object" || obj1 === null || obj2 === null) {
        return false;
      }

      // 检查两个对象的属性数量是否相同
      const keys1 = Object.keys(obj1);
      const keys2 = Object.keys(obj2);

      if (keys1.length !== keys2.length) {
        return false;
      }

      // 递归比较对象的属性值
      for (let key of keys1) {
        if (!keys2.includes(key) || !this.deepEqual(obj1[key], obj2[key])) {
          return false;
        }
      }
      return true;
    },
    SetReadOnly(MODE) {
      this.ReadOnly = MODE
      this.ReadOnly2 = MODE
      this.AddElement.Box.checked = false
    },
    pbConvert(pb, skip) {
      this.ReadOnly = true
      protobufToJson(pb, skip).then(response => {
        try {
          if (response.length > 10240) {
            this.jsonData = {"Data": "此ProtoBuf文件过大,请使用其他工具分析"}
            this.$nextTick(() => {
              this.eventClick()
            })
            return
          }
          this.jsonData = PbJsonConvert(JSON.parse(response))
          this.$nextTick(() => {
            this.eventClick()
          })
        } catch (e) {
        }
      })
    },
    eventClick() {
      this.$nextTick(() => {
        this.$refs.JsonEdit.$expandAll();
        this.$nextTick(() => {
          const header = document.getElementsByClassName("jse-selected-value")[0];
          if (header) {
            header.className = "jse-json-node jse-expanded svelte-1wlxti1 jse-root jse-readonly"
          }
          const header1 = document.getElementsByClassName("jse-selected")[0];
          if (header1) {
            header1.className = "jse-json-node jse-expanded svelte-1wlxti1 jse-root jse-readonly"
          }
        })
      })
    },
    HandleCheckboxChange(event) {
      this.ReadOnly = this.ReadOnly2
      this.jsonData = null
      if (event.target.checked) {
        this.pbConvert(this.Raw, this.AddElement.Input.value)
        this.$nextTick(() => {
          const mm = this.$refs.JsonEdit.$el.getElementsByClassName("jse-contents")
          for (let i = 0; i < mm.length; i++) {
            if (mm[i].innerText.indexOf("空的 JSON 对象") !== -1) {
              mm[i].innerText = "空的 ProtoBuf 对象"
            }
          }
        })
        return
      }

      this.$nextTick(() => {
        const mm = this.$refs.JsonEdit.$el.getElementsByClassName("jse-contents")
        for (let i = 0; i < mm.length; i++) {
          if (mm[i].innerText.indexOf("ProtoBuf") !== -1) {
            mm[i].innerText = "空的 JSON 对象"
          }
        }
      })

      try {
        let _Body = UInt8ToStr(this.Raw, "utf-8")
        if (_Body.indexOf("�") !== -1) {
          _Body = UInt8ToStr(this.Raw, "gbk")
        }
        const c = _Body
        if (!c.includes('{') && !c.includes('[')) {
          return;
        }
        if (c.length > 10240) {
          this.jsonData = {"Data": "此JSON文件过大,请使用其他工具分析"}
          this.SetReadOnly(true)
          this.$nextTick(() => {
            this.eventClick()
          })
          return;
        }
        this.jsonData = JSON.parse(c)
        this.$nextTick(() => {
          this.eventClick()
        })
      } catch (e) {
      }

    },
    HandleSpanClick(event) {
      if (this.AddElement.Box) {
        this.AddElement.Box.click();
      }
    },
    HandleInput(event) {
      if (this.AddElement.Box.checked) {
        this.pbConvert(this.Raw, event.target.value)
      }
    },
    RegComponent() {
      setTimeout(() => {
        const Edit = this.$refs.JsonEdit.$el
        const el = Edit.getElementsByClassName("jse-navigation-bar")[0]
        if (el) {
          {
            if (this.AddElement.ALL != null) {
              el.appendChild(this.AddElement.ALL);
              this.mountedOK = true
              return
            }
            const Element = document.createElement('div');
            Element.style.position = "absolute"
            Element.style.left = "calc(100% - 300px)"
            const checkbox = document.createElement('input');
            checkbox.type = 'checkbox';
            checkbox.style.position = "absolute"
            checkbox.style.left = "calc(50% + 3px)"
            checkbox.style.top = "5px"
            checkbox.addEventListener('change', this.HandleCheckboxChange);
            this.AddElement.Box = checkbox
            Element.appendChild(checkbox)

            const span1 = document.createElement('span');
            span1.innerText = "这可能是ProtoBuf数据,忽略前"
            span1.style.position = "absolute"
            span1.style.left = "calc(50% + 23px)"
            span1.style.width = "200px"
            span1.style.top = "7px"
            span1.style.textAlign = "left"
            span1.addEventListener('click', this.HandleSpanClick);
            this.AddElement.Span1 = span1
            Element.appendChild(span1)
            const input = document.createElement('input');
            input.type = 'number';
            input.value = '0'; // 设置默认值为 123
            input.style.width = "35px"
            input.style.position = "absolute"
            input.style.left = "calc(50% + 221px)"
            input.style.top = "5px"
            input.style.height = "13px"
            input.readOnly = false; // 或者省略该行
            input.addEventListener('input', this.HandleInput);
            this.AddElement.Input = input
            Element.appendChild(input)
            const span2 = document.createElement('span');
            span2.innerText = " 字节"
            span2.style.position = "absolute"
            span2.style.left = "calc(50% + 268px)"
            span2.style.width = "200px"
            span2.style.top = "7px"
            span2.style.textAlign = "left"
            Element.appendChild(span2)
            this.AddElement.ALL = Element
            el.appendChild(Element);
            this.mountedOK = true
          }
        }
        if (!this.mountedOK) {
          this.RegComponent()
        }
      }, 100)
    },
    ReleaseComponent() {
      this.AddElement.Box.removeEventListener('change', this.HandleCheckboxChange);
      this.AddElement.Span1.addEventListener('click', this.HandleSpanClick);
      this.AddElement.Input.addEventListener('input', this.HandleInput);
    }
  },
  mounted() {
    this.RegComponent()

  }
}
</script>
