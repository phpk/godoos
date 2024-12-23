export const cateList: any = ["chat", "translation", "code", "img2txt", "image", "tts", "audio", "embeddings", "reranker"]
export const modelEngines = [
  {
    name: "ollama",
    cpp: "ollama",
    needQuant: true
  },
  // {
  //   name: "llama",
  //   cpp: "llama.cpp",
  //   needQuant: true
  // },
  // {
  //   name: "cortex",
  //   cpp: "cortex.cpp",
  //   needQuant: true
  // },
  // {
  //   name: "llamafile",
  //   cpp: "llamafile",
  //   needQuant: false
  // },
  {
    name: "sd",
    cpp: "stable-diffusion.cpp",
    needQuant: false
  },
  {
    name: "voice",
    cpp: "sherpa.cpp",
    needQuant: false
  }
]
export const netEngines = [

  {
    name: "OpenAI",
    cpp: "openai",
    needID: false,
  },
  {
    name: "GiteeAI",
    cpp: "gitee",
    needID: false,
  },
  {
    name: "CloudflareWorkersAI",
    cpp: "cloudflare",
    needID: true,
  },
  {
    name: "DeepSeek",
    cpp: "deepseek",
    needID: false,
  },
  {
    name: "智谱清言语BigModel",
    cpp: "bigmodel",
    needID: false,
  },
  {
    name: "火山方舟",
    cpp: "volces",
    needID: false,
  },
  {
    name: "阿里通义",
    cpp: "alibaba",
    needID: false,
  },
  {
    name: "Groq",
    cpp: "groq",
    needID: false,
  },
  {
    name: "Mistral",
    cpp: "mistral",
    needID: false,
  },
  {
    name: "Anthropic",
    cpp: "anthropic",
    needID: false,
  },
  {
    name: "llama.family",
    cpp: "llamafamily",
    needID: false,
  },
  {
    name: "硅基流动",
    cpp: "siliconflow",
    needID: false,
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