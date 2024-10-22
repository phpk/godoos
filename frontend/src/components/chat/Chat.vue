<script setup lang="ts">
import { useChatStore } from '@/stores/chat';
import { Search } from '@element-plus/icons-vue';
import {getWorkflowUrl} from '@/system/config'
const store = useChatStore()
const workUrl = getWorkflowUrl()
onMounted(() => {
  store.initChat()
})
</script>
<template>
  <el-container class="container">
    <!--菜单-->
    <el-aside class="menu">
      <chat-menu />
    </el-aside>
    <el-container class="side" v-if="store.currentNavId < 3">
      <!--搜索栏-->
      <el-header class="search" v-if="store.currentNavId < 2">
        <el-input placeholder="搜索" :prefix-icon="Search" class="search-input" v-model="store.search" />
      </el-header>
      <!--好友列表-->
      <el-main class="list">
        <el-scrollbar>
          <chat-msg-list v-if="store.currentNavId == 0" />
          <chat-user-list v-if="store.currentNavId == 1" />
          <!-- <chat-work-list v-if="store.currentNavId == 2" /> -->
        </el-scrollbar>
      </el-main>
    </el-container>
    <el-container class="chat-box">
      <chat-box v-if="store.currentNavId < 2" />
    </el-container>
    <el-container class="chat-setting" v-if="store.currentNavId == 2">
      <iframe class="workflow" :src="workUrl"></iframe>
    </el-container>
    <el-container class="chat-setting" v-if="store.currentNavId == 5">
      <ChatUserSetting />
    </el-container>
  </el-container>

</template>
<style scoped>
.container {
  display: flex;
  height: 100%;
  width: 100%;
  overflow-y: hidden;
  overflow-x: hidden;
}

.menu {
  width: 55px;
  background-color: #2E2E2E;
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

.search {
  width: 100%;
  /* 占据整个宽度 */

  height: 50px;
  padding: 0;
  -webkit-app-region: drag;
}

.search-input {
  width: calc(100% - 20px);
  /* 减去左右边距 */
  margin: 10px;
  -webkit-app-region: no-drag;
  --el-input-placeholder-color: #818181 !important;
  --el-input-icon-color: #5D5D5D !important;
}

.list {
  width: 100%;
  /* 占据整个宽度 */
  padding: 0;
  overflow-y: hidden;
  overflow-x: hidden;
}

.chat-box {
  flex: 3;
  /* 占据剩余宽度的三倍 */
  max-height: max-content;
  background-color: #F5F5F5;
}

.chat-setting {
  width: calc(100% - 65px);
  /* 占据整个宽度 */
  height: 100%;
  overflow: hidden;
}
.workflow {
  width: 100%;
  height: 100%;
  object-fit: contain;
  border: none;
}
</style>