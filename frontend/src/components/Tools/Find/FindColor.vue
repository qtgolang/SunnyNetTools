<template>
  <el-select v-model="value" placeholder="Select" size="small">
    <el-option
        v-for="item in options"
        :key="item.value"
        :label="item.label"
        :value="item.value"
        :disabled="item.disabled"
    />
  </el-select>
</template>

<script>
import {Config_Find, Config_HomeTextMark, Config_IsDark} from "../../config/Config.js";

const DefaultKey = "1";
export default {
  data() {
    return {
      value: "",
      get options() {
        const _options = [];
        for (const [key, value] of Config_HomeTextMark) {
          _options.push({
            value: key,
            label: value["名称"],
          })
        }
        return _options;
      },
    }
  },
  methods: {
    init() {
      const _options = [];
      for (const [key, value] of Config_HomeTextMark) {
        _options.push({
          value: key,
          label: value["名称"],
        })
      }
      if (_options.length > 0) {
        this.value = _options[0].key;
      }
      return _options;
    },
    GetValue() {
      return this.GetColorValue(this.value)
    },
    GetColorValue(key) {
      const o = Config_HomeTextMark.get(key)
      if (o) {
        if (Config_IsDark.value) {
          return o["深色主题"];
        }
        return o["浅色主题"];
      }
      return "#000000";
    },
    GetColorObj(key) {
      const o = Config_HomeTextMark.get(key)
      if (o) {
        return {
          dark: o["深色主题"],
          light: o["浅色主题"]
        };
      }
      return null;
    },
    getColorValue(key, count) {
      requestAnimationFrame(() => {
        if (key === DefaultKey) {
          const o = Config_HomeTextMark.get(key)
          if (o) {
            Config_Find.value.Color = key
            return;
          }
          if (count > 1000) {
            return;
          }
          this.getColorValue(key, count + 1)
          return;
        }
        Config_Find.value.Color = key
      })
    },
    getKeyValue() {
      return this.value;
    }
  },
  watch: {
    value(newVal) {
      this.getColorValue(newVal, 0)
    }
  },
  mounted() {
    this.value = DefaultKey
  }
}

</script>
