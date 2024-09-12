<template>
  <el-container class="container">
    <el-aside class="menu">
      <chat-nav />
    </el-aside>
    <el-container class="side">
      <chat-domain />
    </el-container>
    <el-container class="chat-box">
      <chat-content />
    </el-container>
  </el-container>
  <!-- <el-row justify="space-between">
    <el-col :span="1">
      <chat-nav />
    </el-col>
    <el-col :span="6">
      <chat-domain />
    </el-col>
    <el-col :span="17">
      <chat-content />
    </el-col>
  </el-row> -->
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
.container {
  display: flex;
  height: 100%;
  width: 100%;
  overflow-y: hidden;
  overflow-x: hidden;
}

.menu {
  width: 55px;
  background-color: #f2f2f2;
  overflow-y: hidden;
  overflow-x: hidden;
  -webkit-app-region: drag;
}
.side {
  flex: 1;
  /* 占据剩余宽度 */
  max-height: max-content;
  border-right: 1px solid #edebeb;
  overflow-y: hidden;
  overflow-x: hidden;
  background-color: #F7F7F7;
}
.chat-box {
  flex: 3;
  /* 占据剩余宽度的三倍 */
  max-height: max-content;
  background-color: #F5F5F5;
}


</style>