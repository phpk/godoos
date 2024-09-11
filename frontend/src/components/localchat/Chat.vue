<template>
  <el-row justify="space-between">
    <el-col :span="1">
      <chat-nav />
    </el-col>
    <el-col :span="6">
      <chat-domain />
    </el-col>
    <el-col :span="17">
      <chat-content />
    </el-col>
  </el-row>
</template>

<script setup lang="ts">
import { onMounted } from "vue";

import { useLocalChatStore } from "@/stores/localchat";

const store = useLocalChatStore();


//let source:any;
onMounted(async () => {
  await store.init()
});

</script>

<style lang="scss" scoped>


</style>