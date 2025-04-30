<template>
    <div v-if="
        (store.searchList.length > 0 ? store.searchList : store.chatList)
            .length > 0
    " v-for="item in store.searchList.length > 0
        ? store.searchList
        : store.chatList" :key="item.id">
        <div class="list-item" @click="store.getSessionInfo(item.chatId, item.type)" :style="{
            backgroundColor:
                item.chatId == store.targetChatId ? '#E8F3FF' : '',
        }">
            <el-row>
                <el-col :span="6">
                    <el-avatar v-if="item.type == 'user'" shape="square" :size="40" class="avatar"
                        :src="item.avatar"></el-avatar>

                    <el-avatar v-else-if="item.type == 'group'" shape="square" :size="40" class="avatar" style="
							background-color: #165dff;
							display: flex;
							align-items: center;
							justify-content: center;
						">
                        <div style="
								width: 25px;
								height: 25px;
								background-image: url('@/assets/icons/group.png');
								background-size: cover;
							"></div>
                    </el-avatar>
                    <el-avatar v-else-if="item.type == 'system'" shape="square" :size="40" class="avatar"
                        icon="el-icon-message">
                    </el-avatar>
                </el-col>

                <!-- 在线状态 -->
                <el-icon v-if="item.online && item.type == 'user'" style="
						position: absolute;
						left: 40px;
						bottom: 5px;
						color: #0078d4;
					">
                    <CircleCheckFilled />
                </el-icon>
                <!-- 离线状态 -->
                <el-icon v-else-if="item.type == 'user'" style="
						position: absolute;
						left: 40px;
						bottom: 5px;
						color: #999999;
					">
                    <CircleCloseFilled />
                </el-icon>

                <el-col :span="18" class="preview">
                    <el-row class="preview-content">
                        <el-col :span="18" class="preview-left">
                            <div class="preview-top-content">
                                <div class="previewName">
                                    {{ item.displayName }}
                                </div>

                                <div class="previewTime">
                                    {{ formatTime(item.time) }}
                                </div>
                            </div>

                            <div class="previewChat">
                                <span>{{ item.previewMessage }}</span>
                            </div>
                        </el-col>
                        <el-col :span="6" class="preview-right">
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
import { ref, computed } from "vue";

const store = useChatStore();
const id = ref("1");

// 添加格式化时间的方法
function formatTime(timestamp) {
    const date = new Date(timestamp);
    const month = date.getMonth() + 1; // 月份是从0开始的
    const day = date.getDate();
    return `${month}-${day}`;
}

// 假设消息数据列表在 store 中
const filteredMessages = computed(() => {
    return store.search
        ? store.messages.filter((message) =>
            message.content
                .toLowerCase()
                .includes(store.search.toLowerCase())
        )
        : store.messages;
});
</script>

<style scoped>
.list-item {
    width: 95%;
    height: 60px;
    display: flex;
    transition: all 0.5s;
    margin: 0 auto;
    border-radius: 4px;
    margin-bottom: 5px;
    overflow: hidden;
    margin-top: 5px;
    background-color: #fff;
}

.list-item:hover {
    background-color: #e8f3ff;
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
    margin-left: 5px;
}

.preview-top-content {
    height: 20px;
    display: flex;
    min-width: 110px;
    max-width: 200px;
    align-items: center;
    margin-left: 10px;
    justify-content: space-between;
}

.previewName {
    font-weight: 400;
    font-size: 14px;
    color: #000000;
    font-family: Arial, sans-serif;
    line-height: 1.5;
    overflow: hidden;
    /* 隐藏超出部分 */
    text-overflow: ellipsis;
    /* 显示为省略号 */
    white-space: nowrap;
    /* 不换行 */
}

.previewTime {
    font-size: 12px;
    color: #86909c;
}

.preview-left {
    display: flex;
    flex-direction: column;
    justify-content: center;
    height: 100%;
}

.previewChat {
    line-height: 20px;
    width: 100px;
    height: 20px;
    margin-left: 10px;
    font-size: 12px;
    font-family: Arial, sans-serif;
    color: #86909c;
    overflow: hidden;
    /* 隐藏超出部分 */
    text-overflow: ellipsis;
    /* 显示为省略号 */
    white-space: nowrap;
    /* 不换行 */
}

.preview-right {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
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
