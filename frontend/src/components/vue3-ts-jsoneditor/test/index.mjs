import { defineComponent, inject, ref, computed, watch, onMounted, nextTick, onBeforeUnmount, openBlock, createElementBlock, normalizeClass, normalizeStyle, withModifiers } from "vue";
const propNames = [
  "mode",
  "mainMenuBar",
  "navigationBar",
  "statusBar",
  "readOnly",
  "indentation",
  "tabSize",
  "escapeControlCharacters",
  "escapeUnicodeCharacters",
  "validator",
  "queryLanguages",
  "queryLanguageId",
  "onClassName",
  "onRenderValue",
  "onRenderMenu"
];
const pickDefinedProps = (options = {}, props) => {
  const computedProps = {};
  for (const propName of propNames) {
    const prop = props[propName] !== void 0 ? props[propName] : options[propName];
    if (prop !== void 0) {
      computedProps[propName] = prop;
    }
  }
  return computedProps;
};
const fullWidthIcon = `
  <svg 
    class="fa-icon svelte-1dof0an" 
    viewBox="0 0 1024 1024" 
    version="1.1"
     xmlns="http://www.w3.org/2000/svg" 
     p-id="1927" xmlns:xlink="http://www.w3.org/1999/xlink" 
     width="24" 
     height="24"
  >
    <path d="M63.989383 105.442494l0 268.396843c0 18.935258 15.368012 34.304294 34.304294 34.304294 18.936281 0 34.304294-15.369036 34.304294-34.304294L132.597971 180.156126l218.107483 218.176045c12.82919 12.830213 33.618679 12.830213 46.515407 0 12.830213-12.897751 12.830213-33.686217 0-46.51643l-218.176045-218.107483 193.683211 0c18.935258 0 34.304294-15.369036 34.304294-34.304294 0-18.935258-15.369036-34.304294-34.304294-34.304294L104.331183 65.09967C79.288834 65.09967 63.989383 77.999468 63.989383 105.442494L63.989383 105.442494z" p-id="1928" fill="#e6e6e6"></path><path d="M917.688719 65.09967 649.290853 65.09967c-18.935258 0-34.304294 15.369036-34.304294 34.304294 0 18.936281 15.369036 34.304294 34.304294 34.304294l193.683211 0-218.176045 218.107483c-12.830213 12.82919-12.830213 33.618679 0 46.51643 12.897751 12.830213 33.686217 12.830213 46.515407 0L889.420909 180.156126l0 193.683211c0 18.935258 15.369036 34.304294 34.304294 34.304294 18.936281 0 34.304294-15.369036 34.304294-34.304294L958.029496 105.442494C958.029496 77.999468 942.79963 65.09967 917.688719 65.09967L917.688719 65.09967z" p-id="1929" fill="#e6e6e6"></path>
    <path d="M104.331183 957.013353l268.397866 0c18.935258 0 34.304294-15.368012 34.304294-34.304294 0-18.936281-15.369036-34.304294-34.304294-34.304294L179.045839 888.404766l218.176045-218.107483c12.830213-12.82919 12.830213-33.618679 0-46.515407-12.897751-12.830213-33.686217-12.830213-46.515407 0l-218.107483 218.176045L132.598994 648.27471c0-18.935258-15.368012-34.304294-34.304294-34.304294-18.936281 0-34.304294 15.369036-34.304294 34.304294l0 268.397866C63.989383 944.115602 79.288834 957.013353 104.331183 957.013353L104.331183 957.013353z" p-id="1930" fill="#e6e6e6"></path>
    <path d="M958.029496 916.671553 958.029496 648.27471c0-18.935258-15.368012-34.304294-34.304294-34.304294-18.935258 0-34.304294 15.369036-34.304294 34.304294l0 193.683211L671.313425 623.781876c-12.82919-12.830213-33.618679-12.830213-46.515407 0-12.830213 12.897751-12.830213 33.686217 0 46.515407l218.176045 218.107483L649.290853 888.404766c-18.935258 0-34.304294 15.368012-34.304294 34.304294 0 18.936281 15.369036 34.304294 34.304294 34.304294l268.397866 0C942.79963 957.013353 958.029496 944.115602 958.029496 916.671553L958.029496 916.671553z" p-id="1931" fill="#e6e6e6"></path>
  </svg>
`;
var jseThemeDark = /* @__PURE__ */ (() => ".jse-theme-dark{--jse-theme-color: #2f6dd0;--jse-theme-color-highlight: #467cd2;--jse-background-color: #1e1e1e;--jse-text-color: #d4d4d4;--jse-main-border: 1px solid #4f4f4f;--jse-menu-color: #fff;--jse-modal-background: #2f2f2f;--jse-modal-overlay-background: rgba(0, 0, 0, .5);--jse-modal-code-background: #2f2f2f;--jse-panel-background: #333333;--jse-panel-background-border: 1px solid #464646;--jse-panel-color: var(--jse-text-color);--jse-panel-color-readonly: #737373;--jse-panel-border: 1px solid #3c3c3c;--jse-panel-button-color-highlight: #e5e5e5;--jse-panel-button-background-highlight: #464646;--jse-navigation-bar-background: #656565;--jse-navigation-bar-background-highlight: #7e7e7e;--jse-navigation-bar-dropdown-color: var(--jse-text-color);--jse-context-menu-background: #4b4b4b;--jse-context-menu-background-highlight: #595959;--jse-context-menu-separator-color: #595959;--jse-context-menu-color: var(--jse-text-color);--jse-context-menu-button-background: #737373;--jse-context-menu-button-background-highlight: #818181;--jse-context-menu-button-color: var(--jse-context-menu-color);--jse-key-color: #9cdcfe;--jse-value-color: var(--jse-text-color);--jse-value-color-number: #b5cea8;--jse-value-color-boolean: #569cd6;--jse-value-color-null: #569cd6;--jse-value-color-string: #ce9178;--jse-value-color-url: #ce9178;--jse-delimiter-color: #949494;--jse-edit-outline: 2px solid var(--jse-text-color);--jse-selection-background-color: #464646;--jse-selection-background-light-color: #333333;--jse-hover-background-color: #343434;--jse-collapsed-items-background-color: #333333;--jse-collapsed-items-selected-background-color: #565656;--jse-collapsed-items-link-color: #b2b2b2;--jse-collapsed-items-link-color-highlight: #ec8477;--jse-search-match-color: #724c27;--jse-search-match-outline: 1px solid #966535;--jse-search-match-active-color: #9f6c39;--jse-search-match-active-outline: 1px solid #bb7f43;--jse-tag-background: #444444;--jse-tag-color: #bdbdbd;--jse-input-background: #3d3d3d;--jse-input-border: var(--jse-main-border);--jse-button-background: #808080;--jse-button-background-highlight: #7a7a7a;--jse-button-color: #e0e0e0;--jse-a-color: #55abff;--jse-a-color-highlight: #4387c9;--background: #3d3d3d;--border: 1px solid #4f4f4f;--listBackground: #3d3d3d;--itemHoverBG: #505050;--multiItemBG: #5b5b5b;--inputColor: #d4d4d4;--multiClearBG: #8a8a8a;--listShadow: 0 2px 6px 0 rgba(0, 0, 0, .24);--jse-color-picker-background: #656565;--jse-color-picker-border-box-shadow: #8c8c8c 0 0 0 1px}\n")();
var JsonEditor_vue_vue_type_style_index_0_lang = /* @__PURE__ */ (() => ".vue-ts-json-editor{min-width:300px;width:100%}.vue-ts-json-editor--max-box{position:fixed;top:0;left:0;width:100vw;height:100vh;z-index:10000}.vue-ts-json-editor .jse-menu .jse-full-width{display:flex}.vue-ts-json-editor .jse-menu .jse-full-width--active{background-color:#ffffff38!important;border-color:#fff9!important}\n")();
var _export_sfc = (sfc, props) => {
  const target = sfc.__vccOpts || sfc;
  for (const [key, val] of props) {
    target[key] = val;
  }
  return target;
};
const _sfc_main = defineComponent({
  name: "JsonEditor",
  props: {
    json: [Object, Array, Number, String, Boolean],
    jsonString: String,
    mode: {
      type: String,
      default: "tree",
      validator: (value) => ["tree", "text"].includes(value)
    },
    mainMenuBar: {
      type: Boolean,
      default: void 0
    },
    navigationBar: {
      type: Boolean,
      default: void 0
    },
    statusBar: {
      type: Boolean,
      default: void 0
    },
    readOnly: {
      type: Boolean,
      default: void 0
    },
    indentation: [String, Number],
    tabSize: Number,
    escapeControlCharacters: {
      type: Boolean,
      default: void 0
    },
    escapeUnicodeCharacters: {
      type: Boolean,
      default: void 0
    },
    validator: Function,
    queryLanguages: Array,
    queryLanguageId: String,
    onClassName: Function,
    onRenderValue: Function,
    onRenderMenu: Function,
    height: [String, Number],
    fullWidthButton: {
      type: Boolean,
      default: void 0
    },
    darkTheme: {
      type: Boolean,
      default: void 0
    }
  },
  emits: [
    "update:json",
    "update:jsonString",
    "change",
    "error",
    "change-mode",
    "change-query-language",
    "focus",
    "blur"
  ],
  setup(props, { expose, emit }) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _i, _j, _k, _l;
    const pluginOptions = inject("jsonEditorOptions", {});
    const max = ref(false);
    const container = ref();
    const fullWidthButton = ref(null);
    const editor = ref(null);
    const getHeight = computed(() => {
      const height = props.height || (pluginOptions == null ? void 0 : pluginOptions.height);
      if (height && !max.value) {
        return {
          height: height + "px"
        };
      }
      return {};
    });
    const content = computed(() => {
      return {
        json: props.json,
        text: props.jsonString
      };
    });
    const darkThemeStyle = computed(() => {
      return props.darkTheme || (pluginOptions == null ? void 0 : pluginOptions.darkTheme);
    });
    const removeFullWidthButton = () => {
      if (!fullWidthButton.value)
        return;
      fullWidthButton.value.removeEventListener("click", onButtonClick);
      fullWidthButton.value = null;
    };
    const setFullWidthButton = async () => {
      if (typeof window === "undefined")
        return;
      const oldButton = window == null ? void 0 : window.document.querySelector(".jse-full-width");
      const pluginOptionFlag = (pluginOptions == null ? void 0 : pluginOptions.fullWidthButton) !== void 0 ? pluginOptions == null ? void 0 : pluginOptions.fullWidthButton : true;
      const fullWidthButtonFlag = props.fullWidthButton !== void 0 ? props.fullWidthButton : pluginOptionFlag;
      if (!fullWidthButtonFlag || oldButton)
        return;
      if (fullWidthButton.value) {
        removeFullWidthButton();
      }
      const menu = window == null ? void 0 : window.document.querySelector(".jse-menu");
      fullWidthButton.value = window == null ? void 0 : window.document.createElement("button");
      fullWidthButton.value.classList.add("jse-full-width");
      fullWidthButton.value.classList.add("jse-button");
      fullWidthButton.value.classList.add("svelte-v4jelk");
      fullWidthButton.value.innerHTML += fullWidthIcon;
      menu.appendChild(fullWidthButton.value);
      fullWidthButton.value.addEventListener("click", onButtonClick);
    };
    const onButtonClick = () => {
      var _a2, _b2;
      max.value = !max.value;
      if (max.value) {
        (_a2 = fullWidthButton.value) == null ? void 0 : _a2.classList.add("jse-full-width--active");
      } else {
        (_b2 = fullWidthButton.value) == null ? void 0 : _b2.classList.remove("jse-full-width--active");
      }
    };
    const expandCollapseAll = (value) => {
      var _a2;
      if (props.mode === "text")
        return;
      (_a2 = editor.value) == null ? void 0 : _a2.expand(() => value);
    };
    const onChange = (content2, previousContent, patchResult) => {
      if (!!content2.json) {
        emit("update:json", content2.json);
      }
      if (!!content2.text) {
        emit("update:jsonString", content2.text);
      }
      emit("change", content2, previousContent, patchResult);
    };
    const onError = (err) => {
      emit("error", err);
    };
    const onChangeMode = (mode) => {
      emit("change-mode", mode);
    };
    const onChangeQueryLanguage = (queryLanguageId) => {
      emit("change-query-language", queryLanguageId);
    };
    const onFocus = () => {
      emit("focus");
    };
    const onBlur = () => {
      emit("blur");
    };
    const onRenderMenu = (mode, items) => {
      nextTick(() => {
        setFullWidthButton();
      });
      if (typeof props.onRenderMenu === "function") {
        return props.onRenderMenu(mode, items);
      }
      return items;
    };
    const makeEditorProps = () => {
      const options = { fullWidthButton: true, ...pluginOptions || {} };
      return {
        ...pickDefinedProps(options, props),
        content: content.value,
        onChange,
        onError,
        onChangeMode,
        onChangeQueryLanguage,
        onFocus,
        onBlur,
        onRenderMenu
      };
    };
    const initView = async () => {
      if (typeof window === "undefined")
        return;
      if (!editor.value) {
        const { JSONEditor } = await import("./index.js");
        editor.value = new JSONEditor({
          target: container.value,
          props: makeEditorProps()
        });
      }
      editor.value.focus();
    };
    const updateProps = () => {
      editor.value.updateProps(makeEditorProps());
    };
    const updateContent = () => {
      editor.value.set(content.value);
    };
    const destroyView = () => {
      if (editor.value) {
        editor.value.destroy();
        editor.value = null;
      }
      removeFullWidthButton();
    };
    watch(props, updateProps);
    watch(() => props.json, updateContent, { deep: true });
    watch(() => props.jsonString, updateContent);
    onMounted(() => {
      nextTick(() => {
        initView();
      });
    });
    onBeforeUnmount(() => {
      destroyView();
    });
    expose({
      $collapseAll() {
        expandCollapseAll(false);
      },
      $expandAll() {
        expandCollapseAll(true);
      },
      $expand: (_a = editor.value) == null ? void 0 : _a.expand,
      $get: (_b = editor.value) == null ? void 0 : _b.get,
      $set: (_c = editor.value) == null ? void 0 : _c.set,
      $update: (_d = editor.value) == null ? void 0 : _d.update,
      $updateProps: (_e = editor.value) == null ? void 0 : _e.updateProps,
      $refresh: (_f = editor.value) == null ? void 0 : _f.refresh,
      $focus: (_g = editor.value) == null ? void 0 : _g.focus,
      $patch: (_h = editor.value) == null ? void 0 : _h.patch,
      $transform: (_i = editor.value) == null ? void 0 : _i.transform,
      $scrollTo: (_j = editor.value) == null ? void 0 : _j.scrollTo,
      $findElement: (_k = editor.value) == null ? void 0 : _k.findElement,
      $acceptAutoRepair: (_l = editor.value) == null ? void 0 : _l.acceptAutoRepair
    });
    return {
      max,
      getHeight,
      container,
      content,
      darkThemeStyle
    };
  }
});
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return openBlock(), createElementBlock("div", {
    class: normalizeClass(["vue-ts-json-editor", { "vue-ts-json-editor--max-box": _ctx.max, "jse-theme-dark": _ctx.darkThemeStyle }]),
    style: normalizeStyle(_ctx.getHeight),
    ref: "container",
    onKeydown: _cache[0] || (_cache[0] = withModifiers(() => {
    }, ["stop"]))
  }, null, 38);
}
var VueJsonEditor = /* @__PURE__ */ _export_sfc(_sfc_main, [["render", _sfc_render]]);
var JsonEditorPlugin = {
  ...VueJsonEditor,
  install: (app, params = {}) => {
    app.component(params.componentName || "JsonEditor", VueJsonEditor);
    app.provide("jsonEditorOptions", params.options);
  }
};
export { JsonEditorPlugin as default };
