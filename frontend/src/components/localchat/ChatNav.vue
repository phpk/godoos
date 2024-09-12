<template>
  <el-row>
    <el-avatar shape="square" :size="40" class="userAvatar" src="/logo.png"/>
  </el-row>
  <el-row @click="store.handleSelect(1)">
    <el-icon v-if="store.navId === 1" class="menu-icon-on">
      <ChatLineRound />
    </el-icon>
    <el-icon v-else class="menu-icon">
      <ChatRound />
    </el-icon>
  </el-row>
  <el-row @click="store.handleSelect(2)">
    <el-icon v-if="store.navId === 2" class="menu-icon-on">
      <UserFilled />
    </el-icon>
    <el-icon v-else class="menu-icon">
      <User />
    </el-icon>
  </el-row>
  <!-- <el-space direction="vertical" :size="20" class="win11-chat-nav">
    <div :class="store.navId === item.index ? 'nav-item active' : 'nav-item'" v-for="item in store.navList" :key="item.index">
      <el-icon size="18" @click="store.handleSelect(item.index)">
        <component :is="item.icon" />
      </el-icon>
    </div>
  </el-space> -->
</template>

<script setup lang="ts">
import { useLocalChatStore } from '@/stores/localchat';
const store = useLocalChatStore();
</script>
<style scoped>
.userAvatar {
  margin: 10px auto;
  -webkit-app-region: no-drag
}

.menu-icon {
  margin: 20px auto;
  font-size: 25px;
  color: #8F8F8F;
  cursor: pointer;
}

.menu-icon-on {
  margin: 20px auto;
  font-size: 25px;
  color: #07C160;
  cursor: pointer;
}

</style>
