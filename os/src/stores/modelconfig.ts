export const cateList: any = ["chat", "translation", "code", "img2txt", "image", "tts", "audio", "embeddings", "reranker"]

export const aiEngines = [
  {
    name: "ollama",
    cpp: "ollama",
    enable: false,
    url: "http://localhost:11434/api",
  },
  {
    name: "OpenAI",
    cpp: "openai",
    apikey: "",
    url: "https://api.openai.com/v1",
    enable: false,
  },
  {
    name: "GiteeAI",
    cpp: "gitee",
    apikey: "",
    url: "https://ai.gitee.com/v1",
    enable: false,
  },
  {
    name: "CloudflareWorkersAI",
    cpp: "cloudflare",
    apikey: "",
    url: "https://api.cloudflare.com/client/v4/accounts/{userId}/ai/v1",
    enable: true,
  },
  {
    name: "DeepSeek",
    cpp: "deepseek",
    apikey: "",
    url: "https://api.deepseek.com/v1",
    enable: false,
  },
  {
    name: "智谱清言语BigModel",
    cpp: "bigmodel",
    apikey: "",
    url: "https://open.bigmodel.cn/api/paas/v4",
    enable: false,
  },
  {
    name: "火山方舟",
    cpp: "volces",
    apikey: "",
    url: "https://ark.cn-beijing.volces.com/api/v3",
    enable: false,
  },
  {
    name: "阿里通义",
    cpp: "alibaba",
    apikey: "",
    url: "https://dashscope.aliyuncs.com/compatible-mode/v1",
    enable: false,
  },
  {
    name: "Groq",
    cpp: "groq",
    apikey: "",
    url: "https://api.groq.com/openai/v1",
    enable: false,
  },
  {
    name: "Mistral",
    cpp: "mistral",
    apikey: "",
    url: "https://api.mistral.ai/v1",
    enable: false,
  },
  {
    name: "Anthropic",
    cpp: "anthropic",
    apikey: "",
    url: "https://api.anthropic.com/v1",
    enable: false,
  },
  {
    name: "llama.family",
    cpp: "llamafamily",
    apikey: "",
    url: "https://api.atomecho.cn/v1",
    enable: false,
  },
  {
    name: "硅基流动",
    cpp: "siliconflow",
    apikey: "",
    url: "https://api.siliconflow.cn/v1",
    enable: false,
  },
]
export const llamaQuant = [
  "q2_K",
  "q3_K",
  "q3_K_S",
  "q3_K_M",
  "q3_K_L",
  "q4_0",
  "q4_1",
  "q4_K",
  "q4_K_S",
  "q4_K_M",
  "q5_0",
  "q5_1",
  "q5_K",
  "q5_K_S",
  "q5_K_M",
  "q6_K",
  "q8_0",
  "f16",
]
export const chatInitConfig = {
  chat: {
    key: "chat",
    contextLength: 10,
    num_keep: 5, //保留多少个最有可能的预测结果。这与top_k一起使用，决定模型在生成下一个词时考虑的词汇范围。
    num_predict: 3, //生成多少个预测结果
    top_p: 0.95,
    top_k: 40, //影响生成的随机性。较高的top_k值将使模型考虑更多的词汇
    temperature: 0.7, //影响生成的随机性。较低的温度产生更保守的输出，较高的温度产生更随机的输出。
  },
  translation: {
    key: "translation",
    num_keep: 5,
    num_predict: 1,
    top_k: 40,
    top_p: 0.95,
    temperature: 0.2,
  },
  creation: {
    key: "creation",
    num_keep: 3,
    num_predict: 1,
    top_k: 40,
    top_p: 0.95,
    temperature: 0.2,
  },
  knowledge: {
    key: "knowledge",
    contextLength: 10,
    num_keep: 5,
    num_predict: 1,
    top_k: 40,
    top_p: 0.95,
    temperature: 0.2,
  },
  spoken: {
    key: "spoken",
    contextLength: 10,
    num_keep: 5,
    num_predict: 1,
    top_k: 40,
    top_p: 0.95,
    temperature: 0.2,
  }
}