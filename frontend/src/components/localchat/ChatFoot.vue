<template>
  <footer class="footer-container">
    <!--工具栏-->
    <el-row type="flex" class="toolbar">
      <el-popover
        placement="top"
        popper-class="chat-icon-popover"
        trigger="click"
      >
        <template #reference>
          <div class="emoji-button">
            <img
              width="24"
              height="24"
              class="emoji-image"
              src="/image/chat/emoji.svg"
            />
          </div>
        </template>
        <el-scrollbar class="emoji-scroll">
          <ul class="emoji-list-container">
            <li
              v-for="item in store.emojiList"
              :key="item.title"
              class="p-5px list-none hover:animate-heart-beat animate-count-animated animate-duration-1s cursor-pointer"
              :title="item.title"
            >
              <img
                width="30"
                height="30"
                :src="item.icon"
                @click="selectIcon(item.icon)"
              />
            </li>
          </ul>
        </el-scrollbar>
      </el-popover>
      <div class="upload-button" @click="selectImg()">
        <el-icon :size="22">
          <Picture />
        </el-icon>
      </div>
      <div class="upload-button" @click="selectFile()">
        <el-icon :size="22">
          <Link />
        </el-icon>
      </div>
      <div class="upload-button" @click="store.clearMsg()">
        <el-icon :size="22">
          <Delete />
        </el-icon>
      </div>
    </el-row>
    <ChatEditor></ChatEditor>
  </footer>
</template>

<script setup lang="ts">
import { useLocalChatStore } from "@/stores/localchat";
import { useChooseStore } from "@/stores/choose";

import { toRaw, watch } from "vue";
const store = useLocalChatStore();
const choose = useChooseStore();
//const editor = ref(null)
const imgExt = ["png", "jpg", "jpeg", "gif", "bmp", "webp", "svg"];
const choosetype = ref('image')
// 选择表情
function selectIcon(icon: string) {
  console.log("icon:", icon);
  store.sendInfo +=
    "{-" + icon.replace("/image/chat/emoji/", "").replace(".gif", "") + "-}";
}
function selectImg() {
  choosetype.value = 'image'
  choose.select("选择图片", imgExt);
}
function selectFile() {
  choosetype.value = 'applyfile'
  choose.select("选择文件", "*");
}
watch(
  () => choose.path,
  (newVal, _) => {
    //console.log("userList 变化了:", newVal);
    const paths = toRaw(newVal)
    if(paths.length > 0){
      store.sendInfo = paths;
      choose.path = []
      store.sendMsg(choosetype.value)
    }
  },
  { deep: true } // 添加deep: true以深度监听数组或对象内部的变化
);


</script>

<style lang="scss" scoped>
.footer-container {
  padding-left: 20px;
  padding-right: 20px;
}

.toolbar {
  margin-bottom: 10px;
}
.opacity-0 {
  opacity: 0;
  display: none;
  width: 0;
  height: 0;
}
.emoji-button {
  font-size: 20px;
  cursor: pointer;
}

.emoji-image {
  display: block;
  width: 24px;
  height: 24px;
}

.emoji-scroll {
  height: 150px;
}

.emoji-list {
  padding: 0;
  list-style-type: none;
  flex-wrap: wrap;
}

.emoji-list-item {
  padding: 5px;
  cursor: pointer;
  animation-duration: 1s;
  &.hover:animate-heart-beat {
    animation-name: heart-beat;
  }
}

.emoji-img {
  width: 30px;
  height: 30px;
}
.emoji-list-container {
  margin: 0; /* 对应 m0 */
  padding: 0; /* 对应 p0 */
  display: flex; /* 对应 flex */
  flex-wrap: wrap; /* 对应 flex-wrap */
}
.upload-button {
  margin-left: 15px;
  cursor: pointer;
}

.file-input {
  opacity: 0;
}

.answer-editor {
  /* Assuming this class already exists or you define it elsewhere */
}
</style>