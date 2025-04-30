<script setup lang="ts">
import { useAiChatStore } from "@/stores/aichat";
import { useModelStore } from "@/stores/model";
import { t } from "@/i18n";
import { errMsg } from "@/utils/msg.ts";
import { ElScrollbar } from "element-plus";
import { Vue3Lottie } from "vue3-lottie";
import { isMobileDevice } from "@/utils/device";
import { computed, nextTick, onMounted, ref, watch } from "vue";
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
const knowledgeId = ref(0)
const props = defineProps<{
  win: any;
}>();
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
onMounted(async () => {
  knowledgeId.value = props.win?.props?.knowledgeId || 0
  console.log(knowledgeId.value, 'knowledgeId')
  await chatStore.initChat(knowledgeId.value)
  //await aiStore.initChat()
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
// async function addKnowledge(){
//   const knowledgeId = win?.config?.knowledgeId || 0;
//   let msg = userMessage.value
//   if (knowledgeId > 0) {
//     const askData: any = {
//       id: knowledgeId,
//       input: userMessage.value,
//     }
//     const config = getSystemConfig()
//     const postData: any = {
//       method: "POST",
//       body: JSON.stringify(askData),
//     };
//     const completion = await fetch(config.apiUrl + '/ai/askknowledge', postData);
//     if (!completion.ok) {
//       const errorData = await completion.json();
//       //console.log(errorData)
//       errMsg(errorData.message);
//       isPadding.value = false;
//       return;
//     }
//     const res = await completion.json();
//     console.log(res)
//     // let prompt = await chatStore.getPrompt("knowledge")
//     // if (prompt == '') {
//     //   errMsg("知识库prompt为空")
//     //   return
//     // }
//     if (res && res.data.length > 0) {
//       let context: string = "";
//       res.data.forEach((item: any) => {
//         context += "- " + item.content + "\n";
//       })
//       //prompt = prompt.replace("{content}", context)
//       msg = `请对\n${context}\n的内容进行分析，给出对用户输入的回答:${userMessage.value} `
//     }
//   }
//   return msg
// }
const sendMessage = async () => {
  if (chatStore.activeId < 1) {
    errMsg(t("aichat.selectModel"));
    return;
  }
  if (userMessage.value.trim()) {
    // Add the message to the list
    if (isPadding.value === true) return;
    let saveMessage: any = {
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
    const knowledgeId = chatStore.chatInfo.knowledgeId * 1;
    let postMsg: any = {
      messages: requestMessages.value,
      model: chatStore.chatInfo.model,
      engine: chatStore.chatInfo.engine,
      stream: false,
      webSearch: webSearch.value,
      fileContent: fileContent.value,
      fileName: fileName.value,
      //options: chatConfig,
      knowledgeId: knowledgeId,
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
        knowledgeId: knowledgeId,
      };
    }

    //console.log(postData)
    const res = await modelStore.chatWithModel(postMsg);
    console.log(res)

    //const completion:any = await modelStore.getModel(postData)
    imageData.value = "";
    fileContent.value = "";
    fileName.value = "";
    if (!res) {
      errMsg("模型调用失败，请稍后重试！");
      isPadding.value = false;
      return;
    }
    console.log(res)
    if (res && res.choices && res.choices.length > 0) {
      if (res.choices[0].message.content) {
        const msg = res.choices[0].message.content;
        saveMessage.content = msg;
        if (res.documents && res.documents.length > 0) {
          saveMessage.doc = res.documents;
        }
        if (res.web_search && res.web_search.length > 0) {
          saveMessage.web_search = res.web_search;
        }
        chatStore.messageList.push(saveMessage);
        await chatStore.addMessages(chatStore.activeId, saveMessage);
      }
    }
    if (res && res.message) {
      if (res.message.content.startsWith("<think>\n\n</think>\n\n")) {
        res.message.content = res.message.content.replace("<think>\n\n</think>\n\n", "")
      }
      saveMessage.content = res.message.content;
      chatStore.messageList.push(saveMessage);
      await chatStore.addMessages(chatStore.activeId, saveMessage);
    }
    isPadding.value = false;
  } catch (error: any) {
    isPadding.value = false;
    errMsg(error.message);
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
      errMsg(t("aichat.notEyeModel"));
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
  <el-dialog v-model="chatStore.showInfo" width="600" append-to-body :fullscreen="isMobileDevice() ? true : false">
    <ai-chat-info />
  </el-dialog>
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
      <el-scrollbar class="message-container" ref="messageContainerRef">
        <div ref="messageInnerRef" v-if="chatStore.messageList.length > 0">
          <ai-chat-message v-for="message in chatStore.messageList" :key="message.messageId" :content="message.content"
            :link="message.link" :role="message.role" :createdAt="message.createdAt" :doc="message.doc || []"
            :web_search="message.web_search || []" />
        </div>
        <div class="no-message-container" v-else>
          <Vue3Lottie animationLink="/bot/chat.json" :height="420" :width="420" />
        </div>
      </el-scrollbar>
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
              <el-button v-if="!isPadding" @click="sendMessage" icon="Promotion" type="primary" size="large" circle />
              <el-button type="danger" size="large" v-else loading-icon="Eleme" icon="Loading" loading circle />
            </el-col>
          </el-row>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped lang="scss">
@import "@/styles/aichat.scss";
</style>