<template>
  <el-row>
    <el-avatar shape="square" :size="40" class="userAvatar" src="/logo.png"/>
  </el-row>
  <el-row @click="store.handleSelect(0)">
    <el-icon v-if="store.navId === 0" class="menu-icon-on">
      <Promotion />
    </el-icon>
    <el-icon v-else class="menu-icon">
      <Position />
    </el-icon>
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
