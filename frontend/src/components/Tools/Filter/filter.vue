<script>
import {AgGridVue} from "ag-grid-vue3";
import {
  Config_Find_Range_ALL,
  Config_Find_Window,
  Config_Find_Window_Show,
  Config_Focus_Element,
  Config_Theme_agGrid
} from "../../config/Config.js";
import {AG_GRID_LOCALE_CN} from "../../config/AG_ZH_CN.js";
import {ElMessage} from "element-plus";
import {Tour_Add} from "../../Home/Tour";

export default {
  props: ["Name", "Parent", "column", "apply", "SearchDone", "CancelSearch"],
  components: {'ag-grid-vue': AgGridVue},
  computed: {
    agTheme() {
      return Config_Theme_agGrid.value
    },
  },
  data() {
    return {
      ep: null,
      agGridApi: null,
      Element: null,
      cancelElement: null,
      ParentElement: null,
      defaultColDef: {},
      rowData: [
        {"#": 1, "长度": 2, "序号": 3, "发送长度": 4, "响应长度": 5},//这里只定义数字类型的列,未写的都是文本类型
      ],
      gridOptions: {
        enableAdvancedFilter: true,
        localeText: AG_GRID_LOCALE_CN,
        columnDefs: [],
        rowSelection: {
          mode: 'multiRow',
          checkboxes: false,
          headerCheckbox: false,
          copySelectedRows: false,
          enableClickSelection: true,
        },
        isExternalFilterPresent: this.isExternalFilterPresent,
        doesExternalFilterPass: this.doesExternalFilterPass,
        onFilterChanged: this.onFilterChanged,
      },
    }
  },
  watch: {
    Parent(n, l) {
      this.ParentElement = n;
      if (this.Element !== null) {
        this.init();
      }
    },
    column(n, l) {
      const array = [];
      for (let i = 0; i < n.length; i++) {
        const ms = n[i]
        if ((n[i].headerName + "").indexOf("长度") !== -1) {
          ms["filterParams"] = {filterOptions: ["equals", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"]};
        }
        array.push(ms)
      }
      this.agGridApi.setGridOption("columnDefs", array)
      this.agGridApi.refreshHeader()
    },
  },
  methods: {
    init() {
      this.ParentElement.appendChild(this.Element)
      this.agGridApi.setGridOption("popupParent", document.getElementById("appMain"))
      this.agGridApi.setGridOption("columnDefs", this.column)
      this.ParentElement.appendChild(this.cancelElement)
    },
    setFocusElement() {
      const m = this.Element.querySelector("[data-ref=\"eInput\"]")
      Config_Focus_Element.value = m;
      return m;
    },
    setFilter(Model) {
      try {
        this.agGridApi.setAdvancedFilterModel(JSON.parse(Model));
      } catch (e) {
        this.agGridApi.setAdvancedFilterModel(null);
      }
    },
    SetInput(s) {
      if (this.ep) {
        this.ep.value = s;
        this.ep.dispatchEvent(new Event('input', {bubbles: true}));
        const start = performance.now();
        const tick = () => {
          // 每帧都查一次
          const el = document.querySelector(".ag-autocomplete-list");
          if (el) {
            el.style.display = "none"; // 出现了就隐藏
            return; // 停止
          }
          // 1秒内持续检查
          if (performance.now() - start < 1000) {
            requestAnimationFrame(tick);
          }
        };
        requestAnimationFrame(tick);
      }
    },
    //点击高级搜索
    onClickSearchButton() {
      if (Config_Find_Window.value.isSearchInProgress) {
        ElMessage.error('正在查找其他值,请稍后..')
        return
      }
      Config_Find_Range_ALL.value = this.Name.indexOf("主列表") !== -1;
      Config_Find_Window.value.Title = "🚀 [ " + this.Name + " ] 高级搜索";
      Config_Find_Window.value.CompleteCallback = this.agGridApi.ShowCancelSearch
      Config_Find_Window.value.SearchDone = this.SearchDone
      Config_Find_Window_Show.value()
    },
    onGridReady(params) {
      //向高级过滤器添加一个按钮
      {
        const customButtonContainer = this.$el.getElementsByClassName('ag-advanced-filter ag-advanced-filter-header-cell');
        if (customButtonContainer.length > 0) {
          if (customButtonContainer[0]) {
            const popup = customButtonContainer[0];
            const eApplyFilterButton = popup.querySelector("[data-ref=\"eApplyFilterButton\"]")
            eApplyFilterButton.style.display = "none"
            const eInput = popup.querySelector("[data-ref=\"eInput\"]")
            {

              let isDeleteKey = false;
              const applyButton = (eve) => {
                const Model = eInput.value + "";
                if (Model.trim() === '') {
                  this.apply(null)
                  eApplyFilterButton.click()
                  return
                }
                if (!eApplyFilterButton.disabled) {
                  eApplyFilterButton.click()
                }
              }
              eInput.addEventListener('keydown', (event) => {
                isDeleteKey = event.key === 'Backspace' || event.key === 'Delete';
              });
              eInput.addEventListener("focus", () => {
                const stack = (new Error()).stack + "";
                if (stack.includes("dispatchEvent") && stack.includes("init@")) {
                  return
                }
                const val = eInput.value + ""
                if (val.trim() === "") {
                  if (this.Name === "主列表") {
                    eInput.value = "[全部数据] 包含 \"\""
                  } else {
                    eInput.value = "[数据] 包含 \"\""
                  }
                  eInput.dispatchEvent(new Event('input', {bubbles: true}));
                }
                let ni = 0;
                const ff = () => {
                  const startPos = eInput.value.lastIndexOf("\"");
                  eInput.setSelectionRange(startPos, startPos, startPos);
                  if (ni < 20) {
                    ni++;
                    requestAnimationFrame(ff);
                  }
                }
                requestAnimationFrame(ff);
              });
              eInput.addEventListener("blur", () => {
                if (this.Name === "主列表") {
                  if (eInput.value === "[全部数据] 包含 \"\"") {
                    eInput.value = "";
                    eInput.dispatchEvent(new Event('input', {bubbles: true}));
                  }
                } else {
                  if (eInput.value === "[数据] 包含 \"\"") {
                    eInput.value = "";
                    eInput.dispatchEvent(new Event('input', {bubbles: true}));
                  }
                }
              });
              eInput.addEventListener('input', applyButton);
              {
                //监视“应用”按钮是否可点击，如果是可以点击，则点击一下
                const observer = new MutationObserver((mutations) => {
                  mutations.forEach((mutation) => {
                    if (mutation.attributeName === "disabled") {
                      if (!eApplyFilterButton.disabled) {
                        applyButton()
                      } else {
                        if (eInput.value.indexOf("\"\"") !== -1) {
                          this.apply(null)
                          eApplyFilterButton.click()
                        }
                      }
                    }
                  });
                });
                observer.observe(eApplyFilterButton, {attributes: true});
              }
              this.ep = eInput;
              //确保输入的中文引号可用，应用设置规则时，如果是中文输入法，输入的是中文引号，肉眼不方便观察到,所以将中午引号替换为英文引号
              const nativeValueDescriptor = Object.getOwnPropertyDescriptor(HTMLInputElement.prototype, "value");
              Object.defineProperty(eInput, "value", {
                get() {
                  return (nativeValueDescriptor.get.call(this) + "").replaceAll("“", "\"").replaceAll("”", "\"")
                },
                set(value) {
                  nativeValueDescriptor.set.call(this, value);
                }
              });
            }
            const cButton = document.createElement("button");
            cButton.className = "ag-button ag-standard-button ag-advanced-filter-apply-button";
            cButton.innerHTML = `取消搜索`;
            cButton.style.width = "calc(100% - 20px)";
            cButton.style.color = "crimson";
            cButton.addEventListener("click", () => {
              this.agGridApi.HideCancelSearch()
              this.CancelSearch()
            });
            cButton.style.display = 'none'
            this.cancelElement = cButton
            //popup.parentElement.appendChild(cButton);
            const customButton = document.createElement("button");
            customButton.className = "ag-button ag-advanced-filter-builder-button";
            customButton.innerHTML = `<span aria-hidden="true"><span class="ag-icon ag-icon-group"></span></span><span class="ag-advanced-filter-builder-button-label">高级搜索</span>`;
            if (this.Name === "主列表") {
              const eBuilderFilterButton = popup.querySelector("[data-ref=\"eBuilderFilterButton\"]")

              Tour_Add(eInput, 8, "数据过滤/搜索", "点击这里可以进行快速的数据过滤\n只显示你想看到的")
              if (eBuilderFilterButton) {
                Tour_Add(eBuilderFilterButton, 9, "构建数据过滤表达式", "点击这里可以进行构建数据过滤表达式，更直观的设置过滤逻辑")
              }
              Tour_Add(customButton, 10, "高级搜索", "点击这里可以进行高级搜索\n搜索条件更多")
            }
            customButton.addEventListener("click", () => {
              try {
                if (this.onClickSearchButton) {
                  this.onClickSearchButton()
                }
              } catch (e) {
              }
            });
            // 插入到弹窗内
            popup.appendChild(customButton);
            // 隐藏显示 取消搜索 按钮
            this.agGridApi.HideCancelSearch = () => {
              popup.style.display = ''
              cButton.style.display = 'none'
            }
            //显示 取消搜索 按钮
            this.agGridApi.ShowCancelSearch = () => {
              popup.style.display = 'none'
              cButton.style.display = ''
            }
            this.Element = popup;
            if (this.ParentElement !== null) {
              this.$nextTick(() => {
                this.init()
              })
            }
          }
        }
      }
    },
    onFilterChanged(params) {
      this.apply(params.api.getAdvancedFilterModel())
    },
    isExternalFilterPresent(params) {
      return true;
    },
    doesExternalFilterPass(node) {
      return true;
    },
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.api;
  }
}
</script>

<template>
  <div>
    <!-- 本文件单独提取ag-grid的高级过滤器组件，用来充当过滤器窗口组件 -->
    <!-- 本文件单独提取ag-grid的高级过滤器组件，用来充当过滤器窗口组件 -->
    <!-- 本文件单独提取ag-grid的高级过滤器组件，用来充当过滤器窗口组件 -->
    <div style="width: 0;height: 0;position: absolute" id="1234567891">
      <ag-grid-vue ref="agGrid"
                   :theme="agTheme"
                   style="height: 100%;"
                   :rowData="rowData"
                   :onGridReady="onGridReady"
                   :grid-options="gridOptions"
                   :allowContextMenuWithControlKey="true"
                   :defaultColDef="defaultColDef"
                   :suppressCutToClipboard="true"
      />
    </div>
    <!-- 本文件单独提取ag-grid的高级过滤器组件，用来充当过滤器窗口组件 -->
    <!-- 本文件单独提取ag-grid的高级过滤器组件，用来充当过滤器窗口组件 -->
    <!-- 本文件单独提取ag-grid的高级过滤器组件，用来充当过滤器窗口组件 -->
  </div>
</template>
