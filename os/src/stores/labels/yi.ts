export const yiLabels = {
    name: "yi",
    family: "yi",

    action: "chat",
    models: [
        {
            model: "yi:6b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "3.5GB",
                desk: "4GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
        {
            model: "yi:9b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "5GB",
                desk: "6GB",
                cpu: "32GB",
                gpu: "12GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "Yi是一系列大型语言模型，在3万亿个高质量的语料库上训练，支持英语和汉语。",
    endesc: "Yi is a series of large language models trained on a high-quality corpus of 3 trillion tokens that support both the English and Chinese languages."
}