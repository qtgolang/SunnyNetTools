<script>
import {ElMessage} from 'element-plus'
import {Base64BytesEncode, bytesToString, StringToBytes} from "../CallbackEventsOn.js";
import {ClipboardSetText} from "../../../wailsjs/runtime/runtime.js";

const insertMode = -2;

class MyApply {
  constructor() {
    this.Array = [];
  }

  SetAnsi(bytes, offset, color) {
    let _ansi = ToAnsi(bytes)
    let Length = 0;
    for (const item of this.Array) {
      Length += item.text.length;
    }
    if (Length !== 0) {
      let yu = offset - Length % offset
      if (yu !== 0) {
        let _ansi2 = _ansi.slice(0, yu)
        _ansi = _ansi.slice(yu, _ansi.length)
        this.Array.push({text: _ansi2, color: color})
      }
    }
    let arr = splitBytes(_ansi, offset)
    for (let i = 0; i < arr.length; i++) {
      this.Array.push({text: arr[i], color: color})
    }
  }

  SetHex(bytes, offset, color) {
    let _ansi = bytes
    let Length = 0;
    for (const item of this.Array) {
      Length += item.text.length;
    }
    if (Length !== 0) {
      let yu = offset - Length % offset
      if (yu !== 0) {
        let _ansi2 = _ansi.slice(0, yu)
        _ansi = _ansi.slice(yu, _ansi.length)
        this.Array.push({text: _ansi2, color: color})
      }
    }
    let arr = splitBytes(_ansi, offset)
    for (let i = 0; i < arr.length; i++) {
      this.Array.push({text: arr[i], color: color})
    }
  }

  AnsiApply(offers) {
    let str = ""
    let ix = 0
    for (let i = 0; i < this.Array.length; i++) {
      if (this.Array[i].color.back === "") {
        str += `<span tid="ansi" STYLE="color: ` + this.Array[i].color.text + `; display: inline-block;">` + Escape(this.Array[i].text) + "</span>"
      } else {

        str += `<span tid="ansi" STYLE="color: ` + this.Array[i].color.text + `; background-color: ` + this.Array[i].color.back + `;display: inline-block;">` + Escape(this.Array[i].text) + "</span>"
      }
      ix += this.Array[i].text.length
      if (ix % offers === 0) {
        str += "<br>"
      }
    }
    return str
  }

  HexApply() {
    let str = ""
    let sBack = "";
    let NBack = "";
    for (let i = 0; i < this.Array.length; i++) {
      let hex = Array.from(this.Array[i].text, byte => ('0' + (byte & 0xFF).toString(16)).slice(-2).toUpperCase()).join(' ');
      if (this.Array[i].color.back === "") {
        if (sBack !== "") {
          str += " "
        }
        if (NBack !== this.Array[i].color.text) {
          str += " "
          NBack = this.Array[i].color.text
        }
        NBack = ""
        sBack = ""
        str += `<span tid="hex" STYLE="color: ` + this.Array[i].color.text + `; display: inline-block;">` + hex + "</span>"
      } else {
        if (sBack === "") {
          str += " "
        }
        if (NBack !== this.Array[i].color.text) {
          str += " "
          NBack = this.Array[i].color.text
        }
        sBack = "1"
        str += `<span  tid="hex" STYLE="color: ` + this.Array[i].color.text + `; background-color: ` + this.Array[i].color.back + `;display: inline-block;">` + hex + "</span>"
      }
    }
    return str
  }
}

function Escape(t) {
  let s = t
  s = s.replaceAll("&", "&amp;")
  s = s.replaceAll(" ", "&ensp;")
  s = s.replaceAll("<", "&lt;")
  s = s.replaceAll(">", "&gt;")
  s = s.replaceAll("\"", "&quot;")
  s = s.replaceAll("\t", "&nbsp;")
  return s
}

function hexToBytes(hexString) {
  const length = hexString.length / 2;
  const bytes = new Uint8Array(length);

  for (let i = 0; i < length; i++) {
    const hex = hexString.substr(i * 2, 2);
    bytes[i] = parseInt(hex, 16);
  }

  return bytes;
}

function isHexadecimalString(str) {
  const hexRegex = /^[0-9a-fA-F]+$/;
  return hexRegex.test(RemoveSpaces(str));
}

function RemoveSpaces(str) {
  return str.replace(/[\s\u200B\u00A0]/g, '').replaceAll(" ", "").replaceAll("\t", "").replaceAll("\r", "").replaceAll("\n", "");
}

function deleteBytes(bytes, start, numToDelete) {
  bytes = Array.from(bytes);
  let newBytes = bytes.slice(0, start).concat(bytes.slice(start + numToDelete));
  return new Uint8Array(newBytes);
}

function hexToByte(hex) {
  return parseInt(hex, 16);
}

function byteToHex(byte) {
  return byte.toString(16).padStart(2, '0');
}


function insertBytes(originalBytes, index, insertBytes) {
  const start = originalBytes.subarray(0, index);
  const end = originalBytes.subarray(index);

  const result = new Uint8Array(start.length + insertBytes.length + end.length);
  result.set(start);
  result.set(insertBytes, start.length);
  result.set(end, start.length + insertBytes.length);

  return result;
}


function splitBytes(str, chunkSize) {
  if (chunkSize === 0) {
    return
  }
  const chunks = [];
  let n = 0;
  for (let i = 0; i < str.length; i += chunkSize) {
    n++
    chunks.push(str.slice(i, i + chunkSize));
  }
  return chunks;
}

function ToAnsi(bytes) {
  let str = '';
  for (let i = 0; i < bytes.length; i++) {
    if (bytes[i] < 32 || bytes[i] > 126) {
      str += ".";
    } else {
      str += String.fromCharCode(bytes[i]);
    }
  }
  return str
}

const separatorBytes = new Uint8Array([13, 10, 13, 10]);


function splitBytesBySeparator(bytes, separator) {
  const byteStr = bytes.join(",") + ",";
  const separatorStr = "," + separator.join(",") + ",";
  const chunksStr = byteStr.split(separatorStr);
  return chunksStr.map(chunkStr => {
    const chunkArray = chunkStr.split(",");
    return new Uint8Array(chunkArray.map(byte => parseInt(byte)));
  });
}

