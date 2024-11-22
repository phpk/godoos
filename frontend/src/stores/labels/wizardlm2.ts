export const wizardlm2Labels = {
    name: "wizardlm2",
    family: "llama",
    action: "chat",
    models: [
        {
            model: "wizardlm2:7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": [
                    "USER:",
                    "ASSISTANT:"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "4.1GB",
                desk: "5GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "微软人工智能的最先进的大型语言模型，在复杂聊天、多语言、推理和代理用例方面的性能有所提高。",
    endesc: "State of the art large language model from Microsoft AI with improved performance on complex chat, multilingual, reasoning and agent use cases. "
}