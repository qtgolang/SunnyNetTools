<template>
  <div>
    <select v-model="filterType" @change="onFilterChange">
      <option value="contains">包含</option>
      <option value="notContains">不包含</option>
      <option value="regexMatch">正则匹配</option>
    </select>
    <input
        type="text"
        v-model="filterText"
        placeholder="输入搜索内容..."
        @input="onFilterChange"
    />
  </div>
</template>

<script>
export default {
  data() {
    return {
      filterText: "",
      filterType: "contains",
      params: null,
    };
  },
  methods: {
    init(params) {
      this.params = params;
    },
    doesFilterPass(params) {
      const value = this.params.valueGetter(params.node);
      if (!value) return false;

      const filterText = this.filterText.toLowerCase();

      switch (this.filterType) {
        case "contains":
          return value.toLowerCase().includes(filterText);
        case "notContains":
          return !value.toLowerCase().includes(filterText);
        case "regexMatch":
          try {
            const regex = new RegExp(filterText, "i");
            return regex.test(value);
          } catch (e) {
            return false;
          }
        default:
          return false;
      }
    },
    isFilterActive() {
      return this.filterText !== "";
    },
    getModel() {
      return {
        filterText: this.filterText,
        filterType: this.filterType,
      };
    },
    setModel(model) {
      if (model) {
        this.filterText = model.filterText;
        this.filterType = model.filterType;
      }
    },
    onFilterChange() {
      this.params.filterChangedCallback();
    },
  },
};
</script>
