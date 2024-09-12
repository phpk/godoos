<script setup>
import {ref} from "vue";
import { useChatStore } from "@/stores/chat";
import {Setting as SettingIcon} from "@element-plus/icons-vue";
const store = useChatStore()
</script>

<template>
  <el-row>
    <el-avatar shape="square" :size="40" class="userAvatar" :src="store.userInfo.avatar"/>
  </el-row>
  <el-row @click="store.setCurrentNavId(0)">
    <el-icon v-if="store.currentNavId === 0" class="menu-icon-on">
      <ChatLineRound />
    </el-icon>
    <el-icon v-else class="menu-icon">
      <ChatRound />
    </el-icon>
  </el-row>
  <el-row @click="store.setCurrentNavId(1)">
    <el-icon v-if="store.currentNavId === 1" class="menu-icon-on">
      <UserFilled />
    </el-icon>
    <el-icon v-else class="menu-icon">
      <User />
    </el-icon>
  </el-row>
  <el-row @click="store.setCurrentNavId(2)">
    <el-icon v-if="store.currentNavId === 2" class="menu-icon-on">
      <Platform />
    </el-icon>
    <el-icon v-else class="menu-icon">
      <Monitor />
    </el-icon>
  </el-row>
  <el-row @click="store.setCurrentNavId(5)">
    <el-icon v-if="store.currentNavId === 5" class="menu-icon-on">
      <SettingIcon />
    </el-icon>
    <el-icon v-else class="menu-icon">
      <SettingIcon />
    </el-icon>
  </el-row>
</template>

<style scoped>
.userAvatar {
  margin: 30px auto 10px;
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
