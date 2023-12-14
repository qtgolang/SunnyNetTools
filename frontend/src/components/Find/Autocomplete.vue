<template>
  <el-autocomplete
      ref="auto"
      v-model="state1"
      :fetch-suggestions="querySearch"
      clearable
      class="inline-input w-50"
      placeholder="请输入要搜索的内容"
      @select="handleSelect"
  />
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue'

interface RestaurantItem {
  value: string
  link: string
}

const auto = ref('')
const state1 = ref('')
const data = {
  "UTF8": [],
  "GBK": [],
  "Hex": [],
  "pb": [],
  "整数4": [],
  "整数8": [],
  "浮点数4": [],
  "浮点数8": [],
  "Base64": [],
}

const restaurants = ref<RestaurantItem[]>([])
const querySearch = (queryString: string, cb: any) => {
  const results = queryString
      ? restaurants.value.filter(createFilter(queryString))
      : restaurants.value
  // call callback function to return suggestions
  cb(results)
}
const createFilter = (queryString: string) => {
  return (restaurant: RestaurantItem) => {
    return (
        restaurant.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0
    )
  }
}


const handleSelect = (item: RestaurantItem) => {

}
onMounted(() => {
  restaurants.value = data['UTF8']
})

function SetAutoType(a) {
  restaurants.value = data[a]
}

function Set(Type, value) {
  for (let i = 0; i < data[Type].length; i++) {
    if (data[Type][i].value === value) {
      return
    }
  }
  data[Type].push({value: value})
  restaurants.value = data[Type]
}

function GetValue() {
  return state1.value
}

function focus() {
  auto.value.focus()
}

defineExpose({
  SetAutoType, Set, GetValue, focus
});
</script>
