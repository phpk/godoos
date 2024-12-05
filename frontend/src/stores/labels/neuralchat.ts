export const neuralchatLabels = {
    name: "neural-chat",
    family: "llama",
    action: "chat",
    models: [
        {
            model: "neural-chat:latest",
            params: {
                stream: true,
                "num_ctx": 4096,
                "stop": [
                    "</s>",
                    "<|im_start|>",
                    "<|im_end|>"
                ],            
                "temperature": 0.7,
                "top_k": 5,
                "top_p": 0.8
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
    zhdesc: "基于Mistral的微调模型，具有良好的领域和语言覆盖率。",
    endesc: "A fine-tuned model based on Mistral with good coverage of domain and language. "
}