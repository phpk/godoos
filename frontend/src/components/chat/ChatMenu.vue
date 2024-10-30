<script setup>
import { useChatStore } from "@/stores/chat";
import { Setting as SettingIcon } from "@element-plus/icons-vue";

const store = useChatStore();
</script>

<template>
  <el-row>
    <el-avatar
      shape="square"
      :size="40"
      class="userAvatar"
      :src="store.userInfo.avatar"
    />
  </el-row>
  <el-row @click="store.setCurrentNavId(0)">
    <div class="menu-icon-box">
      <el-icon
        v-if="store.currentNavId === 0"
        class="menu-icon-on"
      >
        <ChatLineRound />
      </el-icon>
      <el-icon
        v-else
        class="menu-icon"
      >
        <ChatRound />
      </el-icon>
    </div>
  </el-row>
  <el-row @click="store.setCurrentNavId(1)">
    <div class="menu-icon-box">
      <el-icon
        v-if="store.currentNavId === 1"
        class="menu-icon-on"
      >
        <UserFilled />
      </el-icon>
      <el-icon
        v-else
        class="menu-icon"
      >
        <User />
      </el-icon>
    </div>
  </el-row>
  <el-row @click="store.setCurrentNavId(2)">
    <div class="menu-icon-box">
      <el-icon
        v-if="store.currentNavId === 2"
        class="menu-icon-on"
      >
        <Platform />
      </el-icon>
      <el-icon
        v-else
        class="menu-icon"
      >
        <Monitor />
      </el-icon>
    </div>
  </el-row>
  <el-row @click="store.setCurrentNavId(5)">
    <div class="menu-icon-box">
      <el-icon
        v-if="store.currentNavId === 5"
        class="menu-icon-on"
      >
        <SettingIcon />
      </el-icon>
      <el-icon
        v-else
        class="menu-icon"
      >
        <SettingIcon />
      </el-icon>
    </div>
  </el-row>
</template>

<style scoped>
.userAvatar {
  margin: 30px auto 10px;
  -webkit-app-region: no-drag;
}

.menu-icon-box {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  margin: 10px auto;
  border-radius: 15%;
  /* background-color: #bae7ff; */
  transition: all 0.5s;
}
.menu-icon-box:hover {
  background-color: #d9d9d9;
}

.menu-icon {
  margin: 0px auto;
  font-size: 25px;
  cursor: pointer;
}

.menu-icon-on {
  margin: 0px auto;
  font-size: 25px;
  color: #1890ff;
  cursor: pointer;
}
</style>