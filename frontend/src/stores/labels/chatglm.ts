export const chatglmLabels = {
    name: "chatglm",
    family: "llama",
    action: "chat",
    models: [
        {
            model: "EntropyYue/chatglm3:6b",
            params: {
                stream: true,
                "stop": [
                    "<|system|>",
                    "<|user|>",
                    "<|assistant|>"            
                ],
                "temperature": 0.7,
                "top_k": 5,
                "top_p": 0.8
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "3.6GB",
                desk: "4GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "ChatGLM是由清华技术成果转化的公司智谱AI发布的开源的、支持中英双语问答的对话语言模型系列，并针对中文进行了优化，该模型基于General Language Model（GLM）架构构建",
    endesc: "ChatGLM is an open-source dialogue language model series released by Zhipu AI, a company that transforms technology achievements from Tsinghua University. It supports bilingual Q&A in both Chinese and English and has been optimized for Chinese. The model is built on the General Language Model (GLM) architecture"
}