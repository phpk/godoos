<template>
  <div class="container">
    <div class="nav">
      <ul>
        <li v-for="(item, index) in items" :key="index" @click="activeIndex = index"
          :class="{ active: index === activeIndex }">
          {{ item }}
        </li>
      </ul>
    </div>
    <div class="setting">
      <div v-if="0 === activeIndex">
        <LocalProxy />
      </div>
      <div v-if="1 === activeIndex">
        <NetProxy />
      </div>
    </div>
</div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import LocalProxy from "./LocalProxy.vue";
const items = ["本地代理", "远程访问"];
const activeIndex = ref(0);
</script>

<style scoped>
@import "./setStyle.css";
</style>