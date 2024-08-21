<template>
  <div class="win11-msg-container">
    <el-scrollbar>
      <div v-if="store.navId < 2" class="user-list-area">
        <el-row
          class="user-list"
          justify="space-around"
          v-for="(msg, key) in store.contentList"
          :key="key"
          v-if="store.contentList.length > 0"
        >
          <!-- <el-col :span="5" class="avatar-col">
            <el-icon :size="22" class="avatar">
              <Place />
            </el-icon>
          </el-col> -->
          <el-col :span="22" :class="msg.targetIp == store.chatTargetIp ? 'user-item active' : 'user-item'" @click="store.setChatId(msg.targetIp)">
            <span class="username">
              
              <el-icon :size="12">
              <Place />
            </el-icon>
              {{ msg.reciperInfo.hostname }}
            </span>
            <span class="msg">{{ msg.content }}</span>
            
            <span class="userip" v-if="msg.readNum > 0">
            <el-badge :value="msg.readNum" :offset="[-3,8]">
            {{ formatChatTime(msg.createdAt) }}</el-badge>
            </span>
            <span class="userip" v-else>
              {{ formatChatTime(msg.createdAt) }}
            </span>
        
          </el-col>
        </el-row>
        <el-empty v-else :image-size="100" description="消息列表为空" />
      </div>
      
      <div v-else class="user-list-area">
        <el-row justify="space-between">
        <el-icon :size="18" @click="store.refreshUserList">
          <RefreshRight />
        </el-icon>
        <el-icon :size="18" @click="store.showAddUser = true">
          <CirclePlusFilled />
        </el-icon>
        </el-row>
        <el-row
          class="user-list"
          justify="space-around"
          v-for="(user, key) in store.userList"
          :key="key"
          v-if="store.userList.length > 0"
        >
          <el-col :span="5" class="avatar-col">
            <el-icon :size="22" class="avatar">
              <Monitor />
            </el-icon>
            <el-icon v-if="user.isOnline" :size="16" class="status-icon online">
              <CircleCheck />
            </el-icon>
            <el-icon v-else :size="16" class="status-icon offline">
              <Warning />
            </el-icon>
          </el-col>
          <el-col :span="19" class="user-item" @click="store.setChatId(user.ip)">
            <span class="username">{{ user.hostname }}</span>
            <span class="userip">{{ user.ip }}</span>
          </el-col>
        </el-row>
        <el-empty v-else :image-size="100" description="好友列表为空" />
      </div>
    </el-scrollbar>
  </div>
  <el-dialog v-model="store.showAddUser" title="添加用户" width="500">
    <el-form>
      <el-form-item :label-width="0">
        <el-input v-model="userIp" autocomplete="off" placeholder="输入用户IP 例如：192.168.1.16"/>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button type="primary" @click="store.addUser(userIp)">
          添加
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { useLocalChatStore } from "@/stores/localchat";
import { formatChatTime } from "@/util/common";

const store = useLocalChatStore();
const userIp = ref('')
</script>
<style scoped lang="scss">
.win11-msg-container {
  height: 100vh; /* 假设顶部有100px的导航栏 */
  width: 100%;
  border-right: 1px solid rgba(238, 238, 238, 0.5); /* 更柔和的边框 */
  border-top: 1px solid rgba(238, 238, 238, 0.5);
  border-radius: 0 12px 0 0; /* Win11风格的圆角 */
  background-color: #f8f8f8; /* 使用更亮的淡灰色背景 */
  overflow: hidden; /* 清除滚动条溢出 */
}

.el-scrollbar__wrap {
  padding: 16px; /* 内容区域的内边距 */
}

.message-item {
  display: flex;
  align-items: flex-start;
  &.self {
    justify-content: flex-end; /* 自身消息右对齐 */
  }
}

.bubble {
  padding: 8px 12px;
  border-radius: 12px; /* 圆角 */
  max-width: 75%;
  background-color: #ffffff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  &.bubble--right {
    background-color: #e8eaed; /* 更接近Win11的发送方消息背景色 */
  }
}
.user-list-area {
  width: 94%;
  margin: 3%;
}
.user-list {
  background-color: #f8f8f8; // 淡灰色背景，与容器区分
  border-radius: 8px; // 圆角边缘，更柔和的外观
  margin-bottom: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); // 温和的阴影效果，增加深度感
  .avatar-col {
    display: flex;
    justify-content: center;
    align-items: center; // 确保图标垂直居中
  }
  .active{
    background-color: #e8eaed;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); 
  }
  .user-item {
    display: flex;
    align-items: flex-start;
    flex-direction: column; /* 改变为列方向布局 */
    gap: 4px; /* 行间间距 */
    padding: 4px 0;
    cursor: pointer;
    transition: background-color 0.3s;

    &:hover {
      background-color: rgba(0, 0, 0, 0.05);
    }

    .avatar {
      border-radius: 50%; /* 头像圆角 */
    }

    .username {
      text-align: left; // 设置username左对齐
      font-size: 12px;
      color: #333;
      overflow: hidden;
      width:95%;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .msg {
      text-align: left; // 设置username左对齐
      font-size: 11px;
      color: #666;
      overflow: hidden;
      width:95%;
      text-overflow: ellipsis;
      white-space: nowrap;
    }


    .userip {
      text-align: right; // 设置userip右对齐
      margin-left: auto; // 这将使userip在容器内右对齐
      font-size: 11px;
      padding-right: 5px;
      color: #666;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
  .status-icon {
    margin-left: 2px; // 与头像保持一定距离
    margin-top:12px;
  }

  .online {
    color: green; // 在线图标颜色
  }

  .offline {
    color: red; // 离线图标颜色
  }
}

/* 空状态调整 */
.el-empty {
  margin: auto;
}
</style>