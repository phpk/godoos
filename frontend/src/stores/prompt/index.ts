import promptsZh from "./prompts-zh.json"
import promptsEn from "./prompts-en.json"
const promptAction = [
    "chat",
    "translation",
    "spoken",
    "creation_system",
    "creation_leader",//生成大纲
    "creation_builder",//根据主题和提纲进行撰写
    "creation_continuation",//续写
    "creation_optimization",//优化
    "creation_proofreading",//纠错
    "creation_summarize",//总结
    "creation_translation",//翻译
    "knowledge",
]

export { promptAction, promptsZh, promptsEn }
