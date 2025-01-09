<script setup lang="ts">
import { useAiChatStore } from "@/stores/aichat";
import { useModelStore } from "@/stores/model";
import { t } from "@/i18n";
import { notifyError } from "@/util/msg.ts";
import { ElScrollbar } from "element-plus";
import { getSystemConfig } from "@/system/config";
import { Vue3Lottie } from "vue3-lottie";
import { file } from "jszip";
import { isMobileDevice } from "@/util/device";
const chatStore = useAiChatStore();
const modelStore = useModelStore();
const isPadding = ref(false); //是否发送中
const webSearch = ref(false);
const imageInput: any = ref(null);
let imageData = ref("");
let fileContent = ref("");
let fileName = ref("");
const messageContainerRef = ref<InstanceType<typeof ElScrollbar>>();
const messageInnerRef = ref<HTMLDivElement>();
// User Input Message
const userMessage = ref("");
const promptMessage = computed(() => {
  return [
    {
      content: chatStore.chatInfo.prompt || "",
      chatType: "chat",
      chatId: chatStore.activeId,
      role: "system",
      id: Date.now(),
      createAt: Date.now(),
    },
  ];
});
const requestMessages = computed(() => {
  const contextLen = modelStore.chatConfig.chat.contextLength;
  //console.log(contextLen)
  if (chatStore.messageList.length <= contextLen) {
    return [...promptMessage.value, ...chatStore.messageList];
  } else {
    // 截取最新的10条信息
    const slicedMessages = chatStore.messageList.slice(-contextLen);
    return [...promptMessage.value, ...slicedMessages];
  }
});
const sendMessage = async () => {
  if (chatStore.activeId < 1) {
    notifyError(t("index.notFindChatModel"));
    return;
  }
  if (userMessage.value) {
    // Add the message to the list
    if (isPadding.value === true) return;
    const saveMessage = {
      content: userMessage.value,
      chatId: chatStore.activeId,
      role: "user",
      id: Date.now(),
      createdAt: Date.now(),
    };
    chatStore.messageList.push(saveMessage);
    await chatStore.addMessages(chatStore.activeId, saveMessage);

    // Clear the input
    userMessage.value = "";

    // Create a completion
    isPadding.value = true;
    await createCompletion();
  }
};

const createCompletion = async () => {
  try {
    const config = getSystemConfig()
    const messageId = Date.now();
    const saveMessage: any = {
      content: "",
      role: "assistant",
      chatType: "chat",
      chatId: chatStore.activeId,
      id: messageId,
      createdAt: messageId,
    };

    const chatConfig = modelStore.chatConfig.chat;
    let postMsg: any = {
      messages: requestMessages.value,
      model: chatStore.chatInfo.model,
      engine: chatStore.chatInfo.engine,
      stream: false,
      webSearch: webSearch.value,
      fileContent: fileContent.value,
      fileName: fileName.value,
      options: chatConfig,
    };
    if (imageData.value != "") {
      const img2txtModel = await modelStore.getModel("img2txt");
      const usermsg = chatStore.messageList[chatStore.messageList.length - 1];
      postMsg = {
        model: img2txtModel.model,
        //"prompt":userMessage.value,
        engine: img2txtModel.info.engine,
        stream: false,
        webSearch: false,
        options: chatConfig,
        messages: [
          {
            role: "user",
            content: usermsg.content,
            images: [imageData.value],
          },
        ],
      };
    }
    const postData: any = {
      method: "POST",
      body: JSON.stringify(postMsg),
    };
    //console.log(postData)
    const completion = await fetch(config.apiUrl + '/ai/chat', postData);
    //const completion:any = await modelStore.getModel(postData)
    imageData.value = "";
    fileContent.value = "";
    fileName.value = "";
    if (!completion.ok) {
      const errorData = await completion.json();
      notifyError(errorData.error.message);
      isPadding.value = false;
      return;
    }
    const res = await completion.json();
    //console.log(res)
    if (res && res.choices && res.choices.length > 0) {
      if (res.choices[0].message.content) {
        const msg = res.choices[0].message.content;
        saveMessage.content = msg;
        chatStore.messageList.push(saveMessage);
        await chatStore.addMessages(chatStore.activeId, saveMessage);
      }
    }
    isPadding.value = false;
  } catch (error: any) {
    isPadding.value = false;
    notifyError(error.message);
  }
};
const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainerRef && messageContainerRef.value) {
      messageContainerRef.value!.setScrollTop(
        messageInnerRef.value!.clientHeight
      );
    }
  });
};
watch(
  () => chatStore.messageList,
  () => {
    scrollToBottom();
  },
  {
    deep: true,
  }
);
const handleKeydown = (e: any) => {
  if (e.key === "Enter" && (e.altKey || e.shiftKey)) {
    // 当同时按下 alt或者shift 和 enter 时，插入一个换行符
    e.preventDefault();
    userMessage.value += "\n";
  } else if (e.key === "Enter") {
    // 当只按下 enter 时，发送消息
    e.preventDefault();
    sendMessage();
  }
};
const selectImage = async () => {
  imageInput.value.click();
};
const uploadImage = async (event: any) => {
  const file = event.target.files[0];
  if (!file) {
    return;
  }
  //console.log(file)
  if (file.type.startsWith('image/')) {
    const img2txtModel = await modelStore.getModel("img2txt");
    if (!img2txtModel) {
      notifyError(t("aichat.notEyeModel"));
      return;
    }
  }
  const reader = new FileReader();
  reader.onload = (e: any) => {
    const fileData = e.target.result.split(",")[1];
    if (file.type.startsWith('image/')) {
      imageData.value = fileData;
    } else {
      fileContent.value = fileData;
      fileName.value = file.name;
    }
    //console.log(fileContent.value)
  };

  reader.readAsDataURL(file);
};