function formatInteger(num) {
  const hex = num.toString(16).toUpperCase().padStart(2, '0');
  return hex.padStart(8, '0');
}

function clipCopy(str) {
  ClipboardSetText(str).then(() => {
    ElMessage({
      message: '已成功复制到剪贴板',
      type: 'success',
    })
  })
      .catch((error) => {
        ElMessage({
          message: '无法复制文本到剪贴板',
          type: 'error',
        })
      });
}


let BeyondLength = 3000;
export default {
  props: ['raw', "readOnly", "backcolor", 'Size'],
  watch: {
    raw(RawBytes) {
      this.RawBytes = RawBytes
      this.LineText = {}
      if (this.RawBytes.length > BeyondLength) {
        this.RawBeyond = true
        this.BeyondBytes = this.RawBytes.slice(BeyondLength);
        this.RawBytes = this.RawBytes.slice(0, BeyondLength);
        this.BeyondText = "当前数据太长,仅显示(" + BeyondLength.toLocaleString() + ")字节,剩余(" + this.BeyondBytes.length.toLocaleString() + ")字节未显示"
      } else {
        this.RawBeyond = false
      }
      this.init()
    },
    backcolor(backcolor) {
      this.color.back = backcolor
    },
    Size(w) {
      const scrollTop = this.scroll.Top
      const scrollLeft = this.scroll.Left
      this.Window.Width = w.w
      this.Window.Height = w.h
      this.$nextTick(() => {
        // 设置其他两个 textarea 的滚动位置
        this.$refs.myIndex.scrollTop = scrollTop;
        this.$refs.myIndex.scrollLeft = scrollLeft;

        this.$refs.myHex.scrollTop = scrollTop;
        this.$refs.myHex.scrollLeft = scrollLeft;

        this.$refs.myAnsi.scrollTop = scrollTop;
        this.$refs.myAnsi.scrollLeft = scrollLeft;
      })
    }
  },
  data() {
    return {
      MainHeight: "0px",
      IsHasModify: false,
      radio: "Ansi",
      MenuColor: {
        textColor: "#303133",
        activeTextColor: "#409EFF",
        backgroundColor: "#FFFFFF",
      },
      isContextMenuVisible: false,
      Menu: null,
      RawBytes: [],
      BeyondBytes: [],
      RawBeyond: false,
      BeyondText: "",
      Body: {
        Index: "",
        Hex: "",
        Ansi: "",
        Count: {
          offset: 0,//每行显示几个字节
          header: 0,//Head有几个字节
          Body: 0,//Body有几个字节
        },
      },
      Style: {Index: "", Hex: "", Ansi: ""},
      CSS: {Index: "", Hex: "", Ansi: ""},
      Window: {Width: 0, Height: 0},
      Select: {Hex: {start: -1, end: -1, cmd: null}, Ansi: {start: -1, end: -1}},
      insertHand: {eve: null, Index: -1},
      LineText: {Line: null, LineIndex: 0, RawIndex: 0, HexOffset: 0},
      AnsiUsingInputMethod: false,
      HexUsingInputMethod: false,
      scroll: {
        Top: 0,
        Left: 0,
      },
      color: {
        back: "#ffffff",
        Index: "#000000",
        Header: "#e8e8e8",
        Body: "#030101",
        prompt: "color: rgba(2,2,5,0.93);position: relative;top:3px;color:#000000",
        Select: {
          back: "#5da2fd", text: "#ffffff"
        }
      },
      get theme() {
        return window.Theme.IsDark
      },
      set theme(newValue) {
        window.Theme.IsDark = newValue
      },
    }
  },
  computed: {
    Theme() {
      if (this.theme) {
        this.color.back = "#181818"
        this.color.Index = "#b2b2b2"
        this.color.Header = "#75a7f5"
        this.color.Body = "#dadada"
        this.color.Select.back = "#5da2fd"
        this.color.text = "#0c0c0c"
        this.color.prompt = "color: rgba(2,2,5,0.93);position: relative;top:" + (this.Window.Height - 23) + "px;color:#ffffff"

        this.MenuColor.textColor = "#f3f3f3"
        this.MenuColor.activeTextColor = "#409EFF"
        this.MenuColor.backgroundColor = "#464646"

      } else {
        this.color.back = "#ffffff"
        this.color.Index = "#000000"
        this.color.Header = "#185aea"
        this.color.Body = "#030101"
        this.color.Select.back = "#5da2fd"
        this.color.text = "#ffffff"
        this.color.prompt = "color: rgba(2,2,5,0.93);position: relative;top:" + (this.Window.Height - 23) + "px;color:#000000"

        this.MenuColor.textColor = "#303133"
        this.MenuColor.activeTextColor = "#409EFF"
        this.MenuColor.backgroundColor = "#FFFFFF"
      }
      this.init();
      return false
    },
    Body_Hex_apply() {
      return this.Body.Hex;
    },
    Body_Hex_Select() {
      if (this.Select.Hex.start === insertMode) {
        this.$nextTick(() => {
          this.init()
          this.$nextTick(() => {
            let span = null;
            let offset = this.insertHand.Index
            if (this.insertHand.Header) {
              span = this.$refs.myHex.querySelector('[hid="' + this.insertHand.Line + '"]');
            } else if (this.insertHand.Body1) {
              span = this.$refs.myHex2;
            } else if (this.insertHand.Body2) {
              span = this.$refs.myHex3.querySelector('[hid="' + this.insertHand.Line + '"]');
            }
            const _Range_ = document.createRange();
            const selection = window.getSelection();
            _Range_.setStart(span.childNodes[0], offset);
            _Range_.collapse(true);
            selection.removeAllRanges();
            selection.addRange(_Range_);


            this.$refs.myIndex.scrollTop = this.scroll.Top;
            this.$refs.myIndex.scrollLeft = this.scroll.Left;

            this.$refs.myHex.scrollTop = this.scroll.Top;
            this.$refs.myHex.scrollLeft = this.scroll.Left;

            this.$refs.myAnsi.scrollTop = this.scroll.Top;
            this.$refs.myAnsi.scrollLeft = this.scroll.Left;
          });
        });
        return
      }

      if (this.Body.Count.offset === 0) {
        return
      }
      const my = new MyApply();
      this.$refs.myHex.scrollTop = this.$refs.myAnsi.scrollTop;
      this.$refs.myHex.scrollLeft = this.$refs.myAnsi.scrollLeft;
      let IsContainsCRLFCRLF = (this.Select.Hex.start < this.Body.Count.header && this.Select.Hex.end > this.Body.Count.header)
      if (IsContainsCRLFCRLF) {
        my.SetHex(this.RawBytes.slice(0, this.Select.Hex.start), this.Body.Count.offset, {
          text: this.color.Header,
          back: ""
        })
        my.SetHex(this.RawBytes.slice(this.Select.Hex.start, this.Select.Hex.end), this.Body.Count.offset, {
          text: this.color.Select.text,
          back: this.color.Select.back
        })
        my.SetHex(this.RawBytes.slice(this.Select.Hex.end, this.RawBytes.length), this.Body.Count.offset, {
          text: this.color.Body,
          back: ""
        })
      } else {
        if (this.Select.Hex.start >= this.Body.Count.header) {
          //选择的是Body
          my.SetHex(this.RawBytes.slice(0, this.Body.Count.header), this.Body.Count.offset, {
            text: this.color.Header,
            back: ""
          })
          my.SetHex(this.RawBytes.slice(this.Body.Count.header, this.Select.Hex.start), this.Body.Count.offset, {
            text: this.color.Body,
            back: ""
          })
          my.SetHex(this.RawBytes.slice(this.Select.Hex.start, this.Select.Hex.end), this.Body.Count.offset, {
            text: this.color.Select.text,
            back: this.color.Select.back
          })
          my.SetHex(this.RawBytes.slice(this.Select.Hex.end, this.RawBytes.length), this.Body.Count.offset, {
            text: this.color.Body,
            back: ""
          })
        } else {
          my.SetHex(this.RawBytes.slice(0, this.Select.Hex.start), this.Body.Count.offset, {
            text: this.color.Header,
            back: ""
          })
          my.SetHex(this.RawBytes.slice(this.Select.Hex.start, this.Select.Hex.end), this.Body.Count.offset, {
            text: this.color.Select.text,
            back: this.color.Select.back
          })
          my.SetHex(this.RawBytes.slice(this.Select.Hex.end, this.Body.Count.header), this.Body.Count.offset, {
            text: this.color.Header,
            back: ""
          })
          my.SetHex(this.RawBytes.slice(this.Body.Count.header, this.RawBytes.length), this.Body.Count.offset, {
            text: this.color.Body,
            back: ""
          })
        }
      }
      return my.HexApply()

    },
    Body_Ansi_Apply() {
      return this.Body.Ansi;
    },
    Body_Ansi_Select() {
      const my = new MyApply();
      this.$refs.myAnsi.scrollTop = this.$refs.myHex.scrollTop;
      this.$refs.myAnsi.scrollLeft = this.$refs.myHex.scrollLeft;
      if (this.Select.Ansi.start === insertMode) {
        this.$nextTick(() => {
          this.init()
          this.$nextTick(() => {
            this.getCurrentElementPosition()
          });
        });
        return
      }
      if (this.Body.Count.offset === 0) {
        return
      }
      let IsContainsCRLFCRLF = (this.Select.Ansi.start < this.Body.Count.header && this.Select.Ansi.end > this.Body.Count.header)
      if (IsContainsCRLFCRLF) {
        my.SetAnsi(this.RawBytes.slice(0, this.Select.Ansi.start), this.Body.Count.offset, {
          text: this.color.Header,
          back: ""
        })
        my.SetAnsi(this.RawBytes.slice(this.Select.Ansi.start, this.Select.Ansi.end), this.Body.Count.offset, {
          text: this.color.Select.text,
          back: this.color.Select.back
        })
        my.SetAnsi(this.RawBytes.slice(this.Select.Ansi.end, this.RawBytes.length), this.Body.Count.offset, {
          text: this.color.Body,
          back: ""
        })
      } else {
        if (this.Select.Ansi.start >= this.Body.Count.header) {
          //选择的是Body
          my.SetAnsi(this.RawBytes.slice(0, this.Body.Count.header), this.Body.Count.offset, {
            text: this.color.Header,
            back: ""
          })
          my.SetAnsi(this.RawBytes.slice(this.Body.Count.header, this.Select.Ansi.start), this.Body.Count.offset, {
            text: this.color.Body,
            back: ""
          })
          my.SetAnsi(this.RawBytes.slice(this.Select.Ansi.start, this.Select.Ansi.end), this.Body.Count.offset, {
            text: this.color.Select.text,
            back: this.color.Select.back
          })
          my.SetAnsi(this.RawBytes.slice(this.Select.Ansi.end, this.RawBytes.length), this.Body.Count.offset, {
            text: this.color.Body,
            back: ""
          })
        } else {
          my.SetAnsi(this.RawBytes.slice(0, this.Select.Ansi.start), this.Body.Count.offset, {
            text: this.color.Header,
            back: ""
          })
          my.SetAnsi(this.RawBytes.slice(this.Select.Ansi.start, this.Select.Ansi.end), this.Body.Count.offset, {
            text: this.color.Select.text,
            back: this.color.Select.back
          })
          my.SetAnsi(this.RawBytes.slice(this.Select.Ansi.end, this.Body.Count.header), this.Body.Count.offset, {
            text: this.color.Header,
            back: ""
          })
          my.SetAnsi(this.RawBytes.slice(this.Body.Count.header, this.RawBytes.length), this.Body.Count.offset, {
            text: this.color.Body,
            back: ""
          })
        }
      }
      return my.AnsiApply(this.Body.Count.offset)
    },
    CSS_HexDiv() {
      if (this.RawBeyond === true) {
        return "background-color: " + this.color.back + ";overflow: hidden; width: " + this.Window.Width + "px;height: calc(" + (this.Window.Height - 25) + "px);position: absolute;"
      }
      return "background-color: " + this.color.back + ";overflow: hidden; width: " + this.Window.Width + "px;height: " + this.Window.Height + "px;position: absolute;"
    },
    CSS_Index() {
      return "color: " + this.color.Index + ";text-align: center;float: left; resize: none;  overflow: auto;  border: none;  outline: none; pointer-events: none;background-color: " + this.color.back + ";" + this.Style.Index
    },
    CSS_Hex() {
      return "float: left;  resize: none;    border: none;   outline: none;   overflow: auto;   overflow-anchor: none;background-color: " + this.color.back + ";" + this.Style.Hex
    },
    CSS_Ansi() {
      return "float: right;   resize: none;   border: none;   outline: none;   overflow-x: hidden;  overflow-anchor: none;  overflow-y: scroll; none;background-color: " + this.color.back + ";" + this.Style.Ansi
    },
    IsWindows() {
      if (window.Theme) {
        return window.Theme.GOOS === "windows" ? "width: 182px;" : "width: 200px;"
      }
      return "width: 182px;"
    }
  },
  mounted() {
    document.addEventListener("selectionchange", this.handleSelect);
    this.LineText = {}
    const el = this.$refs.myHexDiv
    if (el) {
      document.addEventListener('click', this.menuItemLeftClick);
      document.addEventListener('contextmenu', this.menuItemRightClick);
      const showContextMenu = (event) => {
        event.preventDefault();
        if (event.button === 2) {
          this.isContextMenuVisible = true;
          const x = event.clientX;
          const y = event.clientY;
          this.$nextTick(() => {
            const el = this.$refs.menu
            el.style.left = `${x}px`;
            el.style.top = `${y}px`;
            el.style.position = `fixed`;
            el.style.backgroundColor = `#656565`;
            el.style.zIndex = `9999`;
            this.$nextTick(() => {
              const cr = el.getBoundingClientRect()
              let _x = x;
              let _y = y;
              if (x + cr.width > window.innerWidth) {
                _x = window.innerWidth - cr.width
              }
              if (y + cr.height > window.innerHeight) {
                _y = window.innerHeight - cr.height
              }
              el.style.left = `${_x}px`;
              el.style.top = `${_y}px`;
            });
          });
        }
      };
      el.addEventListener('contextmenu', showContextMenu);
      const rootElement = document.documentElement;
      rootElement.style.setProperty('--el-menu-item-height', '31px');
    }

  },
  beforeDestroy() {
    document.removeEventListener("selectionchange", this.handleSelect);
    document.removeEventListener('click', this.menuItemLeftClick);
    document.removeEventListener('contextmenu', this.menuItemRightClick);
  },
  methods: {
    CopyALL() {
      let mergedBytes = []
      if (this.RawBeyond) {
        const bytes1 = new Uint8Array(this.RawBytes);
        const bytes2 = new Uint8Array(this.BeyondBytes);
        mergedBytes = new Uint8Array(bytes1.length + bytes2.length);
        mergedBytes.set(bytes1, 0);
        mergedBytes.set(bytes2, bytes1.length);
      } else {
        mergedBytes = this.RawBytes
      }
      return Base64BytesEncode(mergedBytes)
    },
    menuItemClicked(option) {
      this.isContextMenuVisible = false;
      if (option === "all") {
        let mergedBytes = []
        if (this.RawBeyond) {
          const bytes1 = new Uint8Array(this.RawBytes);
          const bytes2 = new Uint8Array(this.BeyondBytes);
          mergedBytes = new Uint8Array(bytes1.length + bytes2.length);
          mergedBytes.set(bytes1, 0);
          mergedBytes.set(bytes2, bytes1.length);
        } else {
          mergedBytes = this.RawBytes
        }
        if (this.radio === "十六进制") {
          const hexString = Array.from(mergedBytes, byte => byte.toString(16).padStart(2, '0')).join(' ').toUpperCase();
          clipCopy(hexString)
        } else if (this.radio === "Ansi") {
          clipCopy(bytesToString(mergedBytes))
        } else {
          clipCopy(Base64BytesEncode(mergedBytes))
        }
      } else if (option === "select") {
        let SelectBody = []
        if (this.Select.Ansi.start !== -1) {
          SelectBody = this.RawBytes.slice(this.Select.Ansi.start, this.Select.Ansi.end)
        } else if (this.Select.Hex.start !== -1) {
          SelectBody = this.RawBytes.slice(this.Select.Hex.start, this.Select.Hex.end)
        }
        if (this.radio === "十六进制") {
          const hexString = Array.from(SelectBody, byte => byte.toString(16).padStart(2, '0')).join(' ').toUpperCase();
          clipCopy(hexString)
        } else if (this.radio === "Ansi") {
          clipCopy(bytesToString(SelectBody))
        } else {
          clipCopy(Base64BytesEncode(SelectBody))
        }
      } else if (option === "body") {
        let mergedBytes = []
        if (this.RawBeyond) {
          const bytes1 = new Uint8Array(this.RawBytes);
          const bytes2 = new Uint8Array(this.BeyondBytes);
          mergedBytes = new Uint8Array(bytes1.length + bytes2.length);
          mergedBytes.set(bytes1, 0);
          mergedBytes.set(bytes2, bytes1.length);
        } else {
          mergedBytes = this.RawBytes
        }
        if (mergedBytes.length <= this.Body.Count.header) {
          ElMessage({
            message: '该请求无body',
            type: 'error',
          })
          return
        }
        mergedBytes = mergedBytes.slice(-(mergedBytes.length - this.Body.Count.header))
        if (this.radio === "十六进制") {
          const hexString = Array.from(mergedBytes, byte => byte.toString(16).padStart(2, '0')).join(' ').toUpperCase();
          clipCopy(hexString)
        } else if (this.radio === "Ansi") {
          clipCopy(bytesToString(mergedBytes))
        } else {
          clipCopy(Base64BytesEncode(mergedBytes))
        }
      }

    },
    menuItemRightClick(event) {
      let ev = event.target
      while (true) {
        if (ev === null || ev === void 0) {
          break
        }
        if (ev === this.$refs.myHexDiv) {
          return
        }
        ev = ev.parentElement
      }
      this.isContextMenuVisible = false;
    },
    menuItemLeftClick(event) {
      let ev = event.target
      while (true) {
        if (ev === null || ev === void 0) {
          break
        }
        if (ev === this.$refs.menu) {
          return
        }
        ev = ev.parentElement
      }
      this.isContextMenuVisible = false;
    },
    Refresh() {
      this.init()
    },
    //初始化
    init() {
      this.Select.Hex.start = -1
      this.Select.Ansi.start = -1
      const selection = window.getSelection();
      selection.removeAllRanges();
      if (this.$refs.myIndex === undefined || this.$refs.myIndex == null) {
        return
      }
      this.$refs.myIndex.innerText = "00000000\r\n00000001"
      if (this.$refs.myIndex.offsetWidth === undefined || this.$refs.myIndex.offsetWidth == null) {
        return
      }
      this.Style.Index = "height: 100%;width: 76px; word-break:break-all;font-family: monospace;font-size: 12px;";
      this.Style.Hex = "height: 100%;width: 364px";
      this.Style.Ansi = "height: 100%;width: 164px";
      let offset = this.adjustWidth()
      if (offset < 1) {
        return
      }
      let i = 0
      let index = "";
      while (i < this.RawBytes.length) {
        if (i === 0) {
          index = formatInteger(i)
        } else {
          index += "\n" + formatInteger(i)
        }
        i += offset
      }
      const array = splitBytesBySeparator(this.RawBytes, separatorBytes);
      let headerLength = 0;
      if (array.length >= 2) {
        headerLength = array[0].length + 4
      }
      this.Body.Index = index
      this.$refs.myIndex.innerText = index;
      this.Body.Count.offset = offset;
      if (offset === 0) {
        debugger
      }
      this.Body.Count.header = headerLength;
      this.Body.Count.Body = this.RawBytes.length - headerLength;
      const myHex = new MyApply();
      const myAnsi = new MyApply();
      myHex.SetHex(this.RawBytes.slice(0, this.Body.Count.header), this.Body.Count.offset, {
        text: this.color.Header,
        back: ""
      })
      myHex.SetHex(this.RawBytes.slice(this.Body.Count.header, this.RawBytes.length), this.Body.Count.offset, {
        text: this.color.Body,
        back: ""
      })
      this.Body.Hex = myHex.HexApply()
      myAnsi.SetAnsi(this.RawBytes.slice(0, this.Body.Count.header), this.Body.Count.offset, {
        text: this.color.Header,
        back: ""
      })
      myAnsi.SetAnsi(this.RawBytes.slice(this.Body.Count.header, this.RawBytes.length), this.Body.Count.offset, {
        text: this.color.Body,
        back: ""
      })
      this.Body.Ansi = myAnsi.AnsiApply()
    },
    //计算调整宽度
    adjustWidth() {
      const IndexWidth = this.$refs.myIndex.offsetWidth;
      let charWidth = window.Theme.GOOS === "windows" ? 6 : 7.5
      let _Width = this.Window.Width - IndexWidth - 25 //Ansi字符后面由滚动条所以减去16+4 显示滚动条
      let charNum = _Width / charWidth //一共可以显示多少个字符
      let AnsiNum = Math.trunc(charNum / 4)//Ansi字符可以显示几个
      if (AnsiNum < 1) {
        this.Style.Hex = "height: 100%;width: 0px;"
        this.Style.Ansi = "height: 100%;width:  0px;"
        return AnsiNum
      }
      let style1 = "";
      let style2 = "";
      let a = _Width - ((AnsiNum * 3 * charWidth) + (AnsiNum * charWidth));
      if (a < Math.ceil(charWidth)) {
        style1 = (AnsiNum * 3 * charWidth) + 4
        style2 = (AnsiNum * charWidth + 20)
      } else {
        style1 = (AnsiNum * 3 * charWidth) + 4
        style2 = (AnsiNum * charWidth + 20);
      }
      let style3 = "px;word-break:break-all;text-align: left;font-family: monospace;font-size: 12px; color: black;";
      this.Style.Hex = "height: 100%;width: " + style1 + style3;
      this.Style.Ansi = "height: 100%;width: " + style2 + style3;
      return AnsiNum
    },  // 滚动条同步处理方法
    syncScroll(event) {
      const scrollTop = event.target.scrollTop;
      const scrollLeft = event.target.scrollLeft;


      this.scroll.Top = scrollTop;
      this.scroll.Left = scrollLeft;


      // 设置其他两个 textarea 的滚动位置
      this.$refs.myIndex.scrollTop = scrollTop;
      this.$refs.myIndex.scrollLeft = scrollLeft;

      this.$refs.myHex.scrollTop = scrollTop;
      this.$refs.myHex.scrollLeft = scrollLeft;

      this.$refs.myAnsi.scrollTop = scrollTop;
      this.$refs.myAnsi.scrollLeft = scrollLeft;
    },
    handleSelect(event) {
      let Bool = false
      let ev = document.elementFromPoint(window.mouseX, window.mouseY);
      while (true) {
        if (ev === null || ev === void 0) {
          break
        }
        if (ev === this.$refs.menu) {
          return
        }
        if (ev === this.$refs.myHexDiv) {
          Bool = true
          break
        }
        ev = ev.parentElement
      }
      this.isContextMenuVisible = false;
      if (Bool === false) {
        return;
      }
      if (this.AnsiUsingInputMethod || this.HexUsingInputMethod) {
        return;
      }
      this.Select.Ansi.start = -1
      this.Select.Ansi.end = -1
      this.Select.Hex.start = -1
      this.Select.Hex.end = -1
      const selection = window.getSelection();
      let range;
      try {
        range = selection.getRangeAt(0);
      } catch (e) {
        return
      }
      let context = this.getElementsIndex(selection, range)
      if (context === null) {
        return;
      }
      this.LineText.Line = range.startContainer
      this.LineText.RawIndex = context.Start
      this.LineText.LineIndex = range.startOffset
      this.LineText.Type = context.Type
      this.LineText.HexOffset = context.HexOffset
      if (context.Start === context.End || context.End < context.Start) {
        return;
      }
      if (context.Type === "hex") {
        this.Select.Ansi.start = context.Start
        this.Select.Ansi.end = context.End
      } else {
        this.Select.Hex.start = context.Start
        this.Select.Hex.end = context.End
      }
    },
    AnsiInput(event) {
      if (this.AnsiUsingInputMethod || this.readOnly) {
        return;
      }
      this.IsHasModify = true
      const selection = window.getSelection();
      const range = selection.getRangeAt(0);
      if (event.command === 8) {
        //剪切或按下删除键
        const DelNum = this.Select.Hex.start === -1 ? 1 : this.Select.Hex.end - this.Select.Hex.start
        const Cut = this.Select.Hex.start !== -1
        if (Cut) {
          this.LineText.RawIndex += 1
        }
        this.RawBytes = deleteBytes(this.RawBytes, this.LineText.RawIndex - 1, DelNum)
        this.LineText.RawIndex -= 1
        this.init()
        this.$nextTick(() => {
          this.getCurrentElementPosition()
        });
        return;
      }
      let diff = null;
      const DelNum = this.Select.Hex.start === -1 ? -1 : this.Select.Hex.end - this.Select.Hex.start
      if (DelNum !== -1) {
        this.RawBytes = deleteBytes(this.RawBytes, this.LineText.RawIndex, DelNum)
      }
      if (event.command === 9) {
        diff = event
      } else {
        let ss = range.startContainer.parentElement.innerText.substring(this.LineText.LineIndex, range.startOffset).replaceAll(" ", " ")

        diff = {index: this.LineText.LineIndex, Char: StringToBytes(ss)}
      }
      this.RawBytes = insertBytes(this.RawBytes, this.LineText.RawIndex, diff.Char)
      this.LineText.RawIndex += diff.Char.length
      this.init()
      this.$nextTick(() => {
        this.getCurrentElementPosition()
      });
    },
    //设置鼠标光标位置
    getCurrentElementPosition() {
      if (this.LineText.Type !== "ansi" && this.LineText.Type !== "hex") {
        return;
      }
      let childIndex = 0
      let Index = 0
      let child = null;
      let Body1Length = this.Body.Count.offset - this.Body.Count.header % this.Body.Count.offset
      if (Body1Length === this.Body.Count.offset) {
        Body1Length = 0
      }
      if (this.LineText.Type === "ansi") {
        if (this.LineText.RawIndex < this.Body.Count.header || Body1Length === 0) {
          childIndex = Math.floor(this.LineText.RawIndex / this.Body.Count.offset)
          Index = this.LineText.RawIndex % this.Body.Count.offset
          child = this.$refs.myAnsi2.children[childIndex]
          if (child === null || child === undefined) {
            child = this.$refs.myAnsi2.children[childIndex - 1]
            Index = this.Body.Count.offset
          }
        } else if (this.LineText.RawIndex <= this.Body.Count.header + Body1Length) {
          if (this.LineText.RawIndex === this.Body.Count.header + Body1Length) {
            childIndex = Math.floor(this.LineText.RawIndex / this.Body.Count.offset) + 1
            Index = 0
            child = this.$refs.myAnsi2.children[childIndex]
            if (child === null || child === undefined) {
              child = this.$refs.myAnsi2.children[childIndex - 1]
              Index = 0
            }
          } else {
            childIndex = Math.floor(this.LineText.RawIndex / this.Body.Count.offset) + 1
            Index = this.LineText.RawIndex - this.Body.Count.header
            child = this.$refs.myAnsi2.children[childIndex]
            if (child === null || child === undefined) {
              child = this.$refs.myAnsi2.children[childIndex - 1]
              Index = Body1Length
            }
          }

        } else {
          childIndex = Math.floor(this.LineText.RawIndex / this.Body.Count.offset) + 1
          Index = this.LineText.RawIndex % this.Body.Count.offset
          child = this.$refs.myAnsi2.children[childIndex]
          if (child === null || child === undefined) {
            child = this.$refs.myAnsi2.children[childIndex - 1]
            Index = this.Body.Count.offset
          }
        }
        if (Index > child.innerText.length) {
          Index = child.innerText.length
        }
      } else if (this.LineText.Type === "hex") {
        if (this.LineText.RawIndex < this.Body.Count.header || Body1Length === 0) {
          childIndex = Math.floor(this.LineText.RawIndex / this.Body.Count.offset)
          Index = (this.LineText.RawIndex % this.Body.Count.offset) * 3 + this.LineText.HexOffset
          child = this.$refs.myHex2.children[childIndex]
          if (child === null || child === undefined) {
            child = this.$refs.myHex2.children[childIndex - 1]
            Index = child.innerText.length
          }
        } else if (this.LineText.RawIndex <= this.Body.Count.header + Body1Length) {
          if (this.LineText.RawIndex === this.Body.Count.header + Body1Length) {
            childIndex = Math.floor(this.LineText.RawIndex / this.Body.Count.offset) + 1
            Index = this.LineText.HexOffset
            child = this.$refs.myHex2.children[childIndex]
            if (child === null || child === undefined) {
              child = this.$refs.myHex2.children[childIndex - 1]
              Index = child.innerText.length
            }
          } else {
            childIndex = Math.floor(this.LineText.RawIndex / this.Body.Count.offset) + 1
            Index = (this.LineText.RawIndex - this.Body.Count.header) * 3 + this.LineText.HexOffset
            child = this.$refs.myHex2.children[childIndex]
            if (child === null || child === undefined) {
              child = this.$refs.myHex2.children[childIndex - 1]
              Index = child.innerText.length
            }
          }
        } else {
          childIndex = Math.floor(this.LineText.RawIndex / this.Body.Count.offset) + 1
          Index = (this.LineText.RawIndex % this.Body.Count.offset) * 3 + this.LineText.HexOffset
          child = this.$refs.myHex2.children[childIndex]
          if (child === null || child === undefined) {
            child = this.$refs.myHex2.children[childIndex - 1]
            Index = child.innerText.length
          }
        }
        if (Index > child.innerText.length) {
          Index = child.innerText.length
        }
      }
      if (child === null || !child.firstChild) {
        return;
      }
      const _Range_ = document.createRange();
      const selection = window.getSelection()
      _Range_.setStart(child.firstChild, Index);
      _Range_.collapse(true);
      selection.removeAllRanges();
      selection.addRange(_Range_);

      this.$refs.myIndex.scrollTop = this.scroll.Top;
      this.$refs.myIndex.scrollLeft = this.scroll.Left;

      this.$refs.myHex.scrollTop = this.scroll.Top;
      this.$refs.myHex.scrollLeft = this.scroll.Left;

      this.$refs.myAnsi.scrollTop = this.scroll.Top;
      this.$refs.myAnsi.scrollLeft = this.scroll.Left;
    }
    ,
    HexInput(event) {
      if (this.HexUsingInputMethod || this.readOnly) {
        return;
      }
      this.IsHasModify = true
      if (event.command === 9) {
        //剪切或按下删除键
        const DelNum = this.Select.Ansi.start === -1 ? 1 : this.Select.Ansi.end - this.Select.Ansi.start
        const Cut = this.Select.Ansi.start !== -1
        if (Cut) {
          this.LineText.RawIndex += 1
        }
        this.RawBytes = deleteBytes(this.RawBytes, this.LineText.RawIndex - 1, DelNum)
        this.LineText.RawIndex -= 1
        this.init()
        this.$nextTick(() => {
          this.getCurrentElementPosition()
        });
        return;
      }
      const DelNum = this.Select.Ansi.start === -1 ? -1 : this.Select.Ansi.end - this.Select.Ansi.start
      if (DelNum !== -1) {
        this.RawBytes = deleteBytes(this.RawBytes, this.LineText.RawIndex, DelNum)
      }
      if (event.command === 10) {
        this.RawBytes = insertBytes(this.RawBytes, this.LineText.RawIndex, event.key)
        this.LineText.HexOffset = 0
        this.LineText.RawIndex += event.key.length
      } else {
        if (this.RawBytes.length === this.LineText.RawIndex) {
          this.RawBytes = insertBytes(this.RawBytes, this.LineText.RawIndex, new Uint8Array(1))
        }
        let _LodRaw = byteToHex(this.RawBytes[this.LineText.RawIndex])
        let _NewRaw = ""
        if (this.LineText.HexOffset === 0) {
          _NewRaw = event.key + _LodRaw[_LodRaw.length - 1];
          this.LineText.HexOffset++
          this.RawBytes[this.LineText.RawIndex] = hexToByte(_NewRaw)
        } else {
          _NewRaw = _LodRaw[0] + event.key;
          this.RawBytes[this.LineText.RawIndex] = hexToByte(_NewRaw)
          this.LineText.RawIndex++
          this.LineText.HexOffset = 0
        }
      }
      this.init()
      this.$nextTick(() => {
        this.getCurrentElementPosition()
      });
    }
    ,
    handleAnsiPaste(event) {
      event.preventDefault(); // 阻止默认粘贴行为
      const clipboardData = event.clipboardData || window.clipboardData;
      const pastedText = clipboardData.getData('text/plain');
      this.AnsiInput({command: 9, Char: StringToBytes(pastedText)})
    }
    ,
    handleHexPaste(event) {
      if (event.command === 10) {
        if (navigator.clipboard) {
          navigator.clipboard.readText()
              .then(pastedText => {
                if (isHexadecimalString(pastedText)) {
                  let a = RemoveSpaces(pastedText)
                  if (a.length % 2 !== 0) {
                    a = a + "0"
                  }
                  this.HexInput({command: 10, key: hexToBytes(a)})
                } else {
                  console.log("粘贴的字符串不是HEX字符串")
                }
              })
              .catch(error => {
                console.error("获取剪贴板内容出错:", error);
              });
        } else {
          console.error("浏览器不支持 Clipboard API");
        }
        return
      }
      event.preventDefault(); // 阻止默认粘贴行为
      const clipboardData = event.clipboardData || window.clipboardData;
      const pastedText = clipboardData.getData('text/plain');
      if (isHexadecimalString(pastedText)) {
        let a = RemoveSpaces(pastedText)
        if (a.length % 2 !== 0) {
          a = a + "0"
        }
        this.HexInput({command: 10, key: hexToBytes(a)})
      } else {
        console.log("粘贴的字符串不是HEX字符串")
      }
    },
    imStart() {
      this.AnsiUsingInputMethod = true;
    },
    imEnd() {
      this.AnsiUsingInputMethod = false;
      this.Select.Ansi.start = insertMode
    },
    imHexStart() {
      this.HexUsingInputMethod = true;
    },
    imHexEnd() {
      this.HexUsingInputMethod = false;
      this.Select.Hex.start = insertMode
    },
    handleAnsiKeyDown(event) {
      if (this.readOnly) {
        event.preventDefault();
        return;
      }
      if (event.key === 'Enter' || event.key === 'Delete') {
        event.preventDefault(); // 阻止回车键默认行为
      }
      if ((event.ctrlKey && event.key === 'z')) {
        //撤销
        event.preventDefault();
      }
      if ((event.ctrlKey && event.key === 'y')) {
        //重做
        event.preventDefault();
      }
      if (event.key === 'Backspace' || (event.ctrlKey && event.key === 'x')) {
        //退格键 或 剪切
        event.preventDefault();
        this.AnsiInput({command: 8})
      }

    }
    ,
    handleMouseDown(event) {
      if (event.button === 2) {
        return
      }
      if (this.Select.Ansi.start !== -1) {
        this.$nextTick(() => {
          this.Select.Hex.start = -1
          this.Select.Ansi.start = -1
          window.getSelection().removeAllRanges()
        })
      }
      if (this.Select.Hex.start !== -1) {
        this.$nextTick(() => {
          this.Select.Hex.start = -1
          this.Select.Ansi.start = -1
          window.getSelection().removeAllRanges()
        })
      }
    },
    handleHexKeyDown(event) {
      if (this.readOnly) {
        event.preventDefault();
        return;
      }
      if ((event.ctrlKey && (event.key === 'v' || event.key === 'V'))) {
        event.preventDefault();
        this.handleHexPaste({command: 10})
        return;
      }
      if ((event.ctrlKey && (event.key === 'c' || event.key === 'C'))) {
        event.preventDefault();
        return;
      }
      const mKey = ["ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight"]
      if (mKey.includes(event.key)) {
        return
      }
      const key = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'A', 'B', 'C', 'D', 'E', 'F']
      if ((event.ctrlKey && event.key === 'z')) {
        //撤销
        event.preventDefault();
      }
      if ((event.ctrlKey && event.key === 'y')) {
        //重做
        event.preventDefault();
      }
      if (event.key === 'Backspace' || (event.ctrlKey && event.key === 'x')) {
        //退格键 或 剪切
        event.preventDefault();
        this.HexInput({command: 9})
      }
      if (!key.includes(event.key)) {
        event.preventDefault();
      } else {
        const eve = event.key
        event.preventDefault();
        this.$nextTick(() => {
          this.HexInput({command: 7, key: eve})
        });
      }
    }
    ,
    getElementsIndex(selection, range) {
      let index1 = range.startOffset
      let index2 = range.endOffset
      let StartElement = range.startContainer.parentElement;
      if (StartElement.tagName !== "SPAN") {
        StartElement = range.startContainer.nextElementSibling;
        index1 = 0
      }
      let StartIndex = 0;
      try {
        StartIndex = Array.from(StartElement.parentNode.children).indexOf(StartElement);
      } catch (e) {
        return null
      }
      let endElement = range.endContainer.parentElement;
      if (endElement.tagName !== "SPAN") {
        endElement = range.endContainer.nextElementSibling;
        index2 = 0
      }
      let endIndex = Array.from(endElement.parentNode.children).indexOf(endElement);
      let Type = StartElement.getAttribute('tid')
      let Select = {Type: "", Start: 0, End: 0, Next: range.startContainer.nextElementSibling, HexOffset: 0}
      let Body1Line = Math.floor(this.Body.Count.header / this.Body.Count.offset)
      let Body1Length = -1
      if (this.Body.Count.header !== 0) {
        Body1Length = this.Body.Count.offset - (this.Body.Count.header % this.Body.Count.offset)
        if (Body1Length === this.Body.Count.offset) {
          Body1Length = 0
        } else {
          Body1Line++
        }
      }
      if (Type === "ansi") {
        Select.Type = "ansi"
      } else if (Type === "hex") {
        Select.Type = "hex"
        let startOffset = index1
        if (index1 % 3 === 1) {
          startOffset--;
          Select.HexOffset = 1
        } else if (index1 % 3 === 2) {
          startOffset++;
        }
        let EndOffset = range.endOffset
        if (range.endOffset % 3 === 1) {
          EndOffset--;
        } else if (range.endOffset % 3 === 2) {
          EndOffset++;
        }
        index1 = Math.ceil(startOffset / 3)
        index2 = Math.ceil(EndOffset / 3)
      } else {
        return null
      }
      if (StartIndex < Body1Line || Body1Length < 1) {
        Select.Start = this.Body.Count.offset * StartIndex + index1
      } else if (StartIndex === Body1Line) {
        Select.Start = this.Body.Count.offset * StartIndex + index1 - Body1Length
      } else if (StartIndex > Body1Line) {
        Select.Start = this.Body.Count.offset * StartIndex + index1 - this.Body.Count.offset
      }
      if (endIndex < Body1Line || Body1Length < 1) {
        Select.End = this.Body.Count.offset * endIndex + index2
      } else if (endIndex === Body1Line) {
        Select.End = this.Body.Count.offset * endIndex + index2 - Body1Length
      } else if (endIndex > Body1Line) {
        Select.End = this.Body.Count.offset * endIndex + index2 - this.Body.Count.offset
      }
      return Select
    }
  }
}

