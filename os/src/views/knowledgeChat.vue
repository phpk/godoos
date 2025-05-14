<script setup lang="ts">
import { useAiChatStore } from "@/stores/aichat";
import { useModelStore } from "@/stores/model";
import { t } from "@/i18n";
import { errMsg } from "@/utils/msg.ts";
import { ElScrollbar } from "element-plus";
import { Vue3Lottie } from "vue3-lottie";
import { isMobileDevice } from "@/utils/device";
import { computed, nextTick, onMounted, ref, watch } from "vue";
import { askknowledge } from "@/api/knowledge";

const chatStore = useAiChatStore();
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
    knowledgeId.value = props.win?.props?.knowledgeId || 0
    activeId.value = knowledgeId.value
});
const requestMessages = computed(() => {
    //const contextLen = modelStore.chatConfig.chat.contextLength;
    const contextLen = 10;
    //console.log(contextLen)
    if (messageList.value.length <= contextLen) {
        return [...promptMessage.value, ...messageList.value];
    } else {
        // 截取最新的10条信息
        const slicedMessages = messageList.value.slice(-contextLen);
        return [...promptMessage.value, ...slicedMessages];
    }
});

const sendMessage = async () => {
    if (userMessage.value.trim()) {
        // Add the message to the list
        if (isPadding.value === true) return;
        let saveMessage: any = {
            content: userMessage.value,
            chatId: knowledgeId.value,
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
        const messageId = Date.now();
        const saveMessage: any = {
            content: "",
            role: "assistant",
            chatType: "chat",
            chatId: chatStore.activeId,
            id: messageId,
            createdAt: messageId,
        };

        //const knowledgeId = knowledgeId.value * 1;
        let postMsg: any = {
            messages: requestMessages.value,
            model: chatStore.chatInfo.model,
            engine: 'ollama',
            stream: false,
            webSearch: webSearch.value,
            fileContent: fileContent.value,
            fileName: fileName.value,
            //options: chatConfig,
            knowledgeId: knowledgeId.value * 1,
        };
     

        //console.log(postData)
        
        //const res = await modelStore.chatWithModel(postMsg);
        const res:any = await askknowledge(postMsg)
        //console.log(res)

        //const completion:any = await modelStore.getModel(postData)
        imageData.value = "";
        fileContent.value = "";
        fileName.value = "";
        if (!res) {
            errMsg("模型调用失败，请稍后重试！");
            isPadding.value = false;
            return;
        }
        //console.log(res)
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
                //chatStore.messageList.push(saveMessage);
                //await chatStore.addMessages(chatStore.activeId, saveMessage);
            }
        }
        if (res && res.message) {
            if (res.message.content.startsWith("<think>\n\n</think>\n\n")) {
                res.message.content = res.message.content.replace("<think>\n\n</think>\n\n", "")
            }
            saveMessage.content = res.message.content;
            messageList.value.push(saveMessage);
            //await chatStore.addMessages(chatStore.activeId, saveMessage);
        }
        isPadding.value = false;
    } catch (error: any) {
        isPadding.value = false;
        errMsg(error.message);
    }
};
const scrollToBottom = () => {
    nextTick(() => {
        if (messageContainerRef && messageContainerRef.value && messageInnerRef.value) {
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
            <el-scrollbar class="message-container" ref="messageContainerRef">
                <div ref="messageInnerRef" v-if="messageList.length > 0">
                    <ai-chat-message v-for="message in messageList" :key="message.messageId" :content="message.content"
                        :link="message.link" :role="message.role" :createdAt="message.createdAt"
                        :doc="message.documents || []" :web_search="message.web_search || []" />
                </div>
                <div class="no-message-container" v-else>
                    <Vue3Lottie animationLink="/os/bot/chat.json" :height="420" :width="420" />
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
    </div>
</template>

<style scoped lang="scss">
@use '@/styles/aichat.scss';
</style>