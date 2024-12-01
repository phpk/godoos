export const parseAiConfig = (config: any) => {
    if (!config.ollamaUrl) {
        config.ollamaUrl = `${window.location.protocol}//${window.location.hostname}:11434`
    }
    if(!config.ollamaDir) {
        config.ollamaDir = ''
    }
    if (!config.dataDir) {
        config.dataDir = ''
    }
    if (!config.aiUrl) {
        config.aiUrl = config.apiUrl
    }
    //openai
    if (!config.openaiUrl) {
        config.openaiUrl = 'https://api.openai.com/v1'
    }
    if (!config.openaiSecret) {
        config.openaiSecret = ""
    }
    //gitee
    if (!config.giteeSecret) {
        config.giteeSecret = ""
    }
    //cloudflare
    if(!config.cloudflareUserId){
        config.cloudflareUserId = ""
    }
    if(!config.cloudflareSecret){
        config.cloudflareSecret = ""
    }
    if(!config.deepseekSecret) {
        config.deepseekSecret = ""
    }
    if(!config.bigmodelSecret) {
        config.bigmodelSecret = ""
    }
    if(!config.volcesSecret){
        config.volcesSecret = ""
    }
    if(!config.alibabaSecret) {
        config.alibabaSecret = ""
    }
    if(!config.groqSecret) {
        config.groqSecret = ""
    }
    if(!config.mistralSecret) {
        config.mistralSecret = ""
    }
    if(!config.anthropicSecret) {
        config.anthropicSecret = ""
    }
    if(!config.llamafamilySecret) {
        config.llamafamilySecret = ""
    }
    if(!config.siliconflowSecret) {
        config.siliconflowSecret = ""
    }
    return config
};