</script>

<template>

  <div ref="myHexDiv" :style="CSS_HexDiv">
    <div ref="myIndex" :style="CSS_Index" class="Hex_Index_style" contenteditable="true"
         @scroll="syncScroll"></div>
    <div ref="myHex" :style="CSS_Hex" class="Hex_Hex_style" contenteditable="true" @scroll="syncScroll"
         @paste="handleHexPaste" @keydown="handleHexKeyDown" @input="HexInput" @compositionstart="imHexStart"
         @compositionend="imHexEnd" @mousedown="handleMouseDown">
      <div ref="myHex2" v-if="Select.Hex.start===-1" v-html="Body_Hex_apply">
      </div>

      <div v-html="Body_Hex_Select" v-if="Select.Hex.start!==-1"></div>
    </div>
    <div ref="myAnsi" :style="CSS_Ansi" contenteditable="true" @scroll="syncScroll"
         @input="AnsiInput" @paste="handleAnsiPaste" @keydown="handleAnsiKeyDown" @compositionstart="imStart"
         @compositionend="imEnd" @mousedown="handleMouseDown">
      <div ref="myAnsi2" v-if="this.Select.Ansi.start === -1" v-html="Body_Ansi_Apply">
      </div>

      <div v-html="Body_Ansi_Select" v-if="Select.Ansi.start!==-1"></div>
    </div>
  </div>
  <span :style="color.prompt"
        v-if="this.RawBeyond">{{ this.BeyondText }}</span>
  <div v-show="isContextMenuVisible" ref="menu">
    <el-menu
        :active-text-color="MenuColor.activeTextColor"
        :background-color="MenuColor.backgroundColor"
        :text-color="MenuColor.textColor"
        class="el-menu-vertical-demo"
        default-active="2"
        text-color="#fff"
        :style="IsWindows"
    >
      <el-radio-group v-model="radio" size="small">
        <el-radio-button label="十六进制"/>
        <el-radio-button label="Base64"/>
        <el-radio-button label="Ansi"/>
      </el-radio-group>

      <el-menu-item @click="menuItemClicked( 'all')">复制全部</el-menu-item>
      <el-menu-item @click="menuItemClicked( 'select')">复制选中</el-menu-item>
      <el-menu-item @click="menuItemClicked( 'body')">复制全部Body</el-menu-item>

    </el-menu>

  </div>
  <div v-if="Theme"></div>

</template>

<style scoped>
/*-----------------------------------------------------------------------------------------------*/

/* 隐藏滚动条 */
.Hex_Index_style::-webkit-scrollbar {
  width: 0;
  background-color: transparent;
}

/* 隐藏滚动条 */
.Hex_Hex_style::-webkit-scrollbar {
  width: 0;
  background-color: transparent;
}


textarea:focus {
  outline: none;
  border-color: transparent;
}


</style>
