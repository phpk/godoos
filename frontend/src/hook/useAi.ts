import { getSystemConfig, fetchPost } from "@/system/config";
import { useAssistantStore } from '@/stores/assistant';
import { useModelStore } from "@/stores/model";
export async function askAi(question: any, action: string) {
    try {
        const assistantStore = useAssistantStore();
        const modelStore = useModelStore();
        const config = getSystemConfig();
        const model = await modelStore.getModel('chat')
        if (!model) {
            return '请先设置模型'
        }
        let prompt = ""
        if (action === 'creation_ask') {
            if (question.title) {
                prompt = question.title
            } else {
                return ""
            }
            if (question.content || question.content != "") {
                prompt = `${prompt} \n ${question.content}`
            }
        } else {
            prompt = await assistantStore.getPrompt(action)
            if (!prompt) {
                return '请先设置prompt'
            }
            if (question.content) {
                prompt = prompt.replace('{content}', question.content)
            }
            if (question.title) {
                prompt = prompt.replace('{title}', question.title)
            }
            if (question.category) {
                prompt = prompt.replace('{category}', question.category)
            }
        }

        const apiUrl = config.aiUrl + '/ai/chat'
        const postMsg: any = {
            messages: [
                {
                    //role: "assistant",
                    role: "user",
                    content: prompt
                },
            ],
            engine: model.info.engine,
            model: model.model,
            stream: false,
            options: modelStore.chatConfig.creation,
        };
        const complain = await fetchPost(apiUrl, JSON.stringify(postMsg))
        if (!complain.ok) {
            return '请求失败'
        }
        const data = await complain.json()
        return data.choices[0].message.content
    } catch (error) {
        return '请求失败' + error
    }


}
export async function addKnowledge(path: string) {
    const modelStore = useModelStore();
    const config = getSystemConfig();
    const model = await modelStore.getModel('embeddings')
    if (!model) {
        return {
            code:-1,
            message:'请先设置嵌入模型'
        }
    }
    const apiUrl = config.aiUrl + '/ai/addknowledge'
    const postMsg: any = {
        engine: model.info.engine,
        model: model.model,
        file_path: path
    };
    const complain = await fetchPost(apiUrl, JSON.stringify(postMsg))
    if (!complain.ok) {
        return {
            code:-1,
            message:'请求失败'
        }
    }
    return await complain.json()
}