<script setup lang="ts">
import { useAiChatStore } from "@/stores/aichat";
import { Search } from "@element-plus/icons-vue";
import { isMobileDevice } from "@/util/device";
// import { t } from "@/i18n";
// import { notifyInfo,notifySuccess } from "@/util/msg.ts";
const chatStore = useAiChatStore();
const handleSearch = async () => {
  if (chatStore.searchInput != "") {
    const list = await chatStore.getChatList()
    const res = list.filter((item: any) => {
      if (item.title.includes(chatStore.searchInput)) {
        return item;
      }
    })
    if (res.length > 0) {
      chatStore.chatList = res
    }
  } else {
    await chatStore.getChatList()
  }
};

</script>
<template>
  <el-dialog v-model="chatStore.showInfo" width="600" append-to-body :fullscreen="isMobileDevice() ? true : false">
    <ai-chat-info />
  </el-dialog>
  <el-scrollbar>
    <el-header class="search">
      <el-input placeholder="搜索" @keyup.enter="handleSearch" :prefix-icon="Search" class="search-input"
        v-model="chatStore.searchInput" />
      <button class="add-chat" @click="chatStore.showBox(false)">
        <el-icon>
          <Plus />
        </el-icon>
      </button>
    </el-header>
    <div v-for="(item, key) in chatStore.chatList" :key="key">
      <div :class="['list-item', { active: item.id === chatStore.activeId }]">
        <el-row justify="space-around">
          <el-col :span="20" @click.stop="chatStore.setActiveId(item.id)" class="chat-title">
            {{ item.title }}
          </el-col>
          <el-col :span="4" class="iconlist">
            <!-- <el-icon size="15">
              <Edit />
            </el-icon> -->
            <el-icon size="15" @click.stop="chatStore.deleteChat(item.id)">
              <Delete />
            </el-icon>
          </el-col>

        </el-row>
      </div>
    </div>
  </el-scrollbar>
</template>
<style scoped lang="scss">
.search {
  display: flex;
  align-items: center;
  justify-content: space-evenly;
  width: 100%;
  height: 50px;
  padding: 0;
  padding-right: 10px;
  -webkit-app-region: drag;
  border-bottom: 1px solid #e8e8e8;
  border-left: 1px solid #e8e8e8;
}

.add-chat {
  width: 40px;
  height: 30px;
  border: none;
  border-radius: 4px;
  background-color: #f0f0f0;
}

.search-input {
  width: calc(100% - 20px);
  margin: 10px;
  height: 30px;
  font-size: 0.7rem;
  -webkit-app-region: no-drag;
  --el-input-placeholder-color: #bfbfbf !important;
  --el-input-icon-color: #bfbfbf !important;
}

.list-item {
  width: 95%;
  height: 65px;
  transition: all 0.5s;
  margin: 0 auto;
  border-radius: 4px;
  margin-bottom: 5px;
  overflow: hidden;
  margin-top: 5px;
  background-color: #fff;
}

.list-item:hover,
.list-item.active {
  background-color: #e8f3ff;
}

.chat-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  padding-left: 10px;
}

.iconlist {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-right: 10px;
}

.iconlist .el-icon {
  cursor: pointer;
}
</style>