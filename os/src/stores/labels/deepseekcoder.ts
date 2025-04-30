export const deepseekcoderLabels = {
    name: "deepseek-coder",
    family: "gemma",
    action: "code",
    models: [
        {
            model: "deepseek-coder:1.3b-instruct-q8_0",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": []
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "1.4GB",
                desk: "2GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
        {
            model: "deepseek-coder:1.3b-base-q8_0",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": []
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "1.4GB",
                desk: "2GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
        {
            model: "deepseek-coder:6.7b-instruct-q8_0",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": []
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "7.2GB",
                desk: "8GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
        {
            model: "deepseek-coder:6.7b-base-q8_0",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": []
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "7.2GB",
                desk: "8GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
        {
            model: "deepseek-coder-v2:16b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": []
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "8.9GB",
                desk: "9GB",
                cpu: "32GB",
                gpu: "12GB",
                quant: "f16"
            }
        },
    ],
    zhdesc: "DeepSeek Coder是一个基于两万亿代码和自然语言标记的强大编码模型。",
    endesc: "DeepSeek Coder is a capable coding model trained on two trillion code and natural language tokens. "
}