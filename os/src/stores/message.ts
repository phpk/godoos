import { getUrl } from "@/utils/request";
import { defineStore } from "pinia";
import { useLoginStore } from "@/stores/login";
import { useSettingsStore } from "@/stores/settings";
import { useChatStore } from "@/stores/chat";
import { useFileSystemStore } from "@/stores/filesystem";
import { ref } from "vue";
export const useMessageStore = defineStore('messageStore', () => {
    const loginStore = useLoginStore();
    const settingsStore = useSettingsStore();
    const chatStore = useChatStore();
    const fsStore = useFileSystemStore();
    const source: any = ref(null)
    function initMessage() {
        if (!loginStore.isLoginState) return;
        const url = getUrl('/user/message', false);
        source.value = new EventSource(url);

        source.value.onmessage = function (event: any) {
            settingsStore.checkIsLock();
            const data = JSON.parse(event.data);
            //console.log(data)
            handleMessage(data);
        };
        source.value.onerror = function (event: any) {
            console.error('EventSource error:', event);
        };
    }
    function closeMessage() {
        if (source.value) {
            source.value.close();
            source.value = null;
        }
    }
    async function handleMessage(message: any) {
        //console.log(message)
        switch (message.type) {
            case 'update':
                break;
            case 'online':
                chatStore.onlineUserData(message.data)
                break;
            case 'user':
                //console.log(message.data)
                chatStore.userChatMessage(message.data)
                break
            case 'group':
                //console.log(message.data)
                chatStore.groupChatMessage(message.data);
                break;
            case 'update_group':
                chatStore.groupInviteMessage(message.data)
                break;
            case 'system':
                chatStore.systemMessage(message.data)
                break;
            case 'share_refresh':
                console.log(message.data)
                fsStore.currentShareFile = message.data;
                break;
            default:
                console.warn('Unknown message type:', message.type);
        }
    }
    return {
        initMessage,
        closeMessage,
        handleMessage,
        source
    }
})