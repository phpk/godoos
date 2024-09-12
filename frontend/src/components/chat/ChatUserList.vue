<script setup>
import { ref } from "vue";
import { useChatStore } from "@/stores/chat";
const store = useChatStore()

</script>

<template>
  <div v-for="item in store.chatList">
    <div class="list-item" @click="store.changeChatList(item)"
      :style="{ backgroundColor: item.id === store.targetUserId ? '#C4C4C4' : '' }">
      <el-row>
        <el-col :span="6">
          <el-avatar shape="square" :size="40" class="avatar" :src="item.avatar" />
        </el-col>
        <el-col :span="18">
          <el-row>
            <el-col :span="18">
              <div class="previewName">{{ item.name }}</div>
            </el-col>
            <el-col :span="6">
              <div class="previewTime">{{ item.previewTimeFormat }}</div>
            </el-col>
          </el-row>
          <el-row>
            <div v-if="item.previewType === 0" class="previewChat">{{ item.previewMessage }}</div>
            <div v-if="item.previewType === 1" class="previewChat">{{ item.userId === id ? "你" : "\"" + item.name +
              "\""}}撤回了一条消息</div>
          </el-row>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<style scoped>
.list-item {
  width: 240px;
  height: 60px;
}

.list-item:hover {
  background-color: #D0D0D0;
}

.avatar {
  margin: 10px
}

.previewName {
  margin-top: 10px;
  font-size: 15px;
  font-family: Arial, sans-serif;
  line-height: 1.5;
}

.previewChat {
  font-size: 14px;
  font-family: Arial, sans-serif;
  color: #999999;
}

.previewTime {
  float: right;
  margin-top: 10px;
  margin-right: 10px;
  font-size: 12px;
  font-family: Arial, sans-serif;
  color: #999999;
}
</style>
