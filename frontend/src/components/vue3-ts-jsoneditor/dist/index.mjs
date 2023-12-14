import { defineComponent, inject, ref, computed, reactive, watch, onMounted, nextTick, onBeforeUnmount, openBlock, createElementBlock, normalizeClass, normalizeStyle, withModifiers, renderSlot, createCommentVNode } from "vue";
const watchPropNames = [
  "mainMenuBar",
  "navigationBar",
  "statusBar",
  "askToFormat",
  "readOnly",
  "indentation",
  "tabSize",
  "escapeControlCharacters",
  "escapeUnicodeCharacters",
  "flattenColumns",
  "validator",
  "onClassName",
  "onRenderValue",
  "onRenderMenu"
];
const propNames = [
  "mode",
  "mainMenuBar",
  "navigationBar",
  "statusBar",
  "askToFormat",
  "readOnly",
  "indentation",
  "tabSize",
  "escapeControlCharacters",
  "escapeUnicodeCharacters",
  "flattenColumns",
  "validator",
  "parser",
  "validationParser",
  "pathParser",
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
    modelValue: [Object, Array, Number, String, Boolean, String, null],
    value: [Object, Array, Number, String, Boolean, String, null],
    json: [Object, Array, Number, String, Boolean, null],
    text: String,
    jsonString: String,
    mode: {
      type: String,
      default: "tree",
      validator: (value) => ["tree", "text", "table"].includes(value)
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
    askToFormat: {
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
    flattenColumns: {
      type: Boolean,
      default: void 0
    },
    validator: Function,
    parser: Object,
    validationParser: Object,
    pathParser: Object,
    queryLanguagesIds: Array,
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
    "update:modelValue",
    "update:json",
    "update:text",
    "update:jsonString",
    "change",
    "error",
    "change-mode",
    "update:mode",
    "change-query-language",
    "focus",
    "blur"
  ],
  setup(props, { expose, emit }) {
    const pluginOptions = inject("jsonEditorOptions", {});
    const container = ref();
    const fullWidthButton = ref(null);
    const max = ref(false);
    const blockUpdate = ref(false);
    const blockChange = ref(false);
    const mode = ref("tree");
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
    const darkThemeStyle = computed(() => {
      return props.darkTheme || (pluginOptions == null ? void 0 : pluginOptions.darkTheme);
    });
    const queryLanguagesIds = computed(() => {
      return props.queryLanguagesIds || (pluginOptions == null ? void 0 : pluginOptions.queryLanguagesIds);
    });
    const queryLanguageId = computed(() => {
      return props.queryLanguageId || (pluginOptions == null ? void 0 : pluginOptions.queryLanguageId);
    });
    const queryLanguagesBuffer = reactive({});
    const makeQueryLanguages = async () => {
      var _a;
      if (typeof window === "undefined" || typeof queryLanguagesIds.value === "undefined" || !((_a = queryLanguagesIds.value) == null ? void 0 : _a.length)) {
        return;
      }
      for (const languageId of queryLanguagesIds.value) {
        if (!queryLanguagesBuffer[languageId]) {
          switch (languageId) {
            case "javascript": {
              const { javascriptQueryLanguage } = await import("./vanilla-jsoneditor.js");
              queryLanguagesBuffer[languageId] = javascriptQueryLanguage;
              break;
            }
            case "lodash": {
              const { lodashQueryLanguage } = await import("./vanilla-jsoneditor.js");
              queryLanguagesBuffer[languageId] = lodashQueryLanguage;
              break;
            }
            case "jmespath": {
              const { jmespathQueryLanguage } = await import("./vanilla-jsoneditor.js");
              queryLanguagesBuffer[languageId] = jmespathQueryLanguage;
              break;
            }
          }
        }
      }
      const allQueryLanguages = Object.values(queryLanguagesBuffer);
      if (allQueryLanguages.length === 0) {
        return;
      }
      return allQueryLanguages;
    };
    const removeFullWidthButton = () => {
      if (!fullWidthButton.value)
        return;
      fullWidthButton.value.removeEventListener("click", onButtonClick);
      fullWidthButton.value = null;
    };
    const setFullWidthButton = async () => {
      if (typeof window === "undefined")
        return;
      const { getElement, createElement } = await import("./full-width-button-handler.js");
      const { setFullWidthButtonStyle } = await import("./styles-handler.js");
      await setFullWidthButtonStyle();
      const oldButton = getElement(".jse-full-width");
      const pluginOptionFlag = (pluginOptions == null ? void 0 : pluginOptions.fullWidthButton) !== void 0 ? pluginOptions == null ? void 0 : pluginOptions.fullWidthButton : true;
      const fullWidthButtonFlag = props.fullWidthButton !== void 0 ? props.fullWidthButton : pluginOptionFlag;
      if (!fullWidthButtonFlag || oldButton)
        return;
      if (fullWidthButton.value) {
        removeFullWidthButton();
      }
      const menu = getElement(".jse-menu");
      fullWidthButton.value = createElement("button");
      fullWidthButton.value.classList.add("jse-full-width");
      fullWidthButton.value.classList.add("jse-button");
      fullWidthButton.value.classList.add("svelte-497ud4");
      fullWidthButton.value.innerHTML += fullWidthIcon;
      menu.appendChild(fullWidthButton.value);
      fullWidthButton.value.addEventListener("click", onButtonClick);
    };
    const onButtonClick = () => {
      var _a, _b;
      max.value = !max.value;
      if (max.value) {
        (_a = fullWidthButton.value) == null ? void 0 : _a.classList.add("jse-full-width--active");
      } else {
        (_b = fullWidthButton.value) == null ? void 0 : _b.classList.remove("jse-full-width--active");
      }
    };
    const expandCollapseAll = async (value) => {
      var _a;
      if (mode.value !== "tree")
        return;
      await ((_a = editor.value) == null ? void 0 : _a.expand(() => value));
    };
    const onChange = (content, previousContent, status) => {
      if (blockChange.value) {
        blockChange.value = false;
        return;
      }
      blockUpdate.value = true;
      if (!!content.json) {
        emit("update:json", content.json);
        emit("update:modelValue", content.json);
      }
      if (!!content.text) {
        emit("update:text", content.text);
        emit("update:jsonString", content.text);
        emit("update:modelValue", content.text);
      }
      emit("change", content, previousContent, status);
    };
    const onError = (err) => {
      emit("error", err);
    };
    const onChangeMode = (newMode) => {
      mode.value = newMode;
      emit("change-mode", newMode);
      emit("update:mode", newMode);
    };
    const onChangeQueryLanguage = (queryLanguageId2) => {
      emit("change-query-language", queryLanguageId2);
    };
    const onFocus = () => {
      emit("focus");
    };
    const onBlur = () => {
      emit("blur");
    };
    const onRenderMenu = (items, context) => {
      nextTick(() => {
        setFullWidthButton();
      });
      if (typeof props.onRenderMenu === "function") {
        return props.onRenderMenu(items, context);
      }
      return items;
    };
    const makeEditorProps = async () => {
      const options = { fullWidthButton: true, ...pluginOptions || {} };
      const queryLanguages = await makeQueryLanguages();
      return {
        ...pickDefinedProps(options, props),
        queryLanguages,
        queryLanguageId: queryLanguageId.value,
        onChange,
        onError,
        onChangeMode,
        onChangeQueryLanguage,
        onFocus,
        onBlur,
        onRenderMenu
      };
    };
    const fallbackSlot = ref(true);
    const getContent = () => {
      const getJsonContent = (json = {}) => {
        if (json === null || typeof json === "undefined" || typeof json === "number" || typeof json === "bigint" || typeof json === "string" || typeof json === "boolean") {
          return {
            json
          };
        }
        if (Array.isArray(json)) {
          return {
            json: [...json]
          };
        }
        return {
          json: { ...json }
        };
      };
      const getTextContent = (text = "") => {
        return {
          text: text || ""
        };
      };
      const propValue = props.modelValue || props.value;
      if (propValue) {
        if (mode.value === "text") {
          return getTextContent(propValue);
        } else {
          return getJsonContent(propValue);
        }
      }
      if (props.json) {
        return getJsonContent(props.json);
      }
      if (props.text) {
        return getTextContent(props.text);
      }
      if (props.jsonString) {
        return getTextContent(props.jsonString);
      }
      return getTextContent();
    };
    const initView = async () => {
      if (typeof window === "undefined")
        return;
      if (!editor.value) {
        const editorProps = await makeEditorProps();
        const { JSONEditor } = await import("./vanilla-jsoneditor.js");
        fallbackSlot.value = false;
        editor.value = new JSONEditor({
          target: container.value,
          props: editorProps
        });
        await editor.value.set(getContent());
      }
      await editor.value.focus();
    };
    const updateProps = async () => {
      const props2 = await makeEditorProps();
      editor.value.updateProps(props2);
    };
    const updateContent = () => {
      if (blockUpdate.value) {
        blockUpdate.value = false;
        return;
      }
      blockChange.value = true;
      editor.value.update(getContent());
    };
    const destroyView = () => {
      if (editor.value) {
        editor.value.destroy();
        editor.value = null;
      }
      removeFullWidthButton();
    };
    watch([
      ...watchPropNames.map((propName) => {
        return () => props[propName];
      })
    ], updateProps, { deep: true });
    watch([() => props.modelValue, () => props.value, () => props.json, () => props.text, () => props.jsonString], updateContent, {
      deep: true
    });
    watch(() => props.mode, (newMode) => {
      if (newMode !== mode.value) {
        mode.value = newMode;
        updateProps();
      }
    });
    watch(() => darkThemeStyle.value, async (value) => {
      if (!!value) {
        const { setDarkThemeStyle } = await import("./styles-handler.js");
        await setDarkThemeStyle();
      }
    }, { immediate: true });
    onMounted(() => {
      nextTick(() => {
        initView();
      });
    });
    onBeforeUnmount(() => {
      destroyView();
    });
    expose({
      async $collapseAll() {
        await expandCollapseAll(false);
      },
      async $expandAll() {
        await expandCollapseAll(true);
      },
      async $expand(callback) {
        var _a;
        await ((_a = editor.value) == null ? void 0 : _a.expand(callback));
      },
      $get() {
        var _a;
        return (_a = editor.value) == null ? void 0 : _a.get();
      },
      async $set(content) {
        var _a;
        await ((_a = editor.value) == null ? void 0 : _a.set(content));
      },
      async $update(content) {
        var _a;
        await ((_a = editor.value) == null ? void 0 : _a.update(content));
      },
      async $updateProps(props2) {
        var _a;
        await ((_a = editor.value) == null ? void 0 : _a.updateProps(props2));
      },
      async $refresh() {
        var _a;
        await ((_a = editor.value) == null ? void 0 : _a.refresh());
      },
      async $focus() {
        var _a;
        await ((_a = editor.value) == null ? void 0 : _a.focus());
      },
      async $destroy() {
        var _a;
        await ((_a = editor.value) == null ? void 0 : _a.destroy());
      },
      async $patch(operations) {
        var _a;
        return await ((_a = editor.value) == null ? void 0 : _a.patch(operations));
      },
      $transform(args) {
        var _a;
        (_a = editor.value) == null ? void 0 : _a.transform(args);
      },
      async $scrollTo(path) {
        var _a;
        await ((_a = editor.value) == null ? void 0 : _a.scrollTo(path));
      },
      $findElement(path) {
        var _a;
        return (_a = editor.value) == null ? void 0 : _a.findElement(path);
      },
      async $acceptAutoRepair() {
        var _a;
        return await ((_a = editor.value) == null ? void 0 : _a.acceptAutoRepair());
      },
      $validate() {
        var _a;
        return (_a = editor.value) == null ? void 0 : _a.validate();
      }
    });
    return {
      max,
      getHeight,
      container,
      darkThemeStyle,
      fallbackSlot
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
  }, [
    _ctx.fallbackSlot ? renderSlot(_ctx.$slots, "default", { key: 0 }) : createCommentVNode("", true)
  ], 38);
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
