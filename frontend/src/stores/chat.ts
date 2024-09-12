import { defineStore } from 'pinia'
import emojiList from "@/assets/emoji.json"
import { getSystemConfig } from '@/system/config'
export const useChatStore = defineStore('chatStore', () => {
    const userList: any = ref([]) // 用户列表
    const chatList: any = ref([]) // 消息列表
    const msgList: any = ref([]) //聊天消息列表
    const userInfo: any = ref({})
    const showChooseFile = ref(false)
    const currentPage = ref(1)
    const pageSize = ref(50)
    const scrollbarRef = ref()
    const innerRef = ref()
    const config = getSystemConfig()
    const currentNavId = ref(0)
    const message: any = ref()
    const targetUserInfo:any = ref({})
    const targetUserId = ref(0)
    const search = ref('')
    const contextMenu = ref({
        visible: false,
        chatMessageId: 0,
        list: [
            {
                id: 2,
                label: '撤回',
            }
        ],
        x: 0,
        y: 0
    })
    const initChat = () => {
        if(config.userInfo.avatar == ''){
            config.userInfo.avatar = '/logo.png'
        }
        userInfo.value = config.userInfo
    }
    const setCurrentNavId = (id: number) => {
        currentNavId.value = id
    }
    const sendMessage = async () => {
        await setScrollToBottom()
    }
    const setScrollToBottom = async () => {
        await nextTick()
        const max = innerRef.value.clientHeight
        scrollbarRef.value.setScrollTop(max)
    }
    const changeChatList = async (item:any) => {
        // const res = await getChatList(id)
        // chatList.value = res.data
    }
    const handleContextMenu = async (id: number) => {
        contextMenu.value.visible = false;
    }
    const showContextMenu = (event: any, id: number) => {
        contextMenu.value.visible = true;
        contextMenu.value.chatMessageId = id;
        contextMenu.value.x = event.x;
        contextMenu.value.y = event.y;
    }
    return {
        emojiList,
        userList,
        chatList,
        msgList,
        userInfo,
        search,
        showChooseFile,
        currentPage,
        pageSize,
        scrollbarRef,
        innerRef,
        currentNavId,
        targetUserInfo,
        targetUserId,
        message,
        contextMenu,
        initChat,
        setCurrentNavId,
        sendMessage,
        changeChatList,
        handleContextMenu,
        showContextMenu
    }
})