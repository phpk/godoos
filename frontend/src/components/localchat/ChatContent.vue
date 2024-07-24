<script setup lang="ts">
import { onMounted, inject, ref, nextTick, watch } from "vue";
import { useLocalChatStore } from "@/stores/localchat";
import { formatChatTime } from "@/util/common";
import Vditor from "vditor";
import "vditor/dist/index.css";
import { ElScrollbar } from "element-plus";
import { System } from "@/system";
const sys:any = inject<System>("system");
const store = useLocalChatStore();
const messageContainerRef = ref<InstanceType<typeof ElScrollbar>>();
const messageInnerRef = ref<HTMLDivElement>();
let isScrool = false;
onMounted(() => {
  scrollToBottom();
});
const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainerRef && messageContainerRef.value) {
      // messageContainerRef.value!.setScrollTop(
      //   messageInnerRef.value!.clientHeight
      // );
        messageContainerRef.value.setScrollTop(messageInnerRef.value!.clientHeight);
    }
  });
};
watch(
  () => store.msgList,
  (_) => {
    if(!isScrool){
      scrollToBottom();
    }
    
  },
  {
    deep: true,
  }
);
function replaceIconTags(text:any) {
  // 定义正则表达式，匹配 {*任意内容*} 的格式
  text = Vditor.md2html(text);
  const regex = /\{\-(.*?)\-\}/g;

  // 使用正则表达式的replace方法进行替换
  const replacedText = text.replace(regex, (_:any, p1:string) => {
    // p1 是匹配到的内容部分，这里直接构造图片标签
    return `<img src='/image/chat/emoji/${p1}.gif' style='width:30px;height:30px;' />`;
    //return `![avatar](/image/chat/emoji/${p1}.gif)`
  });

  return replacedText;
}
async function scroll({ scrollTop }: { scrollTop: number }) {
  if (store.msgList.length + 1 > store.pageSize && scrollTop < 1) {
    isScrool = true;
    await store.moreMsgList();
    isScrool = false;
  }
}
</script>

<template>
  <div class="chatContentContainer" v-if="store.chatTargetId > 0">
    <div class="message-area">
      <el-scrollbar
       max-height="100%" 
       class="scrollbar-container" 
       @scroll="scroll"
       ref="messageContainerRef">
        <div ref="messageInnerRef" class="message-wrap">
          <div v-for="(item, index) in store.msgList" :key="index" :class="['message-block', item.isMe ? 'mine' : 'theirs']">
            <div class="avatar-container">
              <div class="icon-container">
                <el-icon><component :is="item.isMe ? 'UserFilled' : 'Place'"/></el-icon>
              </div>
            </div>
            <div class="content">
              <div 
              v-if="item.type === 'text'"
              v-html="replaceIconTags(item.content)" 
              class="message-content">
              </div>
              <div v-if="item.type === 'file'">
                <div class="file-bubble">
                  <div class="file-content" v-for="el in item.content" @click="sys.openFile(el.path)">
                    <div class="file-icon"><FileIcon :file="el" /></div>
                    <div class="file-name">{{el.name}}</div>
                  </div>
                </div>
              </div>
              <div class="timestamp text-grey">{{ formatChatTime(item.createdAt) }}</div>
            </div>
          </div>
        </div>
      </el-scrollbar>
    </div>
    <ChatFoot class="mt-20px"></ChatFoot>
  </div>
  <div class="no-message-container" v-else>
    <el-icon :size="180" color="#0078d7">
      <ChatDotRound   />
    </el-icon>
  </div>
</template>
<style scoped lang="scss">
$win10-blue: #0078d7;
$win10-light-blue: #c7e8ff;
$win10-grey: #afafaf;
$win10-light-grey: #f2f2f2;

.chatContentContainer {
  flex: 1;
}

.message-area {
  
  height: 420px;
}

.scrollbar-container {
  width: 100%;

}

.message-wrap {
  padding: 10px;
}

.message-block {
  display: flex;
  margin-bottom: 10px;
  align-items: flex-end; // 保持底部对齐
  
  flex-direction: row;
}

.avatar-container {
  width: 40px;
  height: 40px;
  margin-right: 10px;
}

.icon-container {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  overflow: hidden;
  border: 1px solid $win10-grey;
  display: flex;
  justify-content: center;
  align-items: center;
}

.content {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  max-width: calc(100% - 150px);
}

.message-content {
  background: $win10-light-grey;
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 14px;
  color: $win10-grey;
  word-break: break-word;
}

.timestamp {
  font-size: 12px;
  color: $win10-grey;
}


.mine {
  // 添加这一行来确保'mine'类的消息整体靠右
  justify-content: flex-end;
  .content {
    // 由于'mine'消息块整体靠右，内容区域不需要特别处理对齐
    // justify-content: flex-end; 可以移除，因为它会与外部的justify-content冲突
    align-items: flex-end; // 保持底部对齐
  }
  
  .message-content {
    background-color: $win10-blue;
    color: $win10-light-blue;
    border-radius: 12px 2px 2px 2px; // 保持原有样式
  }
  
  .avatar-container {
    order: 1; 
    margin-left: 10px;
    margin-top:-20px;
  }
}

.theirs {
   .content {
    // 使对方的消息内容靠左对齐
    align-items: flex-end;
  }
  .message-content {
    border-radius: 2px 12px 2px 2px;
  }
}
.no-message-container {
    height: 100%;
    margin: 120px auto;
    text-align: center;
    justify-content: center;
}
.file-bubble {
  background-color: #f0f0f0; /* 背景色，可以根据需要调整 */
  border-radius: 10px; /* 圆角，让框看起来更柔和 */
  padding: 10px; /* 内边距，给内容一些空间 */
  margin-bottom: 10px; /* 气泡间的外边距，使它们看起来不紧凑 */
  max-width: 100%;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* 添加阴影效果，增强立体感 */
}

.file-content {
  display: flex;
  align-items: center;
  height:36px;
  line-height:36px;
  gap: 3px;
}
.file-content:hover {
  background-color: #e0e0e0; /* 改变背景色，悬停时更浅或更深，根据设计调整 */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2); /* 增强阴影效果，使气泡在悬停时更加突出 */
  transition: all 0.3s ease; /* 添加过渡效果，使变化平滑 */
}
.file-icon {
  width: 18px;
  height: 18px;
}

.file-name {
  flex: 1;
  overflow-wrap: break-word;
  word-break: break-all;
  white-space: normal;
  padding-left: 5px;
  /* 可选：限制文件名行数，如果需要 */
  /* max-height: 2em; */
  /* display: -webkit-box; */
  /* -webkit-line-clamp: 2; */
  /* -webkit-box-orient: vertical; */
  /* overflow: hidden; */
}
</style>