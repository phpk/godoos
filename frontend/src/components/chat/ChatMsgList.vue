<template>
  <div v-if="store.chatList.length > 0" v-for="item in store.chatList" :key="item.id">
    <div
      class="list-item"
      @click="store.changeChatList(item.id)"
      :style="{
        backgroundColor: item.id === store.targetUserId ? '#bae7ff' : '',
      }"
    >
      <el-row>
        <el-col :span="6">
          <el-avatar
            shape="square"
            :size="40"
            class="avatar"
            :src="item.avatar"
          />
        </el-col>
        <el-col :span="18" class="preview">
          <el-row class="preview-content">
            <el-col :span="18" class="preview-left">
              <div class="previewName">{{ item.nickname }}</div>
              <div class="previewChat">
                <span v-if="item.previewType === 0">{{ item.previewMessage }}</span>
                <span v-if="item.previewType === 1">
                  {{
                    item.targetUserId === id
                      ? "你"
                      : '"' + item.nickname + '"'
                  }}撤回了一条消息
                </span>
              </div>
            </el-col>
            <el-col :span="6" class="preview-right">
              <div class="previewTime">
                {{ item.previewTimeFormat }}
              </div>
            </el-col>
          </el-row>
        </el-col>
      </el-row>
    </div>
  </div>
  <div v-else class="emptyChat">
    <el-icon :size="60" class="chat-icon">
      <ChatSquare />
    </el-icon>
    <p class="empty-message">暂无数据</p>
  </div>
</template>

<script setup>
import { useChatStore } from "@/stores/chat";
import { ref } from "vue";

const store = useChatStore();
const id = ref("1");

</script>

<style scoped>
	.list-item {
		width: 94%;
		height: 60px;
		display: flex;
    margin: 0 auto;  
    border-radius: 4px; 
    transition: all 0.5s;
    margin-top: 5px;
	}

	.list-item:hover {
		background-color: #bae7ff;
	}

	.avatar {
		margin: 10px;
	}

	.preview {
		display: flex;
		align-items: center;
		justify-content: space-between;
		height: 100%;
	}

	.preview-content {
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: space-between;
		width: 100%;
	}

	.preview-left {
		display: flex;
		flex-direction: column;
		justify-content: center;
		height: 100%;
	}

	.previewName {
		margin-left: 10px;
		font-size: 12px;
		font-family: Arial, sans-serif;
		line-height: 1.5;
		overflow: hidden; /* 隐藏超出部分 */
		text-overflow: ellipsis; /* 显示为省略号 */
		white-space: nowrap; /* 不换行 */
		min-width: 100px; /* 最小宽度 */
		max-width: 100%; /* 最大宽度 */
	}

	.previewChat {
    height: 20px;
		margin-left: 10px;
		font-size: 10px;
		font-family: Arial, sans-serif;
		color: #999999;
		overflow: hidden; /* 隐藏超出部分 */
		text-overflow: ellipsis; /* 显示为省略号 */
		white-space: nowrap; /* 不换行 */
		min-width: 90px; /* 最小宽度 */
		max-width: 100%; /* 最大宽度 */
	}

	.preview-right {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 100%;
	}

	.previewTime {
		font-size: 12px;
		font-family: Arial, sans-serif;
		color: #999999;
		min-width: 30px; /* 最小宽度 */
		max-width: 100%; /* 最大宽度 */
	}

	.emptyChat {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		margin: 120px auto;
		text-align: center;
		font-size: 14px;
	}

	.chat-icon {
		color: #0078d4;
		margin-bottom: 20px;
	}

	.empty-message {
		font-size: 16px;
		color: #666666;
	}
</style>
