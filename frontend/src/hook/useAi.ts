import { getSystemConfig,fetchGet, fetchPost } from "@/system/config";
import { useAssistantStore } from '@/stores/assistant';
export function askAi(question: string,action:string) {
    const assistantStore = useAssistantStore();
    const config = getSystemConfig();


}