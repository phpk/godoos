export const codeqwenLabels = {
    name: "codeqwen",
    family: "gemma",
    action: "code",
    models: [
        {
            model: "qwen2.5-coder:3b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "penalize_newline": false,
                "repeat_penalty": 1,
                "stop": [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "1.9GB",
                desk: "2GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
        {
            model: "codeqwen:7b-chat",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "penalize_newline": false,
                "repeat_penalty": 1,
                "stop": [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "4.2GB",
                desk: "5GB",
                cpu: "12GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
        {
            model: "codeqwen:7b-code",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "penalize_newline": false,
                "repeat_penalty": 1,
                "stop": [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "4.2GB",
                desk: "5GB",
                cpu: "12GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "CodeQwen是一个在大量代码数据上预训练的大型语言模型。",
    endesc: "CodeQwen is a large language model pretrained on a large amount of code data. "
}