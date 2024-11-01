<template>
  <!-- 聊天编辑区域 -->
  <div 
  class="edit-box"
  @dragover.prevent
  @drop.prevent="handleDrop"
  >
    <el-input 
    type="textarea" 
    :rows="4" 
    class="message-input"
    @keydown.enter="keyDown($event)"
    v-model="store.sendInfo" />
      <el-tooltip placement="top" content="按enter键发送，按ctrl+enter键换行">
        <el-icon :size="22" class="win11-chat-send-button" @click="store.sendMsg('text')">
          <Promotion />
        </el-icon>
      </el-tooltip>
  </div>
</template>


<script setup lang="ts">
import { useLocalChatStore } from "@/stores/localchat"
//import { notifyError } from "@/util/msg";

const store = useLocalChatStore()
// 按下回车键
function keyDown(event: any) {
  if(store.sendInfo == '')return
  if (event.ctrlKey && event.keyCode === 13) {
    store.sendInfo = store.sendInfo + "\n"
    
  } else if (event.keyCode === 13) {
    event.preventDefault() // 阻止浏览器默认换行操作
    store.sendMsg('text')
    return false
  }
}
const handleDrop = (event:any) => {
  console.log("handleDrop")
  const frompathArrStr = event?.dataTransfer?.getData('frompath');
  event.preventDefault()
  const files = JSON.parse(frompathArrStr) as string[];
  if (files && files.length > 0) {
    // 处理拖放的文件，例如上传
    //console.log('Files dropped:', files);
    store.sendInfo = files;
    //store.uploadFile(files)
    store.sendMsg('applyfile')

  }
};

</script>


<style scoped>
.edit-box {
  position: relative; /* 为容器添加相对定位，以便子元素可以相对于它进行定位 */
}

.message-input {
  resize: none; /* 禁止调整输入框大小，保持布局稳定 */
}


.win11-chat-send-button {
  position: absolute;
  bottom: 5px;
  right: 5px;
  width: 30px; /* 缩小宽度 */
  height: 30px; /* 减小高度 */
  border-radius: 50%; /* 较小的圆角 */
  background-color: #E8F0FE; /* 浅蓝色，符合Win11的轻量风格 */
  color: #0078D4; /* 使用Win11的强调色作为文字颜色 */
  font-weight: bold;
  border: 1px solid #B3D4FC; /* 添加边框，保持简洁风格 */
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1); /* 轻微阴影 */
  transition: all 0.2s ease; /* 快速过渡效果 */
}

.win11-chat-send-button:hover {
  background-color: #D1E4FF; /* 悬浮时颜色略深，保持浅色调 */
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* 稍微增强阴影 */
}

.win11-chat-send-button:active {
  background-color: #B3D4FC; /* 按下时颜色更深，但依然保持清新 */
  box-shadow: 0 1px 2px rgba(0, 0, 0.1); /* 回复初始阴影 */
  transform: translateY(1px); /* 微小下移，模拟按下 */
}
</style>
