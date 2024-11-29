import { defineStore } from "pinia"
import { db } from "./db.js"
import { t } from "@/i18n/index.ts"
import { useAssistantStore } from "./assistant.ts";
import { useModelStore } from "@/stores/model";
import { ref } from "vue";
export const useAiChatStore = defineStore('aichat', () => {
  const modelStore = useModelStore()
  const promptStore = useAssistantStore()
  const activeId: any = ref(0)
  const currentMessage: any = ref({})
  // 聊天列表
  const chatList: any = ref([])
  const chatInfo: any = ref({})
  const messageList: any = ref([])
  const modelList:any = ref([])
  const promptList: any = ref([])
  const searchInput: any = ref('')
  const showInfo = ref(false)
  const editInfo: any = ref({}); //编辑聊天信息
  const isEditor = ref(true);

  const newChat = async () => {
    const currentModel = await modelStore.getModel('chat')
    if (!currentModel) {
      return false
    }
    const promptData = await promptStore.getPrompt('chat')
    return await addChat(t('aichat.newchat'), currentModel.model, promptData, "")
  }
  const initChat = async () => {
    if (activeId.value === 0) {
      return await newChat()
    }
    modelList.value = await modelStore.getModelAction('chat')
    const promptRes = await promptStore.getPrompts('chat')
    promptList.value = promptRes.list
    chatList.value = await db.getAll('aichatlist')
    if(activeId.value > 0){
      messageList.value = await db.getByField('aichatmsg', 'chatId', activeId.value)
      chatInfo.value = await db.getOne('aichatlist', activeId.value)
    }
  }
  const getActiveChat = async () => {
    chatInfo.value = await db.getOne('aichatlist', activeId.value)
    messageList.value = await db.getByField('aichatmsg', 'chatId', activeId.value)
    chatList.value = await db.getAll('aichatlist')
    return { chatInfo, messageList, chatList }
  }
  const getChatList = async () => {
    chatList.value = await db.getAll('aichatlist')
    return chatList
  }
  // 添加聊天
  async function addChat(title: string, model: any, promptData: any, knowledgeId: string) {
    const newChat = {
      title,
      prompt: promptData.prompt,
      promptId: promptData.id,
      model,
      createdAt: Date.now(),
      knowledgeId
    }
    //console.log(newChat)
    activeId.value = await db.addOne('aichatlist', newChat)
    return activeId.value

  }
  async function setActiveId(newId: number) {
    activeId.value = newId
    chatInfo.value = await db.getOne('aichatlist', newId)
    messageList.value = await db.getByField('aichatmsg', 'chatId', newId)
  }
  // 删除单个聊天
  async function deleteChat(chatId: number) {
    await db.delete('aichatlist', chatId)
    await db.deleteByField('aichatmsg', 'chatId', chatId)
    //如果删除的id是当前id
    let id;
    if (chatId == activeId.value) {
      //
      const list = await db.getAll('aichatlist')
      if (list.length > 0) {
        id = list[0]['id']

      } else {
        id = await newChat()
      }
      setActiveId(id)
    }
    chatList.value = await db.getAll('aichatlist')
  }

  // 更新聊天菜单标题
  async function updateTitle(chatId: number, title: string) {
    await db.update('aichatlist', chatId, { title })
  }

  // 清空所有Chat
  async function clearChat() {
    await db.clear('aichatlist')
    await db.clear('aichatmsg')
  }


  // 新增历史消息
  async function addMessages(chatId: number, message: any) {
    const currentChat = await db.getOne('aichatlist', chatId)
    //console.log(currentChat)
    if (currentChat) {
      return await db.addOne('aichatmsg', message)
    }
  }

  async function getChat(chatId: number) {
    const chats = await db.getOne('aichatlist', chatId)
    //console.log(chats)
    const messages = await db.getByField('aichatmsg', 'chatId', chatId)
    return { chats, messages }
  }

  // 获取指定id的聊天的历史记录

  async function getChatHistory(chatId: number) {
    return await db.getByField('aichatmsg', 'chatId', chatId)
  }

  // 删除指定id的聊天的历史记录
  async function clearChatHistory() {
    if(activeId.value > 0){
      await db.deleteByField('aichatmsg', 'chatId', activeId.value)
      messageList.value = []
    }
    
  }

  // 更新聊天配置
  async function updateChat(config: any, chatId: number) {
    //console.log(config)
    return await db.update('aichatlist', chatId, config)
  }
  const showBox = (flag: any) => {
    isEditor.value = flag;
    if (flag === true) {
      editInfo.value = toRaw(chatInfo.value);
    } else {
      editInfo.value = {
        title: "",
        model: "",
        prompt: "",
        promptId: "",
      };
    }
    showInfo.value = true;
  };
  return {
    activeId,
    chatList,
    messageList,
    chatInfo,
    currentMessage,
    searchInput,
    showInfo,
    editInfo,
    isEditor,
    modelList,
    promptList,
    initChat,
    setActiveId,
    getActiveChat,
    getChatList,
    addChat,
    updateTitle,
    deleteChat,
    clearChat,
    addMessages,
    getChat,
    getChatHistory,
    clearChatHistory,
    updateChat,
    showBox
  }

}, {
  persist: {
    enabled: true,
    strategies: [
      {
        storage: localStorage,
        paths: [
          "activeId"
        ]
      }, // name 字段用localstorage存储
    ],
  }
})
