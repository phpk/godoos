import { getUrl } from "@/system/config";
import { defineStore } from "pinia";
import { useUpgradeStore } from "./upgrade";
export const useMessageStore = defineStore('messageStore', () => {
    const upgradeStore = useUpgradeStore();
    function systemMessage(){
        const url = getUrl('/system/message',false); 
        const source = new EventSource(url);

        source.onmessage = function(event) {
            const data = JSON.parse(event.data);
            //console.log(data)
            handleMessage(data);
        };
        source.onerror = function(event) {
            console.error('EventSource error:', event);
        };
    }
    async function handleMessage(message:any) {
        switch (message.type) {
            case 'update':
                upgradeStore.checkUpdate(message.data)
                break;
            case 'chat':
                
                break;
            default:
                console.warn('Unknown message type:', message.type);
        }
    }
    return {
        systemMessage
    }
})