export const solarLabels = {
    name: "solar",
    family: "llama",
    action: "chat",
    models: [
        {
            model: "solar:10.7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "num_ctx": 4096,
                "stop": [
                    "</s>",
                    "### System:",
                    "### User:",
                    "### Assistant:"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "6.1GB",
                desk: "7GB",
                cpu: "32GB",
                gpu: "12GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "一个紧凑而强大的10.7B大型语言模型，专为单回合对话而设计。",
    endesc: "A compact, yet powerful 10.7B large language model designed for single-turn conversation. "
}