</script>
<template>
  <div class="chat-bot">
    <div class="top-menu">
      <el-icon size="15" @click.stop="chatStore.showBox(true)" class="top-menu-button">
        <Tools />
      </el-icon>
      <el-icon size="15" @click.stop="chatStore.clearChatHistory" class="top-menu-button">
        <DeleteFilled />
      </el-icon>
    </div>
    <div class="messsage-area">
      <el-scrollbar v-if="chatStore.messageList.length > 0" class="message-container" ref="messageContainerRef">
        <div ref="messageInnerRef">
          <ai-chat-message v-for="message in chatStore.messageList" :key="message.messageId" :content="message.content"
            :link="message.link" :role="message.role" :createdAt="message.createdAt" />
        </div>
      </el-scrollbar>
      <div class="no-message-container" v-else>
        <Vue3Lottie animationLink="/bot/chat.json" :height="420" :width="420" />
      </div>
    </div>
    <div class="input-area">
      <div class="input-panel">
        <el-row :gutter="24" style="border-bottom: none;">
          <el-col :span="2">
            <el-button class="file-btn" @click="selectImage" size="large" icon="Paperclip" circle
              :class="{ 'selected-image': imageData != '' || fileContent != '' }" />
            <input type="file" ref="imageInput"
              accept="image/*,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,application/vnd.ms-powerpoint,application/vnd.openxmlformats-officedocument.presentationml.presentation,application/vnd.ms-excel,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,application/pdf"
              style="display: none" @change="uploadImage" />
          </el-col>
          <el-col :span="2">
            <el-button class="websearch-btn" @click="webSearch = !webSearch" size="large" icon="ChromeFilled" circle
              :type="webSearch ? 'primary' : 'default'" />
          </el-col>
          <el-col :span="17">
            <el-input v-model="userMessage" :placeholder="t('aichat.askme')" size="large" clearable
              @keydown="handleKeydown" :autofocus="isMobileDevice() ? false : true" class="ai-input-area" />
          </el-col>
          <el-col :span="2">
            <el-button v-if="!isPadding" @click="sendMessage" icon="Promotion" type="info" size="large" circle />
            <el-button type="primary" size="large" v-else loading-icon="Eleme" icon="Loading" loading circle />
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@import "@/assets/chat.scss";
</style>