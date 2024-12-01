import { defineStore } from "pinia"
import { db } from "./db.js"
import { t } from "@/i18n/index.ts"
import {useAssistantStore} from "./assistant.ts";
import { useModelStore } from "@/stores/model";
import { ref } from "vue";
export const useAiChatStore = defineStore('aichat', () => {
  const modelStore = useModelStore()
  const promptStore = useAssistantStore()
  const activeId: any = ref(0)
  const currentMessage:any = ref({})
  // 聊天列表
  const chatList: any = ref([])
  const chatInfo:any = ref({})
  const messageList : any = ref([])
  const modelList = ref([])
  const promptList:any = ref([])
  const newChat = async() => {
    const currentModel = await modelStore.getModel('chat')
    if (!currentModel) {
      return false
    }
    const promptData = await promptStore.getPrompt('chat')
    return await addChat(t('chat.newchat'), currentModel, promptData, "")
  }
  const initChat = async () => {
    if (activeId.value === 0) {
      return await newChat()
    }
    modelList.value = await modelStore.getModelAction('chat')
    promptList.value = await promptStore.getPrompts('chat')
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
  async function addChat(title: string, modelData: any, promptData: any, knowledgeId:string) {
    const newChat = {
      title,
      prompt: promptData.prompt,
      promptId: promptData.id,
      modelId: modelData.id,
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
    await db.deleteByField('aichatmsg','chatId', chatId)
    //如果删除的id是当前id
    let id;
    if (chatId == activeId.value) {
      //
      const list = await db.getAll('aichatlist')
      if(list.length > 0) {
        id = list[0]['id']
        
      }else{
        id = await newChat()
      }
      setActiveId(id)
    }
    chatList.value = await db.getAll('aichatlist')
  }

  // 更新聊天菜单标题
  async function updateTitle(chatId: number, title: string) {
    await db.update('aichatlist', chatId, {title})
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
  async function clearChatHistory(chatId: number) {
    await db.deleteByField('aichatmsg', 'chatId', chatId)
  }

  // 更新聊天配置
  async function updateChat(config: any, chatId: number) {
    //console.log(config)
    return await db.update('aichatlist',chatId, config)
  }
  return {
    activeId,
    chatList,
    messageList,
    chatInfo,
    currentMessage,
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
    modelList,
    promptList
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
