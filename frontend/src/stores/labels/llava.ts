export const llavaLabels = {
    name: "llava",
    family: "llama",
    action: "img2txt",
    models: [
        {
            model: "llava:7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "num_keep": 24,
                "stop": [
                    "[INST]",
                    "[/INST]"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "4.7GB",
                desk: "5GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
        {
            model: "llava:13b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "num_keep": 24,
                "stop": [
                    "[INST]",
                    "[/INST]"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "8GB",
                desk: "9GB",
                cpu: "32GB",
                gpu: "12GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "LLaVA是一种新颖的端到端训练的大型多模式模型，它结合了视觉编码器和Vicuna，用于通用视觉和语言理解。",
    endesc: "LLaVA is a novel end-to-end trained large multimodal model that combines a vision encoder and Vicuna for general-purpose visual and language understanding. "
} 