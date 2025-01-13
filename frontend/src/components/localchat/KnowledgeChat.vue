<script setup lang="ts">
import { useAiChatStore } from "@/stores/aichat";
import { useModelStore } from "@/stores/model";
import { t } from "@/i18n";
import { notifyError } from "@/util/msg.ts";
import { ElScrollbar } from "element-plus";
import { getSystemConfig } from "@/system/config";
import { Vue3Lottie } from "vue3-lottie";
import { BrowserWindow, useSystem } from "@/system";
import { isMobileDevice } from "@/util/device";
const chatStore = useAiChatStore();
const modelStore = useModelStore();
const isPadding = ref(false); //是否发送中
const webSearch = ref(false);
const imageInput: any = ref(null);
const win: any = inject<BrowserWindow>("browserWindow");
let imageData = ref("");
let fileContent = ref("");
let fileName = ref("");
const messageContainerRef = ref<InstanceType<typeof ElScrollbar>>();
const messageInnerRef = ref<HTMLDivElement>();
// User Input Message
const userMessage = ref("");
const knowledgeId = ref(0)
const system = useSystem()
const messageList: any = ref([])
const chatInfo: any = ref({})
const activeId: any = ref(0)
const promptMessage = computed(() => {
    return [
        {
            content: chatInfo.value.prompt || "",
            chatType: "chat",
            chatId: activeId.value,
            role: "system",
            id: Date.now(),
            createAt: Date.now(),
        },
    ];
});
onMounted(async () => {
    knowledgeId.value = win?.config?.knowledgeId || 0
    await chatStore.initChat(knowledgeId.value)
    const res = await chatStore.getActiveChat()
    messageList.value = res.messageList.value
    chatInfo.value = res.chatInfo.value
    activeId.value = res?.chatInfo?.value?.id || 0
    console.log(messageList.value, chatInfo.value, activeId.value)
    //await aiStore.initChat()
});
const requestMessages = computed(() => {
    const contextLen = modelStore.chatConfig.chat.contextLength;
    //console.log(contextLen)
    if (messageList.value.length <= contextLen) {
        return [...promptMessage.value, ...messageList.value];
    } else {
        // 截取最新的10条信息
        const slicedMessages = messageList.value.slice(-contextLen);
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
//       notifyError(errorData.message);
//       isPadding.value = false;
//       return;
//     }
//     const res = await completion.json();
//     console.log(res)
//     // let prompt = await chatStore.getPrompt("knowledge")
//     // if (prompt == '') {
//     //   notifyError("知识库prompt为空")
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
        notifyError(t("aichat.selectModel"));
        return;
    }
    if (userMessage.value) {
        // Add the message to the list
        if (isPadding.value === true) return;
        let saveMessage: any = {
            content: userMessage.value,
            chatId: activeId.value,
            role: "user",
            id: Date.now(),
            createdAt: Date.now(),
        };
        messageList.value.push(saveMessage);

        await chatStore.addMessages(activeId.value, saveMessage);

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
            chatId: activeId.value,
            id: messageId,
            createdAt: messageId,
        };

        const chatConfig = modelStore.chatConfig.chat;
        const knowledgeId = chatInfo.value.knowledgeId * 1;
        let postMsg: any = {
            messages: requestMessages.value,
            model: chatInfo.value.model,
            engine: chatInfo.value.engine,
            stream: false,
            webSearch: webSearch.value,
            fileContent: fileContent.value,
            fileName: fileName.value,
            options: chatConfig,
            knowledgeId: knowledgeId,
        };
        if (imageData.value != "") {
            const img2txtModel = await modelStore.getModel("img2txt");
            const usermsg = messageList.value[messageList.value.length - 1];
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
                messageList.value.push(saveMessage);
                await chatStore.addMessages(activeId.value, saveMessage);
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
    () => messageList,
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
    <el-dialog v-model="chatStore.showInfo" width="600" append-to-body :fullscreen="isMobileDevice() ? true : false">
        <ai-chat-info />
    </el-dialog>
    <div class="chat-bot">
        <div class="messsage-area">
            <el-scrollbar v-if="messageList.length > 0" class="message-container" ref="messageContainerRef">
                <div ref="messageInnerRef">
                    <ai-chat-message v-for="message in messageList" :key="message.messageId" :content="message.content"
                        :link="message.link" :role="message.role" :createdAt="message.createdAt"
                        :doc="message.doc || []" :web_search="message.web_search || []" :system="system" />
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
                        <el-button class="websearch-btn" @click="webSearch = !webSearch" size="large"
                            icon="ChromeFilled" circle :type="webSearch ? 'primary' : 'default'" />
                    </el-col>
                    <el-col :span="17">
                        <el-input v-model="userMessage" :placeholder="t('aichat.askme')" size="large" clearable
                            @keydown="handleKeydown" :autofocus="isMobileDevice() ? false : true"
                            class="ai-input-area" />
                    </el-col>
                    <el-col :span="2">
                        <el-button v-if="!isPadding" @click="sendMessage" icon="Promotion" type="info" size="large"
                            circle />
                        <el-button type="primary" size="large" v-else loading-icon="Eleme" icon="Loading" loading
                            circle />
                    </el-col>
                </el-row>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">
@import "@/assets/chat.scss";
</style>