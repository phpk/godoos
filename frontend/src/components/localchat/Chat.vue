<template>
  <el-row justify="space-between">
    <el-col :span="2">
      <chat-nav />
    </el-col>
    <el-col :span="6">
      <chat-domain />
    </el-col>
    <el-col :span="16">
      <chat-content />
    </el-col>
  </el-row>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from "vue";

import { useLocalChatStore } from "@/stores/localchat";
import { ElMessage } from "element-plus";
import { getSystemConfig } from "@/system/config";
const store = useLocalChatStore();
const config = getSystemConfig();

let source:any;
onMounted(async () => {
  await store.init()
  init()
});
onUnmounted(() => {
  if (source) {
    source.close();
  }
});
function init() {
  if (typeof EventSource === "undefined") {
    ElMessage.error("您的浏览器不支持SSE");
    return;
  }
  const sseUrl = config.apiUrl + "/localchat/sse";
  source = new EventSource(sseUrl);
  // 当接收到消息时触发
  source.onmessage = async function (event:any) {
    //console.log("has message!");
    const eventData = event.data; // 先保存原始数据
    const jsonData = JSON.parse(eventData); // 解析数据
    //console.log(jsonData);
    if (jsonData.type == "user_list") {
      //store.userList = jsonData;
      store.setUserList(jsonData.content);
      //await nextTick();
    }
    if(jsonData.type == 'text'){
      store.addText(jsonData);
    }
    if(jsonData.type == 'file'){
      store.addFile(jsonData);
    }
  };
  // 当与服务器的连接打开时触发
  source.onopen = function () {
    console.log("Connection opened.");
  };

  // 当与服务器的连接关闭时触发
  source.onerror = function () {
    console.log("Connection closed.");
  };
}

</script>

<style lang="scss" scoped>


</style>