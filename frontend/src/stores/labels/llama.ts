export const llamaLabels = {
    name: "llama",
    family: "llama",
    action: "chat",
    zhdesc: "Llama由Meta Platforms发布,在通用基准测试上优于许多可用的开源聊天模型。",
    endesc: "Llama is released by Meta Platforms, Inc.Llama 3 instruction-tuned models are fine-tuned and optimized for dialogue/chat use cases and outperform many of the available open-source chat models on common benchmarks.",
    models: [
        {
            model: "llama3.2:1b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "num_keep": 24,
                "stop": [
                    "<|start_header_id|>",
                    "<|end_header_id|>",
                    "<|eot_id|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "1.3GB",
                desk: "2GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
        {
            model: "llama3.2:3b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "num_keep": 24,
                "stop": [
                    "<|start_header_id|>",
                    "<|end_header_id|>",
                    "<|eot_id|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "3.2GB",
                desk: "4GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
        {
            model: "llama3:8b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "num_keep": 24,
                "stop": [
                    "<|start_header_id|>",
                    "<|end_header_id|>",
                    "<|eot_id|>"
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
            model: "llama3-chatqa:8b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.3,
                "num_keep": 24,
                "stop": [
                    "<|start_header_id|>",
                    "<|end_header_id|>",
                    "<|eot_id|>"
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
            model: "ollam/unichat-llama3-chinese-8b:q4_0",

            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "num_keep": 24,
                "stop": [
                    "<|start_header_id|>",
                    "<|end_header_id|>",
                    "<|eot_id|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "4.7GB",
                desk: "5GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q8",
            }
        },
        {
            model: "llama2:7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "num_keep": 24,
                "stop": [
                    "[INST]",
                    "[/INST]",
                    "<<SYS>>",
                    "<</SYS>>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "3.8GB",
                desk: "4GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
        {
            model: "llama2:13b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "num_keep": 24,
                "stop": [
                    "[INST]",
                    "[/INST]",
                    "<<SYS>>",
                    "<</SYS>>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "7.4GB",
                desk: "8GB",
                cpu: "32GB",
                gpu: "12GB",
                quant: "q4"
            }
        },
        {
            model: "llama2-chinese:7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "num_keep": 24,
                "stop": [
                    "Name:",
                    "Assistant:"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "3.8GB",
                desk: "4GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
    ],

}