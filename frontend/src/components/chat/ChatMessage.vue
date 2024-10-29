<script setup lang="ts">
import { useChatStore } from '@/stores/chat';
const store = useChatStore();
</script>

<template>
  <div  v-for="item in store.chatHistory" :key="item.id">
    <div v-if="!item.isme" class="chat-item">
      <el-row>
        <el-col :span="8" />
        <el-col :span="14">
          <el-row>
            <el-col :span="24">
              <div class="chat-name-me">{{ item.userInfo.username }}</div>
            </el-col>
          </el-row>
          <div class="bubble-me" @contextmenu.prevent="store.showContextMenu($event, item.id)">
            <div class="chat-font">
              {{ item.content }}
            </div>
          </div>
        </el-col>
        <el-col :span="2">
          <div class="chat-avatar">
            <el-avatar shape="square" style="margin: 0;float: left" :size="32" class="userAvatar"
              :src="item.userInfo.avatar" />
          </div>
        </el-col>
      </el-row>
    </div>
    <div v-else class="chat-item">
      <el-row>
        <el-col :span="2">
          <div class="chat-avatar">
            <el-avatar shape="square" style="margin: 0;float: right" :size="32" class="userAvatar"
              :src="item.userInfo.avatar" />
          </div>
        </el-col>
        <el-col :span="14">
          <el-row>
            <el-col :span="24">
              <div class="chat-name-other">{{ item.userInfo.username }}</div>
            </el-col>
          </el-row>
          <div class="bubble-other">
            <div class="chat-font">
              {{ item.content }}
            </div>
          </div>
        </el-col>
        <el-col :span="8" />
      </el-row>
    </div>
    <div v-if="item.type === 1" class="withdraw">
      {{ item.userInfo.id === store.targetUserId ? "你" : item.userInfo.username }}撤回了一条消息
    </div>
  </div>
  <!--悬浮菜单-->
  <div class="context-menu" v-if="store.contextMenu.visible"
    :style="{ top: `${store.contextMenu.y}px`, left: `${store.contextMenu.x}px` }">
    <div v-for="contextItem in store.contextMenu.list" :key="contextItem.id" class="context-menu-item">
      <div class="context-menu-item-font" @click="store.handleContextMenu(contextItem)">
        {{ contextItem.label }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.bubble-me {
  background-color: #95EC69;
  float: right;
  border-radius: 4px;
  margin-right: 5px;
  margin-top: 5px;
}

.bubble-me:hover {
  background-color: #89D961;
}

.chat-name-me {
  font-size: 14px;
  font-family: Arial, sans-serif;
  line-height: 1.5;
  color: #B2B2B2;
  float: right;
  margin-right: 5px;
}

.bubble-other {
  background-color: #FFFFFF;
  float: left;
  border-radius: 4px;
  margin-left: 5px;
  margin-top: 5px;
}

.bubble-other:hover {
  background-color: #EBEBEB;
}

.chat-name-other {
  font-size: 14px;
  font-family: Arial, sans-serif;
  line-height: 1.5;
  color: #B2B2B2;
  float: left;
  margin-left: 5px;
}

.chat-font {
  margin: 8px;
  font-size: 15px;
  font-family: Arial, sans-serif;
  line-height: 1.5;
}

.chat-avatar {
  margin: 5px;
}

.chat-item {
  margin: 5px;
}

.withdraw {
  text-align: center;
  font-size: 13px;
  font-family: Arial, sans-serif;
  color: #999999;
  line-height: 3.2;
}

.context-menu {
  position: fixed;
  background-color: white;
  z-index: 9999;
  border: 1px solid #cccc;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.context-menu-item {
  width: 80px;
  height: 30px;
}

.context-menu-item:hover {
  background-color: #E2E2E2;
}

.context-menu-item-font {
  font-size: 14px;
  text-align: center;
  font-family: Arial, sans-serif;
  line-height: 2.2;
}
</style>