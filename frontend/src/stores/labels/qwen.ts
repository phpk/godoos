export const qwenLabels = {
    name: "qwen",
    family: "llama",
    action: "chat",
    zhdesc: "Qwen是阿里云基于transformer的一系列大型语言模型，在大量数据上进行预训练，包括网络文本、书籍、代码等。",
    endesc: "Qwen is a series of transformer-based large language models by Alibaba Cloud, pre-trained on a large volume of data, including web texts, books, code, etc.",
    models: [
        {
            model: "qwen2.5:0.5b",
            params: {
                top_p: 0.95,
                stream: true,
                num_keep: 5,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                stop: [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "494MB",
                desk: "1GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
        {
            model: "qwen2.5:1.5b",
            params: {
                top_p: 0.95,
                stream: true,
                num_keep: 5,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                stop: [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "1.54GB",
                desk: "1.6GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
        {
            model: "qwen2.5:3b",
            params: {
                top_p: 0.95,
                stream: true,
                num_keep: 5,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                stop: [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "1.9GB",
                desk: "2GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
        {
            model: "qwen2.5:7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_keep: 5,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                stop: [
                    "<|im_start|>",
                    "<|im_end|>"
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
            model: "qwen2:0.5b",
            params: {
                top_p: 0.95,
                stream: true,
                num_keep: 5,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                stop: [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "352MB",
                desk: "1GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
        {
            model: "qwen2:1.5b",
            params: {
                top_p: 0.95,
                stream: true,
                num_keep: 5,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                stop: [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "935MB",
                desk: "1.5GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
        {
            model: "qwen2:7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_keep: 5,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                stop: [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "4.4GB",
                desk: "6GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
        {
            model: "qwen:0.5b",
            params: {
                top_p: 0.95,
                stream: true,
                num_keep: 5,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                stop: [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "395MB",
                desk: "395MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
    ],
    
}