<template>
  <el-container class="container">
    <el-aside class="menu">
      <chat-nav />
    </el-aside>
    <el-container class="side">
      <chat-domain v-if="store.navId > 0" />
      <ai-chat-left v-else />
    </el-container>
    <el-container class="chat-box">
      <chat-content v-if="store.navId > 0" />
      <ai-chat-main v-else />
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { useAiChatStore } from "@/stores/aichat";
import { useLocalChatStore } from "@/stores/localchat";

const store = useLocalChatStore();
const aiStore = useAiChatStore();

//let source:any;
onMounted(async () => {
  await store.init()
  await aiStore.initChat()
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

@media screen and (max-width: 768px) {
  .container {
    height: calc(100vh - vh(90));
  }

  .menu {
    position: fixed;
    bottom: 0;
    width: 100vw;
    display: flex;
    justify-content: space-evenly;
  }

  .side {
    height: 100%;
    background-color: #fff;
  }

  .chat-box {
    height: 100%;
  }
}
</style>