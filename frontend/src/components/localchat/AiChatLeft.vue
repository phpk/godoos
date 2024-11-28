<script setup lang="ts">
import { useAiChatStore } from "@/stores/aichat";
import { Search } from "@element-plus/icons-vue";
// import { t } from "@/i18n";
// import { notifyInfo,notifySuccess } from "@/util/msg.ts";
const chatStore = useAiChatStore();
const showBox = (flag: any) => {
  chatStore.isEditor = flag;
  if (flag === true) {
    chatStore.editInfo = toRaw(chatStore.chatInfo);
  } else {
    chatStore.editInfo = {
      title: "",
      model: "",
      prompt: "",
      promptName: "",
    };
  }
  chatStore.showInfo = true;
};

</script>
<template>
  <el-dialog v-model="chatStore.showInfo" width="600" append-to-body>
    <ai-chat-info />
  </el-dialog>
  <el-scrollbar>
    <el-header class="search">
      <el-input placeholder="搜索" :prefix-icon="Search" class="search-input" v-model="chatStore.searchInput" />
      <button class="add-chat" @click="showBox(false)">
        <el-icon>
          <Plus />
        </el-icon>
      </button>
    </el-header>
    <div v-for="(item, key) in chatStore.chatList" :key="key">
      <div class="list-item">
        <el-row justify="space-around">
          <el-col :span="17" @click.stop="chatStore.setActiveId(item.id)" class="chat-title">
            {{ item.title }}
          </el-col>
          <el-col :span="7" class="iconlist">
            <el-icon size="15">
              <Edit />
            </el-icon>
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
  height: 60px;
  transition: all 0.5s;
  margin: 0 auto;
  border-radius: 4px;
  margin-bottom: 5px;
  overflow: hidden;
  margin-top: 5px;
  background-color: #fff;
}

.list-item:hover,.list-item .active {
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
.iconlist .el-icon{
  cursor: pointer;
}
</style>